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
