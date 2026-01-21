package test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/sonalgaud12/grater/internal/models"
)

// Options contains configuration for test execution
type Options struct {
	ModulePath string
	Version    string
	Dependents []models.Dependent
	WorkDir    string
	Timeout    time.Duration
	OutputFile string
}

// Runner executes tests for dependents
type Runner struct {
	workDir string
	timeout time.Duration
}

// NewRunner creates a new test runner
func NewRunner(workDir string, timeout time.Duration) *Runner {
	return &Runner{
		workDir: workDir,
		timeout: timeout,
	}
}

// RunAll executes tests for all dependents
func (r *Runner) RunAll(opts Options) (*models.TestSuiteResults, error) {
	fmt.Printf("Testing dependents of %s@%s\n", opts.ModulePath, opts.Version)

	if err := os.MkdirAll(opts.WorkDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create work directory: %w", err)
	}

	results := &models.TestSuiteResults{
		ModulePath: opts.ModulePath,
		Version:    opts.Version,
		Timestamp:  time.Now(),
		Results:    make(map[string]*models.TestResult),
	}

	for i, dep := range opts.Dependents {
		fmt.Printf("[%d/%d] Testing %s...\n", i+1, len(opts.Dependents), dep.Path)
		result := r.runDependentTest(dep, opts.ModulePath, opts.Version)
		results.Results[dep.Path] = result

		status := "PASS"
		if !result.Success {
			status = "FAIL"
		}
		fmt.Printf("  %s (%.2fs)\n", status, result.Duration.Seconds())
	}

	if opts.OutputFile != "" {
		if err := r.saveResults(results, opts.OutputFile); err != nil {
			return nil, err
		}
		fmt.Printf("\nResults written to: %s\n", opts.OutputFile)
	}

	return results, nil
}

func (r *Runner) runDependentTest(dep models.Dependent, targetModule, targetVersion string) *models.TestResult {
	start := time.Now()
	result := &models.TestResult{
		Module: dep.Path,
	}

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	// Get the module
	cmd := exec.CommandContext(ctx, "go", "get", "-t", dep.Path+"@"+dep.Version)
	if output, err := cmd.CombinedOutput(); err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("failed to get module: %v", err)
		result.Output = string(output)
		result.Duration = time.Since(start)
		return result
	}

	// Run tests
	cmd = exec.CommandContext(ctx, "go", "test", "-v", dep.Path+"/...")
	output, err := cmd.CombinedOutput()
	result.Output = string(output)
	result.Duration = time.Since(start)

	if err != nil {
		result.Success = false
		result.Error = err.Error()
	} else {
		result.Success = true
	}

	// Parse test output
	result.TestCount = strings.Count(result.Output, "=== RUN")
	result.FailCount = strings.Count(result.Output, "--- FAIL:")

	return result
}

func (r *Runner) saveResults(results *models.TestSuiteResults, filename string) error {
	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal results: %w", err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write results: %w", err)
	}

	return nil
}

// LoadResults loads test results from a file
func LoadResults(filename string) (*models.TestSuiteResults, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read results file: %w", err)
	}

	var results models.TestSuiteResults
	if err := json.Unmarshal(data, &results); err != nil {
		return nil, fmt.Errorf("failed to parse results: %w", err)
	}

	return &results, nil
}
