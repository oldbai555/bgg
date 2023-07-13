package tool

import (
	"math/rand"
	"time"
)

func GetRandomNum() uint32 {
	//将时间戳设置成种子数
	rand.Seed(time.Now().UnixNano())
	return uint32(rand.Intn(1<<32 - 1))
}
