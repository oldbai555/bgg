# AGENTS.md — 项目操作手册

本文件面向任意 AI 编码工具（Claude Code、Codex、其他 Agent），目标是让你在打开这个仓库的第一分钟就知道：这是什么项目、按什么顺序做事、常见的坑在哪里、什么时候必须停下来问人。

Cursor 用户请注意：本文件内容与 `.cursor/rules/*.mdc` 同源，`.cursor/rules/` 是 Cursor 原生格式（按文件路径自动挂载），本文件是给非 Cursor 工具看的整合版。两边任何一处更新规则，都要同步另一处。

---

## 1. 快速导航

**项目是什么**：前后端分离的管理系统。
- `admin-server`：go-zero 后端，Go 1.24，模块 `postapocgame/admin-server`，统一提供所有 API（后台管理接口 + 公共接口）
- `admin-frontend`：Vite 5 + Vue 3.4 + TypeScript + Element Plus，同时承载后台管理页面（`/admin/*`，需登录）和公共展示页面（`/blog/*`、`/videos/*`，无需登录）
- `script/`：生产环境构建/部署/Supervisor 管理脚本（`admin.sh`），开发环境建议直接用 IDE 运行
- `config/`：Nginx / MySQL / Redis 参考配置
- `docs/`：历史记录类文档（已完成功能、技术决策、关键代码位置索引），**不是规则文件**，规则统一在下面第 2 条的文件里

**规则文件地图**：

| 文件 | 内容 | 何时生效 |
|---|---|---|
| `.cursor/rules/00-workflow.mdc` / 本文件第 2、2.1、5、6、7 节 | 全局工作流、新增模块脚手架、绝对禁止事项 | 全局 |
| `.cursor/rules/05-go-zero-ai-context.mdc` | go-zero AI 上下文（zero-skills 子模块） | 编辑 `admin-server/**` 时 |
| `.cursor/rules/06-mcp-toolchain.mdc` | MCP 工具链使用契约（CodeGraph / Engram / context7 等） | 全局（会话探索代码前） |
| `.cursor/rules/10-go-code-style.mdc` / 本文件第 3 节 | `admin-server/**` 后端规范 | 编辑后端代码时 |
| `.cursor/rules/20-frontend.mdc` / 本文件第 4 节 | `admin-frontend/**` 前端规范 | 编辑前端代码时 |
| `.cursor/rules/21-public-pages.mdc` | 公共展示页样式/交互契约 | 编辑 `views/public/*`、`components/blog/*` 时 |
| `docs/后端开发进度.md` / `docs/前端开发进度.md` | 已完成功能、技术决策记录、关键代码位置 | 需要了解历史背景时查阅，完成新功能后更新 |
| `docs/AI工具链上手.md` | Gentle-AI、CodeGraph、MCP、Engram 换设备/第三人上手 | 新维护者、换设备、改 MCP 清单时 |

**关键目录**（后端）：`admin-server/api/admin.api`（唯一 .api 文件）、`internal/{handler,logic}/<domain>/<module>/`（9 业务域）、`internal/repository/<domain>/`、`internal/model/<domain>/`、`internal/domain/{iam,task}/`（领域服务）、`internal/middleware`、`pkg/errs`、`scripts/generate-*.sh`。维护导航见 [`docs/admin-server-维护导航.md`](docs/admin-server-维护导航.md)。
**关键目录**（前端）：`admin-frontend/src/{api,views,components,stores,composables,directives}`。

---

## 2. 标准开发流程（新增功能模块严格按顺序）

1. 明确功能需求，确定模块名称（snake_case）
2. 评估是否需要数据字典；需要则准备 `db/migrations/dict_{module}_YYYYMMDD.sql`
3. **用户**执行 `admin-server/scripts/generate-sql.sh -group <domain>/<module> -name <name>`（AI 不得代替执行）
4. 补齐 SQL 字段（`created_at`/`updated_at`/`deleted_at`，BIGINT 秒级时间戳）
5. 补齐 `.api` 接口参数和中间件声明（含 `optional` 标签，见第 3 节）
6. **用户**执行 `generate-model.sh <sql_file>`
7. **用户**执行 `generate-api.sh <api_file>`
8. 实现 Repository/Logic 业务逻辑
9. 执行 SQL：字典SQL → 业务表SQL → 权限SQL
10. 启动后端服务测试接口
11. **用户**执行 `generate-ts.sh`
12. 完善前端页面（基于生成的 `.vue` 骨架）
13. 前后端联调测试通过
14. 更新进度文档（`docs/后端开发进度.md` / `docs/前端开发进度.md`）

即：先把 admin-server 一侧（步骤 1-10）做完、验证通过，再推进 admin-frontend 一侧（步骤 11-13）。

### 2.1 新增模块脚手架（工程化能力，标准 CRUD 模块默认走这条路）

admin-server/admin-frontend 从设计之初就是"工程化"的：新增一个标准 CRUD 业务模块不需要从零手写，步骤 3 的 `generate-sql.sh -group <domain>/<module> -name <name>` 一条命令会同时生成四块互相匹配、可直接跑通的骨架：

1. **建表 SQL** `db/create_table_<group>.sql`（自增主键 + `created_at`/`updated_at`/`deleted_at` + 索引，业务字段需手动补）
2. **初始化 SQL** `db/init_<group>.sql`（菜单+按钮、权限 list/create/update/delete、接口 GET/POST/PUT/DELETE、权限-菜单/权限-接口关联，全部自增且幂等；菜单默认挂在 `/temp/<group>` 临时目录，后续手动挪到正式分类）
3. **`.api` 草稿** `api/<group>.api.temp`（类型定义 + List/Create/Update/Delete 服务块），需人工追加进 `admin.api` 后删除
4. **前端页面骨架** `src/views/temp/<GroupUpper>List.vue`（基于 `D2Table`，已含搜索/增删改查，自动对接生成的 `<group>Api`），开箱可用

四块产物的生成模板在 `admin-server/scripts/sqlgen/templates/*.tpl`，可按需定制。**新增标准 CRUD 模块时默认走这条脚手架路径**；只有模块明显偏离"单表 + 列表 CRUD"形态（如聊天、任务调度、WebSocket 类）才手写。

---

## 3. 后端规范（admin-server）

**角色**：编辑后端代码时，你是本项目的资深 Go 后台服务工程师——精通 go-zero 框架与 goctl 代码生成工作流、分层架构、MySQL（squirrel）/ Redis、RBAC 权限模型与中间件链路。对代码生成一致性、线上稳定性、第三人接手成本负责：能用工具生成的绝不手写，严守既有分层与约定，优先简单稳健的方案；拿不准或触及「必须停下来问用户」的事项时先讲清假设/权衡再动手。

- **代码生成优先**：能用 `goctl` 生成的必须用 `goctl` 生成；Model 用自定义模板（`--home .template`），支持软删除、统一时间戳；`internal/handler`、`internal/model`（`*_gen.go`）禁止手改
- **`.api` 文件**：唯一文件 `admin-server/api/admin.api`；Group 格式 `<domain>/<module>`（如 `iam/user`、`blog/article`），goctl 生成嵌套目录；禁止路径参数 `:id`，一律 Query/Body 或 `/xxx/detail` 子路径；所有 Query 参数和可选字段必须加 `optional` 标签
- **目录分层**：handler/logic 按 `internal/<layer>/<domain>/<module>/`；repository/model 按 `internal/<layer>/<domain>/`（package 名 = 域名）；复杂横切逻辑仅在 `internal/domain/iam/`（RBAC）、`internal/domain/task/`（调度）两处引入领域服务
- **ServiceContext**：保留 `Repository *repository.Repository` 与唯一的 `Domain *registry.Domain` 聚合；**禁止**添加多个具名 repository 字段。Logic 优先 `l.svcCtx.Domain.IAM.User`，旧代码可内联 `xxxrepo.NewXxxRepository(svcCtx.Repository)` 直至按域迁移
- **SQL 构建强制用 `squirrel`**（别名 `sq`），Repository 层禁止 `fmt.Sprintf`/字符串拼接。参考实现：`internal/repository/iam/role_permission_repository.go`
- **错误处理**：`pkg/errs.New(code, msg)` / `errs.Wrap(code, msg, err)` / `errs.FromError(err)`；错误码 `CodeOK=0`，1xxxx 通用错误码
- **中间件声明顺序**：`Performance → RateLimit → Auth → Permission → OperationLog`（SDK 系列同理），顺序错误会导致依赖上下文的中间件失效
- **软删除**：业务表一律 `deleted_at` 软删除，禁止物理删除
- **字典枚举**：value 从 1 开始，0 保留表示"全部/不筛选"；字典 SQL 独立增量文件 `db/migrations/dict_{module}_YYYYMMDD.sql`，`ON DUPLICATE KEY UPDATE` 保证幂等
- **常量**：系统级枚举/常量放 `internal/consts`，禁止业务代码里硬编码字符串

详细规范、命名规则、代码示例见 `.cursor/rules/10-go-code-style.mdc`。

---

## 4. 前端规范（admin-frontend）

**角色**：编辑前端代码时，你是本项目的资深前端工程师——精通 Vue 3（Composition API）/ TypeScript / Vite / Element Plus / Pinia，熟悉 goctl 生成 TS API 层的前后端协作方式与 RBAC 动态路由/按钮权限体系。对用户体验、类型安全、第三人接手成本负责：优先复用既有能力（D2Table、脚手架骨架、字典选项、`v-permission`）而不是重复手写，严守生成目录禁止手改等约定；拿不准或与既有约定冲突时先讲清假设/权衡再动手。

- 技术栈：Vite 5 + Vue 3.4（Composition API）+ TS + Element Plus + Pinia + Axios。**不是 Nuxt**——曾有一次 Nuxt SSR 迁移实验，已完全回滚，不要重复尝试
- **新模块起点**：标准 CRUD 模块不要手写页面骨架，用第 2.1 节脚手架生成的 `src/views/temp/<GroupUpper>List.vue`（已基于 D2Table）作为起点，补充业务字段后移到正式业务域目录
- **API 层**：真正的生成入口是 `admin-server/scripts/generate-ts.sh`，产物在 `src/api/generated/`，禁止手改；业务代码统一从 `src/api/*.ts` 的二次封装层导入。注意 `package.json` 里的 `api:gen` 脚本已失效（对应文件不存在），不要使用
- 列表/表单业务优先用 `D2Table`（见 `src/components/common/README.md`）；下拉选项必须来自字典（`useDictOptions` + `stores/dict.ts` 的 `REQUIRED_DICT_CODES`），禁止硬编码
- 权限：`v-permission` 指令 + 路由 `meta.permission`
- 时间字段：后端返回 int64 秒级时间戳，前端负责格式化，不依赖服务端格式化
- 代码风格现状：ESLint 已配置（单引号、无分号），**Prettier 实际未配置**，不要假设有 Prettier 规则

公共展示页（`views/public/*`、`components/blog/*`）的样式/交互契约见 `.cursor/rules/21-public-pages.mdc`：统一复用 `public-list.scss`/`public-detail.scss`、768px 断点、`sessionStorage` 状态恢复、`MetricReporter` 埋点、`IcpFooter`。

详细规范见 `.cursor/rules/20-frontend.mdc`。

---

## 5. 常见错误与修复

| 反例 | 为什么错 | 正确做法 |
|---|---|---|
| 跳过脚本步骤，手写生成类文件 | 破坏生成一致性，后续再生成会冲突 | 能生成的必须生成，AI 只补业务逻辑 |
| 改 `internal/handler`、`internal/model/*_gen.go`、`src/api/generated/*` | 下次重新生成会覆盖 | 只改 Logic/Repository/二次封装层 |
| Group 用驼峰命名 | 影响生成路径和路由前缀 | 必须 snake_case |
| `.api` 用路径参数 `:id` | 与现有全部路由风格不一致 | 用 Query/Body 或 `/xxx/detail` |
| Query/可选字段缺 `optional` 标签 | `httpx.Parse` 报 400 | 加 `,optional`，必填校验放 Logic 层 |
| Repository 用字符串拼 SQL | SQL 注入风险 | 必须用 `squirrel` |
| 业务表物理删除 | 破坏审计/软删除体系 | 一律 `deleted_at` 软删除 |
| 中间件顺序错 | 后置中间件依赖前置写入的上下文 | 按固定顺序声明 |
| 字典 SQL 塞进 `db/data.sql` | 破坏首次部署/增量更新边界 | 独立 `dict_{module}_YYYYMMDD.sql` |
| 字典枚举 value 从 0 开始 | 0 被约定为"不筛选"占位 | value 从 1 开始 |
| 保留旧代码路径/兼容层 | 掩盖真实调用路径 | 确认无引用后直接删除 |

---

## 6. 何时必须停下来问用户

不得替用户跳过或自行决定以下事项：
- 任何 `admin-server/scripts/generate-*.sh` 脚本的执行——会生成/覆盖文件，必须用户亲自运行确认
- 任何数据库 SQL 的实际执行（建表、字典、权限初始化）
- `script/admin.sh` 的构建/打包/Supervisor 部署操作
- 修改 `/etc/work/mysql.json`、`/etc/work/redis.json` 等生产环境配置
- 不确定某个改动是否会破坏既有数据语义或需要新增字典时

---

## 7. 完成的定义

一个功能只有同时满足以下条件才算完成，不要仅凭代码写完就声称任务结束：
1. 前后端联调测试通过
2. `docs/后端开发进度.md` / `docs/前端开发进度.md` 已更新（已完成功能、关键代码位置、技术决策记录）

---

## 8. AI 工具链（Cursor + Claude Code 插件）

详细步骤见 [`docs/AI工具链上手.md`](docs/AI工具链上手.md)。摘要：

**新设备 / 第三人**：`make setup-ai` → 完全重启 Cursor → 完成下方手动清单 → `make sync-claude-mcp-check`

**`make setup-ai` 不会自动完成、必须手动配置**：

| 项 | 操作 |
|---|---|
| `GO_ZERO_MCP_PATH` | 后端开发必需；指向本机 `go-zero-mcp` 可执行文件，写入 `~/.zshrc` |
| `go-lsp` | `go install github.com/isaacphi/mcp-language-server@latest` 且 `go install golang.org/x/tools/gopls@latest` |
| `vue-lsp` / `frontend-ui` | 前端开发：`cd admin-frontend && pnpm install`（`vue-lsp` 与 `go-lsp` 共用 `mcp-language-server`，LSP 为 `@vue/language-server`；`frontend-ui` 使用已提交的 `mcp/dist/`） |
| `mongodb` / `redis` MCP | 可选；需本机服务在跑，否则 `make sync-claude-mcp-check` 中可忽略失败 |
| Claude Code MCP 重连 | 集成终端：`claude` → `/mcp reconnect all` |

**维护者习惯**：

- 改 `.cursor/rules/` → `make sync-claude-rules`
- 在 Cursor 增删 MCP → `make sync-claude-mcp-import` → commit `.mcp.json`
- 跨设备记忆 → `make engram-sync-push` / `make engram-sync-pull`
