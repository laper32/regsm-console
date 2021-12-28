package subcmd

import "github.com/spf13/cobra"

func InitConsoleCMD() *cobra.Command {
	console := &cobra.Command{
		Use: "console",
		Run: func(cmd *cobra.Command, args []string) {
			// 	I still don't really know how to do...
			//
			// 	You know, we cannot use anything about platform-specific,
			// like tmux on Unix, etc.
			//
			// 	That we can only manipulate everything via the internet.
			//
			// 	We now assuming using websocket to handle this.
			//
			// 	After we started the server, the standard IO should has been
			// connected to the coordinator.
			// 	Based on this, when we want to connect to a specific server
			// for further action, we can use something like CLI, webpage,
			// desktop, balabala.
			//
			// 	This is quite similar to something like SSH (Yeah, SSH requires
			// internet connection, that should we say is it the best solution?)
			//
			// Steps:
			// 	1. Search the online server list, find this server.
			// 	2. The frontend establish connection to the coordinator, and the
			// coordinator redirects the frontend's connection to this specific server.
			// 	3. Now everything is set, you can handle this server.
		},
	}
	return console
}
