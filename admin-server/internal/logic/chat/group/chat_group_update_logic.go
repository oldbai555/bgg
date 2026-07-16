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

type ChatGroupUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChatGroupUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatGroupUpdateLogic {
	return &ChatGroupUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// ChatGroupUpdate 薄胶水，实际业务逻辑已经搬进
// services/chat/internal/logic/chatgroupupdatelogic.go。
func (l *ChatGroupUpdateLogic) ChatGroupUpdate(req *types.ChatGroupUpdateReq) (resp *types.Response, err error) {
	if _, ok := jwthelper.FromContext(l.ctx); !ok {
		return nil, errs.New(errs.CodeUnauthorized, "未登录或登录已过期")
	}

	_, err = l.svcCtx.ChatRPC.ChatGroupUpdate(l.ctx, &chatclient.ChatGroupUpdateRequest{
		Id: req.Id, Name: req.Name, Avatar: req.Avatar, Description: req.Description,
	})
	if err != nil {
		return nil, errs.WrapGRPCError("更新群组失败", err)
	}

	return &types.Response{Code: 0, Message: "更新群组成功"}, nil
}
