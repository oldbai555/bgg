package logic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/chat/chat"
	chatconsts "postapocgame/admin-server/services/chat/internal/consts"
	"postapocgame/admin-server/services/chat/internal/svc"
)

type ChatGroupMemberListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewChatGroupMemberListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatGroupMemberListLogic {
	return &ChatGroupMemberListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ChatGroupMemberList 迁移自 internal/logic/chat/group/chat_group_member_list_logic.go，
// 成员信息解析复用 resolveGroupMembers（见 chatgroupdetaillogic.go）。
func (l *ChatGroupMemberListLogic) ChatGroupMemberList(in *chat.ChatGroupMemberListRequest) (*chat.ChatGroupMemberListResponse, error) {
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

	chatUsers, err := l.svcCtx.ChatUser.FindByChatID(l.ctx, c.Id)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "查询群组成员失败", err))
	}

	return &chat.ChatGroupMemberListResponse{List: resolveGroupMembers(l.ctx, l.svcCtx, chatUsers)}, nil
}
