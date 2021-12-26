package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {

	// Process
	// 	1. Connect to coordinator
	// 	2. Redirect stdin to socket connection receive
	// 	3. Start this server
	// 	4. Write start info (PID, etc)
	// 	5. Send to coordinator to make a presistent storage.
	// 	6. DONE.

	type Application struct {
		ID         uint     `json:"server_id"`
		Dir        string   `json:"dir"`
		Executable string   `json:"executable"`
		Args       []string `json:"args"`
	}
	var app Application
	json.Unmarshal([]byte(os.Args[1]), &app)

	// We cannot pass in []strings{} directly.
	// Because for example: if we pass the slice directly, we will found that 'de_dust2"' not found.
	// Yeah, wtf, why occur '"'?
	// Not very sure how golang pass in the executable params, and don't know the implementation detail...
	// Overall, stupid, very stupid
	// Based on this, we need to construct a parameter string
	// And we need also add up a space to make sure it can be executed correctly.
	// Test passed on CS1.6
	serverParamStr := strings.Join(app.Args, " ")
	serverEXE := exec.Command(app.Executable, serverParamStr+" ")
	serverEXE.Stdin = os.Stdin
	serverEXE.Stdout = os.Stdout
	serverEXE.Stderr = os.Stderr
	serverEXE.Dir = app.Dir
	err := serverEXE.Start()
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}
	fmt.Printf("Server started. Server ID: %v(Process ID: %v)", app.ID, serverEXE.Process.Pid)
	type Status struct {
		ServerID  uint
		ProcessID int
		Done      bool
	}
	status := &Status{
		ServerID:  app.ID,
		ProcessID: serverEXE.Process.Pid,
	}
	json.Marshal(status)
}
