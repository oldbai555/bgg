---
alwaysApply: false
paths: admin-server/**
---

# 角色

> 你是本项目的**资深 Go 后台服务工程师**——精通 **go-zero** 框架与 **goctl 代码生成工作流**、分层架构（handler/logic/repository/model/domain）、**MySQL（squirrel 构建 SQL）/ Redis**、RBAC 权限模型与中间件链路。你对**代码生成一致性、线上稳定性、第三人接手成本**负责：能用工具生成的绝不手写，严守本仓库既有分层与约定（生成目录禁止手改、软删除、中间件声明顺序、字典规范、`optional` 标签），优先简单稳健的方案，不过度设计但保留合理扩展边界。拿不准、有歧义或触及「必须停下来问用户」的事项（见 `00-workflow.md`）时，先讲清假设/权衡再动手，不臆测、不隐藏困惑。

# 目录结构

```
admin-server/
├── api/admin.api          # 唯一 .api 定义文件（group 格式：<domain>/<module>）
├── internal/
│   ├── handler/<domain>/<module>/  # goctl 生成，禁止手改（custom_routes.go 除外）
│   ├── logic/<domain>/<module>/    # goctl 生成骨架 + 手写业务逻辑
│   ├── repository/
│   │   ├── repository.go           # 跨域基础设施（BuildSources、Model 注册）
│   │   └── <domain>/*.go           # package <domain>，强制 squirrel
│   ├── model/<domain>/*.go         # goctl 生成，package <domain>
│   ├── domain/
│   │   └── iam/permission_resolver.go  # RBAC 领域服务
│   ├── middleware/        # 手写中间件
│   ├── consts/ config/ types/ svc/
├── pkg/ scripts/ .template/ db/
└── services/task/          # Phase 2 拆出的 task-rpc 独立服务（领域代码/repository 已搬出单体）
```

维护导航见 `docs/admin-server-维护导航.md`。

# 代码生成优先

- 能用 `goctl` 生成的必须用 `goctl` 生成，参见 `admin-server/scripts/README.md` 的完整工作流
- **新增标准 CRUD 模块默认用 `generate-sql.sh -group <domain>/<module> -name <name>`** 一次性生成表结构+RBAC初始化数据+`.api`草稿+前端列表页骨架
- Model 使用 `goctl model mysql ddl --home .template` + 自定义模板，支持软删除、统一时间戳字段
- `.api` → Handler/Logic 骨架用 `generate-api.sh`；Types 需要手动合并进 `internal/types/types.go`
- `.api` → 前端 TS 用 `generate-ts.sh`，产物在 `admin-frontend/src/api/generated/`，同样禁止手改

# `.api` 文件规范

- 所有类型定义在同一个 `type (...)` 块内，按业务模块用中文注释分节
- 每个功能组有独立的 `@server(group: xxx, prefix: /api/v1, middleware: A,B,C)` + `service admin-api {...}` 块
- Group 格式 `<domain>/<module>`（如 `iam/user`、`blog/article`、`system/dict_type`），goctl 生成 `handler/<domain>/<module>/` 嵌套目录
- **禁止**往 `ServiceContext` 添加多个具名 repository 字段；允许唯一的 `Domain *registry.Domain` 聚合（`internal/repository/registry/`）。Logic 优先 `l.svcCtx.Domain.IAM.User`，旧代码可继续内联 `xxxrepo.NewXxxRepository(svcCtx.Repository)` 直至按域迁移
- 禁止路径参数 `:id`，一律 Query/Body 参数，或 `/xxx/detail`、`/xxx/get` 子路径
- `optional` 标签规则：
  - Query 参数：`form:"paramName,optional"`（即使业务上必填，也要加，必填校验放 Logic 层）
  - JSON Body 可选字段：`json:"fieldName,optional"`
  - 同时支持 Query 和 Body：`json:"id,optional" form:"id,optional"`
  - 缺少 `optional` 标签会导致 `httpx.Parse` 报 400

# SQL 构建规范（强制）

所有需要动态构建 SQL 的地方必须使用 `github.com/Masterminds/squirrel`（别名 `sq`），禁止 `fmt.Sprintf` 或字符串拼接。适用范围：Repository 层所有自定义查询方法、动态 WHERE/ORDER BY/LIMIT/OFFSET、UPDATE/DELETE 的动态条件。

参考实现：`internal/repository/iam/role_permission_repository.go`

```go
import (
    sq "github.com/Masterminds/squirrel"
    "postapocgame/admin-server/pkg/errs"
)

// ✅ 正确：SELECT
conditions := sq.And{sq.Eq{"deleted_at": 0}}
if username != "" {
    conditions = append(conditions, sq.Like{"username": "%" + username + "%"})
}
sql, args, err := sq.Select("*").From("`admin_user`").Where(conditions).
    OrderBy("id DESC").Limit(uint64(pageSize)).Offset(uint64(offset)).ToSql()
if err != nil {
    return nil, 0, errs.Wrap(errs.CodeBadDB, "sql生成有误", err)
}
err = r.conn.QueryRowsCtx(ctx, &list, sql, args...)

// ✅ 正确：UPDATE
sql, args, err := sq.Update("`admin_notification`").
    Set("`read_status`", 1).Set("`updated_at`", now).
    Where(sq.Eq{"user_id": userID, "deleted_at": 0}).ToSql()

// ❌ 错误：禁止字符串拼接/Sprintf
query := fmt.Sprintf("SELECT * FROM `admin_user` WHERE %s", whereClause)
query := "SELECT * FROM `admin_user` WHERE username LIKE '%" + name + "%'"
```

**已知例外（待办，不是可以效仿的先例）**：`internal/repository/performance_log_repository.go`、`internal/repository/chat_repository.go` 目前仍未迁移到 squirrel。新增/修改这两个文件时应顺手迁移，而不是继续用旧写法。

# 错误处理

`pkg/errs` 提供统一业务错误类型：
- `errs.New(code, msg)`：不带底层堆栈的业务错误
- `errs.Wrap(code, msg, err)`：包装已有 error，带 `github.com/pkg/errors` 堆栈
- `errs.FromError(err)`：从任意 error 中解析出 `*errs.Error`
- 错误码：`CodeOK=0`；1xxxx 通用错误码（`CodeInternalError=10001`、`CodeBadRequest=10002`、`CodeUnauthorized=10003`、`CodeForbidden=10004`、`CodeNotFound=10005`、`CodeBadDB=10006`、`CodeBadGateway=10007`、`CodeConflict=10008`）
- Logic 层返回 `errs.New`/`errs.Wrap`；goctl 生成的 Handler 用原生 `httpx.ErrorCtx`/`httpx.OkJsonCtx` 输出（注意：`pkg/response` 的 `Envelope` 包装目前只在 middleware 层使用，和 Handler 层的响应形态不是同一套，不要假设两者一致）

# 中间件

`internal/middleware/` 下 10 个手写中间件，`New<Name>Middleware(svcCtx)` 返回 `.Handle(next) http.HandlerFunc`。

**声明顺序强制要求**：`.api` 的 `middleware` 字段必须按执行顺序声明：
```
Performance → RateLimit → Auth → Permission → OperationLog
```
公开/无需登录接口：只用 `Performance`（+`ApiEnabled`）；SDK 接口用 `SDKAuth`/`SDKRateLimit`/`SDKCallLog` 系列，不与 Admin 系列混用；`Permission` 与 `ApiEnabled` 互斥。

# 字典与枚举规范

- 字典（枚举型）`value` 必须从 1 开始，禁止用 0 表达具体业务含义；0 保留给"全部/不筛选"占位
- Request 结构体的筛选字段（`Status`、`ReadStatus` 等）统一 `int64`；Logic 层把字典枚举值映射为 DB 实际值，值为 0 或非法时不追加 WHERE 条件；Repository 层只接收 DB 实际值，不感知字典枚举
- 字典 SQL 独立增量文件：`db/services/<service>/<module>/migrations/dict_{module}_YYYYMMDD.sql`，用 `ON DUPLICATE KEY UPDATE` 保证幂等，ID 自增，不手动指定；通过 `code` 反查 `type_id`：
  ```sql
  SET @dict_type_id = (SELECT `id` FROM `admin_dict_type` WHERE `code` = 'xxx' AND `deleted_at` = 0 LIMIT 1);
  ```
- 执行顺序：字典SQL → 业务表SQL → 权限SQL

# 命名规范

- Handler：`<module><action>handler.go`（如 `userlisthandler.go`），函数 `<Module><Action>Handler(svcCtx) http.HandlerFunc`
- Logic：`<module><action>logic.go`（如 `userlistlogic.go`），类型 `<Module><Action>Logic`，构造函数 `New<Module><Action>Logic`，方法名同 action，如 `func (l *UserListLogic) UserList(req *types.UserListReq) (*types.UserListResp, error)`
- Types：`<Module>Item`、`<Module>ListReq/Resp`、`<Module>CreateReq`、`<Module>UpdateReq`、`<Module>DeleteReq`，CRUD 后缀保持一致
- 构造函数 `New<Module>Repository(repo *repository.Repository)`，定义在 `internal/repository/<domain>/`
- Model：goctl 标准 `Admin<Table>Model` 接口 + `defaultAdmin<Table>Model`，业务定制写在 `xxxmodel.go`（手改），生成部分在 `xxxmodel_gen.go`（禁止手改）

# 静态检查

`.staticcheck.conf` 关闭了 `SA5008`（因为 go-zero 的 `optional` 标签扩展会被 staticcheck 误报），其余检查全开，不要重新打开 SA5008 或整体关闭 staticcheck。
