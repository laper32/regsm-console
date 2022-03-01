package game

import (
	"github.com/laper32/regsm-console/src/lib/conf"
	"github.com/spf13/viper"
)

func CS16ConfigField(cfg *viper.Viper) {
	// If you downloaded via steamCMD, you must make a configuration on it.
	// But if you imported from an existing server (eg: reHLDS), DO NOT ENTER THIS FIELD!
	conf.TestOrInsert(cfg, "server.special.gslt", "")

	// Hmmm? Need I explain?
	conf.TestOrInsert(cfg, "server.special.default_map", "de_dust2")

	// We put in param in the specific game, because different game has different parameters.
	// These are all default parameters, but you can modify it when we saved on disk.
	cfg.Set("server.param", []string{
		"-console", "-game cstrike",
	})

}
