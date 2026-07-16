package progresstools

import (
	"context"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type checklistEntry struct {
	Index        int    `json:"index"`
	Title        string `json:"title"`
	Trigger      string `json:"trigger"`
	Action       string `json:"action"`
	Verification string `json:"verification"`
	Status       string `json:"status"`
}

var (
	checklistHeadingRe = regexp.MustCompile(`(?m)^### (\d+)\s*·\s*(.*)$`)
	statusValueRe      = regexp.MustCompile("`([^`]+)`")
)

func registerQueryDeploymentChecklist(s *server.MCPServer, docsDir string) {
	s.AddTool(mcp.NewTool("query_deployment_checklist",
		mcp.WithDescription("解析 docs/14-production-deployment-checklist.md，返回触发条件/部署步骤/验证方式/状态四字段结构化条目"),
		mcp.WithBoolean("pending_only", mcp.Description("true 时只返回状态非「已执行」的条目，默认 false 返回全部")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		raw, err := os.ReadFile(filepath.Join(docsDir, "14-production-deployment-checklist.md"))
		if err != nil {
			return mcp.NewToolResultErrorFromErr("读取 14-production-deployment-checklist.md 失败", err), nil
		}

		entries := parseChecklistEntries(string(raw))
		pendingOnly := req.GetBool("pending_only", false)
		if pendingOnly {
			filtered := make([]checklistEntry, 0, len(entries))
			for _, e := range entries {
				if e.Status != "已执行" {
					filtered = append(filtered, e)
				}
			}
			entries = filtered
		}

		return mcp.NewToolResultStructured(map[string]any{"entries": entries}, "共 "+itoa(len(entries))+" 条"), nil
	})
}

func parseChecklistEntries(content string) []checklistEntry {
	locs := checklistHeadingRe.FindAllStringSubmatchIndex(content, -1)
	matches := checklistHeadingRe.FindAllStringSubmatch(content, -1)

	entries := make([]checklistEntry, 0, len(matches))
	for i, m := range matches {
		start := locs[i][1]
		end := len(content)
		if i+1 < len(locs) {
			end = locs[i+1][0]
		}
		body := content[start:end]

		trigger, action, verification, rawStatus := extractChecklistFields(body)
		entries = append(entries, checklistEntry{
			Index:        atoi(m[1]),
			Title:        strings.TrimSpace(m[2]),
			Trigger:      trigger,
			Action:       action,
			Verification: verification,
			Status:       classifyStatus(rawStatus),
		})
	}
	return entries
}

// extractChecklistFields 按 **触发条件**/**部署时要做什么**/**如何验证生效**/**状态**
// 四个固定字段标记切分正文，标记之间的全部内容原样返回（trim 首尾空白）。
func extractChecklistFields(body string) (trigger, action, verification, status string) {
	markers := []string{"**触发条件**", "**部署时要做什么**", "**如何验证生效**", "**状态**"}
	positions := make([]int, len(markers))
	for i, m := range markers {
		positions[i] = strings.Index(body, m)
	}

	get := func(i int) string {
		if positions[i] < 0 {
			return ""
		}
		start := positions[i] + len(markers[i])
		start = skipFieldColon(body, start)
		end := len(body)
		for j := i + 1; j < len(markers); j++ {
			if positions[j] >= 0 {
				end = positions[j]
				break
			}
		}
		if end < start {
			end = start
		}
		return strings.TrimSpace(body[start:end])
	}

	return get(0), get(1), get(2), get(3)
}

// skipFieldColon 跳过字段标记后紧跟的全角/半角冒号。
func skipFieldColon(s string, pos int) int {
	i := pos
	for i < len(s) && (s[i] == ' ' || s[i] == '\t') {
		i++
	}
	if strings.HasPrefix(s[i:], "：") {
		return i + len("：")
	}
	if i < len(s) && s[i] == ':' {
		return i + 1
	}
	return i
}

// classifyStatus 从状态字段原文中抠出反引号包裹的三态值之一：TBD / 已就绪，待执行 / 已执行。
func classifyStatus(rawStatus string) string {
	m := statusValueRe.FindStringSubmatch(rawStatus)
	if len(m) < 2 {
		return "TBD"
	}
	switch m[1] {
	case "TBD", "已就绪，待执行", "已执行":
		return m[1]
	default:
		return "TBD"
	}
}

func atoi(s string) int {
	n := 0
	for _, c := range s {
		if c < '0' || c > '9' {
			break
		}
		n = n*10 + int(c-'0')
	}
	return n
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	digits := []byte{}
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	return string(digits)
}
