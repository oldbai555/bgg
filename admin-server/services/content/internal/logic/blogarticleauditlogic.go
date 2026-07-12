package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/content/content"
	contentconsts "postapocgame/admin-server/services/content/internal/consts"
	"postapocgame/admin-server/services/content/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type BlogArticleAuditLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBlogArticleAuditLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlogArticleAuditLogic {
	return &BlogArticleAuditLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 审核

// BlogArticleAudit 迁移自 internal/logic/blog/article_audit/blog_article_audit_logic.go。
// 原实现里的审计日志（pkg/audit.RecordAuditLog）改成回调 IamCallback.RecordAuditLog，
// 见 services/content/internal/logic/audit.go。
func (l *BlogArticleAuditLogic) BlogArticleAudit(in *content.BlogArticleAuditRequest) (*content.Empty, error) {
	if in.Id == 0 {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "文章ID不能为空"))
	}
	if in.Result != contentconsts.BlogArticleAuditStatusPassed && in.Result != contentconsts.BlogArticleAuditStatusRejected {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "审核结果不合法"))
	}

	article, err := l.svcCtx.ArticleService.AuditArticle(l.ctx, in.Id, in.Result, in.Remark, in.OperatorUserId, in.OperatorUsername)
	if err != nil {
		return nil, toGRPCStatus(err)
	}

	recordAuditLog(l.svcCtx, in.OperatorUserId, in.OperatorUsername, contentconsts.AuditTypeBlogArticleAudit, contentconsts.AuditObjectBlogArticle, map[string]any{
		"articleId":     article.Id,
		"title":         article.Title,
		"toAuditStatus": in.Result,
		"remark":        in.Remark,
	})

	return &content.Empty{}, nil
}
