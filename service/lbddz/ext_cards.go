package lbddz

import (
	"math/rand"
	"time"
)

func GetNewCards54() []uint32 {
	cards := make([]uint32, 54)
	for i := uint32(0); i < 54; i++ {
		cards[i] = i + 1
	}

	rand.Seed(time.Now().UnixNano())

	// 洗牌算法
	for i := 53; i >= 0; i-- {
		j := rand.Intn(i + 1)
		if i != j {
			temp := cards[i]
			cards[i] = cards[j]
			cards[j] = temp
		}
	}

	return cards
}
