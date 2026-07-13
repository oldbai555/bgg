package iam

import (
	"context"
	"postapocgame/admin-server/services/iam/internal/repository"

	iammodel "postapocgame/admin-server/services/iam/internal/model/iam"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ApiRepository interface {
	FindByID(ctx context.Context, id uint64) (*iammodel.AdminApi, error)
	FindByMethodAndPath(ctx context.Context, method, path string) (*iammodel.AdminApi, error)
	FindPage(ctx context.Context, page, pageSize int64, name string) ([]iammodel.AdminApi, int64, error)
	Create(ctx context.Context, api *iammodel.AdminApi) error
	Update(ctx context.Context, api *iammodel.AdminApi) error
	DeleteByID(ctx context.Context, id uint64) error
}

type apiRepository struct {
	model iammodel.AdminApiModel
	conn  sqlx.SqlConn
}

func NewApiRepository(repo *repository.Repository) ApiRepository {
	return &apiRepository{model: repo.AdminApiModel, conn: repo.DB}
}

func (r *apiRepository) FindByID(ctx context.Context, id uint64) (*iammodel.AdminApi, error) {
	return r.model.FindOne(ctx, id)
}

func (r *apiRepository) FindByMethodAndPath(ctx context.Context, method, path string) (*iammodel.AdminApi, error) {
	return r.model.FindOneByMethodPath(ctx, method, path)
}

func (r *apiRepository) FindPage(ctx context.Context, page, pageSize int64, name string) ([]iammodel.AdminApi, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	offset := (page - 1) * pageSize

	var (
		list  []iammodel.AdminApi
		total int64
	)

	if name == "" {
		return r.model.FindPage(ctx, page, pageSize)
	}

	// 带名称模糊筛选的自定义查询
	countQuery := "select count(*) from admin_api where deleted_at = 0 and name like ?"
	if err := r.conn.QueryRowCtx(ctx, &total, countQuery, "%"+name+"%"); err != nil {
		return nil, 0, err
	}
	query := "select * from admin_api where deleted_at = 0 and name like ? order by id desc limit ? offset ?"
	if err := r.conn.QueryRowsCtx(ctx, &list, query, "%"+name+"%", pageSize, offset); err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *apiRepository) Create(ctx context.Context, api *iammodel.AdminApi) error {
	result, err := r.model.Insert(ctx, api)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	api.Id = uint64(id)
	return nil
}

func (r *apiRepository) Update(ctx context.Context, api *iammodel.AdminApi) error {
	return r.model.Update(ctx, api)
}

func (r *apiRepository) DeleteByID(ctx context.Context, id uint64) error {
	return r.model.Delete(ctx, id)
}
