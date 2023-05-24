package webtool

import (
	"github.com/oldbai555/lbtool/log"
	"github.com/spf13/viper"
)

const defaultChatGptPrefix = "chatGpt"

type ChatGpt struct {
	Proxy  string `json:"proxy"`
	ApiKey string `json:"api_key"`
}

func NewChatGpt(viper *viper.Viper) *ChatGpt {
	var v ChatGpt
	val := viper.Get(defaultChatGptPrefix)
	err := JsonConvertStruct(val, &v)
	if err != nil {
		log.Errorf("err is %v", err)
		panic(err)
	}
	return &v
}
