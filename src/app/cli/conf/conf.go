package conf

import (
	"fmt"
	"os"
	"time"

	"github.com/laper32/regsm-console/src/lib/log"
)

type Config struct {
	Log *log.Config
}

func Init() *Config {
	logPath := fmt.Sprintf("%v/log/gsm/L%v.log", os.Getenv("GSM_ROOT"), time.Now().Format("20060102"))
	return &Config{
		Log: &log.Config{OutputPath: []string{"stdout", logPath}},
	}
}
