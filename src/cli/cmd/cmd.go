package cmd

// Note: The CMD command level reaches at most 3.

import (
	"fmt"

	"github.com/laper32/regsm-console/src/cli/cmd/coordinator"
	"github.com/laper32/regsm-console/src/cli/cmd/server"
	"github.com/spf13/cobra"
)

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
	return gsm
}
