package conf

import (
	"fmt"
	"github.com/oldbai555/bgg/pkg/webtool"
	"github.com/oldbai555/lbtool/log"
	"github.com/spf13/viper"
	logsys "log"
	"time"
)

var Global *webtool.WebTool // 位置有点不对,待考虑放哪

const (
	// log conf
	LogFlag = logsys.LstdFlags

	// gate conf
	PendingWriteNum        = 2000
	MaxMsgLen       uint32 = 4096
	HTTPTimeout            = 10 * time.Second
	LenMsgLen              = 2
	LittleEndian           = false

	// skeleton conf
	GoLen              = 10000
	TimerDispatcherLen = 10000
	AsynCallLen        = 10000
	ChanRPCLen         = 10000
)

func InitWebTool() {
	v, err := initViper()
	if err != nil {
		log.Errorf("err is %v", err)
		panic(err)
	}
	Global, err = webtool.NewWebTool(v, webtool.OptionWithOrm())

	if err != nil {
		log.Errorf("err:%v", err)
		panic(err)
	}
}

func initViper() (*viper.Viper, error) {
	viper.SetConfigName("application")                                                                         // name of config file (without extension)
	viper.SetConfigType("yaml")                                                                                // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("/etc/work/")                                                                          // path to look for the config file in
	viper.AddConfigPath("./")                                                                                  // optionally look for config in the working directory
	viper.AddConfigPath("./ddz/workers/conf")                                                                  // optionally look for config in the working directory
	viper.AddConfigPath("../cmd/")                                                                             // optionally look for config in the working directory
	viper.AddConfigPath("../resource/")                                                                        // optionally look for config in the working directory
	viper.AddConfigPath("/Users/zhangjianjun/work/lb/github.com/oldbai555/bgg/service/lbddzserver/impl/conf/") // optionally look for config in the working directory
	err := viper.ReadInConfig()                                                                                // Find and read the config file
	if err != nil {                                                                                            // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	return viper.GetViper(), nil
}
