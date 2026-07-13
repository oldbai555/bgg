package iam

import (
	"context"
	"postapocgame/admin-server/services/iam/internal/repository"

	iammodel "postapocgame/admin-server/services/iam/internal/model/iam"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type MenuRepository interface {
	ListAll(ctx context.Context) ([]iammodel.AdminMenu, error)
	FindByID(ctx context.Context, id uint64) (*iammodel.AdminMenu, error)
	Create(ctx context.Context, m *iammodel.AdminMenu) error
	Update(ctx context.Context, m *iammodel.AdminMenu) error
	DeleteByID(ctx context.Context, id uint64) error
}

type menuRepository struct {
	model iammodel.AdminMenuModel
	conn  sqlx.SqlConn
}

func NewMenuRepository(repo *repository.Repository) MenuRepository {
	return &menuRepository{model: repo.AdminMenuModel, conn: repo.DB}
}

func (r *menuRepository) ListAll(ctx context.Context) ([]iammodel.AdminMenu, error) {
	// 直接查询所有未删除的菜单，按 order_num 和 id 排序
	var list []iammodel.AdminMenu
	query := "select id, parent_id, name, path, component, icon, type, order_num, visible, status, created_at, updated_at, deleted_at from admin_menu where deleted_at = 0 order by order_num asc, id asc"
	err := r.conn.QueryRowsCtx(ctx, &list, query)
	return list, err
}

func (r *menuRepository) FindByID(ctx context.Context, id uint64) (*iammodel.AdminMenu, error) {
	return r.model.FindOne(ctx, id)
}

func (r *menuRepository) Create(ctx context.Context, m *iammodel.AdminMenu) error {
	result, err := r.model.Insert(ctx, m)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	m.Id = uint64(id)
	return nil
}

func (r *menuRepository) Update(ctx context.Context, m *iammodel.AdminMenu) error {
	return r.model.Update(ctx, m)
}

func (r *menuRepository) DeleteByID(ctx context.Context, id uint64) error {
	return r.model.Delete(ctx, id)
}
