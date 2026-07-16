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

type NotificationReadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewNotificationReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NotificationReadLogic {
	return &NotificationReadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *NotificationReadLogic) NotificationRead(req *types.NotificationReadReq) (resp *types.Response, err error) {
	if req == nil || req.Id == 0 {
		return nil, errs.New(errs.CodeBadRequest, "请求参数不能为空")
	}

	user, ok := jwthelper.FromContext(l.ctx)
	if !ok {
		return nil, errs.New(errs.CodeUnauthorized, "未登录或登录已过期")
	}

	_, err = l.svcCtx.IamRPC.NotificationRead(l.ctx, &iamclient.NotificationReadRequest{
		Id:     req.Id,
		UserId: user.UserID,
	})
	if err != nil {
		return nil, errs.WrapGRPCError("标记已读失败", err)
	}

	return &types.Response{
		Code:    0,
		Message: "操作成功",
	}, nil
}
