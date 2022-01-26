package windows

import (
	"syscall"
	"unsafe"
)

//sys	_CloseHandle(handle HANDLE) (err error) = kernel32.CloseHandle
//sys	_CreateToolhelp32Snapshot(flags uint32, processId uint32) (handle HANDLE, err error) [failretval==InvalidHandle] = kernel32.CreateToolhelp32Snapshot
//sys	_GetLastError() (lasterr error) = kernel32.GetLastError
//sys	_LoadLibraryA(libname *byte) (handle HANDLE, err error) = kernel32.LoadLibraryA
//sys	_LoadLibraryW(libname *uint16) (handle HANDLE, err error) = kernel32.LoadLibraryW
//sys	_LoadLibraryExA(libname *byte, zero HANDLE, flags uintptr) (handle HANDLE, err error) = kernel32.LoadLibraryExA
//sys	_LoadLibraryExW(libname *uint16, zero HANDLE, flags uintptr) (handle HANDLE, err error) = kernel32.LoadLibraryExW
//		RtlCopyMemory(Destination,Source,Length) = memcpy((Destination),(Source),(Length)) => Not implemented here.
//		_SetConsoleCtrlHandler(handlerRoutine PHANDLER_ROUTINE, add uint) (err error) = kernel32.SetConsoleCtrlHandler
//sys	_SetThreadExecutionState(esFlags EXECUTION_STATE) (ret EXECUTION_STATE, err error) = kernel32.SetThreadExecutionState

var (
	__SetConsoleCtrlHandler = kernel32.NewProc("SetConsoleCtrlHandler")
)

func CreateToolhelp32Snapshot(flags uint32, pid uint32) (HANDLE, error) {
	return _CreateToolhelp32Snapshot(flags, pid)
}

func GetLastError() error {
	return _GetLastError()
}

func LoadLibrary(libname string) (handle HANDLE) {
	if isUnicode() {
		handle = LoadLibraryW(libname)
	} else {
		handle = LoadLibraryA(libname)
	}
	return handle
}

func LoadLibraryA(libname string) HANDLE {
	stringPtr, _ := syscall.BytePtrFromString(libname)
	ret, err := _LoadLibraryA(stringPtr)
	if err != nil {
		panic(err)
	}
	return ret
}

func LoadLibraryW(libname string) HANDLE {
	stringPtr, _ := syscall.UTF16PtrFromString(libname)
	ret, err := _LoadLibraryW(stringPtr)
	if err != nil {
		panic(err)
	}
	return ret
}

func LoadLibraryEx(libname string, zero HANDLE, flags uintptr) (handle HANDLE) {
	if isUnicode() {
		handle = LoadLibraryExW(libname, zero, flags)
	} else {
		handle = LoadLibraryExA(libname, zero, flags)
	}
	return handle
}

func LoadLibraryExA(libname string, zero HANDLE, flags uintptr) HANDLE {
	stringPtr, _ := syscall.BytePtrFromString(libname)
	ret, err := _LoadLibraryExA(stringPtr, zero, flags)
	if err != nil {
		panic(err)
	}
	return ret
}

func LoadLibraryExW(libname string, zero HANDLE, flags uintptr) HANDLE {
	stringPtr, _ := syscall.UTF16PtrFromString(libname)
	ret, err := _LoadLibraryExW(stringPtr, zero, flags)
	if err != nil {
		panic(err)
	}
	return ret
}

func SetConsoleCtrlHandler(handlerRoutine PHANDLER_ROUTINE, add uint) error {
	return _SetConsoleCtrlHandler(handlerRoutine, add)
}

//Adds or removes an application-defined HandlerRoutine function from the list of handler functions for the calling process.
//https://msdn.microsoft.com/en-us/library/windows/desktop/ms686016(v=vs.85).aspx
func _SetConsoleCtrlHandler(handlerRoutine PHANDLER_ROUTINE, add uint) (err error) {
	_, _, err = __SetConsoleCtrlHandler.Call(uintptr(unsafe.Pointer(&handlerRoutine)), uintptr(add))
	if !IsErrSuccess(err) {
		return err
	}
	return nil
}

func SetThreadExecutionState(esFlags EXECUTION_STATE) EXECUTION_STATE {
	ret, _ := _SetThreadExecutionState(esFlags)
	return ret
}
