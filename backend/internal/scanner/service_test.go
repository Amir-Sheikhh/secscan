package scanner

import (
	"testing"
	"time"

	"github.com/amir-sheikh/secscan/backend/internal/model"
	"github.com/amir-sheikh/secscan/backend/internal/storage"
)

func TestBuildSummary(t *testing.T) {
	t.Parallel()

	start := time.Now().UTC()
	end := start.Add(3 * time.Second)
	scan := &model.Scan{
		StartedAt:   &start,
		CompletedAt: &end,
		Modules: map[string]*model.ModuleResult{
			"headers": {
				Name:   "headers",
				Status: model.ModuleCompleted,
				Score:  92,
			},
			"tls": {
				Name:   "tls",
				Status: model.ModuleCompleted,
				Score:  50,
				Findings: []model.Finding{{
					Severity: "high",
				}},
			},
		},
	}

	summary := buildSummary(scan)
	if summary.Score != 71 {
		t.Fatalf("unexpected score: %d", summary.Score)
	}
	if summary.Grade != "C" {
		t.Fatalf("unexpected grade: %s", summary.Grade)
	}
	if summary.High != 1 {
		t.Fatalf("expected 1 high finding, got %d", summary.High)
	}
}

func TestNormalizeModulesDefaults(t *testing.T) {
	t.Parallel()

	service := NewService(storage.NewMemoryStore(), Config{ScanTimeout: 5 * time.Second})
	modules, err := service.normalizeModules(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(modules) != len(model.DefaultModules) {
		t.Fatalf("expected %d modules, got %d", len(model.DefaultModules), len(modules))
	}
}
