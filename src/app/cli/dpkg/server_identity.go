package dpkg

import (
	"fmt"
	"log"
	"os"

	libconf "github.com/laper32/regsm-console/src/lib/conf"
	"github.com/laper32/regsm-console/src/lib/status"
	"github.com/laper32/regsm-console/src/lib/structs"
	"github.com/spf13/viper"
)

type ServerIdentity struct {
	ID              uint   `map:"m_nIndex"`
	Game            string `map:"m_sGame"`
	Deleted         bool   `map:"m_bDeleted"`
	SymlinkServerID uint   `map:"m_nSymlinkServerID"`
}

var (
	ServerIdentityList   []ServerIdentity
	ServerIdentityConfig *viper.Viper
	ServerIdentityMap    []map[string]interface{}
)

func InitServerIdentity() {
	configDirectory := os.Getenv("GSM_ROOT") + "/config/gsm"

	ServerIdentityConfig = libconf.Load(&libconf.Config{
		Name: "server_identity",
		Type: "json",
		Path: []string{configDirectory},
	})

	err := ServerIdentityConfig.UnmarshalKey("server_identity", &ServerIdentityMap)
	if err != nil {
		log.Fatalln("Unexpected error here: ", err)
		return
	}

	for _, content := range ServerIdentityMap {
		var this ServerIdentity
		this.ID = uint(content["m_nIndex"].(float64))
		this.Game = content["m_sGame"].(string)
		this.Deleted = content["m_bDeleted"].(bool)
		this.SymlinkServerID = uint(content["m_nSymlinkServerID"].(float64))
		ServerIdentityList = append(ServerIdentityList, this)
	}
}

func FindIdentifiedServer(serverID uint) *ServerIdentity {
	for _, this := range ServerIdentityList {
		if this.ID == serverID {
			return &this
		}
	}
	return nil
}

func UpdateServerIdentity() {
	ServerIdentityMap = ServerIdentityMap[0:0]
	for _, this := range ServerIdentityList {
		_map, _ := structs.ToMap(this, "map")
		ServerIdentityMap = append(ServerIdentityMap, _map)
	}
	ServerIdentityConfig.Set("server_identity", ServerIdentityMap)
	ServerIdentityConfig.WriteConfig()
}

// Working around to solve recursive referencing issue.
var __builtin_current_id uint = 0
var __builtin_symlink_id uint = 0

func FindRootSymlinkServer(serverID uint) *ServerIdentity {
	// Working around to solve recursive referencing issue.
	if __builtin_current_id == serverID {
		panic(status.CLIInstallSymlinkServerIDFoundRecursiveReferencing.WriteDetail(fmt.Sprintf("Recursive server ID: {%v, %v}", __builtin_current_id, __builtin_symlink_id)))
	}
	for _, this := range ServerIdentityList {
		if this.SymlinkServerID != 0 {
			__builtin_current_id = this.ID
			__builtin_symlink_id = this.SymlinkServerID
			return FindRootSymlinkServer(this.ID)
		} else {
			return &this
		}
	}
	return nil
}
