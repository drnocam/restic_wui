package restic

import (
	"encoding/json"
	"fmt"
)

// InitRepository initializes a new repository
func (r *ResticWrapper) InitRepository() (string, error) {
	return r.RunCommand([]string{"init"})
}

// CheckRepository checks the repository for errors
func (r *ResticWrapper) CheckRepository() (string, error) {
	return r.RunCommand([]string{"check"})
}

// GetRepositoryStats returns statistics about the repository
func (r *ResticWrapper) GetRepositoryStats() (*RepositoryStats, error) {
	output, err := r.RunCommand([]string{"stats", "--mode", "raw-data"})
	if err != nil {
		return nil, err
	}

	var stats RepositoryStats
	if err := json.Unmarshal([]byte(output), &stats); err != nil {
		return nil, fmt.Errorf("failed to parse stats: %w", err)
	}

	return &stats, nil
}

// UnlockRepository removes locks
func (r *ResticWrapper) UnlockRepository() (string, error) {
	return r.RunCommand([]string{"unlock"})
}

// PruneRepository removes unneeded data
func (r *ResticWrapper) PruneRepository() (string, error) {
	return r.RunCommand([]string{"prune"})
}
