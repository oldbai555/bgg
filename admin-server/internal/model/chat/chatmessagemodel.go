package chat

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ChatMessageModel = (*customChatMessageModel)(nil)

type (
	// ChatMessageModel is an interface to be customized, add more methods here,
	// and implement the added methods in customChatMessageModel.
	ChatMessageModel interface {
		chatMessageModel
		// WithSession 返回一个绑定到事务 session 的新 ChatMessageModel，供 Repository.withSession 调用。
		WithSession(session sqlx.Session) ChatMessageModel
	}

	customChatMessageModel struct {
		*defaultChatMessageModel
	}
)

// NewChatMessageModel returns a model for the database table.
func NewChatMessageModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) ChatMessageModel {
	return &customChatMessageModel{
		defaultChatMessageModel: newChatMessageModel(conn, c, opts...),
	}
}

// WithSession 见接口注释。table 字段直接复用，CachedConn 通过 sqlc.CachedConn.WithSession 换绑。
func (m *customChatMessageModel) WithSession(session sqlx.Session) ChatMessageModel {
	return &customChatMessageModel{
		defaultChatMessageModel: &defaultChatMessageModel{
			CachedConn: m.CachedConn.WithSession(session),
			table:      m.table,
		},
	}
}
