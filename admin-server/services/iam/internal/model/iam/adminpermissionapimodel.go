package iam

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ AdminPermissionApiModel = (*customAdminPermissionApiModel)(nil)

type (
	// AdminPermissionApiModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAdminPermissionApiModel.
	AdminPermissionApiModel interface {
		adminPermissionApiModel
		// WithSession 返回一个绑定到事务 session 的新 AdminPermissionApiModel，供 Repository.withSession 调用。
		WithSession(session sqlx.Session) AdminPermissionApiModel
	}

	customAdminPermissionApiModel struct {
		*defaultAdminPermissionApiModel
	}
)

// NewAdminPermissionApiModel returns a model for the database table.
func NewAdminPermissionApiModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) AdminPermissionApiModel {
	return &customAdminPermissionApiModel{
		defaultAdminPermissionApiModel: newAdminPermissionApiModel(conn, c, opts...),
	}
}

// WithSession 见接口注释。table 字段直接复用，CachedConn 通过 sqlc.CachedConn.WithSession 换绑。
func (m *customAdminPermissionApiModel) WithSession(session sqlx.Session) AdminPermissionApiModel {
	return &customAdminPermissionApiModel{
		defaultAdminPermissionApiModel: &defaultAdminPermissionApiModel{
			CachedConn: m.CachedConn.WithSession(session),
			table:      m.table,
		},
	}
}
