package iam

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ AdminUserRoleModel = (*customAdminUserRoleModel)(nil)

type (
	// AdminUserRoleModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAdminUserRoleModel.
	AdminUserRoleModel interface {
		adminUserRoleModel
		// WithSession 返回一个绑定到事务 session 的新 AdminUserRoleModel，供 Repository.withSession 调用。
		WithSession(session sqlx.Session) AdminUserRoleModel
	}

	customAdminUserRoleModel struct {
		*defaultAdminUserRoleModel
	}
)

// NewAdminUserRoleModel returns a model for the database table.
func NewAdminUserRoleModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) AdminUserRoleModel {
	return &customAdminUserRoleModel{
		defaultAdminUserRoleModel: newAdminUserRoleModel(conn, c, opts...),
	}
}

// WithSession 见接口注释。table 字段直接复用，CachedConn 通过 sqlc.CachedConn.WithSession 换绑。
func (m *customAdminUserRoleModel) WithSession(session sqlx.Session) AdminUserRoleModel {
	return &customAdminUserRoleModel{
		defaultAdminUserRoleModel: &defaultAdminUserRoleModel{
			CachedConn: m.CachedConn.WithSession(session),
			table:      m.table,
		},
	}
}
