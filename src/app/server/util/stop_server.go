package util

import (
	"fmt"
	"time"

	"github.com/laper32/regsm-console/src/app/server/conf"
	"github.com/laper32/regsm-console/src/app/server/entity"
	"github.com/laper32/regsm-console/src/lib/sys/windows"
)

func SendToConsole(game, command string) error {
	switch game {
	case "cs1.6", "csgo":
		windows.SendToConsole(entity.Proc.MainWindowHandle, command)
		return nil
	default:
		return fmt.Errorf(fmt.Sprintf("Unsupported game \"%v\"", game))
	}
}

func ForceStopServer(cfg *conf.Config) {
	if entity.Proc.EXE.ProcessState != nil && entity.Proc.EXE.ProcessState.Exited() {
		return
	}
	switch cfg.Server.Game {
	case "cs1.6", "csgo":
		windows.SendToConsole(entity.Proc.MainWindowHandle, "quit")
	default:
		panic("You should not be there.")
	}
	count := 0
	for entity.Proc.EXE.ProcessState == nil {
		if count >= 5 {
			entity.Proc.EXE.Process.Kill()
			break
		}
		time.Sleep(1 * time.Second)
		count++
	}
	// var done bool
	// for {
	// 	if entity.Proc.EXE.ProcessState != nil && entity.Proc.EXE.ProcessState.Exited() {
	// 		done = true
	// 		break
	// 	}
	// 	fmt.Println("Counter:", count)
	// 	if count > 5 {
	// 		done = false
	// 		break
	// 	}
	// 	time.Sleep(1 * time.Second)
	// 	count++
	// }
	// if !done {
	// 	entity.Proc.EXE.Process.Kill()
	// }
}

func KillServer() {

}
