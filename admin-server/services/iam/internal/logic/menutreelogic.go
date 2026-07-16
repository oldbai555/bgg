package logic

import (
	"context"
	"sort"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type MenuTreeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMenuTreeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MenuTreeLogic {
	return &MenuTreeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *MenuTreeLogic) MenuTree(in *iam.Empty) (*iam.MenuTreeResponse, error) {
	cache := l.svcCtx.Repository.BusinessCache
	var cachedResp iam.MenuTreeResponse
	if err := cache.GetMenuTree(l.ctx, &cachedResp); err == nil {
		return &cachedResp, nil
	}

	list, err := l.svcCtx.Domain.IAM.Menu.ListAll(l.ctx)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "查询菜单列表失败", err))
	}

	menuPermissionMap, err := l.svcCtx.Domain.IAM.PermissionMenu.ListMenuPermissionCodes(l.ctx)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "查询菜单权限关联失败", err))
	}

	nodeMap := make(map[uint64]*iam.MenuItem, len(list))
	var roots []*iam.MenuItem

	for _, m := range list {
		var permissionCode string
		if codes, ok := menuPermissionMap[m.Id]; ok && len(codes) > 0 {
			permissionCode = codes[0]
		}

		item := &iam.MenuItem{
			Id:             m.Id,
			ParentId:       m.ParentId,
			Name:           m.Name,
			Path:           m.Path,
			Component:      m.Component,
			Icon:           m.Icon,
			MenuType:       int64(m.Type),
			OrderNum:       m.OrderNum,
			Visible:        m.Visible,
			Status:         m.Status,
			PermissionCode: permissionCode,
		}
		nodeMap[m.Id] = item
	}

	for _, item := range nodeMap {
		if item.ParentId == 0 {
			roots = append(roots, item)
			continue
		}
		if parent, ok := nodeMap[item.ParentId]; ok {
			parent.Children = append(parent.Children, item)
		} else {
			roots = append(roots, item)
		}
	}

	sortMenuItems(roots)
	for _, item := range nodeMap {
		sortMenuItems(item.Children)
	}

	resp := &iam.MenuTreeResponse{List: roots}

	go func() {
		if err := cache.SetMenuTree(context.Background(), resp); err != nil {
			l.Errorf("设置菜单树缓存失败: %v", err)
		}
	}()

	return resp, nil
}

func sortMenuItems(items []*iam.MenuItem) {
	if len(items) == 0 {
		return
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].OrderNum != items[j].OrderNum {
			return items[i].OrderNum < items[j].OrderNum
		}
		return items[i].Id < items[j].Id
	})
}
