package lbddz

import "github.com/oldbai555/lbtool/utils"

func (m *ModelGame) AddCurIndex() {
	m.LastPlayerSeq = m.CurPlayerSeq
	m.CurPlayerSeq++
	if m.CurPlayerSeq > 3 {
		m.CurPlayerSeq = 1 //每次到4就变回1
	}
}

func (m *ModelGame) CheckCurPlayerSeq(seq uint32) bool {
	return m.CurPlayerSeq == seq
}

func (m *BaseGame) GetGamePlayer(pId uint64) (*ModelGamePlayer, error) {
	for i := 0; i < len(m.Gps); i++ {
		if m.Gps[i].PlayerId == pId {
			return m.Gps[i], nil
		}
	}
	return nil, ErrGamePlayerNotFound
}

func (m *BaseGame) GetGamePlayerIds() []uint64 {
	return utils.PluckUint64List(m.Gps, FieldPlayerId)
}

func (m *BaseGame) ResetSave(g *ModelGame, gps []*ModelGamePlayer) {
	m.G = g
	m.Gps = gps
}
