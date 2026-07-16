package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type NotificationClearReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewNotificationClearReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NotificationClearReadLogic {
	return &NotificationClearReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *NotificationClearReadLogic) NotificationClearRead(in *iam.NotificationClearReadRequest) (*iam.Empty, error) {
	if in == nil || in.UserId == 0 {
		return nil, toGRPCStatus(errs.New(errs.CodeUnauthorized, "未登录或登录已过期"))
	}

	if err := l.svcCtx.Domain.System.Notification.ClearRead(l.ctx, in.UserId); err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "清除已读消息失败", err))
	}

	return &iam.Empty{}, nil
}
