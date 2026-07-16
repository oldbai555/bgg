// Package conventiontools 查询本项目的架构约定：Phase 2 服务边界、文件生成策略、
// 中间件声明顺序。数据来源见各文件顶部注释，本工具只读，不修改这些来源文档/文件。
package conventiontools

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"gopkg.in/yaml.v3"
)

type serviceBoundaries struct {
	Services []serviceEntry `yaml:"services"`
	Gateway  struct {
		Name string `yaml:"name"`
		Note string `yaml:"note"`
	} `yaml:"gateway"`
}

type serviceEntry struct {
	Name       string   `yaml:"name"`
	DBSchema   string   `yaml:"db_schema"`
	Domains    []string `yaml:"domains"`
	TableCount int      `yaml:"table_count"`
}

// Register 注册 service_boundary/file_policy/middleware_order 三个约定查询 tool。
// dataDir 是 admin-mcp/data 目录的绝对路径。
func Register(s *server.MCPServer, dataDir string) {
	registerServiceBoundary(s, dataDir)
	registerFilePolicy(s, dataDir)
	registerMiddlewareOrder(s)
}

func registerServiceBoundary(s *server.MCPServer, dataDir string) {
	s.AddTool(mcp.NewTool("query_service_boundary",
		mcp.WithDescription("查询某个业务域（domain）在 Phase 2 拆分后归属哪个 RPC 服务，数据来自 docs/15-service-boundaries.md"),
		mcp.WithString("domain", mcp.Required(), mcp.Description("业务域名，大小写不敏感，支持传入旧 group 前缀（如 blog_article 也能匹配到 blog）")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		domain, err := req.RequireString("domain")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		boundaries, err := loadServiceBoundaries(dataDir)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("加载 service-boundaries.yaml 失败", err), nil
		}

		needle := strings.ToLower(domain)
		for _, svc := range boundaries.Services {
			for _, d := range svc.Domains {
				if strings.HasPrefix(needle, strings.ToLower(d)) {
					siblings := make([]string, 0, len(svc.Domains)-1)
					for _, sib := range svc.Domains {
						if !strings.EqualFold(sib, d) {
							siblings = append(siblings, sib)
						}
					}
					return mcp.NewToolResultStructured(map[string]any{
						"service":         svc.Name,
						"db_schema":       svc.DBSchema,
						"matched_domain":  d,
						"sibling_domains": siblings,
					}, fmt.Sprintf("domain=%s -> service=%s (db_schema=%s)", domain, svc.Name, svc.DBSchema)), nil
				}
			}
		}

		return mcp.NewToolResultError(fmt.Sprintf("未找到 domain=%q 匹配的服务边界，检查拼写或确认 15-service-boundaries.md 是否已更新", domain)), nil
	})
}

func loadServiceBoundaries(dataDir string) (*serviceBoundaries, error) {
	raw, err := os.ReadFile(filepath.Join(dataDir, "service-boundaries.yaml"))
	if err != nil {
		return nil, err
	}
	var b serviceBoundaries
	if err := yaml.Unmarshal(raw, &b); err != nil {
		return nil, err
	}
	return &b, nil
}
