package webtool

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

func GenWebToolByYaml(srv string, option ...Option) *WebTool {
	viper.SetConfigName(srv)           // name of config file (without extension)
	viper.SetConfigType(defaultSuffix) // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(defaultPath)   // path to look for the config file in
	err := viper.ReadInConfig()        // Find and read the config file
	if err != nil {                    // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	global, err := NewWebTool(viper.GetViper(), option...)
	if err != nil {
		log.Errorf("err:%v", err)
		panic(err)
	}

	return global
}

func GenDefaultWebTool(option ...Option) *WebTool {
	viper.SetConfigName(defaultPrefix) // name of config file (without extension)
	viper.SetConfigType(defaultSuffix) // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(defaultPath)   // path to look for the config file in
	err := viper.ReadInConfig()        // Find and read the config file
	if err != nil {                    // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	global, err := NewWebTool(viper.GetViper(), option...)
	if err != nil {
		log.Errorf("err:%v", err)
		panic(err)
	}

	return global
}
