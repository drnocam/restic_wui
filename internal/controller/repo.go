package controller

import (
	"fmt"
	"restic_wui/internal/restic"
)

type RepoController struct {
	Wrapper *restic.ResticWrapper
}

func NewRepoController() *RepoController {
	return &RepoController{}
}

// OpenRepository opens an existing repository
func (c *RepoController) OpenRepository(path, password string) error {
	// Verify it's a valid repo by checking config or running a simple command
	wrapper := restic.NewResticWrapper(path, password)
	_, err := wrapper.CheckRepository()
	if err != nil {
		return fmt.Errorf("failed to open repository: %w", err)
	}
	c.Wrapper = wrapper
	return nil
}

// InitRepository initializes a new repository
func (c *RepoController) InitRepository(path, password string) error {
	wrapper := restic.NewResticWrapper(path, password)
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
