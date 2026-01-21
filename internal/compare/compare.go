package compare

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/sonalgaud12/grater/internal/models"
)

// Service handles comparison of test results
type Service struct{}

// NewService creates a new comparison service
func NewService() *Service {
	return &Service{}
}

// Compare generates a comparison report between two test runs
func (s *Service) Compare(base, new *models.TestSuiteResults) *models.ComparisonReport {
	report := &models.ComparisonReport{
		BaseVersion:     base.Version,
		NewVersion:      new.Version,
		TotalDependents: len(base.Results),
		Details:         make(map[string]*models.ComparisonDetail),
	}

	for module := range base.Results {
		baseResult := base.Results[module]
		newResult := new.Results[module]

		detail := &models.ComparisonDetail{
			Module:      module,
			BaseSuccess: baseResult.Success,
			BaseError:   baseResult.Error,
		}

		if newResult != nil {
			detail.NewSuccess = newResult.Success
			detail.NewError = newResult.Error

			s.categorizeChange(report, module, baseResult, newResult)
		}

		report.Details[module] = detail
	}

	return report
}

func (s *Service) categorizeChange(report *models.ComparisonReport, module string, base, new *models.TestResult) {
	if !base.Success && new.Success {
		report.NewPasses = append(report.NewPasses, module)
	} else if base.Success && !new.Success {
		report.NewFailures = append(report.NewFailures, module)
	} else if !base.Success && !new.Success {
		report.StillFailing = append(report.StillFailing, module)
	} else {
		report.StillPassing = append(report.StillPassing, module)
	}
}

// SaveReport saves a comparison report to a file
func (s *Service) SaveReport(report *models.ComparisonReport, filename string) error {
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal report: %w", err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write report: %w", err)
	}

	return nil
}
