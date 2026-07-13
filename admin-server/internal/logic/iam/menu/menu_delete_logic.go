// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package menu

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iamclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type MenuDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMenuDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MenuDeleteLogic {
	return &MenuDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MenuDeleteLogic) MenuDelete(req *types.MenuDeleteReq) error {
	if req.Id == 0 {
		return errs.New(errs.CodeBadRequest, "菜单ID不能为空")
	}

	_, err := l.svcCtx.IamRPC.MenuDelete(l.ctx, &iamclient.MenuDeleteRequest{Id: req.Id})
	if err != nil {
		return errs.WrapGRPCError("删除菜单失败", err)
	}
	return nil
}
