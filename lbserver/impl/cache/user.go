package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/oldbai555/bgg/client/lbuser"
	webtool "github.com/oldbai555/bgg/pkg/webtoolv2"
	"github.com/oldbai555/lbtool/log"
	utils "github.com/oldbai555/lbtool/pkg/cache_helper"
	"time"
)

var Rdb *ProxyRdb

type ProxyRdb struct {
	*redis.Client
}

func (c *ProxyRdb) SetJson(ctx context.Context, key string, j interface{}, exp time.Duration) error {
	val, err := json.Marshal(j)
	if err != nil {
		log.Errorf("err:%s", err)
		return err
	}
	// 空串这里先不考虑
	if len(val) == 0 {
		return errors.New("unsupported empty value")
	}
	return c.Set(ctx, key, val, exp).Err()
}

func (c *ProxyRdb) GetJson(ctx context.Context, key string, j interface{}) error {
	val, err := c.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return redis.Nil
		}
		log.Errorf("err:%s", err)
		return errors.New("redis exception")
	}
	err = json.Unmarshal(val, j)
	if err != nil {
		log.Errorf("err:%s", err)
		return errors.New("json unmarshal error")
	}
	return nil
}

var UserCache *utils.CacheHelper

func InitCacheHelper(r *webtool.RedisConf) {
	Rdb = &ProxyRdb{
		redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", r.Host, r.Port),
			Password: r.Password,
			DB:       r.Database,
		}),
	}

	UserCache = utils.NewCacheHelper(&utils.NewCacheHelperReq{
		RedisClient: Rdb.Client,
		MType:       &lbuser.ModelUser{},
		FieldNames:  []string{lbuser.FieldId},
	})
}
