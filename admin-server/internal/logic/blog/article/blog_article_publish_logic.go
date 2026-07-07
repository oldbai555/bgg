// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package article

import (
	blogrepo "postapocgame/admin-server/internal/repository/blog"
	"context"
	"time"

	"postapocgame/admin-server/internal/consts"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"

	"github.com/zeromicro/go-zero/core/logx"
)

type BlogArticlePublishLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBlogArticlePublishLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlogArticlePublishLogic {
	return &BlogArticlePublishLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BlogArticlePublishLogic) BlogArticlePublish(req *types.BlogArticlePublishReq) (resp *types.Response, err error) {
	if req.Id == 0 {
		return nil, errs.New(errs.CodeBadRequest, "文章ID不能为空")
	}

	article, err := blogrepo.NewBlogArticleRepository(l.svcCtx.Repository).FindByID(l.ctx, req.Id)
	if err != nil {
		return nil, errs.Wrap(errs.CodeBadDB, "查询文章失败", err)
	}
	if article == nil || article.DeletedAt != 0 {
		return nil, errs.New(errs.CodeNotFound, "文章不存在")
	}

	// 必须审核通过才能上架
	if article.AuditStatus != consts.BlogArticleAuditStatusPassed {
		return nil, errs.New(errs.CodeForbidden, "文章未审核通过，不能上架")
	}

	article.Status = consts.BlogArticleStatusPublished
	if article.PublishTime == 0 {
		article.PublishTime = time.Now().Unix()
	}

	if err := l.svcCtx.Repository.BlogArticleModel.Update(l.ctx, article); err != nil {
		return nil, errs.Wrap(errs.CodeBadDB, "上架失败", err)
	}

	return &types.Response{Code: int(errs.CodeOK), Message: "上架成功"}, nil
}
