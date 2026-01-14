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

	article, err := l.svcCtx.BlogArticleRepository.FindByID(l.ctx, req.Id)
	if err != nil {
		return nil, errs.Wrap(errs.CodeBadDB, "查询文章失败", err)
	}
	if article == nil || article.DeletedAt != 0 {
		return nil, errs.New(errs.CodeNotFound, "文章不存在")
	}

	// 仅待审核可审核
	if article.AuditStatus != consts.BlogArticleAuditStatusPending {
		return nil, errs.New(errs.CodeForbidden, "当前状态不允许审核")
	}

	// 写入审核记录
	u, _ := jwthelper.FromContext(l.ctx)
	auditRecord := &model.BlogArticleAudit{
		ArticleId:   article.Id,
		AuditStatus: req.Result,
		AuditRemark: req.Remark,
		AuditorId:   u.UserID,
		AuditorName: u.Username,
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
		DeletedAt:   0,
	}
	if err := l.svcCtx.BlogArticleAuditRepository.Create(l.ctx, auditRecord); err != nil {
		return nil, err
	}

	// 更新文章审核状态与业务状态
	article.AuditStatus = req.Result
	if req.Result == consts.BlogArticleAuditStatusPassed {
		article.Status = consts.BlogArticleStatusAuditPassed
	} else {
		// 驳回：文章业务状态保持草稿或置为下架？按方案：审核驳回
		// 这里保持 status 不变（仍可作为草稿编辑），但 audit_status=驳回
	}
	if err := l.svcCtx.Repository.BlogArticleModel.Update(l.ctx, article); err != nil {
		return nil, errs.Wrap(errs.CodeBadDB, "更新文章审核状态失败", err)
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
