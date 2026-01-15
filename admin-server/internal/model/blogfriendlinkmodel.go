package model

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
