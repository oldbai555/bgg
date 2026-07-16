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

type ChatGroupMemberAddLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChatGroupMemberAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatGroupMemberAddLogic {
	return &ChatGroupMemberAddLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// ChatGroupMemberAdd 薄胶水，实际业务逻辑已经搬进
// services/chat/internal/logic/chatgroupmemberaddlogic.go。
func (l *ChatGroupMemberAddLogic) ChatGroupMemberAdd(req *types.ChatGroupMemberAddReq) (resp *types.Response, err error) {
	if _, ok := jwthelper.FromContext(l.ctx); !ok {
		return nil, errs.New(errs.CodeUnauthorized, "未登录或登录已过期")
	}
	if len(req.UserIds) == 0 {
		return nil, errs.New(errs.CodeBadRequest, "用户ID列表不能为空")
	}

	_, err = l.svcCtx.ChatRPC.ChatGroupMemberAdd(l.ctx, &chatclient.ChatGroupMemberAddRequest{
		ChatId: req.ChatId, UserIds: req.UserIds,
	})
	if err != nil {
		return nil, errs.WrapGRPCError("添加成员失败", err)
	}

	return &types.Response{Code: 0, Message: "添加成员成功"}, nil
}
