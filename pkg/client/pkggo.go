package client

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/sonalgaud12/grater/internal/models"
)

// Client interface for fetching dependency information
type Client interface {
	FetchDependents(modulePath string, maxResults int) ([]models.Dependent, error)
}

// PkgGoDevClient implements Client using pkg.go.dev
type PkgGoDevClient struct{}

// NewPkgGoDevClient creates a new pkg.go.dev client
func NewPkgGoDevClient() *PkgGoDevClient {
	return &PkgGoDevClient{}
}

// FetchDependents retrieves modules that depend on the given module
// TODO: Implement actual pkg.go.dev API or deps.dev API integration
func (c *PkgGoDevClient) FetchDependents(modulePath string, maxResults int) ([]models.Dependent, error) {
	// Placeholder implementation using go mod
	cmd := exec.Command("go", "list", "-m", "-json", "all")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list modules: %w", err)
	}

	var dependents []models.Dependent
	decoder := json.NewDecoder(strings.NewReader(string(output)))

	for decoder.More() && len(dependents) < maxResults {
		var mod struct {
			Path    string `json:"Path"`
			Version string `json:"Version"`
		}
		if err := decoder.Decode(&mod); err != nil {
			continue
		}
		if mod.Path != "" && mod.Path != modulePath {
			dependents = append(dependents, models.Dependent{
				Path:    mod.Path,
				Version: mod.Version,
			})
		}
	}

	return dependents, nil
}
