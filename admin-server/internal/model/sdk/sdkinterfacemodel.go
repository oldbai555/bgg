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
