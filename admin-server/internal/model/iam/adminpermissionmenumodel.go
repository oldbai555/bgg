package iam

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ AdminPermissionMenuModel = (*customAdminPermissionMenuModel)(nil)

type (
	// AdminPermissionMenuModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAdminPermissionMenuModel.
	AdminPermissionMenuModel interface {
		adminPermissionMenuModel
		// WithSession 返回一个绑定到事务 session 的新 AdminPermissionMenuModel，供 Repository.withSession 调用。
		WithSession(session sqlx.Session) AdminPermissionMenuModel
	}

	customAdminPermissionMenuModel struct {
		*defaultAdminPermissionMenuModel
	}
)

// NewAdminPermissionMenuModel returns a model for the database table.
func NewAdminPermissionMenuModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) AdminPermissionMenuModel {
	return &customAdminPermissionMenuModel{
		defaultAdminPermissionMenuModel: newAdminPermissionMenuModel(conn, c, opts...),
	}
}

// WithSession 见接口注释。table 字段直接复用，CachedConn 通过 sqlc.CachedConn.WithSession 换绑。
func (m *customAdminPermissionMenuModel) WithSession(session sqlx.Session) AdminPermissionMenuModel {
	return &customAdminPermissionMenuModel{
		defaultAdminPermissionMenuModel: &defaultAdminPermissionMenuModel{
			CachedConn: m.CachedConn.WithSession(session),
			table:      m.table,
		},
	}
}
