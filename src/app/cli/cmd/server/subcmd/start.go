package subcmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/laper32/regsm-console/src/app/cli/dpkg"
	"github.com/laper32/regsm-console/src/lib/log"
	"github.com/spf13/cobra"
)

func InitStartCMD() *cobra.Command {
	var serverID uint
	start := &cobra.Command{
		Use: "start",
		Run: func(cmd *cobra.Command, args []string) {

			var thisServer *dpkg.ServerIdentity

			serverExist := func() bool {
				for _, content := range dpkg.ServerIdentityList {
					if serverID == content.ID {
						// We also need to make sure that this server is still existing...
						if !content.Deleted {
							thisServer = &content
							return true
						} else {
							return false
						}
					}
				}
				return false
			}()

			if !serverExist {
				fmt.Println("ERROR: The server does not exist.")
				return
			}

			// We have GSM_ROOT, and GSM_PATH => Make final configuration at daemon.
			os.Setenv("GSM_SERVER_ID", fmt.Sprintf("%v", thisServer.ID))
			var path string
			if runtime.GOOS == "windows" {
				path = fmt.Sprintf("%v/gsm-server.exe", os.Getenv("GSM_PATH"))
			} else if runtime.GOOS == "linux" {
				path = fmt.Sprintf("%v/gsm-server", os.Getenv("GSM_PATH"))
			} else {
				panic("Not supported OS.")
			}
			exe := &exec.Cmd{
				Path:  path,
				Dir:   os.Getenv("GSM_PATH"),
				Env:   os.Environ(),
				Stdin: os.Stdin,
			}
			err := exe.Start()
			if err != nil {
				log.Panic(err)
				return
			}
			log.Info("Server started. Daemon PID:", exe.Process.Pid)
			exe.Process.Release()
		},
	}
	start.Flags().UintVar(&serverID, "server-id", 0, "Server ID to start")
	return start
}
