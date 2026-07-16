package iam

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ AdminRolePermissionModel = (*customAdminRolePermissionModel)(nil)

type (
	// AdminRolePermissionModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAdminRolePermissionModel.
	AdminRolePermissionModel interface {
		adminRolePermissionModel
		// WithSession 返回一个绑定到事务 session 的新 AdminRolePermissionModel，供 Repository.withSession 调用。
		WithSession(session sqlx.Session) AdminRolePermissionModel
	}

	customAdminRolePermissionModel struct {
		*defaultAdminRolePermissionModel
	}
)

// NewAdminRolePermissionModel returns a model for the database table.
func NewAdminRolePermissionModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) AdminRolePermissionModel {
	return &customAdminRolePermissionModel{
		defaultAdminRolePermissionModel: newAdminRolePermissionModel(conn, c, opts...),
	}
}

// WithSession 见接口注释。table 字段直接复用，CachedConn 通过 sqlc.CachedConn.WithSession 换绑。
func (m *customAdminRolePermissionModel) WithSession(session sqlx.Session) AdminRolePermissionModel {
	return &customAdminRolePermissionModel{
		defaultAdminRolePermissionModel: &defaultAdminRolePermissionModel{
			CachedConn: m.CachedConn.WithSession(session),
			table:      m.table,
		},
	}
}
