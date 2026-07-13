package system

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ AdminFileModel = (*customAdminFileModel)(nil)

type (
	// AdminFileModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAdminFileModel.
	AdminFileModel interface {
		adminFileModel
		// WithSession 返回一个绑定到事务 session 的新 AdminFileModel，供 Repository.withSession 调用。
		WithSession(session sqlx.Session) AdminFileModel
	}

	customAdminFileModel struct {
		*defaultAdminFileModel
	}
)

// NewAdminFileModel returns a model for the database table.
func NewAdminFileModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) AdminFileModel {
	return &customAdminFileModel{
		defaultAdminFileModel: newAdminFileModel(conn, c, opts...),
	}
}

// WithSession 见接口注释。table 字段直接复用，CachedConn 通过 sqlc.CachedConn.WithSession 换绑。
func (m *customAdminFileModel) WithSession(session sqlx.Session) AdminFileModel {
	return &customAdminFileModel{
		defaultAdminFileModel: &defaultAdminFileModel{
			CachedConn: m.CachedConn.WithSession(session),
			table:      m.table,
		},
	}
}
