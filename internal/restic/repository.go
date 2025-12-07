package restic

import (
	"encoding/json"
	"fmt"
	"os/exec"
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
	output, err := r.RunCommand([]string{"stats"})
	if err != nil {
		return nil, err
	}

	// Try strict JSON unmarshal first
	var stats RepositoryStats
	if err := json.Unmarshal([]byte(output), &stats); err != nil {
		fmt.Printf("Error unmarshalling stats: %v. Output was: %s\n", err, output)
		return nil, fmt.Errorf("failed to parse stats: %w", err)
	}

	// Double check if TotalSize is 0, log it for debugging
	if stats.TotalSize == 0 {
		fmt.Printf("Warning: TotalSize is 0. Raw output: %s\n", output)
	}

	// Get actual disk size using raw-data mode
	rawOutput, err := r.RunCommand([]string{"stats", "--mode", "raw-data"})
	if err == nil {
		var rawStats struct {
			TotalSize uint64 `json:"total_size"`
		}
		if err := json.Unmarshal([]byte(rawOutput), &rawStats); err == nil {
			stats.DiskSize = rawStats.TotalSize
		}
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

// RebuildIndex rebuilds the index
func (r *ResticWrapper) RebuildIndex() (string, error) {
	return r.RunCommand([]string{"rebuild-index"})
}

// Find searches for a pattern in the repository and returns parsed results
func (r *ResticWrapper) Find(pattern string) ([]FindResult, error) {
	output, err := r.RunCommand([]string{"find", pattern})
	if err != nil {
		return nil, err
	}

	var results []FindResult
	if err := json.Unmarshal([]byte(output), &results); err != nil {
		fmt.Printf("Error unmarshalling find results: %v. Output was: %s\n", err, output)
		return nil, fmt.Errorf("failed to parse find results: %w", err)
	}

	return results, nil
}

// Mount mounts the repository to a specific path
// Note: This is a blocking command usually. We need to handle it carefully.
// For WUI, we might just start it and return success if it starts?
// But 'restic mount' blocks until unmounted.
// We probably want to run this in a way that we can kill it later.
// For now, let's implement a simple version that might block or fail if not handled by caller async.
// Actually, RunCommand waits. We need a way to run background command.
// But the wrapper structure is simple.
// Let's defer complexity of background mount to controller or a new wrapper method later if needed.
// For now, let's implement it but warn it might block.
// Wait, user wants "Mount" feature.
// If I run 'restic mount' it blocks.
// I'll implement it in wrapper, but Controller will need to handle async.
// Mount mounts the repository to a specific path using a background process.
// Returns the command object which can be used to wait for completion or kill the process.
func (r *ResticWrapper) Mount(mountPoint string) (*exec.Cmd, error) {
	if mountPoint == "" {
		return nil, fmt.Errorf("mount point cannot be empty")
	}
	// 'mount' command generally runs until interrupted.
	return r.RunCommandBackground([]string{"mount", mountPoint})
}
