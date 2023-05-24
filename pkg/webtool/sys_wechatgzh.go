package webtool

import (
	"github.com/oldbai555/lbtool/log"
	"github.com/spf13/viper"
)

const defaultWeChatGzhConfPrefix = "wechatConf"

type WxGzhConf struct {
	AppId          string `json:"app_id"`
	AppSecret      string `json:"app_secret"`
	Token          string `json:"token"`
	EncodingAESKey string `json:"encoding_aes_key"`
}

func NewWxGzhConf(viper *viper.Viper) *WxGzhConf {
	var v WxGzhConf
	val := viper.Get(defaultWeChatGzhConfPrefix)
	err := JsonConvertStruct(val, &v)
	if err != nil {
		log.Errorf("err is %v", err)
		panic(err)
	}
	return &v
}
