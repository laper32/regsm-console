package subcmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/laper32/regsm-console/src/app/cli/cmd/server/subcmd/sysprocattr"
	"github.com/laper32/regsm-console/src/app/cli/dpkg"
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
				Path:        path,
				Dir:         os.Getenv("GSM_PATH"),
				SysProcAttr: sysprocattr.SysProcAttr(),
				Env:         os.Environ(),
				Stdout:      os.Stdout,
				Stderr:      os.Stderr,
			}
			err := exe.Start()
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("Server started. Daemon PID:", exe.Process.Pid)
			os.Exit(0)

			// // That's why we want to build a coordinator....
			// // We need a method to see all information...
			// start := make(chan bool)
			// go func() {
			// 	wrapperEXE := &exec.Cmd{
			// 		Path:        wrapperEXEPath,
			// 		Dir:         os.Getenv("GSM_ROOT") + "/bin",
			// 		Args:        []string{wrapperEXEPath, string(ret)},
			// 		SysProcAttr: sysprocattr.SysProcAttr(),
			// 	}
			// 	// TODO: Notify the server that this server is starting.
			// 	err = wrapperEXE.Start()
			// 	if err != nil {
			// 		fmt.Println("ERROR:", err)
			// 		start <- false
			// 		return
			// 	}
			// 	start <- true
			// }()
			// // <-start

		},
	}
	start.Flags().UintVar(&serverID, "server-id", 0, "Server ID to start")
	return start
}
