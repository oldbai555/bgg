// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package blog_article

import (
	"context"
	"postapocgame/admin-server/internal/dict"
	"strings"

	"postapocgame/admin-server/internal/consts"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"

	"github.com/zeromicro/go-zero/core/logx"
)

type BlogArticleUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBlogArticleUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlogArticleUpdateLogic {
	return &BlogArticleUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BlogArticleUpdateLogic) BlogArticleUpdate(req *types.BlogArticleUpdateReq) (resp *types.Response, err error) {
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

	// 编辑权限控制：上架状态（4）不可编辑
	if article.Status == consts.BlogArticleStatusPublished {
		return nil, errs.New(errs.CodeForbidden, "已上架的文章不可编辑，请先下架后再编辑")
	}

	// 记录原状态，用于决定编辑后的状态
	originalStatus := article.Status

	// 更新文章字段
	if req.Title != "" {
		title := strings.TrimSpace(req.Title)
		// 从字典读取文章标题最大长度限制（默认 100 个字符）
		maxTitleLength := dict.GetIntValue(l.ctx, l.svcCtx.Repository, consts.DictCodeBlogArticleTitleMaxLength, 100)
		if err := dict.ValidateLength(title, maxTitleLength, "文章标题"); err != nil {
			return nil, err
		}
		article.Title = title
	}
	if req.Content != "" {
		article.Content = req.Content
	}
	if req.Cover != "" {
		article.Cover = strings.TrimSpace(req.Cover)
	}
	if req.Summary != "" {
		article.Summary = strings.TrimSpace(req.Summary)
	}

	tagIDs := req.TagIds
	if len(tagIDs) == 0 {
		return nil, errs.New(errs.CodeBadRequest, "tagIds 不能为空")
	}

	// 编辑后状态流转规则：
	// - 草稿（1）-> 保存后仍为草稿（1），审核状态保持未提交（1）
	// - 待审核（2）-> 编辑后保存为待审核（2），审核状态重置为待审核（2）
	// - 审核通过-未上架（3）-> 编辑后保存为待审核（2），审核状态重置为待审核（2）
	// - 下架（5）-> 编辑后保存为待审核（2），审核状态重置为待审核（2）
	if originalStatus == consts.BlogArticleStatusDraft {
		// 草稿状态：保持草稿，审核状态保持未提交
		article.Status = consts.BlogArticleStatusDraft
		article.AuditStatus = consts.BlogArticleAuditStatusNotSubmitted
	} else {
		// 其他状态（待审核、审核通过-未上架、下架）：编辑后转为待审核
		article.Status = consts.BlogArticleStatusPendingAudit
		article.AuditStatus = consts.BlogArticleAuditStatusPending
	}

	if err = l.svcCtx.BlogArticleRepository.UpdateWithTags(l.ctx, article, tagIDs); err != nil {
		return nil, err
	}

	return &types.Response{Code: int(errs.CodeOK), Message: "更新成功"}, nil
}
