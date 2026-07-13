package iam

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ AdminMenuModel = (*customAdminMenuModel)(nil)

type (
	// AdminMenuModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAdminMenuModel.
	AdminMenuModel interface {
		adminMenuModel
		// WithSession 返回一个绑定到事务 session 的新 AdminMenuModel，供 Repository.withSession 调用。
		WithSession(session sqlx.Session) AdminMenuModel
	}

	customAdminMenuModel struct {
		*defaultAdminMenuModel
	}
)

// NewAdminMenuModel returns a model for the database table.
func NewAdminMenuModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) AdminMenuModel {
	return &customAdminMenuModel{
		defaultAdminMenuModel: newAdminMenuModel(conn, c, opts...),
	}
}

// WithSession 见接口注释。table 字段直接复用，CachedConn 通过 sqlc.CachedConn.WithSession 换绑。
func (m *customAdminMenuModel) WithSession(session sqlx.Session) AdminMenuModel {
	return &customAdminMenuModel{
		defaultAdminMenuModel: &defaultAdminMenuModel{
			CachedConn: m.CachedConn.WithSession(session),
			table:      m.table,
		},
	}
}
