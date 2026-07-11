// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package menu

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"

	"postapocgame/admin-server/internal/model/iam"
	iamrepo "postapocgame/admin-server/internal/repository/iam"

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

	// 验证菜单层级规则（创建时menuId为0，不需要检查循环引用）
	if err := l.validateMenuHierarchy(0, req.ParentId, req.MenuType, l.svcCtx.Domain.IAM.Menu); err != nil {
		return err
	}

	m := iam.AdminMenu{
		ParentId:  req.ParentId,
		Name:      req.Name,
		Path:      req.Path,
		Component: req.Component,
		Icon:      req.Icon,
		Type:      req.MenuType,
		OrderNum:  req.OrderNum,
		Visible:   req.Visible,
		Status:    req.Status,
	}

	if err := l.svcCtx.Domain.IAM.Menu.Create(l.ctx, &m); err != nil {
		return errs.Wrap(errs.CodeInternalError, "创建菜单失败", err)
	}

	// 清除菜单树缓存
	cache := l.svcCtx.Repository.BusinessCache
	go func() {
		if err := cache.DeleteMenuTree(context.Background()); err != nil {
			l.Errorf("清除菜单树缓存失败: %v", err)
		}
		// 清除所有用户的菜单树缓存（因为菜单变更会影响所有用户）
		// 注意：go-zero Redis 不支持 SCAN，这里只能清除已知的缓存
		// 实际场景中，可以通过定时任务或延迟清除策略来处理
	}()

	return nil
}

// validateMenuHierarchy 验证菜单层级规则
func (l *MenuCreateLogic) validateMenuHierarchy(menuId, parentId uint64, menuType int64, menuRepo iamrepo.MenuRepository) error {
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

	return nil
}
