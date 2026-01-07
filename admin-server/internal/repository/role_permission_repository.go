package repository

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"postapocgame/admin-server/internal/model"
	"postapocgame/admin-server/pkg/errs"
)

type RolePermissionRepository interface {
	ListPermissionIDsByRoleID(ctx context.Context, roleID uint64) ([]uint64, error)
	UpdateRolePermissions(ctx context.Context, roleID uint64, permissionIDs []uint64) error
}

type rolePermissionRepository struct {
	model model.AdminRolePermissionModel
	conn  sqlx.SqlConn
}

func NewRolePermissionRepository(repo *Repository) RolePermissionRepository {
	return &rolePermissionRepository{
		model: repo.AdminRolePermissionModel,
		conn:  repo.DB,
	}
}

// ListPermissionIDsByRoleID 查询角色拥有的权限ID列表
func (r *rolePermissionRepository) ListPermissionIDsByRoleID(ctx context.Context, roleID uint64) ([]uint64, error) {
	var list []model.AdminRolePermission
	sql, args, err := sq.Select("*").From("admin_role_permission").Where(sq.Eq{"role_id": roleID}).ToSql()
	if err != nil {
		return nil, errs.Wrap(errs.CodeBadDB, "sql生成有误", err)
	}
	if err := r.conn.QueryRowsCtx(ctx, &list, sql, args...); err != nil {
		return nil, errs.Wrap(errs.CodeBadDB, "sql执行有误", err)
	}
	ids := make([]uint64, 0, len(list))
	for _, rp := range list {
		ids = append(ids, rp.PermissionId)
	}
	return ids, nil
}

// UpdateRolePermissions 更新角色的权限关联（先物理删除旧的，再添加新的）
func (r *rolePermissionRepository) UpdateRolePermissions(ctx context.Context, roleID uint64, permissionIDs []uint64) error {
	// 先物理删除该角色的所有权限关联
	sql, args, err := sq.Delete("admin_role_permission").Where(sq.Eq{"role_id": roleID}).ToSql()
	if err != nil {
		return errs.Wrap(errs.CodeBadDB, "sql生成有误", err)
	}
	_, err = r.conn.ExecCtx(ctx, sql, args...)
	if err != nil {
		return err
	}
	if len(permissionIDs) == 0 {
		return nil
	}
	// 如果有新的权限，添加关联
	db := sq.Insert("admin_role_permission").Columns("role_id", "permission_id")
	for _, permID := range permissionIDs {
		db = db.Values(roleID, permID)
	}
	sql, args, err = db.ToSql()
	if err != nil {
		return errs.Wrap(errs.CodeBadDB, "sql生成有误", err)
	}
	_, err = r.conn.ExecCtx(ctx, sql, args...)
	if err != nil {
		return errs.Wrap(errs.CodeBadDB, "sql执行有误", err)
	}
	return nil
}
