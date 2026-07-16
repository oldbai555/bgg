// Package consts 复制自 internal/consts/chat.go，按 16-rpc-conventions.md 第 6 节的既定
// 策略：直接复制到各服务自己的 internal/consts，不做成共享包（量很小，维护成本可忽略）。
package consts

const (
	// DefaultGroupChatID 默认企业群组的 chat_id（新用户自动加入）
	DefaultGroupChatID uint64 = 1

	// ChatTypePrivate 私聊
	ChatTypePrivate int64 = 1
	// ChatTypeGroup 群组
	ChatTypeGroup int64 = 2
)
