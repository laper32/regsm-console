package conf

import (
	"github.com/laper32/regsm-console/src/lib/log"
	"github.com/spf13/viper"
)

type Config struct {
	Type string
	Name string
	Path []string
}

// Load is to load the config, and return its pointer.
//
// Note:
//
// 1. If you have a same name but different type config, you must input the fullname of this config
// including its file extension. Check conf_test for more detail.
//
// 2. Check conf_test about how to use this API.
//
// 3. Assuming you have basic knowledge about config structure. See https://github.com/spf13/viper
func Load(conf *Config) *viper.Viper {
	v := viper.New()
	v.SetConfigName(conf.Name)
	for _, val := range conf.Path {
		v.AddConfigPath(val)
	}

	v.SetConfigType(conf.Type)

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			if err = v.SafeWriteConfig(); err != nil {
				log.Error(err)
			}
		} else {

			log.Panic(err)
		}
	}

	return v
}
