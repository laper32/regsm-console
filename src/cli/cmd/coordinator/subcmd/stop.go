package subcmd

import "github.com/spf13/cobra"

func InitStopCMD() *cobra.Command {
	var recursive string
	stop := &cobra.Command{
		Use: "stop",
		Run: func(cmd *cobra.Command, args []string) {
			if recursive == "all" {
				return
			} else if recursive == "server" {
				return
			} else if recursive == "coordinator" {
				return
			} else if recursive == "none" {
				return
			} else {
				return
			}
		},
	}
	stop.Flags().StringVar(&recursive, "recursive", "none", "Recursively stopping servers.")
	stop.Flags().Lookup("recursive").NoOptDefVal = "none"
	return stop
}
