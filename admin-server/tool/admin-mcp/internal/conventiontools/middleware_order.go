package conventiontools

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// middlewareOrder 是硬编码的固定结构，不读文件——这是 AGENTS.md/10-go-code-style.md
// 里明确写死不会变的顺序，编码进 Go 代码比再建一个只有一行数据的 YAML 更直接。
type middlewareOrder struct {
	AdminOrder  []string `json:"admin_order"`
	SDKOrder    []string `json:"sdk_order"`
	PublicOrder []string `json:"public_order"`
	Notes       []string `json:"notes"`
}

func registerMiddlewareOrder(s *server.MCPServer) {
	s.AddTool(mcp.NewTool("query_middleware_order",
		mcp.WithDescription("查询 .api 文件 middleware 字段的强制声明顺序，数据来自 AGENTS.md 第 3 节 / 10-go-code-style.md"),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		order := middlewareOrder{
			AdminOrder:  []string{"Performance", "RateLimit", "Auth", "Permission", "OperationLog"},
			SDKOrder:    []string{"SDKAuth", "SDKRateLimit", "SDKCallLog"},
			PublicOrder: []string{"Performance"},
			Notes: []string{
				"Permission 与 ApiEnabled 互斥，一个 .api 服务块只能选一个",
				"SDK 系列中间件不与 Admin 系列混用",
				"公开/无需登录接口只用 Performance（+ApiEnabled 视情况）",
			},
		}
		return mcp.NewToolResultStructured(order, "见 admin_order/sdk_order/public_order 字段"), nil
	})
}
