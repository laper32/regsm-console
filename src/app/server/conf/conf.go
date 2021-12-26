package conf

import (
	"strconv"

	"github.com/laper32/regsm-console/src/lib/conf"
)

type Config struct {
	Coordinator struct {
		IP   string
		Port uint
	}
	MaxRetryCount int
}

func Init(serverID int) *Config {
	v := conf.Load(&conf.Config{
		Name: "coordinator",
		Type: "toml",
		Path: []string{"../config/gsm" + strconv.Itoa(serverID)},
	})
	return &Config{
		Coordinator: struct {
			IP   string
			Port uint
		}{IP: v.GetString("server.coordinator.ip"), Port: v.GetUint("server.coordinator.port")},
		MaxRetryCount: v.GetInt("server.max_retry"),
	}
}
