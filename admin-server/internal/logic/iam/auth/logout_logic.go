// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package auth

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iamclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type LogoutLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutLogic {
	return &LogoutLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LogoutLogic) Logout(req *types.LogoutReq) error {
	if req == nil {
		return errs.New(errs.CodeBadRequest, "请求参数不能为空")
	}

	_, err := l.svcCtx.IamRPC.Logout(l.ctx, &iamclient.LogoutRequest{
		AccessToken:  req.AccessToken,
		RefreshToken: req.RefreshToken,
	})
	if err != nil {
		return errs.WrapGRPCError("登出失败", err)
	}
	return nil
}
