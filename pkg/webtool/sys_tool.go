package webtool

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/oldbai555/gorm"
	"github.com/oldbai555/lbtool/extpkg/lbconf/bconf"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/storage"
	"github.com/spf13/viper"
	"time"
)

// WebTool 目的 在项目运行中各种中间件都能在此处获取
type WebTool struct {
	ApoC    bconf.Config
	Orm     *gorm.DB
	Rdb     *ProxyRdb
	Storage storage.FileStorageInterface
	V       *viper.Viper
	Sc      *ServerConf
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

func initViper() error {
	viper.SetConfigName("application")  // name of config file (without extension)
	viper.SetConfigType("yaml")         // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("/etc/work/")   // path to look for the config file in
	viper.AddConfigPath("./")           // optionally look for config in the working directory
	viper.AddConfigPath("../resource/") // optionally look for config in the working directory
	err := viper.ReadInConfig()         // Find and read the config file
	if err != nil {                     // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	return nil
}

func NewWebTool(option ...Option) (*WebTool, error) {
	err := initViper()
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}
	lb := &WebTool{
		V: viper.GetViper(),
	}
	option = append(option, OptionWithServer())
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
