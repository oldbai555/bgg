package logic

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/chat/chat"
	chatconsts "postapocgame/admin-server/services/chat/internal/consts"
	"postapocgame/admin-server/services/chat/internal/svc"
)

type ChatGroupUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewChatGroupUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatGroupUpdateLogic {
	return &ChatGroupUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ChatGroupUpdate 迁移自 internal/logic/chat/group/chat_group_update_logic.go。
func (l *ChatGroupUpdateLogic) ChatGroupUpdate(in *chat.ChatGroupUpdateRequest) (*chat.Empty, error) {
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

	updated := false
	if in.Name != "" && in.Name != c.Name {
		c.Name = in.Name
		updated = true
	}
	if in.Avatar != "" && in.Avatar != c.Avatar {
		c.Avatar = in.Avatar
		updated = true
	}
	if in.Description != "" && in.Description != c.Description {
		c.Description = in.Description
		updated = true
	}
	if !updated {
		return &chat.Empty{}, nil
	}

	c.UpdatedAt = time.Now().Unix()
	if err := l.svcCtx.Chat.Update(l.ctx, c); err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "更新群组失败", err))
	}

	return &chat.Empty{}, nil
}
