package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/content/content"
	"postapocgame/admin-server/services/content/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type BlogArticleDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBlogArticleDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlogArticleDeleteLogic {
	return &BlogArticleDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// BlogArticleDelete 迁移自 internal/logic/blog/article/blog_article_delete_logic.go。
func (l *BlogArticleDeleteLogic) BlogArticleDelete(in *content.BlogArticleDeleteRequest) (*content.Empty, error) {
	if in.Id == 0 {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "文章ID不能为空"))
	}
	if err := l.svcCtx.BlogArticle.Delete(l.ctx, in.Id); err != nil {
		return nil, toGRPCStatus(err)
	}
	return &content.Empty{}, nil
}
