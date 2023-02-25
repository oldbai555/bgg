package webtool

import (
	"fmt"
	"github.com/go-redis/redis/v8"
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

func (r *RedisConf) InitConf(viper *viper.Viper) error {
	var v RedisConf
	val := viper.Get(defaultApolloRedisPrefix)
	err := jsonConvertStruct(val, &v)
	if err != nil {
		log.Errorf("err is %v", err)
		return err
	}
	log.Infof("init redis successfully")
	r.Host = v.Host
	r.Database = v.Database
	r.Password = v.Password
	r.Port = v.Port
	return nil
}

func (r *RedisConf) GenConfTool(tool *WebTool) error {
	log.Infof("init rdb engine successfully")
	tool.Rdb = &ProxyRdb{
		Client: redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", r.Host, r.Port),
			Password: r.Password,
			DB:       r.Database,
		}),
	}
	return nil
}
