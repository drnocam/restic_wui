package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"restic_wui/fileop"
	"time"

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
	time.Sleep(350 * time.Millisecond)
	cmd := exec.Command("restic")
	err := cmd.Start()
	var msg string

	if err != nil {
		msg = "Restic Command is not Exist or Configured correctly\n"
		return fmt.Sprint(JsonReturn(Message{0, msg}, ""))
	} else {
		msg = "Restic Command Configured Correctly\n"
	}
	GetSettings()

	data := fmt.Sprintf("{%q:1,\"names\":%s}", "cmd_exists", settings.RepoNickNames())
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

// settings repository id
func (a *App) GetSnapshots(id int) string {

	time.Sleep(350 * time.Millisecond)

	settings_index := settings.FindIndexById(id)

	if settings_index != -1 {

		selected_repository := &settings.Repositories[settings_index]
		cmd := exec.Command("restic", "-r", selected_repository.Path, "snapshots", "--json")

		newEnv := append(os.Environ(), "RESTIC_PASSWORD="+selected_repository.Password)
		cmd.Env = newEnv
		out, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Sprint(JsonReturn(Message{0, "Error occured. No repository found!"}, "{}"))
		}
		return fmt.Sprint(JsonReturn(Message{-1, ""}, string(out)))
	}

	return fmt.Sprint(JsonReturn(Message{0, "Couldnt find a repository"}, "{}"))

}

func (a *App) ReadWriteSettings() {

	GetSettings()

}

func (a *App) AddUpdateRepository(id int, infos string) string {
	time.Sleep(350 * time.Millisecond)
	new_repo := SavedRepository{}
	err := json.Unmarshal([]byte(infos), &new_repo)

	if id == -1 {
		/*  new record */

		if err == nil {

			settings.AddRepoSettings(new_repo)
			/*
				to see what is added
			*/

			data := fmt.Sprintf("{%q:1,\"names\":%s}", "cmd_exists", settings.RepoNickNames())
			return fmt.Sprint(JsonReturn(Message{1, "New Repository Saved Succesfully"}, string(data)))

		}
	} else {
		if err == nil {
			settings.UpdateRepoSettings(id, new_repo)
			data := fmt.Sprintf("{%q:1,\"names\":%s}", "cmd_exists", settings.RepoNickNames())
			return fmt.Sprint(JsonReturn(Message{1, "New Repository Saved Succesfully"}, string(data)))
		}
	}
	return JsonReturn(Message{Message: "Couldnt Save", Type: 0}, infos)

}

func (a *App) DeleteRepositorySetting(id int) string {

	time.Sleep(350 * time.Millisecond)
	if id != -1 {
		settings.DeleteRepoFromSettings(id)
		data := fmt.Sprintf("{%q:1,\"names\":%s}", "cmd_exists", settings.RepoNickNames())
		return fmt.Sprint(JsonReturn(Message{1, "Save Updated"}, string(data)))

	}
	return JsonReturn(Message{Message: "Repository not selected", Type: 0}, "")

}

/*
* delete from disk and settings
 */
func (a *App) DeleteRepositoryFromDisk(id int) string {

	time.Sleep(350 * time.Millisecond)
	if id != -1 {
		settings.DeleteRepoFromSettings(id)
		data := fmt.Sprintf("{%q:1,\"names\":%s}", "cmd_exists", settings.RepoNickNames())
		return fmt.Sprint(JsonReturn(Message{1, "Save Updated"}, string(data)))

	}
	return JsonReturn(Message{Message: "Repository not selected", Type: 0}, "")

}

func (a *App) GetRepoInfo(id int) string {
	return JsonReturn(Message{Message: "", Type: 1}, settings.GetRepoInfo(id))
}
