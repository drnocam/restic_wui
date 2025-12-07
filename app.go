package main

import (
	"context"
	"fmt"
	"restic_wui/internal/controller"
	"restic_wui/internal/restic"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx                context.Context
	repoController     *controller.RepoController
	snapshotController *controller.SnapshotController
}

// NewApp creates a new App application struct
func NewApp() *App {
	repoCtrl := controller.NewRepoController()
	snapshotCtrl := controller.NewSnapshotController(repoCtrl)
	return &App{
		repoController:     repoCtrl,
		snapshotController: snapshotCtrl,
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Repository Methods

func (a *App) OpenRepository(path, password string) error {
	return a.repoController.OpenRepository(path, password)
}

func (a *App) InitRepository(path, password string) error {
	return a.repoController.InitRepository(path, password)
}

func (a *App) GetRepositoryStats() (*restic.RepositoryStats, error) {
	return a.repoController.GetStats()
}

func (a *App) IsRepositoryOpen() bool {
	return a.repoController.IsOpen()
}

func (a *App) GetCurrentRepositoryPath() string {
	return a.repoController.GetCurrentPath()
}

func (a *App) PruneRepository() error {
	if !a.repoController.IsOpen() {
		return fmt.Errorf("no repository open")
	}
	_, err := a.repoController.Wrapper.PruneRepository()
	return err
}

// Snapshot Methods

func (a *App) ListSnapshots() ([]restic.Snapshot, error) {
	return a.snapshotController.ListSnapshots()
}

func (a *App) ForgetSnapshot(id string, prune bool) error {
	return a.snapshotController.ForgetSnapshot(id, prune)
}

func (a *App) RestoreSnapshot(id, targetDir string) error {
	return a.snapshotController.RestoreSnapshot(id, targetDir)
}

func (a *App) ListSnapshotFiles(id string) (string, error) {
	return a.snapshotController.ListSnapshotFiles(id)
}

// Helper Methods

func (a *App) SelectDirectory() (string, error) {
	selection, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Directory",
	})
	if err != nil {
		return "", fmt.Errorf("failed to open directory dialog: %w", err)
	}
	return selection, nil
}
