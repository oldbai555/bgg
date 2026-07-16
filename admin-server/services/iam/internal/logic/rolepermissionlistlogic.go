package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type RolePermissionListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRolePermissionListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RolePermissionListLogic {
	return &RolePermissionListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RolePermissionListLogic) RolePermissionList(in *iam.RolePermissionListRequest) (*iam.RolePermissionListResponse, error) {
	if in.RoleId == 0 {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "角色ID不能为空"))
	}

	if _, err := l.svcCtx.Domain.IAM.Role.FindByID(l.ctx, in.RoleId); err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeBadRequest, "角色不存在", err))
	}

	permissionIDs, err := l.svcCtx.Domain.IAM.RolePermission.ListPermissionIDsByRoleID(l.ctx, in.RoleId)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "查询角色权限失败", err))
	}

	return &iam.RolePermissionListResponse{PermissionIds: permissionIDs}, nil
}
