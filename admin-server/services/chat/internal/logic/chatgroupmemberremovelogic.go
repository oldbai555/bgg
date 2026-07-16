package logic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/chat/chat"
	chatconsts "postapocgame/admin-server/services/chat/internal/consts"
	"postapocgame/admin-server/services/chat/internal/svc"
)

type ChatGroupMemberRemoveLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewChatGroupMemberRemoveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatGroupMemberRemoveLogic {
	return &ChatGroupMemberRemoveLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ChatGroupMemberRemove 迁移自 internal/logic/chat/group/chat_group_member_remove_logic.go。
func (l *ChatGroupMemberRemoveLogic) ChatGroupMemberRemove(in *chat.ChatGroupMemberRemoveRequest) (*chat.Empty, error) {
	c, err := l.svcCtx.Chat.FindByID(l.ctx, in.ChatId)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeNotFound, "群组不存在", err))
	}
	if c.Type != chatconsts.ChatTypeGroup {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "该聊天不是群组"))
	}
	if c.DeletedAt != 0 {
		return nil, toGRPCStatus(errs.New(errs.CodeNotFound, "群组已删除"))
	}

	chatUsers, err := l.svcCtx.ChatUser.FindByChatID(l.ctx, in.ChatId)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "查询群组成员失败", err))
	}
	userInGroup := false
	for _, cu := range chatUsers {
		if cu.UserId == in.UserId {
			userInGroup = true
			break
		}
	}
	if !userInGroup {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "用户不在群组中"))
	}

	if err := l.svcCtx.ChatUser.DeleteByChatIDAndUserID(l.ctx, in.ChatId, in.UserId); err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "移除成员失败", err))
	}

	return &chat.Empty{}, nil
}
