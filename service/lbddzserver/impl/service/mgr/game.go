package mgr

import (
	"github.com/oldbai555/bgg/service/lbddz"
	"sync"
)

type gameMgr struct {
	m *sync.Map // map[uint64]*lbddz.ModelGame
}

func newGameMgr(m *sync.Map) *gameMgr {
	return &gameMgr{m: m}
}

func (m *gameMgr) Set(v *lbddz.BaseGame) {
	m.m.Store(v.G.Id, v)
}

func (m *gameMgr) Del(id uint64) {
	m.m.Delete(id)
}

func (m *gameMgr) Get(id uint64) (*lbddz.BaseGame, bool) {
	value, ok := m.m.Load(id)
	if !ok {
		return nil, false
	}
	v, ok := value.(*lbddz.BaseGame)
	if !ok {
		return nil, false
	}
	return v, true
}
