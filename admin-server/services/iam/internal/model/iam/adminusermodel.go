package iam

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ AdminUserModel = (*customAdminUserModel)(nil)

type (
	// AdminUserModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAdminUserModel.
	AdminUserModel interface {
		adminUserModel
		// WithSession 返回一个绑定到事务 session 的新 AdminUserModel，供 Repository.withSession 调用。
		WithSession(session sqlx.Session) AdminUserModel
	}

	customAdminUserModel struct {
		*defaultAdminUserModel
	}
)

// NewAdminUserModel returns a model for the database table.
func NewAdminUserModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) AdminUserModel {
	return &customAdminUserModel{
		defaultAdminUserModel: newAdminUserModel(conn, c, opts...),
	}
}

// WithSession 见接口注释。table 字段直接复用，CachedConn 通过 sqlc.CachedConn.WithSession 换绑。
func (m *customAdminUserModel) WithSession(session sqlx.Session) AdminUserModel {
	return &customAdminUserModel{
		defaultAdminUserModel: &defaultAdminUserModel{
			CachedConn: m.CachedConn.WithSession(session),
			table:      m.table,
		},
	}
}
