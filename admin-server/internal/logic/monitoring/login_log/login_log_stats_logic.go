// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package login_log

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iamclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogStatsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogStatsLogic {
	return &LoginLogStatsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogStatsLogic) LoginLogStats() (resp *types.LoginLogStatsResp, err error) {
	rpcResp, err := l.svcCtx.IamRPC.LoginLogStats(l.ctx, &iamclient.Empty{})
	if err != nil {
		return nil, errs.WrapGRPCError("查询登录统计失败", err)
	}

	return &types.LoginLogStatsResp{
		TotalCount:      rpcResp.TotalCount,
		SuccessCount:    rpcResp.SuccessCount,
		FailureCount:    rpcResp.FailureCount,
		TodayCount:      rpcResp.TodayCount,
		TodaySuccess:    rpcResp.TodaySuccess,
		TodayFailure:    rpcResp.TodayFailure,
		OnlineUserCount: rpcResp.OnlineUserCount,
	}, nil
}
