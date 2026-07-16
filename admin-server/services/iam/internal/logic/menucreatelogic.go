package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iam"
	iammodel "postapocgame/admin-server/services/iam/internal/model/iam"
	iamrepo "postapocgame/admin-server/services/iam/internal/repository/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type MenuCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMenuCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MenuCreateLogic {
	return &MenuCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *MenuCreateLogic) MenuCreate(in *iam.MenuCreateRequest) (*iam.Empty, error) {
	if in == nil || in.Name == "" || in.MenuType == 0 {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "菜单名称和类型不能为空"))
	}

	if err := l.validateMenuHierarchy(0, in.ParentId, in.MenuType, l.svcCtx.Domain.IAM.Menu); err != nil {
		return nil, toGRPCStatus(err)
	}

	m := iammodel.AdminMenu{
		ParentId:  in.ParentId,
		Name:      in.Name,
		Path:      in.Path,
		Component: in.Component,
		Icon:      in.Icon,
		Type:      in.MenuType,
		OrderNum:  in.OrderNum,
		Visible:   in.Visible,
		Status:    in.Status,
	}

	if err := l.svcCtx.Domain.IAM.Menu.Create(l.ctx, &m); err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "创建菜单失败", err))
	}

	cache := l.svcCtx.Repository.BusinessCache
	go func() {
		if err := cache.DeleteMenuTree(context.Background()); err != nil {
			l.Errorf("清除菜单树缓存失败: %v", err)
		}
	}()

	return &iam.Empty{}, nil
}

// validateMenuHierarchy 验证菜单层级规则
func (l *MenuCreateLogic) validateMenuHierarchy(menuId, parentId uint64, menuType int64, menuRepo iamrepo.MenuRepository) error {
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

	return nil
}
