package iface

import (
	"context"
	"fmt"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/lberr"
	"github.com/oldbai555/lbtool/pkg/routine"
	"sync"
)

var _ IWorker = (*BaseWorker)(nil)

type BaseWorker struct {
	c   chan interface{}
	svr string

	fnM map[ConsumeType]DoHandlerFn
	sync.Mutex
}

func (e *BaseWorker) Send(msgType ConsumeType, is ...interface{}) {
	e.SendItem(NewBaseItem(msgType, is...))
}

func NewBaseWorker(size int, svr string) IWorker {
	return &BaseWorker{
		c:   make(chan interface{}, size),
		svr: svr,
	}
}

func (e *BaseWorker) Start(ctx context.Context, mgr IHandlerMgr) {
	routine.Go(ctx, func(ctx context.Context) error {
		log.Infof(fmt.Sprintf("============ %s starting service ============", e.svr))
		defer func() {
			log.Infof(fmt.Sprintf("============ %s end service ============", e.svr))
		}()
		for {
			select {
			case <-ctx.Done():
				log.Infof(fmt.Sprintf("============ %s service done ============", e.svr))
				return nil
			default:
				receive := e.Receive()
				item, ok := receive.(*BaseItem)
				if !ok {
					log.Warnf("unsupported type")
					continue
				}
				err := mgr.Call(ctx, item.GetMsgType(), item.GetParams()...)
				if err != nil {
					log.Errorf("err:%v", err)
					continue
				}
			}
		}
	})
}

func (e *BaseWorker) Receive() interface{} {
	return <-e.c
}

func (e *BaseWorker) SendItem(v IItem) {
	e.c <- v
}

// =================================================================

var _ IHandlerMgr = (*BaseHandlerMgr)(nil)

type BaseHandlerMgr struct {
	fnM map[ConsumeType]DoHandlerFn
	sync.Mutex
}

func NewBaseHandlerMgr() *BaseHandlerMgr {
	return &BaseHandlerMgr{
		fnM: make(map[ConsumeType]DoHandlerFn),
	}
}

func (e *BaseHandlerMgr) Register(t ConsumeType, fn DoHandlerFn) {
	e.Lock()
	defer e.Unlock()
	_, ok := e.fnM[t]
	if ok {
		panic(fmt.Sprintf("already registered fn , type is %v", t))
		return
	}
	e.fnM[t] = fn
	return
}

func (e *BaseHandlerMgr) Call(ctx context.Context, t ConsumeType, params ...interface{}) error {
	fn, ok := e.fnM[t]
	if !ok {
		log.Warnf("not found fn , type is %v", t)
		return lberr.NewInvalidArg("not found fn , type is %v", t)
	}
	err := fn(ctx, params...)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

// =================================================================

var _ IItem = (*BaseItem)(nil)

type BaseItem struct {
	MsgType ConsumeType
	Params  []interface{}
}

func NewBaseItem(msgType ConsumeType, params ...interface{}) IItem {
	return &BaseItem{MsgType: msgType, Params: params}
}

func (i *BaseItem) GetMsgType() ConsumeType {
	return i.MsgType
}

func (i *BaseItem) GetParams() []interface{} {
	return i.Params
}
