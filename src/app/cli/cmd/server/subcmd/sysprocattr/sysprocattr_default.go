//go:build !windows && !plan9
// +build !windows,!plan9

package sysprocattr

import "syscall"

func SysProcAttr() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{
		Setsid: true,
	}
}
