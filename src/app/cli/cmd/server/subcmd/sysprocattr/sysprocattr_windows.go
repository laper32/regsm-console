//go:build windows
// +build windows

package sysprocattr

import "syscall"

func SysProcAttr() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{
		// HideWindow: true,
	}
}
