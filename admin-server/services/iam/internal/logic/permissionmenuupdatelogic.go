package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type PermissionMenuUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPermissionMenuUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PermissionMenuUpdateLogic {
	return &PermissionMenuUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PermissionMenuUpdateLogic) PermissionMenuUpdate(in *iam.PermissionMenuUpdateRequest) (*iam.Empty, error) {
	if in.PermissionId == 0 {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "权限ID不能为空"))
	}

	if err := l.svcCtx.Domain.IAM.RBAC.UpdatePermissionMenus(l.ctx, in.PermissionId, in.MenuIds); err != nil {
		return nil, toGRPCStatus(err)
	}
	return &iam.Empty{}, nil
}
