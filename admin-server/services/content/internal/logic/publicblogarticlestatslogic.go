package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/content/content"
	"postapocgame/admin-server/services/content/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublicBlogArticleStatsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPublicBlogArticleStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublicBlogArticleStatsLogic {
	return &PublicBlogArticleStatsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// PublicBlogArticleStats 迁移自 internal/logic/blog/public/public_blog_article_stats_logic.go。
func (l *PublicBlogArticleStatsLogic) PublicBlogArticleStats(in *content.PublicBlogGlobalRequest) (*content.PublicBlogArticleStatsResponse, error) {
	total, err := l.svcCtx.BlogArticle.CountPublishedArticles(l.ctx)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeBadDB, "统计文章总数失败", err))
	}
	return &content.PublicBlogArticleStatsResponse{TotalArticles: total}, nil
}
