package subcmd

import (
	"fmt"
	"os"

	"github.com/laper32/regsm-console/src/app/cli/dpkg"
	"github.com/laper32/regsm-console/src/app/cli/misc"
	"github.com/laper32/regsm-console/src/lib/interact"
	"github.com/laper32/regsm-console/src/lib/os/path"
	"github.com/laper32/regsm-console/src/lib/structs"
	"github.com/spf13/cobra"
)

func InitRemoveCMD() *cobra.Command {
	var (
		serverID      uint
		componentKeep []string
	)
	remove := &cobra.Command{
		Use: "remove",
		Run: func(cmd *cobra.Command, args []string) {
			// Removing server, is to remove server files, and update server info config(or database).
			// Basically, we won't keep anything. (But in my view, we should keep logging.)
			//
			// Steps:
			// 	1. Before removing files, we need to let user to do further confirmations.
			// 	2. Retrieve this server directories.
			// 	3. Check which component we want to keep.
			// 	4. Delete. Don't forget to exclude these components that asked to keep.

			// The cobra, if no default value, will throw error, and not execute anything.
			// Based on this, we dont have to do such thing

			deleteServer := func() {
				for i, content := range dpkg.ServerIdentityList {
					if content.ID == serverID {
						fmt.Println("Trigger removal process.")
						content.Deleted = true
						dpkg.ServerIdentityList[i] = content
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

			if !misc.Agree && misc.Decline {
				return
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

			fmt.Println("You want to keep following components: ", componentKeep)
			fmt.Println("Otherwise will be removed.")

			if furtherActionNeeded() {
				fmt.Println("You will going to proceed the removal.")
				result := interact.MakeConfirmation("Proceed?")
				if !result {
					return
				}
			}

			deleteServer()

			// 好tm蠢
			// 要不是现在想不到什么太好的办法...
			updateConfig := func() {
				dpkg.ServerIdentityMap = dpkg.ServerIdentityMap[0:0]
				for _, thisServer := range dpkg.ServerIdentityList {
					thisMap, _ := structs.ToMap(thisServer, "map")
					dpkg.ServerIdentityMap = append(dpkg.ServerIdentityMap, thisMap)
				}
				dpkg.ServerIdentityConfig.Set("server_identity", dpkg.ServerIdentityMap)

				dpkg.ServerIdentityConfig.WriteConfig()

			}
			updateConfig()

			fmt.Printf("server ID: %v has been removed. Component kept: %v. Otherwise are removed.\n", serverID, componentKeep)
		},
	}
	remove.Flags().UintVar(&serverID, "server-id", 0, "")
	remove.MarkFlagRequired("server-id")

	// Which component you want to keep? only accepts {log, config}. Default is {}.
	remove.Flags().StringSliceVar(&componentKeep, "component-keep", nil, "Which component you want to keep, blank means remove everything")

	return remove
}
