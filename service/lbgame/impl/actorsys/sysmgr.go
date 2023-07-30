package actorsys

import (
	"github.com/oldbai555/bgg/service/lbgame/engine"
	"github.com/oldbai555/bgg/service/lbgame/iface"
)

type SysMgr struct {
	owner iface.IActor
	m     map[iface.SystemType]iface.ISystem
}

func NewSysMgr(actor iface.IActor) *SysMgr {
	return &SysMgr{owner: actor}
}

func (m *SysMgr) OnInit() {
	for tye, fn := range engine.SysSet {
		m.m[tye] = fn()
	}
}

func (s *SysMgr) GetAllSys() []iface.ISystem {
	var list []iface.ISystem
	for systemType := range s.m {
		list = append(list, s.m[systemType])
	}
	return list
}

func (s *SysMgr) GetSys(tye iface.SystemType) (iface.ISystem, bool) {
	system, ok := s.m[tye]
	return system, ok
}
