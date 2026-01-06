package model

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
