package subcmd

import (
	"github.com/spf13/cobra"
)

func InitStartCMD() *cobra.Command {
	start := &cobra.Command{
		Use: "start",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	return start
}
