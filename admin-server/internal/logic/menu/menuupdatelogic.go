// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package menu

import (
	"context"

	"postapocgame/admin-server/internal/repository"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"

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

	menuRepo := repository.NewMenuRepository(l.svcCtx.Repository)
	m, err := menuRepo.FindByID(l.ctx, req.Id)
	if err != nil {
		return errs.Wrap(errs.CodeInternalError, "查询菜单失败", err)
	}

	// 验证菜单层级规则
	if err := l.validateMenuHierarchy(req.Id, req.ParentId, req.MenuType, menuRepo); err != nil {
		return err
	}

	m.ParentId = req.ParentId
	m.Name = req.Name
	m.Path = req.Path
	m.Component = req.Component
	m.Icon = req.Icon
	m.Type = req.MenuType
	m.OrderNum = req.OrderNum
	if req.Visible != 0 {
		m.Visible = req.Visible
	}
	if req.Status != 0 {
		m.Status = req.Status
	}

	if err := menuRepo.Update(l.ctx, m); err != nil {
		return errs.Wrap(errs.CodeInternalError, "更新菜单失败", err)
	}

	// 清除菜单树缓存
	cache := l.svcCtx.Repository.BusinessCache
	go func() {
		if err := cache.DeleteMenuTree(context.Background()); err != nil {
			l.Errorf("清除菜单树缓存失败: %v", err)
		}
	}()

	return nil
}

// validateMenuHierarchy 验证菜单层级规则
func (l *MenuUpdateLogic) validateMenuHierarchy(menuId, parentId uint64, menuType int64, menuRepo repository.MenuRepository) error {
	// 规则1：目录（type=1）只能存在根节点下（parent_id = 0）
	if menuType == 1 && parentId != 0 {
		return errs.New(errs.CodeBadRequest, "目录只能存在根节点下")
	}

	// 规则2：菜单（type=2）可以在目录下，也可以在根节点下
	// 规则3：按钮（type=3）只能在菜单下（parent_id必须是菜单的id，且该菜单的type=2）
	if menuType == 3 {
		if parentId == 0 {
			return errs.New(errs.CodeBadRequest, "按钮只能存在于菜单下")
		}
		// 验证父节点必须是菜单（type=2）
		parentMenu, err := menuRepo.FindByID(l.ctx, parentId)
		if err != nil {
			return errs.Wrap(errs.CodeBadRequest, "父菜单不存在", err)
		}
		if parentMenu.Type != 2 {
			return errs.New(errs.CodeBadRequest, "按钮的父节点必须是菜单（type=2）")
		}
	}

	// 规则4：菜单（type=2）如果parent_id不为0，则父节点必须是目录（type=1）
	if menuType == 2 && parentId != 0 {
		parentMenu, err := menuRepo.FindByID(l.ctx, parentId)
		if err != nil {
			return errs.Wrap(errs.CodeBadRequest, "父菜单不存在", err)
		}
		if parentMenu.Type != 1 {
			return errs.New(errs.CodeBadRequest, "菜单的父节点必须是目录（type=1）或根节点")
		}
	}

	// 检查循环引用：不能将菜单移动到自己的子节点下
	if parentId != 0 {
		if err := l.checkCircularReference(menuId, parentId, menuRepo); err != nil {
			return err
		}
	}

	return nil
}

// checkCircularReference 检查循环引用
func (l *MenuUpdateLogic) checkCircularReference(menuId, parentId uint64, menuRepo repository.MenuRepository) error {
	// 如果parentId就是menuId本身，直接返回错误
	if parentId == menuId {
		return errs.New(errs.CodeBadRequest, "不能将菜单设置为自己的父节点")
	}

	// 递归检查parentId是否是menuId的子节点
	currentId := parentId
	for currentId != 0 {
		parentMenu, err := menuRepo.FindByID(l.ctx, currentId)
		if err != nil {
			// 如果查询失败，可能是数据不一致，但不影响循环检查
			break
		}
		if parentMenu.ParentId == menuId {
			return errs.New(errs.CodeBadRequest, "不能将菜单移动到自己的子节点下")
		}
		currentId = parentMenu.ParentId
	}

	return nil
}
