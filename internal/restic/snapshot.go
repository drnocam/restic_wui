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
// Returns the restore summary with files/bytes restored
func (r *ResticWrapper) RestoreSnapshot(snapshotID, targetDir string, includes []string, host string) (*RestoreMessage, error) {
	args := []string{"restore", snapshotID, "--target", targetDir}
	for _, include := range includes {
		if include != "" {
			args = append(args, "--include", include)
		}
	}
	if host != "" {
		args = append(args, "--host", host)
	}

	// Use RunCommand to get JSON output
	output, err := r.RunCommand(args)
	if err != nil {
		return nil, err
	}

	// Parse JSON lines output, looking for the summary message
	var summary *RestoreMessage
	lines := splitLines(output)
	for _, line := range lines {
		if line == "" {
			continue
		}
		var msg RestoreMessage
		if err := json.Unmarshal([]byte(line), &msg); err != nil {
			continue
		}
		// Keep the last summary message
		if msg.MessageType == "summary" {
			summary = &msg
		}
	}

	if summary == nil {
		// No summary found, but command succeeded - return empty summary
		summary = &RestoreMessage{MessageType: "summary"}
	}

	return summary, nil
}

// ListSnapshotFiles lists files in a snapshot and returns parsed nodes
func (r *ResticWrapper) ListSnapshotFiles(snapshotID string) ([]LSNode, error) {
	output, err := r.RunCommand([]string{"ls", snapshotID})
	if err != nil {
		return nil, err
	}

	// ls --json outputs one JSON object per line (JSON lines format)
	// First line is snapshot info, subsequent lines are nodes
	var nodes []LSNode
	lines := splitLines(output)
	for i, line := range lines {
		if line == "" {
			continue
		}
		// Skip first line (snapshot info) - check message_type
		var msgType struct {
			MessageType string `json:"message_type"`
		}
		if err := json.Unmarshal([]byte(line), &msgType); err == nil {
			if msgType.MessageType == "snapshot" {
				continue
			}
		}

		// Try to parse as node
		var node LSNode
		if err := json.Unmarshal([]byte(line), &node); err != nil {
			// Skip unparseable lines
			fmt.Printf("Error parsing line %d: %v\n", i, err)
			continue
		}
		if node.Path != "" {
			nodes = append(nodes, node)
		}
	}

	return nodes, nil
}

// splitLines splits a string by newlines
func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}
