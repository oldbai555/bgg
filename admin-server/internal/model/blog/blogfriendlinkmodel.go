package blog

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ BlogFriendLinkModel = (*customBlogFriendLinkModel)(nil)

type (
	// BlogFriendLinkModel is an interface to be customized, add more methods here,
	// and implement the added methods in customBlogFriendLinkModel.
	BlogFriendLinkModel interface {
		blogFriendLinkModel
		// WithSession 返回一个绑定到事务 session 的新 BlogFriendLinkModel，供 Repository.withSession 调用。
		WithSession(session sqlx.Session) BlogFriendLinkModel
	}

	customBlogFriendLinkModel struct {
		*defaultBlogFriendLinkModel
	}
)

// NewBlogFriendLinkModel returns a model for the database table.
func NewBlogFriendLinkModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) BlogFriendLinkModel {
	return &customBlogFriendLinkModel{
		defaultBlogFriendLinkModel: newBlogFriendLinkModel(conn, c, opts...),
	}
}

// WithSession 见接口注释。table 字段直接复用，CachedConn 通过 sqlc.CachedConn.WithSession 换绑。
func (m *customBlogFriendLinkModel) WithSession(session sqlx.Session) BlogFriendLinkModel {
	return &customBlogFriendLinkModel{
		defaultBlogFriendLinkModel: &defaultBlogFriendLinkModel{
			CachedConn: m.CachedConn.WithSession(session),
			table:      m.table,
		},
	}
}
