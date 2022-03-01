// Different game has different method
// We are now only focusing srcds and hlds

//go:build windows
// +build windows

package oswrap

import (
	"syscall"

	"github.com/laper32/regsm-console/src/lib/sys/windows"
)

func ShowWindow(hwnd syscall.Handle) { windows.ShowWindow(hwnd, windows.SW_NORMAL) }

func HideWindow(hwnd syscall.Handle) { windows.ShowWindow(hwnd, windows.SW_HIDE) }

func SendToConsole(hwnd syscall.Handle, message string) { windows.SendToConsole(hwnd, message) }
