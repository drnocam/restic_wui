package restic

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

// ResticWrapper handles interactions with the restic CLI
type ResticWrapper struct {
	Password string
	RepoPath string
}

// NewResticWrapper creates a new wrapper instance
func NewResticWrapper(repoPath, password string) *ResticWrapper {
	return &ResticWrapper{
		RepoPath: repoPath,
		Password: password,
	}
}

// RunCommand executes a restic command with the configured repo and password
func (r *ResticWrapper) RunCommand(args []string) (string, error) {
	// Base arguments
	cmdArgs := []string{"-r", r.RepoPath, "--json"}
	cmdArgs = append(cmdArgs, args...)

	cmd := exec.Command("restic", cmdArgs...)

	// Set environment variables
	env := os.Environ()
	env = append(env, "RESTIC_PASSWORD="+r.Password)
	cmd.Env = env

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("restic command failed: %s, stderr: %s", err, stderr.String())
	}

	return out.String(), nil
}

// RunRawCommand executes a restic command without forcing JSON output (for some commands)
func (r *ResticWrapper) RunRawCommand(args []string) (string, error) {
	// Base arguments
	cmdArgs := []string{"-r", r.RepoPath}
	cmdArgs = append(cmdArgs, args...)

	cmd := exec.Command("restic", cmdArgs...)

	// Set environment variables
	env := os.Environ()
	env = append(env, "RESTIC_PASSWORD="+r.Password)
	cmd.Env = env

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("restic command failed: %s, stderr: %s", err, stderr.String())
	}

	return out.String(), nil
}
