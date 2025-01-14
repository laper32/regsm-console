package subcmd

import (
	"fmt"
	"net/url"

	"github.com/gorilla/websocket"
	cliconf "github.com/laper32/regsm-console/src/app/cli/conf"
	"github.com/laper32/regsm-console/src/app/cli/misc"
	"github.com/laper32/regsm-console/src/lib/log"
	"github.com/laper32/regsm-console/src/lib/status"
	"github.com/spf13/cobra"
)

func InitStopCMD() *cobra.Command {
	// var recursive string
	stop := &cobra.Command{
		Use: "stop",
		Run: func(cmd *cobra.Command, args []string) {
			// 	Stopping the coordinator.
			// A little bit complicated.
			// 	We need to consider that: What should we handle those connections after we closed
			// this coordinator?
			//
			// 	In order to handle this issue, that we have this solution:
			// When we shut down the coordinator, we will 'notify' connections to do further action.
			// eg: We shut down this coordinator, if this coordinator's connections are all server
			// that we can select 'shut down' all servers (if needed), and then these server are all stopped.
			// Coordinators can also do such things...
			//
			// 	Or, we can do nothing, if we just want to update coordinators, or restarting due to
			// reboot.
			//
			// Steps:
			// 	1. Before we stop the coordinator, we should check the parameter.
			// If, all: We will also shut down all connections.
			// If, server: Will shut down all connections, but only identified as 'server'.
			// If, none: No connections will be terminated.
			// Everything will be manipulated on net, in other words, if it fails to shut down, then
			// shut down process will fail (that you should check connections.)
			// 	2. Shut down coordinator.
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
			defer conn.Close()
			log.Info("Sending stopping command")
			err = conn.WriteJSON(&retGram)
			if err != nil {
				log.Info(err)
				return
			}
		},
	}
	// stop.Flags().StringVar(&recursive, "recursive", "none", "Recursively stopping servers.")
	// stop.Flags().Lookup("recursive").NoOptDefVal = "none"
	return stop
}
