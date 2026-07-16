package logic

import (
	"context"
	"database/sql"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type PermissionUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPermissionUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PermissionUpdateLogic {
	return &PermissionUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PermissionUpdateLogic) PermissionUpdate(in *iam.PermissionUpdateRequest) (*iam.Empty, error) {
	if in == nil || in.Id == 0 {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "权限ID不能为空"))
	}

	p, err := l.svcCtx.Domain.IAM.Permission.FindByID(l.ctx, in.Id)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "查询权限失败", err))
	}

	p.Name = in.Name
	p.Description = sql.NullString{String: in.Description, Valid: in.Description != ""}

	if err := l.svcCtx.Domain.IAM.Permission.Update(l.ctx, p); err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "更新权限失败", err))
	}
	return &iam.Empty{}, nil
}
