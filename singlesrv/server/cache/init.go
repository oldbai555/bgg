package cache

import (
	"github.com/oldbai555/bgg/pkg/syscfg"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/lberr"
	"github.com/oldbai555/micro/bredis"
)

var rdb *bredis.Group

func Rdb() (*bredis.Group, error) {
	if rdb == nil {
		return nil, lberr.NewInvalidArg("not found rdb")
	}
	return rdb, nil
}

func InitCache() (err error) {
	r := syscfg.NewRedisConf("")
	rdb, err = bredis.New(r.Host, r.Port, r.Password)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

func IsNotFound(err error) bool {
	rdb, err1 := Rdb()
	if err1 != nil {
		return false
	}
	return rdb.IsNotFound(err)
}
