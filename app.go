package main

import (
	"context"
	"fmt"
	"restic_wui/internal/config"
	"restic_wui/internal/controller"
	"restic_wui/internal/restic"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx                context.Context
	repoController     *controller.RepoController
	snapshotController *controller.SnapshotController
	config             *config.Config
}

// NewApp creates a new App application struct
func NewApp() *App {
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		// Fallback to default config if load fails
		cfg = &config.Config{ResticCommand: "restic"}
	}

	repoCtrl := controller.NewRepoController(cfg)
	snapshotCtrl := controller.NewSnapshotController(repoCtrl)
	return &App{
		repoController:     repoCtrl,
		snapshotController: snapshotCtrl,
		config:             cfg,
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

func (a *App) CloseRepository() {
	a.repoController.CloseRepository()
}

func (a *App) DeleteRepository() error {
	return a.repoController.DeleteRepository()
}

func (a *App) Backup(source string, excludes []string) (string, error) {
	return a.repoController.Backup(source, excludes)
}

func (a *App) RepairRepository() (string, error) {
	return a.repoController.RepairRepository()
}

func (a *App) CheckRepository() (string, error) {
	return a.repoController.CheckRepository()
}

func (a *App) UnlockRepository() error {
	return a.repoController.UnlockRepository()
}

func (a *App) Find(pattern string) ([]restic.FindResult, error) {
	return a.repoController.Find(pattern)
}

func (a *App) MountRepository(mountPoint string) error {
	return a.repoController.MountRepository(mountPoint)
}

func (a *App) UnmountRepository() error {
	return a.repoController.UnmountRepository()
}

func (a *App) IsMounted() bool {
	return a.repoController.IsMounted()
}

// Snapshot Methods

func (a *App) ListSnapshots() ([]restic.Snapshot, error) {
	return a.snapshotController.ListSnapshots()
}

func (a *App) ForgetSnapshot(id string, prune bool) error {
	return a.snapshotController.ForgetSnapshot(id, prune)
}

func (a *App) RestoreSnapshot(id, targetDir string, paths []string, host string) (*restic.RestoreMessage, error) {
	return a.snapshotController.RestoreSnapshot(id, targetDir, paths, host)
}

func (a *App) ListSnapshotFiles(id string) ([]restic.LSNode, error) {
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

// Settings Methods

func (a *App) GetSettings() (*config.Config, error) {
	return a.config, nil
}

func (a *App) SaveSettings(cfg config.Config) error {
	a.config.ResticCommand = cfg.ResticCommand
	// Update wrapper command if repo is open (optional, but good practice)
	if a.repoController.Wrapper != nil {
		a.repoController.Wrapper.Command = cfg.ResticCommand
	}
	return config.Save(a.config)
}
