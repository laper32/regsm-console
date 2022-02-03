package subcmd

import "github.com/spf13/cobra"

func InitAttachCMD() *cobra.Command {
	attach := &cobra.Command{
		Use: "attach",
		Run: func(cmd *cobra.Command, args []string) {
			// We will move to gsm-server and gsm-coordinator to solve this problem.
			// On windows, we will use C#, aka, we will use C# to complete this part.
			// This is because golang does not have enough ability to manipulate system API.
			// Since it just a simple tool, that we dont have to use c/c++.
			// Also, we want to have chance to manipulate networking (querying server, etc).
			// Based on those information above, C# is fit.

		},
	}
	return attach
}
