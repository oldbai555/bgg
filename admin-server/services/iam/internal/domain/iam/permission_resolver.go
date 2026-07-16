package iam

import (
	"context"

	"postapocgame/admin-server/services/iam/internal/repository"
	iamrepo "postapocgame/admin-server/services/iam/internal/repository/iam"
	"postapocgame/admin-server/pkg/errs"
)

// PermissionResolver RBAC 权限解析领域服务。
type PermissionResolver struct {
	repo *repository.Repository
}

func NewPermissionResolver(repo *repository.Repository) *PermissionResolver {
	return &PermissionResolver{repo: repo}
}

// CanAccess 判断用户是否有权访问指定 HTTP 接口。
func (r *PermissionResolver) CanAccess(ctx context.Context, userID uint64, method, path string) (bool, error) {
	if userID == 1 {
		return true, nil
	}

	apiRepo := iamrepo.NewApiRepository(r.repo)
	api, err := apiRepo.FindByMethodAndPath(ctx, method, path)
	if err != nil {
		return false, errs.New(errs.CodeNotFound, "接口不存在")
	}

	if api.Status != 1 {
		return true, nil
	}

	userRoleRepo := iamrepo.NewUserRoleRepository(r.repo)
	roleIds, err := userRoleRepo.ListRoleIDsByUserID(ctx, userID)
	if err != nil {
		return false, errs.Wrap(errs.CodeInternalError, "获取用户角色失败", err)
	}
	if len(roleIds) == 0 {
		return false, nil
	}

	rolePermissionRepo := iamrepo.NewRolePermissionRepository(r.repo)
	permissionIds := make(map[uint64]bool)
	for _, roleId := range roleIds {
		permIds, err := rolePermissionRepo.ListPermissionIDsByRoleID(ctx, roleId)
		if err != nil {
			continue
		}
		for _, permId := range permIds {
			permissionIds[permId] = true
		}
	}

	if permissionIds[1] {
		return true, nil
	}

	permissionApiRepo := iamrepo.NewPermissionApiRepository(r.repo)
	apiPermissionIds, err := permissionApiRepo.ListPermissionIDsByApiID(ctx, api.Id)
	if err != nil {
		return true, nil
	}

	for _, apiPermId := range apiPermissionIds {
		if permissionIds[apiPermId] {
			return true, nil
		}
	}

	return false, nil
}
