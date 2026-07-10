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
