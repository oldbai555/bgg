// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package blog_friend_link

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"

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

func (l *BlogFriendLinkListLogic) BlogFriendLinkList(req *types.BlogFriendLinkListReq) (resp *types.BlogFriendLinkListResp, err error) {
	// 参数预处理与默认值
	page := req.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.Size
	if pageSize <= 0 {
		pageSize = 20
	}

	// 调用仓储层分页查询
	list, total, err := l.svcCtx.BlogFriendLinkRepository.FindPage(l.ctx, page, pageSize, req.Status, req.Keyword)
	if err != nil {
		return nil, errs.Wrap(errs.CodeBadDB, "查询友情链接列表失败", err)
	}

	items := make([]types.BlogFriendLinkItem, 0, len(list))
	for _, link := range list {
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

	return &types.BlogFriendLinkListResp{
		Page:  page,
		Size:  pageSize,
		Total: total,
		List:  items,
	}, nil
}
