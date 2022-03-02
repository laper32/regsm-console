package dpkg

import (
	"fmt"
	"log"
	"os"

	libconf "github.com/laper32/regsm-console/src/lib/conf"
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
	serverIdentityTable  map[uint][]uint = make(map[uint][]uint)
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

	createIdentityTable()
}

func element_exist(val uint, arr []uint) bool {
	for _, v := range arr {
		if v == val {
			return true
		}
	}
	return false
}

func GetServerChainByID(serverID uint) []uint {
	// search root
	if value, ok := serverIdentityTable[serverID]; ok {
		return append([]uint{serverID}, value...)
	}

	// walking through table
	for k, v := range serverIdentityTable {
		if element_exist(serverID, v) {
			return append([]uint{k}, v...)
		}
	}

	return nil
}
func createIdentityTable() {
	for _, v := range ServerIdentityList {
		if v.SymlinkServerID == 0 {
			serverIdentityTable[v.ID] = make([]uint, 0)
			continue
		}
		if value, ok := serverIdentityTable[v.SymlinkServerID]; ok {
			if !element_exist(v.ID, value) {
				serverIdentityTable[v.SymlinkServerID] = append(value, v.ID)
			}
		} else {
			sv := FindRootSymlinkServer(v.ID)
			if value, ok := serverIdentityTable[sv.ID]; ok {
				if !element_exist(v.ID, value) {
					serverIdentityTable[sv.ID] = append(value, v.ID)
				}
			}
		}
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
var (
	__builtin_last_symlink_id uint = 0
)

// server ID | symlink Server ID
// 		8			7
// 		7			8
//  passing this server ID, then we check the last symlink id
// then we can found the recursive referencing.
func FindRootSymlinkServer(serverID uint) *ServerIdentity {
	if __builtin_last_symlink_id == serverID {
		panic(fmt.Sprintf("recursive server found: %v and %v", __builtin_last_symlink_id, serverID))
	}
	sv := FindIdentifiedServer(serverID)
	for _, v := range ServerIdentityList {
		if v.ID != sv.SymlinkServerID {
			continue
		}
		if v.SymlinkServerID != 0 {
			__builtin_last_symlink_id = v.SymlinkServerID
			return FindRootSymlinkServer(v.ID)
		} else {
			return &v
		}
	}
	return nil
}
