package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/content/content"
	"postapocgame/admin-server/services/content/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublicBlogTagListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPublicBlogTagListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublicBlogTagListLogic {
	return &PublicBlogTagListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// PublicBlogTagList 迁移自 internal/logic/blog/public/public_blog_tag_list_logic.go。
func (l *PublicBlogTagListLogic) PublicBlogTagList(in *content.PublicBlogGlobalRequest) (*content.PublicBlogTagListResponse, error) {
	list, err := l.svcCtx.Tag.FindEnabledList(l.ctx, 1000)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeBadDB, "查询标签列表失败", err))
	}

	items := make([]*content.BlogTagItem, 0, len(list))
	for _, tag := range list {
		items = append(items, &content.BlogTagItem{
			Id:        tag.Id,
			Name:      tag.Name,
			Status:    tag.Status,
			Remark:    tag.Remark,
			CreatedAt: tag.CreatedAt,
			UpdatedAt: tag.UpdatedAt,
		})
	}

	return &content.PublicBlogTagListResponse{List: items}, nil
}
