package logic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/chat/chat"
	chatconsts "postapocgame/admin-server/services/chat/internal/consts"
	"postapocgame/admin-server/services/chat/internal/svc"
)

type ChatGroupDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewChatGroupDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatGroupDeleteLogic {
	return &ChatGroupDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ChatGroupDelete 迁移自 internal/logic/chat/group/chat_group_delete_logic.go。
func (l *ChatGroupDeleteLogic) ChatGroupDelete(in *chat.ChatGroupDeleteRequest) (*chat.Empty, error) {
	c, err := l.svcCtx.Chat.FindByID(l.ctx, in.Id)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeNotFound, "群组不存在", err))
	}
	if c.Type != chatconsts.ChatTypeGroup {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "该聊天不是群组"))
	}
	if c.DeletedAt != 0 {
		return nil, toGRPCStatus(errs.New(errs.CodeNotFound, "群组已删除"))
	}

	if err := l.svcCtx.Chat.DeleteByID(l.ctx, in.Id); err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "删除群组失败", err))
	}

	return &chat.Empty{}, nil
}
