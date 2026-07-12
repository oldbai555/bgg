package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/content/content"
	"postapocgame/admin-server/services/content/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type BlogArticleDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBlogArticleDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlogArticleDetailLogic {
	return &BlogArticleDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// BlogArticleDetail 迁移自 internal/logic/blog/article/blog_article_detail_logic.go。
func (l *BlogArticleDetailLogic) BlogArticleDetail(in *content.BlogArticleDetailRequest) (*content.BlogArticleDetailResponse, error) {
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

	return &content.BlogArticleDetailResponse{
		Id:          article.Id,
		Title:       article.Title,
		Content:     article.Content,
		Status:      article.Status,
		AuditStatus: article.AuditStatus,
		Cover:       article.Cover,
		AuthorId:    article.AuthorId,
		AuthorName:  article.AuthorName,
		PublishTime: article.PublishTime,
		Summary:     article.Summary,
		Tags:        tagItems,
		CreatedAt:   article.CreatedAt,
		UpdatedAt:   article.UpdatedAt,
	}, nil
}
