package conf

import (
	"fmt"
	"os"

	"github.com/laper32/regsm-console/src/app/cli/conf/game"
	"github.com/laper32/regsm-console/src/app/cli/dpkg"
	"github.com/laper32/regsm-console/src/lib/conf"
	"github.com/spf13/viper"
)

func WriteGameConfigSpecial(cfg *viper.Viper, whatGame string) {
	switch whatGame {
	case "csgo":
		game.CSGOConfigField(cfg)
	case "cs1.6":
		game.CS16ConfigField(cfg)
	}
}

func WriteGameConfigField(cfg *viper.Viper, whatGame string) {
	// Server game. DO NOT MODIFY
	conf.TestOrInsert(cfg, "server.game", whatGame)

	// Server IP, I think I don't have to explain what it is.
	conf.TestOrInsert(cfg, "server.ip", "localhost")

	// Server port, I think I don't have to explain what it is.
	conf.TestOrInsert(cfg, "server.port", 27015)

	// This server's name, used for displaying on the frontend
	conf.TestOrInsert(cfg, "server.name", "A server under Game Server Manager")

	// The parameters for server to startup
	// Game specific.
	conf.TestOrInsert(cfg, "server.param", []string{})

	// Maxplayers in a single server
	// Typically, if you are a 'server', you will have a 'maxplayer' param, despite some
	// of them will say they do not have this definition.
	// However, despite these server said they 'don't' have limitation
	// of max players, will you even trust what they are saying? No, isn't it?
	// So, the maxplayer becomes the default configuration param.
	// But to shut these guys mouse up, we define -1 as unlimited max players.
	conf.TestOrInsert(cfg, "server.maxplayer", 24)

	// Allow update when we want to update the server
	// We define this field in order to prevent some case (especially imported server)
	// they don't need to update server, that we don't have to waste any time on it,
	// such as minecraft and cs1.6.
	conf.TestOrInsert(cfg, "server.allow_update", true)

	// Once this has been enabled, if we have found an upgrade, the GSM will stop
	// the server immediately, and them execute update process.
	// This is by design, because we need to handle some cases like CSGO: updating
	// to the lastest version is compulsory.
	conf.TestOrInsert(cfg, "server.auto_update", false)

	// Update the server before we start the server.
	conf.TestOrInsert(cfg, "server.update_on_start", false)

	// Restart when the server is crashed.
	conf.TestOrInsert(cfg, "server.auto_restart", true)

	// When crashed, how many seconds for us to wait to restart this server?
	// -1 means that never retry
	conf.TestOrInsert(cfg, "server.restart_after_delay", 5)

	// How many times for restarting the server when the server crashed?
	// -1 means infinity.
	conf.TestOrInsert(cfg, "server.max_restart_count", -1)

	// Maximum retry count for connecting to the coordinator.
	conf.TestOrInsert(cfg, "server.max_retry_coordinator_startup_connection_count", 5)

	WriteGameConfigSpecial(cfg, whatGame)
}

// Generate a blank config.
func WriteDefaultGameConfig(whatGame, configDirectory string) {
	// Dirty work
	cfg := viper.New()
	cfg.SetConfigName("config")
	cfg.SetConfigType("toml")
	// The location is under ${RootDir}/config/server/${Index}
	cfg.AddConfigPath(configDirectory)

	WriteGameConfigField(cfg, whatGame)

	err := cfg.SafeWriteConfig()
	if err != nil {
		panic(err)
	}
}

func writeStartupArgs(v *viper.Viper) {
	switch v.GetString("server.game") {
	case "csgo":
		game.CSGOStartupConfig(v)
	case "cs1.6":
		game.CS16StartupConfig(v)
	}
}

func StartGameConfiguration(server *dpkg.ServerIdentity) (*viper.Viper, error) {
	configDirectory := fmt.Sprintf("%v/config/server/%v", os.Getenv("GSM_ROOT"), server.ID)
	cfg := conf.Load(&conf.Config{
		Name: "config",
		Type: "toml",
		Path: []string{configDirectory},
	})

	// We need to update config
	WriteGameConfigField(cfg, server.Game)
	err := cfg.WriteConfig()
	if err != nil {
		fmt.Println("Unexpected occured when writing the config:", err)
		return nil, err
	}

	writeStartupArgs(cfg)
	return cfg, nil
}
