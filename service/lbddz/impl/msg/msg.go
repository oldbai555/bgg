package msg

import (
	"github.com/oldbai555/bgg/client/lbddz"
	"github.com/oldbai555/bgg/service/lbddz/impl/json"
)

// 消息注册使用

var Processor *json.Processor

func init() {
	Processor = json.NewProcessor()
	for i := range msgList {
		Processor.Register(msgList[i])
	}
}

var msgList = []interface{}{
	&lbddz.Event{},
	&lbddz.Webhook{},
	&lbddz.Register{},
	&lbddz.Login{},
}
