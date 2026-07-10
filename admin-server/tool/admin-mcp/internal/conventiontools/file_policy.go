package conventiontools

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"gopkg.in/yaml.v3"
)

type filePolicyRules struct {
	Rules []filePolicyRule `yaml:"rules"`
}

type filePolicyRule struct {
	Glob       string   `yaml:"glob"`
	Exceptions []string `yaml:"exceptions"`
	Policy     string   `yaml:"policy"`
	Note       string   `yaml:"note"`
}

func registerFilePolicy(s *server.MCPServer, dataDir string) {
	s.AddTool(mcp.NewTool("query_file_policy",
		mcp.WithDescription("查询某个文件路径的编辑策略（禁止手改的生成产物 / 手改 sibling 文件 / 纯手写 / 必须停下来问用户），数据来自 AGENTS.md 第 3、5 节"),
		mcp.WithString("path", mcp.Required(), mcp.Description("文件路径，相对仓库根或绝对路径均可")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		path, err := req.RequireString("path")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		rules, err := loadFilePolicyRules(dataDir)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("加载 file-policy-rules.yaml 失败", err), nil
		}

		normalized := normalizeRepoPath(path)

		for _, rule := range rules.Rules {
			matched, err := doublestarMatch(rule.Glob, normalized)
			if err != nil {
				return mcp.NewToolResultErrorFromErr(fmt.Sprintf("glob 规则 %q 编译失败", rule.Glob), err), nil
			}
			if !matched {
				continue
			}
			if isException(rule.Exceptions, normalized) {
				continue
			}
			return mcp.NewToolResultStructured(map[string]any{
				"policy":          rule.Policy,
				"note":            rule.Note,
				"matched_glob":    rule.Glob,
				"normalized_path": normalized,
			}, fmt.Sprintf("path=%s -> policy=%s (%s)", normalized, rule.Policy, rule.Note)), nil
		}

		// 理论上不会到这里：规则表最后一条 "**/*" 兜底所有路径。
		return mcp.NewToolResultError("未匹配到任何规则，且兜底规则缺失，检查 file-policy-rules.yaml"), nil
	})
}

func isException(exceptions []string, normalized string) bool {
	for _, ex := range exceptions {
		if ok, _ := doublestarMatch(ex, normalized); ok {
			return true
		}
	}
	return false
}

// doublestarMatch 是 doublestar.Match 的薄包装，供本包内匹配逻辑和测试共用同一入口。
func doublestarMatch(glob, path string) (bool, error) {
	return doublestar.Match(glob, path)
}

// normalizeRepoPath 把绝对路径/相对路径统一成相对仓库根、正斜杠分隔的形式，
// 特殊保留 /etc/work/*.json 这类以 / 开头、代表生产环境绝对路径的规则。
func normalizeRepoPath(path string) string {
	p := filepath.ToSlash(path)
	if strings.HasPrefix(p, "/etc/work/") {
		return p
	}
	p = strings.TrimPrefix(p, "/")
	// 允许调用方传入 admin-server/ 前缀或不带前缀，规则表本身是相对 admin-server/ 写的，
	// 这里统一剥离常见的仓库前缀，容忍两种写法。
	p = strings.TrimPrefix(p, "admin-server/")
	return p
}

func loadFilePolicyRules(dataDir string) (*filePolicyRules, error) {
	raw, err := os.ReadFile(filepath.Join(dataDir, "file-policy-rules.yaml"))
	if err != nil {
		return nil, err
	}
	var r filePolicyRules
	if err := yaml.Unmarshal(raw, &r); err != nil {
		return nil, err
	}
	return &r, nil
}
