package dac

import (
	"time"
)

type Engine struct {
	SourceIndex map[string]EvaluatorService
}

func NewEngine() *Engine {
	evalMap := make(map[string]EvaluatorService)
	return &Engine{
		SourceIndex: evalMap,
	}
}

func (e *Engine) AddEvaluator(name string, ev EvaluatorService) {
	e.SourceIndex[name] = ev
}

func (e Engine) Detect(t Telemetry) ([]DetectionResult, bool) {
	evaluator, ok := e.SourceIndex[t.TelemetryType]
	if !ok {
		return nil, false
	}
	return evaluator.GetMatches(t.Data), true
}

type Telemetry struct {
	TelemetryType string
	Data          any
}

type EvaluatorService interface {
	GetMatches(any) []DetectionResult
}

type DetectionResult struct {
	RuleId         string          `json:"rule_id"`
	IsTuning       bool            `json:"is_tuning"`
	Timestamp      time.Time       `json:"timestamp" yaml:"timestamp"`
	Title          string          `json:"title" yaml:"title"`
	Description    string          `json:"description" yaml:"description"`
	Classification Classification  `json:"classification" yaml:"classification"`
	IsMatch        bool            `json:"is_match" yaml:"is_match"`
	Notes          []DetectionNote `json:"notes" yaml:"notes"`
}

type DetectionNote struct {
	Title       string `json:"title" yaml:"title"`
	Description string `json:"description" yaml:"description"`
}

type Classification string

const (
	INFORMATIONAL Classification = "informational"
	SUSPICIOUS    Classification = "suspicious"
	MALICIOUS     Classification = "malicious"
)

type RuleMeta[T any] struct {
	SourceType            string         `yaml:"source_type"`
	InternalId            string         `yaml:"internal_id"`
	DisplayName           string         `yaml:"display_name"`
	CreatedAt             time.Time      `yaml:"created_at"`
	UpdatedAt             time.Time      `yaml:"updated_at"`
	IsActive              bool           `yaml:"is_active"`
	IsTuning              bool           `yaml:"is_tuning"`
	DefaultDescription    string         `yaml:"default_description"`
	DefaultClassification Classification `yaml:"default_classification"`
	AnalysisSteps         string         `yaml:"analysis_steps"`
	FalsePositives        string         `yaml:"false_positives"`
	References            []string       `yaml:"references"`
	TestBaseFile          string         `yaml:"test_base"`
	Tests                 []TestCase[T]  `yaml:"test"`
}

type TestCase[T any] struct {
	It       string          `yaml:"it"`
	Want     DetectionResult `yaml:"want"`
	Override T               `yaml:"override_data"`
}
