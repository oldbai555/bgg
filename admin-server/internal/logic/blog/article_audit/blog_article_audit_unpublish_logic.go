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

type BlogArticleAuditUnpublishLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBlogArticleAuditUnpublishLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlogArticleAuditUnpublishLogic {
	return &BlogArticleAuditUnpublishLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BlogArticleAuditUnpublishLogic) BlogArticleAuditUnpublish(req *types.BlogArticleAuditUnpublishReq) (resp *types.Response, err error) {
	if req.Id == 0 {
		return nil, errs.New(errs.CodeBadRequest, "文章ID不能为空")
	}

	u, _ := jwthelper.FromContext(l.ctx)
	article, err := l.svcCtx.Domain.Blog.ArticleService.UnpublishArticle(l.ctx, req.Id, req.Remark, u.UserID, u.Username)
	if err != nil {
		return nil, err
	}

	// 审计日志
	audit.RecordAuditLog(l.svcCtx, l.ctx, (&http.Request{Header: make(http.Header)}), consts.AuditTypeBlogArticleUnpublish, consts.AuditObjectBlogArticle, map[string]any{
		"articleId": article.Id,
		"title":     article.Title,
		"remark":    req.Remark,
	})

	return &types.Response{Code: int(errs.CodeOK), Message: "下架成功"}, nil
}
