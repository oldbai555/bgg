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

type UserRoleUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserRoleUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserRoleUpdateLogic {
	return &UserRoleUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserRoleUpdateLogic) UserRoleUpdate(req *types.UserRoleUpdateReq) error {
	if req.UserId == 0 {
		return errs.New(errs.CodeBadRequest, "用户ID不能为空")
	}

	_, err := l.svcCtx.IamRPC.UserRoleUpdate(l.ctx, &iamclient.UserRoleUpdateRequest{
		UserId:  req.UserId,
		RoleIds: req.RoleIds,
	})
	if err != nil {
		return errs.WrapGRPCError("更新用户角色失败", err)
	}
	return nil
}
