package monitoring

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ AdminPerformanceLogModel = (*customAdminPerformanceLogModel)(nil)

type (
	// AdminPerformanceLogModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAdminPerformanceLogModel.
	AdminPerformanceLogModel interface {
		adminPerformanceLogModel
		// WithSession 返回一个绑定到事务 session 的新 AdminPerformanceLogModel，供 Repository.withSession 调用。
		WithSession(session sqlx.Session) AdminPerformanceLogModel
	}

	customAdminPerformanceLogModel struct {
		*defaultAdminPerformanceLogModel
	}
)

// NewAdminPerformanceLogModel returns a model for the database table.
func NewAdminPerformanceLogModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) AdminPerformanceLogModel {
	return &customAdminPerformanceLogModel{
		defaultAdminPerformanceLogModel: newAdminPerformanceLogModel(conn, c, opts...),
	}
}

// WithSession 见接口注释。table 字段直接复用，CachedConn 通过 sqlc.CachedConn.WithSession 换绑。
func (m *customAdminPerformanceLogModel) WithSession(session sqlx.Session) AdminPerformanceLogModel {
	return &customAdminPerformanceLogModel{
		defaultAdminPerformanceLogModel: &defaultAdminPerformanceLogModel{
			CachedConn: m.CachedConn.WithSession(session),
			table:      m.table,
		},
	}
}
