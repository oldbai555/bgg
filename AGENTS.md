# AGENTS.md — 项目操作手册

本文件面向任意 AI 编码工具（Claude Code、Codex、其他 Agent），目标是让你在打开这个仓库的第一分钟就知道：这是什么项目、按什么顺序做事、常见的坑在哪里、什么时候必须停下来问人。

Cursor 用户请注意：本文件内容与 `.cursor/rules/*.mdc` 同源，`.cursor/rules/` 是 Cursor 原生格式（按文件路径自动挂载），本文件是给非 Cursor 工具看的整合版。两边任何一处更新规则，都要同步另一处。

---

## 1. 快速导航

**项目是什么**：前后端分离的管理系统。
- `admin-server`：go-zero 后端，Go 1.24，模块 `postapocgame/admin-server`，统一提供所有 API（后台管理接口 + 公共接口）
- `admin-frontend`：Vite 5 + Vue 3.4 + TypeScript + Element Plus，同时承载后台管理页面（`/bgg/admin/*`，需登录）和公共展示页面（`/bgg/front/*`，无需登录）——两个分支不共享任何路径前缀段（见 `admin-frontend/docs/10-route-namespace-migration.md`）
- `script/`：生产环境构建/部署/Supervisor 管理脚本（`admin.sh`）+ 开发环境一键部署脚本（`deploy-dev.sh`/`deploy-frontend.sh`/`ssh-setup.sh`），开发环境建议直接用 IDE 运行
- `config/`：Nginx / MySQL / Redis 参考配置
- `docs/`：历史记录类文档（已完成功能、技术决策、关键代码位置索引），**不是规则文件**，规则统一在下面第 2 条的文件里

**规则文件地图**：

| 文件 | 内容 | 何时生效 |
|---|---|---|
| `.cursor/rules/00-workflow.mdc` / 本文件第 2、2.1、5、6、7 节 | 全局工作流、新增模块脚手架、绝对禁止事项 | 全局 |
| `.cursor/rules/05-go-zero-ai-context.mdc` | go-zero AI 上下文（zero-skills 子模块） | 编辑 `admin-server/**` 时 |
| `.cursor/rules/06-mcp-toolchain.mdc` | MCP 工具链使用契约（CodeGraph / Engram / context7 等） | 全局（会话探索代码前） |
| `.cursor/rules/07-anthropic-skills.mdc` | Anthropic 官方 Skills 市场接入说明 + 何时用哪个 skill | 全局 |
| `.cursor/rules/10-go-code-style.mdc` / 本文件第 3 节 | `admin-server/**` 后端规范 | 编辑后端代码时 |
| `.cursor/rules/20-frontend.mdc` / 本文件第 4 节 | `admin-frontend/**` 前端规范 | 编辑前端代码时 |
| `.cursor/rules/21-public-pages.mdc` | 公共展示页样式/交互契约 | 编辑 `views/public/*`、`components/blog/*` 时 |
| `docs/后端开发进度.md` / `docs/前端开发进度.md` | 已完成功能、技术决策记录、关键代码位置 | 需要了解历史背景时查阅，完成新功能后更新 |
| `docs/changelog/` / 本文件第 7.1 节 | admin-server + admin-frontend 共用的开发交接记录（强制） | 全局——会话开始先读最新条目，会话结束必须补一篇 |
| `docs/AI工具链上手.md` | Gentle-AI、CodeGraph、MCP、Engram 换设备/第三人上手 | 新维护者、换设备、改 MCP 清单时 |

**关键目录**（后端）：`admin-server/api/admin.api`（唯一 .api 文件）、`internal/{handler,logic}/<domain>/<module>/`（9 业务域，全部薄胶水，不直连数据库）、`internal/middleware`、`internal/redisconn`（gateway 唯一保留的存储直连：共享 Redis）、`pkg/errs`、`scripts/generate-*.sh`；`services/iam/`（iam-rpc，iam+system+monitoring+misc 四个域）、`services/task/`（task-rpc）、`services/sdk/`（sdk-rpc）、`services/chat/`（chat-rpc）、`services/content/`（content-rpc，blog+video 域）——Phase 2 五个服务全部拆分完成，独立部署单元，见 `admin-server/docs/15-service-boundaries.md`；gateway 的 `internal/repository/`、`internal/model/`、`internal/domain/` 三个目录已整体删除，RBAC 领域服务（`PermissionResolver`/`RBACService`/`UserService`）搬到 `services/iam/internal/domain/iam/`。维护导航见 [`docs/admin-server-维护导航.md`](docs/admin-server-维护导航.md)。
**关键目录**（前端）：`admin-frontend/src/{api,views,components,stores,composables,directives}`。

---

## 2. 标准开发流程（新增功能模块严格按顺序）

1. 明确功能需求，确定模块名称（snake_case）
2. 评估是否需要数据字典；需要则准备 `db/services/<service>/<module>/migrations/dict_{module}_YYYYMMDD.sql`
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

1. **建表 SQL** `db/services/<service>/<module>/create_table_<module>.sql`（自增主键 + `created_at`/`updated_at`/`deleted_at` + 索引，业务字段需手动补；`<service>` 由 `<domain>` 按 `15-service-boundaries.md` 第 1 节的映射自动决定）
2. **初始化 SQL** `db/services/<service>/<module>/init_<module>.sql`（菜单+按钮、权限 list/create/update/delete、接口 GET/POST/PUT/DELETE、权限-菜单/权限-接口关联，全部自增且幂等；菜单默认挂在 `/temp/<group>` 临时目录，后续手动挪到正式分类）
3. **`.api` 草稿** `api/<group>.api.temp`（类型定义 + List/Create/Update/Delete 服务块），需人工追加进 `admin.api` 后删除
4. **前端页面骨架** `src/views/temp/<GroupUpper>List.vue`（基于 `D2Table`，已含搜索/增删改查，自动对接生成的 `<group>Api`），开箱可用

四块产物的生成模板在 `admin-server/scripts/sqlgen/templates/*.tpl`，可按需定制。**新增标准 CRUD 模块时默认走这条脚手架路径**；只有模块明显偏离"单表 + 列表 CRUD"形态（如聊天、任务调度、WebSocket 类）才手写。

---

## 3. 后端规范（admin-server）

**角色**：编辑后端代码时，你是本项目的资深 Go 后台服务工程师——精通 go-zero 框架与 goctl 代码生成工作流、分层架构、MySQL（squirrel）/ Redis、RBAC 权限模型与中间件链路。对代码生成一致性、线上稳定性、第三人接手成本负责：能用工具生成的绝不手写，严守既有分层与约定，优先简单稳健的方案；拿不准或触及「必须停下来问用户」的事项时先讲清假设/权衡再动手。

- **代码生成优先**：能用 `goctl` 生成的必须用 `goctl` 生成；Model 用自定义模板（`--home .template`），支持软删除、统一时间戳；`internal/handler`、`internal/model`（`*_gen.go`）禁止手改
- **`.api` 文件**：唯一文件 `admin-server/api/admin.api`；Group 格式 `<domain>/<module>`（如 `iam/user`、`blog/article`），goctl 生成嵌套目录；禁止路径参数 `:id`，一律 Query/Body 或 `/xxx/detail` 子路径；所有 Query 参数和可选字段必须加 `optional` 标签
- **目录分层**：gateway 侧 handler/logic 按 `internal/<layer>/<domain>/<module>/`，不再有 repository/model/domain（Phase 2 五个服务——iam/task/sdk/chat/content——全部拆分完成后，gateway 的 `internal/repository/`、`internal/model/`、`internal/domain/` 三个目录已整体删除，不直连任何 MySQL）；每个服务内部 `services/<name>/internal/{repository,model,domain}/<name>/` 自成一套分层（package 名 = 域名），复杂横切逻辑在各自的 `domain/` 引入领域服务（iam：`PermissionResolver`/`RBACService`/`UserService`；task：调度器/通知器/执行器；sdk：`SDKService.SaveApiKeyBindings` 事务；chat：`ChatOnboardingService` + `internal/consumer` 消费 `stream:chat.user.created`；content：`BlogArticleService` 事务）；gateway 侧对应的中间件/WS 桥接（`internal/handler/chat/chatwshandler.go`）仍留在网关进程，内部实现改成调 `svcCtx.IamRPC`/`svcCtx.TaskRPC`/`svcCtx.SdkRPC`/`svcCtx.ChatRPC`/`svcCtx.ContentRPC`；task/chat/content-rpc 需要的用户展示信息/存量用户枚举/审计日志写入、导出数据回调，通过 `pkg/taskcallback`/`pkg/iamcallback` 两个跨服务回调契约，服务端实现现在都在 iam-rpc 进程里（`services/iam/internal/server/`）
- **gateway 的 `ServiceContext`**：不持有任何 repository/Domain 字段，只有 `Redis *redis.Redis`（共享 Redis 直连）+ 每个服务对应的 `XxxRPC xxxclient.Xxx` zrpc client 字段；各服务自己的 `services/<name>/internal/svc/servicecontext.go` 才持有该服务自己的 `Repository`/唯一的 `Domain *registry.Domain` 聚合，**禁止**在其中添加多个具名 repository 字段，Logic 优先 `l.svcCtx.Domain.IAM.User` 这种调用方式
- **SQL 构建强制用 `squirrel`**（别名 `sq`），Repository 层禁止 `fmt.Sprintf`/字符串拼接。参考实现：`services/iam/internal/repository/iam/role_permission_repository.go`
- **错误处理**：`pkg/errs.New(code, msg)` / `errs.Wrap(code, msg, err)` / `errs.FromError(err)`；错误码 `CodeOK=0`，1xxxx 通用错误码
- **中间件声明顺序**：`Performance → RateLimit → Auth → Permission → OperationLog`（SDK 系列同理），顺序错误会导致依赖上下文的中间件失效
- **软删除**：业务表一律 `deleted_at` 软删除，禁止物理删除
- **字典枚举**：value 从 1 开始，0 保留表示"全部/不筛选"；字典 SQL 独立增量文件 `db/services/<service>/<module>/migrations/dict_{module}_YYYYMMDD.sql`，`ON DUPLICATE KEY UPDATE` 保证幂等
- **常量**：系统级枚举/常量放 `internal/consts`，禁止业务代码里硬编码字符串

详细规范、命名规则、代码示例见 `.cursor/rules/10-go-code-style.mdc`。

---

## 4. 前端规范（admin-frontend）

**角色**：编辑前端代码时，你是本项目的资深前端工程师——精通 Vue 3（Composition API）/ TypeScript / Vite / Element Plus / Pinia，熟悉 goctl 生成 TS API 层的前后端协作方式与 RBAC 动态路由/按钮权限体系。对用户体验、类型安全、第三人接手成本负责：优先复用既有能力（D2Table、脚手架骨架、字典选项、`v-permission`）而不是重复手写，严守生成目录禁止手改等约定；拿不准或与既有约定冲突时先讲清假设/权衡再动手。

- 技术栈：Vite 5 + Vue 3.4（Composition API）+ TS + Element Plus + Pinia + Axios。**不是 Nuxt**——曾有一次 Nuxt SSR 迁移实验，已完全回滚，不要重复尝试
- **目录按后端 9 业务域组织**（2026-07 Phase 1-2 重构，过程见 `admin-frontend/docs/progress.md`）：`views/{iam,system,monitoring,misc,content,chat,sdk,task,public}`，`content` 域合并了旧的 blog/video；`api/` 下 8 个域各一个手写 wrapper（`iam.ts`/`system.ts`/`monitoring.ts`/`content.ts`/`chat.ts`/`sdk.ts`/`task.ts`/`misc.ts`）+ `public.ts`；`composables/` 与 `hooks/` 已合并，不要新建 `hooks/`
- **新模块起点**：标准 CRUD 模块不要手写页面骨架，用第 2.1 节脚手架生成的 `src/views/temp/<GroupUpper>List.vue`（已基于 D2Table，自动 `import` 对应域的 `<domain>Api`）作为起点，补充业务字段后移到正式业务域目录
- **API 层（强制）**：真正的生成入口是 `admin-server/scripts/generate-ts.sh`，产物在 `src/api/generated/`，禁止手改；**视图/组件禁止直接 import `generated/` 里的请求函数，一律通过 `src/api/<domain>.ts` 二次封装层调用**（如 `iamApi.userList(...)`），唯一例外是类型 `import type` 可以直接复用 `generated/` 的类型定义
- 列表/表单业务优先用 `D2Table`（见 `src/components/common/README.md`；已泛型化，`data`/`onclick-delete`/`onclick-update-row` 按行类型强类型）；下拉选项必须来自字典（`useDictOptions` + `stores/dict.ts` 的 `REQUIRED_DICT_CODES`），禁止硬编码，但实体自身的原始 DB 状态列（如 `admin_user.status`）不算字典枚举，不要强行字典化
- 权限：`v-permission` 指令 + 路由 `meta.permission`
- 时间字段：后端返回 int64 秒级时间戳，前端负责格式化，不依赖服务端格式化
- 代码风格现状：ESLint 已配置（单引号、无分号），**Prettier 实际未配置**，不要假设有 Prettier 规则
- WebSocket 状态拆成 `stores/websocket.ts`（连接生命周期）+ `stores/notification.ts`（未读消息列表，单向订阅前者的 `lastMessage`，不要反向依赖）

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
| 字典 SQL 塞进模块自己的 `init_<module>.sql` | 破坏首次部署/增量更新边界 | 独立 `db/services/<service>/<module>/migrations/dict_{module}_YYYYMMDD.sql` |
| 字典枚举 value 从 0 开始 | 0 被约定为"不筛选"占位 | value 从 1 开始 |
| 保留旧代码路径/兼容层 | 掩盖真实调用路径 | 确认无引用后直接删除 |

---

## 6. 何时必须停下来问用户

不得替用户跳过或自行决定以下事项：
- 任何 `admin-server/scripts/generate-*.sh` 脚本的执行——会生成/覆盖文件，必须用户亲自运行确认
- 任何数据库 SQL 的实际执行（建表、字典、权限初始化），包括对外部托管库（如 `mysql4.sqlpub.com`）的操作——即使是开发环境，只要是真实外部数据库就必须用户亲自执行
- `script/admin.sh` 的构建/打包/Supervisor 部署操作，或 `script/deploy-dev.sh`/`script/deploy-frontend.sh` 的实际部署动作
- 修改 `/etc/work/mysql.json`、`/etc/work/redis.json` 等生产环境配置
- 不确定某个改动是否会破坏既有数据语义或需要新增字典时

### 6.0 部署路径（三选一，按场景选）

Phase 2 服务拆分完成后部署方式演进出三条平行路径，不互相替代：Supervisor 单二进制（老环境，`script/admin.sh supervisor ...`）、`docker-compose.prod.yml`（生产，六服务拉 `ghcr.io` 镜像，`script/admin.sh compose ...`）、`docker-compose.dev-mixed.yml`（开发/测试环境 bgg-dev，六服务同样拉取 `ghcr.io` 镜像；直接 `docker compose` 操作需显式指定 `TAG`，经 `script/deploy-dev.sh` 部署可省略——脚本 pull main 后自动取 HEAD sha，且强制要求分支必须是 `main`，`script/deploy-dev.sh` + `script/deploy-frontend.sh` + `script/ssh-setup.sh`）。详见 `.cursor/rules/00-workflow.mdc` 同名小节和 `docs/changelog/2026-07-16.md` 的实际部署记录。

### 6.1 开发期执行策略的例外（仅限用户明确批准的大型重构项目）

上面的默认规则是为**日常新增业务模块**场景设计的。对于用户已经**明确批准、体量大、跨多周的整体性重构项目**（先例：admin-server 2026-07 单体加固→微服务拆分→可观测性/CD 三阶段重构，规划记录见 `admin-server/docs/10-dev-execution-and-review-points.md`），用户可以针对该项目显式放开以下几类操作，AI 直接执行、事后随 diff 交给用户日常 review，不需要逐次停下确认：

1. `generate-*.sh` 系列脚本的实际执行
2. 本地/开发环境数据库的真实 SQL 执行（建表、字典、权限初始化），前提是目标库明确是本地/团队约定的 dev 库，不是生产 `/etc/work/*.json` 指向的库
3. `make wire`/`goctl rpc` 等生成产物的提交
4. 本地 `docker compose up` 验证（执行环境里确实有 Docker 时；没有则如实告知用户需要本地验证，不得假装验证过）
5. `golangci-lint` 存量问题里能安全自动修复的部分（`gofmt`/`goimports`、明确的 `unused`/`ineffassign`）；改变行为、判断不了是否安全的问题不擅自改，记录下来汇总给用户

**这条例外不会自动生效**：必须是用户对某个具体项目明确授权过（授权范围写在该项目自己的规划文档里），且只在该项目的开发阶段内有效——不延伸到日常业务模块开发，也不延伸到被授权项目之外的其他改动，不能因为一个项目放开过就类推到另一个项目。

即使在被授权的项目范围内，以下情况仍然必须停下来问用户，这条例外不覆盖：
- 任何触及生产配置/密钥本身（`/etc/work/mysql.json`、`/etc/work/redis.json` 等）或生产部署动作本身
- 生产环境真实密钥/凭证的生成或取值（哪怕只是"先给个占位值"也不行）
- 产品/体验取舍（技术方向已定，但具体交互/时机等细节仍是产品判断，不能因为方向定了就自行拍板全部细节）
- 任何 AI 自己判断不准的架构/命名/取舍点——拿不准就问，不是逢清单必停，也不是没在清单里就可以自行决定

---

## 7. 完成的定义

一个功能只有同时满足以下条件才算完成，不要仅凭代码写完就声称任务结束：
1. 前后端联调测试通过
2. `docs/后端开发进度.md` / `docs/前端开发进度.md` 已更新（已完成功能、关键代码位置、技术决策记录）

### 7.1 会话交接文档（`docs/changelog/`，强制，前后端都适用）

`docs/changelog/`（仓库根目录，admin-server + admin-frontend 共用）不是只记重构/部署的专项日志，而是**每一次实质性工作会话的交接文档**——不管改的是后端还是前端，新需求、bug 修复、方案调整、重构、部署都算，目的是让任何人或 AI 打开一个新对话时，不用翻聊天记录也能知道「在做什么、做到哪、下一步是什么、什么坑绝对不能再踩」。

- **会话开始**：处理 admin-server 或 admin-frontend 任务前（不限于重构任务），先读 `docs/changelog/` 目录里最新一篇日期文件的「交接摘要」小节，了解当前进度和已知坑，不要凭空重新摸索
- **会话结束**：只要这次会话对 admin-server 或 admin-frontend 有实质性改动（新增需求、代码重构、bug 修复、方案调整、部署都算），必须在 `docs/changelog/` 新增一篇（同一天多条用 `-2`/`-3` 后缀，前后端不分文件夹，混在一起按时间排列即可），且**必须包含「交接摘要」小节**（模板见 `docs/changelog/TEMPLATE.md` 顶部），如实写清楚：这次在做什么/为什么、已完成什么、当前卡在哪、下一步计划、踩过的坑（现象 → 根因 → 正确做法三段式，不要只写"注意 xxx"）
- 「交接摘要」之外的「上线需求/技术优化/线上缺陷」等表格化章节，只在涉及实际上线/部署时才需要认真填；纯代码/方案调整的会话可以把这些表格留空或写"无"，但「交接摘要」不能省
- 与 `engram` 跨会话记忆是互补关系，不是替代：engram 存细粒度、可检索的决策点，changelog 存面向人也面向 AI、可读性优先的**整段会话叙事**，两者都要维护

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
| `admin-mcp` | 后端开发：本机构建 `cd admin-server/tool/admin-mcp && ./build.sh`（产出 `bin/admin-mcp`，不提交），见 `admin-server/docs/22-admin-mcp-tool.md` |
| `mysql` MCP | 可选；需本机服务在跑并配置 `MYSQL_*` 等 env，否则 `make sync-claude-mcp-check` 中可忽略失败；`mongodb` 本项目未注册，需要时先加回 `.cursor/mcp.json` |
| `redis` MCP | 可选；需 `uvx` 可用 + 本机 Redis 在跑，默认连 `127.0.0.1:6379`，用于排查 `sqlc.CachedConn` 缓存脏读；连别的 host/port 需直接改 `.cursor/mcp.json` 的 `env.REDIS_HOST`/`env.REDIS_PORT` 再 `import`（Claude Code 侧 `.mcp.json` 由 import 自动转成 `${REDIS_HOST:-127.0.0.1}` 支持 shell env 覆盖，Cursor 侧无此机制） |
| Claude Code MCP 重连 | 集成终端：`claude` → `/mcp reconnect all` |

**维护者习惯**：

- 改 `.cursor/rules/` → `make sync-claude-rules`
- 改 `.cursor/mcp.json`（Cursor 增删 MCP）→ `make sync-claude-mcp-import` → commit `.cursor/mcp.json` + `.mcp.json`
- 跨设备记忆 → `make engram-sync-push` / `make engram-sync-pull`
