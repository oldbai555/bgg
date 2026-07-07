package blog

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ BlogArticleModel = (*customBlogArticleModel)(nil)

type (
	// BlogArticleModel is an interface to be customized, add more methods here,
	// and implement the added methods in customBlogArticleModel.
	BlogArticleModel interface {
		blogArticleModel
	}

	customBlogArticleModel struct {
		*defaultBlogArticleModel
	}
)

// NewBlogArticleModel returns a model for the database table.
func NewBlogArticleModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) BlogArticleModel {
	return &customBlogArticleModel{
		defaultBlogArticleModel: newBlogArticleModel(conn, c, opts...),
	}
}
