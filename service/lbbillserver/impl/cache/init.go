package cache

import (
	"github.com/oldbai555/bgg/pkg/redis"
	"github.com/oldbai555/bgg/pkg/syscfg"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/lberr"
)

var rdb *redis.Group

func Rdb() (*redis.Group, error) {
	if rdb == nil {
		return nil, lberr.NewInvalidArg("not found rdb")
	}
	return rdb, nil
}

func InitCache() (err error) {
	r := syscfg.NewRedisConf("")
	rdb, err = redis.New(r.Host, r.Port, r.Password)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}
