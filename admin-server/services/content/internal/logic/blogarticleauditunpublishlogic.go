package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/content/content"
	contentconsts "postapocgame/admin-server/services/content/internal/consts"
	"postapocgame/admin-server/services/content/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type BlogArticleAuditUnpublishLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBlogArticleAuditUnpublishLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlogArticleAuditUnpublishLogic {
	return &BlogArticleAuditUnpublishLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// BlogArticleAuditUnpublish 迁移自
// internal/logic/blog/article_audit/blog_article_audit_unpublish_logic.go。
func (l *BlogArticleAuditUnpublishLogic) BlogArticleAuditUnpublish(in *content.BlogArticleAuditUnpublishRequest) (*content.Empty, error) {
	if in.Id == 0 {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "文章ID不能为空"))
	}

	article, err := l.svcCtx.ArticleService.UnpublishArticle(l.ctx, in.Id, in.Remark, in.OperatorUserId, in.OperatorUsername)
	if err != nil {
		return nil, toGRPCStatus(err)
	}

	recordAuditLog(l.svcCtx, in.OperatorUserId, in.OperatorUsername, contentconsts.AuditTypeBlogArticleUnpublish, contentconsts.AuditObjectBlogArticle, map[string]any{
		"articleId": article.Id,
		"title":     article.Title,
		"remark":    in.Remark,
	})

	return &content.Empty{}, nil
}
