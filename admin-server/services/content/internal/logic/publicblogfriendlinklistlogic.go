package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/content/content"
	"postapocgame/admin-server/services/content/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublicBlogFriendLinkListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPublicBlogFriendLinkListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublicBlogFriendLinkListLogic {
	return &PublicBlogFriendLinkListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// PublicBlogFriendLinkList 迁移自
// internal/logic/blog/public/public_blog_friend_link_list_logic.go。
func (l *PublicBlogFriendLinkListLogic) PublicBlogFriendLinkList(in *content.PublicBlogGlobalRequest) (*content.PublicBlogFriendLinkListResponse, error) {
	list, err := l.svcCtx.FriendLink.FindEnabledList(l.ctx)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeBadDB, "查询友情链接列表失败", err))
	}

	items := make([]*content.PublicBlogFriendLinkItem, 0, len(list))
	for _, link := range list {
		items = append(items, &content.PublicBlogFriendLinkItem{
			Id:       link.Id,
			Name:     link.Name,
			Url:      link.Url,
			Remark:   link.Remark,
			OrderNum: link.OrderNum,
		})
	}

	return &content.PublicBlogFriendLinkListResponse{List: items}, nil
}
