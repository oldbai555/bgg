package system

import (
	"postapocgame/admin-server/internal/repository"
	"context"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	systemmodel "postapocgame/admin-server/internal/model/system"
)

type ConfigRepository interface {
	FindByID(ctx context.Context, id uint64) (*systemmodel.AdminConfig, error)
	FindByKey(ctx context.Context, key string) (*systemmodel.AdminConfig, error)
	FindPage(ctx context.Context, page, pageSize int64, group, key string) ([]systemmodel.AdminConfig, int64, error)
	DeleteByID(ctx context.Context, id uint64) error
	Create(ctx context.Context, config *systemmodel.AdminConfig) error
	Update(ctx context.Context, config *systemmodel.AdminConfig) error
}

type configRepository struct {
	model systemmodel.AdminConfigModel
	conn  sqlx.SqlConn
}

func NewConfigRepository(repo *repository.Repository) ConfigRepository {
	return &configRepository{model: repo.AdminConfigModel, conn: repo.DB}
}

func (r *configRepository) FindByID(ctx context.Context, id uint64) (*systemmodel.AdminConfig, error) {
	return r.model.FindOne(ctx, id)
}

func (r *configRepository) FindByKey(ctx context.Context, key string) (*systemmodel.AdminConfig, error) {
	return r.model.FindOneByKey(ctx, key)
}

func (r *configRepository) FindPage(ctx context.Context, page, pageSize int64, group, key string) ([]systemmodel.AdminConfig, int64, error) {
	// 目前生成方法不支持复杂过滤，简单复用生成的分页
	return r.model.FindPage(ctx, page, pageSize)
}

func (r *configRepository) DeleteByID(ctx context.Context, id uint64) error {
	return r.model.Delete(ctx, id)
}

func (r *configRepository) Create(ctx context.Context, config *systemmodel.AdminConfig) error {
	result, err := r.model.Insert(ctx, config)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	config.Id = uint64(id)
	return nil
}

func (r *configRepository) Update(ctx context.Context, config *systemmodel.AdminConfig) error {
	return r.model.Update(ctx, config)
}
