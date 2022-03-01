package windows

import (
	"syscall"
	"time"
	"unsafe"
)

var (
	user32                     = syscall.NewLazyDLL("user32.dll")
	__GetWindow                = user32.NewProc("GetWindow")
	__ShowWindow               = user32.NewProc("ShowWindow")
	__EnumWindows              = user32.NewProc("EnumWindows")
	__PostMessage              = user32.NewProc("PostMessageW")
	__IsWindowVisible          = user32.NewProc("IsWindowVisible")
	__GetWindowThreadProcessId = user32.NewProc("GetWindowThreadProcessId")
)

func EnumWindows(enumFunc, lParam uintptr) error             { return _EnumWindows(enumFunc, lParam) }
func GetWindow(hwnd syscall.Handle, cmd uint) syscall.Handle { return _GetWindow(hwnd, cmd) }

func GetWindowThreadProcessId(hwnd syscall.Handle) (syscall.Handle, int) {
	return _GetWindowThreadProcessId(hwnd)
}

func IsWindowVisible(hwnd syscall.Handle) bool          { return _IsWindowVisible(hwnd) }
func ShowWindow(hwnd syscall.Handle, nCmdShow int) bool { return _ShowWindow(hwnd, nCmdShow) }

func PostMessage(hwnd syscall.Handle, msg uint32, wParam, lParam uintptr) bool {
	return _PostMessage(hwnd, msg, wParam, lParam)
}

// Be aware that the Windows API are all asynchonrous.
//
// What it means that, you must __PAUSE__ to ensure that that thread of input has completed.
//
// And all related to windows API using __time.Sleep()__ are all due to this issue.
func SendToConsole(hwnd syscall.Handle, message string) {
	PostMessage(hwnd, 0x0100, 13, 0)
	time.Sleep(200 * time.Millisecond)
	_byte, _ := syscall.UTF16FromString(message)
	for i := 0; i < len(_byte); i++ {
		if i > 0 && _byte[i] == _byte[i-1] {
			PostMessage(hwnd, 0x0100, 0, 0)
		}
		PostMessage(hwnd, 0x0102, uintptr(_byte[i]), 0)
	}
	time.Sleep(200 * time.Millisecond)
	PostMessage(hwnd, 0x0100, 13, 0)
}

// On Windows, not all servers creating method will call just only console APP.
// Some of them will make a wrapper, for example: hlds, srcds.
//
// To solve this issue, we need to find its main window handle, then send byte
// message to this handle, via windows API.
//
// aka, we need to get its main window handle.
func GetMainWindowHandle(pid int) syscall.Handle {
	var hwnd syscall.Handle
	cb := syscall.NewCallback(func(h syscall.Handle, p uintptr) uintptr {
		isMainWindow := func(_h syscall.Handle) bool {
			return GetWindow(_h, GW_OWNER) == syscall.Handle(0) && IsWindowVisible(_h)
		}
		_, pidFind := GetWindowThreadProcessId(h)
		if pidFind != pid || !isMainWindow(h) {
			return 1
		}
		hwnd = h
		return 0
	})
	EnumWindows(cb, 0)
	return hwnd
}

func _EnumWindows(enumFunc uintptr, lparam uintptr) (err error) {
	r1, _, e1 := __EnumWindows.Call(uintptr(enumFunc), uintptr(lparam), 0)
	if r1 == 0 {
		if e1 != nil {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}
func _GetWindow(hWnd syscall.Handle, cmd uint) syscall.Handle {
	ret, _, _ := __GetWindow.Call(uintptr(hWnd), uintptr(unsafe.Pointer(&cmd)))
	return syscall.Handle(ret)
}
func _GetWindowThreadProcessId(hwnd syscall.Handle) (syscall.Handle, int) {
	var processId int
	ret, _, _ := __GetWindowThreadProcessId.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&processId)))

	return syscall.Handle(ret), processId
}
func _IsWindowVisible(hwnd syscall.Handle) bool {
	ret, _, _ := __IsWindowVisible.Call(uintptr(hwnd))
	return ret != 0
}
func _PostMessage(hwnd syscall.Handle, msg uint32, wParam, lParam uintptr) bool {
	ret, _, _ := __PostMessage.Call(uintptr(hwnd), uintptr(msg), wParam, lParam)
	return ret != 0
}
func _ShowWindow(hwnd syscall.Handle, cmdshow int) bool {
	ret, _, _ := __ShowWindow.Call(uintptr(hwnd), uintptr(cmdshow))
	return ret != 0
}
