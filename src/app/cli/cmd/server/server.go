package server

import (
	"github.com/laper32/regsm-console/src/app/cli/cmd/server/subcmd"
	"github.com/spf13/cobra"
)

func InitCMD() *cobra.Command {
	server := &cobra.Command{
		Use: "server",
	}
	server.AddCommand(
		subcmd.InitBackupCMD(),
		subcmd.InitConsoleCMD(),
		subcmd.InitInstallCMD(),
		subcmd.InitRemoveCMD(),
		subcmd.InitRestartCMD(),
		subcmd.InitStartCMD(),
		subcmd.InitStopCMD(),
		subcmd.InitUpdateCMD(),
		subcmd.InitValidateCMD(),
	)
	return server
}