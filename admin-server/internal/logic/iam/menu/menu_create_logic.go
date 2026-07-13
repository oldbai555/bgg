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

type MenuCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMenuCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MenuCreateLogic {
	return &MenuCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MenuCreateLogic) MenuCreate(req *types.MenuCreateReq) error {
	if req == nil || req.Name == "" || req.MenuType == 0 {
		return errs.New(errs.CodeBadRequest, "菜单名称和类型不能为空")
	}

	_, err := l.svcCtx.IamRPC.MenuCreate(l.ctx, &iamclient.MenuCreateRequest{
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
		return errs.WrapGRPCError("创建菜单失败", err)
	}
	return nil
}
