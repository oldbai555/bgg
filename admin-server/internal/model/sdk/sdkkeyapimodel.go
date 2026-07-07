package sdk

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ SdkKeyApiModel = (*customSdkKeyApiModel)(nil)

type (
	// SdkKeyApiModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSdkKeyApiModel.
	SdkKeyApiModel interface {
		sdkKeyApiModel
	}

	customSdkKeyApiModel struct {
		*defaultSdkKeyApiModel
	}
)

// NewSdkKeyApiModel returns a model for the database table.
func NewSdkKeyApiModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) SdkKeyApiModel {
	return &customSdkKeyApiModel{
		defaultSdkKeyApiModel: newSdkKeyApiModel(conn, c, opts...),
	}
}
