package sdk

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ SdkCallLogModel = (*customSdkCallLogModel)(nil)

type (
	// SdkCallLogModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSdkCallLogModel.
	SdkCallLogModel interface {
		sdkCallLogModel
		// WithSession 返回一个绑定到事务 session 的新 SdkCallLogModel，供 Repository.withSession 调用。
		WithSession(session sqlx.Session) SdkCallLogModel
	}

	customSdkCallLogModel struct {
		*defaultSdkCallLogModel
	}
)

// NewSdkCallLogModel returns a model for the database table.
func NewSdkCallLogModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) SdkCallLogModel {
	return &customSdkCallLogModel{
		defaultSdkCallLogModel: newSdkCallLogModel(conn, c, opts...),
	}
}

// WithSession 见接口注释。table 字段直接复用，CachedConn 通过 sqlc.CachedConn.WithSession 换绑。
func (m *customSdkCallLogModel) WithSession(session sqlx.Session) SdkCallLogModel {
	return &customSdkCallLogModel{
		defaultSdkCallLogModel: &defaultSdkCallLogModel{
			CachedConn: m.CachedConn.WithSession(session),
			table:      m.table,
		},
	}
}
