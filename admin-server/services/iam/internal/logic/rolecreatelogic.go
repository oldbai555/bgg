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

type RoleCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRoleCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleCreateLogic {
	return &RoleCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RoleCreateLogic) RoleCreate(in *iam.RoleCreateRequest) (*iam.Empty, error) {
	if in == nil || in.Name == "" || in.Code == "" {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "角色名称和编码不能为空"))
	}

	role := iammodel.AdminRole{
		Name:        in.Name,
		Code:        in.Code,
		Description: sql.NullString{String: in.Description, Valid: in.Description != ""},
		Status:      in.Status,
	}

	if err := l.svcCtx.Domain.IAM.Role.Create(l.ctx, &role); err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "创建角色失败", err))
	}
	return &iam.Empty{}, nil
}
