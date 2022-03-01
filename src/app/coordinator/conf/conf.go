package conf

import (
	"fmt"
	"os"

	"github.com/laper32/regsm-console/src/lib/conf"
	"github.com/laper32/regsm-console/src/lib/log"
)

type Config struct {
	Log   *log.Config
	Param *Parameter
}

type Parameter struct {
	IP                      string
	Port                    uint
	Pure                    bool
	OtherCoordinatorAddress string
}

func Init() *Config {

	cfg := conf.Load(&conf.Config{
		Name: "coordinator",
		Type: "toml",
		Path: []string{fmt.Sprintf("%v/config/gsm", os.Getenv("GSM_ROOT"))},
	})
	return &Config{
		Log: &log.Config{},
		Param: &Parameter{
			IP:                      cfg.GetString("coordinator.ip"),
			Port:                    cfg.GetUint("coordinator.port"),
			Pure:                    cfg.GetBool("coordinato.pure"),
			OtherCoordinatorAddress: cfg.GetString("coordinator.other_coordinator_address"),
		},
	}
}
