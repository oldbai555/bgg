package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/content/content"
	"postapocgame/admin-server/services/content/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublicBlogArticleNextLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPublicBlogArticleNextLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublicBlogArticleNextLogic {
	return &PublicBlogArticleNextLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// PublicBlogArticleNext 迁移自 internal/logic/blog/public/public_blog_article_next_logic.go。
func (l *PublicBlogArticleNextLogic) PublicBlogArticleNext(in *content.PublicBlogArticlePrevNextRequest) (*content.PublicBlogArticlePrevNextResponse, error) {
	if in.Id == 0 {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "文章ID不能为空"))
	}

	currentArticle, err := l.svcCtx.BlogArticle.FindByID(l.ctx, in.Id)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeBadDB, "查询当前文章失败", err))
	}
	if currentArticle == nil {
		return nil, toGRPCStatus(errs.New(errs.CodeNotFound, "文章不存在"))
	}

	nextArticle, err := l.svcCtx.BlogArticle.FindNextArticle(l.ctx, currentArticle.PublishTime)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeBadDB, "查询下一篇文章失败", err))
	}
	if nextArticle == nil {
		return &content.PublicBlogArticlePrevNextResponse{}, nil
	}

	return &content.PublicBlogArticlePrevNextResponse{
		Id:          nextArticle.Id,
		Title:       nextArticle.Title,
		PublishTime: nextArticle.PublishTime,
	}, nil
}
