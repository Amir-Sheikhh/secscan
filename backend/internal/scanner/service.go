package scanner

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/amir-sheikh/secscan/backend/internal/model"
	"github.com/amir-sheikh/secscan/backend/internal/storage"
)

type Config struct {
	ScanTimeout         time.Duration
	EnableActiveProbes  bool
	EnableExternalIntel bool
}

type Service struct {
	store   *storage.MemoryStore
	modules map[string]Module
	timeout time.Duration
}

type Module interface {
	Name() string
	Run(ctx context.Context, target *Target) model.ModuleResult
}

func NewService(store *storage.MemoryStore, cfg Config) *Service {
	service := &Service{
		store:   store,
		timeout: cfg.ScanTimeout,
		modules: map[string]Module{},
	}

	service.register(NewPortsModule())
	service.register(NewHeadersModule())
	service.register(NewTLSModule())
	service.register(NewFuzzModule())
	service.register(NewXSSModule())
	service.register(NewSQLIModule(cfg.EnableActiveProbes))
	service.register(NewCVEModule(cfg.EnableExternalIntel))

	return service
}

func (s *Service) Start(req model.ScanRequest) (*model.Scan, error) {
	target, err := PrepareTarget(req.URL)
	if err != nil {
		return nil, err
	}

	selected, err := s.normalizeModules(req.Modules)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	scan := &model.Scan{
		ID:         newScanID(),
		URL:        target.RawURL,
		Hostname:   target.Hostname,
		ResolvedIP: target.ResolvedIP.String(),
		Status:     model.ScanQueued,
		CreatedAt:  now,
		Modules:    map[string]*model.ModuleResult{},
		Events: []model.ScanEvent{{
			Type:    "scan",
			Status:  string(model.ScanQueued),
			Message: "scan queued",
			At:      now,
		}},
	}

	for _, name := range selected {
		scan.Modules[name] = &model.ModuleResult{
			Name:   name,
			Status: model.ModulePending,
			Score:  0,
			Details: map[string]any{
				"selected": true,
			},
		}
	}

	if err := s.store.Create(scan); err != nil {
		return nil, err
	}

	go s.run(scan.ID)
	return s.store.Get(scan.ID)
}

func (s *Service) Get(id string) (*model.Scan, error) {
	return s.store.Get(id)
}

func (s *Service) Subscribe(id string) (<-chan model.Scan, func(), error) {
	return s.store.Subscribe(id)
}

func (s *Service) normalizeModules(requested []string) ([]string, error) {
	if len(requested) == 0 {
		return append([]string(nil), model.DefaultModules...), nil
	}

	selected := make([]string, 0, len(requested))
	seen := map[string]struct{}{}
	for _, item := range requested {
		name := strings.ToLower(strings.TrimSpace(item))
		if name == "" {
			continue
		}
		if _, exists := s.modules[name]; !exists {
			return nil, fmt.Errorf("unknown module: %s", name)
		}
		if _, exists := seen[name]; exists {
			continue
		}
		seen[name] = struct{}{}
		selected = append(selected, name)
	}
	if len(selected) == 0 {
		return nil, errors.New("at least one valid module must be selected")
	}
	slices.Sort(selected)
	return selected, nil
}

func (s *Service) register(module Module) {
	s.modules[module.Name()] = module
}

func (s *Service) run(id string) {
	scan, err := s.store.Get(id)
	if err != nil {
		return
	}

	startedAt := time.Now().UTC()
	_, _ = s.store.Update(id, func(scan *model.Scan) error {
		scan.Status = model.ScanRunning
		scan.StartedAt = &startedAt
		scan.Events = append(scan.Events, model.ScanEvent{
			Type:    "scan",
			Status:  string(model.ScanRunning),
			Message: "scan started",
			At:      startedAt,
		})
		return nil
	})

	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	target, err := PrepareTarget(scan.URL)
	if err != nil {
		s.failScan(id, fmt.Errorf("target validation failed at runtime: %w", err))
		return
	}

	selected := make([]string, 0, len(scan.Modules))
	for name := range scan.Modules {
		selected = append(selected, name)
	}
	slices.Sort(selected)

	var wg sync.WaitGroup
	for _, name := range selected {
		module := s.modules[name]
		wg.Add(1)
		go func(name string, module Module) {
			defer wg.Done()
			moduleStart := time.Now().UTC()
			_, _ = s.store.Update(id, func(scan *model.Scan) error {
				result := scan.Modules[name]
				result.Status = model.ModuleRunning
				result.StartedAt = &moduleStart
				scan.Events = append(scan.Events, model.ScanEvent{
					Type:    "module",
					Module:  name,
					Status:  string(model.ModuleRunning),
					Message: name + " started",
					At:      moduleStart,
				})
				return nil
			})

			result := module.Run(ctx, target)
			moduleEnd := time.Now().UTC()
			if result.StartedAt == nil {
				result.StartedAt = &moduleStart
			}
			result.CompletedAt = &moduleEnd
			result.DurationMS = moduleEnd.Sub(*result.StartedAt).Milliseconds()

			_, _ = s.store.Update(id, func(scan *model.Scan) error {
				scan.Modules[name] = &result
				scan.Events = append(scan.Events, model.ScanEvent{
					Type:    "module",
					Module:  name,
					Status:  string(result.Status),
					Message: result.Summary,
					At:      moduleEnd,
				})
				return nil
			})
		}(name, module)
	}
	wg.Wait()

	finishedAt := time.Now().UTC()
	_, _ = s.store.Update(id, func(scan *model.Scan) error {
		scan.Status = model.ScanCompleted
		scan.CompletedAt = &finishedAt
		scan.Summary = buildSummary(scan)
		scan.Events = append(scan.Events, model.ScanEvent{
			Type:    "scan",
			Status:  string(model.ScanCompleted),
			Message: "scan completed",
			At:      finishedAt,
		})
		return nil
	})
}

func (s *Service) failScan(id string, reason error) {
	now := time.Now().UTC()
	_, _ = s.store.Update(id, func(scan *model.Scan) error {
		scan.Status = model.ScanFailed
		scan.CompletedAt = &now
		scan.Events = append(scan.Events, model.ScanEvent{
			Type:    "scan",
			Status:  string(model.ScanFailed),
			Message: reason.Error(),
			At:      now,
		})
		return nil
	})
}

func buildSummary(scan *model.Scan) model.ScanSummary {
	var totalScore int
	var count int
	summary := model.ScanSummary{}

	for _, module := range scan.Modules {
		totalScore += module.Score
		count++
		if module.Score >= 80 && module.Status == model.ModuleCompleted {
			summary.Passed++
		} else {
			summary.Failed++
		}
		for _, finding := range module.Findings {
			summary.Findings++
			switch strings.ToLower(finding.Severity) {
			case "critical":
				summary.Critical++
			case "high":
				summary.High++
			case "medium":
				summary.Medium++
			case "low":
				summary.Low++
			}
		}
	}

	if count > 0 {
		summary.Score = totalScore / count
		summary.ModuleRuns = count
	}
	if scan.StartedAt != nil && scan.CompletedAt != nil {
		summary.DurationMS = scan.CompletedAt.Sub(*scan.StartedAt).Milliseconds()
	}

	switch {
	case summary.Score >= 96:
		summary.Grade = "A+"
	case summary.Score >= 90:
		summary.Grade = "A"
	case summary.Score >= 80:
		summary.Grade = "B"
	case summary.Score >= 70:
		summary.Grade = "C"
	case summary.Score >= 60:
		summary.Grade = "D"
	default:
		summary.Grade = "F"
	}

	switch {
	case summary.Critical > 0:
		summary.RiskLevel = "critical"
	case summary.High > 1:
		summary.RiskLevel = "high"
	case summary.High == 1 || summary.Medium > 2:
		summary.RiskLevel = "medium"
	default:
		summary.RiskLevel = "low"
	}

	return summary
}

func newScanID() string {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err == nil {
		return hex.EncodeToString(buf)
	}
	return fmt.Sprintf("scan-%d", time.Now().UTC().UnixNano())
}
