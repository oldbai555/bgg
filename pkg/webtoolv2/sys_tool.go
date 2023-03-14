package webtool

import (
	"encoding/json"
	"github.com/oldbai555/lbtool/log"
	"github.com/spf13/viper"
)

// WebTool 目的 在项目运行中各种中间件都能在此处获取
type WebTool struct {
	V *viper.Viper

	GormMysqlConf *GormMysqlConf
	RedisConf     *RedisConf
	StorageConf   *StorageConf
	ServerConf    *ServerConf
	WxGzhConf     *WxGzhConf
}

var lb *WebTool

func NewWebTool(viper *viper.Viper, option ...Option) (*WebTool, error) {
	lb = &WebTool{
		V: viper,
	}
	option = append(option, OptionWithServer())
	// 初始化组件
	for _, o := range option {
		o(lb)
	}
	return lb, nil
}

func JsonConvertStruct(re interface{}, out interface{}) error {
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
