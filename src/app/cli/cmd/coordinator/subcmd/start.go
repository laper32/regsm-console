package subcmd

import (
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
		},
	}

	return start
}
