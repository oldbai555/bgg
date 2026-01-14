// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package public_blog

import (
	"context"
	"postapocgame/admin-server/internal/dict"
	"strings"
	"unicode/utf8"

	"github.com/zeromicro/go-zero/core/logx"
	"postapocgame/admin-server/internal/consts"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
)

type PublicBlogArticleListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPublicBlogArticleListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublicBlogArticleListLogic {
	return &PublicBlogArticleListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PublicBlogArticleListLogic) PublicBlogArticleList(req *types.PublicBlogArticleListReq) (resp *types.PublicBlogArticleListResp, err error) {
	page := req.Page
	if page < 1 {
		page = 1
	}
	size := req.Size
	if size <= 0 {
		size = 10
	}

	keyword := strings.TrimSpace(req.Keyword)
	list, total, err := l.svcCtx.BlogArticleRepository.FindPublicPage(l.ctx, page, size, keyword, req.TagId)
	if err != nil {
		return nil, err
	}

	ids := make([]uint64, 0, len(list))
	for _, a := range list {
		ids = append(ids, a.Id)
	}
	tagMap, err := l.svcCtx.BlogArticleTagRepository.FindTagsByArticleIDs(l.ctx, ids)
	if err != nil {
		return nil, err
	}

	// 从字典读取文章摘要截断长度（默认 120 个字符）
	maxLen := dict.GetIntValue(l.ctx, l.svcCtx.Repository, consts.DictCodeBlogArticleSummaryLength, 120)

	items := make([]types.PublicBlogArticleItem, 0, len(list))
	for _, a := range list {
		tagNames := make([]string, 0)
		for _, t := range tagMap[a.Id] {
			tagNames = append(tagNames, t.Name)
		}

		summary := a.Summary
		if summary == "" {
			summary = a.Content
		}
		summary = truncateByRune(summary, maxLen)

		items = append(items, types.PublicBlogArticleItem{
			Id:          a.Id,
			Title:       a.Title,
			Cover:       a.Cover,
			AuthorName:  a.AuthorName,
			Summary:     summary,
			TagNames:    tagNames,
			PublishTime: a.PublishTime,
		})
	}

	return &types.PublicBlogArticleListResp{
		List:  items,
		Page:  page,
		Size:  size,
		Total: total,
	}, nil
}

func truncateByRune(s string, max int) string {
	if max <= 0 || s == "" {
		return ""
	}
	if utf8.RuneCountInString(s) <= max {
		return s
	}
	// 简单截断（按 rune）
	out := make([]rune, 0, max)
	for _, r := range s {
		out = append(out, r)
		if len(out) >= max {
			break
		}
	}
	return string(out)
}
