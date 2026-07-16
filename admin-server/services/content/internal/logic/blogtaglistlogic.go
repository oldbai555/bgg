package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/content/content"
	"postapocgame/admin-server/services/content/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type BlogTagListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBlogTagListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlogTagListLogic {
	return &BlogTagListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// BlogTagList 迁移自 internal/logic/blog/tag/blog_tag_list_logic.go。
func (l *BlogTagListLogic) BlogTagList(in *content.BlogTagListRequest) (*content.BlogTagListResponse, error) {
	page := in.Page
	if page < 1 {
		page = 1
	}
	pageSize := in.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}

	list, total, err := l.svcCtx.Tag.FindPage(l.ctx, page, pageSize, in.Name, in.Status)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeBadDB, "查询标签列表失败", err))
	}

	items := make([]*content.BlogTagItem, 0, len(list))
	for _, t := range list {
		items = append(items, &content.BlogTagItem{
			Id:        t.Id,
			Name:      t.Name,
			Status:    t.Status,
			Remark:    t.Remark,
			CreatedAt: t.CreatedAt,
			UpdatedAt: t.UpdatedAt,
		})
	}

	return &content.BlogTagListResponse{Total: total, List: items}, nil
}
