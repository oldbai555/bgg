// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package public_blog

import (
	"context"
	"postapocgame/admin-server/internal/dict"
	"regexp"
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
			// 如果摘要为空，使用正文内容，但需要去除 Markdown 格式
			summary = stripMarkdown(a.Content)
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
			IsTop:       a.IsTop,
		})
	}

	return &types.PublicBlogArticleListResp{
		List:  items,
		Page:  page,
		Size:  size,
		Total: total,
	}, nil
}

// stripMarkdown 去除 Markdown 格式，转换为纯文本
func stripMarkdown(text string) string {
	if text == "" {
		return ""
	}

	// 1. 去除代码块（```code```）
	codeBlockRegex := regexp.MustCompile("(?s)```.*?```")
	text = codeBlockRegex.ReplaceAllString(text, "")

	// 2. 去除行内代码（`code`）
	inlineCodeRegex := regexp.MustCompile("`[^`]+`")
	text = inlineCodeRegex.ReplaceAllString(text, "")

	// 3. 去除图片（![alt](url)）
	imageRegex := regexp.MustCompile(`!\[([^\]]*)\]\([^\)]+\)`)
	text = imageRegex.ReplaceAllString(text, "$1")

	// 4. 去除链接，保留文本（[text](url)）
	linkRegex := regexp.MustCompile(`\[([^\]]+)\]\([^\)]+\)`)
	text = linkRegex.ReplaceAllString(text, "$1")

	// 5. 去除标题符号（# ## ### 等）
	headingRegex := regexp.MustCompile(`^#{1,6}\s+`)
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		lines[i] = headingRegex.ReplaceAllString(line, "")
	}
	text = strings.Join(lines, "\n")

	// 6. 去除加粗（**text** 或 __text__），先处理加粗避免与斜体冲突
	// 使用非贪婪匹配处理 **text** 格式
	boldRegex := regexp.MustCompile(`\*\*([^\*]+?)\*\*`)
	text = boldRegex.ReplaceAllString(text, "$1")
	// 处理 __text__ 格式
	boldRegex2 := regexp.MustCompile(`__([^_]+?)__`)
	text = boldRegex2.ReplaceAllString(text, "$1")

	// 7. 去除斜体（*text* 或 _text_），只处理单个符号的情况（前后不能是相同符号）
	italicRegex := regexp.MustCompile(`([^\*]|^)\*([^\*\n]+?)\*([^\*]|$)`)
	text = italicRegex.ReplaceAllString(text, "$1$2$3")
	italicRegex2 := regexp.MustCompile(`([^_]|^)_([^_\n]+?)_([^_]|$)`)
	text = italicRegex2.ReplaceAllString(text, "$1$2$3")

	// 8. 去除列表符号（- * + 1. 等）
	listRegex := regexp.MustCompile(`^\s*[-*+]\s+|^\s*\d+\.\s+`)
	lines = strings.Split(text, "\n")
	for i, line := range lines {
		lines[i] = listRegex.ReplaceAllString(line, "")
	}
	text = strings.Join(lines, "\n")

	// 9. 去除引用符号（> ）
	quoteRegex := regexp.MustCompile(`^\s*>\s+`)
	lines = strings.Split(text, "\n")
	for i, line := range lines {
		lines[i] = quoteRegex.ReplaceAllString(line, "")
	}
	text = strings.Join(lines, "\n")

	// 10. 去除水平线（--- 或 ***）
	hrRegex := regexp.MustCompile(`^[-*]{3,}\s*$`)
	lines = strings.Split(text, "\n")
	for i, line := range lines {
		lines[i] = hrRegex.ReplaceAllString(line, "")
	}
	text = strings.Join(lines, "\n")

	// 11. 去除表格符号（|）
	text = strings.ReplaceAll(text, "|", " ")

	// 12. 合并多个连续空白字符为单个空格
	spaceRegex := regexp.MustCompile(`\s+`)
	text = spaceRegex.ReplaceAllString(text, " ")

	// 13. 去除首尾空白
	text = strings.TrimSpace(text)

	return text
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
