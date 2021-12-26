package subcmd

import (
	"github.com/spf13/cobra"
)

func InitStopCMD() *cobra.Command {
	var serverID uint
	stop := &cobra.Command{
		Use: "stop",
		Run: func(cmd *cobra.Command, args []string) {
			// Process
			//	1. Send 'STOP' to server
			// 	2. Server retrieved 'STOP' signal
			// 	3. Stop the server
		},
	}
	stop.Flags().UintVar(&serverID, "server-id", 0, "")
	stop.MarkFlagRequired("server-id")
	return stop
}
