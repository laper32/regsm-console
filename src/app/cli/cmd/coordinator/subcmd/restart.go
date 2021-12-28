package subcmd

import "github.com/spf13/cobra"

func InitRestartCMD() *cobra.Command {
	restart := &cobra.Command{
		Use: "restart",
		Run: func(cmd *cobra.Command, args []string) {
			// The restart is just only the combination of stop, and start
			// Step:
			// 	1. Stop the coordinator
			// 	2. Start the coordinator
		},
	}
	return restart
}
