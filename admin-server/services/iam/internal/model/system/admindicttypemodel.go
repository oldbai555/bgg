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
		// WithSession 返回一个绑定到事务 session 的新 AdminDictTypeModel，供 Repository.withSession 调用。
		WithSession(session sqlx.Session) AdminDictTypeModel
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

// WithSession 见接口注释。table 字段直接复用，CachedConn 通过 sqlc.CachedConn.WithSession 换绑。
func (m *customAdminDictTypeModel) WithSession(session sqlx.Session) AdminDictTypeModel {
	return &customAdminDictTypeModel{
		defaultAdminDictTypeModel: &defaultAdminDictTypeModel{
			CachedConn: m.CachedConn.WithSession(session),
			table:      m.table,
		},
	}
}
