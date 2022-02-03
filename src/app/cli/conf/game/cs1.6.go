package game

import (
	"fmt"

	"github.com/laper32/regsm-console/src/lib/conf"
	"github.com/spf13/viper"
)

func CS16ConfigField(cfg *viper.Viper) {
	// If you downloaded via steamCMD, you must make a configuration on it.
	// But if you imported from an existing server (eg: reHLDS), DO NOT ENTER THIS FIELD!
	conf.TestOrInsert(cfg, "server.special.gslt", "")

	// Hmmm? Need I explain?
	conf.TestOrInsert(cfg, "server.special.default_map", "")

	// We put in param in the specific game, because different game has different parameters.
	// These are all default parameters, but you can modify it when we saved on disk.
	conf.TestOrInsert(cfg, "server.param", []string{
		"-console ", "-game cstrike",
	})
}

func CS16StartupConfig(cfg *viper.Viper) {
	paramList := cfg.GetStringSlice("server.param")
	defaultMap := cfg.GetString("server.special.default_map")
	ip := cfg.GetString("server.ip")
	port := cfg.GetUint("server.port")
	maxplayer := cfg.GetInt("server.maxplayer")
	paramList = append(paramList, []string{
		fmt.Sprintf("-ip %v", ip), fmt.Sprintf("-port %v", port), fmt.Sprintf("-maxplayers %v", maxplayer),
		fmt.Sprintf("+map %s", defaultMap)}...,
	)

	// Will only append it when you requested a GSLT for HLDS
	// If your server is imported (eg, ReHLDS, etc), you SHOULD NOT input the GSLT token!
	if cfg.GetString("server.special.gslt") != "" {
		paramList = append(paramList, fmt.Sprintf("+sv_setsteamaccount %v", cfg.GetString("server.special.gslt")))
	}
	cfg.Set("server.param", paramList)
}
