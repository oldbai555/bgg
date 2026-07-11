package system

import (
	"postapocgame/admin-server/internal/repository"
	"context"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	systemmodel "postapocgame/admin-server/internal/model/system"
)

type DictTypeRepository interface {
	FindByID(ctx context.Context, id uint64) (*systemmodel.AdminDictType, error)
	FindByCode(ctx context.Context, code string) (*systemmodel.AdminDictType, error)
	FindPage(ctx context.Context, page, pageSize int64, name, code string) ([]systemmodel.AdminDictType, int64, error)
	DeleteByID(ctx context.Context, id uint64) error
	Create(ctx context.Context, dictType *systemmodel.AdminDictType) error
	Update(ctx context.Context, dictType *systemmodel.AdminDictType) error
}

type dictTypeRepository struct {
	model systemmodel.AdminDictTypeModel
	conn  sqlx.SqlConn
}

func NewDictTypeRepository(repo *repository.Repository) DictTypeRepository {
	return &dictTypeRepository{model: repo.AdminDictTypeModel, conn: repo.DB}
}

func (r *dictTypeRepository) FindByID(ctx context.Context, id uint64) (*systemmodel.AdminDictType, error) {
	return r.model.FindOne(ctx, id)
}

func (r *dictTypeRepository) FindByCode(ctx context.Context, code string) (*systemmodel.AdminDictType, error) {
	return r.model.FindOneByCode(ctx, code)
}

func (r *dictTypeRepository) FindPage(ctx context.Context, page, pageSize int64, name, code string) ([]systemmodel.AdminDictType, int64, error) {
	// 目前生成方法不支持复杂过滤，简单复用生成的分页
	return r.model.FindPage(ctx, page, pageSize)
}

func (r *dictTypeRepository) DeleteByID(ctx context.Context, id uint64) error {
	return r.model.Delete(ctx, id)
}

func (r *dictTypeRepository) Create(ctx context.Context, dictType *systemmodel.AdminDictType) error {
	result, err := r.model.Insert(ctx, dictType)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	dictType.Id = uint64(id)
	return nil
}

func (r *dictTypeRepository) Update(ctx context.Context, dictType *systemmodel.AdminDictType) error {
	return r.model.Update(ctx, dictType)
}
