// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package blog_article_audit

import (
	"context"
	"net/http"
	"time"

	"postapocgame/admin-server/internal/consts"
	"postapocgame/admin-server/internal/model"
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

	article, err := l.svcCtx.BlogArticleRepository.FindByID(l.ctx, req.Id)
	if err != nil {
		return nil, errs.Wrap(errs.CodeBadDB, "查询文章失败", err)
	}
	if article == nil || article.DeletedAt != 0 {
		return nil, errs.New(errs.CodeNotFound, "文章不存在")
	}

	if article.Status != consts.BlogArticleStatusPublished {
		return nil, errs.New(errs.CodeForbidden, "仅已上架文章可下架")
	}

	article.Status = consts.BlogArticleStatusUnpublished
	if err := l.svcCtx.Repository.BlogArticleModel.Update(l.ctx, article); err != nil {
		return nil, errs.Wrap(errs.CodeBadDB, "下架失败", err)
	}

	// 记录一条审核记录（便于追踪谁下架）
	u, _ := jwthelper.FromContext(l.ctx)
	auditRecord := &model.BlogArticleAudit{
		ArticleId:   article.Id,
		AuditStatus: consts.BlogArticleAuditStatusRejected, // 复用字段，语义为“审核操作类记录”
		AuditRemark: req.Remark,
		AuditorId:   u.UserID,
		AuditorName: u.Username,
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
		DeletedAt:   0,
	}
	_ = l.svcCtx.BlogArticleAuditRepository.Create(l.ctx, auditRecord)

	// 审计日志
	audit.RecordAuditLog(l.svcCtx, l.ctx, (&http.Request{Header: make(http.Header)}), consts.AuditTypeBlogArticleUnpublish, consts.AuditObjectBlogArticle, map[string]any{
		"articleId": article.Id,
		"title":     article.Title,
		"remark":    req.Remark,
	})

	return &types.Response{Code: int(errs.CodeOK), Message: "下架成功"}, nil
}
