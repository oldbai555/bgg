package task

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ AdminTaskModel = (*customAdminTaskModel)(nil)

type (
	// AdminTaskModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAdminTaskModel.
	AdminTaskModel interface {
		adminTaskModel
		// WithSession 返回一个绑定到事务 session 的新 AdminTaskModel，供 Repository.withSession 调用。
		WithSession(session sqlx.Session) AdminTaskModel
	}

	customAdminTaskModel struct {
		*defaultAdminTaskModel
	}
)

// NewAdminTaskModel returns a model for the database table.
func NewAdminTaskModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) AdminTaskModel {
	return &customAdminTaskModel{
		defaultAdminTaskModel: newAdminTaskModel(conn, c, opts...),
	}
}

// WithSession 见接口注释。table 字段直接复用，CachedConn 通过 sqlc.CachedConn.WithSession 换绑。
func (m *customAdminTaskModel) WithSession(session sqlx.Session) AdminTaskModel {
	return &customAdminTaskModel{
		defaultAdminTaskModel: &defaultAdminTaskModel{
			CachedConn: m.CachedConn.WithSession(session),
			table:      m.table,
		},
	}
}
