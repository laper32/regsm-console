package dpkg

import (
	"log"
	"os"

	libconf "github.com/laper32/regsm-console/src/lib/conf"
	"github.com/spf13/viper"
)

type ServerInfo struct {
	ID        uint   `map:"m_nIndex"`
	Game      string `map:"m_sGame"`
	Deleted   bool   `map:"m_bDeleted"`
	RelatedTo int    `map:"m_nRelatedTo"`
}

var (
	ServerInfoList   []ServerInfo
	ServerInfoConfig *viper.Viper
	ServerInfoMap    []map[string]interface{}
)

func InitServerInfo() {
	configDirectory := os.Getenv("GSM_ROOT") + "/config/gsm"

	ServerInfoConfig = libconf.Load(&libconf.Config{
		Name: "server_info",
		Type: "json",
		Path: []string{configDirectory},
	})

	err := ServerInfoConfig.UnmarshalKey("server_info", &ServerInfoMap)
	if err != nil {
		log.Fatalln("Unexpected error here: ", err)
		return
	}

	for _, content := range ServerInfoMap {
		var this ServerInfo
		this.ID = uint(content["m_nIndex"].(float64))
		this.Game = content["m_sGame"].(string)
		this.Deleted = content["m_bDeleted"].(bool)
		this.RelatedTo = int(content["m_nRelatedTo"].(float64))
		ServerInfoList = append(ServerInfoList, this)
	}
}

func FindServerInfoServer(serverID uint) *ServerInfo {
	for _, this := range ServerInfoList {
		if this.ID == serverID {
			return &this
		}
	}
	return nil
}

func FindRootRelatedServer(serverID uint) *ServerInfo {
	for _, this := range ServerInfoList {
		if this.RelatedTo != -1 {
			return FindRootRelatedServer(this.ID)
		} else {
			return &this
		}
	}
	return nil
}
