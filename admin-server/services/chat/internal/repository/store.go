// Package repository 从 internal/repository/chat/ 原样搬迁而来。唯一的结构性改动是把三个
// repository（ChatRepository/ChatUserRepository/ChatMessageRepository）原来共享的
// *repository.Repository（单体聚合了全部 9 个业务域 Model 的大句柄）换成这里的 *Store——
// chat-rpc 从第一天起只有 chat/chat_user/chat_message 三张表，不该也不能继续持有指向其它
// 域的句柄。和 services/sdk/internal/repository/store.go 同一个模式。
package repository

import (
	"context"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"

	chatmodel "postapocgame/admin-server/services/chat/internal/model/chat"
)

// Store 聚合 chat-rpc 自己需要的全部 Model，供 ChatRepository/ChatUserRepository/
// ChatMessageRepository 共用。
type Store struct {
	DB               sqlx.SqlConn
	ChatModel        chatmodel.ChatModel
	ChatUserModel    chatmodel.ChatUserModel
	ChatMessageModel chatmodel.ChatMessageModel
}

func NewStore(conn sqlx.SqlConn, cacheConf cache.CacheConf) *Store {
	return &Store{
		DB:               conn,
		ChatModel:        chatmodel.NewChatModel(conn, cacheConf),
		ChatUserModel:    chatmodel.NewChatUserModel(conn, cacheConf),
		ChatMessageModel: chatmodel.NewChatMessageModel(conn, cacheConf),
	}
}

// Transact 在单个 MySQL 事务内执行 fn，用法和 internal/repository/repository.go 的
// Repository.Transact 完全同构（chat-rpc 自己的小号版本，只换绑这三个 Model）。
func (s *Store) Transact(ctx context.Context, fn func(ctx context.Context, txStore *Store) error) error {
	return s.DB.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		return fn(ctx, s.withSession(session))
	})
}

func (s *Store) withSession(session sqlx.Session) *Store {
	return &Store{
		DB:               sqlx.NewSqlConnFromSession(session),
		ChatModel:        s.ChatModel.WithSession(session),
		ChatUserModel:    s.ChatUserModel.WithSession(session),
		ChatMessageModel: s.ChatMessageModel.WithSession(session),
	}
}
