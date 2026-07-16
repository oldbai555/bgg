# admin-server 重构任务书总纲（Phase 1-3）

> 本文档是 `admin-server/docs/` 整个文档集（23 篇 + `progress.md`）的入口和索引。执行者（Cursor / Claude Code）在开始任何一个 Phase 之前，应先完整读一遍本文档，再按需跳转到具体的分篇文档。文档风格对齐仓库已有的 `docs/admin-server-ddd-refactor-prompt.md`：可直接执行的任务说明，不是设计讨论稿。

## 0. 这一批文档要解决什么

`admin-server`（go-zero 单体，Go 1.24，模块 `postapocgame/admin-server`）在上一轮 DDD-lite 重构（见 `docs/admin-server-ddd-refactor-prompt.md`，已完成：目录按 9 域拆分、`internal/domain/{iam,task}` 建立、`ServiceContext` 精简）和 Wire 组合根引入（见最近提交 `36fbda9`）之后，架构改造实际只完成了皮层。本轮的目标不是继续修修补补，而是产出一套**详细到可以让 AI 编码工具独立按文档执行**的任务书文档集，覆盖：单体内部加固 → 微服务拆分 → 可观测性/API 文档/CD。**本轮产出的是文档，不是代码**——`00`~`22` 全部 23 篇 + `progress.md` 已随本轮交付完成，`internal/` 下没有任何代码改动；从下一轮开始按 Phase 1 Week 1 顺序实际动代码。

## 1. 现状核查（为什么要做这一轮）

以下问题是本轮规划前对代码库的实地核查结果，不是推测：

| 问题 | 证据 |
|---|---|
| 规模 | 512 个 `.go` 文件，约 4.26 万行代码 |
| 零测试/CI/Dockerfile | 仓库内无 `*_test.go`（生成代码除外的手写测试）、无 `.github/workflows`、无 `Dockerfile`；`.template/docker/docker.tpl` 存在但从未被使用 |
| DDD-lite 是皮层 | `internal/domain/` 下只有 `iam/permission_resolver.go` 和 `task/`（`scheduler.go`/`executors/`/`notifier.go`）两个域有真正的领域服务；其余 7 个域（blog/video/chat/sdk/monitoring/system/misc）的 logic 文件直接调用 repository，没有领域层 |
| DI 迁移几乎没推进 | `internal/repository/registry/domain.go` 的 `registry.Domain` 聚合已经搭好（9 个子域、约 40 个 repository 字段），但全仓库 161 个 logic 文件里只有 1 个（`permissionmiddleware.go` 间接经由 `iamdomain.NewPermissionResolver`）真正用上，其余仍是方法内内联 `xxxrepo.NewXxxRepository(l.svcCtx.Repository)` |
| 零事务 | `internal/repository/repository.go`、`internal/repository/registry/domain.go` 都没有任何 `Transact`/事务相关方法；跨表写入靠"手写 squirrel、显式不进事务"的方式绕过，例如 `internal/repository/blog/blog_article_repository.go` 的 `CreateWithTags`/`UpdateWithTags` 源码注释直接写着"使用 squirrel 手动插入，避免依赖事务 session API" |
| 旗舰级 bug：`user_create_logic.go` | `internal/logic/iam/user/user_create_logic.go` 的 `initChatForNewUser`：① 用 `userRepo.FindPage(ctx, 1, 10000, "")` 一次性拉全量用户做 N+1 私聊初始化；② IAM 域直接 import `chatrepo "postapocgame/admin-server/internal/repository/chat"`，跨域越界；③ 建用户→加入默认群→逐个建私聊全程没有事务，任何一步失败都不回滚，当前处理方式是 `logx.Errorf` 记日志、不中断主流程 |
| 密钥硬编码 | `etc/admin-api.yaml` 的 `JWT.AccessSecret`/`JWT.RefreshSecret` 是明文占位符字符串，随源码提交仓库 |
| 中间件/Wire 是半成品 | `internal/wire/providers.go` 的 `buildMiddlewareBundle(svcCtx)` 手工现场 `New` 全部 11 个中间件，因为每个中间件构造函数都吃整个 `*svc.ServiceContext`，不是真正的 Wire 装配 |

这是用户自己的项目，**尚未上线，没有外部用户，不需要考虑兼容性**。目标是把系统做成"对标开源项目"的工程质量，之后主要靠 Cursor / Claude Code 按文档独立推进。

## 2. 已确认的范围与决策（执行时不用再问）

1. **不对外发布**——"开源级别"是内部质量对标，不做 LICENSE/品牌脱敏/README 门面工作。
2. **项目未上线，无兼容性负担**——可以自由重新设计 schema、拆分数据库、重组代码目录，不需要保留任何历史兼容层。
3. **本轮全部纳入，无 descope 保留项**：DI/Domain 迁移 + 跨域越界修复、事务安全、测试覆盖、CI/CD、全部 9 个域补齐领域层、真正的微服务拆分、结构化日志 + trace-id、Swagger/OpenAPI 文档生成、CD 自动部署。
4. **`admin-server/scripts` 代码生成脚本需要规范化**（含新增的 RPC/Swagger 生成脚本），见 `12-scripts-standardization.md`。
5. **`.cursor/rules` 需要跟着 `AGENTS.md` 同步更新**，见 `13-rules-sync-checklist.md`。
6. **开发期执行策略覆盖 `AGENTS.md` 第 6 节的默认口径**：开发阶段 AI 可以直接跑 `generate-*.sh`、直接执行本地/开发库 SQL，不必停下来等用户；只有真正的产品/业务判断才停下来问。完整规则见 `10-dev-execution-and-review-points.md`，**这是一份项目级别的临时澄清，不是对 `AGENTS.md` 的静默修改**——是否固化回 `AGENTS.md`，留到 Phase 3 收尾时问用户。
7. **时间线约 14 周（约 3 个月）**，不追求 1 个月内做完全部。

## 3. 三阶段结构与时间线

```
Phase 1  单体内部加固           Week 1-5    admin-server 内部改造，不新增部署单元
Phase 2  微服务拆分             Week 6-12   6 个独立部署单元：gateway + 5 个 rpc 服务
Phase 3  可观测性 + API 文档 + CD  Week 13-14+  Telemetry、结构化日志、Swagger、生产 compose
```

**关键前提**：Phase 1 产出的 `internal/domain/<domain>` 包边界从一开始就按 Phase 2 最终的 5 个服务分组组织（见 `01-architecture-target.md` A.2），不是先按 9 域做、之后再推翻重做——Phase 1 是 Phase 2 的地基，不是过渡阶段。

### Phase 1（Week 1-5，单体加固）

- **Week 1（地基）**：`Repository.Transact`/`withSession`、`registry.Transact`、中间件收窄 + 真正的 Wire 装配、密钥管理（`conf.UseEnv()`）、CI/Docker 骨架起步、自建 `admin-mcp` 工具（越早可用，后续 13 周越受益）。对应 `01`~`03`、`09` 前半部分、`22`。
- **Week 2**：IAM + Chat 联合改造（含 `user_create_logic.go` 修复：改造成 `UserDomainService.CreateUser` 并包一层 `Repository.Transact`、`FindPage` 全量拉取改 `FindChunk` 分批、初始化聊天数据改异步尽力而为）+ Task 域补测试。对应 `04`、`05`。
- **Week 3**：简单域改造（blog/video/sdk → monitoring/system/misc）+ 全仓库扫尾核查跨域越界 import。对应 `06`、`07`。
- **Week 4-5（加固）**：事务审计（确认所有 ≥2 表写方法都过了 `registry.Transact`）、测试覆盖补漏、CI 集成测试跑绿、`golangci-lint` 扩大范围、Phase 1 文档收尾。对应 `08`、`09`。

### Phase 2（Week 6-12，微服务拆分）

6 个独立部署单元：`gateway`（HTTP 唯一入口，无状态）+ `iam-rpc`（iam+system+monitoring+misc）+ `content-rpc`（blog+video）+ `chat-rpc`（chat，含 hub）+ `task-rpc`（task）+ `sdk-rpc`（sdk）。拆分顺序按风险从低到高：

- **Week 6-7**：task-rpc（最小风险，跑通 `goctl rpc` 脚手架、`generate-rpc.sh`、数据库拆分、Dockerfile、docker-compose 全套机制）。
- **Week 8**：sdk-rpc。
- **Week 9-10**：chat-rpc（含 WebSocket ↔ gRPC 双向流桥接）。
- **Week 11**：content-rpc。
- **Week 12**：iam-rpc（放最后，此时能根据前 4 次拆分的真实使用情况确定 `BatchGetUserProfiles` 等 RPC 接口的形状）。

对应 `15`~`18`。

### Phase 3（Week 13-14+，可观测性 + API 文档 + CD）

全部 6 个服务接入 `Telemetry` + 结构化 JSON 日志（`trace_id`/`span_id`/`service`/`user_id`）；`generate-swagger.sh` 产出 gateway 的 OpenAPI 文档；docker-compose 生产切换 + 按服务 CI 镜像构建；`AGENTS.md`/`.cursor/rules` 做一次实质性重写（不是小修）。对应 `19`~`21`。

## 4. 如何使用这套文档

- 每份 Phase 文档（`02` 起）开头写"前置依赖"（读哪些文档、代码库要处于什么状态才能开始），结尾写"完成的定义"（`go build ./...` 通过、测试通过、冒烟步骤）。
- `01-architecture-target.md` 是 Phase 1 的技术总纲，`02`~`09` 是它的可执行拆解，遇到需要复核设计决策的地方回指 `01`，不要在分篇文档里重新推导一遍。
- `15-service-boundaries.md` 是 Phase 2 的技术总纲，作用与 `01` 对 Phase 1 类似。
- 不要跨 Phase 并行改动；每个 Phase / 每个子任务结束都要跑 `go build ./...`（涉及运行时行为的改动，额外做人工冒烟）。
- 任何不确定是否属于"可以直接做"还是"必须停下来问用户"的判断，先查 `10-dev-execution-and-review-points.md`，拿不准就问，不要逢清单必停也不要逢清单都不停。
- **`admin-server/docs/progress.md` 和仓库根目录的 `docs/后端开发进度.md` 是两份不同的记录，不要互相替代**：`progress.md` 是本轮重构 Phase 1-3 期间的过程记录（阶段/周次/关键决策），`docs/后端开发进度.md` 是 `AGENTS.md` 第 7 节"完成的定义"要求的、跨越整个项目生命周期的功能级进度索引。本轮重构过程中只更新 `progress.md`；但每当 Phase 1-3 里的某个改造让某个业务模块的功能行为发生了实质变化（不是纯内部重构，比如某个接口的返回结构变了、某个功能从同步变异步），按 `AGENTS.md` 第 7 节的要求，同步更新 `docs/后端开发进度.md` 对应条目——这是既有规则，本轮不豁免。

## 5. 全部 23 篇文档索引

| 文件 | 一句话用途 | 状态 |
|---|---|---|
| `00-refactor-overview.md` | 总纲，覆盖 Phase 1-3 全貌，本文档 | 已产出 |
| `01-architecture-target.md` | Part A 技术决策正文：事务方案、领域服务分层、Wire 中间件、密钥管理 | 已产出 |
| `02-transactions-and-uow.md` | `Transact`/`withSession` 实现细节 + 逐 Model 文件清单（可执行版） | 已产出 |
| `03-wire-and-middleware.md` | 中间件收窄 + Wire provider 化（可执行版） | 已产出 |
| `04-domain-iam-chat.md` | IAM + Chat 联合改造，含 `user_create_logic.go` 修复 | 已产出 |
| `05-domain-task.md` | Task 域：对齐 `registry.Transact` + 补测试 | 已产出 |
| `06-domain-blog-video-sdk.md` | blog/video/sdk 三个简单域改造 | 已产出 |
| `07-domain-monitoring-system-misc.md` | monitoring/system/misc 三个最薄的域改造 | 已产出 |
| `08-testing-strategy.md` | 测试范围、`sqlmock` 约定、集成测试套件 | 已产出 |
| `09-ci-cd-and-deployability.md` | Phase 1 版 Dockerfile/docker-compose/CI 骨架 | 已产出 |
| `10-dev-execution-and-review-points.md` | 开发期直接执行 vs 事后 review vs 真正停下的边界 | 已产出 |
| `11-descoped.md` | 本轮之外仍不做的事项（K8s、多 module、etcd 服务发现等） | 已产出 |
| `12-scripts-standardization.md` | `scripts/` 规范化（含 sqlgen.exe 清理、`generate-rpc.sh`/`generate-swagger.sh`） | 已产出 |
| `13-rules-sync-checklist.md` | `AGENTS.md` ↔ `.cursor/rules` 同步清单 | 已产出 |
| `14-production-deployment-checklist.md` | 持续追加的上线部署清单 | 已产出 |
| `15-service-boundaries.md` | Part B.1-B.2：6 服务划分 + 数据归属 + 拆分理由 | 已产出 |
| `16-rpc-conventions.md` | Part B.3：proto 风格、`generate-rpc.sh`、zrpc 接入约定 | 已产出 |
| `17-async-eventing.md` | Part B.4：`TaskCallback` + Redis Streams + 同步/异步判断规则表 | 已产出 |
| `18-service-extraction-runbook.md` | 可重复执行 5 次的拆分 checklist + 每次的差异附录 | 已产出 |
| `19-observability.md` | Part C.1：`Telemetry` 配置 + 结构化日志规范 | 已产出 |
| `20-api-docs-generation.md` | Part C.2：`generate-swagger.sh` 流程 | 已产出 |
| `21-cd-and-deployment.md` | Part C.3：镜像构建/推送、compose 生产切换 | 已产出 |
| `22-admin-mcp-tool.md` | Part A.8：自建 `admin-mcp` 工具设计（封装 generate 脚本、项目约定查询、进度查询） | 已产出 |
| `progress.md` | 贯穿 Phase 1-3 的唯一进度记录，不分叉 | 已产出（种子条目） |

## 6. 明确不做的事（详见 `11-descoped.md`，此处先列最重要的几条）

- 不上 Kubernetes——docker-compose 对独立维护者场景足够。
- 不做多 module monorepo——单 module、多 `main` 二进制。
- 不引入 etcd 服务发现——Phase 2 初期用静态 zrpc target。
- 不引入 Kafka/RabbitMQ——跨服务异步场景全部用已有的 Redis Streams。
- 不做日志聚合系统（ELK/Loki）——结构化 JSON 日志先落地，聚合是后续可选项。
- CI 不设全局覆盖率百分比门槛。

## 7. 当前进度

本次交付：`00`~`22` 全部 23 篇任务书文档 + `progress.md` 种子条目，一次性完成，`internal/` 下没有任何代码改动。从下一次会话开始，按 Phase 1 Week 1（见第 3 节）实际动代码，每个子任务完成后回来更新 `progress.md`（追加条目，不要重写），涉及部署动作的同步维护 `14-production-deployment-checklist.md`。
