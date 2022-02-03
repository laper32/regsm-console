package conf

import (
	"fmt"
	"os"

	"github.com/laper32/regsm-console/src/lib/conf"
	"github.com/spf13/viper"
)

func WriteCoordinatorConfigField(cfg *viper.Viper) {
	conf.TestOrInsert(cfg, "coordinator.ip", "localhost")
	conf.TestOrInsert(cfg, "coordinator.port", 3484)
}

func CoordinatorConfiguration() (*viper.Viper, error) {
	configDirectory := fmt.Sprintf("%v/config/gsm", os.Getenv("GSM_ROOT"))
	cfg := conf.Load(&conf.Config{
		Name: "coordinator",
		Type: "toml",
		Path: []string{configDirectory},
	})
	WriteCoordinatorConfigField(cfg)
	err := cfg.WriteConfig()
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	return cfg, nil
}
