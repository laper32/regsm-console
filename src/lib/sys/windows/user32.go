package windows

import (
	"syscall"
	"unsafe"
)

// TRICK: char[] -> const char*
// Knowing that we can pass in char[] to const char* directly.
// (char[] is a sequence of character, but it can decay to pointer.)
// https://stackoverflow.com/questions/60738379/why-does-const-char-get-converted-to-const-char-all-the-time
// Now, based on this, we have know that the original function convention is:
// int MessageBoxA(HWND hWnd, LPCSTR lpText, LPCSTR lpCaption, UINT uType)
// where HWND, LPCSTR, UINT denotes unsigned long*, const char*, unsigned int, respectively.
// Now according to the previous analysis we stated, we can treat go's string as char[]
// And then we pass it in.
// That is, string -> *byte

//		_EnumWindows(lpEnumFunc WNDENUMPROC, lParam HWND) bool = user32.EnumWindows
//		_EnumChildWindows(hWndParent HWND, lpEnumFunc WNDENUMPROC, lParam LPARAM) bool = user32.EnumChildWindows
//sys	_FindWindowA(classname *byte, windowname *byte) (ret HWND, err error) = user32.FindWindowA
//sys	_FindWindowW(classname *uint16, windowname *uint16) (ret HWND, err error) = user32.FindWindowW
//sys	_FindWindoWExA(hWndParent HWND, hWndChildAfter HWND, lpszClass *byte, lpszWindow *byte) (ret HWND, err error) = user32.FindWindowExA
//sys	_FindWindoWExW(hWndParent HWND, hWndChildAfter HWND, lpszClass *uint16, lpszWindow *uint16) (ret HWND, err error) = user32.FindWindowExW
//sys	_GetClassNameA(hWnd HWND, lpClassName *byte, nMaxCount int) (len int32, err error) = user32.GetClassNameA
//sys	_GetClassNameW(hWnd HWND, lpClassName *uint16, nMaxCount int) (len int32, err error) = user32.GetClassNameW
//sys	_GetWindowTextA(hwnd HWND, str *byte, maxCount int32) (len int32, err error) = user32.GetWindowTextA
//sys	_GetWindowTextW(hwnd HWND, str *uint16, maxCount int32) (len int32, err error) = user32.GetWindowTextW
//sys	_GetWindowTextLength(hwnd HWND) (ret int32, err error) = user32.GetWindowTextLength
//sys	_GetWindow(hWnd HWND, uCmd uint) (ret HWND, err error) = user32.GetWindow
//sys	_GetWindowThreadProcessId(hwnd HWND, pid *uint32) (tid uint32, err error) = user32.GetWindowThreadProcessId
//sys	_GetTopWindow(hWnd HWND) = (ret HWND, err error) = user32.GetTopWindow
//sys	_IsWindowVisible(hwnd HWND) (ret BOOL, err error) = user32.IsWindowVisible
//sys	_MessageBoxA(hwnd HWND, text *byte, caption *byte, boxtype uint32) (ret int32, err error) = user32.MessageBoxA
//sys	_MessageBoxW(hwnd HWND, text *uint16, caption *uint16, boxtype uint32) (ret int32, err error) = user32.MessageBoxW
//sys	_PostMessageA(hwnd HWND, msg uint32, wParam WPARAM, lParam LPARAM) (ret BOOL, err error) = user32.PostMessageA
//sys	_PostMessageW(hwnd HWND, msg uint32, wParam WPARAM, lParam LPARAM) (ret BOOL, err error) = user32.PostMessageW
//sys	_ShowWindow(hwnd HWND, cmdshow int) (ret BOOL, err error) = user32.ShowWindow
//sys	_SetForegroundWindow(hWnd HWND) = (ret BOOL, err error) = user32.SetForegroundWindow
//sys	_SendMessageA(hwnd HWND, msg uint32, wParam WPARAM, lParam LPARAM) (ret LRESULT, err error) = user32.SendMessageA
//sys	_SendMessageW(hwnd HWND, msg uint32, wParam WPARAM, lParam LPARAM) (ret LRESULT, err error) = user32.SendMessageW

var (
	__EnumWindows      = user32.NewProc("EnumWindows")
	__EnumChildWindows = user32.NewProc("EnumChildWindows")
)

func EnumWindows(lpEnumFunc WNDENUMPROC, lParam HWND) bool {
	return _EnumWindows(lpEnumFunc, lParam)
}

func _EnumWindows(lpEnumFunc WNDENUMPROC, lParam HWND) bool {
	ret, _, _ := __EnumWindows.Call(uintptr(unsafe.Pointer(&lpEnumFunc)), uintptr(lParam))
	return ret != 0
}

func EnumChildWindows(hWndParent HWND, lpEnumFunc WNDENUMPROC, lParam LPARAM) bool {
	return _EnumChildWindows(hWndParent, lpEnumFunc, lParam)
}

func _EnumChildWindows(hWndParent HWND, lpEnumFunc WNDENUMPROC, lParam LPARAM) bool {
	ret, _, _ := __EnumChildWindows.Call(uintptr(hWndParent), uintptr(syscall.NewCallback(lpEnumFunc)), uintptr(lParam))
	return ret != 0
}

func MessageBox(hwnd HWND, text, caption string, boxtype uint32) int32 {
	var ret int32
	if isUnicode() {
		ret = MessageBoxW(hwnd, text, caption, boxtype)
	} else {
		ret = MessageBoxA(hwnd, text, caption, boxtype)
	}
	return ret
}

func MessageBoxA(hwnd HWND, text, caption string, flags uint32) int32 {
	textPtr, _ := syscall.BytePtrFromString(text)
	captionPtr, _ := syscall.BytePtrFromString(caption)
	ret, _ := _MessageBoxA(hwnd, textPtr, captionPtr, flags)
	return ret
}

func MessageBoxW(hwnd HWND, text, caption string, flags uint32) int32 {
	textPtr, _ := syscall.UTF16PtrFromString(text)
	captionPtr, _ := syscall.UTF16PtrFromString(caption)
	ret, _ := _MessageBoxW(hwnd, textPtr, captionPtr, flags)
	return ret
}

func PostMessage(hwnd HWND, msg uint32, wParam WPARAM, lParam LPARAM) (ret bool) {
	if isUnicode() {
		ret = PostMessageW(hwnd, msg, wParam, lParam)
	} else {
		ret = PostMessageA(hwnd, msg, wParam, lParam)
	}
	return ret
}

func PostMessageA(hwnd HWND, msg uint32, wParam WPARAM, lParam LPARAM) bool {
	val, _ := _PostMessageA(hwnd, msg, wParam, lParam)
	return BOOLToBool(val)
}

func PostMessageW(hwnd HWND, msg uint32, wParam WPARAM, lParam LPARAM) bool {
	val, _ := _PostMessageW(hwnd, msg, wParam, lParam)
	return BOOLToBool(val)
}

func SetForegroundWindow(hwnd HWND) bool {
	ret, _ := _SetForegroundWindow(hwnd)
	return BOOLToBool(ret)
}

func SendMessage(hwnd HWND, msg uint32, wParam WPARAM, lParam LPARAM) (ret LRESULT) {
	if isUnicode() {
		ret = SendMessageW(hwnd, msg, wParam, lParam)
	} else {
		ret = SendMessageA(hwnd, msg, wParam, lParam)
	}
	return ret
}

func SendMessageA(hwnd HWND, msg uint32, wParam WPARAM, lParam LPARAM) LRESULT {
	ret, _ := _SendMessageA(hwnd, msg, wParam, lParam)
	return ret
}

func SendMessageW(hwnd HWND, msg uint32, wParam WPARAM, lParam LPARAM) LRESULT {
	ret, _ := _SendMessageW(hwnd, msg, wParam, lParam)
	return ret
}
