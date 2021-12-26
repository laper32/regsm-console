package path

import (
	"os"
	"path/filepath"
	"strings"
)

func Exist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func GetFileDir(thisFile string) string {
	ret, _ := filepath.Abs(filepath.Dir(thisFile))
	// requires that all dir MUST be / rather than \\
	return strings.Replace(ret, "\\", "/", -1)
}
