package daemon

import (
	"fmt"
	"testing"
	"time"
)

func TestDate(t *testing.T) {
	now := time.Now()
	// y, m, d := now.Date()
	_str := now.Format("20060102") // how weird...
	fmt.Println("L" + _str + ".log")
}
