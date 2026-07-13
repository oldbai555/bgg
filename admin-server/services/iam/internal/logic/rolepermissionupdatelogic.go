package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type RolePermissionUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRolePermissionUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RolePermissionUpdateLogic {
	return &RolePermissionUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RolePermissionUpdateLogic) RolePermissionUpdate(in *iam.RolePermissionUpdateRequest) (*iam.Empty, error) {
	if in.RoleId == 0 {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "角色ID不能为空"))
	}

	if err := l.svcCtx.Domain.IAM.RBAC.UpdateRolePermissions(l.ctx, in.RoleId, in.PermissionIds); err != nil {
		return nil, toGRPCStatus(err)
	}

	cache := l.svcCtx.Repository.BusinessCache
	go func() {
		if err := cache.DeleteMenuTree(context.Background()); err != nil {
			l.Errorf("清除菜单树缓存失败: %v", err)
		}
	}()

	return &iam.Empty{}, nil
}
