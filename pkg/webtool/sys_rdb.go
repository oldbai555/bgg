package webtool

import (
	"github.com/oldbai555/lbtool/log"
	"github.com/spf13/viper"
)

const defaultApolloRedisPrefix = "redis"

type RedisConf struct {
	Database int    `json:"database"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
}

func NewRedisConf(viper *viper.Viper) *RedisConf {
	var v RedisConf
	val := viper.Get(defaultApolloRedisPrefix)
	err := JsonConvertStruct(val, &v)
	if err != nil {
		log.Errorf("err is %v", err)
		panic(err)
	}
	return &v
}
