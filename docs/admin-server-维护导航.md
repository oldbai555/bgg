# admin-server 维护导航

> 配合 DDD-lite 重构 + Phase 2 微服务拆分后的目录结构。找代码时先确定**业务域**，再判断这个域是否已经拆分成独立服务；Phase 2 五个服务（iam/task/sdk/chat/content）已经全部拆分完成，**gateway 不再持有任何 repository/model/domain，也不直连任何 MySQL**，只剩 `internal/{handler,logic}/<domain>/<module>/` 薄胶水 + 共享 Redis 直连。

## 9 个业务域速查

| 域 | 路径前缀 | 典型功能 |
|----|----------|----------|
| **iam** | gateway 薄胶水 `internal/{handler,logic}/iam/` + `services/iam/`（iam-rpc，2026-07-13 拆分） | 登录、用户、角色、权限、菜单、部门、接口；权限判定核心 `services/iam/internal/domain/iam/permission_resolver.go` |
| **system** | gateway 薄胶水 `internal/{handler,logic}/system/` + `services/iam/`（system 域并入 iam-rpc） | 配置、字典、文件、公告、通知 |
| **monitoring** | gateway 薄胶水 `internal/{handler,logic}/monitoring/` + `services/iam/`（monitoring 域并入 iam-rpc） | 监控、指标、操作/登录/审计/性能日志 |
| **misc** | gateway 薄胶水 `internal/{handler,logic}/misc/` + `services/iam/`（misc 域并入 iam-rpc） | ping、demo、每日一句、公共字典 |
| **blog** | gateway 薄胶水 `internal/{handler,logic}/blog/` + `services/content/`（content-rpc，2026-07-12 拆分） | 文章、标签、审核、友链、社交信息、公开博客页 API |
| **video** | gateway 薄胶水 `internal/{handler,logic}/video/` + `services/content/`（video 域并入 content-rpc） | 视频 CRUD、M3U8、采集、公开视频页 |
| **chat** | gateway 薄胶水 `internal/{handler,logic}/chat/` + `services/chat/`（chat-rpc，2026-07-12 拆分） | 聊天、群组、消息、WebSocket↔gRPC 桥接（`internal/handler/chat/chatwshandler.go`），Hub/领域服务/消费者在 `services/chat/internal/{hub,domain,consumer}/` |
| **sdk** | gateway 薄胶水 `internal/logic/sdk/` + `services/sdk/`（sdk-rpc，2026-07-12 拆分） | API Key、接口绑定、调用日志、公开 SDK 上传，领域服务/repository 在 `services/sdk/internal/{domain,repository}/` |
| **task** | gateway 薄胶水 `internal/logic/task/` + `services/task/`（task-rpc，2026-07-11 拆分） | 异步任务 CRUD、调度器、Excel 导出执行器，领域服务/repository 在 `services/task/internal/{domain,repository}/` |

**iam-rpc 是唯一一个"挂了就全站不可用"的服务**：`AuthMiddleware`/`PermissionMiddleware`/`ApiEnabledMiddleware`/`OperationLogMiddleware`/`PerformanceMiddleware` 五个中间件全部调 `IamRPC`，几乎所有接口都要经过这几个中间件之一。

## 「我要改 X」决策树

```
改登录/JWT/用户/角色/权限/菜单？
  → iam 域已拆分成独立服务 services/iam/（iam-rpc）；gateway 只剩薄胶水 internal/logic/iam/；
    权限判定核心在 services/iam/internal/domain/iam/permission_resolver.go

改配置/字典/文件/公告/通知？
  → system 域已并入 services/iam/（iam-rpc）；gateway 只剩薄胶水 internal/logic/system/

改监控/指标/操作日志/登录日志/审计日志/性能日志？
  → monitoring 域已并入 services/iam/（iam-rpc）；gateway 只剩薄胶水 internal/logic/monitoring/

改健康检查/demo/每日一句？
  → misc 域已并入 services/iam/（iam-rpc）；gateway 只剩薄胶水 internal/logic/misc/

改博客文章/标签/友链/社交信息？
  → blog 域已拆分成独立服务 services/content/（content-rpc）；gateway 只剩薄胶水 internal/logic/blog/

改视频或 M3U8？
  → video 域已并入 services/content/（content-rpc）；gateway 只剩薄胶水 internal/logic/video/；
    M3u8Proxy（纯 HTTP 代理）/VideoCollectOptions（CORS 预检占位）两个端点不涉及域数据，留在 gateway 不接入 RPC

改聊天或 WebSocket？
  → chat 域已拆分成独立服务 services/chat/（chat-rpc）；gateway 只剩薄胶水 internal/logic/chat/ +
    WS↔gRPC 桥接 internal/handler/chat/chatwshandler.go；Hub/领域服务在 services/chat/internal/hub/

改 SDK 对外开放？
  → sdk 域已拆分成独立服务 services/sdk/（sdk-rpc）；gateway 只剩薄胶水 internal/logic/sdk/

改异步任务/定时调度？
  → task 域已拆分成独立服务 services/task/（task-rpc）；调度器在
    services/task/internal/domain/task/scheduler.go
```

## 标准调用模式（Phase 2 五个服务拆分完成后）

**gateway 侧**：Logic 一律调对应服务的 zrpc client，不直连任何数据库：

```go
// gateway 的 internal/logic/<domain>/<module>/xxx_logic.go
rpcResp, err := l.svcCtx.IamRPC.UserList(l.ctx, &iamclient.UserListRequest{...})
// 或 l.svcCtx.TaskRPC / l.svcCtx.SdkRPC / l.svcCtx.ChatRPC / l.svcCtx.ContentRPC
```

**各服务内部**（`services/<name>/internal/`）：仍然是重构后统一的 Repository/Domain 分层，新代码优先用 Domain 聚合：

```go
// services/iam/internal/logic/xxxlogic.go
userRepo := l.svcCtx.Domain.IAM.User
```

**禁止**在 gateway 的 `ServiceContext` 上挂任何 repository/Domain 字段——gateway 只持有 `Redis *redis.Redis` + 5 个 `XxxRPC` client。**禁止**在各服务自己的 `ServiceContext` 上挂多个具名 repo 字段；允许唯一的 `Domain *registry.Domain` 聚合（如 `services/iam/internal/repository/registry/domain.go`）。

**Model 引用**（按域分包，包名 = 域名，在各服务自己的 `internal/model/<domain>/` 下）：

```go
import "postapocgame/admin-server/services/content/internal/model/blog"

article := &blog.BlogArticle{...}
```

**新增 CRUD 模块**（脚手架 group 格式，产出的 handler/logic 骨架进 gateway，业务落在对应已拆分服务时需要手工补 RPC 调用）：

```bash
# 用户执行
./scripts/generate-sql.sh -group blog/coupon -name coupon
# 生成：handler/logic 骨架在 gateway 的 internal/{handler,logic}/blog/coupon/
```

## RBAC 改动检查清单

动权限相关功能时，按顺序检查：

1. `services/iam/internal/domain/iam/permission_resolver.go` — 权限判定逻辑（`CanAccess`）
2. `services/iam/internal/logic/checkpermissionlogic.go` — RPC 入口，包一层 `CanAccess`
3. `internal/middleware/permissionmiddleware.go` — gateway 薄适配层（调 `IamRPC.CheckPermission`），一般不需改业务
4. `services/iam/internal/repository/iam/` — 用户角色、角色权限、权限-接口关联
5. `api/admin.api` 中对应 `@server` 的 `middleware` 声明顺序

## api group 命名

`admin.api` 中 group 格式为 `<domain>/<module>`，例如：

- `iam/user`、`blog/article`、`system/dict_type`、`misc/ping`

goctl 会生成嵌套目录 `internal/handler/<domain>/<module>/`（都在 gateway 里，不管这个域是否已拆分成独立服务——拆分只影响 Logic 内部实现，不影响 Handler 骨架的生成位置）。

## 启动依赖（Google Wire）

组合根依赖由 Wire 编译期注入，**与 goctl 的 `generate-*.sh` 职责分离**：

| 步骤 | 命令 / 文件 |
|------|-------------|
| 改 Provider | 编辑 `internal/wire/providers.go` 或 `wire.go` |
| 重新生成 | `cd admin-server && make wire` |
| 提交产物 | `internal/wire/wire_gen.go`（**禁止手改**）+ 若变动的 `go.mod`/`go.sum` |
| 首次安装 CLI | `make init` |

合并 `wire_gen.go` 冲突时：任选一侧后执行 `make wire` 重生。

Handler/Logic 仍通过 `*svc.ServiceContext` 访问依赖；gateway 的 `ServiceContext` 现在只有 `Redis`/`IamRPC`/`IamCallbackRPC`/`TaskRPC`/`SdkRPC`/`ChatRPC`/`ContentRPC`/一批中间件字段（见下面关键文件索引）。

## 关键文件索引

| 用途 | 路径 |
|------|------|
| API 真源 | `api/admin.api` |
| 路由注册 | `internal/handler/routes.go`（goctl 生成） |
| 自定义路由 | `internal/handler/custom_routes.go`（WebSocket、静态文件） |
| gateway DI 容器 | `internal/svc/servicecontext.go`（`Redis` + 5 个 `XxxRPC` + 横切组件，不持有任何 repository/Domain） |
| gateway 共享 Redis 构造 | `internal/redisconn/redisconn.go` |
| 组合根 Wire | `internal/wire/`（改 Provider 后执行 `make wire`，提交 `wire_gen.go`） |
| iam-rpc 本体 | `services/iam/`（`iam.go` 主入口，同一进程注册 `Iam`/`TaskCallback`/`IamCallback` 三个 gRPC service） |
| iam-rpc 领域 Repo 聚合 | `services/iam/internal/repository/registry/domain.go` |
| iam-rpc 跨域 Model 注册 | `services/iam/internal/repository/repository.go` |
| 任务调度 | `services/task/internal/domain/task/scheduler.go`（task-rpc） |
| 权限解析 | `services/iam/internal/domain/iam/permission_resolver.go` |
| 跨服务回调契约 | `pkg/taskcallback/`、`pkg/iamcallback/`（proto + client，服务端实现都在 `services/iam/internal/server/`） |

## 相关文档

- DDD-lite 分层重构历史背景（Spike 结论 + 重构任务书，原两篇文档已归档删除）：[`docs/changelog/archive-backend.md`](changelog/archive-backend.md) §15
- Phase 2 服务拆分完整过程：`admin-server/docs/progress.md`（按时间顺序的唯一进度记录）
