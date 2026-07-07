package blog

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ BlogTagModel = (*customBlogTagModel)(nil)

type (
	// BlogTagModel is an interface to be customized, add more methods here,
	// and implement the added methods in customBlogTagModel.
	BlogTagModel interface {
		blogTagModel
	}

	customBlogTagModel struct {
		*defaultBlogTagModel
	}
)

// NewBlogTagModel returns a model for the database table.
func NewBlogTagModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) BlogTagModel {
	return &customBlogTagModel{
		defaultBlogTagModel: newBlogTagModel(conn, c, opts...),
	}
}
