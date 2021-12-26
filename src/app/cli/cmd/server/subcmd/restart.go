package subcmd

import "github.com/spf13/cobra"

func InitRestartCMD() *cobra.Command {
	restart := &cobra.Command{
		Use: "restart",
		Run: func(cmd *cobra.Command, args []string) {
			// restart...
			// yeah... stop and start...
			InitStopCMD().Execute()
			InitStartCMD().Execute()
		},
	}
	return restart
}
