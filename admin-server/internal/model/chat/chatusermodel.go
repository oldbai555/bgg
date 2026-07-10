package chat

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ChatUserModel = (*customChatUserModel)(nil)

type (
	// ChatUserModel is an interface to be customized, add more methods here,
	// and implement the added methods in customChatUserModel.
	ChatUserModel interface {
		chatUserModel
		// WithSession 返回一个绑定到事务 session 的新 ChatUserModel，供 Repository.withSession 调用。
		WithSession(session sqlx.Session) ChatUserModel
	}

	customChatUserModel struct {
		*defaultChatUserModel
	}
)

// NewChatUserModel returns a model for the database table.
func NewChatUserModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) ChatUserModel {
	return &customChatUserModel{
		defaultChatUserModel: newChatUserModel(conn, c, opts...),
	}
}

// WithSession 见接口注释。table 字段直接复用，CachedConn 通过 sqlc.CachedConn.WithSession 换绑。
func (m *customChatUserModel) WithSession(session sqlx.Session) ChatUserModel {
	return &customChatUserModel{
		defaultChatUserModel: &defaultChatUserModel{
			CachedConn: m.CachedConn.WithSession(session),
			table:      m.table,
		},
	}
}
