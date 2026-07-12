package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/content/content"
	contentconsts "postapocgame/admin-server/services/content/internal/consts"
	"postapocgame/admin-server/services/content/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type BlogArticleSubmitLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBlogArticleSubmitLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlogArticleSubmitLogic {
	return &BlogArticleSubmitLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// BlogArticleSubmit 迁移自 internal/logic/blog/article/blog_article_submit_logic.go。
// 同 BlogArticlePublish，原直连 Model 技术债顺手改成走 BlogArticleRepository.Update。
func (l *BlogArticleSubmitLogic) BlogArticleSubmit(in *content.BlogArticleSubmitRequest) (*content.Empty, error) {
	if in.Id == 0 {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "文章ID不能为空"))
	}

	article, err := l.svcCtx.BlogArticle.FindByID(l.ctx, in.Id)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeBadDB, "查询文章失败", err))
	}
	if article == nil || article.DeletedAt != 0 {
		return nil, toGRPCStatus(errs.New(errs.CodeNotFound, "文章不存在"))
	}

	if !(article.Status == contentconsts.BlogArticleStatusDraft || article.AuditStatus == contentconsts.BlogArticleAuditStatusRejected) {
		return nil, toGRPCStatus(errs.New(errs.CodeForbidden, "当前状态不允许提交审核"))
	}

	article.Status = contentconsts.BlogArticleStatusPendingAudit
	article.AuditStatus = contentconsts.BlogArticleAuditStatusPending

	if err := l.svcCtx.BlogArticle.Update(l.ctx, article); err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeBadDB, "提交审核失败", err))
	}

	return &content.Empty{}, nil
}
