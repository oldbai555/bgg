package logic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"postapocgame/admin-server/pkg/errs"
	iamcallbackpb "postapocgame/admin-server/pkg/iamcallback/pb"
	"postapocgame/admin-server/services/chat/chat"
	"postapocgame/admin-server/services/chat/internal/svc"
)

type ChatMessageListAdminLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewChatMessageListAdminLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatMessageListAdminLogic {
	return &ChatMessageListAdminLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ChatMessageListAdmin 迁移自 internal/logic/chat/message/chat_message_list_admin_logic.go，
// 和 ChatMessageList 逻辑相同（chatId==0 时查全部消息，管理页面用），gateway 侧两个不同的
// HTTP 路由都映射到这两个独立的 RPC 方法，和拆分前 admin.api 的 group 划分保持一致。
func (l *ChatMessageListAdminLogic) ChatMessageListAdmin(in *chat.ChatMessageListRequest) (*chat.ChatMessageListResponse, error) {
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
