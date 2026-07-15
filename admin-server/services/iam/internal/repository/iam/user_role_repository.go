package iam

import (
	"context"
	"postapocgame/admin-server/services/iam/internal/repository"

	iammodel "postapocgame/admin-server/services/iam/internal/model/iam"
	"postapocgame/admin-server/pkg/errs"

	sq "github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type UserRoleRepository interface {
	ListRoleIDsByUserID(ctx context.Context, userID uint64) ([]uint64, error)
	UpdateUserRoles(ctx context.Context, userID uint64, roleIDs []uint64) error
	// ListRoleNamesByUserIDs 批量查询多个用户的角色名称，一次 JOIN 查询代替逐用户查询，
	// 供 UserList 组装 UserItem.RoleNames 使用（避免 N+1）
	ListRoleNamesByUserIDs(ctx context.Context, userIDs []uint64) (map[uint64][]string, error)
}

type userRoleRepository struct {
	model iammodel.AdminUserRoleModel
	conn  sqlx.SqlConn
}

func NewUserRoleRepository(repo *repository.Repository) UserRoleRepository {
	return &userRoleRepository{model: repo.AdminUserRoleModel, conn: repo.DB}
}

func (r *userRoleRepository) ListRoleIDsByUserID(ctx context.Context, userID uint64) ([]uint64, error) {
	var list []iammodel.AdminUserRole
	query := "select * from admin_user_role where user_id = ?"
	if err := r.conn.QueryRowsCtx(ctx, &list, query, userID); err != nil {
		return nil, err
	}
	ids := make([]uint64, 0, len(list))
	for _, ur := range list {
		ids = append(ids, ur.RoleId)
	}
	return ids, nil
}

// UpdateUserRoles 更新用户的角色关联（先物理删除旧的，再添加新的）
func (r *userRoleRepository) UpdateUserRoles(ctx context.Context, userID uint64, roleIDs []uint64) error {
	// 先物理删除该用户的所有角色关联
	_, err := r.conn.ExecCtx(ctx, "delete from admin_user_role where user_id = ?", userID)
	if err != nil {
		return err
	}

	// 如果有新的角色，添加关联
	if len(roleIDs) > 0 {
		for _, roleID := range roleIDs {
			newUR := &iammodel.AdminUserRole{
				UserId: userID,
				RoleId: roleID,
			}
			_, err := r.model.Insert(ctx, newUR)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *userRoleRepository) ListRoleNamesByUserIDs(ctx context.Context, userIDs []uint64) (map[uint64][]string, error) {
	result := make(map[uint64][]string, len(userIDs))
	if len(userIDs) == 0 {
		return result, nil
	}

	sql, args, err := sq.Select("ur.user_id AS user_id", "r.name AS role_name").
		From("admin_user_role ur").
		Join("admin_role r ON r.id = ur.role_id").
		Where(sq.And{
			sq.Eq{"ur.user_id": userIDs},
			sq.Eq{"r.deleted_at": 0},
		}).
		ToSql()
	if err != nil {
		return nil, errs.Wrap(errs.CodeBadDB, "sql生成有误", err)
	}

	var rows []struct {
		UserId   uint64 `db:"user_id"`
		RoleName string `db:"role_name"`
	}
	if err := r.conn.QueryRowsCtx(ctx, &rows, sql, args...); err != nil {
		return nil, err
	}
	for _, row := range rows {
		result[row.UserId] = append(result[row.UserId], row.RoleName)
	}
	return result, nil
}
