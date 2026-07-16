// Package progresstools 解析 progress.md / 14-production-deployment-checklist.md，
// 把自由格式 Markdown 切分成结构化条目返回；摘要/理解交给调用方的 LLM 做，本工具只负责
// 结构化切分,不做二次摘要。
package progresstools

import (
	"context"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type progressEntry struct {
	Date  string `json:"date"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

var progressHeadingRe = regexp.MustCompile(`(?m)^## (\d{4}-\d{2}-\d{2})[：:]\s*(.*)$`)

// Register 注册 query_progress/query_deployment_checklist 两个 tool。
// docsDir 是 admin-server/docs 目录的绝对路径。
func Register(s *server.MCPServer, docsDir string) {
	registerQueryProgress(s, docsDir)
	registerQueryDeploymentChecklist(s, docsDir)
}

func registerQueryProgress(s *server.MCPServer, docsDir string) {
	s.AddTool(mcp.NewTool("query_progress",
		mcp.WithDescription("解析 docs/progress.md，按日期倒序返回全部条目（贯穿 Phase 1-3 的重构进度记录）"),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		raw, err := os.ReadFile(filepath.Join(docsDir, "progress.md"))
		if err != nil {
			return mcp.NewToolResultErrorFromErr("读取 progress.md 失败", err), nil
		}

		entries := parseProgressEntries(string(raw))
		return mcp.NewToolResultStructured(map[string]any{"entries": entries}, summarizeEntries(entries)), nil
	})
}

func parseProgressEntries(content string) []progressEntry {
	locs := progressHeadingRe.FindAllStringSubmatchIndex(content, -1)
	matches := progressHeadingRe.FindAllStringSubmatch(content, -1)

	entries := make([]progressEntry, 0, len(matches))
	for i, m := range matches {
		start := locs[i][1] // 标题行结束位置
		end := len(content)
		if i+1 < len(locs) {
			end = locs[i+1][0]
		}
		entries = append(entries, progressEntry{
			Date:  m[1],
			Title: strings.TrimSpace(m[2]),
			Body:  strings.TrimSpace(content[start:end]),
		})
	}

	sort.Slice(entries, func(i, j int) bool { return entries[i].Date > entries[j].Date })
	return entries
}

func summarizeEntries(entries []progressEntry) string {
	if len(entries) == 0 {
		return "progress.md 未解析出任何条目"
	}
	return entries[0].Date + "：" + entries[0].Title
}
