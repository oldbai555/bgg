# 22 — 自建 `admin-mcp` 工具（Phase 1 Week 1 基础设施投入）

> 本文档是可直接执行的任务说明。这是 Part A.8 的落地细节，属于 **Phase 1 Week 1**（和 `registry.Transact`/`WithSession`、中间件收窄、密钥管理、CI/Docker 骨架同一批地基工作一起做），不是 Week 2 才开始的域改造。建好之后立刻在后续 13 周的开发里天天用，所以放最前面，越早可用收益越大。

## 0. 前置依赖

本文档**不依赖** `01-architecture-target.md`/`02-transactions-and-uow.md`/`03-wire-and-middleware.md` 的任何实现（那三份是单体内部代码改造，本文档是独立于 `internal/` 之外的配套工具），但依赖以下几份文档/文件作为**数据来源**（本工具只读它们，不修改它们）：

- [`AGENTS.md`](../../AGENTS.md) 第 3、5 节——`query_file_policy` 的规则来源。
- [`15-service-boundaries.md`](./15-service-boundaries.md) 第 1、2 节——`query_service_boundary` 的种子数据来源（本轮已随本文档集一起产出，直接读取即可，不存在"还没写完"的情况）。
- [`progress.md`](./progress.md)、[`14-production-deployment-checklist.md`](./14-production-deployment-checklist.md)——`query_progress`/`query_deployment_checklist` 的解析对象。
- `admin-server/scripts/generate-sql.sh`/`generate-model.sh`/`generate-api.sh`/`generate-ts.sh`——被封装的脚本，本工具不改这四个脚本本身。
- [`10-dev-execution-and-review-points.md`](./10-dev-execution-and-review-points.md)——"开发期 AI 可以直接执行 `generate-*.sh`"这条策略是本工具存在的前提；本工具不改变"谁能执行"这条政策，只是把"AI 自己拼命令行"变成"AI 调一个有明确 schema 的 tool"。

## 1. 目标与非目标

**目标**：一个跑在本机、通过 stdio 暴露给 Cursor/Claude Code 的项目专属 MCP server，覆盖计划文档 A.8 确认的三项能力（脚本封装、约定查询、进度查询）。

**非目标（明确不做，不是"以后再做"的隐含承诺）**：
- 不做鉴权、不做多用户——单机本地工具，和 `go-lsp`/`codegraph` 这些已接入的 MCP server 同等信任级别（谁能启动 Claude Code/Cursor，谁就能调这个工具）。
- 不替代 `scripts/generate-*.sh` 本身——生成逻辑仍然是这四个 shell 脚本 + `scripts/sqlgen`，本工具只是包一层结构化调用接口，脚本改了本工具的 wrapper 也要跟着改，不是反过来。
- 不追求覆盖 `AGENTS.md` 的每一条规则——第一版只覆盖计划文档点名的三项确认能力（脚本封装、`query_service_boundary`/`query_file_policy`/`query_middleware_order`、`query_progress`/`query_deployment_checklist`），后续真的需要别的查询能力再加，不要在 Week 1 一次性把 `AGENTS.md` 全部规则都编码进去。
- 不做代码生成之外的"智能"——不做自然语言到 SQL/Go 代码的生成，不做代码审查，纯粹是"脚本调用 + 结构化数据查询"两类工具。
- `generate_rpc`/`generate_swagger` 两个 tool 在 Phase 1 阶段只是**占位 stub**（对应脚本 `scripts/generate-rpc.sh`/`scripts/generate-swagger.sh` 要到 Phase 2/3 才会被 `16-rpc-conventions.md`/`20-api-docs-generation.md` 真正建出来）——调用时返回明确的"未实现，计划于 Phase 2/3，见 16/20 文档"结构化错误，不是让 MCP client 收到一个 404 式的黑盒失败。

## 2. 模块骨架（新建 `admin-server/tool/admin-mcp/`）

参考仓库里已有的 `admin-server/scripts/sqlgen/` 子模块模式（独立 `go.mod`，不并入主 module `postapocgame/admin-server`，避免这个纯开发期工具的依赖污染主二进制的 `go.sum`）：

```
admin-server/tool/admin-mcp/
├── go.mod                      # module admin-mcp，go 1.24（与主项目一致）
├── go.sum
├── main.go                     # 入口：构造 stdio MCP server，注册全部 tool，Serve()
├── build.sh                    # go build -o bin/admin-mcp .（参考 sqlgen/build.sh 的单行风格）
├── internal/
│   ├── exec/
│   │   └── script.go           # 共享的"跑 shell 脚本 + 自动确认 + git diff 收集产物列表"逻辑，三组 tool 共用
│   ├── scripttools/
│   │   └── generate.go         # 第 4 节：generate_sql / generate_model / generate_api / generate_ts / generate_rpc(stub) / generate_swagger(stub)
│   ├── conventiontools/
│   │   ├── service_boundary.go # query_service_boundary
│   │   ├── file_policy.go      # query_file_policy
│   │   └── middleware_order.go # query_middleware_order
│   └── progresstools/
│       ├── progress.go         # query_progress
│       └── checklist.go        # query_deployment_checklist
└── data/
    ├── service-boundaries.yaml # 第 5.1 节 schema，种子数据来自 doc 15
    └── file-policy-rules.yaml  # 第 5.2 节 schema，种子数据来自 AGENTS.md 第 3/5 节
```

`go.mod` 内容：

```go
module admin-mcp

go 1.24
```

## 3. MCP SDK 选择与 `main.go` 骨架

用 `github.com/mark3labs/mcp-go`（社区最成熟的 Go MCP SDK，stdio transport 开箱支持，`Tool`/`ToolHandler` API 简单，不需要自己实现 JSON-RPC 帧协议）。`go.mod` 加一条 require（具体版本执行时 `go get github.com/mark3labs/mcp-go@latest` 锁定，本文档不写死版本号）。另外需要 `gopkg.in/yaml.v3` 解析第 5 节的两份数据文件。

```go
// main.go
package main

import (
	"admin-mcp/internal/conventiontools"
	"admin-mcp/internal/progresstools"
	"admin-mcp/internal/scripttools"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	s := server.NewMCPServer("admin-mcp", "0.1.0")

	scripttools.Register(s)      // 第 4 节 6 个 tool
	conventiontools.Register(s)  // 第 5 节 3 个 tool
	progresstools.Register(s)    // 第 6 节 2 个 tool

	if err := server.ServeStdio(s); err != nil {
		panic(err)
	}
}
```

每个 `Register(s *server.MCPServer)` 函数内部用 `s.AddTool(mcp.NewTool(name, mcp.WithDescription(...), mcp.WithString("group", mcp.Required(), ...)), handlerFunc)` 的模式声明参数 schema——具体字段级参数见第 4/5/6 节各 tool 的入参表，写代码时照着表逐字段加 `mcp.WithString`/`mcp.WithBoolean`。

## 4. Tool 组 1：封装 `scripts/generate-*.sh`

### 4.1 共享执行逻辑（`internal/exec/script.go`）

四个已有脚本和后续两个都有同一个障碍：**它们全部用 `read -p "确认...? (y/N): "` 交互式确认**（`generate-sql.sh:125`、`generate-model.sh:158`、`generate-api.sh:133`、`generate-ts.sh:144`）。`exec.Command` 跑起来的子进程默认没有 TTY，`read -p` 会读到 EOF 直接判定为"未确认"并退出。本工具的运行前提是 `10-dev-execution-and-review-points.md` 已经明确"开发期 AI 可以直接执行 `generate-*.sh`"，所以这里的自动确认不是绕开政策，而是把已经拍板的政策落到代码里：

```go
package exec

import (
	"bytes"
	"os/exec"
)

type ScriptResult struct {
	Success        bool     `json:"success"`
	ExitCode       int      `json:"exit_code"`
	Stdout         string   `json:"stdout"`
	Stderr         string   `json:"stderr"`
	GeneratedFiles []string `json:"generated_files"`
}

// RunWithAutoConfirm 跑一个 generate-*.sh 脚本，向 stdin 喂 "y\n" 自动通过脚本自带的确认提示，
// 用 git status 前后对比推导出这次调用新增/修改了哪些文件（比逐脚本硬编码"输出文件名规律"更稳健，
// 脚本内部命名规则变了也不用同步改这里）。
func RunWithAutoConfirm(repoRoot, scriptPath string, args []string, watchDir string) (*ScriptResult, error) {
	before, err := gitStatusSnapshot(repoRoot, watchDir)
	if err != nil {
		return nil, err
	}

	cmd := exec.Command(scriptPath, args...)
	cmd.Dir = repoRoot
	cmd.Stdin = bytes.NewBufferString("y\n")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	runErr := cmd.Run()

	after, err := gitStatusSnapshot(repoRoot, watchDir)
	if err != nil {
		return nil, err
	}

	result := &ScriptResult{
		Stdout:         stdout.String(),
		Stderr:         stderr.String(),
		GeneratedFiles: diffFiles(before, after),
	}
	if exitErr, ok := runErr.(*exec.ExitError); ok {
		result.ExitCode = exitErr.ExitCode()
	} else if runErr == nil {
		result.ExitCode = 0
	} else {
		return nil, runErr // 脚本本身没找到/没有执行权限等，属于工具环境错误，不是业务失败
	}
	result.Success = result.ExitCode == 0
	return result, nil
}
```

`gitStatusSnapshot`/`diffFiles` 内部就是 `git status --porcelain=v1 -- <watchDir>` 跑两次做差集（新增的 `??` 条目 + 状态从无到有的 `M`/`A` 条目）。**`repoRoot` 固定传 `admin-server/`**（这个工具本身也在 `admin-server/tool/admin-mcp/` 下，`cmd.Dir` 设成 `admin-server/` 和 `scripts/generate-*.sh` 一贯的"在 `admin-server/` 下运行"约定一致），`watchDir` 因此必须是相对 `admin-server/` 的路径——`generate_sql`/`generate_model`/`generate_api` 三个的产物都在 `admin-server/` 内部，直接写 `db/services`/`internal/model`/`internal/handler internal/types` 没问题；**`generate_ts` 例外**：产物在 `admin-frontend/src/api/generated`，是 `admin-server/` 的**同级目录**，不是子目录，相对 `repoRoot=admin-server/` 必须写成 `../admin-frontend/src/api/generated`（多一层 `..`），直接写 `admin-frontend/src/api/generated` 会被 `git status` 解析成 `admin-server/admin-frontend/src/api/generated`，一个不存在的路径，`generated_files` 会一直是空数组。

**已知脚本 bug，不在本工具里掩盖**：`12-scripts-standardization.md` 已经点出 `generate-sql.sh` 里 `rm -f sqlgen` 之后立刻判断 `$?`，实际拿到的是 `rm` 的退出码而不是 `sqlgen` 程序本身的——这意味着 `ExitCode`/`Success` 字段在这一个脚本上不完全可信，**这是脚本本身的 bug，按 12 文档的任务修复，本工具不做绕过式的特殊处理**；`GeneratedFiles` 字段（基于 git diff，不依赖脚本退出码）在这个 bug 修复前是判断"这次调用到底生成没生成文件"更可靠的信号，调用方（AI）应该优先看 `generated_files` 是否非空，而不是只看 `success`。

### 4.2 六个 tool 的参数映射

| Tool 名 | 对应脚本 | 参数 | 说明 |
|---|---|---|---|
| `generate_sql` | `scripts/generate-sql.sh` | `group`（必填，`<domain>/<module>` snake_case，`<domain>` 决定落进哪个服务目录）、`name`（必填，中文功能名）、`parent_id`（可选）、`parent_path`（可选，默认 `/temp`） | 映射到 `-group -name -parent-id -parent-path`；`watchDir=db/services` |
| `generate_model` | `scripts/generate-model.sh` | `migration_file`（必填，相对 `db/services/<service>/<module>/` 或绝对路径）、`dir`（可选，默认 `internal/model`）、`cache`（可选 bool，默认 true） | 映射到位置参数 + `-d`/`-c`；`watchDir=internal/model` |
| `generate_api` | `scripts/generate-api.sh` | `api_file`（必填，相对 `api/` 或绝对路径） | 映射到位置参数；`watchDir="internal/handler internal/types"` |
| `generate_ts` | `scripts/generate-ts.sh` | `api_file`（可选，默认 `admin.api`） | 映射到位置参数；`watchDir=../admin-frontend/src/api/generated`（注意开头的 `../`，见上方说明） |
| `generate_rpc` | `scripts/generate-rpc.sh`（**尚不存在**） | 同上占位 | 见 4.3 节，返回 stub |
| `generate_swagger` | `scripts/generate-swagger.sh`（**尚不存在**） | 同上占位 | 见 4.3 节，返回 stub |

`generate_api` 调用完之后，工具返回里要在 `stdout` 之外额外拼一句提醒（不是脚本自带的，是本工具加的）："生成的 Types 定义在临时文件里，需要手动合并进 `internal/types/types.go`"——这是脚本自己在 usage 里写的注意事项（`generate-api.sh:60`），容易被 AI 调用完就忘。`generate_ts` 同理提醒"禁止手改 `generated/` 目录，如果路径包含 `/auth` 前缀需要在二次封装时修正"（`generate-ts.sh:166`）。

### 4.3 `generate_rpc`/`generate_swagger` 占位实现

```go
func generateRpcStub() *mcp.CallToolResult {
	return mcp.NewToolResultError(
		"generate_rpc 尚未实现：scripts/generate-rpc.sh 要到 Phase 2（约 Week 6-7，task-rpc 拆分时）才会创建，" +
			"设计见 admin-server/docs/16-rpc-conventions.md。当前调用是提前占位，不代表功能已就绪。",
	)
}
// generate_swagger 同构，指向 scripts/generate-swagger.sh + docs/20-api-docs-generation.md，Phase 3。
```

两个 stub 的 tool schema（参数列表）此时先按"预期未来会长成什么样"写（`generate_rpc` 大概率是 `service_name`/`proto_path` 之类，`generate_swagger` 大概率无参数），到 Phase 2/3 真正实现对应脚本时再核对调整，不要求现在就精确。

## 5. Tool 组 2：项目约定查询

### 5.1 `query_service_boundary(domain string)`

`data/service-boundaries.yaml` schema：

```yaml
services:
  - name: iam-rpc
    db_schema: admin_platform
    domains: [iam, system, monitoring, misc]
    table_count: 22   # 10 + 6 + 5 + 2，见 doc 15 第 2 节实测表清单
  - name: content-rpc
    db_schema: admin_content
    domains: [blog, video]
    table_count: 7    # 6 + 1
  - name: chat-rpc
    db_schema: admin_chat
    domains: [chat]
    table_count: 3
  - name: task-rpc
    db_schema: admin_task
    domains: [task]
    table_count: 1
  - name: sdk-rpc
    db_schema: admin_sdk
    domains: [sdk]
    table_count: 4
gateway:
  name: gateway
  note: HTTP 唯一入口，无状态，不持有自己的数据库，不在 domains 列表里
```

种子数据直接来自 [`15-service-boundaries.md`](./15-service-boundaries.md) 第 1 节的服务划分和第 2 节实测表清单（`iam-rpc` 的 22 张表 = iam 10 + system 6 + monitoring 5（含 15 文档发现的 `metric_daily_stats` 漏计项）+ misc 2）。**doc 15 是人读的权威版本，这份 YAML 是机器读的同步副本**，两者手动保持一致——15 文档如果后续因为拆分实际执行发现边界调整（例如某个域中途改了归属），要同步改这份 YAML，反过来也一样，改 YAML 前先确认 15 文档没有过时。

`query_service_boundary("blog")` 返回：
```json
{"service": "content-rpc", "db_schema": "admin_content", "sibling_domains": ["video"]}
```
`domain` 参数大小写不敏感、支持传入还没拆分前的旧 group 前缀（如 `blog_article` 也应该能匹配到 `blog`——匹配逻辑做前缀匹配，不要求精确相等，因为调用方经常是拿着一个 `.api` 里的 group 名字符串来查）。

### 5.2 `query_file_policy(path string)`

`data/file-policy-rules.yaml`——按 glob 规则从上到下匹配，第一条命中的规则生效（越靠前越具体）：

```yaml
rules:
  - glob: "internal/handler/**"
    exceptions: ["**/custom_routes.go"]
    policy: forbidden_generated
    note: "goctl 生成，禁止手改（custom_routes.go 除外）"
  - glob: "internal/model/**/*_gen.go"
    policy: forbidden_generated
    note: "goctl 生成，禁止手改"
  - glob: "internal/model/**/*.go"
    policy: editable_sibling
    note: "手改 sibling 文件（业务定制扩展 goctl 生成的 Model）"
  - glob: "admin-frontend/src/api/generated/**"
    policy: forbidden_generated
    note: "goctl 生成前端 TS，禁止手改，二次封装写在 src/api/*.ts"
  - glob: "internal/logic/**"
    policy: editable_generated_skeleton
    note: "goctl 生成骨架 + 手写业务逻辑，允许改方法体"
  - glob: "internal/repository/**"
    policy: editable_handwritten
    note: "手写，强制用 squirrel 构建 SQL，禁止 fmt.Sprintf 拼接"
  - glob: "internal/domain/**"
    policy: editable_handwritten
    note: "手写领域服务，Phase 1 起按 5 服务分组组织（iam 含 system/monitoring/misc，content 含 blog/video，chat/task/sdk 各自独立），不是按 9 个业务域各建一个包，见 01-architecture-target.md A.2"
  - glob: "api/admin.api"
    policy: editable_handwritten
    note: "唯一 .api 定义文件，手改后需重新执行 generate-api.sh/generate-ts.sh"
  - glob: "etc/admin-api.yaml"
    policy: editable_handwritten
    note: "开发环境配置，允许改；生产环境配置需要用户亲自确认，见 AGENTS.md 第 6 节"
  - glob: "/etc/work/*.json"
    policy: stop_and_ask
    note: "生产环境配置，AGENTS.md 明令必须停下来问用户，AI 不得直接改"
  - glob: "**/*"
    policy: normal_business_code
    note: "未命中以上任何规则，视为普通业务代码，正常手改"
```

`policy` 枚举值固定 5 种：`forbidden_generated`（禁止手改的生成产物）、`editable_sibling`（生成产物旁边的手改 sibling 文件）、`editable_generated_skeleton`（生成骨架+手写方法体混合文件）、`editable_handwritten`（纯手写）、`stop_and_ask`（AGENTS.md 第 6 节"必须停下来问用户"的路径）。查询逻辑：把输入 `path` 规范化成相对仓库根的路径后按顺序 glob 匹配，返回命中规则的 `policy`+`note`，找不到命中就落到最后的兜底规则。

### 5.3 `query_middleware_order()`

无参数，直接返回硬编码的固定结构（不读文件，因为这是 `AGENTS.md`/`10-go-code-style.mdc` 里明确写死不会变的顺序，编码进 Go 代码比再建一个只有一行数据的 YAML 更直接）：

```json
{
  "admin_order": ["Performance", "RateLimit", "Auth", "Permission", "OperationLog"],
  "sdk_order": ["SDKAuth", "SDKRateLimit", "SDKCallLog"],
  "public_order": ["Performance"],
  "notes": [
    "Permission 与 ApiEnabled 互斥，一个 .api 服务块只能选一个",
    "SDK 系列中间件不与 Admin 系列混用",
    "公开/无需登录接口只用 Performance（+ApiEnabled 视情况）"
  ]
}
```

## 6. Tool 组 3：进度查询

### 6.1 `query_progress()`

解析 [`progress.md`](./progress.md)：文件结构是 `## <日期>：<标题>` 一级标题分节，正文是自由格式 Markdown。解析逻辑：按 `^## ` 正则切分成条目，每条目提取 `date`（标题里 `YYYY-MM-DD` 部分）、`title`（标题剩余部分）、`body`（该标题到下一个 `## ` 之间的全部内容，原样返回，不做二次摘要——摘要交给调用方的 LLM 做，本工具只负责结构化切分）。返回按日期倒序排列的条目数组，调用方通常只关心最新一条，但全部返回，不在工具层截断。

```json
{
  "entries": [
    {"date": "2026-07-10", "title": "文档集编写完成（Phase 1-3 尚未开始实际代码改动）", "body": "..."}
  ]
}
```

### 6.2 `query_deployment_checklist(pending_only bool)`

解析 [`14-production-deployment-checklist.md`](./14-production-deployment-checklist.md)：文件结构是 `### <序号> · <标题>` 三级标题，正文固定含 `**触发条件**：`、`**部署时要做什么**：`、`**如何验证生效**：`、`**状态**：` 四个字段（14 文档第 9-21 行已经把这个格式写成强制约定）。解析逻辑：按 `^### ` 切分条目，每条目内用 `**触发条件**：(.*?)\n\n**部署时要做什么**` 这类跨字段的非贪婪正则依次抠出四个字段的文本（`状态` 字段进一步解析出 `TBD`/`已就绪，待执行`/`已执行` 三态之一，用于 `pending_only` 过滤）。

`pending_only=true` 时只返回状态不是 `已执行` 的条目（即 `TBD` + `已就绪，待执行`）；`pending_only=false` 返回全部条目。返回结构：

```json
{
  "entries": [
    {
      "index": 1,
      "title": "JWT 密钥改用环境变量注入",
      "trigger": "...",
      "action": "...",
      "verification": "...",
      "status": "已就绪，待执行"
    }
  ]
}
```

## 7. 注册进 `.mcp.json`

仓库现有流程（`AGENTS.md` 第 8 节、`script/sync_claude_mcp.sh`）是：维护者在仓库内提交 git 的 **`.cursor/mcp.json`**（Cursor 通道团队 SSOT，按项目按需加载）里增删 MCP server，然后跑 `make sync-claude-mcp-import` 把 Cursor 配置规范化路径后同步进项目内提交的 `.mcp.json`（Claude Code 读这份）。`admin-mcp` 走同一条路径，不新开一条注册机制：

1. 本地构建二进制：`cd admin-server/tool/admin-mcp && ./build.sh`（产出 `bin/admin-mcp`，参考 `sqlgen/build.sh` 的单行 `go build` 风格；`bin/` 目录加进 `.gitignore`，二进制不提交）。
2. 在仓库内 `.cursor/mcp.json` 的 `mcpServers` 里新增一条，**用 Cursor 的 `${workspaceFolder}` 变量指向仓库内的构建产物路径**（不要写死本机绝对路径，`import-cursor` 的 `normalize_cursor_mcp_json` 会把 `${workspaceFolder}` 整体替换成 `${CLAUDE_PROJECT_DIR:-.}`，这个替换是对文件原始文本做字符串替换、发生在 JSON 解析之前，所以直接写在 `command` 字段里也能被正确替换，不需要额外处理）：
   ```json
   "admin-mcp": {
     "command": "${workspaceFolder}/admin-server/tool/admin-mcp/bin/admin-mcp"
   }
   ```
3. 完全重启 Cursor 或执行 `/mcp reconnect all` 确认 Cursor 侧能连上。
4. 维护者跑 `make sync-claude-mcp-import`（内部是 `script/sync_claude_mcp.sh import-cursor`），确认输出里出现新的 `admin-mcp` 条目、且 `command` 字段已经变成 `${CLAUDE_PROJECT_DIR:-.}/admin-server/tool/admin-mcp/bin/admin-mcp`。
5. **commit `.cursor/mcp.json` + `.mcp.json`**（`AGENTS.md` 明确要求这一步，两者都是团队 SSOT）。
6. 在集成终端里 `claude mcp list` 确认 `admin-mcp` 出现且状态正常；`make sync-claude-mcp-check` 也会跑这条命令，可以复用做验证。

不需要 `env` 字段——`admin-mcp` 不像 `go-lsp`/`mongodb`/`redis` 那样需要指向外部资源的连接参数，它只在自己所在的仓库内跑 `git status`/`exec.Command` 调本仓库内的脚本，天然是仓库相对路径，不需要环境变量注入。

## 8. sqlmock/单元测试范围

本工具是纯基础设施代码，不接触业务表，**不适用**其余文档里"每个开 `Transact` 的方法配 sqlmock happy-path+rollback-path"的测试要求（这条约定是给 `admin-server` 主 module 的领域服务定的，`tool/admin-mcp` 是独立 module，不共享这条约束）。测试范围收窄到：

- `internal/exec/script.go` 的 `diffFiles`：纯函数，给定两个 `git status --porcelain` 快照字符串，断言差集计算正确（新增/修改/删除三种 case）。
- `internal/conventiontools/file_policy.go` 的 glob 匹配顺序：给几个真实路径（`internal/handler/iam/user/user_create_handler.go`、`internal/repository/iam/user_repository.go`、`/etc/work/mysql.json`）断言命中第 5.2 节对应的规则,不需要真的跑一个 MCP client 做端到端测试。
- `internal/progresstools/checklist.go` 的正则解析：拿 `14-production-deployment-checklist.md` 当前真实内容跑一次，断言至少解析出第 1 条（"JWT 密钥改用环境变量注入"）且 `status` 字段正确识别为 `已就绪，待执行`。

不测试 `scripttools` 六个 tool 的真实脚本执行路径（会真的跑 `generate-*.sh` 并对仓库产生副作用，不适合放进自动化测试）——这几个 tool 的验证方式是第 9 节的人工冒烟测试。

## 9. 完成的定义

1. `cd admin-server/tool/admin-mcp && go build ./...` 通过（独立 module，不影响主 `admin-server` 的 `go build ./...`）。
2. `go test ./...`（在 `tool/admin-mcp` 目录下）覆盖第 8 节列出的三类用例，全部通过。
3. `.mcp.json` 里出现 `admin-mcp` 条目，`git diff` 能看到这一处新增（按第 7 节步骤走完并 commit）。
4. 人工冒烟测试（在 Claude Code/Cursor 里实际调用，不是跑 Go 测试）：
   - 调 `query_middleware_order()`，确认返回的顺序和 `AGENTS.md` 第 3 节一致。
   - 调 `query_file_policy(path="internal/handler/iam/user/user_create_handler.go")`，确认返回 `forbidden_generated`。
   - 调 `query_service_boundary(domain="blog")`，确认返回 `content-rpc`。
   - 调 `query_progress()`，确认能读到 `progress.md` 当前最新一条记录的标题。
   - 调 `query_deployment_checklist(pending_only=true)`，确认返回条目数与手工数 `14-production-deployment-checklist.md` 里状态非"已执行"的条目数一致。
   - 在一个**测试用**的临时 `-group system/mcp_smoke_test -name MCP冒烟测试`上调一次 `generate_sql`，确认自动确认生效（没有卡在等待交互输入）、`generated_files` 字段里出现 `db/services/iam/mcp_smoke_test/create_table_mcp_smoke_test.sql`/`init_mcp_smoke_test.sql`，之后手动删除这两个测试文件（不要把冒烟测试产物留在仓库里）。
   - 调 `generate_rpc`，确认返回的是第 4.3 节描述的"未实现"结构化错误，而不是进程崩溃或空响应。
