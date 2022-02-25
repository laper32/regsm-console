package subcmd

import (
	"fmt"
	"os"
	"os/exec"

	cliconf "github.com/laper32/regsm-console/src/app/cli/conf"
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

			// If windows: gsm-coordinator.exe
			// otherwise gsm-coordinator
			_, err := cliconf.CoordinatorConfiguration()
			if err != nil {
				fmt.Println("ERROR:", err)
				return
			}
			exedir := fmt.Sprintf("%v/gsm-coordinator.exe", os.Getenv("GSM_PATH"))
			exe := &exec.Cmd{
				Path:   exedir,
				Dir:    os.Getenv("GSM_PATH"),
				Env:    os.Environ(),
				Stderr: os.Stderr,
				Stdout: os.Stdout,
			}
			exe.Run()

			// How could we know when the coordinator starts failed?
			// start := make(chan bool)
			// var pid int
			// go func() {
			// 	exe := &exec.Cmd{
			// 		Path: exedir,
			// 		Dir:  os.Getenv("GSM_PATH"),
			// 		Env:  os.Environ(),
			// 	}
			// 	err = exe.Start()
			// 	if err != nil {
			// 		fmt.Println("ERROR:", err)
			// 		start <- false
			// 		return
			// 	}
			// 	pid = exe.Process.Pid
			// 	start <- true
			// }()
			// result := <-start
			// if result {
			// 	fmt.Println("Coordinator started. Process ID:", pid)
			// 	return
			// } else {
			// 	fmt.Println("Coordinator starting failed. Message:", err)
			// 	return
			// }
		},
	}

	return start
}
