package system

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ AdminNotificationModel = (*customAdminNotificationModel)(nil)

type (
	// AdminNotificationModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAdminNotificationModel.
	AdminNotificationModel interface {
		adminNotificationModel
		// WithSession 返回一个绑定到事务 session 的新 AdminNotificationModel，供 Repository.withSession 调用。
		WithSession(session sqlx.Session) AdminNotificationModel
	}

	customAdminNotificationModel struct {
		*defaultAdminNotificationModel
	}
)

// NewAdminNotificationModel returns a model for the database table.
func NewAdminNotificationModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) AdminNotificationModel {
	return &customAdminNotificationModel{
		defaultAdminNotificationModel: newAdminNotificationModel(conn, c, opts...),
	}
}

// WithSession 见接口注释。table 字段直接复用，CachedConn 通过 sqlc.CachedConn.WithSession 换绑。
func (m *customAdminNotificationModel) WithSession(session sqlx.Session) AdminNotificationModel {
	return &customAdminNotificationModel{
		defaultAdminNotificationModel: &defaultAdminNotificationModel{
			CachedConn: m.CachedConn.WithSession(session),
			table:      m.table,
		},
	}
}
