# Phase 1 单体加固 — 技术决策正文（Part A）

> 本文档是 Phase 1（Week 1-5）的技术总纲，`02`~`09` 是它的可执行拆解版本。分篇文档遇到需要复核设计决策的地方应回指本文档，不要重新推导。本文档不含任务清单和"完成的定义"——那些在对应的可执行版文档里。

## 前置阅读

执行任何一节之前，先确认已经读过：`internal/repository/repository.go`（`Repository` 聚合结构）、`internal/repository/registry/domain.go`（`registry.Domain` 聚合结构）、`internal/domain/iam/permission_resolver.go`（当前唯一的领域服务实现）、`internal/wire/providers.go`（当前 Wire 组合根）。本文档的所有代码示例都是在读过这四个文件之后、对照当前真实代码写出的，不是抽象设计稿。

---

## A.1 事务方案

### 现状

`Repository`（`internal/repository/repository.go`）聚合了 37 个 goctl 生成的 Model 字段（`AdminUserModel`、`BlogArticleModel`……），全部由 `NewRepository` 在启动时用同一个 `sqlx.SqlConn conn` 一次性构造。`Repository` 本身、以及 `registry.Domain`（`internal/repository/registry/domain.go`，`NewDomain(repo *repository.Repository) *Domain` 按 9 个域把 repository 分组）都没有任何事务相关方法。跨表写入目前的处理方式是"不进事务，靠代码注释承认风险"——例如 `internal/repository/blog/blog_article_repository.go` 的 `CreateWithTags`/`UpdateWithTags`：

```go
// 使用 squirrel 手动插入，避免依赖事务 session API
```

这行注释就是本节要解决的问题：go-zero 原生就有事务 session API，只是这个仓库从没用过。

### 方案：复用 go-zero 原生 `sqlc.CachedConn.WithSession` + `sqlx.NewSqlConnFromSession`

不引入 `UnitOfWork` 接口等新抽象。已对照 vendored 源码 `$GOMODCACHE/github.com/zeromicro/go-zero@v1.9.3` 验证以下签名：

```go
// core/stores/sqlc/cachedsql.go:269
func (cc CachedConn) WithSession(session sqlx.Session) CachedConn

// core/stores/sqlx/sqlconn.go:142
func NewSqlConnFromSession(session Session) SqlConn

// core/stores/sqlx/sqlconn.go:37-44 —— SqlConn 接口本身就内嵌了 TransactCtx，
// Repository.DB 字段的类型就是 sqlx.SqlConn，可以直接调用，不需要类型断言
type SqlConn interface {
	Session
	RawDB() (*sql.DB, error)
	Transact(fn func(Session) error) error
	TransactCtx(ctx context.Context, fn func(context.Context, Session) error) error
}
```

每个 goctl 生成的 Model（如 `AdminUserModel`）都是"手改 sibling 文件定义的接口 + `defaultXxxModel` 内嵌 `sqlc.CachedConn`"的结构（见 `internal/model/iam/adminusermodel.go` + `adminusermodel_gen.go`）。给 `sqlc.CachedConn` 换绑 session 拿到的是裸的 `CachedConn` 值，不是 `AdminUserModel` 接口——所以每个手改 sibling 文件要新增一个 `WithSession` 方法，把换绑后的 `CachedConn` 重新包装回 `AdminUserModel` 接口。逐 Model 的机械改法和完整文件清单见 `02-transactions-and-uow.md`；这里只定设计。

### `internal/repository/repository.go` 新增内容

```go
package repository

import (
	"context"
	// ...已有 import
)

// Transact 在单个 MySQL 事务内执行 fn。
// fn 收到的 txRepo 是 r 的克隆：DB 字段与全部 *Model 字段都已经换绑到本次事务的 session 上。
// fn 内部必须只通过 txRepo 访问数据，不能继续闭包引用外层 r —— 否则读写会跳出事务边界。
// 不做嵌套事务检测：调用方需保证不在已经开启的 Transact 内部再次调用 Transact。
func (r *Repository) Transact(ctx context.Context, fn func(ctx context.Context, txRepo *Repository) error) error {
	return r.DB.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		return fn(ctx, r.withSession(session))
	})
}

// withSession 返回一个新的 *Repository，DB 字段与全部 37 个 *Model 字段都换绑到给定的事务 session。
// CacheConf/Redis/BusinessCache 保持不变——事务内不使用查询缓存，这是 go-zero WithSession
// 官方注释里明确警告过的（未提交数据不应该被缓存，也不应该用缓存查未提交数据）。
func (r *Repository) withSession(session sqlx.Session) *Repository {
	return &Repository{
		DB:            sqlx.NewSqlConnFromSession(session),
		CacheConf:     r.CacheConf,
		Redis:         r.Redis,
		BusinessCache: r.BusinessCache,

		AdminUserModel:           r.AdminUserModel.WithSession(session),
		AdminRoleModel:           r.AdminRoleModel.WithSession(session),
		// ... 其余 35 个 *Model 字段逐个调用 .WithSession(session)，
		// 完整字段列表见 02-transactions-and-uow.md
	}
}
```

### `internal/repository/registry/domain.go` 新增内容

```go
package registry

import (
	"context"

	"postapocgame/admin-server/internal/repository"
)

// Transact 在事务内执行 fn。fn 收到的 txDomain 是用换绑过事务 session 的 Repository
// 重新调用 NewDomain 构造出来的 —— 这样做是因为每个 <domain>repo.NewXxxRepository(repo)
// 都是在构造时才把 repo.DB / repo.XxxModel 捕获进内部字段（例如
// internal/repository/iam/user_repository.go 的 userRepository.conn = repo.DB），
// 事务场景下必须用换绑后的 repo 重新构造一遍，不能只换 Repository 本身。
func Transact(ctx context.Context, repo *repository.Repository, fn func(ctx context.Context, txDomain *Domain) error) error {
	return repo.Transact(ctx, func(ctx context.Context, txRepo *repository.Repository) error {
		return fn(ctx, NewDomain(txRepo))
	})
}
```

### 两条使用路径，不要混用（避免循环 import）

`registry` 包要 import 各域的 `internal/domain/<domain>` 包（`Domain` 结构体要引用 `*iamdomain.UserDomainService` 这类类型才能挂领域服务字段，见 A.2），所以**领域服务本身绝不能反过来 import `internal/repository/registry`**——否则 `internal/domain/iam` → `internal/repository/registry` → `internal/domain/iam` 两个包互相 import，编译直接失败。据此明确分两条路径，`04-domain-iam-chat.md` 已经按这个规则写好了范例，本文档后续所有领域服务代码示例都遵守同一规则：

- **领域服务内部**（`internal/domain/<domain>/*.go` 里的方法）：只允许用 `Repository.Transact(ctx, fn)`（上面刚定义的底层原语，`internal/repository` 包，不是 `registry` 包），回调里按需 `iamrepo.NewXxxRepository(txRepo)` 自己构造需要的子仓储，和现有 `PermissionResolver`、`internal/domain/task` 的写法完全一致（它们目前也只 import `internal/repository`，不 import `internal/repository/registry`）。
- **没有专属领域服务的 logic 文件**：直接用 `registry.Transact(ctx, repo, fn)`，回调拿到一个事务内的 `*registry.Domain`，可以像平时一样 `txDomain.IAM.User.XXX(...)` 跨多个已有 repository 操作，不需要现造一个领域服务。

### 使用示例

**没有专属领域服务的 logic 文件**（直接用 `registry.Transact`）：

```go
// 示例：批量调整角色权限的 logic 文件，没有专属领域服务，直接用 registry.Transact
func (l *RolePermissionReplaceLogic) RolePermissionReplace(ctx context.Context, roleID uint64, permIDs []uint64) error {
	return registry.Transact(ctx, l.svcCtx.Repository, func(ctx context.Context, txDomain *registry.Domain) error {
		if err := txDomain.IAM.RolePermission.DeleteByRoleID(ctx, roleID); err != nil {
			return err
		}
		return txDomain.IAM.RolePermission.BatchCreate(ctx, roleID, permIDs)
	})
}
```

**领域服务内部**（只用 `Repository.Transact`，不 import `registry`）：

```go
// internal/domain/iam/user_service.go（示意，完整模板见下方 A.2，含跨域 onboarding 依赖）
func (s *UserDomainService) CreateUser(ctx context.Context, req CreateUserInput) (*iammodel.AdminUser, error) {
	var created *iammodel.AdminUser
	err := s.repo.Transact(ctx, func(ctx context.Context, txRepo *repository.Repository) error {
		user := &iammodel.AdminUser{ /* ... */ }
		if err := iamrepo.NewUserRepository(txRepo).Create(ctx, user); err != nil {
			return err
		}
		created = user
		return nil
	})
	return created, err
}
```

两条路径内部都是同一个 `sqlx.SqlConn.TransactCtx`：任何一步返回 `error` 自动 `ROLLBACK`，`fn` 正常返回 `nil` 才 `COMMIT`，不需要手写回滚逻辑。

### 明确不做

- 不做通用 `UnitOfWork` 接口——`registry.Transact` 已经是够用的抽象，多一层接口只是为了"看起来更 DDD"，没有实际收益。
- 不做分布式事务/Saga——当前阶段单 MySQL 实例足够；Phase 2 拆库之后跨服务一致性走异步事件方案（`17-async-eventing.md`），不是分布式事务。

---

## A.2 领域服务落地位置与分层原则

### 放置位置：不新增 `ServiceContext` 字段

`AGENTS.md` 第 3 节明令禁止往 `ServiceContext` 加具名 repository 字段；领域服务同理不新增字段，而是作为字段挂在 `registry.Domain` 内部各域子结构体上。当前 `registry.Domain`：

```go
type Domain struct {
	IAM        IAMDomain
	Blog       BlogDomain
	Chat       ChatDomain
	SDK        SDKDomain
	Video      VideoDomain
	Task       TaskDomain
	Monitoring MonitoringDomain
	System     SystemDomain
	Misc       MiscDomain
}

type IAMDomain struct {
	User           iamrepo.UserRepository
	Role           iamrepo.RoleRepository
	// ... 11 个 repository 字段
}
```

改造后，`IAMDomain` 这类需要领域服务的子结构体新增服务字段，与现有 repository 字段并列：

```go
type IAMDomain struct {
	User           iamrepo.UserRepository
	Role           iamrepo.RoleRepository
	// ... 原有 11 个 repository 字段不变
	UserService       *iamdomain.UserDomainService // 新增：跨表/跨域编排
	PermissionResolver *iamdomain.PermissionResolver // 已存在，挪进来统一管理
}
```

`registry.NewDomain(repo *repository.Repository) *Domain` 按依赖顺序构造：先构造域内所有 repository 字段，再用这些 repository（或 `repo` 本身）构造该域的领域服务字段。领域服务构造函数继续吃 `*repository.Repository`（跟 `PermissionResolver` 现在的写法一致：`func NewPermissionResolver(repo *repository.Repository) *PermissionResolver`），不吃 `registry.Domain` 自身——避免构造顺序的循环依赖。

### 跨域调用规则

跨域调用必须通过目标域暴露的**窄接口**，禁止跨域直接 import 对方的 `internal/repository/<domain>` 包。当前 `user_create_logic.go` 的写法（`iamrepo` 直接 import `chatrepo`）就是要修掉的反例，完整修复方案见 `04-domain-iam-chat.md`。窄接口的形状（以 Chat 域给 IAM 用为例）：

```go
// internal/domain/chat/onboarding.go —— 目录是 internal/domain/chat，package 名是 chatdomain
// （两者不同，Go 允许；04-domain-iam-chat.md 已经这样落地，本文档跟随同一约定）。
package chatdomain

// Onboarding 是 Chat 域暴露给其他域的窄接口，只暴露"新用户入群/建私聊"这一个能力，
// 不暴露 ChatRepository/ChatUserRepository 本身。
type Onboarding interface {
	InitNewUser(ctx context.Context, userID uint64) error
}
```

IAM 域的 `UserService`（类型 `*iamdomain.UserDomainService`）依赖这个接口类型，而不是 `chatrepo.ChatRepository`；`registry.NewDomain` 构造时把 `Chat.Onboarding`（实现了 `Onboarding`）传给 `IAM.UserService`。这样 IAM 包的 import 列表里不会出现 `internal/repository/chat`。

### 分层判断标准：哪些 logic 文件需要领域服务

不是给 161 个 logic 文件全部套模板。只有满足以下任一条件的方法才需要领域服务，估算约 35-40 个文件符合条件：

1. **跨 ≥2 表写**（哪怕都在同一个 repository 里）——例如 `blogArticleRepository.CreateWithTags`（写 `blog_article` + `blog_article_tag` 两张表）目前完全没有事务保护，是这一类的典型反例。
2. **跨仓储读写**——一个业务方法内需要调用 2 个以上不同的 `XxxRepository`。
3. **跨域**——例如 `user_create_logic.go` 的 `initChatForNewUser` 同时触达 IAM 和 Chat。
4. **非平凡业务规则**——RBAC 授权判定（`PermissionResolver.CanAccess`）、密码/认证、任务状态机（`internal/domain/task/scheduler.go`）。

不满足以上任何一条的方法（单表单仓储的纯读、简单写）继续走 `svcCtx.Domain.X.Y` 直调，**不要**为了架构统一强行包一层领域服务——那是过度工程，会让 Cursor/Claude Code 在后续维护时多一层无意义的间接。

### 复杂域模板：IAM / User

```go
// internal/domain/iam/user_service.go
package iam

import (
	"context"

	"postapocgame/admin-server/internal/repository"
	iamrepo "postapocgame/admin-server/internal/repository/iam"
	chatdomain "postapocgame/admin-server/internal/domain/chat"
	iammodel "postapocgame/admin-server/internal/model/iam"
)

// UserDomainService 承载用户创建等跨表/跨域编排逻辑；单表 CRUD（FindByID/FindPage 等）
// 不搬进这里，继续留在 iamrepo.UserRepository 上，Logic 层直接调用即可。
// 命名与完整实现见 04-domain-iam-chat.md，本文档只展示事务用法。
type UserDomainService struct {
	repo       *repository.Repository
	onboarding chatdomain.Onboarding // 窄接口，不是 chatrepo.ChatRepository，不 import internal/repository/chat
}

func NewUserDomainService(repo *repository.Repository, onboarding chatdomain.Onboarding) *UserDomainService {
	return &UserDomainService{repo: repo, onboarding: onboarding}
}

func (s *UserDomainService) CreateUser(ctx context.Context, user *iammodel.AdminUser) error {
	// 领域服务内部只用 Repository.Transact，不 import internal/repository/registry（见上一节）。
	err := s.repo.Transact(ctx, func(ctx context.Context, txRepo *repository.Repository) error {
		return iamrepo.NewUserRepository(txRepo).Create(ctx, user)
	})
	if err != nil {
		return err
	}
	// 入群/建私聊是尽力而为的异步操作，不在建用户事务内，失败不影响主流程。
	// Phase 1 通过 internal/domain/task 现有调度器异步派发（细节见 04）。
	go func() {
		if err := s.onboarding.InitNewUser(context.Background(), user.Id); err != nil {
			logx.Errorf("初始化新用户聊天数据失败: %v", err)
		}
	}()
	return nil
}
```

`iamrepo.NewUserRepository(repo)` 这类单表调用不受影响，`user_create_logic.go` 之外的其余 IAM logic 文件（如 `user_update_logic.go`、`user_list_logic.go`）继续直调 `svcCtx.Domain.IAM.User`。

### 简单域模板：Blog / Article

`CreateWithTags` 满足"跨表写"标准，需要一个薄的领域服务包一层事务；`Article` 域其余方法（`FindByID`/`FindPage`/`Delete`）不需要：

**文件名、类型名、方法名以 `06-domain-blog-video-sdk.md` §2.2 的完整实现为准**（本文档只示意"简单域该长什么样"，不重复定义一份不一致的版本）：

```go
// internal/domain/content/blog_service.go（目录按下方"包边界"一节归入 content 分组）
package content

import (
	"context"

	"postapocgame/admin-server/internal/repository"
	blogrepo "postapocgame/admin-server/internal/repository/blog"
	blogmodel "postapocgame/admin-server/internal/model/blog"
)

// BlogArticleService 只承载 CreateArticle/UpdateArticle（跨 blog_article+blog_article_tag 两表写）
// 和 AuditArticle（跨 blog_article_audit+blog_article 两表写）这类方法，完整实现见 06 文档。
type BlogArticleService struct {
	repo *repository.Repository
}

func NewBlogArticleService(repo *repository.Repository) *BlogArticleService {
	return &BlogArticleService{repo: repo}
}

func (s *BlogArticleService) CreateArticle(ctx context.Context, article *blogmodel.BlogArticle, tagIDs []uint64) error {
	// 领域服务内部只用 Repository.Transact，不 import internal/repository/registry。
	return s.repo.Transact(ctx, func(ctx context.Context, txRepo *repository.Repository) error {
		return blogrepo.NewBlogArticleRepository(txRepo).CreateWithTags(ctx, article, tagIDs)
	})
}
```

`registry.Domain` 里这个服务挂在哪个字段、`internal/domain/content` 包名和 `registry.Domain.Blog`/`registry.Domain.Video` 字段路径为什么不需要保持一致，见 06 文档第 5 节，本文档不重复。

`blog_article_create_logic.go` 改调 `svcCtx.Domain.Blog.ArticleService.CreateArticle(...)`（领域服务方法名是 `CreateArticle`，内部才调用 repository 层的 `CreateWithTags`，两者不是同一个名字，照抄时不要混用）；`blog_article_detail_logic.go`/`blog_article_list_logic.go` 等纯读方法保持直调 `svcCtx.Domain.Blog.Article.FindByID(...)` 不变。

### 关键新增要求：包边界按 Phase 2 最终的 5 服务分组组织

这是本节最重要、其他文档会依赖的决策，**不能软化**。`internal/domain/<domain>` 的目录从一开始就要按 Phase 2（`15-service-boundaries.md`）最终的 5 个 RPC 服务边界组织，不按现在的 9 个域目录一一对应：

- **不单独建 `internal/domain/monitoring`**——监控相关领域服务代码放进 `internal/domain/iam`（或紧邻的子包，如 `internal/domain/iam/monitoring/`），因为最终它们会合并进同一个 `iam-rpc` 服务。
- **`blog` 和 `video` 的领域服务代码放在同一个分组下**（即使现在数据表还是分开的两套），对应未来的 `content-rpc`——例如 `internal/domain/content/blog_service.go`、`internal/domain/content/video_service.go` 同属 `package content` 或至少同一父目录。
- **`system`、`misc` 同样并入 `iam` 分组**。
- **`chat`、`task`、`sdk` 各自独立分组**，直接对应未来的 `chat-rpc`、`task-rpc`、`sdk-rpc`，维持现状（`internal/domain/chat`、`internal/domain/task`、`internal/domain/sdk`）即可。

目标目录形状：

```
internal/domain/
├── iam/            # 承载 iam + system + monitoring + misc 的领域服务
│   ├── user_service.go
│   ├── permission_resolver.go
│   ├── system/     # 或直接平铺在 iam/ 下，命名不冲突即可，Phase 1 阶段不强求子包
│   └── monitoring/
├── content/        # 承载 blog + video 的领域服务
│   └── blog_service.go   # 只有 Blog 有真正的领域服务，Video 没有，见 06 文档第 3 节
├── chat/
├── task/
└── sdk/
```

这样 Phase 2 把每个服务拆出去时，接近"挪目录 + 加一个 `main.go` + 包一层 `.proto`"，而不是重新设计包结构。`04`~`07` 四篇分域文档在建领域服务包时必须遵守这个分组，不要按方便就地建 `internal/domain/monitoring`、`internal/domain/system` 这种未来要拆掉的目录。

---

## A.4 中间件接入 Wire

### 现状

11 个中间件构造函数（`internal/middleware/*.go`）全部签名统一为 `func NewXxxMiddleware(svcCtx *svc.ServiceContext) *XxxMiddleware`，但实测每个中间件真正用到的字段远比整个 `ServiceContext` 窄：

| 中间件文件 | 实际用到的字段 |
|---|---|
| `authmiddleware.go` | `svcCtx.Repository`（`iamrepo.NewTokenBlacklistRepository`）+ `svcCtx.Config.JWT.AccessSecret` |
| `permissionmiddleware.go` | `svcCtx.Repository`（`iamdomain.NewPermissionResolver`） |
| `apienabledmiddleware.go` | `svcCtx.Repository`（`iamrepo.NewApiRepository`） |
| `operationlogmiddleware.go` | `svcCtx.Repository`（`monitoringrepo.NewOperationLogRepository`） |
| `publicoperationlogmiddleware.go` | `svcCtx.Repository`（同上） |
| `performancemiddleware.go` | `svcCtx.Repository`（`monitoringrepo.NewPerformanceLogRepository`）+ `svcCtx.Config.Database.SlowQueryThreshold` |
| `ratelimitmiddleware.go` | `svcCtx.Repository.Redis` + `svcCtx.Config.RateLimit` |
| `sdkauthmiddleware.go` | `svcCtx.Repository`（`sdkrepo.NewSdkRepository`） |
| `sdkcalllogmiddleware.go` | `svcCtx.Repository`（同上） |
| `sdkratelimitmiddleware.go` | `svcCtx.Repository`（同上，含 `.Redis`） |
| `corsmiddleware.go` | 无（`NewCorsMiddleware()` 已经是零参数） |

没有任何一个中间件用到 `svcCtx.Domain`、`svcCtx.ChatHub`、`svcCtx.TaskScheduler`，或任何一个 `xxxMiddleware` 字段本身（不存在中间件依赖另一个中间件）。之所以现在都吃整个 `*svc.ServiceContext`，是因为构造时机上 `ServiceContext` 本身还没构造完（中间件字段本身就是 `ServiceContext` 的一部分），`internal/wire/providers.go` 的 `buildMiddlewareBundle(svcCtx)` 先拿到一个"Repository/Domain 已经赋值、中间件字段还是空"的 `svcCtx`，再手工 `New` 出全部中间件——这是能跑但不是真正 Wire 装配的折中写法，Wire 无法把这一步纳入编译期依赖图。

### 修复：narrowed 构造函数 + Wire provider

把每个中间件构造函数的参数收窄到只声明它真正依赖的类型（`config.Config` 的相关子结构，或直接吃整个 `config.Config`；`*repository.Repository`）。这三类依赖（`config.Config`、`*repository.Repository`、`*registry.Domain`）已经是独立的 Wire 节点（分别由 Wire 的隐式输入参数、`provideRepository`、`provideDomain` 提供），narrow 之后 Wire 才能在不依赖 `svcCtx` 自身的情况下把每个中间件构造出来，从而打破"中间件依赖 `svcCtx`、`svcCtx` 又包含中间件字段"的循环。

先在 `AuthMiddleware`/`PermissionMiddleware` 上跑通这个模式（这两个中间件后续也是 Phase 2 里 gateway 调用 `iam-rpc.CheckPermission` 的直接前身，提前收窄有额外价值）：

```go
// internal/middleware/authmiddleware.go
type AuthMiddleware struct {
	repo      *repository.Repository
	jwtConfig config.JWTConf
}

func NewAuthMiddleware(cfg config.Config, repo *repository.Repository) *AuthMiddleware {
	return &AuthMiddleware{repo: repo, jwtConfig: cfg.JWT}
}
```

```go
// internal/middleware/permissionmiddleware.go
type PermissionMiddleware struct {
	repo *repository.Repository
}

func NewPermissionMiddleware(repo *repository.Repository) *PermissionMiddleware {
	return &PermissionMiddleware{repo: repo}
}
```

其余 9 个中间件按同样方式收窄，逐个签名见 `03-wire-and-middleware.md`。`svcCtx.AuthMiddleware` 等字段本身保持扁平结构不变（不引入 `svcCtx.Middleware.Auth` 嵌套，上一轮重构已经明确排除这个改法）——变的只是"谁来构造这些中间件实例"，不是 `ServiceContext` 的对外形状。Wire provider 的具体注册方式（`wire.NewSet` 里加哪些条目、`MiddlewareBundle` 的组装方式）是可执行细节，见 `03-wire-and-middleware.md`，本文档只定"收窄依赖"这一条设计原则。

---

## A.5 密钥管理

go-zero 1.9.3 的 `conf.MustLoad` 原生支持 `conf.UseEnv()`（对 YAML 内容做 `os.ExpandEnv()`），无需新依赖、无需手写 `${VAR}` 解析逻辑。当前 `admin.go` 第 55 行：

```go
var c config.Config
conf.MustLoad(*configFile, &c)
```

MySQL/Redis 配置已经通过 `config.MergeExternalConfig(&c, *mysqlConfigFile, *redisConfigFile, *middlewareConfigFile)`（`admin.go` 第 58 行）走外部文件、不提交仓库，这部分没问题。问题只在 `etc/admin-api.yaml` 里硬编码的 JWT 占位符：

```yaml
JWT:
  AccessSecret: "replace-with-secure-access-secret"
  RefreshSecret: "replace-with-secure-refresh-secret"
```

### 修复

1. `etc/admin-api.yaml`：

```yaml
JWT:
  AccessSecret: "${JWT_ACCESS_SECRET}"
  RefreshSecret: "${JWT_REFRESH_SECRET}"
```

2. `admin.go`：

```go
var c config.Config
conf.MustLoad(*configFile, &c, conf.UseEnv())

if c.JWT.AccessSecret == "" || c.JWT.RefreshSecret == "" {
	log.Fatalf("JWT_ACCESS_SECRET / JWT_REFRESH_SECRET 未设置，拒绝以空密钥启动")
}
```

fail-fast 检查放在 `conf.MustLoad` 之后、`appwire.InitializeApp(c)` 之前，避免用空密钥跑完整个启动流程才在运行时报错。

3. 本地开发用的 JWT 密钥可以由 AI 自动生成并写进本地 `.env`/环境变量说明，注明"仅供本地开发"；生产环境真实密钥值必须由用户设置，不能由 AI 代为生成——这条属于 `10-dev-execution-and-review-points.md` 里"必须真正停下来问用户"的第 3 类情形，本文档只定技术改法。这一条要素也要记入 `14-production-deployment-checklist.md`。
