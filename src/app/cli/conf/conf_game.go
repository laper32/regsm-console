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
	conf.TestOrInsert(cfg, "server.maxplayer", -1)

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
	conf.TestOrInsert(cfg, "server.restart_after_seconds", 5)

	// The maximum retry count.
	// Shared for both following two configurations.
	conf.TestOrInsert(cfg, "coordinator.retry_count", 5)

	// Determine whether should we allow retry when starting up the server.
	// This is used by the case that when we cannot establish connection to the coordinator,
	// we met some trouble of opening the daemon executable, etc.
	// Overall, this is to determine that what should we do when we cannot starting up the server
	// correctly.
	conf.TestOrInsert(cfg, "coordinator.allow_retry_at_startup", true)

	// Determine whether should we allow trying to reconnect to the coordinator once it has been
	// disconnected.
	// This is used for the case when the coordinator offline, there are some issues of current
	// daemon's websocket connection, etc.
	// If we cannot reconnect to the coordinator, then everything will be stored as a local file.
	// Once we have been reconnected to the coordinator, the data will be packed, and send to the
	// coordinator.
	// Noting that if disconnected, everything what is related to the coordinator ARE ALL DISABLED!
	conf.TestOrInsert(cfg, "coordinator.allow_reconnect_when_running", true)

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
