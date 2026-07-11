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

type PublicBlogTagListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPublicBlogTagListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublicBlogTagListLogic {
	return &PublicBlogTagListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PublicBlogTagListLogic) PublicBlogTagList() (resp *types.PublicBlogTagListResp, err error) {
	// 查询启用的标签列表（使用较大的limit，实际业务中可根据需要调整）
	tagList, err := l.svcCtx.Domain.Blog.Tag.FindEnabledList(l.ctx, 1000)
	if err != nil {
		return nil, errs.Wrap(errs.CodeBadDB, "查询标签列表失败", err)
	}

	items := make([]types.BlogTagItem, 0, len(tagList))
	for _, tag := range tagList {
		items = append(items, types.BlogTagItem{
			Id:        tag.Id,
			Name:      tag.Name,
			Status:    tag.Status,
			Remark:    tag.Remark,
			CreatedAt: tag.CreatedAt,
			UpdatedAt: tag.UpdatedAt,
		})
	}

	return &types.PublicBlogTagListResp{
		List: items,
	}, nil
}
