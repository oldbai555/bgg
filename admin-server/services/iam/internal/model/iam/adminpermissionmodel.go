package iam

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ AdminPermissionModel = (*customAdminPermissionModel)(nil)

type (
	// AdminPermissionModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAdminPermissionModel.
	AdminPermissionModel interface {
		adminPermissionModel
		// WithSession 返回一个绑定到事务 session 的新 AdminPermissionModel，供 Repository.withSession 调用。
		WithSession(session sqlx.Session) AdminPermissionModel
	}

	customAdminPermissionModel struct {
		*defaultAdminPermissionModel
	}
)

// NewAdminPermissionModel returns a model for the database table.
func NewAdminPermissionModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) AdminPermissionModel {
	return &customAdminPermissionModel{
		defaultAdminPermissionModel: newAdminPermissionModel(conn, c, opts...),
	}
}

// WithSession 见接口注释。table 字段直接复用，CachedConn 通过 sqlc.CachedConn.WithSession 换绑。
func (m *customAdminPermissionModel) WithSession(session sqlx.Session) AdminPermissionModel {
	return &customAdminPermissionModel{
		defaultAdminPermissionModel: &defaultAdminPermissionModel{
			CachedConn: m.CachedConn.WithSession(session),
			table:      m.table,
		},
	}
}
