package logic

import (
	"context"
	"strings"

	"postapocgame/admin-server/services/content/content"
	"postapocgame/admin-server/services/content/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublicBlogArticleListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPublicBlogArticleListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublicBlogArticleListLogic {
	return &PublicBlogArticleListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 公共博客展示

// PublicBlogArticleList 迁移自 internal/logic/blog/public/public_blog_article_list_logic.go。
// 摘要截断长度原来读字典 blog_article_summary_length（物理属于 iam 域），改成
// svcCtx.Config.Limits.BlogArticleSummaryLength 静态配置。
func (l *PublicBlogArticleListLogic) PublicBlogArticleList(in *content.PublicBlogArticleListRequest) (*content.PublicBlogArticleListResponse, error) {
	page := in.Page
	if page < 1 {
		page = 1
	}
	size := in.Size
	if size <= 0 {
		size = 10
	}

	keyword := strings.TrimSpace(in.Keyword)
	list, total, err := l.svcCtx.BlogArticle.FindPublicPage(l.ctx, page, size, keyword, in.TagId)
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

	maxLen := l.svcCtx.Config.Limits.BlogArticleSummaryLength

	items := make([]*content.PublicBlogArticleItem, 0, len(list))
	for _, a := range list {
		tagNames := make([]string, 0)
		for _, t := range tagMap[a.Id] {
			tagNames = append(tagNames, t.Name)
		}

		summary := a.Summary
		if summary == "" {
			summary = stripMarkdown(a.Content)
		}
		summary = truncateByRune(summary, maxLen)

		items = append(items, &content.PublicBlogArticleItem{
			Id:          a.Id,
			Title:       a.Title,
			Cover:       a.Cover,
			AuthorName:  a.AuthorName,
			Summary:     summary,
			TagNames:    tagNames,
			PublishTime: a.PublishTime,
			IsTop:       a.IsTop,
		})
	}

	return &content.PublicBlogArticleListResponse{List: items, Page: page, Size: size, Total: total}, nil
}
