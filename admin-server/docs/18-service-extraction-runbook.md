# 18 — 服务拆分执行手册（重复执行 5 次的 Runbook）

> 前置依赖：已读 `15-service-boundaries.md`（服务边界、schema 归属、db/services 目录结构）、`16-rpc-conventions.md`（proto 约定、`generate-rpc.sh`、gateway 变薄机制、Wire 装配）、`17-async-eventing.md`（TaskCallback、Redis Streams）。本文档不重复解释"为什么这样拆"，只给"具体怎么拆"的可执行步骤。

## 0. 使用方式

第 1 节是**通用 checklist，写一次，按 `15-service-boundaries.md` B.6 顺序（task-rpc → sdk-rpc → chat-rpc → content-rpc → iam-rpc）执行 5 遍**。每次执行严格走完全部步骤再进入下一个服务，不要并行拆两个服务——单体到微服务的迁移过程中,前一个服务的拆分方式和踩过的坑会直接影响下一个服务怎么拆得更快,顺序执行本身是这套 Runbook 设计的一部分（B.6 的顺序本身就是"先拆风险最小、最能验证机制的",不是任意顺序）。

第 2 节是 5 个服务各自的**差异附录**——每次执行到 checklist 某一步如果这个服务有特殊情况，回来查对应附录。

## 1. 通用 Checklist（每个服务执行一遍）

### 1.1 服务骨架

- [ ] `services/<name>/` 目录建好，子目录 `rpc/`、`internal/{config,svc,logic,server}`、`etc/`。
- [ ] 编写 `services/<name>/rpc/<name>.proto`（方法集合来自该域现有 `.api` 定义的接口 1:1 对应，参考 `16-rpc-conventions.md` 第 4 节的 `iam.proto` 写法）。
- [ ] 跑 `scripts/generate-rpc.sh <name>`，确认生成了 `pb/`、`<name>client/`、`internal/logic/`、`internal/server/`、`internal/config/config.go`、`internal/svc/servicecontext.go`、`<name>.go`（main 入口）、`etc/<name>.yaml`。

### 1.2 领域代码搬迁

- [ ] 把 `internal/domain/<domain>`（Phase 1 已经按 5 服务分组组织好的领域服务包）整体移动到 `services/<name>/internal/domain/`，import 路径更新。
- [ ] 把 `internal/logic/<domain>/**/*.go` 的手写方法体，按 goctl 生成的新骨架（`services/<name>/internal/logic/`）逐个搬迁——这一步是"复制方法体 + 删旧目录"，不是原地改名，要文件级 diff 核对不要漏搬（同上一轮 DDD-lite 任务书 Phase 3 的搬迁方式，原文档已删除，经验证过一次）。
- [ ] 把 `internal/repository/<domain>/*.go`、`internal/model/<domain>/*.go` 整体移动到 `services/<name>/internal/repository/`、`services/<name>/internal/model/`，包名不变（域名本身已经是包名，不需要再改）。
- [ ] 编译该服务自己的 `go build ./services/<name>/...` 通过。

### 1.3 数据库拆分

- [ ] 按 `15-service-boundaries.md` 第 4 节的 `db/services/<service>/` 目录，把该服务名下所有模块的 `create_table_*.sql`/`init_*.sql`/`migrations/` 确认齐全（如果 Phase 1 已经提前做完这一层目录重组，这一步只是核对，不需要重新搬文件）。
- [ ] 建对应的 MySQL schema（`CREATE DATABASE admin_<service> ...`），跑该服务名下全部建表 SQL。
- [ ] 该服务的 `etc/<name>.yaml` 配置连接到新 schema（本地/开发环境执行，按 `10-dev-execution-and-review-points.md` 的政策，AI 可以直接执行本地/开发库 SQL，不用等用户）。
- [ ] 从旧的单体库里把这部分表标记为待清理（不立即删除，等这个服务在生产稳定运行一段时间后再清理，具体时机记入 `14-production-deployment-checklist.md`，不在开发期执行）。

### 1.4 跨服务调用接入

- [ ] 该服务如果是 RPC 调用方（如 task-rpc 需要回调 iam-rpc 取导出数据），按 `17-async-eventing.md` 的契约接入 `taskcallback` client。
- [ ] 该服务如果是 Streams 生产者/消费者，按 `17-async-eventing.md` 第 2 节接入对应 Stream 的生产/消费逻辑，包含幂等和死信处理。
- [ ] 该服务如果暴露给 gateway 调用（几乎所有服务都是），在 `internal/wire/providers.go`（gateway 侧）新增对应的 `provide<Name>RPC` provider，`config.Config` 新增对应的 `zrpc.RpcClientConf` 字段，`etc/gateway.yaml` 新增静态 `Endpoints` 配置（参考 `16-rpc-conventions.md` 第 5 节）。
- [ ] gateway 侧原来这个域的 handler/logic 改成"解析请求 → 调 zrpc client → 映射响应"的薄胶水，原来的业务方法体应该已经在 1.2 步骤搬空了,这一步只是清理调用方式。

### 1.5 部署配置

- [ ] 写 `services/<name>/Dockerfile`（复用 `.template/docker/docker.tpl` 的模式，Phase 1 已经给单体写过一版，改成多阶段构建 + 只打包这一个服务的二进制）。
- [ ] `docker-compose.yml` 新增这个服务的条目（镜像、`etc/<name>.yaml` 挂载、依赖的 MySQL schema、依赖的共享 Redis、`depends_on` 声明）。
- [ ] 更新 Supervisor 配置 → docker-compose 服务条目的对应关系（如果还在用 Supervisor 过渡，两者并存一段时间；最终按 Part C.3 的节奏切到纯 compose）。
- [ ] 更新 `14-production-deployment-checklist.md`，追加本次拆分的具体条目：新增的 schema 建库/建表 SQL、新的服务间凭证/配置项、镜像 tag、compose/Supervisor 切换步骤，格式统一为"触发条目的改动是什么 → 部署时具体要执行的命令/SQL/配置变更 → 如何验证生效"。
- [ ] 更新 `progress.md`：记录这个服务拆分做了什么、遇到的问题、和 Runbook 描述有出入的地方（下一个服务拆分前回来看一眼，看有没有需要反馈回本文档的经验）。

### 1.6 验证

- [ ] `go build ./...`（整个 repo，含 gateway + 已拆分的服务 + 还未拆分、仍在单体里的部分）通过。
- [ ] 该服务的单元测试（Phase 1 已经写好的领域服务测试，搬迁后应该原样能跑）通过。
- [ ] 本地 `docker compose up` 起 gateway + 该服务 + 依赖的 MySQL/Redis，人工冒烟测试这个服务对外暴露的每一类接口（至少各一个：一个读、一个写、如果有跨服务调用则至少验证一次）。

## 2. 各服务差异附录

### 2.1 task-rpc（第一个拆，验证机制）

只有 4 个 logic 文件（`task_list_logic.go`、`task_detail_logic.go`、`task_cancel_logic.go`、`task_recent_logic.go`），已经有独立的 `internal/domain/task` 领域包（`scheduler.go`/`executors/`/`notifier.go`），调度器的 Redis 锁（`acquireLock`/`releaseLock`）已经是多副本安全的，不需要为拆分重新设计并发模型——这是选它第一个拆的直接原因。

这一轮的目的不是"快速拆完"，而是**把整套机制完整跑通一遍**：`generate-rpc.sh` 是否好用、`db/services/task/` 目录是否真的能对应搬迁、Dockerfile/compose/Supervisor→compose 切换的每一步是否顺畅。预留比后续服务更宽裕的时间（`15-service-boundaries.md` 引用的时间线是 Week 6-7，两周,而 sdk-rpc 只有 1 周），机制本身跑不通就先把 `16-rpc-conventions.md`/本文档改对，不要在有疑问的情况下硬着头皮往下拆第二个服务。

特殊点：
- `task_recent_logic.go` 对 `system` 域字典表的读取（见 `15-service-boundaries.md` 第 5 节末尾、`17-async-eventing.md` disposition 表最后一行）——执行到 1.4 步骤时现场判断：改成 task-rpc 自己的静态配置（推荐,省一次 RPC 依赖）还是保留同步 RPC + 失败降级。两种做法都要在 `progress.md` 里记清楚选了哪个、为什么。
- `ExcelExportExecutor` 的改造（`executors map[int]interfaces.TaskExecutor` → `executors map[int]taskcallback.Client`）是这一轮里唯一需要动"跨服务契约设计"的部分，`TaskCallback` proto（`17-async-eventing.md` 第 1.3 节）在这一步落地，但此时 `iam-rpc`/`sdk-rpc` 还没拆分——`TaskCallback` server 端实现要**在单体里先跑通**（单体里的 `iam` 域代码临时实现一个 `TaskCallback` server,监听在一个本地端口），等 iam-rpc/sdk-rpc 真正拆分时把这个 server 实现原样搬过去，不需要重新设计契约。

### 2.2 sdk-rpc

13 个文件，边界干净（`internal/repository/sdk/` 只有 2 个手写 repository 文件，`sdk_admin_repository.go` + `sdk_repository.go`，对应"admin 管理面 + public 调用面 + 调用日志"三块职责已经在现有代码里靠这两个文件分开），唯一跨域依赖是对 task-rpc 的调用（`sdk_call_log_export_logic.go`），第 1 步 task-rpc 拆分时已经验证过这条路径的机制（`TaskCallback`/`SubmitTask`）。

执行到 1.2 步骤时验证一件事：`sdk_admin_repository.go`（对应后台管理这几张表的 CRUD：`sdk_key`/`sdk_interface`/`sdk_key_api`）和 `sdk_repository.go`（对应外部调用方视角：`SDKAuthMiddleware`/`SDKRateLimitMiddleware` 用到的鉴权查询、`sdk_call_log` 写入）这个二分是否干净地对应到 `db/services/sdk/sdk/` 单一模块目录下——`15-service-boundaries.md` 里 sdk 只有一个模块目录（不像 iam 拆成 user/role/permission 等多个），执行时确认这个"粗粒度模块"假设站得住,如果发现 admin/public 两块实际上有交叉引用需要拆得更细,回来更新 15 文档的目录树。

`SDKAuthMiddleware`/`SDKRateLimitMiddleware`/`SDKCallLogMiddleware` 这三个中间件本身要不要搬到 sdk-rpc 侧，还是留在 gateway 侧只调 sdk-rpc 做鉴权判断——推荐做法是中间件本身留在 gateway（因为 gateway 才是收 HTTP 请求的地方,限流/鉴权判断需要在请求最早期发生），内部实现改成调 `sdk-rpc.VerifyAPIKey`/`CheckRateLimit` 之类的 RPC 方法，`sdk_call_log` 的写入按 `17-async-eventing.md` 的思路评估是否也要走 Streams（如果判断标准是"失败不影响本次调用" → 异步；SDK 调用日志目前的写入方式是同步 `Create`，需要在这一步现场核查实际代码判断，不在本文档预先下结论）。

### 2.3 chat-rpc

多留一周（`15-service-boundaries.md` 引用的时间线 Week 9-10）做 WS↔gRPC 流桥接——`16-rpc-conventions.md` 第 7 节已经给出了完整的 `Stream` RPC 方法签名和网关侧桥接 goroutine 的形状，执行时对照那一节的代码骨架实现，`ChatHub`（`internal/hub/chathub.go`）的连接表数据结构原样搬到 chat-rpc,只是"连接"的类型从 `*websocket.Conn` 换成 `grpc.ServerStream`。

这一步**依赖 `17-async-eventing.md` 的 Streams 机制已经被第 1 个服务（task-rpc 拆分时其实还没有真正用到 Streams,Streams 第一次真正投入使用是这里）验证过**——严格说 Streams 机制在 task-rpc 拆分阶段并没有被使用（task-rpc 用的是 `TaskCallback` 同步 RPC，不是 Streams），所以 chat-rpc 拆分实际上是 `stream:chat.user.created` 第一次真正上生产路径。执行到 1.4 步骤时，`iam-rpc` 那一侧（用户创建发布事件的生产者）此时也还没拆分——和 task-rpc 附录里 `TaskCallback` 的处理方式一样：**先在单体里的 `iam` 域代码里加上"发布 `stream:chat.user.created` 事件"这一行**（这是一个很小的改动，不需要等 iam-rpc 真正拆分），chat-rpc 侧的消费者逻辑按 `17-async-eventing.md` 第 2.1 节实现并接入,完整验证一次用户创建 → 群加入/私聊建立的端到端事件流,再进入 content-rpc/iam-rpc 的拆分。

### 2.4 content-rpc

文件数最多（blog 域约 34 个 logic 文件 + video 域约 7 个），但架构最简单——**明确定性为"机械，不是有风险"**：blog/video 内部都是单表/几张关联表的标准 CRUD，没有 chat 那种状态、没有 task 那种跨服务回调、没有 sdk 那种独立信任边界，1.2 步骤的"复制方法体到新骨架"会花比其他服务更多的时间，但每一步都是重复劳动，不需要现场设计任何新机制。

唯一需要现场判断的跨服务点是 `blog → iam`（`public_blog_author_info_logic.go` 目前写死查 `userID=1`）——按 `16-rpc-conventions.md` 第 4 节，这个位置改成调 `iam-rpc.GetUserProfile`；如果此时 iam-rpc 还没拆分（content-rpc 排在 iam-rpc 前面），处理方式与前两个附录一致：先在单体里补一个可以被临时调用的最小实现或者直接保留内联查询到 iam-rpc 真正拆分时再切换，两种都可以，选择依据是"content-rpc 拆分时切换成本高不高"，现场判断,记入 `progress.md`。

### 2.5 iam-rpc（最后拆，故意的）

放最后**是故意的**：所有其他 4 个服务都依赖 iam-rpc 做权限校验/用户资料查询，如果第一个就拆 iam-rpc，`CheckPermission`/`BatchGetUserProfiles` 这些方法的参数、返回结构、缓存策略全部要凭空设计，容易设计出"当时觉得够用、实际接入才发现不够用"的契约,需要返工。放到最后，前 4 次拆分（task/sdk/chat/content）各自已经真实调用过 iam-rpc 的某个子集能力（chat/content 调 `BatchGetUserProfiles`/`GetUserProfile`，task/sdk 调 `TaskCallback` 侧的 iam-rpc 实现），到这一步时 `iam.proto` 的方法集合和字段设计已经被 4 轮真实用例验证过，不是拍脑袋定的。

拆到 iam-rpc 之前，"尚未拆分的部分"继续留在收缩中的单体里充当 gateway 的角色——不需要第一天就有一个成型的瘦 gateway 二进制，`cmd/gateway/`（或者过渡阶段仍然是 `admin.go`）在前 4 次拆分完成后,自身持有的数据库连接已经只剩 `admin_platform` schema（iam+system+monitoring+misc），本质上就是"iam-rpc 的代码还没被物理搬出去，但已经是唯一剩下的有状态部分"。

这一步的 1.2 步骤（领域代码搬迁）体量最大（iam 域本身约 36 个 logic 文件 + system 约 30 个 + monitoring 约 20 个 + misc 约 8 个），但因为 Phase 1 已经把这些域的领域服务包统一组织在 `internal/domain/iam/`（B.1 节要求的分组），实际搬迁工作量比文件数字面看起来小——真正花时间的是把前 4 轮暴露出来的 `iam.proto` 补充需求（如果有）在这一步一次性落地到位，而不是分批修契约。

`CheckPermission`/`BatchGetUserProfiles` 内部接的权限缓存（`CacheKeyUserPermissions`）、Token 黑名单直连 Redis 的设计（`16-rpc-conventions.md` 第 6 节）在这一步落地为真实代码，之前 4 次拆分里 gateway 侧的 `AuthMiddleware`/`PermissionMiddleware` 一直是直接读单体里的 `Domain`/`Repository`（因为 iam 还没搬走），这一步是这两个中间件第一次真正切换成调 `iam-rpc` 的 zrpc client。

## 3. 非目标

- 不要求 5 个服务拆分过程中始终保持"单体和微服务同时能跑"的双轨制——每个服务拆完就切换调用方式，不维护旧代码路径（`AGENTS.md` 第 5 节"保留旧代码路径/兼容层"本来就是明确的反例）。
- 不做拆分过程中的"灰度流量切换"机制（如按用户 ID 灰度切到新服务）——项目未上线,没有真实用户,不需要灰度,直接切换。
- 不为拆分过程本身写自动化回归测试套件——复用 Phase 1 已经写好的领域服务单元测试（搬迁后应该原样能跑），加一次跑通即可，不新增专门为"验证拆分没搬错"而写的测试。

## 4. 完成的定义（对每一次服务拆分）

- 第 1 节 checklist 的全部条目勾选完成。
- `go build ./...`（整个仓库）通过，该服务的单元测试通过。
- 本地 `docker compose up` 冒烟测试通过（gateway 能正确路由到新拆出的服务）。
- `14-production-deployment-checklist.md`、`progress.md` 都已经更新（含这次拆分特有的经验，供下一次拆分参考）。
- 5 次全部完成后：整个仓库不再有指向已删除的 `internal/{handler,logic,repository,model,domain}/<已拆分域>/` 的引用,`internal/` 顶层只剩 gateway 自己的关切（handler/types/middleware/wire），`cmd/gateway/` 是唯一保留的单体式入口,进入 Part C（可观测性/API 文档/CD）。
