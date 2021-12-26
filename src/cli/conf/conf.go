package conf

import (
	"github.com/laper32/regsm-console/src/lib/log"
)

type Config struct {
	log *log.Config
}

func Init() *Config {
	return &Config{
		log: &log.Config{},
	}
}
