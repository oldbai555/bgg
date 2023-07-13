package mgr

import "github.com/name5566/leaf/gate"

type agentMgr struct {
	ap map[gate.Agent]uint64
	pa map[uint64]gate.Agent
}

func (m *agentMgr) Set(a gate.Agent, playerId uint64) {
	m.pa[playerId] = a
	m.ap[a] = playerId

}

func (m *agentMgr) GetAgent(playerId uint64) gate.Agent {
	return m.pa[playerId]
}

func (m *agentMgr) GetAgentList(pIds ...uint64) []gate.Agent {
	var as []gate.Agent
	for _, pId := range pIds {
		as = append(as, m.GetAgent(pId))
	}
	return as
}

func (m *agentMgr) GetPlayerId(a gate.Agent) uint64 {
	return m.ap[a]
}

func (m *agentMgr) Del(a gate.Agent) {
	id := m.GetPlayerId(a)
	delete(m.pa, id)
	delete(m.ap, a)
}
