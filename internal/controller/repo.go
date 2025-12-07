package controller

import (
	"fmt"
	"os"
	"os/exec"
	"restic_wui/internal/config"
	"restic_wui/internal/restic"
)

type RepoController struct {
	Wrapper  *restic.ResticWrapper
	Config   *config.Config
	MountCmd *exec.Cmd
}

func NewRepoController(cfg *config.Config) *RepoController {
	return &RepoController{
		Config: cfg,
	}
}

// OpenRepository opens an existing repository
func (c *RepoController) OpenRepository(path, password string) error {
	// Verify it's a valid repo by checking config or running a simple command
	wrapper := restic.NewResticWrapper(path, password, c.Config.ResticCommand)
	_, err := wrapper.CheckRepository()
	if err != nil {
		return fmt.Errorf("failed to open repository: %w", err)
	}
	c.Wrapper = wrapper
	return nil
}

// InitRepository initializes a new repository
func (c *RepoController) InitRepository(path, password string) error {
	wrapper := restic.NewResticWrapper(path, password, c.Config.ResticCommand)
	_, err := wrapper.InitRepository()
	if err != nil {
		return fmt.Errorf("failed to init repository: %w", err)
	}
	c.Wrapper = wrapper
	return nil
}

// GetStats returns repository statistics
func (c *RepoController) GetStats() (*restic.RepositoryStats, error) {
	if c.Wrapper == nil {
		return nil, fmt.Errorf("no repository open")
	}
	return c.Wrapper.GetRepositoryStats()
}

// IsOpen checks if a repository is currently open
func (c *RepoController) IsOpen() bool {
	return c.Wrapper != nil
}

// GetCurrentPath returns the path of the open repository
func (c *RepoController) GetCurrentPath() string {
	if c.Wrapper == nil {
		return ""
	}
	return c.Wrapper.RepoPath
}

// CloseRepository closes the current repository
func (c *RepoController) CloseRepository() {
	if c.IsMounted() {
		c.UnmountRepository()
	}
	c.Wrapper = nil
}

// DeleteRepository deletes the current repository directory
func (c *RepoController) DeleteRepository() error {
	if c.Wrapper == nil {
		return fmt.Errorf("no repository open")
	}
	path := c.Wrapper.RepoPath
	c.CloseRepository()
	c.CloseRepository()
	return os.RemoveAll(path)
}

// Backup runs the backup command
func (c *RepoController) Backup(source string, excludes []string) (string, error) {
	if c.Wrapper == nil {
		return "", fmt.Errorf("no repository open")
	}
	return c.Wrapper.Backup(source, excludes)
}

// RepairRepository runs rebuild-index
func (c *RepoController) RepairRepository() (string, error) {
	if c.Wrapper == nil {
		return "", fmt.Errorf("no repository open")
	}
	return c.Wrapper.RebuildIndex()
}

// CheckRepository runs check
func (c *RepoController) CheckRepository() (string, error) {
	if c.Wrapper == nil {
		return "", fmt.Errorf("no repository open")
	}
	return c.Wrapper.CheckRepository()
}

// UnlockRepository runs unlock
func (c *RepoController) UnlockRepository() error {
	if c.Wrapper == nil {
		return fmt.Errorf("no repository open")
	}
	_, err := c.Wrapper.UnlockRepository()
	return err
}

// Find searches for a pattern
func (c *RepoController) Find(pattern string) ([]restic.FindResult, error) {
	if c.Wrapper == nil {
		return nil, fmt.Errorf("no repository open")
	}
	return c.Wrapper.Find(pattern)
}

// MountRepository mounts the repo
func (c *RepoController) MountRepository(mountPoint string) error {
	if c.Wrapper == nil {
		return fmt.Errorf("no repository open")
	}
	if c.MountCmd != nil && c.MountCmd.Process != nil {
		// Check if process is still running
		if c.MountCmd.ProcessState == nil {
			return fmt.Errorf("repository already mounted")
		}
		// Process has exited, clean up
		c.MountCmd = nil
	}

	// Validate mount point exists
	info, err := os.Stat(mountPoint)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("mount point does not exist: %s", mountPoint)
		}
		return fmt.Errorf("cannot access mount point: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("mount point is not a directory: %s", mountPoint)
	}

	cmd, err := c.Wrapper.Mount(mountPoint)
	if err != nil {
		return err
	}
	c.MountCmd = cmd
	return nil
}

// UnmountRepository stops the mount process
func (c *RepoController) UnmountRepository() error {
	if c.MountCmd != nil {
		if c.MountCmd.Process != nil {
			// Send interrupt signal
			c.MountCmd.Process.Signal(os.Interrupt)
			// Wait for process to exit
			go c.MountCmd.Wait()
		}
		c.MountCmd = nil
	}
	return nil
}

// IsMounted checks if mounted (verifies process is still running)
func (c *RepoController) IsMounted() bool {
	if c.MountCmd == nil {
		return false
	}
	if c.MountCmd.Process == nil {
		c.MountCmd = nil
		return false
	}
	// Check if process has exited
	if c.MountCmd.ProcessState != nil {
		c.MountCmd = nil
		return false
	}
	return true
}
