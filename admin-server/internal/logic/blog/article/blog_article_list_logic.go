// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package article

import (
	"context"
	"strings"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type BlogArticleListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBlogArticleListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlogArticleListLogic {
	return &BlogArticleListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BlogArticleListLogic) BlogArticleList(req *types.BlogArticleListReq) (resp *types.BlogArticleListResp, err error) {
	page := req.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}

	list, total, err := l.svcCtx.Domain.Blog.Article.FindPage(
		l.ctx,
		page,
		pageSize,
		strings.TrimSpace(req.Title),
		req.Status,
		req.AuditStatus,
		req.TagId,
		req.StartTime,
		req.EndTime,
	)
	if err != nil {
		return nil, err
	}

	ids := make([]uint64, 0, len(list))
	for _, a := range list {
		ids = append(ids, a.Id)
	}

	tagMap, err := l.svcCtx.Domain.Blog.ArticleTag.FindTagsByArticleIDs(l.ctx, ids)
	if err != nil {
		return nil, err
	}

	items := make([]types.BlogArticleItem, 0, len(list))
	for _, a := range list {
		tags := tagMap[a.Id]
		tagIDs := make([]uint64, 0, len(tags))
		tagNames := make([]string, 0, len(tags))
		for _, t := range tags {
			tagIDs = append(tagIDs, t.Id)
			tagNames = append(tagNames, t.Name)
		}
		items = append(items, types.BlogArticleItem{
			Id:          a.Id,
			Title:       a.Title,
			Status:      a.Status,
			AuditStatus: a.AuditStatus,
			Cover:       a.Cover,
			AuthorId:    a.AuthorId,
			AuthorName:  a.AuthorName,
			TagIds:      tagIDs,
			TagNames:    tagNames,
			PublishTime: a.PublishTime,
			IsTop:       a.IsTop,
			CreatedAt:   a.CreatedAt,
			UpdatedAt:   a.UpdatedAt,
		})
	}

	return &types.BlogArticleListResp{
		Total: total,
		List:  items,
	}, nil
}
