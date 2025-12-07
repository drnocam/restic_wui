package controller

import (
	"fmt"
	"restic_wui/internal/restic"
)

type SnapshotController struct {
	repoCtrl *RepoController
}

func NewSnapshotController(repoCtrl *RepoController) *SnapshotController {
	return &SnapshotController{
		repoCtrl: repoCtrl,
	}
}

func (c *SnapshotController) ListSnapshots() ([]restic.Snapshot, error) {
	if !c.repoCtrl.IsOpen() {
		return nil, fmt.Errorf("no repository open")
	}
	return c.repoCtrl.Wrapper.ListSnapshots()
}

func (c *SnapshotController) ForgetSnapshot(id string, prune bool) error {
	if !c.repoCtrl.IsOpen() {
		return fmt.Errorf("no repository open")
	}
	_, err := c.repoCtrl.Wrapper.ForgetSnapshot(id, prune)
	return err
}

func (c *SnapshotController) RestoreSnapshot(id, targetDir string, paths []string, host string) (*restic.RestoreMessage, error) {
	if !c.repoCtrl.IsOpen() {
		return nil, fmt.Errorf("no repository open")
	}
	return c.repoCtrl.Wrapper.RestoreSnapshot(id, targetDir, paths, host)
}

func (c *SnapshotController) ListSnapshotFiles(id string) ([]restic.LSNode, error) {
	if !c.repoCtrl.IsOpen() {
		return nil, fmt.Errorf("no repository open")
	}
	return c.repoCtrl.Wrapper.ListSnapshotFiles(id)
}
