package logic

import (
	"context"
	"time"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/content/content"
	contentconsts "postapocgame/admin-server/services/content/internal/consts"
	"postapocgame/admin-server/services/content/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type BlogArticlePublishLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBlogArticlePublishLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlogArticlePublishLogic {
	return &BlogArticlePublishLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// BlogArticlePublish 迁移自 internal/logic/blog/article/blog_article_publish_logic.go。
// 原实现直连 l.svcCtx.Repository.BlogArticleModel.Update（11-descoped.md 第 10 条记录的
// 直连 Model 技术债），content-rpc 拆分后 gateway 的 *repository.Repository 已经不存在，
// 顺手改成走 BlogArticleRepository.Update（该方法早就存在，行为完全一致）。
func (l *BlogArticlePublishLogic) BlogArticlePublish(in *content.BlogArticlePublishRequest) (*content.Empty, error) {
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

	if article.AuditStatus != contentconsts.BlogArticleAuditStatusPassed {
		return nil, toGRPCStatus(errs.New(errs.CodeForbidden, "文章未审核通过，不能上架"))
	}

	article.Status = contentconsts.BlogArticleStatusPublished
	if article.PublishTime == 0 {
		article.PublishTime = time.Now().Unix()
	}

	if err := l.svcCtx.BlogArticle.Update(l.ctx, article); err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeBadDB, "上架失败", err))
	}

	return &content.Empty{}, nil
}
