package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/content/content"
	"postapocgame/admin-server/services/content/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type BlogTagOptionsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBlogTagOptionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlogTagOptionsLogic {
	return &BlogTagOptionsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// BlogTagOptions 迁移自 internal/logic/blog/tag/blog_tag_options_logic.go。
func (l *BlogTagOptionsLogic) BlogTagOptions(in *content.BlogTagOptionsRequest) (*content.BlogTagOptionsResponse, error) {
	limit := int64(1000)
	if in != nil && in.Limit > 0 {
		limit = in.Limit
	}
	list, err := l.svcCtx.Tag.FindEnabledList(l.ctx, limit)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeBadDB, "查询标签选项失败", err))
	}
	items := make([]*content.BlogTagOptionItem, 0, len(list))
	for _, t := range list {
		items = append(items, &content.BlogTagOptionItem{Id: t.Id, Name: t.Name})
	}
	return &content.BlogTagOptionsResponse{List: items}, nil
}
