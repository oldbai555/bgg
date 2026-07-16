// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package public

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/content/contentclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublicBlogFriendLinkListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPublicBlogFriendLinkListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublicBlogFriendLinkListLogic {
	return &PublicBlogFriendLinkListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// PublicBlogFriendLinkList 薄胶水，实际业务逻辑已经搬进
// services/content/internal/logic/publicblogfriendlinklistlogic.go。
func (l *PublicBlogFriendLinkListLogic) PublicBlogFriendLinkList() (resp *types.PublicBlogFriendLinkListResp, err error) {
	rpcResp, err := l.svcCtx.ContentRPC.PublicBlogFriendLinkList(l.ctx, &contentclient.PublicBlogGlobalRequest{})
	if err != nil {
		return nil, errs.WrapGRPCError("查询友情链接列表失败", err)
	}
	items := make([]types.PublicBlogFriendLinkItem, 0, len(rpcResp.List))
	for _, link := range rpcResp.List {
		items = append(items, types.PublicBlogFriendLinkItem{
			Id:       link.Id,
			Name:     link.Name,
			Url:      link.Url,
			Remark:   link.Remark,
			OrderNum: link.OrderNum,
		})
	}
	return &types.PublicBlogFriendLinkListResp{List: items}, nil
}
