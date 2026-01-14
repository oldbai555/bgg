package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ BlogArticleTagModel = (*customBlogArticleTagModel)(nil)

type (
	// BlogArticleTagModel is an interface to be customized, add more methods here,
	// and implement the added methods in customBlogArticleTagModel.
	BlogArticleTagModel interface {
		blogArticleTagModel
	}

	customBlogArticleTagModel struct {
		*defaultBlogArticleTagModel
	}
)

// NewBlogArticleTagModel returns a model for the database table.
func NewBlogArticleTagModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) BlogArticleTagModel {
	return &customBlogArticleTagModel{
		defaultBlogArticleTagModel: newBlogArticleTagModel(conn, c, opts...),
	}
}
