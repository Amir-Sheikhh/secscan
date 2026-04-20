package model

import "time"

type ScanState string

const (
	ScanQueued    ScanState = "queued"
	ScanRunning   ScanState = "running"
	ScanCompleted ScanState = "completed"
	ScanFailed    ScanState = "failed"
)

type ModuleState string

const (
	ModulePending   ModuleState = "pending"
	ModuleRunning   ModuleState = "running"
	ModuleCompleted ModuleState = "completed"
	ModuleFailed    ModuleState = "failed"
	ModuleSkipped   ModuleState = "skipped"
)

var DefaultModules = []string{"ports", "headers", "tls", "fuzz", "xss", "sqli", "cve"}

type ScanRequest struct {
	URL     string   `json:"url" binding:"required"`
	Modules []string `json:"modules"`
}

type Scan struct {
	ID          string                   `json:"id"`
	URL         string                   `json:"url"`
	Hostname    string                   `json:"hostname"`
	ResolvedIP  string                   `json:"resolvedIp"`
	Status      ScanState                `json:"status"`
	CreatedAt   time.Time                `json:"createdAt"`
	StartedAt   *time.Time               `json:"startedAt,omitempty"`
	CompletedAt *time.Time               `json:"completedAt,omitempty"`
	Modules     map[string]*ModuleResult `json:"modules"`
	Events      []ScanEvent              `json:"events"`
	Summary     ScanSummary              `json:"summary"`
}

type ModuleResult struct {
	Name        string         `json:"name"`
	Status      ModuleState    `json:"status"`
	StartedAt   *time.Time     `json:"startedAt,omitempty"`
	CompletedAt *time.Time     `json:"completedAt,omitempty"`
	DurationMS  int64          `json:"durationMs"`
	Score       int            `json:"score"`
	Severity    string         `json:"severity"`
	Summary     string         `json:"summary"`
	Findings    []Finding      `json:"findings"`
	Details     map[string]any `json:"details,omitempty"`
	Error       string         `json:"error,omitempty"`
}

type Finding struct {
	Title          string `json:"title"`
	Severity       string `json:"severity"`
	Category       string `json:"category"`
	Description    string `json:"description"`
	Recommendation string `json:"recommendation"`
	Evidence       string `json:"evidence,omitempty"`
}

type ScanEvent struct {
	Type    string    `json:"type"`
	Module  string    `json:"module,omitempty"`
	Status  string    `json:"status"`
	Message string    `json:"message"`
	At      time.Time `json:"at"`
}

type ScanSummary struct {
	Score      int    `json:"score"`
	Grade      string `json:"grade"`
	RiskLevel  string `json:"riskLevel"`
	Findings   int    `json:"findings"`
	Critical   int    `json:"critical"`
	High       int    `json:"high"`
	Medium     int    `json:"medium"`
	Low        int    `json:"low"`
	Passed     int    `json:"passed"`
	Failed     int    `json:"failed"`
	ModuleRuns int    `json:"moduleRuns"`
	DurationMS int64  `json:"durationMs"`
}
