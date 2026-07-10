package monitoring

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ AdminLoginLogModel = (*customAdminLoginLogModel)(nil)

type (
	// AdminLoginLogModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAdminLoginLogModel.
	AdminLoginLogModel interface {
		adminLoginLogModel
		// WithSession 返回一个绑定到事务 session 的新 AdminLoginLogModel，供 Repository.withSession 调用。
		WithSession(session sqlx.Session) AdminLoginLogModel
	}

	customAdminLoginLogModel struct {
		*defaultAdminLoginLogModel
	}
)

// NewAdminLoginLogModel returns a model for the database table.
func NewAdminLoginLogModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) AdminLoginLogModel {
	return &customAdminLoginLogModel{
		defaultAdminLoginLogModel: newAdminLoginLogModel(conn, c, opts...),
	}
}

// WithSession 见接口注释。table 字段直接复用，CachedConn 通过 sqlc.CachedConn.WithSession 换绑。
func (m *customAdminLoginLogModel) WithSession(session sqlx.Session) AdminLoginLogModel {
	return &customAdminLoginLogModel{
		defaultAdminLoginLogModel: &defaultAdminLoginLogModel{
			CachedConn: m.CachedConn.WithSession(session),
			table:      m.table,
		},
	}
}
