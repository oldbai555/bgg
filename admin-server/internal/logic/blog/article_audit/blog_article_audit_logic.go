// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package article_audit

import (
	"context"
	"net/http"

	"postapocgame/admin-server/internal/consts"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/audit"
	"postapocgame/admin-server/pkg/errs"
	jwthelper "postapocgame/admin-server/pkg/jwt"

	"github.com/zeromicro/go-zero/core/logx"
)

type BlogArticleAuditLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBlogArticleAuditLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlogArticleAuditLogic {
	return &BlogArticleAuditLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BlogArticleAuditLogic) BlogArticleAudit(req *types.BlogArticleAuditReq) (resp *types.Response, err error) {
	if req.Id == 0 {
		return nil, errs.New(errs.CodeBadRequest, "文章ID不能为空")
	}
	if req.Result != consts.BlogArticleAuditStatusPassed && req.Result != consts.BlogArticleAuditStatusRejected {
		return nil, errs.New(errs.CodeBadRequest, "审核结果不合法")
	}

	u, _ := jwthelper.FromContext(l.ctx)
	article, err := l.svcCtx.Domain.Blog.ArticleService.AuditArticle(l.ctx, req.Id, req.Result, req.Remark, u.UserID, u.Username)
	if err != nil {
		return nil, err
	}

	// 记录审计日志（audit_type 使用字典 value）
	audit.RecordAuditLog(l.svcCtx, l.ctx, (&http.Request{Header: make(http.Header)}), consts.AuditTypeBlogArticleAudit, consts.AuditObjectBlogArticle, map[string]any{
		"articleId":     article.Id,
		"title":         article.Title,
		"toAuditStatus": req.Result,
		"remark":        req.Remark,
	})

	return &types.Response{Code: int(errs.CodeOK), Message: "审核成功"}, nil
}
