// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package public

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"

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

func (l *PublicBlogFriendLinkListLogic) PublicBlogFriendLinkList() (resp *types.PublicBlogFriendLinkListResp, err error) {
	// 查询启用的友情链接列表
	list, err := l.svcCtx.Domain.Blog.FriendLink.FindEnabledList(l.ctx)
	if err != nil {
		return nil, errs.Wrap(errs.CodeBadDB, "查询友情链接列表失败", err)
	}

	items := make([]types.PublicBlogFriendLinkItem, 0, len(list))
	for _, link := range list {
		items = append(items, types.PublicBlogFriendLinkItem{
			Id:       link.Id,
			Name:     link.Name,
			Url:      link.Url,
			Remark:   link.Remark,
			OrderNum: link.OrderNum,
		})
	}

	return &types.PublicBlogFriendLinkListResp{
		List: items,
	}, nil
}
