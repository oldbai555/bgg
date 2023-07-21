package cache

import (
	"fmt"
	"github.com/oldbai555/bgg/service/lbddz"
)

type KeyBuilder struct {
	Svr    string
	Prefix string
}

func NewKeyBuilder(svr string, prefix string) *KeyBuilder {
	return &KeyBuilder{Svr: svr, Prefix: prefix}
}

func (m *KeyBuilder) BuildKey(key interface{}) string {
	return fmt.Sprintf("%s:%s:%v", m.Svr, m.Prefix, key)
}

var (
	RoomKeyBuilder   = NewKeyBuilder(lbddz.ServerName, "room")
	PlayerKeyBuilder = NewKeyBuilder(lbddz.ServerName, "player")
)
