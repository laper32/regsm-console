package windows

import (
	"os/exec"
	"syscall"
	"time"
)

type WinProcess struct {
	EXE              *exec.Cmd // use exec.cmd as process for simplify
	MainWindowHandle syscall.Handle
}

func (w *WinProcess) Start() error {
	err := w.EXE.Start()
	if err != nil {
		return err
	}
	time.Sleep(time.Millisecond * 200)
	w.MainWindowHandle = GetMainWindowHandle(w.EXE.Process.Pid)
	return nil
}
