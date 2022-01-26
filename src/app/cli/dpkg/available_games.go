package dpkg

import (
	"log"
	"os"

	"github.com/laper32/regsm-console/src/lib/conf"
)

type AvailableGame struct {
	Name     string
	Short    string
	Specific map[string]interface{}
}

var (
	availableGameList []AvailableGame
)

func InitAvailableGameData() {
	cfg := conf.Load(&conf.Config{
		Name: "available_games",
		Type: "json",
		Path: []string{os.Getenv("GSM_ROOT") + "/config/gsm"},
	})
	var availableGameMap []map[string]interface{}
	err := cfg.UnmarshalKey("AvailableGames", &availableGameMap)
	if err != nil {
		log.Fatalf("FATAL: %v\n", err)
		return
	}
	for _, content := range availableGameMap {
		var this AvailableGame
		this.Name = content["name"].(string)
		this.Short = content["short"].(string)
		this.Specific = content["specific"].(map[string]interface{})
		availableGameList = append(availableGameList, this)
	}
}

func AvailableGames() []AvailableGame {
	return availableGameList
}

// Accept both its name or short name
func FindGame(game string) *AvailableGame {
	for _, content := range availableGameList {
		if content.Name == game || content.Short == game {
			return &content
		}
	}
	return nil
}
