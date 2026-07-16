package iam

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ AdminUserThirdPartyModel = (*customAdminUserThirdPartyModel)(nil)

type (
	// AdminUserThirdPartyModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAdminUserThirdPartyModel.
	AdminUserThirdPartyModel interface {
		adminUserThirdPartyModel
		// WithSession 返回一个绑定到事务 session 的新 AdminUserThirdPartyModel，供 Repository.withSession 调用。
		WithSession(session sqlx.Session) AdminUserThirdPartyModel
	}

	customAdminUserThirdPartyModel struct {
		*defaultAdminUserThirdPartyModel
	}
)

// NewAdminUserThirdPartyModel returns a model for the database table.
func NewAdminUserThirdPartyModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) AdminUserThirdPartyModel {
	return &customAdminUserThirdPartyModel{
		defaultAdminUserThirdPartyModel: newAdminUserThirdPartyModel(conn, c, opts...),
	}
}

// WithSession 见接口注释。table 字段直接复用，CachedConn 通过 sqlc.CachedConn.WithSession 换绑。
func (m *customAdminUserThirdPartyModel) WithSession(session sqlx.Session) AdminUserThirdPartyModel {
	return &customAdminUserThirdPartyModel{
		defaultAdminUserThirdPartyModel: &defaultAdminUserThirdPartyModel{
			CachedConn: m.CachedConn.WithSession(session),
			table:      m.table,
		},
	}
}
