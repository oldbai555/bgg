# admin-server 重构进度记录

> 本文档是 `admin-server` 单体加固 → 微服务拆分 → 可观测性/CD 三阶段重构的**唯一进度记录**，贯穿 Phase 1-3，不分叉。风格与仓库根目录 `docs/后端开发进度.md` 一致：按时间顺序追加条目，记录做了什么、关键决策、关键文件位置。

> **维护方式：只追加，不重写。** 每次实际工作会话结束后，在文末新增一条日期条目，不要回头改写已有条目（除非是修正明确的事实错误）。历史决策即使后来被推翻，也保留原条目 + 新增一条说明推翻原因，不做静默删除。

---

## 2026-07-10：文档集编写完成（Phase 1-3 尚未开始实际代码改动）

本次会话产出的是**规划文档集**，不是代码改动——`admin-server/docs/` 下新增/补全 `00`~`22` 共 23 份任务书文档 + 本进度记录，`internal/`、`services/`、`cmd/` 等实际代码目录未被触碰。（修正：本条目最初写作"22 份"，会话中途用户追加了自建 `admin-mcp` 工具的需求，新增 `22-admin-mcp-tool.md`，最终是 23 份，此处直接订正而非另起条目，因为是同一次会话内的事实修正。）

**三阶段结构与约 14 周时间线概要**：

- **Phase 1（Week 1-5，单体内部加固）**：在现有单体进程内把事务方案（`registry.Transact` + `WithSession`）、领域服务分层（仅对约 35-40 个跨表/跨域/有复杂业务规则的文件引入 `internal/domain/<domain>`）、中间件 Wire 化、JWT 密钥外部化、测试策略（`sqlmock` 覆盖领域服务 happy-path + rollback-path）、CI/Docker 骨架落地。关键约束：`internal/domain/<domain>` 的包边界从一开始就按 Phase 2 最终的 5 个 RPC 服务分组组织（`iam` 吸收 `monitoring`/`system`/`misc`，`content` 吸收 `blog`/`video`），不是按现在的 9 个域机械对应，为 Phase 2 少走返工。Week 1 同时搭建自建 `admin-mcp` 工具（`admin-server/tool/admin-mcp/`，独立 Go module），封装 `generate-*.sh`、项目约定查询、进度查询三类能力，越早可用后续 13 周越受益。
- **Phase 2（Week 6-12，微服务拆分）**：把单体拆成 6 个独立部署单元——`gateway`（HTTP 唯一入口，无状态）+ `iam-rpc`/`content-rpc`/`chat-rpc`/`task-rpc`/`sdk-rpc` 五个 zrpc 服务，每个从第一天起有自己的 MySQL schema（`admin_platform`/`admin_content`/`admin_chat`/`admin_task`/`admin_sdk`）。拆分顺序按风险从低到高：`task-rpc` 先跑通全套机制 → `sdk-rpc` → `chat-rpc`（含 WS↔gRPC 流桥接）→ `content-rpc` → `iam-rpc` 放最后。跨服务一致性复用现有 Redis（Redis Streams 做尽力而为的异步副作用，task-rpc 升级为通用任务队列），不引入 Kafka/RabbitMQ；RBAC 校验方案是 gateway 同步调用 `iam-rpc.CheckPermission` 走缓存，不把权限塞进 JWT。
- **Phase 3（Week 13-14+，可观测性 + API 文档 + CD）**：六个服务接入 go-zero 原生 `Telemetry` 配置 + 新增 `pkg/logging` 统一 JSON 结构化日志（trace_id/span_id/service/user_id），解决拆分后"一个错误横跨四个进程、四份日志"的排查退化问题；新增 `scripts/generate-swagger.sh`（`goctl-swagger` 插件）从 `admin.api` 生成 `docs/openapi/admin-api.json`，5 个 RPC 服务不做单独文档工具，`.proto` 本身即文档；部署机制从"单二进制 + Supervisor"演进到"六个容器镜像 + docker-compose"（不上 Kubernetes），CI 自动构建推送镜像到 `ghcr.io`，生产环境的 `compose deploy` 保持人工触发。

**贯穿全程的原则**：项目未上线、无兼容性负担，可自由重新设计；开发期 `generate-*.sh` 执行、本地/开发库 SQL、`make wire`/`goctl` 生成产物提交等允许 AI 直接做、事后 review，只有真正触及生产部署动作、生产密钥、产品/体验取舍的地方才停下来问用户（完整规则见 `10-dev-execution-and-review-points.md`）。

**关键文件位置**：
- 文档集全貌与索引：`admin-server/docs/00-refactor-overview.md`
- 单体架构目标（Phase 1 正文）：`admin-server/docs/01-architecture-target.md`
- 服务边界设计（Phase 2 核心决策）：`admin-server/docs/15-service-boundaries.md`
- 可观测性/API 文档/CD（Phase 3，本次新增）：`admin-server/docs/19-observability.md`、`20-api-docs-generation.md`、`21-cd-and-deployment.md`
- 本轮之外明确不做的事项：`admin-server/docs/11-descoped.md`
- 持续追加的上线部署清单：`admin-server/docs/14-production-deployment-checklist.md`

**下一步**：用户过一遍 `00`/`01`/`15` 三份核心文档确认技术方案，再从 Phase 1 Week 1 开始实际动代码。本文件从下一次真实工作会话起追加新条目。

---

## 2026-07-10（续）：Phase 1 Week 1 地基工作全部落地

本次会话完成 Week 1 全部五项地基工作（对应 `02`、`03`、`01 A.5`、`09`、`22`），`go build ./...`、`go test ./...` 全绿，`admin-mcp` 独立 module 已构建并做过端到端 stdio 冒烟测试。

**1. 事务方案（`02-transactions-and-uow.md`）**：
- 37 个 Model 手改 sibling 文件（`internal/model/**`，`*_gen.go` 未触碰）全部新增 `WithSession(session sqlx.Session) XxxModel` 方法；批量部分用脚本机械生成，`admindicttypemodel.go`/`admindictitemmodel.go` 两个接口里带自定义方法的文件手工补的，全部核对过 `go build ./internal/model/...` 通过。
- `internal/repository/repository.go` 新增 `Repository.Transact`/`withSession`；`internal/repository/registry/domain.go` 新增 `registry.Transact`。
- 新增 `internal/repository/repository_test.go`：用 `sqlmock` + `miniredis`（后者用于满足 `NewRepository` 对非空 Redis 客户端/`CachedConn` 缓存节点的要求，`cache.New` 对空 `CacheConf` 会 `log.Fatal`，找到的可行组合是 miniredis 实例 + 单节点 `CacheConf`）验证 `Transact` 的 happy-path（提交）和 rollback-path（部分执行后回滚），两个测试都通过。
- `go.mod` 新增 `github.com/DATA-DOG/go-sqlmock`、`github.com/alicebob/miniredis/v2` 为直接依赖。

**2. 中间件收窄 + Wire（`03-wire-and-middleware.md`）**：
- 全部 11 个中间件构造函数按文档收窄：10 个吃 `*repository.Repository`（+ 需要的 `config.Config` 子结构），`PermissionMiddleware` 吃 `*iamdomain.PermissionResolver`。
- 配套改动：`registry.Domain.IAM` 新增 `PermissionResolver *iamdomain.PermissionResolver` 字段，`NewDomain` 里一次性构造（之前是每次请求在中间件里现场 `iamdomain.NewPermissionResolver(...)`，这是文档点名的真实 bug，现在全进程只构造一次）——这一小步提前借用了 `01-architecture-target.md` A.2 的"领域服务挂在 registry.Domain"模式，不等 Week 2 `04-domain-iam-chat.md` 才做，因为 `03` 文档明确说 `PermissionMiddleware` 的正确收窄依赖这个前提。
- `internal/wire/providers.go` 重写：删除 `buildMiddlewareBundle`，新增 `providePermissionMiddleware`（适配函数）+ `provideMiddlewareBundle`（assembler），`ProviderSet` 注册全部 11 个中间件构造函数。`make wire` 重新生成 `internal/wire/wire_gen.go` 成功。
- 人工冒烟：编译通过 + 启动到"读取 `/etc/work/mysql.json`"这一步（本机没有该文件，属预期，验证到 JWT fail-fast 检查已经通过、Wire 装配没有在启动阶段报错）。完整的登录/权限/限流真实请求级冒烟测试留给用户在有本地 MySQL/Redis 的环境里做。

**3. 密钥管理（`01-architecture-target.md` A.5）**：
- `etc/admin-api.yaml` 的 `JWT.AccessSecret`/`RefreshSecret` 改成 `${JWT_ACCESS_SECRET}`/`${JWT_REFRESH_SECRET}`；`admin.go` 的 `conf.MustLoad` 加 `conf.UseEnv()`，紧跟着加空值 fail-fast 检查。
- 实测验证：不设环境变量 → 进程在读完配置后立即 `log.Fatalf` 退出，不往下走；设置本地开发值后 → 通过 JWT 检查，继续往下走到 MySQL 外部配置加载这一步（符合预期，本机没有 `/etc/work/mysql.json`）。
- `docs/14-production-deployment-checklist.md` 条目 1 状态从 `TBD` 更新为 `已就绪，待执行`。

**4. CI/Docker 骨架（`09-ci-cd-and-deployability.md`）**：
- `goctl docker -go admin.go` 生成 `admin-server/Dockerfile`，人工核对/修正：Go 版本对齐 `go.mod` 的 `1.24`（`golang:1.24-alpine`）、`BaseImage` 从生成默认的 `scratch` 改成文档要求的 `alpine:latest`、补上 `EXPOSE 20000`（对齐 `etc/admin-api.yaml` 的 `Port`）。
- 新增 `admin-server/.dockerignore`、`admin-server/docker-compose.yml`（MySQL + Redis + app，健康检查齐全）、`admin-server/.golangci.yml`（范围收窄到 `internal/repository`/`internal/domain`，`golangci-lint run` 已实测跑通，发现一个跟本轮改动无关的既有问题：`internal/repository/chat/chat_repository.go:88` 的 `ineffassign`，留给该文件后续迁移 squirrel 时顺手清理，不在本轮处理）、仓库根 `.github/workflows/ci.yml`（三个 job：lint-build/unit-test/integration-test）。
- **数据库初始化顺序问题（文档已预见、本次实测确认）**：`db/` 目录下 `migrations/` 是子目录，MySQL 官方镜像的 `docker-entrypoint-initdb.d` 只处理挂载目录**顶层**的 `.sql`/`.sh` 文件，不递归子目录，直接整个挂载 `db/` 会导致 `migrations/*.sql` 全部不被执行。解决方式：新增 `admin-server/db/docker-init.sh`，显式按 `tables.sql → migrations/create_table_*.sql → data.sql → migrations/dict_*.sql → migrations/init_*.sql` 顺序逐个跑；`docker-compose.yml` 把整个 `db/` 挂到 `/db`（不是 `/docker-entrypoint-initdb.d`），只把这一个脚本单独挂到 `/docker-entrypoint-initdb.d/00-init.sh`。
- 本机没有安装 Docker，`docker build`/`docker compose up` 未实际验证，需要用户在有 Docker 的环境里跑一遍确认。

**5. `admin-mcp` 工具（`22-admin-mcp-tool.md`）**：
- 新建 `admin-server/tool/admin-mcp/`（独立 module，`go.mod` 声明 `go 1.25.5`——`go get github.com/mark3labs/mcp-go@latest` 自动把 toolchain 提到这个版本，文档写的是"与主项目一致的 1.24"，实测该 SDK 版本的间接依赖需要更高版本，是本工具作为独立 module 的合理妥协，不影响主 `admin-server` 的 `go 1.24.0`）。
- 三组共 8 个 tool 全部落地：`generate_sql`/`generate_model`/`generate_api`/`generate_ts`（真实封装）+ `generate_rpc`/`generate_swagger`（占位 stub）；`query_service_boundary`/`query_file_policy`/`query_middleware_order`；`query_progress`/`query_deployment_checklist`。
- `internal/exec/script.go` 的自动确认机制（喂 `y\n` 给 stdin）+ git status 前后 diff 推导产物列表，端到端实测跑通（见下）。
- 测试覆盖 `diffFiles`（新增/无变化/rename 三种 case）、`query_file_policy` 的 glob 匹配顺序（9 个真实路径 case，含 `exceptions` 命中后正确 fallback 到兜底规则）、`query_deployment_checklist` 的正则解析（对着本文件当前真实内容跑，断言第 1 条标题/状态正确解析），全部通过。
- **端到端 stdio 冒烟测试（不是走 Go test，是真的起进程发 JSON-RPC）**：`query_middleware_order`/`query_file_policy`/`query_service_boundary`/`query_progress`/`query_deployment_checklist`/`generate_rpc` 六个 tool 逐一调用验证返回符合预期；`generate_sql` 用临时 group `mcp_smoke_test` 真实跑通一次（自动确认生效、`generated_files` 正确列出新增文件），跑完手动删除了生成的 4 个文件（含前端 `admin-frontend/src/views/temp/McpSmokeTestList.vue`），未留痕迹。
- **意外发现并顺手修复**：`admin-server/scripts/*.sh`（`generate-sql.sh`/`generate-model.sh`/`generate-api.sh`/`generate-ts.sh`/`migrate-menu.sh`）在当前工作区里都没有可执行权限（`-rw-r--r--`），但 `scripts/README.md` 文档的用法是 `./scripts/generate-sql.sh ...` 直接执行——这会导致 `admin-mcp` 的 `exec.Command` 直接调用失败（`permission denied`），也会导致用户本地直接 `./scripts/xxx.sh` 一样失败。已执行 `chmod +x` 恢复全部 5 个脚本的可执行位，几个文件在 git diff 里只有权限位变化。
- **意外发现，未处理（记录留给后续）**：`generate-sql.sh` 的 `-group` 参数实测**不支持** `<domain>/<module>` 斜杠格式（`00-workflow.md`/`10-go-code-style.md` 描述的 `<domain>/<module>` 格式实际是 `.api` 文件 `@server(group:...)` 字段的格式，和这个脚本的 `-group` CLI 参数是两个不同的"group"概念）——传 `temp/mcp_smoke_test` 这种带斜杠的值会因为中间目录不存在而在建表 SQL 生成阶段失败（`open .../create_table_temp/mcp_smoke_test.sql: no such file or directory`）。现有的 `create_table_*.sql` 文件命名（`blog`/`blog_extension`/`metric`/`task`）也印证了实际约定是扁平 snake_case，不带斜杠。`admin-mcp` 的 `generate_sql` tool 描述文本目前原样照抄了 `22-admin-mcp-tool.md` 里"`<domain>/<module>` snake_case"的参数说明，跟脚本真实行为不一致，需要在后续会话核实到底是文档表述有误还是脚本该支持嵌套，再决定改文档还是改脚本（本次不擅自改动 `generate-sql.sh` 本身）。
- **未完成的收尾步骤（需要用户执行）**：把 `admin-mcp` 注册进 `~/.cursor/mcp.json` 属于用户本机 Cursor 配置，AI 不代为修改；用户完成注册后跑 `make sync-claude-mcp-import` 并 commit `.mcp.json`（`06-mcp-toolchain.md` 规定的标准流程）。

**遗留/需要用户关注的点**：
1. `internal/repository/chat/chat_repository.go`、`internal/repository/performance_log_repository.go` 仍未迁移到 squirrel（`10-go-code-style.md` 里的已知例外），本轮未处理。
2. Docker/compose 未在本机实际跑通（环境无 Docker），需要用户验证。
3. `generate-sql.sh` 的 `-group` 斜杠格式问题（见上）需要下次会话确认口径。
4. Week 2（`04-domain-iam-chat.md` + `05-domain-task.md`：IAM+Chat 联合改造，含 `user_create_logic.go` 修复；Task 域补测试）尚未开始。

**下一步**：进入 Phase 1 Week 2，按 `04-domain-iam-chat.md` 做 IAM+Chat 联合改造（`UserDomainService.CreateUser` + `Repository.Transact`、`FindPage` 改 `FindChunk`、chat onboarding 异步化——注意 chat onboarding 具体触发时机等体验细节按 `10-dev-execution-and-review-points.md` 第 2 节第 1 条仍需跟用户确认一次，不能直接照单实现），随后 `05-domain-task.md` 补 Task 域测试。

---

## 2026-07-10（续二）：提交前 Cursor 自动代码审查（Gentleman Guardian Angel）发现问题修复

Git 提交前钩子跑了一次基于 Cursor 的自动代码审查（对照 `AGENTS.md`/`.cursor/rules`），第一次提交被拦截（`CODE REVIEW FAILED`）。逐项核实后处理如下：

**已修复（判断为真问题）**：
1. `apienabledmiddleware.go`/`sdkauthmiddleware.go` 三处 `Status != 1` 硬编码改为 `Status != consts.Open`（`admin.go` 的 `syncRoutesToAdminAPI` 已经在用 `consts.Open`，属于既有约定，这几处中间件之前没跟上，不是本轮改造引入的新代码但顺手修了）。
2. `ratelimitmiddleware.go` 的限流响应把 HTTP 状态码 `http.StatusTooManyRequests`（429）误当业务错误码传给 `errs.New`——新增 `pkg/errs.CodeTooManyRequests = 10009` 并改用它；`w.WriteHeader(http.StatusTooManyRequests)` 保留（`response.ErrorCtx` 对业务错误统一写 400，靠这一行手动写在前面让客户端仍然收到 429，Go `net/http` 对同一响应重复 `WriteHeader` 只认第一次）。
3. `.github/workflows/ci.yml` 集成测试库初始化只跑了 `tables.sql`，没跟上 `db/docker-init.sh` 的完整顺序（建表→增量建表→基础数据→字典增量→模块初始化）——改成内联同样顺序的循环。
4. `docker-compose.yml` 的 `app` 服务补充说明：`admin.go` 强制要求 `-mysql-config`/`-redis-config` 指向的 JSON 文件存在（`MergeExternalConfig` 找不到直接报错退出，本次 JWT fail-fast 冒烟测试已经实测到这一步），新增 `etc/dev-mysql.json`/`etc/dev-redis.json`（仅供本地 compose 使用，host 指向 compose 内部服务名）+ `command:` 覆盖传入这两个 flag。
5. `tool/admin-mcp` 的 `generate_sql` 工具 `group` 参数描述文本改正，去掉之前照抄 `22-admin-mcp-tool.md` 的 `<domain>/<module>` 措辞（上一条目已经发现这是和 `.api` group 字段不同的概念，这次把工具自身的 schema 描述也同步改掉，不让调用方被误导）。

**核实后判断为超出本轮范围，未修复（保留分歧记录）**：
1. 审查要求把 `internal/model/system/admindictitemmodel.go` 的 `fmt.Sprintf` 动态 SQL 迁移到 squirrel。核实 `.cursor/rules/10-go-code-style.mdc`「SQL 构建规范」原文：「适用范围：**Repository 层**所有自定义查询方法」，且「已知例外」清单点名的是两个 `internal/repository/**` 文件，不含任何 `internal/model/**` 文件——规则文本本身没有覆盖 Model 层手改 sibling 文件。这个文件本次会话只新增了 `WithSession`，`FindPageByTypeId` 的实现是历史遗留代码，不在 Week 1 任何一份任务书文档的范围内，贸然改写有引入回归的风险，本轮不做，留给后续单独评估。
2. 审查要求本次改动同步更新根目录 `docs/后端开发进度.md`。核实 `00-refactor-overview.md` 第 4 节原文：「本轮重构过程中只更新 `progress.md`；但每当 Phase 1-3 里的某个改造让某个业务模块的功能行为发生了实质变化...同步更新 `docs/后端开发进度.md`」——Week 1 的五项工作（事务机制、中间件依赖收窄、密钥读取方式、CI/Docker、内部工具）都是纯内部重构，对外 API 行为、请求/响应结构均无变化，按这条既有规则不触发 `docs/后端开发进度.md` 更新。
3. 审查提到的硬编码 CORS 路径、SDK 限流 Redis 键前缀、操作日志类型字符串（`"create"`/`"update"` 等）——均为逐字保留的历史代码，本轮中间件收窄只改了构造函数签名和字段访问，不改这些方法体内部逻辑，视为与本次改造无关的既有代码风格问题，不顺手夹带修改。

---

## 2026-07-11：Phase 1 Week 2（`04-domain-iam-chat.md` + `05-domain-task.md`）全部完成

**产品/体验取舍确认**：`UserDomainService.CreateUser` 建完用户后触发 Chat onboarding 的方式，用户选择方案 A——进程内 `goroutine` 直调（维持"建用户请求返回后立刻后台触发、对用户完全不可见、不产生任务记录"的现有语义），不是 `internal/domain/task` 调度器派发方案。按此实现，未使用 `04` 文档里"可选"的调度器派发路径。

**1. IAM 域（`04-domain-iam-chat.md` 任务 1/2/3/4/5，任务 3/4 在 Week 1 已顺带完成，本轮补齐 1/2/5）**：
- 新增 `internal/domain/iam/user_service.go`（`UserDomainService.CreateUser`）：用户名唯一性校验 + 密码加密 + 落库包在 `Repository.Transact` 里，成功后 `go func` 异步触发 `chatdomain.Onboarding.InitNewUser`（`recover()` 兜底，失败只记日志）。
- 新增 `internal/domain/chat/onboarding.go`（`package chat`）：`Onboarding`/`UserLister` 两个窄接口 + `ChatOnboardingService`（`joinDefaultGroup` 加默认群、`createPrivateChatsForExistingUsers` 用 `UserLister.FindChunk` 分批建历史私聊、`createPrivateChat` 三条写（1 个 chat + 2 个 chat_user）包在一个 `Transact` 里）。IAM→Chat、Chat→IAM 两个方向都不再互相 `import internal/repository/<domain>`，全部通过窄接口 + `registry.NewDomain`（组合根）里的适配器（`iamUserListerAdapter`）桥接。
- 新增 `internal/domain/iam/rbac_service.go`（`RBACService`）：`UpdateRolePermissions`/`UpdateUserRoles`/`UpdatePermissionMenus`/`UpdatePermissionApis` 四个方法，把原来 4 个 logic 文件里"先删后插两步裸写、无事务保护"的逻辑收进领域服务、包一层 `Transact`；4 个 logic 文件（`role_permission_update_logic.go`/`user_role_update_logic.go`/`permission_menu_update_logic.go`/`permission_api_update_logic.go`）改为薄委托 `l.svcCtx.Domain.IAM.RBAC.UpdateXxx(...)`，缓存失效的 `go func(){...}` 尽力而为写法保留在 logic 层不变。
- `internal/repository/registry/domain.go`：`IAMDomain` 新增 `UserService`/`RBAC` 字段，`ChatDomain` 新增 `Onboarding` 字段；`NewDomain` 里 `chatOnboarding` 提到局部变量先构造，再分别注入 `IAM.UserService` 和 `Chat.Onboarding`，避免出现两个不同实例。
- `user_create_logic.go`（`internal/logic/iam/user/user_create_logic.go`）改为薄委托，原先手写的 `initChatForNewUser`（IAM 直接 `import internal/repository/chat` 的越界代码、`FindPage(1,10000,"")` 全表扫描存量用户）整体删除。
- **顺手修复的真实 bug（不在任务书范围内，但直接阻塞本轮功能正确性）**：`internal/repository/iam/user_repository.go` 的 `Create` 原来丢弃了 `model.Insert` 返回的 `sql.Result`，导致调用方拿到的 `user.Id` 恒为 `0`——`UserDomainService.CreateUser` 依赖正确的 `user.Id` 才能把新用户 ID 传给 Chat onboarding，这个 bug 在重构前就存在（原 `initChatForNewUser` 同样吃的是恒为 0 的 `user.Id`，只是从未被测试或人工验证过），本轮顺手补上 `LastInsertId()` 赋值。

**2. Task 域（`05-domain-task.md`）**：审计结论与任务书一致——`scheduler.go`/`notifier.go`/`excel_export_executor.go` 三个文件均不需要 `Transact` 改造（状态迁移是独立操作、通知失败尽力而为、文件系统失败已有补偿删除）。唯一代码改动是任务书里标记为"可选"的清理：`scanAsyncTasks`/`scanScheduledTasks` 删掉未使用的 `taskRepo taskrepo.TaskRepository` 参数（连带清理 `scanAndExecute` 里的调用点）。

**3. 测试（`08-testing-strategy.md` 约定：`sqlmock` 测 SQL/事务语义，`mockery` 测跨接口边界协作对象）**：
- 新增依赖：`mockery` CLI（`go install github.com/vektra/mockery/v2@latest`，已装在 `~/go/bin`，不进 `go.mod`）；`go.mod` 新增间接依赖 `github.com/stretchr/objx`（`mockery` 生成物依赖 `testify/mock` 带来的）。
- `chatdomain.Onboarding`/`chatdomain.UserLister`、`interfaces.TaskExecutor` 三个接口顶部加 `//go:generate mockery` 指令，`go generate` 产物落 `internal/mocks/chat/`、`internal/mocks/interfaces/`（按生成代码规则处理，不手改）。
- `internal/domain/iam/user_service_test.go`：`CreateUser` happy path（含等待异步 onboarding 触发、断言 `user.Id` 正确）+ 两个 rollback path（用户名已存在 / INSERT 失败）。
- `internal/domain/iam/rbac_service_test.go`：`UpdateRolePermissions` happy path + 两个 rollback path（角色不存在 / INSERT 失败）。
- `internal/domain/chat/onboarding_test.go`（`package chat` 白盒）：`createPrivateChat` happy/rollback path；`internal/domain/chat/onboarding_pagination_test.go`（`package chat_test` 外部测试包，避免 `chat` ↔ `internal/mocks/chat` 的 test-only import cycle）：`InitNewUser` 分页边界（page1 恰好 100 条触发第二次 `FindChunk`，page2 不足 100 条提前结束），用 `mockery` 生成的 `UserLister` mock 断言调用次数和 `lastID` 参数。
- `internal/repository/registry/domain_test.go`：`iamUserListerAdapter.FindChunk` 过滤禁用/已删除用户的单测。
- `internal/domain/task/scheduler_test.go`：`scanAsyncTasks`/`scanScheduledTasks` 的 WHERE 条件断言 + `executeTask` 的 4 个分支（成功/锁已被持有/执行器未注册/执行器 panic）；`internal/domain/task/notifier_test.go`：三种状态通知 + 通知落库失败不影响主流程 + `chatHub` 为 nil 不 panic；`internal/domain/task/executors/excel_export_executor_test.go`：不支持的导出模块 + 一次完整导出成功路径（含读取真实落盘 CSV 内容校验）+ `generateCSVFile` 的 DB 写入失败补偿删除 + 命中已有文件记录复用。
- 全部新增测试 `go test ./internal/domain/... ./internal/repository/...` 通过；`go build ./...`、`go vet ./...` 全仓库无告警。

**4. 完成的定义核对（对照 `04` 文档第 7 节的可执行验收标准）**：
- `grep -rn "internal/repository/chat" internal/logic/iam`：无匹配。
- `grep -n "iamdomain.NewPermissionResolver" internal/middleware/permissionmiddleware.go`：无匹配（Week 1 已完成，本轮复核仍成立）。
- 人工冒烟测试（建用户接口几十毫秒内返回、权限中间件放行/拒绝、角色权限保存失败后原权限不丢）留给用户在有本地 MySQL/Redis 的环境里做，本轮未执行（无本地数据库）。

**遗留/需要用户关注的点**：
1. `06-domain-blog-video-sdk.md`/`07-domain-monitoring-system-misc.md`（Phase 1 Week 3）尚未开始。
2. `04` 文档任务 6（`AuthDomainService.Login` 下沉，可选、优先级低于本轮已完成项）未做，按文档建议可以推迟到 Week 4-5 测试覆盖补漏阶段。
3. 人工冒烟测试（建用户/权限校验/角色权限回滚）需要用户在有本地 MySQL/Redis 的环境里验证一遍。

**下一步**：进入 Phase 1 Week 3，按 `06-domain-blog-video-sdk.md` + `07-domain-monitoring-system-misc.md` 继续领域服务改造。

**第二轮审查又发现一处真问题，已修复**：上面第 2 条只改了 Admin 侧 `RateLimitMiddleware`，`SDKRateLimitMiddleware` 的限流响应仍是 `errs.CodeForbidden`（未写 429），审查指出两处限流响应应该保持一致——已同步改成 `errs.CodeTooManyRequests` + `w.WriteHeader(http.StatusTooManyRequests)`，和 Admin 侧写法完全对齐。

---

## 2026-07-11（续）：Phase 1 Week 3（`06-domain-blog-video-sdk.md` + `07-domain-monitoring-system-misc.md`）全部完成，含全仓库 Day 5 扫尾

**1. Blog 域（`06` 第 2 节）**：现状核查与文档假设完全一致（`blog_article_repository.go` 的 `CreateWithTags`/`UpdateWithTags` 确认是唯一需要事务保护的多表写，`blog_article_audit_logic.go:82` 确认直连 `BlogArticleModel.Update`）。
- 新增 `internal/domain/content/blog_service.go`（`package content`）：`BlogArticleService.CreateArticle`/`UpdateArticle`/`AuditArticle` 三个方法均为文档给出的代码原样落地。**新增 `UnpublishArticle`（文档未给出具体代码，但 06 第 2.1 节表格和 07 第 4.2 节都把 `blog_article_audit_unpublish_logic.go` 列为要迁移的文件）**，按 `AuditArticle` 同一模式实现：查文章→改状态→写审核记录三步包进一个 `Transact`，并把原代码里被忽略错误的 `_ = ...Create(auditRecord)`（fire-and-forget）改成传播错误触发回滚——这是从"无事务保护、失败静默"改成"事务保护、失败整体回滚"的真实语义变化，符合本轮改造的目标。
- `BlogArticleRepository` 接口按文档新增纯 `Update(ctx, article)` 方法（不碰标签，薄封装转发 `articleModel.Update`）。
- `CreateWithTags`/`UpdateWithTags` 内部原来的手动补偿删除（`_ = r.articleModel.Delete(...)`）按文档 2.3 节"可以顺手删掉"的建议实际删除了（不是留着），因为留着会让 sqlmock 测试必须额外 mock 一条注定不会被真正执行到的补偿 SQL，删掉后测试和生产代码都更干净。
- `blog_article_create_logic.go`/`blog_article_update_logic.go`/`blog_article_audit_logic.go`/`blog_article_audit_unpublish_logic.go` 四个 logic 文件改为薄委托 `l.svcCtx.Domain.Blog.ArticleService.Xxx(...)`；`blog_article_update_logic.go` 里保留的 `FindByID` 只读调用顺手也改成 `l.svcCtx.Domain.Blog.Article.FindByID`（该字段本来就在 registry 里，不新增依赖）。
- `public_blog_author_info_logic.go` 按文档 2.4 节只加 `TODO(phase2-content-rpc)` 注释标记，不修跨域 import。

**2. SDK 域（`06` 第 4 节）**：现状核查同样与文档假设一致（`SaveBindings` 的"先删后插"确认是唯一需要事务保护的地方，`SdkAdminRepository` 本身持有完整 `*repository.Repository` 字段，构造在事务绑定过的 Repository 上即可让内部 SQL 自动共享事务）。
- 新增 `internal/domain/sdk/sdk_service.go`（`package sdk`）：`SDKService.SaveApiKeyBindings` 照文档代码落地。
- `sdk_api_key_bind_save_logic.go` 改为薄委托 `l.svcCtx.Domain.SDK.Service.SaveApiKeyBindings(...)`。
- `sdk_call_log_export_logic.go` 按文档 4.1/4.3 节只加 `TODO(phase2-task-rpc)` 注释标记（原文档描述成 IAM 跨域，实际读代码后是跨到 Task 域创建导出任务，注释按实际情况标注为 `phase2-task-rpc`，不是文档字面的措辞）。

**3. `registry.Domain` 接线（`06` 第 5 节）**：`BlogDomain` 新增 `ArticleService *contentdomain.BlogArticleService` 字段，`SDKDomain` 新增 `Service *sdkdomain.SDKService` 字段，均挂在各自原有子结构体上，不新建 `ContentDomain` 聚合体（文档明确要求本轮不做，等 Phase 2 拆 `content-rpc` 时再考虑）。

**4. 测试（`06` 第 7 节）**：新增 `internal/domain/content/blog_service_test.go`、`internal/domain/sdk/sdk_service_test.go`，`sqlmock` 覆盖：
- `BlogArticleService.CreateArticle`：happy path + rollback path（第二条标签 INSERT 失败，断言无需任何手动补偿逻辑）。
- `BlogArticleService.AuditArticle`：happy path + "状态不是待审核"分支（无任何写）+ 审核记录写入成功但文章状态更新失败的 rollback path。
- `BlogArticleService.UnpublishArticle`（新增方法，文档未要求但补齐）：happy path + "非上架状态不可下架"分支。
- `SDKService.SaveApiKeyBindings`：happy path + rollback path（第二条绑定 INSERT 失败，断言旧绑定"看起来"没被清空）。
- 全部通过；调试过程中发现的唯一坑是 squirrel 生成的 `INSERT INTO` 是大写，和 goctl 生成 Model 的 `insert into` 小写不一致，两处 SQL 分别要用不同大小写的正则匹配。

**5. Monitoring / System / Misc 域（`07` 第 1-3 节）**：现状核查结果与文档完全吻合（17/30/10 个文件，多仓储 grep 结果与文档给出的表格逐条一致），确认三个域都不需要领域服务，不新建 `internal/domain/monitoring|system|misc/` 目录，无代码改动。`notice_create_logic.go`/`notice_update_logic.go`（跨域读 IAM 用户）和 4 个 `monitoring/*_export_logic.go`（跨域写 Task 域创建导出任务）按文档要求只加 `TODO(phase2-iam-rpc)`/`TODO(phase2-task-rpc)` 注释标记。

**6. Day 5 全仓库扫尾（`07` 第 4 节）**：`l.svcCtx.Repository` → `svcCtx.Domain.X.Y` 机械替换。
- 写了一次性脚本（未入库，仅本次会话用）按 `NewDomain` 的字段映射表批量替换 `<pkg>repo.New<Xxx>Repository(l.svcCtx.Repository)` 调用点，并对"`varName := l.svcCtx.Domain.X.Y` 之后只被当简单别名使用"的情况做了内联清理（删掉局部变量声明，调用点直接写全路径），再用 `goimports -w` 清理因此产生的未使用 import。
- **脚本的内联清理步骤有 bug，过程中发现并修复**：内联清理是基于变量名的全文件正则替换，没有理解 Go 的函数/块作用域——命中了两类真实语料：① 同名变量在同一文件的不同函数里代表不同东西（`menu_create_logic.go`/`menu_update_logic.go` 的 `menuRepo` 既是外层局部变量、又是内部 `validateMenuHierarchy`/`checkCircularReference` 私有方法的形参名；`sdk_api_key_create_logic.go` 的 `repo` 同理是 `generateUniqueKeyPair` 的形参名），脚本把这些形参声明也错误替换成了 `l.svcCtx.Domain.X.Y`，产生语法错误（`missing ',' in parameter list`），`goimports -w`/`go build` 立刻报出，手工核对后把这 3 处形参声明和函数体内对应用法改回原变量名；② 同一文件里两个不同作用域（外层函数体 + 内部 `go func(){}()` 闭包）各自声明了同名局部变量（`dict_item_create_logic.go` 的 `dictTypeRepo`），脚本先内联了外层声明并把内层声明的变量名也一起错误替换，产生 `非法 :=` 左值（`l.svcCtx.Domain.System.DictType := l.svcCtx.Domain.System.DictType`），`go build` 报错定位，修复方式是直接删掉内层这条完全冗余的重复声明（本来就和外层指向同一个 `Domain` 字段，不需要在闭包里再声明一次）。修复后用 `grep -rn "l\.svcCtx\.Domain\.[A-Za-z.]* :="` 和逐文件参数名唯一性检查确认没有同类漏网实例。
- 最终统计：`grep -rl "l.svcCtx.Repository" internal/logic | wc -l` 从起点 155（Week1/2 已完成 6 个文件迁移后的基线）降到 **35**；`grep -rl "svcCtx.Domain" internal/logic | wc -l` 从 11 涨到 **159**。35 个剩余文件里：18 个是文档预期的 `BusinessCache` 例外（不是仓储，无对应 `Domain` 字段）；剩余 17 个是文档写作时未穷举到的、同样合法在扫尾范围之外的用法——10 个文件是 `dict.GetIntValue(l.ctx, l.svcCtx.Repository, ...)`（把整个 `*repository.Repository` 传给字典读取工具函数，不是 `xxxrepo.NewXxxRepository(...)` 构造调用，不满足机械替换规则第 1 条）、4 个文件是直接访问 `l.svcCtx.Repository.DB`/`.Redis`（health check、监控统计原生 SQL、Redis 直连，同样不是仓储构造调用）、3 个文件是文档 4.3 节点名的"直连 Model"技术债的新发现（`blog_article_publish/unpublish/submit_logic.go` 直连 `BlogArticleModel.Update`，按规则记录不改，已写入 `11-descoped.md` 第 10 条）。
- `go build ./...`、`go vet ./...`、`go test ./...` 全仓库通过（含本轮新增的 Blog/SDK 领域服务测试）。

**7. 完成的定义核对（对照 `07` 文档第 6 节）**：
1. `go build ./...` 通过 ✅。
2. 扫尾最终数字（35 = 18 BusinessCache + 17 已归类例外）与例外清单已写入本条目 ✅。
3. `internal/domain/` 下没有出现 `monitoring/`、`system/`、`misc/` 子目录 ✅（只新增了 `content/`、`sdk/` 两个 06 文档要求的目录）。
4. 人工冒烟测试（操作日志导出、字典项编辑缓存刷新、公告全员通知、`/api/v1/misc/ping` 健康检查、Blog 文章创建/审核/下架、SDK Key 绑定替换、视频回归）留给用户在有本地 MySQL/Redis 的环境里验证，本轮未执行。

**遗留/需要用户关注的点**：
1. 人工冒烟测试（见上，Blog/SDK/Monitoring/System/Misc 全部留给用户在本地环境验证一遍）。
2. `04` 文档任务 6（`AuthDomainService.Login` 下沉）仍未做，按文档建议优先级低，可继续推迟。
3. `blog_article_publish/unpublish/submit_logic.go` 三处直连 Model 技术债（见 `11-descoped.md` 第 10 条）留给后续专门会话清理。

**下一步**：Phase 1 Week 3 是 `00-refactor-overview.md` 里 Phase 1（Week 1-5）的最后一个域改造周；Week 4-5 按文档集应为测试覆盖补漏 + Phase 1 整体验收，需要用户确认下一步具体读哪份文档（`00-refactor-overview.md` 目录里 Week 4-5 对应文档尚未在本次会话核对）。

---

## 2026-07-11（续二）：Phase 1 Week 4-5（`08-testing-strategy.md` + `09-ci-cd-and-deployability.md` 后半部分）全部完成

**核对结论**：`00-refactor-overview.md` 第 3 节明确 Week 4-5 = 事务审计 + 测试覆盖补漏 + CI 集成测试跑绿 + `golangci-lint` 扩大范围 + Phase 1 文档收尾，对应 `08`/`09` 两篇文档（`09` 只取 Week1 未覆盖的后半部分：`.golangci.yml` 范围扩大 + CI 集成测试）。用户确认按 1→5 顺序全部做完，不中途拆分会话。

**1. 全仓库事务审计（`08` 完成的定义第 2 条）**：写 Python 脚本扫描 `internal/repository/**`、`internal/logic/**` 里所有函数内的多次写调用（`Insert`/`Update`/`Delete`/`Create` 组合），核对每一处是否已经过 `Repository.Transact`/`registry.Transact`。Repository 层本身全部是单写方法（Week1-3 已把真正的多表写方法收进领域服务），扫描命中的是 **Logic 层两处未被此前三周审计覆盖到的多写缺口**：
- `chat_group_create_logic.go`（`ChatGroupCreate`）：建群组 + 拉创建人入群两条写完全没有事务保护，加入失败会留下没有任何成员的孤儿群组——和 Week2 修复的 `user_create_logic.go` 是同一类"未被发现的旗舰级 bug"。改法：`registry.Transact` 直接包在 logic 文件里（没有专属领域服务，按 `01` 文档"两条路径"里第二条处理）。
- `blog_article_top_logic.go`（`BlogArticleTop`）：置顶数量超限时"取消最早置顶"和"设置新文章置顶"两次 `UpdateTopStatus` 独立执行、无事务。新增 `BlogArticleService.SetArticleTop`（`internal/domain/content/blog_service.go`），把判断+两次写全部收进一个 `Transact`，`blog_article_top_logic.go` 改为薄委托。
- 两处均补了 sqlmock happy-path + rollback-path 测试（`internal/logic/chat/group/chat_group_create_logic_test.go` 新增；`internal/domain/content/blog_service_test.go` 追加 `TestBlogArticleService_SetArticleTop_*`）。

**2. 事务审计过程中衍生发现：全仓库 23 个文件的 `Create` 方法丢弃 `LastInsertId`（同一 bug 类，Week2 `user_repository.go` 已经修过一次，这次是系统性排查剩余同类代码）**：
- **4 处确认是当前活跃故障**（不是理论风险，是真的在错误运行）：
  - `chat_repository.go`（`Chat.Create`）：`chat.Id` 恒为 0，直接导致上面第 1 条修的 `ChatGroupCreate` 即使补了事务，`chat_user.chat_id` 依然写成 0——事务审计和这个 bug 互相牵连，先修事务包裹不修这个 ID 丢失，效果等于没修。
  - `chat_message_repository.go`（`ChatMessage.Create`）：`message.Id` 恒为 0，`chat_message_send_logic.go` 把它用在 WebSocket 广播的 `MessageID` 字段和 HTTP 响应的 `Id` 字段，两处都一直是错的。
  - `file_repository.go`（`File.Create`）：`file.Id` 恒为 0，`file_upload_logic.go` 的上传响应 `Id` 字段一直返回 0。
  - `notice_repository.go`（`Notice.Create`）：`notice.Id` 恒为 0，`notice_create_logic.go` 用它触发 `createNotificationsForAllUsers(notice.Id, ...)` 给全员发公告通知，实际发出去的通知全部关联到 `notice_id=0`。
  - 4 处修法一致：`result, err := r.model.Insert(...)` 之后 `result.LastInsertId()` 回填 `.Id`，和 Week2 `user_repository.go` 的修复完全同构。
- **19 个文件同类但当前未被下游读取、判定为"预防性修复"**（不清零就是留给下一个人意外撞见的地雷）：`demo`、`daily_short_sentence`、`chat_user`、`notification`、`config`、`dict_item`、`department`、`role`、`dict_type`、`permission`、`menu`、`api`、`performance_log`、`operation_log`（含 `Create` 和 `BatchCreate` 两处）、`audit_log`、`blog_friend_link`、`blog_article_audit`、`blog_tag`、`blog_social_info` 各自的 `Create` 方法，同一模式批量修复。association 批量插入循环（`user_role_repository.go`/`permission_api_repository.go`/`permission_menu_repository.go`）**特意不修**——关联表自身的自增 ID 在当前代码里任何地方都不会被用到（关联记录靠外键对查，不靠自身 ID），修了也是无意义的噪音。
- `go build ./...`、`go vet ./...`、`go test ./...` 全绿。

**3. 测试覆盖补漏（`08` 第 4 节）**：
- 登录 4 分支（`internal/logic/iam/auth/login_logic_test.go`，新增）：用户不存在、密码错误、账号禁用、登录成功，均用 sqlmock + miniredis 打桩真实 `svc.ServiceContext`。异步的 `recordLoginLog`/`createUnreadNoticeNotifications`（登录成功后触发）不在断言范围内——sqlmock 对未预期调用只返回 error 被内部 `logx` 记录，不影响测试通过/失败，`-race -count=3` 验证过没有数据竞争。
- Refresh Token 3 分支（`internal/logic/iam/auth/refresh_logic_test.go`，新增），**测试编写过程中发现一个真实安全漏洞**：`refresh_logic.go` 原实现只校验 JWT 签名和过期时间，从未检查黑名单——`authmiddleware.go` 对 access token 有黑名单校验，`logout_logic.go` 退出登录时会把 access/refresh token 都拉黑，但 `Refresh` 端点完全不查，意味着退出登录后原 refresh token 在自然过期前仍能不断换发新的 access token，退出登录形同虚设。修复：`Refresh` 里补一次 `TokenBlacklist.IsBlacklisted` 校验，和 `authmiddleware.go` 对齐。3 个分支：token 过期、token 在黑名单（种子数据直接在 miniredis 里种黑名单 key）、正常刷新。
- 权限校验（`08` 第 4 节"权限不足时返回正确错误码"的烟测 + 第 2 节"值得测"的 `PermissionResolver.CanAccess`）：`internal/domain/iam/permission_resolver_test.go`（新增，5 个分支：超级管理员直通、接口不存在、无角色拒绝、命中权限放行、权限不匹配拒绝）+ `internal/middleware/permissionmiddleware_test.go`（新增，中间件层烟测：未登录 401、无权限 403，断言 `pkg/response.ErrorCtx` 统一走 HTTP 400 + 业务码的既有约定）。
- 全部新增测试 `go test ./... -race` 通过。

**4. 集成测试套件（`08` 第 5 节，`//go:build integration`）**：新增 `internal/integration/` 包（`setup_test.go` 共享 `TEST_MYSQL_DSN`/`TEST_REDIS_ADDR` 环境变量打底，缺一个就 `t.Skip`）+ `internal/domain/task/scheduler_integration_test.go`（同包直调未导出的 `scanAndExecute`，不用等真实 ticker）。6 个场景全部按文档清单写完：
1. 登录 e2e（`login_test.go`）。
2. RBAC 允许的请求（`rbac_test.go`，真实拼 `AuthMiddleware.Handle(PermissionMiddleware.Handle(...))` 请求链，断言 200）。
3. RBAC 拒绝的请求（同文件，断言业务码 403）。
4. IAM 用户创建 → chat onboarding 全链路（`chat_onboarding_test.go`，轮询代替固定 sleep，等新用户出现在默认企业群组成员列表里）。
5. Task 调度器完整周期（`scheduler_integration_test.go`，用一个只在本文件里注册的测试专用 `task_type=999999` 执行器，不污染真实字典）。
6.（可选，文档标记可选）`Repository.Transact` 在真实 MySQL 上的回滚验证（`transaction_rollback_test.go`）。
- **本机没有 Docker/MySQL/Redis，这 6 个测试只验证了 `go build -tags=integration ./...`、`go vet -tags=integration ./...` 编译通过，从未真正连接过真实数据库执行——运行时正确性完全没有被验证过**，需要用户在有 Docker 的环境跑 `docker compose up` 起 MySQL/Redis 后手动执行一遍 `go test -tags=integration ./... -count=1`（或者直接等下次 CI 触发），如果跑出编译时看不出来的逻辑错误（字段名、SQL 匹配、时序假设）需要回来修。

**5. `golangci-lint` 扩大范围（`09` 第 3 节 + `00` 第 3 节 Week4-5 任务项）**：`golangci-lint run ./...` 扫全仓库，只有 13 个既有问题（比预想的"512 个文件太吵"轻得多），逐个核实后全部修掉，不是简单 `_ =` 消音了事：
- `internal/logic/system/file/file_upload_logic.go`：`file.Seek(0, 0)` 未检查错误——这不是无害噪音，Seek 失败会导致后面 `io.Copy` 从错误偏移量开始拷贝，落盘文件内容和用于命名的 MD5 对不上，改成真正返回错误。
- `pkg/monitor/performance.go`：`SA9003` 空分支——是一段被注释掉的死代码（`else {}` 里只剩注释），直接删除整个空分支。
- `internal/repository/chat/chat_repository.go`：`ineffassign`——是 Week1 CI 审查时就记录过的已知问题（一段被废弃的 JOIN+GROUP BY 查询字符串赋值后立刻被子查询版本覆盖），这次顺手删掉死代码。
- `internal/hub/chathub.go`（8 处）+ `internal/config/loader.go`（1 处）：WebSocket 读写超时设置、Redis DB 编号解析失败均为已知的"失败即用默认值/静默降级"既有语义，显式 `_ = ` 丢弃，不改变行为。
- `internal/handler/chat/chatwshandler.go`：`BroadcastChatMessage` 错误未检查——改成和 `chat_message_send_logic.go` 一致的 `logx.Errorf` 记日志，不是静默丢弃。
- 修复后 `golangci-lint run ./...` **全仓库 0 issue**。`.github/workflows/ci.yml` 的 `lint-build` job 从 Week1 的 `./internal/repository/... ./internal/domain/...` 改成 `./...`，不再限制路径。

**6. CI 集成测试配置核对**：`ci.yml` 的 `integration-test` job（DB 初始化顺序、`TEST_MYSQL_DSN`/`TEST_REDIS_ADDR` 环境变量名）核对下来和本次新增的 `internal/integration/` 包完全对得上，不需要改动——这部分本来就是 Week1 审查回合修复过的。`go test -tags=integration ./... -count=1` 会自动带上本轮新增的 6 个集成测试，不需要额外配置。

**完成的定义核对（对照 `08`/`09` 两篇的验收标准）**：
1. `go build ./...`、`go vet ./...`、`go test ./...`（含 `-race`）全绿 ✅。
2. `go build -tags=integration ./...`、`go vet -tags=integration ./...` 编译通过 ✅；`go test -tags=integration ./...` 实际执行 ❌（本机无 Docker，需要用户验证，见上）。
3. 每个新增的 `Transact`/`registry.Transact` 包裹方法都有 happy-path + rollback-path 测试 ✅（`SetArticleTop`、`ChatGroupCreate`）。
4. `golangci-lint run ./...` 全仓库 0 issue ✅。
5. CI 三个 job 的配置本身自洽（`lint-build`/`unit-test`/`integration-test`）✅；实际跑绿需要真实触发一次 GitHub Actions，本次会话没有推送验证。

**7. 提交前 Cursor 自动代码审查（Gentleman Guardian Angel）发现问题修复**：`git commit` 触发的 `gga` pre-commit hook 第一次审查 FAILED，逐项核实后处理如下。

**已修复（判断为真问题）**：
1. `sdk_api_key_update_logic.go`：`if req.IpWhitelist != "" || req.IpWhitelist == ""` 是永真式（不管取值都会执行），导致每次更新 API Key 都无条件覆盖 `IpWhitelist`，哪怕调用方没传这个字段——和同函数里 `Name`/`Status`/`ExpireAt`"未提供则不改"的既有模式不一致。改成 `if req.IpWhitelist != ""`，和其余字段的模式对齐。
2. **`admin_notification.read_status` 与 `admin_login_log.status` 存在两套冲突的取值约定，是本次审查发现的最严重问题、预存在于本轮重构之前**：建表 DDL 注释和三处通知创建代码（`login_logic.go`/`notice_create_logic.go`/`notice_update_logic.go`）一直用的是"布尔式"取值（`read_status`: 0=未读,1=已读；`status`: 0=失败,1=成功），但字典种子数据（`db/data.sql` 的 `read_status`/`login_status` 两个字典类型）按项目"字典枚举 value 从 1 开始"的约定写成了另一套（1=未读,2=已读；1=成功,2=失败），列表筛选（`notification_list_logic.go`/`login_log_list_logic.go`/`login_log_repository.go`）和 Task 域通知创建（`internal/domain/task/notifier.go`，Week1-2 期间新写的代码，已经是对的）走的都是字典这套。两套取值同时存在导致：标记单条已读实际写成"未读"的值、`MarkAllAsRead` 的 `WHERE read_status=0` 从来匹配不到任何行（真实存量数据不会是 0）、`ClearRead` 按"已读"的旗号实际删除的是未读消息、登录日志按成功/失败筛选和统计（`login_log_stats_logic.go`）失败次数永远查不到。用户确认现在就修，统一到字典这套取值：
   - 代码：`login_logic.go`（通知创建 `ReadStatus`、登录日志 `status`）、`notice_create_logic.go`/`notice_update_logic.go`（`ReadStatus`）、`notification_read_logic.go`（标记已读写 2 不是 1）、`notification_repository.go`（`MarkAllAsRead`/`ClearRead` 目标值和 WHERE 条件）、`login_log_stats_logic.go`（失败次数统计传 2 不是 0）。
   - `db/tables.sql`：两个字段的 DDL 注释和 `DEFAULT` 值同步改成新取值（`DEFAULT 1`，不再用 0 表示任何具体业务状态，和项目"0 只作不筛选占位"的字典规则对齐）。
   - `db/data.sql`：两条通知种子数据的 `read_status` 从 0 改成 1，全新部署从一开始就是对的取值，不需要额外迁移。
   - 新增 `db/migrations/fix_status_semantics_20260711.sql`：给已经用旧版 `data.sql` 跑起来、积累了真实业务数据的库做一次性存量数据迁移（不是幂等脚本，只应执行一次；全新部署不需要跑，也刻意不匹配 `docker-init.sh`/CI 的 `create_table_*`/`dict_*`/`init_*` 通配符，不会被自动流程误跑）。

**核实后判断为超出本轮范围、不修（保留分歧记录，和 2026-07-10 续二 条目同样的处理方式）**：
1. 3 处 `blog_article_publish/submit/unpublish_logic.go` 直连 `BlogArticleModel.Update`——已经是 `11-descoped.md` 第 10 条记录的已知技术债，本轮不顺带清理。
2. `chatwshandler.go` 内联构造 `iamrepo.NewTokenBlacklistRepository`——和 `authmiddleware.go` 已有的同款写法一致（WebSocket 升级握手没有走标准 REST 中间件链，两处都是不得不在非 Logic 层做鉴权的特例），不是新引入的越界。
3. 审查提到多个仓储文件"裸 SQL 违反 squirrel 规范"——核实后大多是全参数化的静态查询（`chat_user_repository.go`/`menu_repository.go`/`api_repository.go`），没有拼接用户输入，不属于规则真正防的动态 SQL 拼接风险；`chat_repository.go` 本身就是 `10-go-code-style.mdc` 写明的已知例外。`performance_log_repository.go` 审查说"仍未迁移"经核实是错的，已经在用 squirrel（规则文档里的例外清单本身没更新，不是代码问题）。`permission_repository.go` 的 IN 子句用 `strings.Repeat` 拼 `?` 占位符是风格不一致但无注入风险，本轮不顺带改。
4. 硬编码业务常量（`Status: 1`、`Type: 2`、`"notice"` 字符串等）、`Code: 0` 字面量 vs `errs.CodeOK`——全仓库大量既有代码都是这个风格，本轮没有引入新的此类代码，清理是独立的、范围大得多的任务，不顺带做。
5. `chat_user_repository.go` 的 `DeleteByChatIDAndUserID` 物理删除——是否等同 RBAC 关联表那样允许物理删除需要产品确认语义，本轮不擅自改。
6. `internal/consts/consts.go` 的 `PathTaskCancel = "/api/v1/tasks/:id/cancel"`——核实是从未被引用的死常量（`grep` 全仓库零命中），不影响任何真实路由，本轮不清理。
7. `metric_report_options_logic.go`/`video_collect_options_logic.go` 是 goctl 脚手架留下的 `todo` 占位、从未实现——核实是本轮之前就存在的未完成功能，不是 Week1-5 任务书范围内的模块，不顺带补。

`go build`/`go vet`/`go test -race`/`golangci-lint run ./...` 修复后重新全部跑绿。

**第二轮审查（修复后重跑）又发现一处漏网的同类真问题，已修复**：`internal/repository/video/video_repository.go` 的 `Create` 用 squirrel 手写 `INSERT`（不走 goctl 生成的 `model.Insert`，第一轮排查按 `r.\w*model.Insert(ctx` 的模式扫描没有覆盖到这种写法），同样丢弃了 `LastInsertId`——`video_collect_logic.go`（SDK 视频采集接口）建完视频后立刻用 `v.Id` 构造响应，实际一直返回 `Id: 0`。按同一模式修复（`ExecCtx` 拿到 `result` 后回填 `video.Id`）。全仓库重新扫了一遍 `sq.Insert(...)+ExecCtx` 组合，确认没有第三个漏网的。

**第三轮审查又指出一处同类问题（风险标注为低，仍顺手修复）**：`internal/repository/monitoring/login_log_repository.go` 的 `Create` 其实已经调用了 `result.LastInsertId()`，但只用于打一行调试日志，从未把结果赋回 `log.Id`——当前 `recordLoginLog` 是异步 fire-and-forget、调用方不读返回值，所以是预防性修复而不是活跃故障，但既然同一个 bug 类已经系统排查过，这处也顺手补上，不留漏网之鱼。

**同一轮审查另补了根目录 `docs/后端开发进度.md` 第 16 节**，记录本轮"对外行为发生实质变化"的部分（群聊创建此前不可用、多处响应 ID 恒为 0、Refresh Token 黑名单漏洞、`read_status`/`login_status` 筛选语义修复），不重复 `progress.md` 的完整过程记录。

`go build`/`go vet`/`go test -race`/`golangci-lint run ./...` 再次全绿，`git commit` 通过。

**遗留/需要用户关注的点**：
1. `internal/integration/` 全部 6 个 + `scheduler_integration_test.go` 集成测试从未连接真实 MySQL/Redis 跑过，只验证了编译通过，需要用户在有 Docker 的环境实测一遍（`docker compose up` 或连接开发库，`export TEST_MYSQL_DSN=... TEST_REDIS_ADDR=... && go test -tags=integration ./... -count=1`）。
2. `04` 文档任务 6（`AuthDomainService.Login` 下沉）仍未做，本轮继续判断优先级低，推迟到 Phase 1 收尾之后再评估要不要做。
3. ~~`docs/后端开发进度.md`（仓库根目录）本轮未更新~~——已在第三轮审查回合补上第 16 节，记录本轮对外行为有实质变化的部分。
4. `db/migrations/fix_status_semantics_20260711.sql` 是非幂等的一次性存量数据迁移，只有已经用旧版 `data.sql` 跑起来、积累了真实业务数据的库需要手动执行一次；全新部署不需要。
5. 本次会话已执行 `git commit`（单个提交，涵盖 Week2-5 全部改动），未 `push`。

**下一步**：Phase 1 Week 4-5 已完成，Phase 1（Week 1-5）整体验收条件基本具备（`01`~`09` 全部落地），下一步是 Phase 1 收尾确认（人工冒烟 + Week4-5 遗留的集成测试真实验证）后，决定是否进入 Phase 2（`15-service-boundaries.md` 起）。
