package subcmd

import "github.com/spf13/cobra"

func InitRestartCMD() *cobra.Command {
	restart := &cobra.Command{
		Use: "restart",
		Run: func(cmd *cobra.Command, args []string) {
			// Just stop and start.
			// Attention to the sequence!

			InitStopCMD().Execute()
			InitStartCMD().Execute()
		},
	}
	return restart
}
