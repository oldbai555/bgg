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

type ProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProfileLogic {
	return &ProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProfileLogic) Profile() (resp *types.ProfileResp, err error) {
	user, ok := jwthelper.FromContext(l.ctx)
	if !ok {
		return nil, errs.New(errs.CodeUnauthorized, "未登录或登录已过期")
	}

	rpcResp, err := l.svcCtx.IamRPC.Profile(l.ctx, &iamclient.ProfileRequest{UserId: user.UserID})
	if err != nil {
		return nil, errs.WrapGRPCError("获取个人信息失败", err)
	}

	return &types.ProfileResp{
		Id:          rpcResp.Id,
		Username:    rpcResp.Username,
		Nickname:    rpcResp.Nickname,
		Avatar:      rpcResp.Avatar,
		Signature:   rpcResp.Signature,
		Permissions: rpcResp.Permissions,
	}, nil
}
