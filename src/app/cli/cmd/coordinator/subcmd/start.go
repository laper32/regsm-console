package subcmd

import (
	"fmt"
	"os"
	"os/exec"

	cliconf "github.com/laper32/regsm-console/src/app/cli/conf"
	"github.com/laper32/regsm-console/src/lib/log"
	"github.com/spf13/cobra"
)

type RetGram struct {
	Role    string                 `json:"role"`
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Detail  map[string]interface{} `json:"detail"`
}

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
				Path: exedir,
				Dir:  os.Getenv("GSM_PATH"),
				Env:  os.Environ(),
			}
			err = exe.Start()
			if err != nil {
				log.Panic(err)
				return
			}
			log.Info("Coordinator started. PID:", exe.Process.Pid)
			exe.Process.Release()
		},
	}

	return start
}
