package conf

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/laper32/regsm-console/src/lib/conf"
	"github.com/laper32/regsm-console/src/lib/log"
)

type Config struct {
	Log         *log.Config
	Server      *Server
	Coordinator *Coordinator
}

type Server struct {
	ID                                        uint
	AllowUpdate                               bool
	AutoRestart                               bool
	AutoUpdate                                bool
	Game                                      string
	IP                                        string
	MaxRestartCount                           int
	MaxRetryCoordinatorStartupConnectionCount int
	MaxPlayer                                 int
	Name                                      string
	Args                                      []string
	Port                                      uint
	RestartAfterDelay                         int
	UpdateOnStart                             bool
	Special                                   map[string]interface{}
}

type Coordinator struct {
	IP   string
	Port uint
}

/*
[server]
  allow_update = false
  auto_restart = true
  auto_update = false
  game = "cs1.6"
  ip = "localhost"
  max_restart_count = -1
  max_retry_coordinator_startup_connection_count = 5
  maxplayer = 32
  name = "A server under Game Server Manager"
  param = ["-console", "-game cstrike"]
  port = 23333
  restart_after_delay = 5
  update_on_start = false

  [server.special]
    default_map = "de_dust2"
    gslt = ""
*/
func Init() *Config {
	cfg_server := conf.Load(&conf.Config{
		Name: "config",
		Type: "toml",
		Path: []string{fmt.Sprintf("%v/config/server/%v", os.Getenv("GSM_ROOT"), os.Getenv("GSM_SERVER_ID"))},
	})
	cfg_coordinator := conf.Load(&conf.Config{
		Name: "coordinator",
		Type: "toml",
		Path: []string{fmt.Sprintf("%v/config/gsm", os.Getenv("GSM_ROOT"))},
	})
	buildStartupParams(cfg_server)
	var args []string
	for _, v := range cfg_server.GetStringSlice("server.param") {
		args = append(args, strings.Split(v, " ")...)
	}
	serverID, err := strconv.ParseUint(os.Getenv("GSM_SERVER_ID"), 10, 64)
	if err != nil {
		panic(err)
	}
	return &Config{
		Log: &log.Config{},
		Server: &Server{
			ID:              uint(serverID),
			AllowUpdate:     cfg_server.GetBool("server.allow_update"),
			AutoRestart:     cfg_server.GetBool("server.auto_restart"),
			AutoUpdate:      cfg_server.GetBool("server.auto_update"),
			Game:            cfg_server.GetString("server.game"),
			IP:              cfg_server.GetString("server.ip"),
			MaxRestartCount: cfg_server.GetInt("server.max_restart_count"),
			MaxRetryCoordinatorStartupConnectionCount: cfg_server.GetInt("server.max_retry_coordinator_startup_connection_count"),
			MaxPlayer:         cfg_server.GetInt("server.maxplayer"),
			Name:              cfg_server.GetString("server.name"),
			Args:              args,
			Port:              cfg_server.GetUint("server.port"),
			RestartAfterDelay: cfg_server.GetInt("server.restart_after_delay"),
			UpdateOnStart:     cfg_server.GetBool("server.update_on_start"),
			Special:           cfg_server.GetStringMap("server.special"),
		},
		Coordinator: &Coordinator{
			IP:   cfg_coordinator.GetString("coordinator.ip"),
			Port: cfg_coordinator.GetUint("coordinator.port"),
		},
	}
}
