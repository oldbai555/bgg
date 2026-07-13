package logic

import (
	"context"
	"database/sql"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoleUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRoleUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleUpdateLogic {
	return &RoleUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RoleUpdateLogic) RoleUpdate(in *iam.RoleUpdateRequest) (*iam.Empty, error) {
	if in == nil || in.Id == 0 {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "角色ID不能为空"))
	}

	role, err := l.svcCtx.Domain.IAM.Role.FindByID(l.ctx, in.Id)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "查询角色失败", err))
	}

	role.Name = in.Name
	if in.Description != "" {
		role.Description = sql.NullString{String: in.Description, Valid: true}
	}
	if in.Status == 0 || in.Status == 1 {
		role.Status = in.Status
	}

	if err := l.svcCtx.Domain.IAM.Role.Update(l.ctx, role); err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "更新角色失败", err))
	}
	return &iam.Empty{}, nil
}
