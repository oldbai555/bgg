package webtool

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/oldbai555/gorm"
	"github.com/oldbai555/lbtool/extpkg/lbconf/bconf"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/storage"
	"time"
)

// WebTool 目的 在项目运行中各种中间件都能在此处获取
type WebTool struct {
	ApoC    bconf.Config
	Orm     *gorm.DB
	Rdb     *ProxyRdb
	Storage storage.FileStorageInterface
}

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

// NewWebTool 只支持 apollo
func NewWebTool(conf *ApolloConf, option ...Option) (*WebTool, error) {
	var err error
	lb := &WebTool{}

	// 初始化 apollo 配置中心
	apollo, err := initApollo(conf)
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}
	lb.ApoC = apollo

	// 初始化组件
	for _, o := range option {
		o(lb)
	}
	return lb, nil
}

func (s *WebTool) GetJson4Apollo(key string, out interface{}) error {
	re, err := s.ApoC.Get(key)
	if err != nil {
		log.Errorf("err is : %v", err)
		return err
	}
	marshal, err := json.Marshal(re)
	if err != nil {
		log.Errorf("err is : %v", err)
		return err
	}
	err = json.Unmarshal(marshal, out)
	if err != nil {
		log.Errorf("err is : %v", err)
		return err
	}
	return nil
}
