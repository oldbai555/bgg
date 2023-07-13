package mgr

import (
	"github.com/oldbai555/bgg/client/lbddz"
	"sync"
)

type playerMgr struct {
	m *sync.Map // map[uint64]*entity.Player
}

func newPlayerMgr(m *sync.Map) *playerMgr {
	return &playerMgr{m: m}
}

func (m *playerMgr) Set(v *lbddz.ModelPlayer) {
	m.m.Store(v.Id, v)
}

func (m *playerMgr) Del(id uint64) {
	m.m.Delete(id)
}

func (m *playerMgr) Get(id uint64) (*lbddz.ModelPlayer, bool) {
	value, ok := m.m.Load(id)
	if !ok {
		return nil, false
	}
	v, ok := value.(*lbddz.ModelPlayer)
	if !ok {
		return nil, false
	}
	return v, true
}

func (m *playerMgr) LoadPlayerInfo(v *lbddz.ModelPlayer) {
	// 加载更新玩家信息
	m.Set(v)

	// 可以推送信息给客户端更新
}
