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

type MenuUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMenuUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MenuUpdateLogic {
	return &MenuUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MenuUpdateLogic) MenuUpdate(req *types.MenuUpdateReq) error {
	if req == nil || req.Id == 0 {
		return errs.New(errs.CodeBadRequest, "菜单ID不能为空")
	}

	_, err := l.svcCtx.IamRPC.MenuUpdate(l.ctx, &iamclient.MenuUpdateRequest{
		Id:        req.Id,
		ParentId:  req.ParentId,
		Name:      req.Name,
		Path:      req.Path,
		Component: req.Component,
		Icon:      req.Icon,
		MenuType:  req.MenuType,
		OrderNum:  req.OrderNum,
		Visible:   req.Visible,
		Status:    req.Status,
	})
	if err != nil {
		return errs.WrapGRPCError("更新菜单失败", err)
	}
	return nil
}
