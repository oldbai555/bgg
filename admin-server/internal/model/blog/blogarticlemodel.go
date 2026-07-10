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
		// WithSession 返回一个绑定到事务 session 的新 BlogArticleModel，供 Repository.withSession 调用。
		WithSession(session sqlx.Session) BlogArticleModel
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

// WithSession 见接口注释。table 字段直接复用，CachedConn 通过 sqlc.CachedConn.WithSession 换绑。
func (m *customBlogArticleModel) WithSession(session sqlx.Session) BlogArticleModel {
	return &customBlogArticleModel{
		defaultBlogArticleModel: &defaultBlogArticleModel{
			CachedConn: m.CachedConn.WithSession(session),
			table:      m.table,
		},
	}
}
