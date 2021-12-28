package subcmd

import (
	"github.com/spf13/cobra"
)

func InitStopCMD() *cobra.Command {
	var serverID uint
	stop := &cobra.Command{
		Use: "stop",
		Run: func(cmd *cobra.Command, args []string) {
			// We want to stop the server.
			// Still, we need to shut down gracefully, what it means that we should shut down inside the game,
			// for example, someone type 'exit' in valve game console, etc.
			//
			// This is aiming to provide enough time for server to stop it correctly.
			// eg: save the world, save user data, etc.
			//
			// 	Also, this action is also related to the coordinator, for searching whether
			// this server is online.
			//
			// Steps
			// 	1. Retrieve online servers from the coordinator, and find this server by index.
			// 	2. Send exit game server command, by websocket interact.
			// 	eg: exit of Valve games, etc.
			// 	3. The game server process has been closed, now this server wrapper should also
			// closed.
			// 	Also, provide this process's exit code to the coordinator is required.
			// 	4. Pop this server out of online servers list.
		},
	}
	stop.Flags().UintVar(&serverID, "server-id", 0, "")
	stop.MarkFlagRequired("server-id")
	return stop
}
