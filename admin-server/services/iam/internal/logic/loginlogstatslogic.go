package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/chat/chatclient"
	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogStatsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogStatsLogic {
	return &LoginLogStatsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogStatsLogic) LoginLogStats(in *iam.Empty) (*iam.LoginLogStatsResponse, error) {
	successCount, err := l.svcCtx.Domain.Monitoring.LoginLog.CountByStatus(l.ctx, 1)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "查询成功次数失败", err))
	}
	failureCount, err := l.svcCtx.Domain.Monitoring.LoginLog.CountByStatus(l.ctx, 2)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "查询失败次数失败", err))
	}
	totalCount := successCount + failureCount

	todayCount, err := l.svcCtx.Domain.Monitoring.LoginLog.CountToday(l.ctx)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "查询今日登录次数失败", err))
	}

	todaySuccess, err := l.svcCtx.Domain.Monitoring.LoginLog.CountTodayByStatus(l.ctx, 1)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "查询今日成功次数失败", err))
	}

	todayFailure, err := l.svcCtx.Domain.Monitoring.LoginLog.CountTodayByStatus(l.ctx, 2)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "查询今日失败次数失败", err))
	}

	onlineUserCount := int64(0)
	if onlineResp, err := l.svcCtx.ChatRPC.GetOnlineUserCount(l.ctx, &chatclient.Empty{}); err == nil {
		onlineUserCount = onlineResp.Count
	}

	return &iam.LoginLogStatsResponse{
		TotalCount:      totalCount,
		SuccessCount:    successCount,
		FailureCount:    failureCount,
		TodayCount:      todayCount,
		TodaySuccess:    todaySuccess,
		TodayFailure:    todayFailure,
		OnlineUserCount: onlineUserCount,
	}, nil
}
