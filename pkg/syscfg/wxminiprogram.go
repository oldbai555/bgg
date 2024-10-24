package syscfg

import (
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/lberr"
	"github.com/spf13/viper"
)

const defaultApolloWxMiniProgramPrefix = "wxMiniProgram"

type WxMiniProgramConf struct {
	AppId  string `json:"appId"`
	Secret string `json:"secret"`
}

func NewWxMiniProgramConf(viper *viper.Viper) *WxMiniProgramConf {
	var v WxMiniProgramConf
	val := viper.Get(defaultApolloWxMiniProgramPrefix)
	err := JsonConvertStruct(val, &v)
	if err != nil {
		log.Errorf("err is %v", err)
		panic(err)
	}
	return &v
}

func GetWxMiniProgramConf() (*WxMiniProgramConf, error) {
	if Global == nil {
		return nil, lberr.NewInvalidArg("not found Global")
	}
	if Global.WxMiniProgramConf == nil {
		return nil, lberr.NewInvalidArg("not found wx mini program conf")
	}
	return Global.WxMiniProgramConf, nil
}
