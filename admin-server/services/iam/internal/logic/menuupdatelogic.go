package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iam"
	iamrepo "postapocgame/admin-server/services/iam/internal/repository/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type MenuUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMenuUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MenuUpdateLogic {
	return &MenuUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *MenuUpdateLogic) MenuUpdate(in *iam.MenuUpdateRequest) (*iam.Empty, error) {
	if in == nil || in.Id == 0 {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "菜单ID不能为空"))
	}

	m, err := l.svcCtx.Domain.IAM.Menu.FindByID(l.ctx, in.Id)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "查询菜单失败", err))
	}

	if err := l.validateMenuHierarchy(in.Id, in.ParentId, in.MenuType, l.svcCtx.Domain.IAM.Menu); err != nil {
		return nil, toGRPCStatus(err)
	}

	m.ParentId = in.ParentId
	m.Name = in.Name
	m.Path = in.Path
	m.Component = in.Component
	m.Icon = in.Icon
	m.Type = in.MenuType
	m.OrderNum = in.OrderNum
	if in.Visible != 0 {
		m.Visible = in.Visible
	}
	if in.Status != 0 {
		m.Status = in.Status
	}

	if err := l.svcCtx.Domain.IAM.Menu.Update(l.ctx, m); err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "更新菜单失败", err))
	}

	cache := l.svcCtx.Repository.BusinessCache
	go func() {
		if err := cache.DeleteMenuTree(context.Background()); err != nil {
			l.Errorf("清除菜单树缓存失败: %v", err)
		}
	}()

	return &iam.Empty{}, nil
}

func (l *MenuUpdateLogic) validateMenuHierarchy(menuId, parentId uint64, menuType int64, menuRepo iamrepo.MenuRepository) error {
	if menuType == 1 && parentId != 0 {
		return errs.New(errs.CodeBadRequest, "目录只能存在根节点下")
	}

	if menuType == 3 {
		if parentId == 0 {
			return errs.New(errs.CodeBadRequest, "按钮只能存在于菜单下")
		}
		parentMenu, err := menuRepo.FindByID(l.ctx, parentId)
		if err != nil {
			return errs.Wrap(errs.CodeBadRequest, "父菜单不存在", err)
		}
		if parentMenu.Type != 2 {
			return errs.New(errs.CodeBadRequest, "按钮的父节点必须是菜单（type=2）")
		}
	}

	if menuType == 2 && parentId != 0 {
		parentMenu, err := menuRepo.FindByID(l.ctx, parentId)
		if err != nil {
			return errs.Wrap(errs.CodeBadRequest, "父菜单不存在", err)
		}
		if parentMenu.Type != 1 {
			return errs.New(errs.CodeBadRequest, "菜单的父节点必须是目录（type=1）或根节点")
		}
	}

	if parentId != 0 {
		if err := l.checkCircularReference(menuId, parentId, menuRepo); err != nil {
			return err
		}
	}

	return nil
}

func (l *MenuUpdateLogic) checkCircularReference(menuId, parentId uint64, menuRepo iamrepo.MenuRepository) error {
	if parentId == menuId {
		return errs.New(errs.CodeBadRequest, "不能将菜单设置为自己的父节点")
	}

	currentId := parentId
	for currentId != 0 {
		parentMenu, err := menuRepo.FindByID(l.ctx, currentId)
		if err != nil {
			break
		}
		if parentMenu.ParentId == menuId {
			return errs.New(errs.CodeBadRequest, "不能将菜单移动到自己的子节点下")
		}
		currentId = parentMenu.ParentId
	}

	return nil
}
