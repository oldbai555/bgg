package system

import (
	"context"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ AdminDictTypeModel = (*customAdminDictTypeModel)(nil)

type (
	// AdminDictTypeModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAdminDictTypeModel.
	AdminDictTypeModel interface {
		adminDictTypeModel
		FindIdByCode(ctx context.Context, code string) (uint64, error)
	}

	customAdminDictTypeModel struct {
		*defaultAdminDictTypeModel
	}
)

// NewAdminDictTypeModel returns a model for the database table.
func NewAdminDictTypeModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) AdminDictTypeModel {
	return &customAdminDictTypeModel{
		defaultAdminDictTypeModel: newAdminDictTypeModel(conn, c, opts...),
	}
}

// FindIdByCode 返回字典类型 ID，若不存在返回 ErrNotFound
func (m *customAdminDictTypeModel) FindIdByCode(ctx context.Context, code string) (uint64, error) {
	data, err := m.FindOneByCode(ctx, code)
	if err != nil {
		return 0, err
	}
	return data.Id, nil
}
