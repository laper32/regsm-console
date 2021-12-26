package subcmd

import "github.com/spf13/cobra"

func InitConsoleCMD() *cobra.Command {
	console := &cobra.Command{
		Use: "console",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
	return console
}
