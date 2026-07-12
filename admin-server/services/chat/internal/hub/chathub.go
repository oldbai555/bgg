// Package hub 从 internal/hub/chathub.go 移植而来。连接表数据结构（clients map[uint64]*Client）
// 原样保留，唯一的结构性改动是"连接"的类型从 *websocket.Conn 换成 gRPC 双向流
// （chat.Chat_StreamServer）——gateway 继续终结 WebSocket 连接，chat-rpc 不再需要
// gorilla/websocket 依赖，也不需要原来的 ReadPump/WritePump（读写超时、ping/pong 帧这些
// TCP/WS 协议层细节现在是 gateway 桥接 handler 的职责，见
// internal/handler/chat/chatwshandler.go），Hub 这一层只关心"往哪个用户的 Send 队列塞一帧"。
//
// ChatMessage 类型定义原样从 internal/hub/chathub.go 复制（JSON 字段名/tag 完全一致），
// 保证 WS wire 格式不变，前端不需要跟着这次拆分改任何解析逻辑，见
// services/chat/rpc/chat.proto 里 MessageFrame 的注释。
package hub

import (
	"encoding/json"
	"sync"

	"github.com/zeromicro/go-zero/core/logx"

	"postapocgame/admin-server/services/chat/chat"
)

// Client 表示一条注册进 Hub 的 gRPC 流式连接。
type Client struct {
	Hub      *ChatHub
	Stream   chat.Chat_StreamServer
	Send     chan *chat.ServerFrame
	UserID   uint64
	Username string
}

// ChatHub 管理所有在线连接和消息广播。
type ChatHub struct {
	clients    map[uint64]*Client
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

func NewChatHub() *ChatHub {
	return &ChatHub{
		clients:    make(map[uint64]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *ChatHub) Register() chan<- *Client   { return h.register }
func (h *ChatHub) Unregister() chan<- *Client { return h.unregister }

// Run 启动 Hub 的注册/注销事件循环，需要在独立 goroutine 里跑（和原 ChatHub.Run 用法一致）。
func (h *ChatHub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.UserID] = client
			h.mu.Unlock()
			logx.Infof("chat-rpc 连接注册: UserID=%d, Username=%s", client.UserID, client.Username)

		case client := <-h.unregister:
			h.mu.Lock()
			if existing, ok := h.clients[client.UserID]; ok && existing == client {
				delete(h.clients, client.UserID)
				close(client.Send)
			}
			h.mu.Unlock()
			logx.Infof("chat-rpc 连接注销: UserID=%d, Username=%s", client.UserID, client.Username)
		}
	}
}

// SendToUser 向指定用户的在线连接发送一帧，用户不在线时返回 false（调用方按"尽力而为"处理，
// 不当作错误）。
func (h *ChatHub) SendToUser(userID uint64, frame *chat.ServerFrame) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	client, ok := h.clients[userID]
	if !ok {
		return false
	}
	select {
	case client.Send <- frame:
		return true
	default:
		return false // Send 队列已满，丢弃这一帧，不阻塞 Hub
	}
}

// BroadcastToChat 向指定聊天的多个成员在线连接分别发送同一帧。
func (h *ChatHub) BroadcastToChat(userIDs []uint64, frame *chat.ServerFrame) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, userID := range userIDs {
		if client, ok := h.clients[userID]; ok {
			select {
			case client.Send <- frame:
			default:
			}
		}
	}
}

func (h *ChatHub) IsUserOnline(userID uint64) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	_, ok := h.clients[userID]
	return ok
}

// OnlineUserCount 供 GetOnlineUserCount RPC 使用，取代原来 gateway 直接读
// ChatHub.GetOnlineUsers() 再取 len 的写法。
func (h *ChatHub) OnlineUserCount() int64 {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return int64(len(h.clients))
}

// ChatMessage 是 WS wire 格式（JSON），字段名/tag 与拆分前的 internal/hub.ChatMessage
// 逐一对齐，保证前端不用改。
type ChatMessage struct {
	Type      string `json:"type"`
	FromID    uint64 `json:"fromId"`
	FromName  string `json:"fromName"`
	ToID      uint64 `json:"toId"`
	RoomID    string `json:"roomId"`
	ChatID    uint64 `json:"chatId"`
	Content   string `json:"content"`
	MessageID uint64 `json:"messageId"`
	CreatedAt int64  `json:"createdAt"`
	TaskID    string `json:"taskId,omitempty"`
	TaskName  string `json:"taskName,omitempty"`
	Progress  int    `json:"progress,omitempty"`
	Status    string `json:"status,omitempty"`
	Title     string `json:"title,omitempty"`
	Level     string `json:"level,omitempty"`
}

// NewMessageFrame 把 ChatMessage 编码进 MessageFrame.PayloadJson，编码失败时返回 nil
// （调用方应记日志跳过，不阻塞主流程，与原 BroadcastToChat 遇到序列化错误时"忽略"的既有
// 语义一致）。
func NewMessageFrame(msg *ChatMessage) *chat.ServerFrame {
	payload, err := json.Marshal(msg)
	if err != nil {
		logx.Errorf("ChatMessage 序列化失败: %v", err)
		return nil
	}
	return &chat.ServerFrame{Payload: &chat.ServerFrame_Message{
		Message: &chat.MessageFrame{PayloadJson: string(payload)},
	}}
}
