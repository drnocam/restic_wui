/*
TODO: merge restic commands for one function
time waiting for some animations may remove at future

*/

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"restic_wui/fileop"
	"strings"
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

	runtime.EventsEmit(a.ctx, "event_test", EventTest())

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

	//time.Sleep(350 * time.Millisecond)

	settings_index := settings.FindIndexById(id)

	if settings_index != -1 {

		selected_repository := &settings.Repositories[settings_index]
		cmd := exec.Command("restic", "-r", selected_repository.Path, "snapshots", "--json")

		newEnv := append(os.Environ(), "RESTIC_PASSWORD="+selected_repository.Password)
		cmd.Env = newEnv
		out, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Sprint(JsonReturn(Message{0, "Error occured. No repository found!"}, "{\"error\":1}"))
		}

		return fmt.Sprint(JsonReturn(Message{-1, ""}, string(out)))
	}

	return fmt.Sprint(JsonReturn(Message{0, "Couldnt find a repository"}, "{\"error\":1}"))

}

func (a *App) ReadWriteSettings() {

	GetSettings()

}

func (a *App) AddUpdateRepository(id int, infos string) string {
	time.Sleep(350 * time.Millisecond)
	new_repo := SavedRepository{}
	err_json := json.Unmarshal([]byte(infos), &new_repo)
	/*
	   TODO dont let same names
	*/
	if err_json == nil {
		if id == -1 {
			/*  new record */
			/*
				if folder exists and empty then create init.
				if folder exists and not empty then check config file exists so this is restic repo
				else give error and do not add settings.
			*/
			if ok, err := fileop.Exists(new_repo.Path); ok {

				if ok, err_empty := fileop.IsEmpty(new_repo.Path); ok {

					// restic init pass....
					cmd := exec.Command("restic", "-r", new_repo.Path, "init", "--json")

					newEnv := append(os.Environ(), "RESTIC_PASSWORD="+new_repo.Password)
					cmd.Env = newEnv
					out, err := cmd.CombinedOutput()
					if err != nil {
						return fmt.Sprint(JsonReturn(Message{0, "Error occured. No repository found!"}, fmt.Sprintf("%q", err_empty.Error())))
					}
					//return fmt.Sprint(JsonReturn(Message{1, fmt.Sprintf("%q", string(out))}, "{}"))
					data := fmt.Sprintf("{%q:1,\"names\":%s}", "cmd_exists", settings.RepoNickNames())
					return fmt.Sprint(JsonReturn(Message{1, fmt.Sprintf("New Repository Settings Saved Succesfully<br>%q", out)}, fmt.Sprintf("%q", string(data))))

				} else {

					if ok, err := fileop.Exists(new_repo.Path + "/config"); ok {

						settings.AddRepoSettings(new_repo)
						/*
							to see what is added
						*/

						data := fmt.Sprintf("{%q:1,\"names\":%s}", "cmd_exists", settings.RepoNickNames())
						return fmt.Sprint(JsonReturn(Message{1, "New Repository Settings Saved Succesfully"}, fmt.Sprintf("%q", string(data))))

					} else {
						err_msg := fmt.Sprintf("Couldnt open the config file : %s ", err.Error())
						return JsonReturn(Message{Message: err_msg, Type: 0}, "{\"error\":1}")

					}
				}

			} else {

				err_msg := fmt.Sprintf("Couldnt open the folder : %s ", err.Error())
				return JsonReturn(Message{Message: err_msg, Type: 0}, "{\"error\":1}")

			}

		} else {
			/*
				update info
			*/
			settings.UpdateRepoSettings(id, new_repo)
			data := fmt.Sprintf("{%q:1,\"names\":%s}", "cmd_exists", settings.RepoNickNames())
			return fmt.Sprint(JsonReturn(Message{1, "New Repository Saved Succesfully"}, fmt.Sprintf("%q", string(data))))
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

func (a *App) GetRepoStats(id int) string {

	time.Sleep(350 * time.Millisecond)

	settings_index := settings.FindIndexById(id)

	if settings_index != -1 {

		selected_repository := &settings.Repositories[settings_index]

		cmd := exec.Command("restic", "-r", selected_repository.Path, "stats", "--json")

		newEnv := append(os.Environ(), "RESTIC_PASSWORD="+selected_repository.Password)
		cmd.Env = newEnv
		out, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Sprint(JsonReturn(Message{0, "Error occured. No repository found!"}, "{\"error\":1}"))
		}
		return fmt.Sprint(JsonReturn(Message{-1, ""}, string(out)))
	}
	return fmt.Sprint(JsonReturn(Message{0, "Error occured. No repository found!"}, "{\"error\":1}"))
}

func (a *App) CheckRepoErrors(id int) string {

	time.Sleep(350 * time.Millisecond)

	settings_index := settings.FindIndexById(id)

	if settings_index != -1 {

		selected_repository := &settings.Repositories[settings_index]

		cmd := exec.Command("restic", "-r", selected_repository.Path, "check", "--json")

		newEnv := append(os.Environ(), "RESTIC_PASSWORD="+selected_repository.Password)
		cmd.Env = newEnv
		out, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Sprint(JsonReturn(Message{0, "Error occured. No repository found!"}, "{\"error\":1}"))
		}
		return fmt.Sprint(JsonReturn(Message{-1, ""}, fmt.Sprintf("%q", string(out))))
	}
	return fmt.Sprint(JsonReturn(Message{0, "Error occured. No repository found!"}, "{\"error\":1}"))
}

func (a *App) SearchInRepo(id int, search_text string) string {

	settings_index := settings.FindIndexById(id)

	if settings_index != -1 {

		selected_repository := &settings.Repositories[settings_index]

		cmd := exec.Command("restic", "-r", selected_repository.Path, "--json", "find", search_text)

		newEnv := append(os.Environ(), "RESTIC_PASSWORD="+selected_repository.Password)
		cmd.Env = newEnv

		fmt.Println(cmd.String())

		out, err := cmd.CombinedOutput()

		if err != nil {
			return fmt.Sprint(JsonReturn(Message{0, "Error occured. No repository found!"}, "{\"error\":1}"))
		}

		return fmt.Sprint(JsonReturn(Message{-1, ""}, string(out)))
	}
	return fmt.Sprint(JsonReturn(Message{0, "Error occured. No repository found!"}, "{\"error\":1}"))
}

func (a *App) RestoreRepo(id int, snapshot_id string) string {

	settings_index := settings.FindIndexById(id)

	if settings_index != -1 {

		selected_repository := &settings.Repositories[settings_index]

		cwd, _ := os.Getwd()
		target_dir, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
			Title:            "Choose Directory To Restore",
			DefaultDirectory: cwd,
		})
		if err != nil {
			return fmt.Sprint(JsonReturn(Message{0, "Directory couldnt opened!"}, "{\"error\":1}"))
		}

		cmd := exec.Command("restic", "-r", selected_repository.Path, "restore", snapshot_id, "--target", fmt.Sprintf("%q", target_dir))

		newEnv := append(os.Environ(), "RESTIC_PASSWORD="+selected_repository.Password)
		cmd.Env = newEnv
		out, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Sprint(JsonReturn(Message{0, "Error occured. No repository found!"}, "{\"error\":1}"))
		}
		return fmt.Sprint(JsonReturn(Message{-1, ""}, fmt.Sprintf("%q", string(out))))
	}
	return fmt.Sprint(JsonReturn(Message{0, "Error occured. No repository found!"}, "{\"error\":1}"))
}

func (a *App) ListFilesInSnapshots(id int, snapshot_id string) string {

	settings_index := settings.FindIndexById(id)

	if settings_index != -1 {

		selected_repository := &settings.Repositories[settings_index]

		cmd := exec.Command("restic", "-r", selected_repository.Path, "ls", snapshot_id, "--json")

		newEnv := append(os.Environ(), "RESTIC_PASSWORD="+selected_repository.Password)
		cmd.Env = newEnv
		out, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Sprint(JsonReturn(Message{0, "Error occured. No repository found!"}, "{\"error\":1}"))
		}
		add_comma := strings.Replace(string(out), "}\n{", "},{", -1)
		return fmt.Sprint(JsonReturn(Message{-1, ""}, fmt.Sprintf("[%s]", add_comma)))
	}
	return fmt.Sprint(JsonReturn(Message{0, "Error occured. No repository found!"}, "{\"error\":1}"))
}

func (a *App) RestoreFilesInSnapshots(id int, snapshot_id string, files_json string) string {

	settings_index := settings.FindIndexById(id)

	if settings_index != -1 {

		selected_repository := &settings.Repositories[settings_index]

		//files := [string]
		var arr []string
		_ = json.Unmarshal([]byte(files_json), &arr)

		cwd, _ := os.Getwd()
		target_dir, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
			Title:            "Choose Directory To Restore",
			DefaultDirectory: cwd,
		})
		if err != nil {
			return fmt.Sprint(JsonReturn(Message{0, "Directory couldnt opened!"}, "{\"error\":1}"))
		}

		cmd := exec.Command("restic", "-r", selected_repository.Path, "restore", snapshot_id, "--target", target_dir, "--json")
		for i := 0; i < len(arr); i++ {
			cmd.Args = append(cmd.Args, "--include")
			cmd.Args = append(cmd.Args, arr[i])
		}
		// fmt.Println(cmd.String())
		// fmt.Printf("%v", cmd.Args)
		newEnv := append(os.Environ(), "RESTIC_PASSWORD="+selected_repository.Password)
		cmd.Env = newEnv
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println(string(out))
			return fmt.Sprint(JsonReturn(Message{0, "Error occured. No repository found!"}, "{\"error\":1}"))
		}
		return fmt.Sprint(JsonReturn(Message{-1, ""}, fmt.Sprintf("%q", string(out))))
	}
	return fmt.Sprint(JsonReturn(Message{0, "Error occured. No repository found!"}, "{\"error\":1}"))
}

func EventTest() string {
	return fmt.Sprint("event test pushed a value")
}
