package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type MenuMyTreeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMenuMyTreeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MenuMyTreeLogic {
	return &MenuMyTreeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *MenuMyTreeLogic) MenuMyTree(in *iam.MenuMyTreeRequest) (*iam.MenuTreeResponse, error) {
	if in == nil || in.UserId == 0 {
		return nil, toGRPCStatus(errs.New(errs.CodeUnauthorized, "未登录或登录已过期"))
	}
	userID := in.UserId

	// 超级管理员（user_id=1）默认拥有最高权限，直接返回完整菜单树
	if userID == 1 {
		treeLogic := NewMenuTreeLogic(l.ctx, l.svcCtx)
		return treeLogic.MenuTree(&iam.Empty{})
	}

	cache := l.svcCtx.Repository.BusinessCache
	var cachedResp iam.MenuTreeResponse
	if err := cache.GetUserMenuTree(l.ctx, userID, &cachedResp); err == nil {
		return &cachedResp, nil
	}

	roleIDs, err := l.svcCtx.Domain.IAM.UserRole.ListRoleIDsByUserID(l.ctx, userID)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "查询用户角色失败", err))
	}

	perms, err := l.svcCtx.Domain.IAM.Permission.ListByRoleIDs(l.ctx, roleIDs)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "查询用户权限失败", err))
	}

	permCodes := make([]string, 0, len(perms))
	permSet := make(map[string]struct{}, len(perms))
	for _, p := range perms {
		permCodes = append(permCodes, p.Code)
		permSet[p.Code] = struct{}{}
	}

	hasSuperAdmin := false
	for _, code := range permCodes {
		if code == "*" {
			hasSuperAdmin = true
			break
		}
	}

	if hasSuperAdmin {
		treeLogic := NewMenuTreeLogic(l.ctx, l.svcCtx)
		return treeLogic.MenuTree(&iam.Empty{})
	}

	menuPermissionMap, err := l.svcCtx.Domain.IAM.PermissionMenu.ListMenuPermissionCodes(l.ctx)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "查询菜单权限关联失败", err))
	}

	treeLogic := NewMenuTreeLogic(l.ctx, l.svcCtx)
	fullTree, err := treeLogic.MenuTree(&iam.Empty{})
	if err != nil {
		return nil, toGRPCStatus(err)
	}

	var filterMenu func(item *iam.MenuItem) *iam.MenuItem
	filterMenu = func(item *iam.MenuItem) *iam.MenuItem {
		if codes, hasPerms := menuPermissionMap[item.Id]; hasPerms && len(codes) > 0 {
			hasAccess := false
			for _, code := range codes {
				if _, ok := permSet[code]; ok {
					hasAccess = true
					break
				}
			}
			if !hasAccess {
				return nil
			}
		}

		filtered := &iam.MenuItem{
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
		}
		if codes, hasPerms := menuPermissionMap[item.Id]; hasPerms && len(codes) > 0 {
			filtered.PermissionCode = codes[0]
		}
		for _, child := range item.Children {
			if filteredChild := filterMenu(child); filteredChild != nil {
				filtered.Children = append(filtered.Children, filteredChild)
			}
		}

		// 目录节点本身不是可访问页面，只是子菜单的容器；子菜单全部被过滤掉后
		// 目录不能再保留在树里——否则会渲染成一个指向空路由的死链接（点击 404）。
		if item.MenuType == 1 && len(filtered.Children) == 0 {
			return nil
		}

		if item.MenuType == 3 {
			if codes, hasPerms := menuPermissionMap[item.Id]; hasPerms && len(codes) > 0 {
				hasAccess := false
				for _, code := range codes {
					if _, ok := permSet[code]; ok {
						hasAccess = true
						break
					}
				}
				if !hasAccess {
					return nil
				}
			}
		}
		return filtered
	}

	filteredRoots := make([]*iam.MenuItem, 0)
	for _, root := range fullTree.List {
		if filtered := filterMenu(root); filtered != nil {
			filteredRoots = append(filteredRoots, filtered)
		}
	}

	resp := &iam.MenuTreeResponse{List: filteredRoots}

	go func() {
		if err := cache.SetUserMenuTree(context.Background(), userID, resp); err != nil {
			l.Errorf("设置用户菜单树缓存失败: userId=%d, error=%v", userID, err)
		}
	}()

	return resp, nil
}
