package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type NotificationDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewNotificationDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NotificationDeleteLogic {
	return &NotificationDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *NotificationDeleteLogic) NotificationDelete(in *iam.NotificationDeleteRequest) (*iam.Empty, error) {
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
		return nil, toGRPCStatus(errs.New(errs.CodeForbidden, "无权删除该消息通知"))
	}

	if err := l.svcCtx.Domain.System.Notification.DeleteByID(l.ctx, in.Id); err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "删除消息通知失败", err))
	}

	return &iam.Empty{}, nil
}
