package syscfg

import (
	"fmt"
	"github.com/oldbai555/lbtool/log"
	"github.com/spf13/viper"
)

const (
	defaultPath   = "/etc/work/"
	defaultPrefix = "application"
	defaultSuffix = "yaml"
)

func LoadSysCfgByYaml(srv, path string, option ...Option) *SysCfg {
	if path == "" {
		path = defaultPath
	}

	if srv == "" {
		srv = defaultPrefix
	}

	viper.SetConfigName(srv)           // name of config file (without extension)
	viper.SetConfigType(defaultSuffix) // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(path)          // path to look for the config file in
	err := viper.ReadInConfig()        // Find and read the config file
	if err != nil {                    // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	global, err := New(viper.GetViper(), option...)
	if err != nil {
		log.Errorf("err:%v", err)
		panic(err)
	}

	return global
}
