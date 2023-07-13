package iface

import "context"

type IHandlerMgr interface {
	Register(e ConsumeType, fn DoHandlerFn)
	Call(ctx context.Context, e ConsumeType, params ...interface{}) error
}
