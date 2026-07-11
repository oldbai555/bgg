// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package article

import (
	"context"

	"postapocgame/admin-server/internal/consts"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"

	"github.com/zeromicro/go-zero/core/logx"
)

type BlogArticleUnpublishLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBlogArticleUnpublishLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlogArticleUnpublishLogic {
	return &BlogArticleUnpublishLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BlogArticleUnpublishLogic) BlogArticleUnpublish(req *types.BlogArticleUnpublishReq) (resp *types.Response, err error) {
	if req.Id == 0 {
		return nil, errs.New(errs.CodeBadRequest, "文章ID不能为空")
	}

	article, err := l.svcCtx.Domain.Blog.Article.FindByID(l.ctx, req.Id)
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

	return &types.Response{Code: int(errs.CodeOK), Message: "下架成功"}, nil
}
