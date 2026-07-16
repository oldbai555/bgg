package system

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ AdminConfigModel = (*customAdminConfigModel)(nil)

type (
	// AdminConfigModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAdminConfigModel.
	AdminConfigModel interface {
		adminConfigModel
		// WithSession 返回一个绑定到事务 session 的新 AdminConfigModel，供 Repository.withSession 调用。
		WithSession(session sqlx.Session) AdminConfigModel
	}

	customAdminConfigModel struct {
		*defaultAdminConfigModel
	}
)

// NewAdminConfigModel returns a model for the database table.
func NewAdminConfigModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) AdminConfigModel {
	return &customAdminConfigModel{
		defaultAdminConfigModel: newAdminConfigModel(conn, c, opts...),
	}
}

// WithSession 见接口注释。table 字段直接复用，CachedConn 通过 sqlc.CachedConn.WithSession 换绑。
func (m *customAdminConfigModel) WithSession(session sqlx.Session) AdminConfigModel {
	return &customAdminConfigModel{
		defaultAdminConfigModel: &defaultAdminConfigModel{
			CachedConn: m.CachedConn.WithSession(session),
			table:      m.table,
		},
	}
}
