package models

import "time"

// Dependent represents a module that depends on the target module
type Dependent struct {
	Path    string `json:"path"`
	Version string `json:"version"`
}

// TestResult represents the outcome of running tests for a dependent
type TestResult struct {
	Module    string        `json:"module"`
	Success   bool          `json:"success"`
	Output    string        `json:"output"`
	Duration  time.Duration `json:"duration"`
	TestCount int           `json:"test_count"`
	FailCount int           `json:"fail_count"`
	Error     string        `json:"error,omitempty"`
}

// TestSuiteResults contains all test results for a module version
type TestSuiteResults struct {
	ModulePath string                 `json:"module_path"`
	Version    string                 `json:"version"`
	Timestamp  time.Time              `json:"timestamp"`
	Results    map[string]*TestResult `json:"results"`
}

// ComparisonReport contains statistical summary of test differences
type ComparisonReport struct {
	BaseVersion     string                       `json:"base_version"`
	NewVersion      string                       `json:"new_version"`
	TotalDependents int                          `json:"total_dependents"`
	NewFailures     []string                     `json:"new_failures"`
	NewPasses       []string                     `json:"new_passes"`
	StillFailing    []string                     `json:"still_failing"`
	StillPassing    []string                     `json:"still_passing"`
	Details         map[string]*ComparisonDetail `json:"details"`
}

// ComparisonDetail contains detailed comparison for a single dependent
type ComparisonDetail struct {
	Module      string `json:"module"`
	BaseSuccess bool   `json:"base_success"`
	NewSuccess  bool   `json:"new_success"`
	BaseError   string `json:"base_error,omitempty"`
	NewError    string `json:"new_error,omitempty"`
}
