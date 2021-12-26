package subcmd

import "github.com/spf13/cobra"

func InitValidateCMD() *cobra.Command {
	validate := &cobra.Command{
		Use: "validate",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
	return validate
}
