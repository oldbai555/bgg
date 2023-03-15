package conf

import (
	"fmt"
	webtool "github.com/oldbai555/bgg/pkg/webtoolv2"
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
		webtool.OptionWithOrm(),
		webtool.OptionWithRdb(),
		webtool.OptionWithStorage(),
		webtool.OptionWithServer(),
		webtool.OptionWithChatGpt(),
		webtool.OptionWithWxGzh())

	if err != nil {
		log.Errorf("err:%v", err)
		panic(err)
	}
}

func initViper() (*viper.Viper, error) {
	viper.SetConfigName("application")          // name of config file (without extension)
	viper.SetConfigType("yaml")                 // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("/etc/work/")           // path to look for the config file in
	viper.AddConfigPath("./")                   // optionally look for config in the working directory
	viper.AddConfigPath("./lbserver/cmd/")      // optionally look for config in the working directory
	viper.AddConfigPath("../cmd/")              // optionally look for config in the working directory
	viper.AddConfigPath("./lbserver/resource/") // optionally look for config in the working directory
	viper.AddConfigPath("../resource/")         // optionally look for config in the working directory
	err := viper.ReadInConfig()                 // Find and read the config file
	if err != nil {                             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	return viper.GetViper(), nil
}
