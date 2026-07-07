# admin-server 维护导航

> 配合 DDD-lite 重构后的目录结构。找代码时先确定**业务域**，再进对应子目录。

## 9 个业务域速查

| 域 | 路径前缀 | 典型功能 |
|----|----------|----------|
| **iam** | `internal/{handler,logic,repository,model}/iam/` | 登录、用户、角色、权限、菜单、部门、接口 |
| **blog** | `.../blog/` | 文章、标签、审核、友链、社交信息、公开博客页 API |
| **video** | `.../video/` | 视频 CRUD、M3U8、采集、公开视频页 |
| **chat** | `.../chat/` | 聊天、群组、消息、WebSocket（`handler/chat/chatwshandler.go`） |
| **sdk** | `.../sdk/` | API Key、接口绑定、调用日志、公开 SDK 上传 |
| **task** | `.../task/` + `internal/domain/task/` | 异步任务 CRUD、调度器、Excel 导出执行器 |
| **monitoring** | `.../monitoring/` | 监控、指标、操作/登录/审计/性能日志 |
| **system** | `.../system/` | 配置、字典、文件、公告、通知 |
| **misc** | `.../misc/` | ping、demo、每日一句、公共字典 |

## 「我要改 X」决策树

```
改登录/JWT/用户/角色/权限/菜单？
  → iam/（权限判定核心：internal/domain/iam/permission_resolver.go）

改博客文章/标签/友链？
  → blog/

改视频或 M3U8？
  → video/

改聊天或 WebSocket？
  → chat/（Hub：internal/hub/chathub.go）

改 SDK 对外开放？
  → sdk/

改异步任务/定时调度？
  → task/（领域服务：internal/domain/task/scheduler.go）

改日志/监控/指标？
  → monitoring/

改字典/配置/文件/公告？
  → system/

改健康检查/demo？
  → misc/
```

## 标准调用模式（重构后统一）

**Repository 构造**（Logic 层内联，禁止再往 ServiceContext 加具名字段）：

```go
import blogrepo "postapocgame/admin-server/internal/repository/blog"

repo := blogrepo.NewArticleRepository(l.svcCtx.Repository)
```

**Model 引用**（按域分包，包名 = 域名）：

```go
import "postapocgame/admin-server/internal/model/blog"

article := &blog.BlogArticle{...}
```

**新增 CRUD 模块**（脚手架 group 格式）：

```bash
# 用户执行
./scripts/generate-sql.sh -group blog/coupon -name coupon
# 生成：handler/logic 在 blog/coupon/，model 在 model/blog/，repository 在 repository/blog/
```

## RBAC 改动检查清单

动权限相关功能时，按顺序检查：

1. `internal/domain/iam/permission_resolver.go` — 权限判定逻辑（`CanAccess`）
2. `internal/middleware/permissionmiddleware.go` — 薄适配层，一般不需改业务
3. `internal/repository/iam/` — 用户角色、角色权限、权限-接口关联
4. `api/admin.api` 中对应 `@server` 的 `middleware` 声明顺序

## api group 命名

`admin.api` 中 group 格式为 `<domain>/<module>`，例如：

- `iam/user`、`blog/article`、`system/dict_type`、`misc/ping`

goctl 会生成嵌套目录 `internal/handler/<domain>/<module>/`。

## 关键文件索引

| 用途 | 路径 |
|------|------|
| API 真源 | `api/admin.api` |
| 路由注册 | `internal/handler/routes.go`（goctl 生成） |
| 自定义路由 | `internal/handler/custom_routes.go`（WebSocket、静态文件） |
| DI 容器 | `internal/svc/servicecontext.go`（仅 `Repository` 句柄 + 横切组件） |
| 跨域 Model 注册 | `internal/repository/repository.go` |
| 任务调度 | `internal/domain/task/scheduler.go` |
| 权限解析 | `internal/domain/iam/permission_resolver.go` |

## 相关文档

- Spike 结论：[`admin-server-phase0-goctl-spike.md`](admin-server-phase0-goctl-spike.md)
- 重构任务书：[`admin-server-ddd-refactor-prompt.md`](admin-server-ddd-refactor-prompt.md)
