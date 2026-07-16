package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type PermissionApiUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPermissionApiUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PermissionApiUpdateLogic {
	return &PermissionApiUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PermissionApiUpdateLogic) PermissionApiUpdate(in *iam.PermissionApiUpdateRequest) (*iam.Empty, error) {
	if in.PermissionId == 0 {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "权限ID不能为空"))
	}

	if err := l.svcCtx.Domain.IAM.RBAC.UpdatePermissionApis(l.ctx, in.PermissionId, in.ApiIds); err != nil {
		return nil, toGRPCStatus(err)
	}
	return &iam.Empty{}, nil
}
