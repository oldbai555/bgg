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
		// WithSession 返回一个绑定到事务 session 的新 BlogTagModel，供 Repository.withSession 调用。
		WithSession(session sqlx.Session) BlogTagModel
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

// WithSession 见接口注释。table 字段直接复用，CachedConn 通过 sqlc.CachedConn.WithSession 换绑。
func (m *customBlogTagModel) WithSession(session sqlx.Session) BlogTagModel {
	return &customBlogTagModel{
		defaultBlogTagModel: &defaultBlogTagModel{
			CachedConn: m.CachedConn.WithSession(session),
			table:      m.table,
		},
	}
}
