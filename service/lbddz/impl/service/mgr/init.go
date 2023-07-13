package mgr

import (
	"github.com/name5566/leaf/gate"
	"sync"
)

var RoomMgr = newRoomMgr(new(sync.Map))

var PlayerMgr = newPlayerMgr(new(sync.Map))

var GameMgr = newGameMgr(new(sync.Map))

var AgentMgr = &agentMgr{
	ap: make(map[gate.Agent]uint64),
	pa: make(map[uint64]gate.Agent),
}
