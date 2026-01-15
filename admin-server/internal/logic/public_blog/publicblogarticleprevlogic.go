// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package public_blog

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublicBlogArticlePrevLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPublicBlogArticlePrevLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublicBlogArticlePrevLogic {
	return &PublicBlogArticlePrevLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PublicBlogArticlePrevLogic) PublicBlogArticlePrev(req *types.PublicBlogArticlePrevReq) (resp *types.PublicBlogArticlePrevResp, err error) {
	if req.Id == 0 {
		return nil, errs.New(errs.CodeBadRequest, "文章ID不能为空")
	}

	// 查询当前文章信息
	currentArticle, err := l.svcCtx.BlogArticleRepository.FindByID(l.ctx, req.Id)
	if err != nil {
		return nil, errs.Wrap(errs.CodeBadDB, "查询当前文章失败", err)
	}
	if currentArticle == nil {
		return nil, errs.New(errs.CodeNotFound, "文章不存在")
	}

	// 查询上一篇文章
	prevArticle, err := l.svcCtx.BlogArticleRepository.FindPrevArticle(l.ctx, currentArticle.PublishTime)
	if err != nil {
		return nil, errs.Wrap(errs.CodeBadDB, "查询上一篇文章失败", err)
	}

	// 如果没有上一篇文章，返回空值
	if prevArticle == nil {
		return &types.PublicBlogArticlePrevResp{}, nil
	}

	return &types.PublicBlogArticlePrevResp{
		Id:          prevArticle.Id,
		Title:       prevArticle.Title,
		PublishTime: prevArticle.PublishTime,
	}, nil
}
