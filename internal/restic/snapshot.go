package restic

import (
	"encoding/json"
	"fmt"
)

// ListSnapshots lists all snapshots
func (r *ResticWrapper) ListSnapshots() ([]Snapshot, error) {
	output, err := r.RunCommand([]string{"snapshots"})
	if err != nil {
		return nil, err
	}

	var snapshots []Snapshot
	if err := json.Unmarshal([]byte(output), &snapshots); err != nil {
		return nil, fmt.Errorf("failed to parse snapshots: %w", err)
	}

	return snapshots, nil
}

// ForgetSnapshot removes a snapshot
func (r *ResticWrapper) ForgetSnapshot(snapshotID string, prune bool) (string, error) {
	args := []string{"forget", snapshotID}
	if prune {
		args = append(args, "--prune")
	}
	return r.RunCommand(args)
}

// RestoreSnapshot restores a snapshot to a target directory
func (r *ResticWrapper) RestoreSnapshot(snapshotID, targetDir string) (string, error) {
	return r.RunCommand([]string{"restore", snapshotID, "--target", targetDir})
}

// ListSnapshotFiles lists files in a snapshot
func (r *ResticWrapper) ListSnapshotFiles(snapshotID string) (string, error) {
	// ls command output is not a single JSON array, but a stream of JSON objects.
	// For now, we might return the raw string and handle parsing in the controller or frontend if needed,
	// or we can wrap it in a JSON array.
	// Restic ls --json outputs one JSON object per line/file.

	output, err := r.RunCommand([]string{"ls", snapshotID})
	if err != nil {
		return "", err
	}
	return output, nil
}
