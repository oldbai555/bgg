package iam

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ AdminRoleModel = (*customAdminRoleModel)(nil)

type (
	// AdminRoleModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAdminRoleModel.
	AdminRoleModel interface {
		adminRoleModel
		// WithSession 返回一个绑定到事务 session 的新 AdminRoleModel，供 Repository.withSession 调用。
		WithSession(session sqlx.Session) AdminRoleModel
	}

	customAdminRoleModel struct {
		*defaultAdminRoleModel
	}
)

// NewAdminRoleModel returns a model for the database table.
func NewAdminRoleModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) AdminRoleModel {
	return &customAdminRoleModel{
		defaultAdminRoleModel: newAdminRoleModel(conn, c, opts...),
	}
}

// WithSession 见接口注释。table 字段直接复用，CachedConn 通过 sqlc.CachedConn.WithSession 换绑。
func (m *customAdminRoleModel) WithSession(session sqlx.Session) AdminRoleModel {
	return &customAdminRoleModel{
		defaultAdminRoleModel: &defaultAdminRoleModel{
			CachedConn: m.CachedConn.WithSession(session),
			table:      m.table,
		},
	}
}
