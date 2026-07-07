// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package article

import (
	blogrepo "postapocgame/admin-server/internal/repository/blog"
	"context"

	"postapocgame/admin-server/internal/consts"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"

	"github.com/zeromicro/go-zero/core/logx"
)

type BlogArticleSubmitLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBlogArticleSubmitLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlogArticleSubmitLogic {
	return &BlogArticleSubmitLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BlogArticleSubmitLogic) BlogArticleSubmit(req *types.BlogArticleSubmitReq) (resp *types.Response, err error) {
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

	// 草稿 或 审核驳回 才允许提交审核
	if !(article.Status == consts.BlogArticleStatusDraft || article.AuditStatus == consts.BlogArticleAuditStatusRejected) {
		return nil, errs.New(errs.CodeForbidden, "当前状态不允许提交审核")
	}

	article.Status = consts.BlogArticleStatusPendingAudit
	article.AuditStatus = consts.BlogArticleAuditStatusPending

	if err := l.svcCtx.Repository.BlogArticleModel.Update(l.ctx, article); err != nil {
		return nil, errs.Wrap(errs.CodeBadDB, "提交审核失败", err)
	}

	return &types.Response{Code: int(errs.CodeOK), Message: "已提交审核"}, nil
}
