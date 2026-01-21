package gather

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/sonalgaud12/grater/internal/models"
	"github.com/sonalgaud12/grater/pkg/client"
)

// Options contains configuration for gathering dependents
type Options struct {
	ModulePath string
	MaxResults int
	OutputFile string
	Client     client.Client
}

// Service handles gathering of dependents
type Service struct {
	client client.Client
}

// NewService creates a new gather service
func NewService(c client.Client) *Service {
	return &Service{client: c}
}

// Gather retrieves and saves the list of dependents
func (s *Service) Gather(opts Options) ([]models.Dependent, error) {
	fmt.Printf("Gathering dependents for: %s\n", opts.ModulePath)

	dependents, err := s.client.FetchDependents(opts.ModulePath, opts.MaxResults)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch dependents: %w", err)
	}

	fmt.Printf("Found %d dependents\n", len(dependents))

	if opts.OutputFile != "" {
		if err := s.saveDependents(dependents, opts.OutputFile); err != nil {
			return nil, err
		}
		fmt.Printf("Dependents list written to: %s\n", opts.OutputFile)
	}

	return dependents, nil
}

func (s *Service) saveDependents(dependents []models.Dependent, filename string) error {
	data, err := json.MarshalIndent(dependents, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal dependents: %w", err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}

	return nil
}

// LoadDependents loads dependents from a file
func LoadDependents(filename string) ([]models.Dependent, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read dependents file: %w", err)
	}

	var dependents []models.Dependent
	if err := json.Unmarshal(data, &dependents); err != nil {
		return nil, fmt.Errorf("failed to parse dependents: %w", err)
	}

	return dependents, nil
}
