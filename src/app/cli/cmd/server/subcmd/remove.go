package subcmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/laper32/regsm-console/src/app/cli/dpkg"
	"github.com/laper32/regsm-console/src/lib/os/path"
	"github.com/laper32/regsm-console/src/lib/structs"
	"github.com/spf13/cobra"
)

func InitRemoveCMD() *cobra.Command {
	var serverID uint
	var componentKeep []string
	var confirm bool
	remove := &cobra.Command{
		Use: "remove",
		Run: func(cmd *cobra.Command, args []string) {
			// The cobra, if no default value, will throw error, and not execute anything.
			// Based on this, we dont have to do such thing

			deleteServer := func() {
				for i, content := range dpkg.ServerInfoList {
					if content.ID == serverID {
						fmt.Println("Trigger removal process.")
						content.Deleted = true
						dpkg.ServerInfoList[i] = content
					}
				}

				makeDirectory := func(serverID uint) (string, string, string) {
					rootDirectory := os.Getenv("GSM_ROOT")
					thisServerDirectory := fmt.Sprintf("%v/server/%v", rootDirectory, serverID)
					thisConfigDirectory := fmt.Sprintf("%v/config/server/%v", rootDirectory, serverID)
					thisLogDirectory := fmt.Sprintf("%v/log/server/%v", rootDirectory, serverID)
					return thisServerDirectory, thisConfigDirectory, thisLogDirectory
				}

				// The logging files will keep for possible use in the future.
				// Otherwise will be deleted immediately.
				serverDirectory, configDirectory, logDirectory := makeDirectory(serverID)

				if len(componentKeep) == 0 {
					if path.Exist(serverDirectory) {
						os.RemoveAll(serverDirectory)
					}

					if path.Exist(configDirectory) {
						os.RemoveAll(configDirectory)
					}

					if path.Exist(logDirectory) {
						os.RemoveAll(logDirectory)
					}

				} else {
					if path.Exist(serverDirectory) {
						os.RemoveAll(serverDirectory)
					}

					for _, content := range componentKeep {
						if content == "log" {
							os.RemoveAll(logDirectory)
						} else if content == "config" {
							os.RemoveAll(configDirectory)
						} else {
							fmt.Println("Unknown component:", content)
						}
					}
				}
			}

			makeConfirmation := func() {
				if !confirm {
					fmt.Println("WARNING: You are now removing a server!")
					fmt.Printf("The server(ID=%v)'s file will be deleted permanently, but these components listed below will be kept: %v\n", serverID, componentKeep)

					// TODO: When we input something like 'CTRL+C', 'CTRL+Z', etc, this will not behave expectly...
					// Perhaps we need to make detiled implementation.
					reader := bufio.NewScanner(os.Stdin)
					for {
						fmt.Printf("Are you really sure you want to remove this server? [Y/N] (case insensitively): ")
						reader.Scan()
						text := reader.Text()
						if strings.ToLower(text) == "y" {
							confirm = true
							break
						} else if strings.ToLower(text) == "n" {
							confirm = false
							break
						} else {
							fmt.Println("Unknown input. Your answer must either Y or N, which is case insensitively.")
						}
					}
				}
			}

			makeConfirmation()

			if confirm {
				deleteServer()

				// 好tm蠢
				// 要不是现在想不到什么太好的办法...
				updateConfig := func() {
					dpkg.ServerInfoMap = dpkg.ServerInfoMap[0:0]
					for _, thisServer := range dpkg.ServerInfoList {
						thisMap, _ := structs.ToMap(thisServer, "map")
						dpkg.ServerInfoMap = append(dpkg.ServerInfoMap, thisMap)
					}
					dpkg.ServerInfoConfig.Set("server_info", dpkg.ServerInfoMap)

					dpkg.ServerInfoConfig.WriteConfig()
				}
				updateConfig()

				fmt.Printf("server ID: %v has been removed. Component kept: %v. Otherwise are removed.\n", serverID, componentKeep)
			}
		},
	}
	remove.Flags().UintVar(&serverID, "server-id", 0, "")
	remove.MarkFlagRequired("server-id")

	// Which component you want to keep? only accepts {log, config}. Default is {}.
	remove.Flags().StringSliceVar(&componentKeep, "component-keep", nil, "Which component you want to keep, blank means remove everything")
	remove.Flags().BoolVarP(&confirm, "yes", "y", false, "")

	return remove
}
