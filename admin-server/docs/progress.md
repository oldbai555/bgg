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

**第二轮审查又发现一处真问题，已修复**：上面第 2 条只改了 Admin 侧 `RateLimitMiddleware`，`SDKRateLimitMiddleware` 的限流响应仍是 `errs.CodeForbidden`（未写 429），审查指出两处限流响应应该保持一致——已同步改成 `errs.CodeTooManyRequests` + `w.WriteHeader(http.StatusTooManyRequests)`，和 Admin 侧写法完全对齐。
