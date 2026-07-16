// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package notification

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	jwthelper "postapocgame/admin-server/pkg/jwt"
	"postapocgame/admin-server/services/iam/iamclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type NotificationClearReadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewNotificationClearReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NotificationClearReadLogic {
	return &NotificationClearReadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *NotificationClearReadLogic) NotificationClearRead() (resp *types.Response, err error) {
	user, ok := jwthelper.FromContext(l.ctx)
	if !ok {
		return nil, errs.New(errs.CodeUnauthorized, "未登录或登录已过期")
	}

	_, err = l.svcCtx.IamRPC.NotificationClearRead(l.ctx, &iamclient.NotificationClearReadRequest{UserId: user.UserID})
	if err != nil {
		return nil, errs.WrapGRPCError("清除已读消息失败", err)
	}

	return &types.Response{
		Code:    0,
		Message: "操作成功",
	}, nil
}
