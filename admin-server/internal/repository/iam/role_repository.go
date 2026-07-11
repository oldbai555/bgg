package iam

import (
	"postapocgame/admin-server/internal/repository"
	"context"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	iammodel "postapocgame/admin-server/internal/model/iam"
)

type RoleRepository interface {
	FindByID(ctx context.Context, id uint64) (*iammodel.AdminRole, error)
	FindByCode(ctx context.Context, code string) (*iammodel.AdminRole, error)
	FindPage(ctx context.Context, page, pageSize int64, name string) ([]iammodel.AdminRole, int64, error)
	FindChunk(ctx context.Context, limit int64, lastId uint64) ([]iammodel.AdminRole, uint64, error)
	DeleteByID(ctx context.Context, id uint64) error
	Create(ctx context.Context, role *iammodel.AdminRole) error
	Update(ctx context.Context, role *iammodel.AdminRole) error
}

type roleRepository struct {
	model iammodel.AdminRoleModel
	conn  sqlx.SqlConn
}

func NewRoleRepository(repo *repository.Repository) RoleRepository {
	return &roleRepository{model: repo.AdminRoleModel, conn: repo.DB}
}

func (r *roleRepository) FindByID(ctx context.Context, id uint64) (*iammodel.AdminRole, error) {
	return r.model.FindOne(ctx, id)
}

func (r *roleRepository) FindByCode(ctx context.Context, code string) (*iammodel.AdminRole, error) {
	return r.model.FindOneByCode(ctx, code)
}

// FindPage 分页查询角色列表（符合新规范）
func (r *roleRepository) FindPage(ctx context.Context, page, pageSize int64, name string) ([]iammodel.AdminRole, int64, error) {
	// 目前生成方法不支持模糊过滤，简单复用生成的分页
	return r.model.FindPage(ctx, page, pageSize)
}

// FindChunk 分片查询角色列表（基于lastId，适用于大数据量分批处理）
func (r *roleRepository) FindChunk(ctx context.Context, limit int64, lastId uint64) ([]iammodel.AdminRole, uint64, error) {
	return r.model.FindChunk(ctx, limit, lastId)
}

func (r *roleRepository) DeleteByID(ctx context.Context, id uint64) error {
	return r.model.Delete(ctx, id)
}

func (r *roleRepository) Create(ctx context.Context, role *iammodel.AdminRole) error {
	result, err := r.model.Insert(ctx, role)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	role.Id = uint64(id)
	return nil
}

func (r *roleRepository) Update(ctx context.Context, role *iammodel.AdminRole) error {
	return r.model.Update(ctx, role)
}
