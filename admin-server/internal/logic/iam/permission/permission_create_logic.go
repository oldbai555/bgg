// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package permission

import (
	"context"
	"database/sql"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"

	"github.com/zeromicro/go-zero/core/logx"
	"postapocgame/admin-server/internal/model/iam"
	iamrepo "postapocgame/admin-server/internal/repository/iam"
)

type PermissionCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPermissionCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PermissionCreateLogic {
	return &PermissionCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PermissionCreateLogic) PermissionCreate(req *types.PermissionCreateReq) error {
	if req == nil || req.Name == "" || req.Code == "" {
		return errs.New(errs.CodeBadRequest, "权限名称和编码不能为空")
	}

	permissionRepo := iamrepo.NewPermissionRepository(l.svcCtx.Repository)
	p := iam.AdminPermission{
		Name:        req.Name,
		Code:        req.Code,
		Description: sql.NullString{String: req.Description, Valid: req.Description != ""},
	}
	if err := permissionRepo.Create(l.ctx, &p); err != nil {
		return errs.Wrap(errs.CodeInternalError, "创建权限失败", err)
	}
	return nil
}
