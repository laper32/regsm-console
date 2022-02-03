package subcmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func InitStartCMD() *cobra.Command {
	start := &cobra.Command{
		Use: "start",
		Run: func(cmd *cobra.Command, args []string) {
			// Starts the coordinator.
			// 	The coordinator, which is the 'hub' to interact with front-end with backend.
			// We use this because we do not have any idea to do something like tmux in linux...
			//
			// The step for starting the coordinator:
			// 	1. Check whether the coordinator has been started. (In fact, checking the port is OK)
			// 	2. Establishing websocket server.
			// 	3. Waiting servers/CLI connect.

			// The coordinator is at ${GSM_PATH}/gsm-coordinator.exe

			// How could we know when the coordinator starts failed?
			start := make(chan bool)
			go func() {
				cmd := &exec.Cmd{
					Path: fmt.Sprintf("%v/gsm-coordinator.exe", os.Getenv("GSM_PATH")),
					Dir:  os.Getenv("GSM_PATH"),
					Env:  os.Environ(),
				}
				err := cmd.Start()
				if err != nil {
					fmt.Println("ERROR:", err)
					start <- false
					return
				}
				fmt.Println("Coordinator has been started. Process ID:", cmd.Process.Pid)
				start <- true
			}()
			<-start
		},
	}

	return start
}
