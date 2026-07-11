// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package message

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChatMessageDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChatMessageDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatMessageDeleteLogic {
	return &ChatMessageDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChatMessageDeleteLogic) ChatMessageDelete(req *types.ChatMessageDeleteReq) (resp *types.Response, err error) {
	if req == nil || req.Id == 0 {
		return nil, errs.New(errs.CodeBadRequest, "消息ID不能为空")
	}

	// 检查消息是否存在
	message, err := l.svcCtx.Domain.Chat.ChatMessage.FindByID(l.ctx, req.Id)
	if err != nil {
		return nil, errs.Wrap(errs.CodeNotFound, "消息不存在", err)
	}
	if message == nil {
		return nil, errs.New(errs.CodeNotFound, "消息不存在")
	}

	// 删除消息（软删除）
	if err := l.svcCtx.Domain.Chat.ChatMessage.DeleteByID(l.ctx, req.Id); err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, "删除消息失败", err)
	}

	return &types.Response{
		Code:    0,
		Message: "删除成功",
	}, nil
}
