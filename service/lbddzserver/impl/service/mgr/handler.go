package mgr

import (
	"github.com/oldbai555/bgg/service/lbddz/iface"
)

var OrmHandlerMgr = iface.NewBaseHandlerMgr()
var GameHandlerMgr = iface.NewBaseHandlerMgr()
var WebhookHandlerMgr = iface.NewBaseHandlerMgr()