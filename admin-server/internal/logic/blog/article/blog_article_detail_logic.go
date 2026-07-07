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

type BlogArticleDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBlogArticleDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlogArticleDetailLogic {
	return &BlogArticleDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BlogArticleDetailLogic) BlogArticleDetail(req *types.BlogArticleDetailReq) (resp *types.BlogArticleDetailResp, err error) {
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

	tags, err := blogrepo.NewBlogArticleTagRepository(l.svcCtx.Repository).FindTagsByArticleID(l.ctx, req.Id)
	if err != nil {
		return nil, err
	}

	tagItems := make([]types.BlogTagItem, 0, len(tags))
	for _, t := range tags {
		tagItems = append(tagItems, types.BlogTagItem{
			Id:        t.Id,
			Name:      t.Name,
			Status:    t.Status,
			Remark:    t.Remark,
			CreatedAt: t.CreatedAt,
			UpdatedAt: t.UpdatedAt,
		})
	}

	return &types.BlogArticleDetailResp{
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
