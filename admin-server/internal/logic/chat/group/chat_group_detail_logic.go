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

type ChatGroupDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChatGroupDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatGroupDetailLogic {
	return &ChatGroupDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// ChatGroupDetail 薄胶水，成员信息（用户名/昵称/部门名/角色名）解析已经搬进
// services/chat/internal/logic/chatgroupdetaillogic.go（回调 IamCallback）。
func (l *ChatGroupDetailLogic) ChatGroupDetail(req *types.ChatGroupDetailReq) (resp *types.ChatGroupDetailResp, err error) {
	if _, ok := jwthelper.FromContext(l.ctx); !ok {
		return nil, errs.New(errs.CodeUnauthorized, "未登录或登录已过期")
	}

	rpcResp, err := l.svcCtx.ChatRPC.ChatGroupDetail(l.ctx, &chatclient.ChatGroupDetailRequest{Id: req.Id})
	if err != nil {
		return nil, errs.WrapGRPCError("查询群组详情失败", err)
	}

	members := make([]types.ChatGroupMemberItem, 0, len(rpcResp.Members))
	for _, m := range rpcResp.Members {
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

	return &types.ChatGroupDetailResp{
		Id:          rpcResp.Id,
		Name:        rpcResp.Name,
		Avatar:      rpcResp.Avatar,
		Description: rpcResp.Description,
		CreatedBy:   rpcResp.CreatedBy,
		CreatedAt:   rpcResp.CreatedAt,
		MemberCount: rpcResp.MemberCount,
		Members:     members,
	}, nil
}
