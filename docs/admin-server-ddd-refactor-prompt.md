# admin-server 领域分层重构任务书（DDD-lite，单体不拆微服务）

> 本文档是可直接执行的任务说明。执行者（Cursor）在改动前应完整阅读一遍，按 Phase 顺序推进，每个 Phase 结束都要跑通 `go build ./...`（Phase 3/4 额外要求人工冒烟测试），不要跨 Phase 并行改动。

## 0. 背景与目标

`admin-server`（go-zero 单体应用，MySQL 单库）当前 504 个 .go 文件、约 4.28 万行代码，`internal/handler`、`internal/logic`、`internal/repository`、`internal/model` 均按 ~44-46 个 API 模块扁平组织，没有领域边界，`internal/svc/servicecontext.go` 有持续膨胀风险。

**目标**：不引入微服务、不新增部署单元、不引入服务间 RPC，仅在现有单体进程内，按 9 个业务领域重新组织代码边界；只在真正有业务复杂度的 IAM（RBAC 权限解析）和 Task（任务调度）两个域引入领域服务层（`internal/domain/`），其余 7 个域保持简单的 repository+logic，不套用聚合根/值对象/领域事件等重型 DDD 战术模式。

**非目标**：不拆分数据库、不拆分部署进程、不引入 gRPC/HTTP 域间调用、不要求同时重写 `internal/task/scheduler.go` 内部实现（squirrel 直连改造为通过 repository 的调用是可选的后续优化，不在本次范围）。

## 1. 领域边界（9 个域，对应 api/admin.api 现有 group）

| 领域 | 现有 group（旧名） | 说明 |
|---|---|---|
| **iam** | auth, user, role, permission, permission_api, permission_menu, department, menu, api, user_role, role_permission | RBAC，含权限解析领域服务 |
| **blog** | tag, article(blog_article), article_audit(audit), friend_link, social_info, public(public_blog) | 博客/CMS |
| **video** | video, m3u8, collect(video_collect), public(public_video) | 视频 |
| **chat** | chat, group(chat_group), message(chat_message) | 聊天（含 internal/hub/chathub.go） |
| **sdk** | sdk, public(sdk_public) | 外部 API Key/调用日志 |
| **task** | task, public(task_public) | 任务调度，含领域服务 |
| **monitoring** | metric, metric_admin, operation_log, login_log, audit_log, performance_log, monitor | 监控/日志 |
| **system** | dict_type, dict_item, dict, config, file, notice, notification | 系统配置类 |
| **misc** | demo, daily_short_sentence, ping, public | 无明确领域归属的工具端点，允许保留在顶层，不强求归入以上 8 个域 |

## 2. 目标目录结构

```
internal/
├── handler/<domain>/<module>/     # goctl 生成，来自 admin.api 的 group: <domain>/<module>
├── logic/<domain>/<module>/       # goctl 生成骨架 + 手写方法体
├── repository/
│   ├── repository.go              # 不变：聚合 Repository 结构体 + BuildSources()，顶层公共基础设施
│   ├── sql_conn.go                # 不变：连接管理
│   ├── cache_conf.go              # 不变
│   └── <domain>/*.go               # 按域拆分的 39 个手写 repository 文件，package 名改为 <domain>
├── model/
│   └── <domain>/*.go               # goctl 生成 + 手写扩展，package 名改为 <domain>
├── domain/
│   ├── iam/permission_resolver.go  # 新增：从 permissionmiddleware.go 提取的权限解析领域服务
│   └── task/                        # internal/task/* 整体迁移到这里（scheduler.go, executors/, notifier.go）
├── svc/servicecontext.go           # 精简：删除 7 个具名 repository 字段，保留 Repository 共享句柄 + 中间件字段不变
├── middleware/, consts/, dict/, types/, config/, hub/
└── interfaces/task.go
```

## 3. api/admin.api group 命名规则

现有 44 个 `@server(group: xxx)` 全部改为 `group: <domain>/<module>`（映射见第 1 节表格，模块名如与域名重复则去重，例如 `blog_article` → `blog/article` 而非 `blog/blog_article`）。

**⚠️ 执行前必须先做验证 spike，不要直接对全部 44 个 group 批量改名：**

1. 安装与 `internal/handler/routes.go` 文件头注释一致版本的 goctl（当前标注 `goctl 1.9.2`）。
2. 复制一份 `admin.api` 到 scratch 文件，只把其中一个 group（建议 `ping` → `misc/ping`）改成嵌套形式，运行 `goctl api go -api <scratch>.api -dir <scratch目录>`，检查：
   - 是否生成了 `handler/misc/ping/`、`logic/misc/ping/` 嵌套目录；
   - 生成的 `routes.go` 里 import 别名规则（是否取路径最后一段作为别名，是否会和其他域下同名模块冲突，例如 `iam/permission` 和其他域如果也有 `permission` 模块）。
3. 如验证通过：按第 1 节映射表对全部 44 个 group 一次性改名，重新生成。
4. 如验证有别名冲突或行为不符预期：改用 fallback —— 要么保持现有扁平唯一 group 名不变、生成后手工把 handler/logic 目录挪进 `<domain>/` 下（放弃利用 goctl 原生嵌套能力，但风险最低），要么用 `iam_auth` 这种下划线拼接的扁平但带域前缀的 group 名（不生成嵌套目录，但保留可读性和未来迁移的可能）。**把最终选择的方案记录进本文档所在目录（更新本文件或在 docs/ 下新增一条记录）。**

`scripts/generate-api.sh`、`scripts/generate-model.sh`（已支持 `-d <dir>`）、`.template/*.tpl` 均为通用模板，未硬编码扁平路径假设，预期不需要改动；`scripts/sqlgen` 不生成 repository 代码，只需要把 AGENTS.md 里 `-group` 参数的说明改为要求传 `<domain>/<module>` 形式。

## 4. Repository / Model 层拆分

- `internal/model/*.go`（约 50 个文件，goctl 生成 `*_gen.go` + 手写扩展）按第 1 节映射迁移到 `internal/model/<domain>/`，每个目录下 `package model` 改为 `package <domain>`。**不要**在这一步顺带重命名类型（如把 `AdminUserModel` 改成 `UserModel`）——保持类型名不变，只挪包路径，让这一步保持纯机械、可 diff review。
- `internal/repository/*_repository.go`（39 个手写文件，`repository.go`/`sql_conn.go`/`cache_conf.go` 除外）按第 1 节映射迁移到 `internal/repository/<domain>/`，`package repository` 改为 `package <domain>`。
- `repository.go`（留在顶层 `internal/repository` 包）需要更新对 model 类型的引用，从 `model.AdminXModel` 改为 `<domain>.AdminXModel`（导入对应的 `internal/model/<domain>` 包）。它本身不属于任何单一领域，继续作为跨域的连接/model 注册表基础设施包。
- 全代码库中对 `repository.NewXRepository(...)` 与 `model.AdminXModel` 的调用点，机械替换为 `<domain>.NewXRepository(...)` / `<domain>.AdminXModel`。

## 5. 领域服务层（仅 IAM、Task 两个域）

- 新建 `internal/domain/iam/permission_resolver.go`：把 `internal/middleware/permissionmiddleware.go` 里内联的权限解析逻辑（当前直接调用 `repository.NewApiRepository`、`NewUserRoleRepository`、`NewRolePermissionRepository`、`NewPermissionApiRepository` 并手写权限集合计算，约 90 行）提取为：
  ```go
  type PermissionResolver struct{ repo *repository.Repository }
  func (r *PermissionResolver) CanAccess(ctx context.Context, userID uint64, method, path string) (bool, error)
  ```
  内部改用 `iam.NewApiRepository(r.repo)` 等新路径。`permissionmiddleware.go` 改为薄适配层：构造 resolver → 调 `CanAccess` → 转换成 HTTP 响应。这是行为保持不变的纯提取，不要在这一步顺带修改权限判定逻辑本身。
- 把 `internal/task/`（`scheduler.go`、`executors/`、`notifier.go`）整体移动到 `internal/domain/task/`，更新 `internal/svc/servicecontext.go` 里的 import 路径。**不要**顺带把 scheduler 内部直连 squirrel 的 SQL 改成走 repository 接口——这是可选的后续优化，本次只做目录迁移，避免和纯重构目标混在一起。
- 其余 7 个域（blog/video/chat/sdk/monitoring/system/misc）**不要**建 `internal/domain/<x>/` 包，logic 文件继续直接调用 `repository/<domain>` 的构造函数，只是目录换了地方，不引入领域服务这层。

## 6. ServiceContext 精简

`internal/svc/servicecontext.go` 现有 7 个具名 repository 字段（`BlogTagRepository`、`BlogArticleRepository`、`BlogArticleTagRepository`、`BlogArticleAuditRepository`、`BlogFriendLinkRepository`、`BlogSocialInfoRepository`、`UserRepository`）只被 Blog 域约 20 个 logic 文件使用，是历史遗留的例外写法，其余 127+ 处 logic 文件都是方法内内联 `repository.NewXRepository(l.svcCtx.Repository)` 构造。

- **删除**这 7 个具名字段。
- 把 Blog 域这约 20 个 logic 文件的调用方式改成和其它域一致的内联构造：`repo := blog.NewArticleRepository(l.svcCtx.Repository)`。
- `Repository *repository.Repository` 字段保持不变，作为唯一共享句柄。
- 11 个 Middleware 字段维持扁平结构，**不要**引入嵌套 struct（如 `svcCtx.Middleware.Auth`）——它们都在启动时构造一次、被路由装配直接引用，符合 go-zero 惯例，嵌套会增加不必要的偏离，不做这个改动。
- 之后新增领域**禁止**再往 `ServiceContext` 加具名 repository 字段，统一走内联构造——这条约束要写进第 8 节要更新的文档里。

## 7. 迁移执行顺序（保证每步可编译）

不要按"域"一次性做完 handler→logic→repository→model 全链路再换下一个域；而是按"层"从依赖图底部往上做，每层结束都跑 `go build ./...` 检查点：

1. **Phase 0 验证 spike**：按第 3 节做 goctl 嵌套 group 验证，确定最终 group 命名方案；同时确认本文档第 9 节要改的文档列表准确。
2. **Phase 1 Model 层拆分**：先做 model（依赖图最底层），改完后全仓库搜索 `model\.` 引用点做机械替换，编译通过为止。
3. **Phase 2 Repository 层拆分**：repository.go 更新引用 + 39 个文件迁移 + 全部调用点替换，编译通过为止。
4. **Phase 3 api/admin.api 改名 + 重新生成 handler/logic**：按验证过的方案批量改 44 个 group，重新生成后，逐个把旧的手写 Logic 方法体从旧扁平目录搬到新生成的骨架文件里（goctl 只会在新路径生成骨架，不会带走旧文件里的手写逻辑，因此这一步是"复制方法体 + 删旧目录"，不是原地改名，要文件级 diff 核对不要漏搬）。搬完确认编译通过，并人工冒烟测试登录、一个受权限保护的接口、一个公开接口（如 `/api/v1/ping`、`/blog` 相关公开页面）。
5. **Phase 4 IAM/Task 领域服务提取**：按第 5 节做，做完要跑通登录、一次权限校验通过、一次权限校验被拒绝、一次任务调度执行的端到端验证。
6. **Phase 5 ServiceContext 精简**：按第 6 节做，做完跑 Blog 域 CRUD 和公开博客页面（前端可见）验证。
7. **Phase 6 文档同步**：见第 9 节。

## 8. 前后对照示例

**简单域（system/dict）：**
```
之前:
internal/handler/dict/dictgethandler.go        package dict
internal/logic/dict/dictgetlogic.go            package dict
internal/repository/dict_type_repository.go    package repository (type DictTypeRepository)
internal/model/admindicttypemodel.go           package model (type AdminDictTypeModel)
api/admin.api:  group: dict_type / dict_item / dict

之后:
internal/handler/system/dict/dictgethandler.go       package dict（basename/包名不变，只是嵌套路径变了）
internal/logic/system/dict/dictgetlogic.go           package dict
internal/repository/system/dict_type_repository.go   package system (type DictTypeRepository)
internal/model/system/admindicttypemodel.go          package system (type AdminDictTypeModel)
api/admin.api:  group: system/dict_type / system/dict_item / system/dict
```

**复杂域（iam）：**
```
之前:
internal/middleware/permissionmiddleware.go 内联调用
    repository.NewApiRepository / NewUserRoleRepository / NewRolePermissionRepository / NewPermissionApiRepository
    + ~90 行手写权限集合计算逻辑
internal/repository/user_repository.go             package repository
internal/model/adminusermodel.go                    package model

之后:
internal/domain/iam/permission_resolver.go
    type PermissionResolver struct{ repo *repository.Repository }
    func (r *PermissionResolver) CanAccess(ctx, userID, method, path) (bool, error)
    内部改用 iam.NewApiRepository(r.repo) 等
internal/middleware/permissionmiddleware.go
    resolver := iamdomain.NewPermissionResolver(m.svcCtx.Repository)
    ok, err := resolver.CanAccess(ctx, user.UserID, method, path)
internal/repository/iam/user_repository.go             package iam
internal/model/iam/adminusermodel.go                    package iam
api/admin.api: group: iam/auth, iam/user, iam/role, iam/permission, iam/permission_api,
               iam/permission_menu, iam/department, iam/menu, iam/api, iam/user_role, iam/role_permission
```

## 9. 需要同步更新的文档

改动落地后，必须同步更新（根目录 CLAUDE.md 要求 AGENTS.md 与 .cursor/rules/*.mdc 保持同步，两处都要改）：

- `AGENTS.md`：目录结构图、脚手架流程（`-group` 参数要求写成 `<domain>/<module>`）、后端规范里 Repository/Model 路径命名表。
- `.cursor/rules/00-workflow.mdc`、`.cursor/rules/10-go-code-style.mdc`：同步以上变更。
- `admin-server/README.md`：目录结构图更新为按域分层。
- `docs/后端开发进度.md`：按仓库自身"完成的定义"要求补一条本次重构的记录（做了什么、为什么、关键决策，如最终选定的 group 命名方案）。
- `.template/README.md`（如存在对旧扁平路径的说明）一并检查。

## 10. 风险与注意事项

- goctl 嵌套 group 生成行为未经验证，Phase 0 spike 是强制前置步骤，不能跳过直接对 44 个 group 批量改名。
- Phase 3（迁移手写 Logic 方法体到新生成骨架）是全流程里唯一非机械、需要人工逐文件核对的一步，约 46 个旧目录，要留足 review 时间，避免漏搬或搬错方法体。
- `pkg/audit/audit.go` 当前直接依赖 `internal/repository`、`internal/svc`、`internal/model`，且其审计类型（`permission_assign`、`role_change`）明显是 IAM 域概念，却放在通用 `pkg/` 目录下——本次不要求处理，但建议记录为后续可选清理项，不要在本次重构里顺带改。
- 全程不引入新的部署单元、不新增服务间网络调用、不拆数据库——任何偏离这一约束的改动方案都应该先暂停并确认，而不是自行扩大范围。
