package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/content/content"
	"postapocgame/admin-server/services/content/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type BlogArticleUntopLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBlogArticleUntopLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlogArticleUntopLogic {
	return &BlogArticleUntopLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// BlogArticleUntop 迁移自 internal/logic/blog/article/blog_article_untop_logic.go。
func (l *BlogArticleUntopLogic) BlogArticleUntop(in *content.BlogArticleUntopRequest) (*content.Empty, error) {
	if in.Id == 0 {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "文章ID不能为空"))
	}
	if err := l.svcCtx.BlogArticle.UpdateTopStatus(l.ctx, in.Id, 0); err != nil {
		return nil, toGRPCStatus(err)
	}
	return &content.Empty{}, nil
}
