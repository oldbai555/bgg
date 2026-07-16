package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/content/content"
	contentconsts "postapocgame/admin-server/services/content/internal/consts"
	"postapocgame/admin-server/services/content/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type BlogArticleUnpublishLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBlogArticleUnpublishLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlogArticleUnpublishLogic {
	return &BlogArticleUnpublishLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// BlogArticleUnpublish 迁移自 internal/logic/blog/article/blog_article_unpublish_logic.go
// （管理端"我要下架自己文章"的简化入口，不写审核记录；带审核记录的下架走
// BlogArticleAuditUnpublish/UnpublishArticle）。同样把原直连 Model 技术债改成走
// BlogArticleRepository.Update。
func (l *BlogArticleUnpublishLogic) BlogArticleUnpublish(in *content.BlogArticleUnpublishRequest) (*content.Empty, error) {
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

	if article.Status != contentconsts.BlogArticleStatusPublished {
		return nil, toGRPCStatus(errs.New(errs.CodeForbidden, "仅已上架文章可下架"))
	}

	article.Status = contentconsts.BlogArticleStatusUnpublished
	if err := l.svcCtx.BlogArticle.Update(l.ctx, article); err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeBadDB, "下架失败", err))
	}

	return &content.Empty{}, nil
}
