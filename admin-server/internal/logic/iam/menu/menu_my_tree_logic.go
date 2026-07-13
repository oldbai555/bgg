// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package menu

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	jwthelper "postapocgame/admin-server/pkg/jwt"
	"postapocgame/admin-server/services/iam/iamclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type MenuMyTreeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMenuMyTreeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MenuMyTreeLogic {
	return &MenuMyTreeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MenuMyTreeLogic) MenuMyTree() (resp *types.MenuTreeResp, err error) {
	user, ok := jwthelper.FromContext(l.ctx)
	if !ok {
		return nil, errs.New(errs.CodeUnauthorized, "未登录或登录已过期")
	}

	rpcResp, err := l.svcCtx.IamRPC.MenuMyTree(l.ctx, &iamclient.MenuMyTreeRequest{UserId: user.UserID})
	if err != nil {
		return nil, errs.WrapGRPCError("查询菜单树失败", err)
	}

	return &types.MenuTreeResp{List: convertMenuItems(rpcResp.List)}, nil
}
