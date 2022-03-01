package game

import (
	"fmt"

	"github.com/spf13/viper"
)

func CSGOStartupConfig(cfg *viper.Viper) {
	paramList := cfg.GetStringSlice("server.param")
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
