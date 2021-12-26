package subcmd

import (
	"fmt"
	"os"

	cliconf "github.com/laper32/regsm-console/src/cli/conf"
	"github.com/laper32/regsm-console/src/cli/dpkg"
	"github.com/laper32/regsm-console/src/lib/structs"
	"github.com/spf13/cobra"
)

func InitInstallCMD() *cobra.Command {
	var (
		game            string
		importFromExist bool
		importFromDir   string
		relatedTo       int
	)

	install := &cobra.Command{
		Use: "install",
		Run: func(cmd *cobra.Command, args []string) {
			// 	In order to install a new server, we need to make a basic identification, which
			// has been shown above.
			// 	Here, the process
			// 		1. Check the game from the database, if found, then continue, otherwise interrupt.
			// 	We don't have to check whether has been input something, because we have been marked it
			// 	as required.
			// 		2. Check whether the game is imported. If this game is imported, then we just
			// 	copy and paste from the external directory.
			// 	Noting that once this server is imported from outside, we cannot link it to any existed
			// 	servers!
			// 	Also, don't forget to check the directory whether exist!
			// 		3. Once this server is related to an existing server, we just simply create a symbolic
			// 	link of them.
			// 	But keep in mind that if you want to 'link' a server, you need to link the 'root' server,
			// 	despite this server you want to link with 'probably' that not the root server (thinking
			// 	it as a linked list).
			// 		Or, if this server not related to any game, you just need to install from a downloader,
			// 	eg: steamcmd.

			if result := dpkg.FindGame(game); result == nil {
				fmt.Println("ERROR: This game does not found in database: ", game)
				return
			}

			if importFromExist {
				if importFromDir == "" {
					fmt.Println("ERROR: You must explicitly declare where to import the game!")
					return
				}
			}

			makeDirectory := func(serverID uint) (string, string, string) {
				rootDirectory := os.Getenv("GSM_ROOT")
				thisServerDirectory := fmt.Sprintf("%v/server/%v", rootDirectory, serverID)
				thisConfigDirectory := fmt.Sprintf("%v/config/server/%v", rootDirectory, serverID)
				thisLogDirectory := fmt.Sprintf("%v/log/server/%v", rootDirectory, serverID)
				return thisServerDirectory, thisConfigDirectory, thisLogDirectory
			}

			updateConfig := func() {
				dpkg.ServerInfoMap = dpkg.ServerInfoMap[0:0]
				for _, thisServer := range dpkg.ServerInfoList {
					thisMap, _ := structs.ToMap(thisServer, "map")
					dpkg.ServerInfoMap = append(dpkg.ServerInfoMap, thisMap)
				}
				dpkg.ServerInfoConfig.Set("server_info", dpkg.ServerInfoMap)
				dpkg.ServerInfoConfig.WriteConfig()
			}

			var serverWrapper *dpkg.ServerInfo
			if relatedTo > -1 {
				// 	When we know thi server is related to an existing server, waow, very nice news to us. This
				// means that we just need to create a simple symbolic link, rather than a planty of massive of
				// shit.
				// 	Based on this, we just need to do these steps:
				// 	1. Check whether this server is imported from an existing server, if so, we should interrupt
				// immediately.
				// 	2. Check this related server, or in other words, parent server, whether exist.
				// 	3. Check whether the game we want to install matches the parent server.
				// 	4. Check server info, if we have deleted server, then use their ID, otherwise insert a new one.
				// 	Don't forget generate directory for this server!
				// 	5. Update config.

				if importFromExist {
					fmt.Println("ERROR: Your game server is related to an existing server, that you CANNOT enable 'import-from-exist'!")
					return
				}

				// Check the related server whether exists
				// In this term, the wrapper server ptr is this 'related to' server.
				serverExist := func() bool {
					for _, this := range dpkg.ServerInfoList {
						if relatedTo == int(this.ID) {
							serverWrapper = &this
							return true
						}
					}
					return false
				}()

				if !serverExist {
					fmt.Println("ERROR: The server what is related to does not exist.")
					return
				}

				// To check this newly installed game is same as the related server game, avoid embrassment.
				serverGameWhetherEqual := func() bool {
					for _, this := range dpkg.ServerInfoList {
						if game != this.Game {
							return false
						}
					}
					return true
				}()

				if !serverGameWhetherEqual {
					fmt.Printf("ERROR: Game matching failed. Install game: %v, related server game: %v", game, serverWrapper.Game)
					return
				}

				rootServer := dpkg.FindRootRelatedServer(serverWrapper.ID)

				// The server ID input is the newly installed server ID
				generateServerDirectory := func(serverID uint) {
					// The newly installed server distributed files folder string.
					thisServerDirectory, thisConfigDirectory, thisLogDirectory := makeDirectory(serverID)

					// Root server dir.
					rootServerDirectory, _, _ := makeDirectory(rootServer.ID)

					// 	If you want to implement error, you will get this massive of shit.
					// Or just os.Mkdir, balabala.
					// 	Why we just only link server files, but log and config still mkdir?
					// Well...logging, I think I don't have to say anything
					// 	But config, yes, we can say everything are same. But however, we also
					// need to modify something like: Port...
					// In this term, we have to mkdir rather create a symlink......
					// 	There still also has a solution to handle this case, but this is restricted
					// in commercial version.
					err := os.Mkdir(thisLogDirectory, os.ModePerm)
					if err != nil {
						fmt.Println("ERROR:", err)
						return
					}
					err = os.Symlink(rootServerDirectory, thisServerDirectory)
					if err != nil {
						fmt.Println("ERROR:", err)
						return
					}
					err = os.Mkdir(thisConfigDirectory, os.ModePerm)
					if err != nil {
						fmt.Println("ERROR:", err)
						return
					}
				}

				// Check the info of reuse id
				// If we have found the id reuse, then we will modify the wrapper to this deleted server.
				hasReuseID := func() bool {
					for i, this := range dpkg.ServerInfoList {
						if this.Deleted {
							this.Deleted = false
							dpkg.ServerInfoList[i] = this
							serverWrapper = &this
							generateServerDirectory(this.ID)
							return true
						}
					}
					return false
				}()

				// Copy, paste and modify of code below
				// Because the case is slightly different

				pushNewServer := func() {
					thisServer := &dpkg.ServerInfo{
						ID:        uint(len(dpkg.ServerInfoList)) + 1,
						Game:      game,
						Deleted:   false,
						RelatedTo: relatedTo,
					}
					generateServerDirectory(thisServer.ID)
					dpkg.ServerInfoList = append(dpkg.ServerInfoList, *thisServer)
				}

				if hasReuseID {
					updateConfig()
				} else {
					pushNewServer()
					updateConfig()
				}
			} else {

				// Step
				// 	1. Read configuration file, to check installed server.
				// 	2. Identify deleted server.
				// 		2.1. If, we have found an index of server what have been deleted, then
				// 		we relocate this index to this newly-installed server.
				// 		2.2 Otherwise, we increment the index.
				// 	3. Generate server folder, config folder, and log folder.

				generateServerDirectory := func(serverID uint) {
					thisServerDirectory, thisConfigDirectory, thisLogDirectory := makeDirectory(serverID)
					os.Mkdir(thisServerDirectory, os.ModePerm)
					os.Mkdir(thisConfigDirectory, os.ModePerm)
					os.Mkdir(thisLogDirectory, os.ModePerm)
				}

				// This forloop is aiming to re-use all deleted server.

				hasReuseID := func() bool {
					for i, content := range dpkg.ServerInfoList {
						if content.Deleted {
							content.Deleted = false
							dpkg.ServerInfoList[i] = content
							serverWrapper = &content
							generateServerDirectory(content.ID)
							return true
						}
					}
					return false
				}()

				pushNewServer := func() {
					thisServer := &dpkg.ServerInfo{
						ID:        uint(len(dpkg.ServerInfoList)) + 1,
						Game:      game,
						Deleted:   false,
						RelatedTo: relatedTo,
					}
					generateServerDirectory(thisServer.ID)
					serverWrapper = thisServer
					dpkg.ServerInfoList = append(dpkg.ServerInfoList, *thisServer)
				}

				if hasReuseID {
					updateConfig()
				} else {
					pushNewServer()
					updateConfig()
				}

				thisServerDirectory, thisConfigDirectory, _ := makeDirectory(serverWrapper.ID)
				fmt.Println(thisServerDirectory)
				cliconf.WriteDefaultGameConfig(serverWrapper.Game, thisConfigDirectory)

				// Then, check the installation requirement, and install.
				execute := func() {
					gameData := dpkg.FindGame(game)
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
							fmt.Println(appid, modName, custom)

							dpkg.SteamCMDInstall(platformList, thisServerDirectory, appid, modName, true, custom)

							return
						}
					}
				}
				execute()
			}
		},
	}

	install.Flags().StringVar(&game, "game", "", "Game to install")
	install.MarkFlagRequired("game")

	install.Flags().BoolVar(&importFromExist, "import-from-exist", false, "Install server from an existed local server package")
	install.Flags().Lookup("import-from-exist").NoOptDefVal = "false"

	install.Flags().StringVar(&importFromDir, "import-from-dir", "", "Where to import")
	install.Flags().Lookup("import-from-dir").NoOptDefVal = ""

	install.Flags().IntVar(&relatedTo, "related-to", -1, "Related to")
	install.Flags().Lookup("related-to").NoOptDefVal = "-1"

	return install
}
