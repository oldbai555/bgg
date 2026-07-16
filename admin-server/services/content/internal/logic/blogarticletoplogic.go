package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/content/content"
	"postapocgame/admin-server/services/content/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type BlogArticleTopLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBlogArticleTopLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlogArticleTopLogic {
	return &BlogArticleTopLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// BlogArticleTop 迁移自 internal/logic/blog/article/blog_article_top_logic.go。置顶数量上限
// 原来读字典 blog_article_top_max_count（物理属于 iam 域），改成
// svcCtx.Config.Limits.BlogArticleTopMaxCount 静态配置。
func (l *BlogArticleTopLogic) BlogArticleTop(in *content.BlogArticleTopRequest) (*content.Empty, error) {
	if in.Id == 0 {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "文章ID不能为空"))
	}

	maxCount := l.svcCtx.Config.Limits.BlogArticleTopMaxCount
	if err := l.svcCtx.ArticleService.SetArticleTop(l.ctx, in.Id, maxCount); err != nil {
		return nil, toGRPCStatus(err)
	}

	return &content.Empty{}, nil
}
