package subcmd

import (
	"fmt"
	"net/url"
	"os"

	"github.com/gorilla/websocket"
	cliconf "github.com/laper32/regsm-console/src/app/cli/conf"
	"github.com/laper32/regsm-console/src/app/cli/dpkg"
	"github.com/laper32/regsm-console/src/app/cli/misc"
	"github.com/laper32/regsm-console/src/lib/conf"
	"github.com/laper32/regsm-console/src/lib/interact"
	"github.com/laper32/regsm-console/src/lib/log"
	"github.com/laper32/regsm-console/src/lib/status"
	"github.com/spf13/cobra"
)

func InitUpdateCMD() *cobra.Command {
	var serverID uint
	update := &cobra.Command{
		Use: "update",
		Run: func(cmd *cobra.Command, args []string) {
			// Procedure:
			// 1. Stop current running server (Gracefully).
			// 2. When server stopped, obtain current server directories(log, config, distributed, etc)
			// 3. Check whether allowed update, if allowed then continue, otherwise stop.
			// 4. Update server according the installation method.

			thisServer := dpkg.FindIdentifiedServer(serverID)
			if thisServer == nil {
				fmt.Printf("ERROR: Cannot found server %v.\n", serverID)
				return
			}

			if thisServer.Deleted {
				fmt.Printf("Server has been deleted.\n")
				return
			}

			cfg, err := cliconf.CoordinatorConfiguration()
			if err != nil {
				log.Error(err)
				return
			}

			makeDirectory := func(serverID uint) (string, string, string) {
				rootDirectory := os.Getenv("GSM_ROOT")
				thisServerDirectory := fmt.Sprintf("%v/server/%v", rootDirectory, serverID)
				thisConfigDirectory := fmt.Sprintf("%v/config/server/%v", rootDirectory, serverID)
				thisLogDirectory := fmt.Sprintf("%v/log/server/%v", rootDirectory, serverID)
				return thisServerDirectory, thisConfigDirectory, thisLogDirectory
			}
			serverDirectory, configDirectory, _ := makeDirectory(serverID)

			cfggame := conf.Load(&conf.Config{
				Name: "config",
				Type: "toml",
				Path: []string{configDirectory},
			})
			allowUpdate := cfggame.GetBool("server.allow_update")
			if !allowUpdate {
				fmt.Println("This server does not allow update.")
				return
			}

			startUpdate := func(serverDirectory string) {
				gameData := dpkg.FindGame(thisServer.Game)
				if gameData != nil {
					// Trust me, this part will become a massive of shit.
					installVia := gameData.Specific["install_via"]
					if installVia == "steamcmd" {
						appid, modName, custom := int64(gameData.Specific["appid"].(float64)), "", ""
						handleUpdate := func() {
							if value, ok := gameData.Specific["mod"].(string); ok {
								modName = value
							}

							if value, ok := gameData.Specific["custom"]; ok {
								custom = value.(string)
							}

							var platformList []string
							for _, this := range gameData.Specific["platform"].([]interface{}) {
								platformList = append(platformList, this.(string))
							}

							dpkg.SteamCMDInstall(platformList, serverDirectory, appid, modName, true, custom)
						}
						latestBuild := dpkg.CheckLatestBuild(appid)
						localBuild, err := dpkg.CheckLocalBuild(serverID, appid)
						if err != nil {
							fmt.Println("Seems that appmanifest file has been corrupted.")
							fmt.Println("Forcely update the server.")
							fmt.Println("If this message occurs multiple times, please contact for supportance.")
							fmt.Println("Error message:", err)
							handleUpdate()
							return
						}
						fmt.Println("Local build:", localBuild)
						fmt.Println("Latest build:", latestBuild)
						if latestBuild != localBuild {
							fmt.Println("Difference detected.")
							fmt.Println("Starting update.")
							handleUpdate()
							return
						}
						fmt.Println("Current build is latest. No further action needed.")
					}
				}
			}

			execute := func(serverID uint) {
				startUpdate(serverDirectory)
			}
			furtherActionNeeded := func() bool {
				if misc.Agree && !misc.Decline {
					return false
				} else if !misc.Agree && misc.Decline {
					return false
				} else {
					return true
				}
			}

			if !misc.Agree && misc.Decline {
				return
			}
			if furtherActionNeeded() {
				fmt.Println("Noting that updating server means that you may have to stop multiple servers.")
				fmt.Println("Make sure that you have known all possible consequences.")
				confirm := interact.MakeConfirmation("Proceed?")
				if !confirm {
					return
				}
			}
			// Noting that when updating server, all relavent server are required to shutdown
			// to update
			// especially you are running symlink
			chain := dpkg.GetServerChainByID(serverID)
			url := url.URL{Scheme: "ws", Host: fmt.Sprintf("%v:%v", cfg.GetString("coordinator.ip"), cfg.GetUint("coordinator.port"))}

			stopServer := func(v uint) {
				conn, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
				if err != nil {
					log.Error(err)
					return
				}
				retGram := &RetGram{
					Role:    misc.Role,
					Code:    status.ServerConnectedCoordinatorAndLoggingIn.ToInt(),
					Message: status.ServerConnectedCoordinatorAndLoggingIn.Message(),
					Detail:  map[string]interface{}{"server_id": v, "command": "update"},
				}
				err = conn.WriteJSON(&retGram)
				if err != nil {
					return
				}
				err = conn.ReadJSON(&retGram)
				if err != nil {
					return
				}
				conn.Close()
			}
			isStopped := func(v uint) error {
				conn, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
				if err != nil {
					log.Error(err)
					return err
				}
				retGram := &RetGram{
					Role:    misc.Role,
					Code:    status.ServerConnectedCoordinatorAndLoggingIn.ToInt(),
					Message: status.ServerConnectedCoordinatorAndLoggingIn.Message(),
					Detail:  map[string]interface{}{"server_id": v, "command": "update"},
				}
				err = conn.WriteJSON(&retGram)
				if err != nil {
					return err
				}
				err = conn.ReadJSON(&retGram)
				if err != nil {
					return err
				}
				if retGram.Code == status.CoordinatorServerOffline.ToInt() {
					return nil
				} else {
					// 有毛病吧 大小写还要你管了？
					return fmt.Errorf("re-send command required")
				}
			}

			for _, v := range chain {
				stopServer(v)
			}

			for _, v := range chain {
				for {
					err := isStopped(v)
					if err == nil {
						break
					} else {
						stopServer(v)
					}
				}
			}

			// When this server is related to another server, all config field
			// about updating are all disabled.
			if thisServer.SymlinkServerID > 0 {
				rootServer := dpkg.FindRootSymlinkServer(thisServer.ID)
				if rootServer == nil {
					fmt.Println("ERROR: we cannot found the root server.")
					return
				}

				if rootServer.Deleted {
					fmt.Println("Root server has been deleted.")
					return
				}

				execute(rootServer.ID)
			} else {
				execute(thisServer.ID)
			}
		},
	}
	update.Flags().UintVar(&serverID, "server-id", 0, "")
	update.MarkFlagRequired("server-id")

	return update
}
