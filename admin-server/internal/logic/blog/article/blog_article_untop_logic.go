// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package article

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"

	"github.com/zeromicro/go-zero/core/logx"
)

type BlogArticleUntopLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBlogArticleUntopLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlogArticleUntopLogic {
	return &BlogArticleUntopLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BlogArticleUntopLogic) BlogArticleUntop(req *types.BlogArticleUntopReq) (resp *types.Response, err error) {
	if req.Id == 0 {
		return nil, errs.New(errs.CodeBadRequest, "文章ID不能为空")
	}

	// 直接执行更新，UpdateTopStatus 内部会检查文章是否存在
	// 这样可以避免缓存问题，直接从数据库查询最新状态
	if err := l.svcCtx.Domain.Blog.Article.UpdateTopStatus(l.ctx, req.Id, 0); err != nil {
		return nil, err
	}

	return &types.Response{
		Code:    0,
		Message: "取消置顶成功",
	}, nil
}
