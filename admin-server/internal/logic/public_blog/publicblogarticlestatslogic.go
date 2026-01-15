// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package public_blog

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublicBlogArticleStatsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPublicBlogArticleStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublicBlogArticleStatsLogic {
	return &PublicBlogArticleStatsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PublicBlogArticleStatsLogic) PublicBlogArticleStats() (resp *types.PublicBlogArticleStatsResp, err error) {
	// 统计已发布文章总数
	total, err := l.svcCtx.BlogArticleRepository.CountPublishedArticles(l.ctx)
	if err != nil {
		return nil, errs.Wrap(errs.CodeBadDB, "统计文章总数失败", err)
	}

	return &types.PublicBlogArticleStatsResp{
		TotalArticles: total,
	}, nil
}
