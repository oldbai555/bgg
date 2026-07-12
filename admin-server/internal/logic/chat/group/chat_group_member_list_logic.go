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

type ChatGroupMemberListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChatGroupMemberListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatGroupMemberListLogic {
	return &ChatGroupMemberListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// ChatGroupMemberList 薄胶水，实际业务逻辑已经搬进
// services/chat/internal/logic/chatgroupmemberlistlogic.go。
func (l *ChatGroupMemberListLogic) ChatGroupMemberList(req *types.ChatGroupMemberListReq) (resp *types.ChatGroupMemberListResp, err error) {
	if _, ok := jwthelper.FromContext(l.ctx); !ok {
		return nil, errs.New(errs.CodeUnauthorized, "未登录或登录已过期")
	}

	rpcResp, err := l.svcCtx.ChatRPC.ChatGroupMemberList(l.ctx, &chatclient.ChatGroupMemberListRequest{Id: req.Id})
	if err != nil {
		return nil, errs.WrapGRPCError("查询群组成员失败", err)
	}

	members := make([]types.ChatGroupMemberItem, 0, len(rpcResp.List))
	for _, m := range rpcResp.List {
		members = append(members, types.ChatGroupMemberItem{
			UserId:         m.UserId,
			Username:       m.Username,
			Nickname:       m.Nickname,
			Avatar:         m.Avatar,
			DepartmentName: m.DepartmentName,
			RoleNames:      m.RoleNames,
			JoinedAt:       m.JoinedAt,
		})
	}

	return &types.ChatGroupMemberListResp{List: members}, nil
}
