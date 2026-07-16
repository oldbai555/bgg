// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package message

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/chat/chatclient"

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

// ChatMessageDelete 薄胶水，实际业务逻辑已经搬进
// services/chat/internal/logic/chatmessagedeletelogic.go。
func (l *ChatMessageDeleteLogic) ChatMessageDelete(req *types.ChatMessageDeleteReq) (resp *types.Response, err error) {
	if req == nil || req.Id == 0 {
		return nil, errs.New(errs.CodeBadRequest, "消息ID不能为空")
	}

	_, err = l.svcCtx.ChatRPC.ChatMessageDelete(l.ctx, &chatclient.ChatMessageDeleteRequest{Id: req.Id})
	if err != nil {
		return nil, errs.WrapGRPCError("删除消息失败", err)
	}

	return &types.Response{Code: 0, Message: "删除成功"}, nil
}
