package mgr

import (
	"github.com/oldbai555/bgg/client/lbddz"
	"sync"
)

type roomMgr struct {
	m *sync.Map // map[uint64]*entity.Room
}

func newRoomMgr(m *sync.Map) *roomMgr {
	return &roomMgr{m: m}
}

func (m *roomMgr) Set(v *lbddz.ModelRoom) {
	m.m.Store(v.Id, v)
}

func (m *roomMgr) Del(id uint64) {
	m.m.Delete(id)
}

func (m *roomMgr) Get(id uint64) (*lbddz.ModelRoom, bool) {
	value, ok := m.m.Load(id)
	if !ok {
		return nil, false
	}
	v, ok := value.(*lbddz.ModelRoom)
	if !ok {
		return nil, false
	}
	return v, true
}
