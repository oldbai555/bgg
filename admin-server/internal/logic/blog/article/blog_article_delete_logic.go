// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package article

import (
	blogrepo "postapocgame/admin-server/internal/repository/blog"
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"

	"github.com/zeromicro/go-zero/core/logx"
)

type BlogArticleDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBlogArticleDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlogArticleDeleteLogic {
	return &BlogArticleDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BlogArticleDeleteLogic) BlogArticleDelete(req *types.BlogArticleDeleteReq) (resp *types.Response, err error) {
	if req.Id == 0 {
		return nil, errs.New(errs.CodeBadRequest, "文章ID不能为空")
	}

	if err := blogrepo.NewBlogArticleRepository(l.svcCtx.Repository).Delete(l.ctx, req.Id); err != nil {
		return nil, err
	}

	return &types.Response{Code: int(errs.CodeOK), Message: "删除成功"}, nil
}
