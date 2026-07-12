package logic

import (
	"context"
	"strings"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/content/content"
	contentconsts "postapocgame/admin-server/services/content/internal/consts"
	"postapocgame/admin-server/services/content/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type BlogArticleUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBlogArticleUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlogArticleUpdateLogic {
	return &BlogArticleUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// BlogArticleUpdate 迁移自 internal/logic/blog/article/blog_article_update_logic.go。
func (l *BlogArticleUpdateLogic) BlogArticleUpdate(in *content.BlogArticleUpdateRequest) (*content.Empty, error) {
	if in.Id == 0 {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "文章ID不能为空"))
	}

	article, err := l.svcCtx.BlogArticle.FindByID(l.ctx, in.Id)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeBadDB, "查询文章失败", err))
	}
	if article == nil || article.DeletedAt != 0 {
		return nil, toGRPCStatus(errs.New(errs.CodeNotFound, "文章不存在"))
	}

	if article.Status == contentconsts.BlogArticleStatusPublished {
		return nil, toGRPCStatus(errs.New(errs.CodeForbidden, "已上架的文章不可编辑，请先下架后再编辑"))
	}

	originalStatus := article.Status

	if in.Title != "" {
		title := strings.TrimSpace(in.Title)
		if err := validateLength(title, l.svcCtx.Config.Limits.BlogArticleTitleMaxLength, "文章标题"); err != nil {
			return nil, toGRPCStatus(err)
		}
		article.Title = title
	}
	if in.Content != "" {
		article.Content = in.Content
	}
	if in.Cover != "" {
		article.Cover = strings.TrimSpace(in.Cover)
	}
	if in.Summary != "" {
		article.Summary = strings.TrimSpace(in.Summary)
	}

	tagIDs := in.TagIds
	if len(tagIDs) == 0 {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "tagIds 不能为空"))
	}

	// 编辑后状态流转规则：
	// - 草稿（1）-> 保存后仍为草稿（1），审核状态保持未提交（1）
	// - 待审核（2）/审核通过-未上架（3）/下架（5）-> 编辑后保存为待审核（2），审核状态重置为待审核（2）
	if originalStatus == contentconsts.BlogArticleStatusDraft {
		article.Status = contentconsts.BlogArticleStatusDraft
		article.AuditStatus = contentconsts.BlogArticleAuditStatusNotSubmitted
	} else {
		article.Status = contentconsts.BlogArticleStatusPendingAudit
		article.AuditStatus = contentconsts.BlogArticleAuditStatusPending
	}

	if err := l.svcCtx.ArticleService.UpdateArticle(l.ctx, article, tagIDs); err != nil {
		return nil, toGRPCStatus(err)
	}

	return &content.Empty{}, nil
}
