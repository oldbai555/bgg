package internal

import (
	"context"
	"github.com/name5566/leaf/gate"
	"github.com/oldbai555/bgg/client/lbddz"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/utils"
	"reflect"
)

func handleMsg(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func init() {
	handleMsg(&lbddz.Event{}, handleEvent)
}

func handleEvent(args []interface{}) {
	ctx := context.Background()
	log.SetLogHint(utils.GenUUID())

	m := args[0].(*lbddz.Event)
	a := args[1].(gate.Agent)
	log.Infof("event msg is %v", m)
	fn, ok := fnMap[m.Type]
	if !ok {
		log.Errorf("not found event handle")
		return
	}
	fn(ctx, a, m)
}
