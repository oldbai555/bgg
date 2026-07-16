package logic

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/chat/chat"
	chatmodel "postapocgame/admin-server/services/chat/internal/model/chat"
	"postapocgame/admin-server/services/chat/internal/svc"

	chathub "postapocgame/admin-server/services/chat/internal/hub"
)

type ChatMessageSendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewChatMessageSendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatMessageSendLogic {
	return &ChatMessageSendLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ChatMessageSend 迁移自 internal/logic/chat/chat/chat_message_send_logic.go：验证成员资格 →
// 持久化 → 广播给聊天里所有在线成员，逻辑原样保留。operator_user_id/operator_username 由
// gateway 侧显式传入（chat-rpc 不解析 JWT，见 services/chat/rpc/chat.proto 顶部注释）。
func (l *ChatMessageSendLogic) ChatMessageSend(in *chat.ChatMessageSendRequest) (*chat.ChatMessageSendResponse, error) {
	if in.Content == "" {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "消息内容不能为空"))
	}
	if in.ChatId == 0 {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "聊天ID不能为空"))
	}

	c, err := l.svcCtx.Chat.FindByID(l.ctx, in.ChatId)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeBadRequest, "聊天不存在", err))
	}
	if c.DeletedAt != 0 {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "聊天已删除"))
	}

	chatUsers, err := l.svcCtx.ChatUser.FindByChatID(l.ctx, in.ChatId)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "查询聊天成员失败", err))
	}
	hasPermission := false
	for _, cu := range chatUsers {
		if cu.UserId == in.OperatorUserId {
			hasPermission = true
			break
		}
	}
	if !hasPermission {
		return nil, toGRPCStatus(errs.New(errs.CodeForbidden, "您不在该聊天中"))
	}

	messageType := in.MessageType
	if messageType == 0 {
		messageType = 1
	}

	now := time.Now().Unix()
	message := &chatmodel.ChatMessage{
		ChatId:      in.ChatId,
		FromUserId:  in.OperatorUserId,
		Content:     in.Content,
		MessageType: messageType,
		Status:      1,
		CreatedAt:   now,
		UpdatedAt:   now,
		DeletedAt:   0,
	}
	if err := l.svcCtx.ChatMessage.Create(l.ctx, message); err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "发送消息失败", err))
	}

	userIDs := make([]uint64, 0, len(chatUsers))
	for _, cu := range chatUsers {
		userIDs = append(userIDs, cu.UserId)
	}
	frame := chathub.NewMessageFrame(&chathub.ChatMessage{
		Type:      "chat",
		FromID:    in.OperatorUserId,
		FromName:  in.OperatorUsername,
		ChatID:    in.ChatId,
		Content:   in.Content,
		MessageID: message.Id,
		CreatedAt: message.CreatedAt,
	})
	if frame != nil {
		l.svcCtx.Hub.BroadcastToChat(userIDs, frame)
	}

	return &chat.ChatMessageSendResponse{Id: message.Id}, nil
}
