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

type ChatGroupCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChatGroupCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatGroupCreateLogic {
	return &ChatGroupCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// ChatGroupCreate 薄胶水，建群组+拉创建人入群的事务、初始成员校验等实际业务逻辑已经搬进
// services/chat/internal/logic/chatgroupcreatelogic.go。
func (l *ChatGroupCreateLogic) ChatGroupCreate(req *types.ChatGroupCreateReq) (resp *types.Response, err error) {
	user, ok := jwthelper.FromContext(l.ctx)
	if !ok {
		return nil, errs.New(errs.CodeUnauthorized, "未登录或登录已过期")
	}

	_, err = l.svcCtx.ChatRPC.ChatGroupCreate(l.ctx, &chatclient.ChatGroupCreateRequest{
		Name:           req.Name,
		Avatar:         req.Avatar,
		Description:    req.Description,
		UserIds:        req.UserIds,
		OperatorUserId: user.UserID,
	})
	if err != nil {
		return nil, errs.WrapGRPCError("创建群组失败", err)
	}

	return &types.Response{Code: 0, Message: "创建群组成功"}, nil
}
