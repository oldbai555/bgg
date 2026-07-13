// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package role_permission

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iamclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type RolePermissionUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRolePermissionUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RolePermissionUpdateLogic {
	return &RolePermissionUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RolePermissionUpdateLogic) RolePermissionUpdate(req *types.RolePermissionUpdateReq) error {
	if req.RoleId == 0 {
		return errs.New(errs.CodeBadRequest, "角色ID不能为空")
	}

	_, err := l.svcCtx.IamRPC.RolePermissionUpdate(l.ctx, &iamclient.RolePermissionUpdateRequest{
		RoleId:        req.RoleId,
		PermissionIds: req.PermissionIds,
	})
	if err != nil {
		return errs.WrapGRPCError("更新角色权限失败", err)
	}
	return nil
}
