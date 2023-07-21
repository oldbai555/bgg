package lbddz

func (m *ModelGamePlayer) RemoveCards(cards []uint32) (canRemove, hasOut bool) {
	var haveHard bool
	for i, length := 0, len(cards); i < length; i++ {
		haveHard = false

		for j := 0; j < len(m.CurCards); j++ {
			if cards[i] == m.CurCards[j] { // 清除对应的牌
				haveHard = true
				m.CurCards[j] = 0
				break
			}
		}

		if !haveHard {
			return
		}
	}

	// 检查是否出完牌了
	hasOut = true
	for j := 0; j < len(m.CurCards); j++ {
		if 0 != m.CurCards[j] {
			hasOut = false
			break
		}
	}

	return
}
