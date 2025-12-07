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
	Command  string
}

// NewResticWrapper creates a new wrapper instance
func NewResticWrapper(repoPath, password, command string) *ResticWrapper {
	if command == "" {
		command = "restic"
	}
	return &ResticWrapper{
		RepoPath: repoPath,
		Password: password,
		Command:  command,
	}
}

// RunCommand executes a restic command with the configured repo and password
func (r *ResticWrapper) RunCommand(args []string) (string, error) {
	// Base arguments
	cmdArgs := []string{"-r", r.RepoPath, "--json"}
	cmdArgs = append(cmdArgs, args...)

	cmd := exec.Command(r.Command, cmdArgs...)

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

	cmd := exec.Command(r.Command, cmdArgs...)

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

// Backup runs the backup command
func (r *ResticWrapper) Backup(source string, excludes []string) (string, error) {
	args := []string{"backup", source}
	for _, exclude := range excludes {
		args = append(args, "--exclude", exclude)
	}
	return r.RunCommand(args)
}

// RunCommandBackground executes a restic command in the background and returns the command object
func (r *ResticWrapper) RunCommandBackground(args []string) (*exec.Cmd, error) {
	// Base arguments
	cmdArgs := []string{"-r", r.RepoPath}
	cmdArgs = append(cmdArgs, args...)

	cmd := exec.Command(r.Command, cmdArgs...)

	// Set environment variables
	env := os.Environ()
	env = append(env, "RESTIC_PASSWORD="+r.Password)
	cmd.Env = env

	// Start the command
	err := cmd.Start()
	if err != nil {
		return nil, fmt.Errorf("failed to start restic command: %w", err)
	}

	return cmd, nil
}
