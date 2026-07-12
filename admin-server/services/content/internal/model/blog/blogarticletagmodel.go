package blog

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
		// WithSession 返回一个绑定到事务 session 的新 BlogArticleTagModel，供 Repository.withSession 调用。
		WithSession(session sqlx.Session) BlogArticleTagModel
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

// WithSession 见接口注释。table 字段直接复用，CachedConn 通过 sqlc.CachedConn.WithSession 换绑。
func (m *customBlogArticleTagModel) WithSession(session sqlx.Session) BlogArticleTagModel {
	return &customBlogArticleTagModel{
		defaultBlogArticleTagModel: &defaultBlogArticleTagModel{
			CachedConn: m.CachedConn.WithSession(session),
			table:      m.table,
		},
	}
}
