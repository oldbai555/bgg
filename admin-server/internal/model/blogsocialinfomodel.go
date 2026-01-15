package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ BlogSocialInfoModel = (*customBlogSocialInfoModel)(nil)

type (
	// BlogSocialInfoModel is an interface to be customized, add more methods here,
	// and implement the added methods in customBlogSocialInfoModel.
	BlogSocialInfoModel interface {
		blogSocialInfoModel
	}

	customBlogSocialInfoModel struct {
		*defaultBlogSocialInfoModel
	}
)

// NewBlogSocialInfoModel returns a model for the database table.
func NewBlogSocialInfoModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) BlogSocialInfoModel {
	return &customBlogSocialInfoModel{
		defaultBlogSocialInfoModel: newBlogSocialInfoModel(conn, c, opts...),
	}
}
