package chat

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ChatModel = (*customChatModel)(nil)

type (
	// ChatModel is an interface to be customized, add more methods here,
	// and implement the added methods in customChatModel.
	ChatModel interface {
		chatModel
		// WithSession 返回一个绑定到事务 session 的新 ChatModel，供 Repository.withSession 调用。
		WithSession(session sqlx.Session) ChatModel
	}

	customChatModel struct {
		*defaultChatModel
	}
)

// NewChatModel returns a model for the database table.
func NewChatModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) ChatModel {
	return &customChatModel{
		defaultChatModel: newChatModel(conn, c, opts...),
	}
}

// WithSession 见接口注释。table 字段直接复用，CachedConn 通过 sqlc.CachedConn.WithSession 换绑。
func (m *customChatModel) WithSession(session sqlx.Session) ChatModel {
	return &customChatModel{
		defaultChatModel: &defaultChatModel{
			CachedConn: m.CachedConn.WithSession(session),
			table:      m.table,
		},
	}
}
