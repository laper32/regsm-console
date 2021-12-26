package conf

import (
	"fmt"
	"testing"
)

func TestLoadJSON(t *testing.T) {
	conf := &Config{
		Name: "config",
		Type: "json",
		Path: []string{"."},
	}
	v := Load(conf)
	fmt.Println(v.GetString("Test.Section1.Value"))
}

// Don't really know why this part loaded config.json
// Looks like it is a bug here...
// maybe someone open an issue about how to fix/workaround on this?
func TestLoadTOML(t *testing.T) {
	conf := &Config{
		// if not explicitly declare its suffix, then goes error
		// aka, if declare "config" rather than "config.toml" then goes error
		Name: "config.toml",
		Type: "toml",
		Path: []string{"."},
	}
	v := Load(conf)
	fmt.Println(v.GetString("Test.Section1.Value"))
}

func TestLoadYAML(t *testing.T) {
	conf := &Config{
		Name: "config",
		Type: "yaml",
		Path: []string{"."},
	}
	v := Load(conf)
	fmt.Println(v.GetString("Test.Section1.Value"))
}
