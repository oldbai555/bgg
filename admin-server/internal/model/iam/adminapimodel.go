package iam

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ AdminApiModel = (*customAdminApiModel)(nil)

type (
	// AdminApiModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAdminApiModel.
	AdminApiModel interface {
		adminApiModel
		// WithSession 返回一个绑定到事务 session 的新 AdminApiModel，供 Repository.withSession 调用。
		WithSession(session sqlx.Session) AdminApiModel
	}

	customAdminApiModel struct {
		*defaultAdminApiModel
	}
)

// NewAdminApiModel returns a model for the database table.
func NewAdminApiModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) AdminApiModel {
	return &customAdminApiModel{
		defaultAdminApiModel: newAdminApiModel(conn, c, opts...),
	}
}

// WithSession 见接口注释。table 字段直接复用，CachedConn 通过 sqlc.CachedConn.WithSession 换绑。
func (m *customAdminApiModel) WithSession(session sqlx.Session) AdminApiModel {
	return &customAdminApiModel{
		defaultAdminApiModel: &defaultAdminApiModel{
			CachedConn: m.CachedConn.WithSession(session),
			table:      m.table,
		},
	}
}
