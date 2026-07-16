// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package friend_link

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/content/contentclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type BlogFriendLinkListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBlogFriendLinkListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlogFriendLinkListLogic {
	return &BlogFriendLinkListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// BlogFriendLinkList 薄胶水，实际业务逻辑已经搬进
// services/content/internal/logic/blogfriendlinklistlogic.go。
func (l *BlogFriendLinkListLogic) BlogFriendLinkList(req *types.BlogFriendLinkListReq) (resp *types.BlogFriendLinkListResp, err error) {
	rpcResp, err := l.svcCtx.ContentRPC.BlogFriendLinkList(l.ctx, &contentclient.BlogFriendLinkListRequest{
		Page:    req.Page,
		Size:    req.Size,
		Status:  req.Status,
		Keyword: req.Keyword,
	})
	if err != nil {
		return nil, errs.WrapGRPCError("查询友情链接列表失败", err)
	}

	items := make([]types.BlogFriendLinkItem, 0, len(rpcResp.List))
	for _, link := range rpcResp.List {
		items = append(items, types.BlogFriendLinkItem{
			Id:        link.Id,
			Name:      link.Name,
			Url:       link.Url,
			Remark:    link.Remark,
			Status:    link.Status,
			OrderNum:  link.OrderNum,
			CreatedAt: link.CreatedAt,
			UpdatedAt: link.UpdatedAt,
		})
	}

	return &types.BlogFriendLinkListResp{Page: rpcResp.Page, Size: rpcResp.Size, Total: rpcResp.Total, List: items}, nil
}
