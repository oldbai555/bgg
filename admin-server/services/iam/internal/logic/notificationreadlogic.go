package logic

import (
	"context"
	"time"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type NotificationReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewNotificationReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NotificationReadLogic {
	return &NotificationReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *NotificationReadLogic) NotificationRead(in *iam.NotificationReadRequest) (*iam.Empty, error) {
	if in == nil || in.Id == 0 {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "请求参数不能为空"))
	}
	if in.UserId == 0 {
		return nil, toGRPCStatus(errs.New(errs.CodeUnauthorized, "未登录或登录已过期"))
	}

	notification, err := l.svcCtx.Domain.System.Notification.FindByID(l.ctx, in.Id)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeNotFound, "消息通知不存在", err))
	}
	if notification.UserId != in.UserId {
		return nil, toGRPCStatus(errs.New(errs.CodeForbidden, "无权操作该消息通知"))
	}

	now := time.Now().Unix()
	notification.ReadStatus = 2
	notification.ReadAt = now
	notification.UpdatedAt = now

	if err := l.svcCtx.Domain.System.Notification.Update(l.ctx, notification); err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "标记已读失败", err))
	}

	return &iam.Empty{}, nil
}
