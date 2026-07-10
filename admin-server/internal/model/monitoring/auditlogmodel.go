package monitoring

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ AuditLogModel = (*customAuditLogModel)(nil)

type (
	// AuditLogModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAuditLogModel.
	AuditLogModel interface {
		auditLogModel
		// WithSession 返回一个绑定到事务 session 的新 AuditLogModel，供 Repository.withSession 调用。
		WithSession(session sqlx.Session) AuditLogModel
	}

	customAuditLogModel struct {
		*defaultAuditLogModel
	}
)

// NewAuditLogModel returns a model for the database table.
func NewAuditLogModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) AuditLogModel {
	return &customAuditLogModel{
		defaultAuditLogModel: newAuditLogModel(conn, c, opts...),
	}
}

// WithSession 见接口注释。table 字段直接复用，CachedConn 通过 sqlc.CachedConn.WithSession 换绑。
func (m *customAuditLogModel) WithSession(session sqlx.Session) AuditLogModel {
	return &customAuditLogModel{
		defaultAuditLogModel: &defaultAuditLogModel{
			CachedConn: m.CachedConn.WithSession(session),
			table:      m.table,
		},
	}
}
