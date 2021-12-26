package conf

import (
	"fmt"
	"os"

	"github.com/laper32/regsm-console/src/app/cli/conf/game"
	"github.com/laper32/regsm-console/src/lib/conf"
	"github.com/spf13/viper"
)

func writeDefaultSpecial(cfg *viper.Viper, whatGame string) {
	switch whatGame {
	case "csgo":
		game.CSGODefaultConfig(cfg)
	case "cs1.6":
		game.CS16DefaultConfig(cfg)
	}
}

// Generate a blank config.
func WriteDefaultGameConfig(whatGame, configDirectory string) {
	// Dirty work
	cfg := viper.New()
	cfg.SetConfigName("config")
	cfg.SetConfigType("toml")
	// The location is under ${RootDir}/config/server/${Index}
	cfg.AddConfigPath(configDirectory)
	// Server game. DO NOT MODIFY
	cfg.Set("server.game", whatGame)

	// Server IP, I think I don't have to explain what it is.
	cfg.Set("server.ip", "localhost")
	// Server port, I think I don't have to explain what it is.
	cfg.Set("server.port", 27015)

	// This server's name, used for displaying on the frontend
	cfg.Set("server.name", "A server under Game Server Manager")

	// The parameters for server to startup
	// Game specific.
	cfg.Set("server.param", []string{})

	// Maxplayers in one server
	// Typically, if you are a 'server', you will have a 'maxplayer' param, despite some
	// of them will say they do not have this definition.
	// However, despite these server said they 'don't' have limitation
	// of max players, will you even trust what they are saying? No, isn't it?
	// So, the maxplayer becomes the default configuration param.
	// But to shut these guys mouse up, we define -1 as unlimited max players.
	cfg.Set("server.maxplayer", -1)

	// Allow update when we want to update the server
	// We define this field in order to prevent some case (especially imported server)
	// they don't need to update server, that we don't have to waste any time on it,
	// such as minecraft and cs1.6.
	cfg.Set("server.allow_update", true)

	// Once this has been enabled, if we have found an upgrade, the GSM will stop
	// the server immediately, and them execute update process.
	// This is by design, because we need to handle some cases like CSGO: updating
	// to the lastest version is compulsory.
	cfg.Set("server.auto_update", false)

	// Nothing to say, right?
	cfg.Set("server.update_on_start", false)

	writeDefaultSpecial(cfg, whatGame)
	err := cfg.SafeWriteConfig()
	if err != nil {
		panic(err)
	}
}

func implementSpecific(v *viper.Viper) {
	switch v.GetString("server.game") {
	case "csgo":
		game.CSGOStartupConfig(v)
	case "cs1.6":
		game.CS16StartupConfig(v)
	}
}

func StartupConfiguration(serverID uint) *viper.Viper {
	configDirectory := fmt.Sprintf("%v/config/server/%v", os.Getenv("GSM_ROOT"), serverID)
	cfg := conf.Load(&conf.Config{
		Name: "config",
		Type: "toml",
		Path: []string{configDirectory},
	})
	implementSpecific(cfg)
	return cfg
}
