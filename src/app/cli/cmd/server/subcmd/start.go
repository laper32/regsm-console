package subcmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/laper32/regsm-console/src/app/cli/cmd/server/subcmd/sysprocattr"
	cliconf "github.com/laper32/regsm-console/src/app/cli/conf"
	"github.com/laper32/regsm-console/src/app/cli/dpkg"
	"github.com/spf13/cobra"
)

func InitStartCMD() *cobra.Command {
	var serverID uint
	start := &cobra.Command{
		Use: "start",
		Run: func(cmd *cobra.Command, args []string) {
			// 		Starting server can be done through a daemon.
			// 		The key concept is that we need to interact console, in other words,
			// we need to use system API, eg: On windows, we need to use PostMessage to
			// send message to the server console, and also we need to retrieve its output.
			// Since we **MUST** retrieve the true windows handle, it is very hard
			// to solve with just only go language (Even they do not provide a suitable
			// OS related package...).
			// 		Based on this, and also, this is just a CLI, we just need to send some
			// information to daemon, then everything leave to it to resolve.
			// 		Now, to solve this problem, we use C# to do the job: It provide
			// a very great method for us to interact with system API, and we don't have
			// to take too much care about it ----- At least on Windows.
			// 		Based on this, if we need to support Linux, I think C# could also help
			// us to do the job well.
			//
			// Now, according to the info above, procedures:
			// 		1. Check server whether exist.
			// 		2. Make final startup configuration.
			// 		3. Obtain the executable directory.
			// 	Noting that different game could have different executable files.
			// 	Currently we just write an inline override function to modify it.
			// 	Considering write the executable file path to the config file.
			// 		4. Locate the executable file.
			//	We need to locate the true executable file, since different game
			// will place their final exectuable file by their own method.
			// Considering write this to the config file.
			// 		5. Write JSON string, identify method to let daemon know
			// what we are doing, and write the information what we stated above.
			// 		6. Create a subprocess (and detach) (In golang we should call it goroutine).
			// 			6.1. We need to notify the coordinator that this server is starting.
			//			6.2. We startup the daemon, and let the daemon to startup the final
			// 		server application.

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

			// We export server.param for passing in 'Args' when we executing commands.
			cfg, err := cliconf.StartGameConfiguration(thisServer)

			if err != nil {
				fmt.Println("ERROR:", err)
				return
			}

			// By default, the executable file is under ${GSM_ROOT}/server/${ServerIndex}
			// Some of them may be under ${GSM_ROOT}/server/${ServerIndex}/bin, or whatever
			// If this case occured, will override via 'overrideServerExecutablePath'

			serverExecutableDir := fmt.Sprintf("%v/server/%v", os.Getenv("GSM_ROOT"), thisServer.ID)

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

			coordinator_cfg, err := cliconf.CoordinatorConfiguration()

			if err != nil {
				fmt.Println("ERROR:", err)
				return
			}

			// We need to pass parameters to gsm-server
			// That for convience, we send params via JSON
			// Leave everything to gsm-server to handle.
			// In this term, our task is completed.

			application_policy := make(map[string]interface{})
			application_policy["allow_update"] = cfg.GetBool("server.allow_update")
			application_policy["auto_update"] = cfg.GetBool("server.auto_update")
			application_policy["update_on_start"] = cfg.GetBool("server.update_on_start") // for the possibility that this server crashed.
			application_policy["auto_restart"] = cfg.GetBool("server.auto_restart")
			application_policy["restart_after_delay"] = cfg.GetInt("server.restart_after_delay")
			application_policy["max_restart_count"] = cfg.GetInt("server.max_restart_count")

			application := make(map[string]interface{})
			application["server_id"] = thisServer.ID
			application["game"] = cfg.GetString("server.game")
			application["dir"] = serverExecutableDir
			application["executable"] = serverExecutablePath
			application["args"] = strings.Join((cfg.GetStringSlice("server.param")), " ")
			application["policy"] = application_policy

			coordinator_policy := make(map[string]interface{})
			coordinator_policy["retry_count"] = cfg.GetUint("coordinator.retry_count")
			coordinator_policy["allow_retry_at_startup"] = cfg.GetBool("coordinator.allow_retry_at_startup")
			coordinator_policy["allow_reconnect_when_running"] = cfg.GetBool("coordinator.allow_reconnect_when_running")

			coordinator := make(map[string]interface{})
			coordinator["ip"] = coordinator_cfg.GetString("coordinator.ip")
			coordinator["port"] = coordinator_cfg.GetUint("coordinator.port")
			coordinator["policy"] = coordinator_policy

			data := make(map[string]interface{})
			data["application"] = application
			data["coordinator"] = coordinator

			msg := make(map[string]interface{})
			msg["command"] = "start"
			msg["message"] = data

			ret, err := json.Marshal(msg)
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
				Path:        wrapperEXEPath,
				Dir:         os.Getenv("GSM_ROOT") + "/bin",
				Args:        []string{wrapperEXEPath, string(ret)},
				SysProcAttr: sysprocattr.SysProcAttr(),
				Env:         os.Environ(),
				Stdin:       os.Stdin,
				Stdout:      os.Stdout,
				Stderr:      os.Stderr,
			}
			err = wrapperEXE.Run()
			if err != nil {
				fmt.Println("ERROR:", err)
				return
			}

			// That's why we want to build a coordinator....
			// We need a method to see all information...
			start := make(chan bool)
			go func() {
				wrapperEXE := &exec.Cmd{
					Path:        wrapperEXEPath,
					Dir:         os.Getenv("GSM_ROOT") + "/bin",
					Args:        []string{wrapperEXEPath, string(ret)},
					SysProcAttr: sysprocattr.SysProcAttr(),
				}
				// TODO: Notify the server that this server is starting.
				err = wrapperEXE.Start()
				if err != nil {
					fmt.Println("ERROR:", err)
					start <- false
					return
				}
				start <- true
			}()
			// <-start

		},
	}
	start.Flags().UintVar(&serverID, "server-id", 0, "Server ID to start")
	return start
}
