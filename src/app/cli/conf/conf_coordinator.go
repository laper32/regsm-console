package conf

import (
	"fmt"
	"os"

	"github.com/laper32/regsm-console/src/lib/conf"
	"github.com/spf13/viper"
)

func WriteCoordinatorConfigField(cfg *viper.Viper) {
	// This coordinator's IP
	conf.TestOrInsert(cfg, "coordinator.ip", "localhost")

	// This coordinator's port
	conf.TestOrInsert(cfg, "coordinator.port", 3484)

	// Determine whether this coordinator accepts only other coordinators', or servers' connection.
	// This is extermely useful when you have a plenty of servers.
	conf.TestOrInsert(cfg, "coordinator.pure", false)

	// Other coordinator's address.
	// Once you set this field, we will trying to connect to other coordinators.
	conf.TestOrInsert(cfg, "coordinator.other_coordinator_address", "")
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
