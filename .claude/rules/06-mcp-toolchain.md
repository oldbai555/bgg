# MCP 工具链（必读）

本项目配置了 **9 个 MCP server**（清单 SSOT：仓库根 `.mcp.json`；上手见 `docs/AI工具链上手.md`）。

| 通道 | MCP 配置来源 |
|------|----------------|
| **Cursor**（Chat / Agent / Composer） | `~/.cursor/mcp.json`（本机） |
| **Claude Code 插件** | `.mcp.json` + `.claude/settings.local.json` 的 `env` |

两通道 **server 名称一致**，但进程独立；终端 `claude mcp list` 全绿 ≠ 插件已连通。

## 会话开始（强制）

有 MCP 可用时，**第一条工具链动作**（在大量 Read/Grep 之前）：

1. **engram** `mem_current_project` → 无项目则 `mem_save` 登记
2. **engram** `mem_context` 或 `mem_search`（关键词：当前任务、模块名、近期决策）
3. 做出决策 / 修复 bug / 发现约定后 **立即** `mem_save`（不要等用户提醒）

## 按场景选 MCP（优先于裸 Read/Grep）

| 场景 | 优先 MCP | 工具 / 用法 |
|------|-----------|-------------|
| 理解代码结构、调用链、改代码前摸底 | **codegraph** | `codegraph_explore`（自然语言或符号名）；**一次调用**通常够，返回带行号的源码 + 调用路径 |
| 查第三方库 / 框架 API（Vue、Element Plus、go-zero 官方用法等） | **context7** | 先 `resolve-library-id`，再 `query-docs` |
| Go 符号定义、引用、重命名（`admin-server/**`） | **go-lsp** | `definition` / `references` / `hover`（编辑前确认签名与引用面） |
| Vue/TS 符号定义、引用、诊断（`admin-frontend/**`） | **vue-lsp** | `definition` / `references` / `diagnostics` / `hover` |
| 本项目 UI 组件（D2Table、layout、blog 组件） | **frontend-ui** | `ui_get_component` / `ui_list_components` / `ui_get_patterns` |
| go-zero 框架概念、goctl 参数、配置校验 | **mcp-zero** | `query_docs`、`validate_config`、`analyze_project` |
| 查本地 Mongo / Redis 运行数据（调试联调） | **mongodb** / **redis** | 仅在本机服务已启动时使用 |
| 跨会话决策、踩坑、团队约定 | **engram** | `mem_search` / `mem_save` / `mem_context` |

## 与本项目工作流的边界（硬约束）

| MCP 能做 | MCP **不能**替代 |
|----------|------------------|
| `codegraph_explore` 替代大范围 grep+Read | `admin-server/scripts/generate-*.sh`（**必须用户执行**） |
| `query_docs` / `mcp-zero` 查 go-zero 文档 | 直接 `goctl api go` 覆盖 `internal/handler`、`internal/model` |
| `mcp-zero` `analyze_project` 分析结构 | 用 mcp-zero 生成并提交与本仓库脚手架冲突的代码 |
| engram 记决策与约定 | 用 engram 代替更新 `docs/后端开发进度.md` |

新增标准 CRUD 模块仍走 `00-workflow.mdc` 的 `generate-sql.sh` 脚手架；mcp-zero 仅作**查阅与校验**，不绕过项目生成脚本。

## 反模式（降低 MCP 命中率时自查）

| 反模式 | 正确做法 |
|--------|----------|
| 一上来 `Grep` + `Read` 扫全仓库 | 先 `codegraph_explore` |
| 凭记忆写 Vue/go-zero API | 先 **context7** `query-docs` |
| 改 Go 函数不看引用 | 先 **go-lsp** `references` |
| 改 Vue/TS 符号不看引用 | 先 **vue-lsp** `references` |
| 手写 D2Table 用法不查文档 | 先 **frontend-ui** `ui_get_component` |
| 重复犯同类错误不记 engram | 修复后 `mem_save` |
| codegraph 已返回的源码再 `Read` 同一文件 | 直接基于 codegraph 返回的行号编辑 |
| MCP 报 not connected 仍假装能用 | 告知用户执行 `make sync-claude-mcp-check`，并降级用内置工具 |

## MCP 不可用时的降级

1. 告知用户哪个 server 未连接（`make sync-claude-mcp-check`）
2. **codegraph** 不可用 → `Grep` + `Read`（范围尽量小）
3. **context7** 不可用 → 查项目内 `docs/` 与 `.ai-context/zero-skills/`
4. **engram** 不可用 → 读 `docs/后端开发进度.md` 历史决策
5. **mongodb/redis** 不可用 → 不查库，改让用户确认或启动本地服务

## 维护者

Cursor 增删 MCP 后：`make sync-claude-mcp-import` → commit `.mcp.json`。改本规则后：`make sync-claude-rules`。
