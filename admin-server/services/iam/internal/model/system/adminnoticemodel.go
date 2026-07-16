package system

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ AdminNoticeModel = (*customAdminNoticeModel)(nil)

type (
	// AdminNoticeModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAdminNoticeModel.
	AdminNoticeModel interface {
		adminNoticeModel
		// WithSession 返回一个绑定到事务 session 的新 AdminNoticeModel，供 Repository.withSession 调用。
		WithSession(session sqlx.Session) AdminNoticeModel
	}

	customAdminNoticeModel struct {
		*defaultAdminNoticeModel
	}
)

// NewAdminNoticeModel returns a model for the database table.
func NewAdminNoticeModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) AdminNoticeModel {
	return &customAdminNoticeModel{
		defaultAdminNoticeModel: newAdminNoticeModel(conn, c, opts...),
	}
}

// WithSession 见接口注释。table 字段直接复用，CachedConn 通过 sqlc.CachedConn.WithSession 换绑。
func (m *customAdminNoticeModel) WithSession(session sqlx.Session) AdminNoticeModel {
	return &customAdminNoticeModel{
		defaultAdminNoticeModel: &defaultAdminNoticeModel{
			CachedConn: m.CachedConn.WithSession(session),
			table:      m.table,
		},
	}
}
