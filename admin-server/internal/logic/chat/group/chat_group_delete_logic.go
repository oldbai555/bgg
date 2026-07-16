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

type ChatGroupDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChatGroupDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatGroupDeleteLogic {
	return &ChatGroupDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// ChatGroupDelete 薄胶水，实际业务逻辑已经搬进
// services/chat/internal/logic/chatgroupdeletelogic.go。
func (l *ChatGroupDeleteLogic) ChatGroupDelete(req *types.ChatGroupDeleteReq) (resp *types.Response, err error) {
	if _, ok := jwthelper.FromContext(l.ctx); !ok {
		return nil, errs.New(errs.CodeUnauthorized, "未登录或登录已过期")
	}

	_, err = l.svcCtx.ChatRPC.ChatGroupDelete(l.ctx, &chatclient.ChatGroupDeleteRequest{Id: req.Id})
	if err != nil {
		return nil, errs.WrapGRPCError("删除群组失败", err)
	}

	return &types.Response{Code: 0, Message: "删除群组成功"}, nil
}
