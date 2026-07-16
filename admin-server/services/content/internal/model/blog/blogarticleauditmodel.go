package blog

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ BlogArticleAuditModel = (*customBlogArticleAuditModel)(nil)

type (
	// BlogArticleAuditModel is an interface to be customized, add more methods here,
	// and implement the added methods in customBlogArticleAuditModel.
	BlogArticleAuditModel interface {
		blogArticleAuditModel
		// WithSession 返回一个绑定到事务 session 的新 BlogArticleAuditModel，供 Repository.withSession 调用。
		WithSession(session sqlx.Session) BlogArticleAuditModel
	}

	customBlogArticleAuditModel struct {
		*defaultBlogArticleAuditModel
	}
)

// NewBlogArticleAuditModel returns a model for the database table.
func NewBlogArticleAuditModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) BlogArticleAuditModel {
	return &customBlogArticleAuditModel{
		defaultBlogArticleAuditModel: newBlogArticleAuditModel(conn, c, opts...),
	}
}

// WithSession 见接口注释。table 字段直接复用，CachedConn 通过 sqlc.CachedConn.WithSession 换绑。
func (m *customBlogArticleAuditModel) WithSession(session sqlx.Session) BlogArticleAuditModel {
	return &customBlogArticleAuditModel{
		defaultBlogArticleAuditModel: &defaultBlogArticleAuditModel{
			CachedConn: m.CachedConn.WithSession(session),
			table:      m.table,
		},
	}
}
