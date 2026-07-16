package logic

import (
	"context"
	"strings"

	"postapocgame/admin-server/services/content/content"
	"postapocgame/admin-server/services/content/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type BlogArticleListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBlogArticleListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlogArticleListLogic {
	return &BlogArticleListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 文章

// BlogArticleList 迁移自 internal/logic/blog/article/blog_article_list_logic.go。
func (l *BlogArticleListLogic) BlogArticleList(in *content.BlogArticleListRequest) (*content.BlogArticleListResponse, error) {
	page := in.Page
	if page < 1 {
		page = 1
	}
	pageSize := in.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}

	list, total, err := l.svcCtx.BlogArticle.FindPage(
		l.ctx, page, pageSize, strings.TrimSpace(in.Title), in.Status, in.AuditStatus, in.TagId, in.StartTime, in.EndTime,
	)
	if err != nil {
		return nil, toGRPCStatus(err)
	}

	ids := make([]uint64, 0, len(list))
	for _, a := range list {
		ids = append(ids, a.Id)
	}

	tagMap, err := l.svcCtx.BlogArticleTag.FindTagsByArticleIDs(l.ctx, ids)
	if err != nil {
		return nil, toGRPCStatus(err)
	}

	items := make([]*content.BlogArticleItem, 0, len(list))
	for _, a := range list {
		tags := tagMap[a.Id]
		tagIDs := make([]uint64, 0, len(tags))
		tagNames := make([]string, 0, len(tags))
		for _, t := range tags {
			tagIDs = append(tagIDs, t.Id)
			tagNames = append(tagNames, t.Name)
		}
		items = append(items, &content.BlogArticleItem{
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

	return &content.BlogArticleListResponse{Total: total, List: items}, nil
}
