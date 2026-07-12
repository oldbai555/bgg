// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package login_log

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/chat/chatclient"

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
	// 总登录次数（查询所有状态，使用一个大的查询）
	// 字典值（login_status）：1=成功，2=失败
	var totalCount int64
	// 由于 CountByStatus 不支持 -1，我们分别查询成功和失败，然后相加
	successCount, err := l.svcCtx.Domain.Monitoring.LoginLog.CountByStatus(l.ctx, 1)
	if err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, "查询成功次数失败", err)
	}
	failureCount, err := l.svcCtx.Domain.Monitoring.LoginLog.CountByStatus(l.ctx, 2)
	if err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, "查询失败次数失败", err)
	}
	totalCount = successCount + failureCount

	// 今日登录次数
	todayCount, err := l.svcCtx.Domain.Monitoring.LoginLog.CountToday(l.ctx)
	if err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, "查询今日登录次数失败", err)
	}

	// 今日成功次数
	todaySuccess, err := l.svcCtx.Domain.Monitoring.LoginLog.CountTodayByStatus(l.ctx, 1)
	if err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, "查询今日成功次数失败", err)
	}

	// 今日失败次数
	todayFailure, err := l.svcCtx.Domain.Monitoring.LoginLog.CountTodayByStatus(l.ctx, 2)
	if err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, "查询今日失败次数失败", err)
	}

	// 当前在线用户数（chat 域拆分成独立服务后，从 ChatRPC.GetOnlineUserCount 获取，
	// 不再直接读 ChatHub——连接表已经搬进 chat-rpc 自己的进程）
	onlineUserCount := int64(0)
	if onlineResp, err := l.svcCtx.ChatRPC.GetOnlineUserCount(l.ctx, &chatclient.Empty{}); err == nil {
		onlineUserCount = onlineResp.Count
	}

	return &types.LoginLogStatsResp{
		TotalCount:      totalCount,
		SuccessCount:    successCount,
		FailureCount:    failureCount,
		TodayCount:      todayCount,
		TodaySuccess:    todaySuccess,
		TodayFailure:    todayFailure,
		OnlineUserCount: onlineUserCount,
	}, nil
}
