package subcmd

import "github.com/spf13/cobra"

func InitValidateCMD() *cobra.Command {
	validate := &cobra.Command{
		Use: "validate",
		Run: func(cmd *cobra.Command, args []string) {
			// Validating this server's file integrity.
			//
			// In general, this is quite similar to install/update.
			// I won't do further introduction here.
		},
	}
	return validate
}
