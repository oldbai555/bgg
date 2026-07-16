// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package permission

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iamclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type PermissionDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPermissionDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PermissionDeleteLogic {
	return &PermissionDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PermissionDeleteLogic) PermissionDelete(req *types.PermissionDeleteReq) error {
	if req.Id == 0 {
		return errs.New(errs.CodeBadRequest, "权限ID不能为空")
	}

	_, err := l.svcCtx.IamRPC.PermissionDelete(l.ctx, &iamclient.PermissionDeleteRequest{Id: req.Id})
	if err != nil {
		return errs.WrapGRPCError("删除权限失败", err)
	}
	return nil
}
