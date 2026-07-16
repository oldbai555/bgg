package misc

import (
	"context"
	"postapocgame/admin-server/services/iam/internal/repository"

	miscmodel "postapocgame/admin-server/services/iam/internal/model/misc"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type DemoRepository interface {
	FindByID(ctx context.Context, id uint64) (*miscmodel.Demo, error)
	FindPage(ctx context.Context, page, pageSize int64, name string) ([]miscmodel.Demo, int64, error)
	DeleteByID(ctx context.Context, id uint64) error
	Create(ctx context.Context, demo *miscmodel.Demo) error
	Update(ctx context.Context, demo *miscmodel.Demo) error
}

type demoRepository struct {
	model miscmodel.DemoModel
	conn  sqlx.SqlConn
}

func NewDemoRepository(repo *repository.Repository) DemoRepository {
	return &demoRepository{model: repo.DemoModel, conn: repo.DB}
}

func (r *demoRepository) FindByID(ctx context.Context, id uint64) (*miscmodel.Demo, error) {
	return r.model.FindOne(ctx, id)
}

func (r *demoRepository) FindPage(ctx context.Context, page, pageSize int64, name string) ([]miscmodel.Demo, int64, error) {
	// 目前生成方法不支持复杂过滤，简单复用生成的分页
	return r.model.FindPage(ctx, page, pageSize)
}

func (r *demoRepository) DeleteByID(ctx context.Context, id uint64) error {
	return r.model.Delete(ctx, id)
}

func (r *demoRepository) Create(ctx context.Context, demo *miscmodel.Demo) error {
	result, err := r.model.Insert(ctx, demo)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	demo.Id = uint64(id)
	return nil
}

func (r *demoRepository) Update(ctx context.Context, demo *miscmodel.Demo) error {
	return r.model.Update(ctx, demo)
}
