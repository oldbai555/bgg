package iam

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ AdminDepartmentModel = (*customAdminDepartmentModel)(nil)

type (
	// AdminDepartmentModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAdminDepartmentModel.
	AdminDepartmentModel interface {
		adminDepartmentModel
		// WithSession 返回一个绑定到事务 session 的新 AdminDepartmentModel，供 Repository.withSession 调用。
		WithSession(session sqlx.Session) AdminDepartmentModel
	}

	customAdminDepartmentModel struct {
		*defaultAdminDepartmentModel
	}
)

// NewAdminDepartmentModel returns a model for the database table.
func NewAdminDepartmentModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) AdminDepartmentModel {
	return &customAdminDepartmentModel{
		defaultAdminDepartmentModel: newAdminDepartmentModel(conn, c, opts...),
	}
}

// WithSession 见接口注释。table 字段直接复用，CachedConn 通过 sqlc.CachedConn.WithSession 换绑。
func (m *customAdminDepartmentModel) WithSession(session sqlx.Session) AdminDepartmentModel {
	return &customAdminDepartmentModel{
		defaultAdminDepartmentModel: &defaultAdminDepartmentModel{
			CachedConn: m.CachedConn.WithSession(session),
			table:      m.table,
		},
	}
}
