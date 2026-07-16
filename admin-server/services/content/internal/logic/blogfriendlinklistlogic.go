package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/content/content"
	"postapocgame/admin-server/services/content/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type BlogFriendLinkListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBlogFriendLinkListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlogFriendLinkListLogic {
	return &BlogFriendLinkListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 友情链接

// BlogFriendLinkList 迁移自 internal/logic/blog/friend_link/blog_friend_link_list_logic.go。
func (l *BlogFriendLinkListLogic) BlogFriendLinkList(in *content.BlogFriendLinkListRequest) (*content.BlogFriendLinkListResponse, error) {
	page := in.Page
	if page < 1 {
		page = 1
	}
	pageSize := in.Size
	if pageSize <= 0 {
		pageSize = 20
	}

	list, total, err := l.svcCtx.FriendLink.FindPage(l.ctx, page, pageSize, in.Status, in.Keyword)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeBadDB, "查询友情链接列表失败", err))
	}

	items := make([]*content.BlogFriendLinkItem, 0, len(list))
	for _, link := range list {
		items = append(items, &content.BlogFriendLinkItem{
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

	return &content.BlogFriendLinkListResponse{Page: page, Size: pageSize, Total: total, List: items}, nil
}
