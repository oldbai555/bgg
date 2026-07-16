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

type ChatListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChatListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatListLogic {
	return &ChatListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// ChatList 薄胶水：解析 HTTP 请求 -> 拼一次 ChatRPC 请求 -> 映射响应，chat 域的实际业务
// 逻辑（私聊对方用户信息回调 IamCallback 展开）已经搬进
// services/chat/internal/logic/chatlistlogic.go。
func (l *ChatListLogic) ChatList(req *types.ChatListReq) (resp *types.ChatListResp, err error) {
	user, ok := jwthelper.FromContext(l.ctx)
	if !ok {
		return nil, errs.New(errs.CodeUnauthorized, "未登录或登录已过期")
	}

	rpcResp, err := l.svcCtx.ChatRPC.ChatList(l.ctx, &chatclient.ChatListRequest{OperatorUserId: user.UserID})
	if err != nil {
		return nil, errs.WrapGRPCError("查询聊天列表失败", err)
	}

	items := make([]types.ChatItem, 0, len(rpcResp.List))
	for _, c := range rpcResp.List {
		items = append(items, types.ChatItem{
			ChatId:         c.ChatId,
			Name:           c.Name,
			ChatType:       c.ChatType,
			Avatar:         c.Avatar,
			Description:    c.Description,
			UserId:         c.UserId,
			Username:       c.Username,
			Nickname:       c.Nickname,
			DepartmentName: c.DepartmentName,
			RoleNames:      c.RoleNames,
			UnreadCount:    c.UnreadCount,
			LastMessage:    c.LastMessage,
			LastMessageAt:  c.LastMessageAt,
		})
	}

	return &types.ChatListResp{List: items}, nil
}
