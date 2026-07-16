package logic

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

// stripMarkdown 去除 Markdown 格式，转换为纯文本，逐字迁移自
// internal/logic/blog/public/public_blog_article_list_logic.go。
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
	boldRegex := regexp.MustCompile(`\*\*([^\*]+?)\*\*`)
	text = boldRegex.ReplaceAllString(text, "$1")
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

func truncateByRune(s string, max int64) string {
	if max <= 0 || s == "" {
		return ""
	}
	if int64(utf8.RuneCountInString(s)) <= max {
		return s
	}
	out := make([]rune, 0, max)
	for _, r := range s {
		out = append(out, r)
		if int64(len(out)) >= max {
			break
		}
	}
	return string(out)
}
