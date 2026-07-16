package logic

import (
	"context"

	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type PingLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PingLogic {
	return &PingLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Ping 供 gateway /api/v1/ping 探测 iam-rpc 是否可达（gateway 自己不再直连 MySQL，
// 原来的 "SELECT 1" 探活挪到这里；Redis 探活 gateway 直接查共享 Redis，不经过这个 RPC）。
func (l *PingLogic) Ping(in *iam.Empty) (*iam.PingResponse, error) {
	var result int
	err := l.svcCtx.Repository.DB.QueryRowCtx(l.ctx, &result, "SELECT 1")
	return &iam.PingResponse{Ok: err == nil}, nil
}
