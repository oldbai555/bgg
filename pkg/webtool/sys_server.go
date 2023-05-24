package webtool

import (
	"github.com/oldbai555/lbtool/log"
	"github.com/spf13/viper"
)

const defaultApolloServerPrefix = "server"

type ServerConf struct {
	Name string `json:"name"`
	Port int    `json:"port"`
}

func NewServerConf(viper *viper.Viper) *ServerConf {
	var v ServerConf
	val := viper.Get(defaultApolloServerPrefix)
	err := JsonConvertStruct(val, &v)
	if err != nil {
		log.Errorf("err is %v", err)
		panic(err)
	}
	return &v
}
