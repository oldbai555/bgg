package sdk

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ SdkInterfaceModel = (*customSdkInterfaceModel)(nil)

type (
	// SdkInterfaceModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSdkInterfaceModel.
	SdkInterfaceModel interface {
		sdkInterfaceModel
		// WithSession 返回一个绑定到事务 session 的新 SdkInterfaceModel，供 Repository.withSession 调用。
		WithSession(session sqlx.Session) SdkInterfaceModel
	}

	customSdkInterfaceModel struct {
		*defaultSdkInterfaceModel
	}
)

// NewSdkInterfaceModel returns a model for the database table.
func NewSdkInterfaceModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) SdkInterfaceModel {
	return &customSdkInterfaceModel{
		defaultSdkInterfaceModel: newSdkInterfaceModel(conn, c, opts...),
	}
}

// WithSession 见接口注释。table 字段直接复用，CachedConn 通过 sqlc.CachedConn.WithSession 换绑。
func (m *customSdkInterfaceModel) WithSession(session sqlx.Session) SdkInterfaceModel {
	return &customSdkInterfaceModel{
		defaultSdkInterfaceModel: &defaultSdkInterfaceModel{
			CachedConn: m.CachedConn.WithSession(session),
			table:      m.table,
		},
	}
}
