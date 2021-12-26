package daemon

import (
	"log"
	"os"
	"os/exec"
	"time"
)

func startProc(args, env []string) (*exec.Cmd, error) {
	filename := "L" + time.Now().Format("20060102") + ".log"
	cmd := &exec.Cmd{
		Path:        args[0],
		Args:        args,
		Env:         env,
		SysProcAttr: NewSysProcAttr(),
	}
	stdout, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Println(os.Getegid(), ": Unexpected error occured when opening log file: ", err)
		return nil, err
	}
	cmd.Stderr = stdout
	cmd.Stdout = stdout

	err = cmd.Start()
	if err != nil {
		return nil, err
	}
	return cmd, nil
}
