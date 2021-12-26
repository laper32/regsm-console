package dpkg

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/laper32/regsm-console/src/lib/archive/zip"
	libnet "github.com/laper32/regsm-console/src/lib/net"
	"github.com/laper32/regsm-console/src/lib/os/path"
)

// On Windows, the SteamCMD will be installed at ${GSMDir}/bin/steamcmd
// On Linux/Mac, will check dpkg to validate whether it is installed.

//lint:ignore U1000 used in template.
func steamCMDParamList(steamCMDExecutable, installDir string, appID int64, modName string, validate bool, custom string) []string {
	ret := []string{}
	// Remember that the exec.Cmd is actually called as: ${ExecutablePath} ${ParamPack, split as ' '}
	// That the first param will always be the path.
	// To resolve this issue, we need to provide an extra parameter passin to let the path as first.
	ret = append(ret, steamCMDExecutable)

	// Login as anonymous
	// I don't really know why, but WindowsGSM why also provide login as non-anonymous?
	// Over-designed, buddy.
	ret = append(ret, "+force_install_dir \""+installDir+"\"")
	ret = append(ret, "+login anonymous")
	//sb.Append(!string.IsNullOrWhiteSpace(modName) ? $" +app_set_config {appId} mod \"{modName}\"" : string.Empty);
	modName = strings.TrimSpace(modName)
	if modName != "" {
		str := fmt.Sprintf("+app_set_config %v mod \"%v\"", appID, modName)
		ret = append(ret, str)
	}

	// hl requires 4 more times.
	// For this, we recommend you using 'import' instead of 'install'
	if appID == 90 {
		for i := 0; i < 4; i++ {
			str := fmt.Sprintf("+app_update %v", appID)
			if validate {
				str += " validate"
			} else {
				str += ""
			}
			ret = append(ret, str)
		}
	} else {
		str := fmt.Sprintf("+app_update %v", appID)
		str += " " + strings.TrimSpace(custom)
		ret = append(ret, str)
	}
	// Last term we need to exit
	ret = append(ret, "+quit")
	return ret
}

// Check documentation above for more details.
//lint:ignore U1000 used in template.
func windowsCheckSteamCMD(steamCMDDirectory string, steamCMDExecutable string) {
	// If steamCMD does not exist, then we create.
	if path.Exist(steamCMDDirectory) {
		os.Mkdir(steamCMDDirectory, os.ModePerm)
	}

	// if steamcmd.exe not exist
	// Download via steam CDN or somewhere we are self-hosting.
	// But you must sure that if you want to download via steam cdn, the hosts has been configured
	// so that you can connect to it.
	if !path.Exist(steamCMDExecutable) {
		// steamcmd.zip, should I say something?
		steamCMDZipFile := steamCMDDirectory + "/steamcmd.zip"
		// Then, download
		err := libnet.DownloadFile(steamCMDZipFile, "https://steamcdn-a.akamaihd.net/client/installer/steamcmd.zip")
		if err != nil {
			log.Fatal("Unable to download steamCMD: ", err)
			return
		}
		zip.Unzip(steamCMDZipFile, steamCMDDirectory)
		os.Remove(steamCMDZipFile)
	}
}

func windowsInstallation(serverDirectory string, appID int64, modName string, validate bool, custom string) {
	binDirectory := os.Getenv("GSM_PATH")

	steamCMDDirectory := binDirectory + "/steamcmd"
	steamCMDExecutable := steamCMDDirectory + "/steamcmd.exe"
	windowsCheckSteamCMD(steamCMDDirectory, steamCMDExecutable)

	command := steamCMDParamList(steamCMDExecutable, serverDirectory, appID, modName, validate, custom)
	fmt.Println(command)
	cmd := &exec.Cmd{
		Path: steamCMDExecutable,
		Args: command,
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func SteamCMDInstall(platform []string, serverDirectory string, appID int64, modName string, validate bool, custom string) {
	checkPlatform := func() bool {
		for _, this := range platform {
			if runtime.GOOS == this {
				return true
			}
		}
		return false
	}()
	if !checkPlatform {
		log.Fatalln("We cannot provide any installation because this game does not support your platform.")
		return
	}

	if runtime.GOOS == "windows" {
		windowsInstallation(serverDirectory, appID, modName, validate, custom)
		return
	}

	if runtime.GOOS == "linux" {
		// Check: Ubuntu/Centos/Arch/etc

		return
	}
}
