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

// Check if command exist
// * repository list
func (a *App) CmdCheck() string {

	cmd := exec.Command("restic")
	err := cmd.Start()
	var msg string
	names := make(map[int]string)
	if err != nil {
		msg = "Restic Command is not Exist or Configured correctly\n"
		return fmt.Sprint(JsonReturn(Message{0, msg}, ""))
	} else {
		msg = "Restic Command Configured Correctly\n"
	}
	GetSettings()
	for _, v := range settings.Repositories {
		names[v.Id] = v.Name
	}
	nms, _ := json.Marshal(names)
	data := fmt.Sprintf("%q:1,\"names\":%s", "cmd_exists", string(nms))
	return fmt.Sprint(JsonReturn(Message{1, msg}, string(data)))
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

	GetSettings()

}

func (a *App) AddUpdateRepository(id int, infos string) string {

	if id == 0 {
		/*  new record */
		new_repo := SavedRepository{}
		err := json.Unmarshal([]byte(infos), &new_repo)
		if err == nil {

			settings.AddRepository(new_repo)
			/*
				to see what is added
			*/
			if s, ok := json.Marshal(settings); ok == nil {

				return JsonReturn(Message{Message: "New Repository Saved Succesfully", Type: 1}, string(s))
			} else {

				return JsonReturn(Message{Message: "Couldnt Save", Type: 0}, infos)
			}

		}
		fmt.Printf(err.Error())
	}
	return JsonReturn(Message{Message: "Couldnt Save", Type: 0}, infos)

}
