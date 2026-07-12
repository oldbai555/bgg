package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/content/content"
	contentconsts "postapocgame/admin-server/services/content/internal/consts"
	"postapocgame/admin-server/services/content/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublicBlogArticleDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPublicBlogArticleDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublicBlogArticleDetailLogic {
	return &PublicBlogArticleDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// PublicBlogArticleDetail 迁移自
// internal/logic/blog/public/public_blog_article_detail_logic.go。
func (l *PublicBlogArticleDetailLogic) PublicBlogArticleDetail(in *content.PublicBlogArticleDetailRequest) (*content.PublicBlogArticleDetailResponse, error) {
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

	if article.AuditStatus != contentconsts.BlogArticleAuditStatusPassed || article.Status != contentconsts.BlogArticleStatusPublished {
		return nil, toGRPCStatus(errs.New(errs.CodeForbidden, "文章不可访问"))
	}

	tags, err := l.svcCtx.BlogArticleTag.FindTagsByArticleID(l.ctx, in.Id)
	if err != nil {
		return nil, toGRPCStatus(err)
	}
	tagItems := make([]*content.BlogTagItem, 0, len(tags))
	for _, t := range tags {
		tagItems = append(tagItems, &content.BlogTagItem{
			Id:        t.Id,
			Name:      t.Name,
			Status:    t.Status,
			Remark:    t.Remark,
			CreatedAt: t.CreatedAt,
			UpdatedAt: t.UpdatedAt,
		})
	}

	return &content.PublicBlogArticleDetailResponse{
		Id:          article.Id,
		Title:       article.Title,
		Content:     article.Content,
		Cover:       article.Cover,
		AuthorName:  article.AuthorName,
		PublishTime: article.PublishTime,
		Tags:        tagItems,
	}, nil
}
