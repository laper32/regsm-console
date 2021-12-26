package path

import (
	"fmt"
	"testing"
)

func TestPathExist(t *testing.T) {
	fmt.Println("Exist: ", Exist("./path.go"))
}
