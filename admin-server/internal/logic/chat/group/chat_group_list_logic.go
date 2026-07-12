// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package group

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

type ChatGroupListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChatGroupListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatGroupListLogic {
	return &ChatGroupListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// ChatGroupList 薄胶水，实际业务逻辑已经搬进
// services/chat/internal/logic/chatgrouplistlogic.go。
func (l *ChatGroupListLogic) ChatGroupList(req *types.ChatGroupListReq) (resp *types.ChatGroupListResp, err error) {
	if _, ok := jwthelper.FromContext(l.ctx); !ok {
		return nil, errs.New(errs.CodeUnauthorized, "未登录或登录已过期")
	}
	req.Page, req.PageSize = logicutil.NormalizePage(req.Page, req.PageSize, 10, 100)

	rpcResp, err := l.svcCtx.ChatRPC.ChatGroupList(l.ctx, &chatclient.ChatGroupListRequest{
		Page: req.Page, PageSize: req.PageSize, Name: req.Name,
	})
	if err != nil {
		return nil, errs.WrapGRPCError("查询群组列表失败", err)
	}

	items := make([]types.ChatGroupItem, 0, len(rpcResp.List))
	for _, g := range rpcResp.List {
		items = append(items, types.ChatGroupItem{
			Id:          g.Id,
			Name:        g.Name,
			Avatar:      g.Avatar,
			Description: g.Description,
			CreatedBy:   g.CreatedBy,
			CreatedAt:   g.CreatedAt,
			MemberCount: g.MemberCount,
		})
	}

	return &types.ChatGroupListResp{Total: rpcResp.Total, List: items}, nil
}
