package sdk

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ SdkKeyModel = (*customSdkKeyModel)(nil)

type (
	// SdkKeyModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSdkKeyModel.
	SdkKeyModel interface {
		sdkKeyModel
		// WithSession 返回一个绑定到事务 session 的新 SdkKeyModel，供 Repository.withSession 调用。
		WithSession(session sqlx.Session) SdkKeyModel
	}

	customSdkKeyModel struct {
		*defaultSdkKeyModel
	}
)

// NewSdkKeyModel returns a model for the database table.
func NewSdkKeyModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) SdkKeyModel {
	return &customSdkKeyModel{
		defaultSdkKeyModel: newSdkKeyModel(conn, c, opts...),
	}
}

// WithSession 见接口注释。table 字段直接复用，CachedConn 通过 sqlc.CachedConn.WithSession 换绑。
func (m *customSdkKeyModel) WithSession(session sqlx.Session) SdkKeyModel {
	return &customSdkKeyModel{
		defaultSdkKeyModel: &defaultSdkKeyModel{
			CachedConn: m.CachedConn.WithSession(session),
			table:      m.table,
		},
	}
}
