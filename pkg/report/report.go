package report

import (
	"fmt"

	"github.com/sonalgaud12/grater/internal/models"
)

// PrintTestSummary outputs a human-readable summary of test results
func PrintTestSummary(results *models.TestSuiteResults) {
	passed, failed := 0, 0
	for _, r := range results.Results {
		if r.Success {
			passed++
		} else {
			failed++
		}
	}

	total := len(results.Results)
	fmt.Printf("\n========== Test Summary ==========\n")
	fmt.Printf("Module: %s@%s\n", results.ModulePath, results.Version)
	fmt.Printf("Total Dependents: %d\n", total)
	fmt.Printf("Passed: %d (%.1f%%)\n", passed, float64(passed)/float64(total)*100)
	fmt.Printf("Failed: %d (%.1f%%)\n", failed, float64(failed)/float64(total)*100)
}

// PrintComparisonSummary outputs a human-readable comparison summary
func PrintComparisonSummary(report *models.ComparisonReport) {
	fmt.Printf("\n========== Comparison Summary ==========\n")
	fmt.Printf("Base Version: %s\n", report.BaseVersion)
	fmt.Printf("New Version:  %s\n", report.NewVersion)
	fmt.Printf("Total Dependents: %d\n", report.TotalDependents)
	fmt.Printf("\n")

	if len(report.NewFailures) > 0 {
		fmt.Printf("⚠️  NEW FAILURES (%d):\n", len(report.NewFailures))
		for _, m := range report.NewFailures {
			fmt.Printf("  - %s\n", m)
		}
		fmt.Printf("\n")
	}

	if len(report.NewPasses) > 0 {
		fmt.Printf("✓ NEW PASSES (%d):\n", len(report.NewPasses))
		for _, m := range report.NewPasses {
			fmt.Printf("  - %s\n", m)
		}
		fmt.Printf("\n")
	}

	fmt.Printf("Still Passing: %d\n", len(report.StillPassing))
	fmt.Printf("Still Failing: %d\n", len(report.StillFailing))

	if len(report.NewFailures) > 0 {
		fmt.Printf("\n⚠️  WARNING: This change introduced %d new test failures\n", len(report.NewFailures))
	} else {
		fmt.Printf("\n✓ No new test failures introduced\n")
	}
}
