// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package user_role

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iamclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserRoleListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserRoleListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserRoleListLogic {
	return &UserRoleListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserRoleListLogic) UserRoleList(req *types.UserRoleListReq) (resp *types.UserRoleListResp, err error) {
	if req.UserId == 0 {
		return nil, errs.New(errs.CodeBadRequest, "用户ID不能为空")
	}

	rpcResp, err := l.svcCtx.IamRPC.UserRoleList(l.ctx, &iamclient.UserRoleListRequest{UserId: req.UserId})
	if err != nil {
		return nil, errs.WrapGRPCError("查询用户角色失败", err)
	}

	return &types.UserRoleListResp{RoleIds: rpcResp.RoleIds}, nil
}
