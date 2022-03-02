package subcmd

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"

	"github.com/gorilla/websocket"
	cliconf "github.com/laper32/regsm-console/src/app/cli/conf"
	"github.com/laper32/regsm-console/src/app/cli/misc"
	"github.com/laper32/regsm-console/src/lib/log"
	"github.com/laper32/regsm-console/src/lib/status"
	"github.com/spf13/cobra"
)

func InitRestartCMD() *cobra.Command {
	restart := &cobra.Command{
		Use: "restart",
		Run: func(cmd *cobra.Command, args []string) {
			// The restart is just only the combination of stop, and start
			// Step:
			// 	1. Stop the coordinator
			// 	2. Start the coordinator
			cfg, err := cliconf.CoordinatorConfiguration()
			if err != nil {
				fmt.Println("ERROR:", err)
				return
			}
			url := &url.URL{
				Scheme: "ws",
				Host:   fmt.Sprintf("%v:%v", cfg.GetString("coordinator.ip"), cfg.GetUint("coordinator.port")),
			}
			fmt.Printf("[%v] Connecting to the coordinator...", url.String())
			conn, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
			if err != nil {
				fmt.Println()
				log.Warn("Seems that the connection is closed. Message:", err)
				return
			}
			fmt.Println("OK")
			retGram := &RetGram{
				Role:    misc.Role,
				Code:    status.CLICoordinatorSendStopSignal.ToInt(),
				Message: status.CLICoordinatorSendStopSignal.Message(),
				Detail:  map[string]interface{}{},
			}
			log.Info("Sending stopping command")
			err = conn.WriteJSON(&retGram)
			if err != nil {
				log.Info(err)
				return
			}
			for {
				_, _, err := conn.ReadMessage()
				if err != nil {
					break
				}
			}
			exedir := fmt.Sprintf("%v/gsm-coordinator.exe", os.Getenv("GSM_PATH"))
			ok := make(chan bool)
			go func() {
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
				log.Info("Coordinator started. Process ID:", exe.Process.Pid)
				ok <- true
			}()
			<-ok
		},
	}
	return restart
}
