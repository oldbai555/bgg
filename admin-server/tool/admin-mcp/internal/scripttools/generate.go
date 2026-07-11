// Package scripttools 封装 admin-server/scripts/generate-*.sh 系列脚本，把"AI 自己拼命令行"
// 变成"AI 调一个有明确 schema 的 tool"。不替代脚本本身——生成逻辑仍然是这四个 shell 脚本 +
// scripts/sqlgen，脚本改了这里的 wrapper 也要跟着改，不是反过来。
package scripttools

import (
	"context"
	"path/filepath"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"admin-mcp/internal/exec"
)

// Register 注册全部 6 个脚本封装 tool。repoRoot 是 admin-server/ 目录的绝对路径。
func Register(s *server.MCPServer, repoRoot string) {
	s.AddTool(mcp.NewTool("generate_sql",
		mcp.WithDescription("封装 scripts/generate-sql.sh：一次性生成表结构 SQL、RBAC 初始化数据 SQL、.api 草稿、前端列表页骨架"),
		mcp.WithString("group", mcp.Required(), mcp.Description("功能组名，格式 <domain>/<module>（如 iam/user、blog/article），与 .api 文件 "+
			"@server(group:...) 的格式一致。<domain> 决定落进哪个服务的 db/services/<service>/ 目录（iam/system/monitoring/misc→iam，"+
			"blog/video→content，chat/task/sdk 各自独立），module 为 snake_case")),
		mcp.WithString("name", mcp.Required(), mcp.Description("功能中文名称，如 用户管理")),
		mcp.WithString("parent_id", mcp.Description("父菜单 ID，可选，优先级最高")),
		mcp.WithString("parent_path", mcp.Description("前端父目录路径，可选，默认 /temp")),
	), handleGenerateSQL(repoRoot))

	s.AddTool(mcp.NewTool("generate_model",
		mcp.WithDescription("封装 scripts/generate-model.sh：从建表 SQL 生成 goctl Model 代码"),
		mcp.WithString("migration_file", mcp.Required(), mcp.Description("建表 SQL 文件路径，相对 db/services/<service>/<module>/ 或绝对路径")),
		mcp.WithString("dir", mcp.Description("输出目录，可选，默认 internal/model")),
		mcp.WithBoolean("cache", mcp.Description("是否启用缓存，可选，默认 true")),
	), handleGenerateModel(repoRoot))

	s.AddTool(mcp.NewTool("generate_api",
		mcp.WithDescription("封装 scripts/generate-api.sh：从 .api 文件生成 Handler/Logic 骨架。"+
			"注意：生成的 Types 定义在临时文件里，需要手动合并进 internal/types/types.go"),
		mcp.WithString("api_file", mcp.Required(), mcp.Description("API 文件路径，相对 api/ 或绝对路径")),
	), handleGenerateAPI(repoRoot))

	s.AddTool(mcp.NewTool("generate_ts",
		mcp.WithDescription("封装 scripts/generate-ts.sh：从 .api 文件生成前端 TypeScript 代码。"+
			"注意：禁止手改 generated/ 目录，如果路径包含 /auth 前缀需要在二次封装时修正"),
		mcp.WithString("api_file", mcp.Description("API 文件路径，可选，默认 admin.api")),
	), handleGenerateTS(repoRoot))

	s.AddTool(mcp.NewTool("generate_rpc",
		mcp.WithDescription("封装 scripts/generate-rpc.sh：从 .proto 文件生成 services/<service_name>/ 下的 pb/client/server/logic 骨架"),
		mcp.WithString("service_name", mcp.Required(), mcp.Description("服务名，如 iam/content/chat/task/sdk")),
		mcp.WithString("proto_file", mcp.Description("proto 文件路径，可选，默认 services/<service_name>/rpc/<service_name>.proto")),
	), handleGenerateRPC(repoRoot))

	s.AddTool(mcp.NewTool("generate_swagger",
		mcp.WithDescription("占位 stub：scripts/generate-swagger.sh 尚未实现，计划于 Phase 3"),
	), handleGenerateSwaggerStub)
}

func handleGenerateSQL(repoRoot string) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		group, err := req.RequireString("group")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		name, err := req.RequireString("name")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		args := []string{"-group", group, "-name", name}
		if v := req.GetString("parent_id", ""); v != "" {
			args = append(args, "-parent-id", v)
		}
		if v := req.GetString("parent_path", ""); v != "" {
			args = append(args, "-parent-path", v)
		}

		result, err := exec.RunWithAutoConfirm(repoRoot, filepath.Join(repoRoot, "scripts", "generate-sql.sh"), args, "db/services")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("执行 generate-sql.sh 失败", err), nil
		}
		return mcp.NewToolResultStructured(result, result.Stdout), nil
	}
}

func handleGenerateModel(repoRoot string) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		migrationFile, err := req.RequireString("migration_file")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		args := []string{migrationFile}
		if dir := req.GetString("dir", ""); dir != "" {
			args = append(args, "-d", dir)
		}
		// generate-model.sh 没有"关闭缓存"的开关，-c 本身就是脚本的默认行为；
		// cache 参数仅用于让调用方显式表达意图，这里始终追加 -c。
		args = append(args, "-c")

		result, err := exec.RunWithAutoConfirm(repoRoot, filepath.Join(repoRoot, "scripts", "generate-model.sh"), args, "internal/model")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("执行 generate-model.sh 失败", err), nil
		}
		return mcp.NewToolResultStructured(result, result.Stdout), nil
	}
}

func handleGenerateAPI(repoRoot string) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		apiFile, err := req.RequireString("api_file")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		result, err := exec.RunWithAutoConfirm(repoRoot, filepath.Join(repoRoot, "scripts", "generate-api.sh"), []string{apiFile}, "internal/handler internal/types")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("执行 generate-api.sh 失败", err), nil
		}
		note := "提醒：生成的 Types 定义在临时文件里，需要手动合并进 internal/types/types.go（generate-api.sh 自带的 usage 提示，容易被调用完就忘）。\n\n" + result.Stdout
		return mcp.NewToolResultStructured(result, note), nil
	}
}

func handleGenerateTS(repoRoot string) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var args []string
		if apiFile := req.GetString("api_file", ""); apiFile != "" {
			args = append(args, apiFile)
		}

		// generate_ts 产物在 admin-frontend/src/api/generated，是 admin-server/ 的同级目录，
		// 不是子目录：相对 repoRoot=admin-server/ 必须写成 ../admin-frontend/...（多一层 ..），
		// 直接写 admin-frontend/... 会被 git status 解析成 admin-server/admin-frontend/...，
		// 一个不存在的路径，generated_files 会一直是空数组。
		result, err := exec.RunWithAutoConfirm(repoRoot, filepath.Join(repoRoot, "scripts", "generate-ts.sh"), args, "../admin-frontend/src/api/generated")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("执行 generate-ts.sh 失败", err), nil
		}
		note := "提醒：禁止手改 generated/ 目录，如果路径包含 /auth 前缀需要在二次封装时修正（generate-ts.sh 自带的注意事项）。\n\n" + result.Stdout
		return mcp.NewToolResultStructured(result, note), nil
	}
}

func handleGenerateRPC(repoRoot string) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		serviceName, err := req.RequireString("service_name")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		args := []string{serviceName}
		if v := req.GetString("proto_file", ""); v != "" {
			args = append(args, "-f", v)
		}

		result, err := exec.RunWithAutoConfirm(repoRoot, filepath.Join(repoRoot, "scripts", "generate-rpc.sh"), args, filepath.Join("services", serviceName))
		if err != nil {
			return mcp.NewToolResultErrorFromErr("执行 generate-rpc.sh 失败", err), nil
		}
		return mcp.NewToolResultStructured(result, result.Stdout), nil
	}
}

func handleGenerateSwaggerStub(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return mcp.NewToolResultError(
		"generate_swagger 尚未实现：scripts/generate-swagger.sh 要到 Phase 3 才会创建，" +
			"设计见 admin-server/docs/20-api-docs-generation.md。当前调用是提前占位，不代表功能已就绪。",
	), nil
}
