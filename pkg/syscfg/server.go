package syscfg

import (
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/lberr"
	"github.com/spf13/viper"
)

const defaultApolloServerPrefix = "server"

type ServerConf struct {
	Name     string `json:"name"`
	Port     uint32 `json:"port"`
	Ip       string `json:"ip"`
	GatePort uint32 `json:"gatePort"`
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

func GetServerConf() (*ServerConf, error) {
	if Global == nil {
		return nil, lberr.NewInvalidArg("not found Global")
	}
	if Global.ServerConf == nil {
		return nil, lberr.NewInvalidArg("not found server conf")
	}
	return Global.ServerConf, nil
}
