// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package article

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

	// 从字典读取置顶最大数量限制（默认1篇）
	maxCount := int64(dict.GetIntValue(l.ctx, l.svcCtx.Repository, consts.DictCodeBlogArticleTopMaxCount, 1))

	if err := l.svcCtx.Domain.Blog.ArticleService.SetArticleTop(l.ctx, req.Id, maxCount); err != nil {
		return nil, err
	}

	return &types.Response{
		Code:    0,
		Message: "置顶成功",
	}, nil
}
