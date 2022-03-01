package conf

import (
	"fmt"

	"github.com/laper32/regsm-console/src/app/server/conf/game"
	"github.com/spf13/viper"
)

func buildStartupParams(v *viper.Viper) {
	switch v.GetString("server.game") {
	case "cs1.6":
		game.CS16StartupConfig(v)
	case "csgo":
		game.CSGOStartupConfig(v)
	default:
		panic(fmt.Sprintf("Not supported game: \"%v\"", v.GetString("server.game")))
	}
}
