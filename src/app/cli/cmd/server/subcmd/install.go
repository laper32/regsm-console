package subcmd

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"

	cliconf "github.com/laper32/regsm-console/src/app/cli/conf"
	"github.com/laper32/regsm-console/src/app/cli/dpkg"
	"github.com/laper32/regsm-console/src/app/cli/misc"
	"github.com/laper32/regsm-console/src/lib/interact"
	"github.com/laper32/regsm-console/src/lib/os/path"
	"github.com/laper32/regsm-console/src/lib/os/shutil"
	"github.com/laper32/regsm-console/src/lib/status"
	"github.com/spf13/cobra"
)

func InitInstallCMD() *cobra.Command {
	// gsm server install [--game] [--symlink --symlink-server-id] [--import --import-from-dir]

	var (
		game            string
		isImport        bool
		importDir       string
		symlink         bool
		symlinkServerID uint
	)

	install := &cobra.Command{
		Use: "install",
		Run: func(cmd *cobra.Command, args []string) {
			/*
				1. Check symlink and import flag. Reject if both occur.
				2. Process symlink case
					2.1. Check symlink ID
					2.2. Find the root server for symlink, if have.
					2.3. Check whether deleted.
					2.4.
				3. Process import case
				4. Process none of them case.
			*/

			makeDirString := func(serverID uint) (string, string, string) {
				rootDirectory := os.Getenv("GSM_ROOT")
				thisServerDirectory := fmt.Sprintf("%v/server/%v", rootDirectory, serverID)
				thisConfigDirectory := fmt.Sprintf("%v/config/server/%v", rootDirectory, serverID)
				thisLogDirectory := fmt.Sprintf("%v/log/server/%v", rootDirectory, serverID)
				return thisServerDirectory, thisConfigDirectory, thisLogDirectory
			}

			// https://stackoverflow.com/questions/10303319/aligning-text-output-by-the-console
			generateInstallationInfo := func(key, value []string) {
				if len(key) != len(value) {
					return
				}
				// Noting here:
				// If we just let orderedKey=key
				// This is the fact that in cpp:
				// void* ptr_new = &ptr;
				// and you should also know that if we modify ptr_new
				// we will also modify *ptr, not satisfy our assumption.
				// Based on this, we need to make its copy, then implement everything on it.
				orderedKey := make([]string, len(key))
				copy(orderedKey, key)
				sort.Strings(orderedKey)

				padWidth := len(orderedKey[len(orderedKey)-1]) + 10
				for i := 0; i < len(key); i++ {
					paddedString := fmt.Sprintf("%-"+strconv.Itoa(padWidth)+"v", key[i])
					fmt.Println(paddedString + value[i])
				}
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

			// Can't let them both occur
			if symlink && isImport {
				fmt.Println(status.CLIInstallNotAllowedBothSetSymlinkAndInstall.WriteDetail(""))
				return
			}

			var thisServer *dpkg.ServerIdentity
			if symlink {

				// Check existance first.
				symlinkServer := dpkg.FindIdentifiedServer(symlinkServerID)
				if symlinkServer == nil {
					fmt.Println(status.CLIInstallSymlinkServerNotExist.WriteDetail(""))
					return
				}

				// Indeed, the intermidiate server chain with possible deleted.
				// But it does not affect the root result.
				// Root server not deleted=everything is OK.
				if symlinkServer.SymlinkServerID != 0 {
					symlinkServer = dpkg.FindRootSymlinkServer(symlinkServer.ID)
					if symlinkServer == nil {
						fmt.Println(status.CLIInstallRootSymlinkServerNotExist.WriteDetail(""))
						return
					}
				}

				if symlinkServer.Deleted {
					fmt.Println(status.CLIInstallSymlinkServerDeleted.WriteDetail(""))
					return
				}

				// The newly installed server distributed files folder string.
				generateServerDirectory := func(serverID uint) {

					thisServerDirectory, thisConfigDirectory, thisLogDirectory := makeDirString(serverID)
					// Root server dir.
					rootServerDirectory, _, _ := makeDirString(symlinkServer.ID)

					// 	If you want to implement error, you will get this massive of shit.
					// Or just os.Mkdir, balabala.
					// 	Why we just only link server files, but log and config still mkdir?
					// Well...logging, I think I don't have to say anything
					// 	But config, yes, we can say everything are same. But however, we also
					// need to modify something like: Port...
					// In this term, we have to mkdir rather create a symlink......
					// 	There still also has a solution to handle this case, but this is restricted
					// in commercial version.
					fmt.Printf("Creating log directory...")
					err := os.Mkdir(thisLogDirectory, os.ModePerm)
					if err != nil {
						fmt.Println(status.CLIInstallServerDirectoryAlreadyExist.WriteDetail(err))
						return
					}
					fmt.Println("OK")

					fmt.Printf("[%v -> %v] Creating symbolic link...", rootServerDirectory, thisServerDirectory)
					err = os.Symlink(rootServerDirectory, thisServerDirectory)
					if err != nil {
						fmt.Println(status.CLIInstallUnableToCreateSymlink.WriteDetail(err))
						return
					}
					fmt.Println("OK")

					fmt.Printf("Creating config directory...")
					err = os.Mkdir(thisConfigDirectory, os.ModePerm)
					if err != nil {
						fmt.Println(status.CLIInstallServerDirectoryAlreadyExist.WriteDetail(err))
						return
					}
					fmt.Println("OK")
				}

				// Check the info of reuse id
				// If we have found the id reuse, then we will modify the wrapper to this deleted server.

				canReuseID := func() bool {
					for i, this := range dpkg.ServerIdentityList {
						if this.Deleted {
							this.Deleted = false
							this.Game = symlinkServer.Game
							this.SymlinkServerID = symlinkServerID
							dpkg.ServerIdentityList[i] = this
							thisServer = &this
							return true
						}
					}
					return false
				}()
				if !canReuseID {
					thisServer = &dpkg.ServerIdentity{
						ID:              uint(len(dpkg.ServerIdentityList)) + 1,
						Game:            symlinkServer.Game,
						Deleted:         false,
						SymlinkServerID: symlinkServerID,
					}
				}

				if !misc.Agree && misc.Decline {
					return
				}

				symlinkServerDir, _, _ := makeDirString(thisServer.SymlinkServerID)
				serverDir, configDir, logDir := makeDirString(thisServer.ID)
				fmt.Println("Installation information")
				_key := []string{
					"Server ID:",
					"Game:",
					"Symbolic server ID:",
					"Symbolic server path:",
					"Server path:",
					"Server config path:",
					"Server logging path:",
				}

				_value := []string{
					fmt.Sprintf("%v", thisServer.ID),
					dpkg.FindGame(thisServer.Game).Name,
					fmt.Sprintf("%v", thisServer.SymlinkServerID),
					symlinkServerDir,
					serverDir,
					configDir,
					logDir,
				}

				generateInstallationInfo(_key, _value)

				if furtherActionNeeded() {
					fmt.Println("You are now going to proceed this installation.")
					result := interact.MakeConfirmation("Proceed?")
					if !result {
						return
					}
				}

				// We need to check symlink dir whether exists...
				if path.Exist(symlinkServerDir) {
					start := time.Now()
					if !canReuseID {
						dpkg.ServerIdentityList = append(dpkg.ServerIdentityList, *thisServer)
					}
					dpkg.UpdateServerIdentity()
					fmt.Println("Generating server directories...")
					generateServerDirectory(thisServer.ID)
					fmt.Println("OK")

					fmt.Printf("Writing default configuration...")
					cliconf.WriteDefaultGameConfig(thisServer.Game, configDir)
					fmt.Println("OK")

					elapsed := time.Since(start)
					fmt.Println("Installation complete. Time elapsed: ", elapsed)
				} else {
					fmt.Println("Unexpected error occured: symlink server dir has been deleted!")
					return
				}

				return
			}

			if isImport {
				// Always 0
				symlinkServerID = 0

				if len(importDir) == 0 {
					fmt.Println("Must explicitly declare the import directory.")
					return
				}

				if !path.Exist(importDir) {
					fmt.Println("Directory not found: \"", importDir, "\"")
					return
				}

				if result := dpkg.FindGame(game); result == nil {
					fmt.Println("ERROR: Can't find the game!")
					return
				}

				generateServerDirectory := func(serverID uint) {
					thisServerDirectory, thisConfigDirectory, thisLogDirectory := makeDirString(serverID)
					fmt.Printf("Creating log directory...")
					err := os.Mkdir(thisLogDirectory, os.ModePerm)
					if err != nil {
						fmt.Println(status.CLIInstallServerDirectoryAlreadyExist.WriteDetail(err))
						return
					}
					fmt.Println("OK")

					fmt.Printf("Creating server directory...")
					err = os.Mkdir(thisServerDirectory, os.ModePerm)
					if err != nil {
						fmt.Println(status.CLIInstallServerDirectoryAlreadyExist.WriteDetail(err))
						return
					}
					fmt.Println("OK")

					fmt.Printf("Creating config directory...")
					err = os.Mkdir(thisConfigDirectory, os.ModePerm)
					if err != nil {
						fmt.Println(status.CLIInstallServerDirectoryAlreadyExist.WriteDetail(err))
						return
					}
					fmt.Println("OK")
				}

				// Check the info of reuse id
				// If we have found the id reuse, then we will modify the wrapper to this deleted server.
				canReuseID := func() bool {
					for i, this := range dpkg.ServerIdentityList {
						if this.Deleted {
							this.Deleted = false
							this.Game = game
							this.SymlinkServerID = symlinkServerID
							dpkg.ServerIdentityList[i] = this
							thisServer = &this
							return true
						}
					}
					return false
				}()
				if !canReuseID {
					thisServer = &dpkg.ServerIdentity{
						ID:              uint(len(dpkg.ServerIdentityList)) + 1,
						Game:            game,
						Deleted:         false,
						SymlinkServerID: symlinkServerID,
					}
				}

				if !misc.Agree && misc.Decline {
					return
				}

				serverDir, configDir, logDir := makeDirString(thisServer.ID)
				fmt.Println("Installation information")
				_key := []string{
					"Server ID:",
					"Game:",
					"External server package directory:",
					"Server path:",
					"Server config path:",
					"Server logging path:",
				}

				_value := []string{
					fmt.Sprintf("%v", thisServer.ID),
					dpkg.FindGame(thisServer.Game).Name,
					importDir,
					serverDir,
					configDir,
					logDir,
				}

				generateInstallationInfo(_key, _value)

				if furtherActionNeeded() {
					fmt.Println("You are now going to proceed this installation.")
					result := interact.MakeConfirmation("Proceed?")
					if !result {
						return
					}
				}

				start := time.Now()
				if !canReuseID {
					dpkg.ServerIdentityList = append(dpkg.ServerIdentityList, *thisServer)
				}
				dpkg.UpdateServerIdentity()

				fmt.Printf("Generating server directories...")
				generateServerDirectory(thisServer.ID)
				fmt.Println("OK")

				fmt.Printf("Copying files...")
				shutil.CopyDir(importDir, serverDir)
				fmt.Println("OK")

				fmt.Printf("Writing default configurations...")
				cliconf.WriteDefaultGameConfig(game, configDir)
				fmt.Println("OK")

				elapsed := time.Since(start)
				fmt.Println("Installation complete. Time elapsed: ", elapsed)

				return
			}

			// Newly installed server, we set 0 to avoid mistyping
			symlinkServerID = 0
			fmt.Println("Newly installed server.")

			if len(game) == 0 {
				fmt.Println(status.CLIInstallExplicitlyDeclareWhichGameToInstall.WriteDetail(""))
				return
			}

			result := dpkg.FindGame(game)
			if result == nil {
				fmt.Println("Unable to find the game.")
				return
			}

			generateServerDirectory := func(serverID uint) {
				thisServerDirectory, thisConfigDirectory, thisLogDirectory := makeDirString(serverID)
				fmt.Println()
				fmt.Printf("Creating log directory...")
				err := os.Mkdir(thisLogDirectory, os.ModePerm)
				if err != nil {
					fmt.Println(status.CLIInstallServerDirectoryAlreadyExist.WriteDetail(err))
					return
				}
				fmt.Println("OK")

				fmt.Printf("Creating server directory...")
				err = os.Mkdir(thisServerDirectory, os.ModePerm)
				if err != nil {
					fmt.Println(status.CLIInstallServerDirectoryAlreadyExist.WriteDetail(err))
					return
				}
				fmt.Println("OK")

				fmt.Printf("Creating config directory...")
				err = os.Mkdir(thisConfigDirectory, os.ModePerm)
				if err != nil {
					fmt.Println(status.CLIInstallServerDirectoryAlreadyExist.WriteDetail(err))
					return
				}
				fmt.Println("OK")
			}

			// Check the info of reuse id
			// If we have found the id reuse, then we will modify the wrapper to this deleted server.
			canReuseID := func() bool {
				for i, this := range dpkg.ServerIdentityList {
					if this.Deleted {
						this.Deleted = false
						this.Game = game
						this.SymlinkServerID = symlinkServerID
						dpkg.ServerIdentityList[i] = this
						thisServer = &this
						return true
					}
				}
				return false
			}()
			if !canReuseID {
				thisServer = &dpkg.ServerIdentity{
					ID:              uint(len(dpkg.ServerIdentityList)) + 1,
					Game:            game,
					Deleted:         false,
					SymlinkServerID: symlinkServerID,
				}
			}

			if !misc.Agree && misc.Decline {
				return
			}

			serverDir, configDir, logDir := makeDirString(thisServer.ID)
			fmt.Println("Installation information")
			_key := []string{
				"Server ID:",
				"Game:",
				"Server path:",
				"Server config path:",
				"Server logging path:",
			}

			_value := []string{
				fmt.Sprintf("%v", thisServer.ID),
				dpkg.FindGame(thisServer.Game).Name,
				serverDir,
				configDir,
				logDir,
			}

			generateInstallationInfo(_key, _value)

			if furtherActionNeeded() {
				fmt.Println("You are now going to proceed this installation.")
				result := interact.MakeConfirmation("Proceed?")
				if !result {
					return
				}
			}

			start := time.Now()
			if !canReuseID {
				dpkg.ServerIdentityList = append(dpkg.ServerIdentityList, *thisServer)
			}
			dpkg.UpdateServerIdentity()

			fmt.Printf("Generating server directories...")
			generateServerDirectory(thisServer.ID)
			fmt.Println("OK")

			executeInstallation := func() {
				installVia := result.Specific["install_via"]
				if installVia == "steamcmd" {
					appid, modName, custom := int64(result.Specific["appid"].(float64)), "", ""

					if value, ok := result.Specific["mod"].(string); ok {
						modName = value
					}

					if value, ok := result.Specific["custom"]; ok {
						custom = value.(string)
					}

					var platformList []string
					for _, this := range result.Specific["platform"].([]interface{}) {
						platformList = append(platformList, this.(string))
					}
					fmt.Println(appid, modName, custom)

					dpkg.SteamCMDInstall(platformList, serverDir, appid, modName, true, custom)
					return
				}
			}

			fmt.Println("Installing server...")
			executeInstallation()
			fmt.Println("OK")

			fmt.Printf("Writing default configurations...")
			cliconf.WriteDefaultGameConfig(game, configDir)
			fmt.Println("OK")

			elapsed := time.Since(start)
			fmt.Println("Installation complete. Time elapsed: ", elapsed)

		},
	}
	install.Flags().StringVar(&game, "game", "", "Game to install. Must explicitly declared unless \"--symlink\" is set.")
	install.Flags().Lookup("game").NoOptDefVal = ""

	install.Flags().BoolVar(&isImport, "import", false, "Determine whether this installation is import.")
	install.Flags().Lookup("import").NoOptDefVal = "true"

	install.Flags().StringVar(&importDir, "import-dir", "", "The directory of external package.")
	install.Flags().Lookup("import-dir").NoOptDefVal = ""

	install.Flags().BoolVar(&symlink, "symlink", false, "Determine whether this game is symlink.")
	install.Flags().Lookup("symlink").NoOptDefVal = "true"

	install.Flags().UintVar(&symlinkServerID, "symlink-server-id", 0, "Symbolic link server id.")
	install.Flags().Lookup("symlink-server-id").NoOptDefVal = "0"

	return install
}
