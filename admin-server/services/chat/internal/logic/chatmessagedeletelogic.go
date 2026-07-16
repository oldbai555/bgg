package logic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/chat/chat"
	"postapocgame/admin-server/services/chat/internal/svc"
)

type ChatMessageDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewChatMessageDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatMessageDeleteLogic {
	return &ChatMessageDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ChatMessageDelete 迁移自 internal/logic/chat/message/chat_message_delete_logic.go。
func (l *ChatMessageDeleteLogic) ChatMessageDelete(in *chat.ChatMessageDeleteRequest) (*chat.Empty, error) {
	if in.Id == 0 {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "消息ID不能为空"))
	}

	message, err := l.svcCtx.ChatMessage.FindByID(l.ctx, in.Id)
	if err != nil || message == nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeNotFound, "消息不存在", err))
	}

	if err := l.svcCtx.ChatMessage.DeleteByID(l.ctx, in.Id); err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "删除消息失败", err))
	}

	return &chat.Empty{}, nil
}
