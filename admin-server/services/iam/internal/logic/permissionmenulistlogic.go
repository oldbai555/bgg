package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type PermissionMenuListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPermissionMenuListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PermissionMenuListLogic {
	return &PermissionMenuListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// RBAC 关联关系
func (l *PermissionMenuListLogic) PermissionMenuList(in *iam.PermissionMenuListRequest) (*iam.PermissionMenuListResponse, error) {
	if in.PermissionId == 0 {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "权限ID不能为空"))
	}

	if _, err := l.svcCtx.Domain.IAM.Permission.FindByID(l.ctx, in.PermissionId); err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeBadRequest, "权限不存在", err))
	}

	menuIDs, err := l.svcCtx.Domain.IAM.PermissionMenu.ListMenuIDsByPermissionID(l.ctx, in.PermissionId)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "查询权限菜单失败", err))
	}

	return &iam.PermissionMenuListResponse{MenuIds: menuIDs}, nil
}
