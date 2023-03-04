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

func (r *ServerConf) InitConf(viper *viper.Viper) error {
	var v ServerConf
	val := viper.Get(defaultApolloServerPrefix)
	err := JsonConvertStruct(val, &v)
	if err != nil {
		log.Errorf("err is %v", err)
		return err
	}
	log.Infof("init server successfully")
	r.Name = v.Name
	r.Port = v.Port
	return nil
}

func (r *ServerConf) GenConfTool(tool *WebTool) error {
	log.Infof("init rdb engine successfully")
	tool.Sc = r
	return nil
}
