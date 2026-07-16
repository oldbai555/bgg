// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package group

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	jwthelper "postapocgame/admin-server/pkg/jwt"
	"postapocgame/admin-server/services/chat/chatclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChatGroupMemberRemoveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChatGroupMemberRemoveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatGroupMemberRemoveLogic {
	return &ChatGroupMemberRemoveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// ChatGroupMemberRemove 薄胶水，实际业务逻辑已经搬进
// services/chat/internal/logic/chatgroupmemberremovelogic.go。
func (l *ChatGroupMemberRemoveLogic) ChatGroupMemberRemove(req *types.ChatGroupMemberRemoveReq) (resp *types.Response, err error) {
	if _, ok := jwthelper.FromContext(l.ctx); !ok {
		return nil, errs.New(errs.CodeUnauthorized, "未登录或登录已过期")
	}

	_, err = l.svcCtx.ChatRPC.ChatGroupMemberRemove(l.ctx, &chatclient.ChatGroupMemberRemoveRequest{
		ChatId: req.ChatId, UserId: req.UserId,
	})
	if err != nil {
		return nil, errs.WrapGRPCError("移除成员失败", err)
	}

	return &types.Response{Code: 0, Message: "移除成员成功"}, nil
}
