package game

import (
	"fmt"

	"github.com/spf13/viper"
)

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
