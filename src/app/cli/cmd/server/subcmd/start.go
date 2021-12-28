package subcmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"

	cliconf "github.com/laper32/regsm-console/src/app/cli/conf"
	"github.com/laper32/regsm-console/src/app/cli/dpkg"
	"github.com/spf13/cobra"
)

func InitStartCMD() *cobra.Command {
	var serverID uint
	start := &cobra.Command{
		Use: "start",
		Run: func(cmd *cobra.Command, args []string) {
			// 	Starting a server is a easy job, but the difficulty is that
			// what should we do if we also want to interact the game server
			// console.
			// 	It is not an easy job, since the Windows do not have tmux,
			// and we want to build a cross-platform project.
			//
			// 	Based on this, we need to redirect the game server's standard IO
			// to somewhere, for example, websocket, and connect it to a public
			// server, for we can manipulate this server remotely.
			//
			// Steps
			// 	1. Check this server whether exists.
			// 	2. Configure this game, do the final configuration.
			// 	3. Check the executable's folder.
			// 	4. Connect this server to the coordinator.
			// In this term, we just state it as 'coonected'.
			// Don't forget that we need also provide a standard IO
			// by websocket.
			// 	6. Start the server. Redirect the standard IO
			// what have been established by websocket to this
			// process.
			// 	7. Now everything is OK. The server's state is 'OK'.
			// Send the server's process ID, startup time, etc to the
			// coordinator.
			// 	8. Now everything is set, and the server is ready to go.

			serverExist := func() bool {
				for _, content := range dpkg.ServerInfoList {
					if serverID == content.ID {
						return true
					}
				}
				return false
			}()
			if !serverExist {
				fmt.Println("ERROR: The server does not exist.")
				return
			}

			// We export server.param for passing in 'Args' when we executing commands.
			cfg := cliconf.StartupConfiguration(serverID)

			// By default, the executable file is under ${GSM_ROOT}/server/${ServerIndex}
			// Some of them may be under ${GSM_ROOT}/server/${ServerIndex}/bin, or whatever
			// If this case occured, will override via 'overrideServerExecutablePath'

			serverExecutableDir := fmt.Sprintf("%v/server/%v", os.Getenv("GSM_ROOT"), serverID)

			overrideServerExecutablePath := func() {

			}

			overrideServerExecutablePath()
			var serverExecutablePath string
			locateExecutable := func() {
				game := cfg.Get("server.game")
				// Source 0
				// This is dirty work, which is NOT avoidable!
				if game == "cs1.6" || game == "czero" {
					serverExecutablePath = serverExecutableDir + "/hlds.exe"
				}
			}

			locateExecutable()

			// We need to pass parameters to gsm-server
			// That for convience, we send params via JSON
			type Application struct {
				ID         uint     `json:"server_id"`
				Dir        string   `json:"dir"`
				Executable string   `json:"executable"`
				Args       []string `json:"args"`
			}
			app := &Application{
				ID:         serverID,
				Dir:        serverExecutableDir,
				Executable: serverExecutablePath,
				Args:       cfg.GetStringSlice("server.param"),
			}
			ret, err := json.Marshal(app)
			if err != nil {
				fmt.Println("ERROR:", err)
				return
			}

			wrapperEXEPath := os.Getenv("GSM_ROOT") + "/bin"
			if runtime.GOOS == "windows" {
				wrapperEXEPath += "/gsm-server.exe"
			} else if runtime.GOOS == "linux" {
				wrapperEXEPath += "/gsm-server"
			} else {
				fmt.Println("Unknown OS, aborting")
				return
			}
			wrapperEXE := &exec.Cmd{
				Path:   wrapperEXEPath,
				Dir:    os.Getenv("GSM_ROOT") + "/bin",
				Args:   []string{wrapperEXEPath, string(ret)},
				Stdin:  os.Stdin,
				Stdout: os.Stdout,
				Stderr: os.Stderr,
			}
			err = wrapperEXE.Start()
			if err != nil {
				fmt.Println(err)
				return
			}
		},
	}
	start.Flags().UintVar(&serverID, "server-id", 0, "Server ID to start")
	return start
}