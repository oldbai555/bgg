// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package role_permission

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"

	"github.com/zeromicro/go-zero/core/logx"
)

type RolePermissionListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRolePermissionListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RolePermissionListLogic {
	return &RolePermissionListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RolePermissionListLogic) RolePermissionList(req *types.RolePermissionListReq) (resp *types.RolePermissionListResp, err error) {
	if req.RoleId == 0 {
		return nil, errs.New(errs.CodeBadRequest, "角色ID不能为空")
	}

	// 验证角色是否存在
	_, err = l.svcCtx.Domain.IAM.Role.FindByID(l.ctx, req.RoleId)
	if err != nil {
		return nil, errs.Wrap(errs.CodeBadRequest, "角色不存在", err)
	}

	permissionIDs, err := l.svcCtx.Domain.IAM.RolePermission.ListPermissionIDsByRoleID(l.ctx, req.RoleId)
	if err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, "查询角色权限失败", err)
	}

	return &types.RolePermissionListResp{
		PermissionIds: permissionIDs,
	}, nil
}
