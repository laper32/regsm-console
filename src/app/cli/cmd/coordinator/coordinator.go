package coordinator

import (
	"github.com/laper32/regsm-console/src/app/cli/cmd/coordinator/subcmd"
	"github.com/spf13/cobra"
)

/*
gsm coordinator start
gsm coordinator stop
*/

func InitCMD() *cobra.Command {
	coordinator := &cobra.Command{
		Use: "coordinator",
	}
	coordinator.AddCommand(
		subcmd.InitRestartCMD(),
		subcmd.InitStartCMD(),
		subcmd.InitStopCMD(),
	)
	return coordinator
}
