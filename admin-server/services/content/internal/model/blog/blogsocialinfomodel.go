package blog

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
		// WithSession 返回一个绑定到事务 session 的新 BlogSocialInfoModel，供 Repository.withSession 调用。
		WithSession(session sqlx.Session) BlogSocialInfoModel
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

// WithSession 见接口注释。table 字段直接复用，CachedConn 通过 sqlc.CachedConn.WithSession 换绑。
func (m *customBlogSocialInfoModel) WithSession(session sqlx.Session) BlogSocialInfoModel {
	return &customBlogSocialInfoModel{
		defaultBlogSocialInfoModel: &defaultBlogSocialInfoModel{
			CachedConn: m.CachedConn.WithSession(session),
			table:      m.table,
		},
	}
}
