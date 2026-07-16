// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package monitor

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iamclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type MonitorStatsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMonitorStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MonitorStatsLogic {
	return &MonitorStatsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MonitorStatsLogic) MonitorStats() (resp *types.MonitorStatsResp, err error) {
	rpcResp, err := l.svcCtx.IamRPC.MonitorStats(l.ctx, &iamclient.Empty{})
	if err != nil {
		return nil, errs.WrapGRPCError("查询监控统计失败", err)
	}

	return &types.MonitorStatsResp{
		UserCount:         rpcResp.UserCount,
		RoleCount:         rpcResp.RoleCount,
		PermissionCount:   rpcResp.PermissionCount,
		MenuCount:         rpcResp.MenuCount,
		OnlineUserCount:   rpcResp.OnlineUserCount,
		OperationLogCount: rpcResp.OperationLogCount,
		LoginLogCount:     rpcResp.LoginLogCount,
	}, nil
}
