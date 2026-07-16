package logic

import (
	"context"
	"database/sql"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iam"
	iammodel "postapocgame/admin-server/services/iam/internal/model/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type PermissionCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPermissionCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PermissionCreateLogic {
	return &PermissionCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PermissionCreateLogic) PermissionCreate(in *iam.PermissionCreateRequest) (*iam.Empty, error) {
	if in == nil || in.Name == "" || in.Code == "" {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "权限名称和编码不能为空"))
	}

	p := iammodel.AdminPermission{
		Name:        in.Name,
		Code:        in.Code,
		Description: sql.NullString{String: in.Description, Valid: in.Description != ""},
	}
	if err := l.svcCtx.Domain.IAM.Permission.Create(l.ctx, &p); err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "创建权限失败", err))
	}
	return &iam.Empty{}, nil
}
