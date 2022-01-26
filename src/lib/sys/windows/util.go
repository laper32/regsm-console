package windows

import (
	"os"
	"strconv"
	"syscall"
)

var (
	unicode = os.Getenv("UNICODE")
)

func isUnicode() bool {
	ret, err := strconv.ParseBool(unicode)
	if err != nil {
		ret = false
	}
	return ret
}

func BoolToBOOL(in bool) BOOL {
	if in {
		return 1
	}
	return 0
}

func BOOLToBool(in BOOL) bool {
	return in == 1
}

// IsErrSuccess checks if an "error" returned is actually the
// success code 0x0 "The operation completed successfully."
//
// This is the optimal approach since the error messages are
// localized depending on the OS language.
func IsErrSuccess(err error) bool {
	if errno, ok := err.(syscall.Errno); ok {
		if errno == 0 {
			return true
		}
	}
	return false
}
