package util

import (
	"fmt"
	"os"
	"runtime"

	"github.com/laper32/regsm-console/src/app/server/conf"
)

func CombineSeverPath(cfg *conf.Config) (exeDir string, exeName string) {
	exeDir = fmt.Sprintf("%v/server/%v", os.Getenv("GSM_ROOT"), cfg.Server.ID)
	switch cfg.Server.Game {
	case "csgo", "insurgency", "l4d", "l4d2", "css":
		if runtime.GOOS == "windows" {
			exeName = "srcds.exe"
		} else if runtime.GOOS == "linux" {
			exeName = "srcds_linux"
		} else {
			panic("Not supported")
		}
	case "cs1.6":
		if runtime.GOOS == "windows" {
			exeName = "hlds.exe"
		} else {
			panic("Not supported.")
		}
	default:
		panic(fmt.Sprintf("Unsupported game \"%v\"", cfg.Server.Game))
	}
	return
}
