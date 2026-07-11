// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package permission_api

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"

	"github.com/zeromicro/go-zero/core/logx"
)

type PermissionApiUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPermissionApiUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PermissionApiUpdateLogic {
	return &PermissionApiUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PermissionApiUpdateLogic) PermissionApiUpdate(req *types.PermissionApiUpdateReq) error {
	if req.PermissionId == 0 {
		return errs.New(errs.CodeBadRequest, "权限ID不能为空")
	}

	return l.svcCtx.Domain.IAM.RBAC.UpdatePermissionApis(l.ctx, req.PermissionId, req.ApiIds)
}
