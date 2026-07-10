package monitoring

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ AdminOperationLogModel = (*customAdminOperationLogModel)(nil)

type (
	// AdminOperationLogModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAdminOperationLogModel.
	AdminOperationLogModel interface {
		adminOperationLogModel
		// WithSession 返回一个绑定到事务 session 的新 AdminOperationLogModel，供 Repository.withSession 调用。
		WithSession(session sqlx.Session) AdminOperationLogModel
	}

	customAdminOperationLogModel struct {
		*defaultAdminOperationLogModel
	}
)

// NewAdminOperationLogModel returns a model for the database table.
func NewAdminOperationLogModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) AdminOperationLogModel {
	return &customAdminOperationLogModel{
		defaultAdminOperationLogModel: newAdminOperationLogModel(conn, c, opts...),
	}
}

// WithSession 见接口注释。table 字段直接复用，CachedConn 通过 sqlc.CachedConn.WithSession 换绑。
func (m *customAdminOperationLogModel) WithSession(session sqlx.Session) AdminOperationLogModel {
	return &customAdminOperationLogModel{
		defaultAdminOperationLogModel: &defaultAdminOperationLogModel{
			CachedConn: m.CachedConn.WithSession(session),
			table:      m.table,
		},
	}
}
