package cmd

import (
	"fmt"

	"github.com/laper32/regsm-console/src/app/cli/cmd/coordinator"
	"github.com/laper32/regsm-console/src/app/cli/cmd/server"
	"github.com/laper32/regsm-console/src/app/cli/misc"
	"github.com/spf13/cobra"
)

// Note: The CMD command level reaches at most 3.

func initUpdateCMD() *cobra.Command {
	var content string
	update := &cobra.Command{
		Use: "update",
		Run: func(cmd *cobra.Command, args []string) {
			if content == "all" {
				fmt.Println("Update everything.")
				return
			} else if content == "cli" {
				fmt.Println("Update cli.")
				return
			} else if content == "coordinator" {
				fmt.Println("Update coordinator.")
				return
			} else if content == "server" {
				fmt.Println("Update server.")
				return
			} else {
				fmt.Println("Unknown content input: ", content)
				return
			}
		},
	}
	update.Flags().StringVar(&content, "content", "all", "Which content you want to update")
	update.Flags().Lookup("content").NoOptDefVal = "all"
	return update
}

func InitCMD() *cobra.Command {
	gsm := &cobra.Command{
		Use: "gsm",
	}
	gsm.AddCommand(
		coordinator.InitCMD(),
		server.InitCMD(),
		initUpdateCMD(),
	)
	gsm.PersistentFlags().BoolVarP(&misc.Agree, "yes", "y", false, "Confirm.")
	gsm.PersistentFlags().Lookup("yes").NoOptDefVal = "true"
	gsm.PersistentFlags().MarkHidden("yes")

	gsm.PersistentFlags().BoolVarP(&misc.Decline, "no", "n", false, "Decline.")
	gsm.PersistentFlags().Lookup("no").NoOptDefVal = "true"
	gsm.PersistentFlags().MarkHidden("no")

	return gsm
}
