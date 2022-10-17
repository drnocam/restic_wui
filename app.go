package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"restic_wui/fileop"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

var selected_repository string

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {

	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) CmdCheck() string {

	cmd := exec.Command("restic")
	err := cmd.Start()
	if err != nil {
		return fmt.Sprint(JsonReturn(Message{}, "0"))
	}
	return fmt.Sprint(JsonReturn(Message{}, "1"))
}

// repository choose
func (a *App) ChooseRepository() string {
	cwd, _ := os.Getwd()
	selection, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title:            "Choose Repository",
		DefaultDirectory: cwd,
	})
	if err != nil {
		log.Print("Error opening directory")

	}
	selected_repository = selection

	size, _ := fileop.DirSize(selected_repository)
	fmt.Printf("Repository %s , Size: %d", selection, size)

	return fmt.Sprint(selection)
}

func (a *App) GetSnapshots(password string) string {

	cmd := exec.Command("restic", "-r", selected_repository, "snapshots", "--json")

	newEnv := append(os.Environ(), "RESTIC_PASSWORD="+password)
	cmd.Env = newEnv
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}

	return fmt.Sprintf("%s\n", string(out))

}

func (a *App) ReadWriteSettings() {

	if data, err := json.Marshal(GetSettings()); err == nil {
		fileop.WriteSettings(data)
	}

	fileop.ReadSettings()

	fmt.Printf("Dosya yazıldı.")

}

func (a *App) AddUpdateRepository(id int, infos string) string {

	if id == 0 {
		/*  new record */
		new_repo := SavedRepository{}
		err := json.Unmarshal([]byte(infos), &new_repo)
		if err == nil {
			/*
				TODOS: get last id
			*/

			new_repo.id = 5
			if data, e := json.Marshal(new_repo); e == nil {
				fileop.WriteSettings(data)
				return JsonReturn(Message{Message: "New Repository Saved Succesfully", Type: 1}, infos)
			} else {
				return JsonReturn(Message{Message: "Couldnt Save", Type: 0}, infos)
			}

		}
		fmt.Printf(err.Error())
	}
	return JsonReturn(Message{Message: "Couldnt Save", Type: 0}, infos)

}
