package main

import (
	"context"
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"os"
	"log"
	"os/exec"     
)




var selected_repository string;
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
		return fmt.Sprint(JsonReturn(Message{},"0"))
	}
	return fmt.Sprint(JsonReturn(Message{},"1"))
}

// repository choose
func (a *App) ChooseRepository() string {
	cwd,_ := os.Getwd();
	selection, err := runtime.OpenDirectoryDialog(a.ctx,runtime.OpenDialogOptions{
		Title:        "Choose Repository",
		DefaultDirectory : cwd,
	})
	if(err!=nil){
		log.Print("Error opening directory")

	}
	selected_repository = selection;
	return fmt.Sprintf("Repository %s", selection)
}

func (a *App) GetSnapshots(password string) string {

	cmd := exec.Command("restic", "-r", selected_repository, "snapshots")
	
	newEnv := append(os.Environ(), "RESTIC_PASSWORD=" + password)
	cmd.Env = newEnv
	out, err := cmd.CombinedOutput()
    if err != nil {
        log.Fatalf("cmd.Run() failed with %s\n", err)
    }
    return fmt.Sprintf("%s\n", string(out))


}

