package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ AdminTaskModel = (*customAdminTaskModel)(nil)

type (
	// AdminTaskModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAdminTaskModel.
	AdminTaskModel interface {
		adminTaskModel
	}

	customAdminTaskModel struct {
		*defaultAdminTaskModel
	}
)

// NewAdminTaskModel returns a model for the database table.
func NewAdminTaskModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) AdminTaskModel {
	return &customAdminTaskModel{
		defaultAdminTaskModel: newAdminTaskModel(conn, c, opts...),
	}
}
