package game

import (
	"fmt"

	"github.com/spf13/viper"
)

func CSGODefaultConfig(cfg *viper.Viper) {
	cfg.Set("server.special.gslt", "")
	cfg.Set("server.special.default_map", "")
	// We only provides default config
	// You should make further configruation by yourself, rather than us.
	cfg.Set("server.param", []string{
		"-console", "-game csgo", "+game_type 0", "+game_mode 0",
		"+mapgroup mg_active", "-tickrate 64", "-usercon",
	})
}

func CSGOStartupConfig(cfg *viper.Viper) {
	paramList := cfg.Get("server.param").([]string)
	GSLTToken := cfg.GetString("server.special.gslt")
	defaultMap := cfg.GetString("server.special.default_map")
	ip := cfg.GetString("server.ip")
	port := cfg.GetUint("server.port")
	maxplayer := cfg.GetInt("server.maxplayer")

	paramList = append(paramList, []string{
		fmt.Sprintf("-ip %v", ip), fmt.Sprintf("-port %v", port), fmt.Sprintf("-maxplayers_override %v", maxplayer),
		fmt.Sprintf("+map %v", defaultMap), fmt.Sprintf("+sv_setsteamaccount %v", GSLTToken)}...)

	cfg.Set("server.param", paramList)
}
