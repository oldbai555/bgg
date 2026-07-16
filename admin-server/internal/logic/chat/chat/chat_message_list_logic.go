// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package chat

import (
	"context"
	"postapocgame/admin-server/internal/logic/logicutil"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	jwthelper "postapocgame/admin-server/pkg/jwt"
	"postapocgame/admin-server/services/chat/chatclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChatMessageListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChatMessageListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatMessageListLogic {
	return &ChatMessageListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// ChatMessageList 薄胶水：解析 HTTP 请求 -> 拼一次 ChatRPC 请求 -> 映射响应，实际业务逻辑
// 已经搬进 services/chat/internal/logic/chatmessagelistlogic.go。
func (l *ChatMessageListLogic) ChatMessageList(req *types.ChatMessageListReq) (resp *types.ChatMessageListResp, err error) {
	if req == nil {
		return nil, errs.New(errs.CodeBadRequest, "请求参数不能为空")
	}
	req.Page, req.PageSize = logicutil.NormalizePage(req.Page, req.PageSize, 20, 100)

	if _, ok := jwthelper.FromContext(l.ctx); !ok {
		return nil, errs.New(errs.CodeUnauthorized, "未登录或登录已过期")
	}

	rpcResp, err := l.svcCtx.ChatRPC.ChatMessageList(l.ctx, &chatclient.ChatMessageListRequest{
		Page: req.Page, PageSize: req.PageSize, ChatId: req.ChatId,
	})
	if err != nil {
		return nil, errs.WrapGRPCError("查询聊天消息列表失败", err)
	}

	items := make([]types.ChatMessageItem, 0, len(rpcResp.List))
	for _, msg := range rpcResp.List {
		items = append(items, types.ChatMessageItem{
			Id:           msg.Id,
			ChatId:       msg.ChatId,
			FromUserId:   msg.FromUserId,
			FromUserName: msg.FromUserName,
			Content:      msg.Content,
			MessageType:  msg.MessageType,
			Status:       msg.Status,
			CreatedAt:    msg.CreatedAt,
		})
	}

	return &types.ChatMessageListResp{Total: rpcResp.Total, List: items}, nil
}
