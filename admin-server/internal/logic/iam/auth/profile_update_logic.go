// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package auth

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	jwthelper "postapocgame/admin-server/pkg/jwt"
	"postapocgame/admin-server/services/iam/iamclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProfileUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewProfileUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProfileUpdateLogic {
	return &ProfileUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProfileUpdateLogic) ProfileUpdate(req *types.ProfileUpdateReq) error {
	if req == nil {
		return errs.New(errs.CodeBadRequest, "请求参数不能为空")
	}

	user, ok := jwthelper.FromContext(l.ctx)
	if !ok {
		return errs.New(errs.CodeUnauthorized, "未登录或登录已过期")
	}

	_, err := l.svcCtx.IamRPC.ProfileUpdate(l.ctx, &iamclient.ProfileUpdateRequest{
		UserId:    user.UserID,
		Nickname:  req.Nickname,
		Avatar:    req.Avatar,
		Signature: req.Signature,
	})
	if err != nil {
		return errs.WrapGRPCError("更新个人信息失败", err)
	}
	return nil
}
