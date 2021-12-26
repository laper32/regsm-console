package conf

import (
	"github.com/laper32/regsm-console/src/lib/log"
)

type Config struct {
	Log *log.Config
}

func Init() *Config {
	return &Config{
		Log: &log.Config{},
	}
}
