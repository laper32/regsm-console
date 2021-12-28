package subcmd

import (
	"fmt"
	"os"

	"github.com/laper32/regsm-console/src/app/cli/dpkg"
	"github.com/laper32/regsm-console/src/lib/conf"
	"github.com/spf13/cobra"
)

func InitUpdateCMD() *cobra.Command {
	var serverID uint
	var confirm bool
	update := &cobra.Command{
		Use: "update",
		Run: func(cmd *cobra.Command, args []string) {
			// Step:
			// 	1. Find this server from server_info.json
			// 	2. Check its relationship
			// 		1. If, this server does not related to any servers, then update it directly
			// 		2. If, this server does related to some servers, then we found its root server.
			// 	And we update this root server.
			// 		(Noting that in this term all servers except root server with update related configs
			// 		are all disabled.)
			// 	3. Update. (In fact the remaining are quite similar to installation.)

			thisServer := dpkg.FindServerInfoServer(serverID)
			if thisServer == nil {
				fmt.Printf("ERROR: Cannot found server %v.\n", serverID)
				return
			}

			makeDirectory := func(serverID uint) (string, string, string) {
				rootDirectory := os.Getenv("GSM_ROOT")
				thisServerDirectory := fmt.Sprintf("%v/server/%v", rootDirectory, serverID)
				thisConfigDirectory := fmt.Sprintf("%v/config/server/%v", rootDirectory, serverID)
				thisLogDirectory := fmt.Sprintf("%v/log/server/%v", rootDirectory, serverID)
				return thisServerDirectory, thisConfigDirectory, thisLogDirectory
			}

			startUpdate := func(serverDirectory string) {
				gameData := dpkg.FindGame(thisServer.Game)
				if gameData != nil {
					// Trust me, this part will become a massive of shit.
					installVia := gameData.Specific["install_via"]
					if installVia == "steamcmd" {
						appid, modName, custom := int64(gameData.Specific["appid"].(float64)), "", ""

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
				}
			}

			execute := func(serverID uint) {
				serverDirectory, configDirectory, _ := makeDirectory(serverID)

				cfg := conf.Load(&conf.Config{
					Name: "config",
					Type: "toml",
					Path: []string{configDirectory},
				})
				allowUpdate := cfg.GetBool("server.allow_update")
				if !allowUpdate {
					fmt.Println("ERROR: This server does not allow update.")
					return
				}

				startUpdate(serverDirectory)
			}

			// When this server is related to another server, all config field
			// about updating are all disabled.
			if thisServer.RelatedTo > -1 {
				rootServer := dpkg.FindRootRelatedServer(thisServer.ID)
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

	update.Flags().BoolVarP(&confirm, "yes", "y", false, "")
	update.Flags().MarkHidden("yes")

	return update
}
