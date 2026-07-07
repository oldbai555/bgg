// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package public

import (
	blogrepo "postapocgame/admin-server/internal/repository/blog"
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublicBlogArticleNextLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPublicBlogArticleNextLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublicBlogArticleNextLogic {
	return &PublicBlogArticleNextLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PublicBlogArticleNextLogic) PublicBlogArticleNext(req *types.PublicBlogArticleNextReq) (resp *types.PublicBlogArticleNextResp, err error) {
	if req.Id == 0 {
		return nil, errs.New(errs.CodeBadRequest, "文章ID不能为空")
	}

	// 查询当前文章信息
	currentArticle, err := blogrepo.NewBlogArticleRepository(l.svcCtx.Repository).FindByID(l.ctx, req.Id)
	if err != nil {
		return nil, errs.Wrap(errs.CodeBadDB, "查询当前文章失败", err)
	}
	if currentArticle == nil {
		return nil, errs.New(errs.CodeNotFound, "文章不存在")
	}

	// 查询下一篇文章
	nextArticle, err := blogrepo.NewBlogArticleRepository(l.svcCtx.Repository).FindNextArticle(l.ctx, currentArticle.PublishTime)
	if err != nil {
		return nil, errs.Wrap(errs.CodeBadDB, "查询下一篇文章失败", err)
	}

	// 如果没有下一篇文章，返回空值
	if nextArticle == nil {
		return &types.PublicBlogArticleNextResp{}, nil
	}

	return &types.PublicBlogArticleNextResp{
		Id:          nextArticle.Id,
		Title:       nextArticle.Title,
		PublishTime: nextArticle.PublishTime,
	}, nil
}
