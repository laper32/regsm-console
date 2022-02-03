// windows: go build -o gsm.exe
// linux: go build -o gsm

package main

import (
	"os"
	stdstring "strings"

	"github.com/laper32/regsm-console/src/app/cli/cmd"
	"github.com/laper32/regsm-console/src/app/cli/conf"
	"github.com/laper32/regsm-console/src/app/cli/dpkg"
	"github.com/laper32/regsm-console/src/lib/log"
	"github.com/laper32/regsm-console/src/lib/os/path"
	libstring "github.com/laper32/regsm-console/src/lib/strings"
)

// check readme for details
func initializeBaseDirectory() {
	// When we are running gsm.exe, we are at ${GSMDir}/bin
	// In order to get the root dir, we need to get its parent dir.
	dir := os.Getenv("GSM_ROOT")

	writeBackupDir := func() {
		if !path.Exist(dir + "/backup") {
			os.Mkdir(dir+"/backup", os.ModePerm)
		}
	}

	writeConfigDir := func() {
		if !path.Exist(dir + "/config") {
			os.Mkdir(dir+"/config", os.ModePerm)
		}
		if !path.Exist(dir + "/config/gsm") {
			os.Mkdir(dir+"/config/gsm", os.ModePerm)
		}
		if !path.Exist(dir + "/config/server") {
			os.Mkdir(dir+"/config/server", os.ModePerm)
		}
	}

	writeLogDir := func() {
		if !path.Exist(dir + "/log") {
			os.Mkdir(dir+"/log", os.ModePerm)
		}
		if !path.Exist(dir + "/log/gsm") {
			os.Mkdir(dir+"/log/gsm", os.ModePerm)
		}
		if !path.Exist(dir + "/log/server") {
			os.Mkdir(dir+"/log/server", os.ModePerm)
		}
	}

	writeServerDir := func() {
		if !path.Exist(dir + "/server") {
			os.Mkdir(dir+"/server", os.ModePerm)
		}
	}

	writeDir := func() {
		writeBackupDir()
		writeConfigDir()
		writeLogDir()
		writeServerDir()
	}

	writeDir()
}

func initializeEnv() {
	// The executable file is at ${ROOT_DIR}/bin
	// that we need to get its parent dir to retrieve the root dir.
	// We can't do it directly.
	binDirectory, rootDirectory := func() (string, string) {
		_binDir, _ := os.Getwd()
		_binDir = stdstring.Replace(_binDir, "\\", "/", -1)
		_rootDir := libstring.Subtract(_binDir, 0, stdstring.LastIndex(_binDir, "/"))
		return _binDir, _rootDir
	}()

	os.Setenv("GSM_ROOT", rootDirectory)
	os.Setenv("GSM_PATH", binDirectory)
}

func initializeModule() {

	cfg := conf.Init()

	log.Init(cfg.Log)
}

func initMisc() {
	dpkg.InitAvailableGameData()
	dpkg.InitServerIdentity()
}

func initAll() {
	initializeEnv()
	initializeModule()
	initializeBaseDirectory()
	initMisc()
}

func runCommand() {
	cmd.InitCMD().Execute()
}

func start() {
	initAll()
	runCommand()
}

func main() {

	start()
}
