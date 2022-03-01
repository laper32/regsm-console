package log

import (
	"fmt"
	"testing"
	"time"
)

func TestInit(t *testing.T) {
	logfilename := time.Now().Format("20060102") + ".log"
	fmt.Println(logfilename)
	tests := []struct {
		name string
		args Config
	}{
		{
			"pro",
			Config{Debug: false, OutputPath: []string{"stdout", logfilename}},
		},
		{
			"debug",
			Config{Debug: true, OutputPath: []string{"stdout", logfilename}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Init(&tt.args)
			defer func() {
				Close()
				Warn(recover())
			}()

			Info("This is a info message")
			Warn("This is a warn message")
			Error("This a error message")
			Debug("This a debug message")
			Panic("This a panic message")
		})
	}
}
