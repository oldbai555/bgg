package conf

import (
	"fmt"
	"github.com/oldbai555/bgg/pkg/webtool"
	"github.com/oldbai555/lbtool/log"
	"github.com/spf13/viper"
)

var Global *webtool.WebTool // 位置有点不对,待考虑放哪

func InitWebTool() {
	v, err := initViper()
	if err != nil {
		log.Errorf("err is %v", err)
		panic(err)
	}
	Global, err = webtool.NewWebTool(v,
		webtool.OptionWithServer(),
	)

	if err != nil {
		log.Errorf("err:%v", err)
		panic(err)
	}
}

func initViper() (*viper.Viper, error) {
	viper.SetConfigName("application")                                                                        // name of config file (without extension)
	viper.SetConfigType("yaml")                                                                               // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("/etc/work/")                                                                         // path to look for the config file in
	viper.AddConfigPath("/Users/zhangjianjun/work/lb/github.com/oldbai555/bgg/service/gateway/internal/conf") // optionally look for config in the working directory
	err := viper.ReadInConfig()                                                                               // Find and read the config file
	if err != nil {                                                                                           // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	return viper.GetViper(), nil
}
