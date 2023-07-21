package cc

import (
	"context"
	"github.com/petermattis/goid"
)

func CtxWithConnId(ctx context.Context, connId uint64) context.Context {
	goId := goid.Get()
	return context.WithValue(ctx, goId, connId)
}

func GetCtxConnId(ctx context.Context) uint64 {
	goId := goid.Get()
	connId, ok := ctx.Value(goId).(uint64)
	if !ok {
		return 0
	}
	return connId
}
