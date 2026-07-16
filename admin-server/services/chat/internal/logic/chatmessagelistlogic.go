package logic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"postapocgame/admin-server/pkg/errs"
	iamcallbackpb "postapocgame/admin-server/pkg/iamcallback/pb"
	"postapocgame/admin-server/services/chat/chat"
	"postapocgame/admin-server/services/chat/internal/svc"
)

type ChatMessageListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewChatMessageListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatMessageListLogic {
	return &ChatMessageListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ChatMessageList 迁移自 internal/logic/chat/chat/chat_message_list_logic.go。
func (l *ChatMessageListLogic) ChatMessageList(in *chat.ChatMessageListRequest) (*chat.ChatMessageListResponse, error) {
	page, pageSize := normalizePage(in.Page, in.PageSize, 20, 100)

	list, total, err := l.svcCtx.ChatMessage.FindByChatID(l.ctx, page, pageSize, in.ChatId)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "查询聊天消息列表失败", err))
	}

	items := make([]*chat.ChatMessageItem, 0, len(list))
	for _, msg := range list {
		fromUserName := ""
		if profile, err := l.svcCtx.IamCallback.GetUserProfile(l.ctx, &iamcallbackpb.GetUserProfileRequest{UserId: msg.FromUserId}); err == nil && profile.Exists {
			fromUserName = profile.Username
		}

		items = append(items, &chat.ChatMessageItem{
			Id:           msg.Id,
			ChatId:       msg.ChatId,
			FromUserId:   msg.FromUserId,
			FromUserName: fromUserName,
			Content:      msg.Content,
			MessageType:  msg.MessageType,
			Status:       msg.Status,
			CreatedAt:    msg.CreatedAt,
		})
	}

	return &chat.ChatMessageListResponse{Total: total, List: items}, nil
}

// normalizePage 迁移自 gateway internal/logic/logicutil.NormalizePage——chat-rpc 的
// internal/ 不反向依赖 gateway 的 internal/logicutil（见 16-rpc-conventions.md 的服务边界
// 原则，task-rpc/sdk-rpc 拆分时已经确立的先例：分页兜底逻辑各服务自己内联一份，不共享）。
func normalizePage(page, pageSize, defaultSize, maxSize int64) (int64, int64) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = defaultSize
	}
	if pageSize > maxSize {
		pageSize = maxSize
	}
	return page, pageSize
}
