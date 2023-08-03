package gate

import (
	"github.com/oldbai555/bgg/service/lbddz"
	"github.com/oldbai555/bgg/service/lbddzserver/impl/moude/game"
	"github.com/oldbai555/bgg/service/lbddzserver/impl/moude/login"
	"github.com/oldbai555/bgg/service/lbddzserver/impl/msg"
)

// 注册路由

func init() {
	msg.Processor.SetRouter(&lbddz.Register{}, login.ChanRPC)
	msg.Processor.SetRouter(&lbddz.Login{}, login.ChanRPC)
	msg.Processor.SetRouter(&lbddz.Event{}, game.ChanRPC)
}