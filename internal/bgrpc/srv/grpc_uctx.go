package srv

import (
	"context"
	"github.com/oldbai555/bgg/pkg/uctx"
)

var _ uctx.IUCtx = (*GrpcUCtx)(nil)

func NewGrpcUCtx(ctx context.Context) *GrpcUCtx {
	return &GrpcUCtx{
		Context:  ctx,
		BaseUCtx: uctx.NewBaseUCtx(),
	}
}

type GrpcUCtx struct {
	context.Context
	*uctx.BaseUCtx
}
