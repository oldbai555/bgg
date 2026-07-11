// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package permission_menu

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"

	"github.com/zeromicro/go-zero/core/logx"
)

type PermissionMenuUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPermissionMenuUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PermissionMenuUpdateLogic {
	return &PermissionMenuUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PermissionMenuUpdateLogic) PermissionMenuUpdate(req *types.PermissionMenuUpdateReq) error {
	if req.PermissionId == 0 {
		return errs.New(errs.CodeBadRequest, "权限ID不能为空")
	}

	return l.svcCtx.Domain.IAM.RBAC.UpdatePermissionMenus(l.ctx, req.PermissionId, req.MenuIds)
}
