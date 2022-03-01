package game

import (
	"github.com/laper32/regsm-console/src/lib/conf"
	"github.com/spf13/viper"
)

func CSGOConfigField(cfg *viper.Viper) {
	conf.TestOrInsert(cfg, "server.special.gslt", "")
	conf.TestOrInsert(cfg, "server.special.default_map", "de_dust2")
	// We only provides default config
	// You should make further configruation by yourself, rather than us.
	cfg.Set("server.param", []string{
		"-console", "-game csgo", "+game_type 0", "+game_mode 0",
		"+mapgroup mg_active", "-tickrate 64", "-usercon",
	})
	// conf.TestOrInsert(cfg, "server.param", []string{
	// })
}
