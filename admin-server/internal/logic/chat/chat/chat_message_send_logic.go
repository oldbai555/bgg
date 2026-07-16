// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package chat

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	jwthelper "postapocgame/admin-server/pkg/jwt"
	"postapocgame/admin-server/services/chat/chatclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChatMessageSendLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChatMessageSendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatMessageSendLogic {
	return &ChatMessageSendLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// ChatMessageSend 薄胶水：解析 HTTP 请求 -> 拼一次 ChatRPC 请求 -> 映射响应，持久化+广播的
// 实际业务逻辑已经搬进 services/chat/internal/logic/chatmessagesendlogic.go。
func (l *ChatMessageSendLogic) ChatMessageSend(req *types.ChatMessageSendReq) (resp *types.ChatMessageSendResp, err error) {
	if req == nil {
		return nil, errs.New(errs.CodeBadRequest, "请求参数不能为空")
	}

	user, ok := jwthelper.FromContext(l.ctx)
	if !ok {
		return nil, errs.New(errs.CodeUnauthorized, "未登录或登录已过期")
	}

	rpcResp, err := l.svcCtx.ChatRPC.ChatMessageSend(l.ctx, &chatclient.ChatMessageSendRequest{
		ChatId:           req.ChatId,
		Content:          req.Content,
		MessageType:      req.MessageType,
		OperatorUserId:   user.UserID,
		OperatorUsername: user.Username,
	})
	if err != nil {
		return nil, errs.WrapGRPCError("发送消息失败", err)
	}

	return &types.ChatMessageSendResp{Id: rpcResp.Id}, nil
}
