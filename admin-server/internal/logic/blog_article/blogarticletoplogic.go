// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package blog_article

import (
	"context"

	"postapocgame/admin-server/internal/consts"
	"postapocgame/admin-server/internal/dict"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"

	"github.com/zeromicro/go-zero/core/logx"
)

type BlogArticleTopLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBlogArticleTopLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlogArticleTopLogic {
	return &BlogArticleTopLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BlogArticleTopLogic) BlogArticleTop(req *types.BlogArticleTopReq) (resp *types.Response, err error) {
	if req.Id == 0 {
		return nil, errs.New(errs.CodeBadRequest, "文章ID不能为空")
	}

	// 查询文章是否存在
	article, err := l.svcCtx.BlogArticleRepository.FindByID(l.ctx, req.Id)
	if err != nil {
		return nil, errs.Wrap(errs.CodeBadDB, "查询文章失败", err)
	}
	if article == nil {
		return nil, errs.New(errs.CodeNotFound, "文章不存在")
	}

	// 如果已经是置顶状态，直接返回成功
	if article.IsTop == 1 {
		return &types.Response{
			Code:    0,
			Message: "文章已置顶",
		}, nil
	}

	// 从字典读取置顶最大数量限制（默认1篇）
	maxCount := int64(dict.GetIntValue(l.ctx, l.svcCtx.Repository, consts.DictCodeBlogArticleTopMaxCount, 1))

	// 查询当前置顶数量
	currentTopCount, err := l.svcCtx.BlogArticleRepository.FindTopCount(l.ctx)
	if err != nil {
		return nil, errs.Wrap(errs.CodeBadDB, "查询置顶文章数量失败", err)
	}

	// 如果已达到最大数量，取消最早置顶的文章
	if currentTopCount >= maxCount {
		oldestTopArticle, err := l.svcCtx.BlogArticleRepository.FindOldestTopArticle(l.ctx)
		if err != nil {
			return nil, errs.Wrap(errs.CodeBadDB, "查询最早置顶文章失败", err)
		}
		if oldestTopArticle != nil {
			// 取消最早置顶的文章
			if err := l.svcCtx.BlogArticleRepository.UpdateTopStatus(l.ctx, oldestTopArticle.Id, 0); err != nil {
				return nil, errs.Wrap(errs.CodeBadDB, "取消最早置顶文章失败", err)
			}
		}
	}

	// 设置新文章为置顶
	if err := l.svcCtx.BlogArticleRepository.UpdateTopStatus(l.ctx, req.Id, 1); err != nil {
		return nil, errs.Wrap(errs.CodeBadDB, "设置文章置顶失败", err)
	}

	return &types.Response{
		Code:    0,
		Message: "置顶成功",
	}, nil
}
