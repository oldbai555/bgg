package iam

import (
	"context"

	"postapocgame/admin-server/services/iam/internal/repository"
	iamrepo "postapocgame/admin-server/services/iam/internal/repository/iam"
	"postapocgame/admin-server/pkg/errs"
)

// RBACService 承载角色/权限/菜单/接口的授权分配类写操作——
// 这些方法之所以需要领域服务，不是因为跨了很多表（大多数只写一张关联表），
// 而是因为它们是"非平凡业务规则"（RBAC 授权变更），且现有的"先删后插"两步写法
// 本身就需要事务保护（删除成功、插入失败会导致权限被清空且无法自动恢复）。
//
// 注意：这里的"物理删除"不违反项目软删除规则——admin_role_permission/admin_user_role/
// admin_permission_menu/admin_permission_api 是纯关联表，不承载需要审计恢复的业务实体历史，
// 本服务只加事务，不改删除策略本身。
type RBACService struct {
	repo *repository.Repository
}

func NewRBACService(repo *repository.Repository) *RBACService {
	return &RBACService{repo: repo}
}

// UpdateRolePermissions 校验角色/权限存在性 + 事务内先删后插。
func (s *RBACService) UpdateRolePermissions(ctx context.Context, roleID uint64, permissionIDs []uint64) error {
	return s.repo.Transact(ctx, func(ctx context.Context, txRepo *repository.Repository) error {
		roleRepo := iamrepo.NewRoleRepository(txRepo)
		if _, err := roleRepo.FindByID(ctx, roleID); err != nil {
			return errs.Wrap(errs.CodeBadRequest, "角色不存在", err)
		}
		permRepo := iamrepo.NewPermissionRepository(txRepo)
		permissions, err := permRepo.ListByIds(ctx, permissionIDs)
		if err != nil {
			return errs.Wrap(errs.CodeBadRequest, "权限查询有误", err)
		}
		if len(permissions) != len(permissionIDs) {
			return errs.New(errs.CodeBadRequest, "权限不存在")
		}
		rpRepo := iamrepo.NewRolePermissionRepository(txRepo)
		if err := rpRepo.UpdateRolePermissions(ctx, roleID, permissionIDs); err != nil {
			return errs.Wrap(errs.CodeInternalError, "更新角色权限失败", err)
		}
		return nil
	})
}

// UpdateUserRoles 校验用户/角色存在性（禁止分配超级管理员角色）+ 事务内先删后插。
func (s *RBACService) UpdateUserRoles(ctx context.Context, userID uint64, roleIDs []uint64) error {
	return s.repo.Transact(ctx, func(ctx context.Context, txRepo *repository.Repository) error {
		userRepo := iamrepo.NewUserRepository(txRepo)
		if _, err := userRepo.FindByID(ctx, userID); err != nil {
			return errs.Wrap(errs.CodeBadRequest, "用户不存在", err)
		}
		roleRepo := iamrepo.NewRoleRepository(txRepo)
		for _, roleID := range roleIDs {
			if roleID == 1 {
				return errs.New(errs.CodeBadRequest, "不允许分配超级管理员角色")
			}
			if _, err := roleRepo.FindByID(ctx, roleID); err != nil {
				return errs.Wrap(errs.CodeBadRequest, "角色不存在", err)
			}
		}
		urRepo := iamrepo.NewUserRoleRepository(txRepo)
		if err := urRepo.UpdateUserRoles(ctx, userID, roleIDs); err != nil {
			return errs.Wrap(errs.CodeInternalError, "更新用户角色失败", err)
		}
		return nil
	})
}

// UpdatePermissionMenus 校验权限/菜单存在性 + 事务内先删后插。
func (s *RBACService) UpdatePermissionMenus(ctx context.Context, permissionID uint64, menuIDs []uint64) error {
	return s.repo.Transact(ctx, func(ctx context.Context, txRepo *repository.Repository) error {
		permRepo := iamrepo.NewPermissionRepository(txRepo)
		if _, err := permRepo.FindByID(ctx, permissionID); err != nil {
			return errs.Wrap(errs.CodeBadRequest, "权限不存在", err)
		}
		menuRepo := iamrepo.NewMenuRepository(txRepo)
		for _, menuID := range menuIDs {
			if _, err := menuRepo.FindByID(ctx, menuID); err != nil {
				return errs.Wrap(errs.CodeBadRequest, "菜单不存在", err)
			}
		}
		pmRepo := iamrepo.NewPermissionMenuRepository(txRepo)
		if err := pmRepo.UpdatePermissionMenus(ctx, permissionID, menuIDs); err != nil {
			return errs.Wrap(errs.CodeInternalError, "更新权限菜单失败", err)
		}
		return nil
	})
}

// UpdatePermissionApis 校验权限/接口存在性 + 事务内先删后插。
func (s *RBACService) UpdatePermissionApis(ctx context.Context, permissionID uint64, apiIDs []uint64) error {
	return s.repo.Transact(ctx, func(ctx context.Context, txRepo *repository.Repository) error {
		permRepo := iamrepo.NewPermissionRepository(txRepo)
		if _, err := permRepo.FindByID(ctx, permissionID); err != nil {
			return errs.Wrap(errs.CodeBadRequest, "权限不存在", err)
		}
		apiRepo := iamrepo.NewApiRepository(txRepo)
		for _, apiID := range apiIDs {
			if _, err := apiRepo.FindByID(ctx, apiID); err != nil {
				return errs.Wrap(errs.CodeBadRequest, "接口不存在", err)
			}
		}
		paRepo := iamrepo.NewPermissionApiRepository(txRepo)
		if err := paRepo.UpdatePermissionApis(ctx, permissionID, apiIDs); err != nil {
			return errs.Wrap(errs.CodeInternalError, "更新权限接口失败", err)
		}
		return nil
	})
}
