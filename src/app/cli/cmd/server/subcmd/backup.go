package subcmd

import (
	"fmt"
	"os"

	"github.com/laper32/regsm-console/src/lib/os/shutil"
	"github.com/spf13/cobra"
)

func InitBackupCMD() *cobra.Command {
	var serverID uint
	backup := &cobra.Command{
		Use: "backup",
		Run: func(cmd *cobra.Command, args []string) {
			// Backup servers is just simply copy and paste server related files.
			// I think we should not say anything about it.
			//
			// Step:
			// 	1. Get the server related directories.
			// 	2. Copy and paste.
			// btw, we should have 'restore' command to restore files in backup.
			// but not this version.
			// we do this in the future.

			makeDirectory := func(serverID uint) (string, string, string) {
				rootDirectory := os.Getenv("GSM_ROOT")
				thisServerDirectory := fmt.Sprintf("%v/server/%v", rootDirectory, serverID)
				thisConfigDirectory := fmt.Sprintf("%v/config/server/%v", rootDirectory, serverID)
				thisLogDirectory := fmt.Sprintf("%v/log/server/%v", rootDirectory, serverID)
				return thisServerDirectory, thisConfigDirectory, thisLogDirectory
			}
			serverDirectory, configDirectory, logDirectory := makeDirectory(serverID)
			backupDirectory := fmt.Sprintf("%v/backup/%v", os.Getenv("GSM_ROOT"), serverID)
			shutil.CopyDir(serverDirectory, fmt.Sprintf("%v/server", backupDirectory))
			shutil.CopyDir(configDirectory, fmt.Sprintf("%v/config", backupDirectory))
			shutil.CopyDir(logDirectory, fmt.Sprintf("%v/log", backupDirectory))
		},
	}
	backup.Flags().UintVar(&serverID, "server-id", 0, "")
	backup.MarkFlagRequired("server-id")
	return backup
}
