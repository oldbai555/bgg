package internal

import (
	"github.com/name5566/leaf/gate"
	"github.com/oldbai555/bgg/service/lbddz"
	"github.com/oldbai555/bgg/service/lbddzserver/impl/service/mgr"
	"github.com/oldbai555/bgg/service/lbddzserver/impl/service/workers"
	"github.com/oldbai555/lbtool/log"
)

func init() {
	skeleton.RegisterChanRPC("NewAgent", rpcNewAgent)
	skeleton.RegisterChanRPC("CloseAgent", rpcCloseAgent)
}

// agent 被创建时
func rpcNewAgent(args []interface{}) {
	a := args[0].(gate.Agent)
	log.Infof("new agent local addr %s, remote addr %s", a.LocalAddr(), a.RemoteAddr())
}

// agent 被关闭时
func rpcCloseAgent(args []interface{}) {
	a := args[0].(gate.Agent)

	playerId := mgr.AgentMgr.GetPlayerId(a)

	// 更新下线状态
	workers.OrmWorker.Send(lbddz.OrmConsumeTypePlayerLogout, playerId, a.RemoteAddr().String())

	// 移除
	mgr.AgentMgr.Del(a)

	log.Infof("close agent local addr %s, remote addr %s", a.LocalAddr(), a.RemoteAddr())
}
