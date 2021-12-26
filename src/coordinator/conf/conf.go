package conf

import (
	"github.com/laper32/regsm-console/src/lib/conf"
	"github.com/laper32/regsm-console/src/lib/log"
)

type Config struct {
	Coordinator struct {
		Port           uint
		Pure           bool
		ConnectAddress string
	}
	Log *log.Config
}

func Init() *Config {
	v := conf.Load(&conf.Config{
		Name: "coordinator",
		Type: "toml",
		Path: []string{"../config/gsm"},
	})
	return &Config{
		Coordinator: struct {
			Port           uint
			Pure           bool
			ConnectAddress string
		}{
			Port:           v.GetUint("coordinator.port"),
			Pure:           v.GetBool("coordinator.pure"),
			ConnectAddress: v.GetString("coordinator.connect_address"),
		},
		Log: &log.Config{},
	}
}
