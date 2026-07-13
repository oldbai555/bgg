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

type MenuTreeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMenuTreeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MenuTreeLogic {
	return &MenuTreeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MenuTreeLogic) MenuTree() (resp *types.MenuTreeResp, err error) {
	rpcResp, err := l.svcCtx.IamRPC.MenuTree(l.ctx, &iamclient.Empty{})
	if err != nil {
		return nil, errs.WrapGRPCError("查询菜单树失败", err)
	}

	return &types.MenuTreeResp{List: convertMenuItems(rpcResp.List)}, nil
}

// convertMenuItems 把 iamclient.MenuItem（Children 是指针切片）转成 types.MenuItem
// （Children 是值切片），goctl api 生成的类型历史上就是值切片形状，这里保持不变。
func convertMenuItems(items []*iamclient.MenuItem) []types.MenuItem {
	result := make([]types.MenuItem, 0, len(items))
	for _, item := range items {
		result = append(result, types.MenuItem{
			Id:             item.Id,
			ParentId:       item.ParentId,
			Name:           item.Name,
			Path:           item.Path,
			Component:      item.Component,
			Icon:           item.Icon,
			MenuType:       item.MenuType,
			OrderNum:       item.OrderNum,
			Visible:        item.Visible,
			Status:         item.Status,
			PermissionCode: item.PermissionCode,
			Children:       convertMenuItems(item.Children),
		})
	}
	return result
}
