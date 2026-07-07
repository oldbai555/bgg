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
