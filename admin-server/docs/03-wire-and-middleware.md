# 中间件收窄 + Wire provider 化

## 前置依赖

- 先读 `01-architecture-target.md` 的 A.4 节——本文档是它的可执行拆解。
- 代码库状态：`internal/wire/{providers.go,wire.go,wire_gen.go,middleware_bundle.go}`、`internal/middleware/*.go` 处于当前 HEAD（`36fbda9`）未改动状态。
- 不依赖 `02-transactions-and-uow.md` 的改动，两者可以并行/任意顺序做。

---

## 1. 现状：11 个中间件的真实依赖面

全部 11 个中间件构造函数当前签名统一是 `func NewXxxMiddleware(svcCtx *svc.ServiceContext) *XxxMiddleware`，结构体统一是 `struct { svcCtx *svc.ServiceContext }`。逐文件核查（`grep -n "svcCtx\."` 全部 11 个文件）实际用到的字段：

| 文件 | 当前构造模式 | 实际依赖 |
|---|---|---|
| `authmiddleware.go` | 存 `svcCtx`，`Handle` 内现场 `iamrepo.NewTokenBlacklistRepository(m.svcCtx.Repository)` + 读 `m.svcCtx.Config.JWT.AccessSecret` | `*repository.Repository` + `config.JWTConf` |
| `permissionmiddleware.go` | 存 `svcCtx`，`Handle` 内现场 `iamdomain.NewPermissionResolver(m.svcCtx.Repository)`——**每次请求都重新构造一次 `PermissionResolver`，这本身就是要修的 bug**，见下方脚注 | 目标依赖是 `*iamdomain.PermissionResolver`（不是 `*repository.Repository`），见下方脚注 |
| `apienabledmiddleware.go` | 存 `svcCtx`，`Handle` 内现场 `iamrepo.NewApiRepository(m.svcCtx.Repository)` | `*repository.Repository` |
| `operationlogmiddleware.go` | 存 `svcCtx`，构造时启动 `logWriter` goroutine，`logWriter`/`Handle` 内用 `monitoringrepo.NewOperationLogRepository(m.svcCtx.Repository)` | `*repository.Repository` |
| `publicoperationlogmiddleware.go` | 同上模式 | `*repository.Repository` |
| `performancemiddleware.go` | 构造时读 `svcCtx.Config.Database.SlowQueryThreshold` 算慢阈值，`Handle` 内用 `monitoringrepo.NewPerformanceLogRepository(m.svcCtx.Repository)` | `*repository.Repository` + `config.DatabaseConf`（只需 `SlowQueryThreshold` 字段） |
| `ratelimitmiddleware.go` | 构造时用 `svcCtx.Repository.Redis` 建限流器，`Handle` 内读 `m.svcCtx.Config.RateLimit` | `*repository.Repository`（`.Redis`） + `config.RateLimitConf` |
| `sdkauthmiddleware.go` | 存 `svcCtx`，`Handle` 内现场 `sdkrepo.NewSdkRepository(m.svcCtx.Repository)` | `*repository.Repository` |
| `sdkcalllogmiddleware.go` | 同上模式 | `*repository.Repository` |
| `sdkratelimitmiddleware.go` | 同上模式，另加 `m.svcCtx.Repository.Redis` | `*repository.Repository` |
| `corsmiddleware.go` | `NewCorsMiddleware()` 零参数，结构体 `struct{}` | 无 |

**没有任何一个中间件的*现状代码*用到 `svcCtx.Domain`、`svcCtx.ChatHub`、`svcCtx.TaskExecutors`、`svcCtx.TaskScheduler`，或彼此的 `xxxMiddleware` 字段。** 现状里 10 个非 Cors 中间件的依赖面只有两类：`*repository.Repository`（全部都要）、`config.Config` 的某个子结构（只有 `AuthMiddleware`/`PerformanceMiddleware`/`RateLimitMiddleware` 三个要）。

**唯一例外，且是本轮改造之后才成立的目标状态，不是现状**：`PermissionMiddleware`。`04-domain-iam-chat.md` 任务 3 把 `PermissionResolver` 的构造从"每次请求现场 new"挪进了 `registry.NewDomain`（这正是要修的 bug——现状每次请求都重新构造一次，见上表脚注），挪完之后 `PermissionMiddleware` 该依赖的就不再是 `*repository.Repository`（不然又要在 `Handle` 里重新 `iamdomain.NewPermissionResolver(...)` 一次，等于没修），而是已经构造好的 `*iamdomain.PermissionResolver`。**处理方式是"Wire provider 函数传 `*registry.Domain`，中间件构造函数只吃取出来的 `*iamdomain.PermissionResolver`"，不是让中间件结构体持有整个 `*registry.Domain`**——那样会让一个只做权限判定的中间件拿到读写任意域的能力，超出实际需要。具体代码见 04 文档任务 4，下方 2.2 节的 `PermissionMiddleware` 示例已经按这个目标状态写，不是按上表"现状"列写的。除这一个中间件外，`01-architecture-target.md` A.4 节"不需要给中间件暴露 `*registry.Domain`"的结论对其余 10 个中间件仍然成立。

`internal/wire/providers.go` 现状：

```go
func provideServiceContext(...) (*svc.ServiceContext, func()) {
	svcCtx := &svc.ServiceContext{ /* Config/Repository/Domain/... 先赋值 */ }
	mw := buildMiddlewareBundle(svcCtx)   // 拿到已经部分赋值的 svcCtx，手工现场 New 出全部中间件
	svcCtx.AuthMiddleware = mw.Auth
	// ...
}

func buildMiddlewareBundle(svcCtx *svc.ServiceContext) *MiddlewareBundle {
	return &MiddlewareBundle{
		Auth: middleware.NewAuthMiddleware(svcCtx).Handle,
		// ... 11 个
	}
}
```

`buildMiddlewareBundle` 不是 Wire provider（没有注册进 `wire.NewSet`），是普通函数调用，Wire 编译期依赖图里看不到这一步——这就是"不是真正 Wire 装配"的具体含义。

---

## 2. 目标：narrowed 构造函数

### 2.1 `AuthMiddleware`（先跑通模式）

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

`Handle` 内部把 `m.svcCtx.Repository` 改成 `m.repo`，`m.svcCtx.Config.JWT.AccessSecret` 改成 `m.jwtConfig.AccessSecret`。

### 2.2 `PermissionMiddleware`（第二个跑通模式，也是唯一不吃 `*repository.Repository` 的中间件）

**不要**把这个中间件收窄成"存 `*repository.Repository`、`Handle` 内改成 `iamdomain.NewPermissionResolver(m.repo)`"——那只是把"每次请求现场 new 一次 `PermissionResolver`"的构造代码从吃 `svcCtx` 改成吃 `repo`，没有解决"现场 new"这个真正的问题。正确做法是配合 04 文档任务 3（`PermissionResolver` 挪进 `registry.NewDomain`，全进程只构造一次）：中间件直接持有构造好的 `*iamdomain.PermissionResolver`，请求处理时不再有任何构造动作：

```go
// internal/middleware/permissionmiddleware.go
type PermissionMiddleware struct {
	resolver *iamdomain.PermissionResolver
}

func NewPermissionMiddleware(resolver *iamdomain.PermissionResolver) *PermissionMiddleware {
	return &PermissionMiddleware{resolver: resolver}
}
```

`Handle` 内 `iamdomain.NewPermissionResolver(m.svcCtx.Repository).CanAccess(...)`（每次请求都构造）改成 `m.resolver.CanAccess(...)`（直接用已构造好的实例）。对应的 Wire provider 函数入参是 `*registry.Domain`（`providePermissionMiddleware(domain *registry.Domain) *middleware.PermissionMiddleware { return middleware.NewPermissionMiddleware(domain.IAM.PermissionResolver) }`），但中间件构造函数本身的签名只声明 `*iamdomain.PermissionResolver`，不声明 `*registry.Domain`——provider 函数负责"从 Domain 里取出中间件真正需要的那一小块"，中间件不需要知道 `registry.Domain` 整体的存在。这是全部 11 个中间件里唯一一个需要额外包一层 `providePermissionMiddleware` 适配函数（而不是把 `NewPermissionMiddleware` 直接注册进 `wire.NewSet`）的，3.4 节的 `ProviderSet` 会体现这个差异。

### 2.3 其余 8 个吃 `*repository.Repository` 的中间件

```go
// internal/middleware/apienabledmiddleware.go
type ApiEnabledMiddleware struct{ repo *repository.Repository }
func NewApiEnabledMiddleware(repo *repository.Repository) *ApiEnabledMiddleware {
	return &ApiEnabledMiddleware{repo: repo}
}

// internal/middleware/operationlogmiddleware.go
type OperationLogMiddleware struct {
	repo  *repository.Repository
	logCh chan *monitoring.AdminOperationLog
}
func NewOperationLogMiddleware(repo *repository.Repository) *OperationLogMiddleware {
	m := &OperationLogMiddleware{repo: repo, logCh: make(chan *monitoring.AdminOperationLog, 1000)}
	go m.logWriter()
	return m
}

// internal/middleware/publicoperationlogmiddleware.go —— 同 OperationLogMiddleware 模式

// internal/middleware/sdkauthmiddleware.go
type SDKAuthMiddleware struct{ repo *repository.Repository }
func NewSDKAuthMiddleware(repo *repository.Repository) *SDKAuthMiddleware {
	return &SDKAuthMiddleware{repo: repo}
}

// internal/middleware/sdkcalllogmiddleware.go —— 同上模式
// internal/middleware/sdkratelimitmiddleware.go —— 同上模式
```

### 2.4 `PerformanceMiddleware`（需要 Config 子结构）

```go
type PerformanceMiddleware struct {
	repo          *repository.Repository
	monitor       *monitor.PerformanceMonitor
	slowThreshold int64
}

func NewPerformanceMiddleware(cfg config.Config, repo *repository.Repository) *PerformanceMiddleware {
	slowThreshold := int64(2000)
	if cfg.Database.SlowQueryThreshold > 0 {
		slowThreshold = int64(cfg.Database.SlowQueryThreshold) * 2
	}
	return &PerformanceMiddleware{
		repo:          repo,
		monitor:       monitor.NewPerformanceMonitor(slowThreshold),
		slowThreshold: slowThreshold,
	}
}
```

### 2.5 `RateLimitMiddleware`（需要 Config 子结构 + Redis）

```go
type RateLimitMiddleware struct {
	repo   *repository.Repository
	config config.RateLimitConf
	// ipLimiters/userLimiters/apiLimiters/globalLimiter/mu 等字段不变
}

func NewRateLimitMiddleware(cfg config.Config, repo *repository.Repository) *RateLimitMiddleware {
	// 原本用 svcCtx.Repository.Redis 构造限流器的逻辑改用 repo.Redis；
	// 原本 Handle 内 m.svcCtx.Config.RateLimit 改成 m.config
}
```

### 2.6 `CorsMiddleware`

不用改，`NewCorsMiddleware()` 已经零参数。

---

## 3. `internal/wire/providers.go` 改动

### 3.1 10 个 narrowed 构造函数本身就是 Wire provider，`PermissionMiddleware` 例外

10 个中间件（除 `PermissionMiddleware`）参数只剩 `config.Config`（Wire 隐式输入，`InitializeApp(c config.Config)` 的入参）和 `*repository.Repository`（`provideRepository` 提供）——这两者都已经是 Wire 图里的节点，所以 Wire 可以直接把 `middleware.NewXxxMiddleware` 函数本身注册进 `wire.NewSet`，不需要额外包一层适配函数。

**`PermissionMiddleware` 不满足这个条件**：它的构造函数吃 `*iamdomain.PermissionResolver`，但 `*iamdomain.PermissionResolver` 本身不是一个独立的 Wire 节点——它只是 `*registry.Domain` 结构体里的一个字段（`domain.IAM.PermissionResolver`），Wire 没法凭空"生产"出一个游离的 `*iamdomain.PermissionResolver`。所以 `PermissionMiddleware` 需要 2.2 节提到的 `providePermissionMiddleware(domain *registry.Domain) *middleware.PermissionMiddleware` 适配函数：Wire 用已有的 `*registry.Domain` 节点调用这个适配函数，适配函数内部取出字段、调用 `middleware.NewPermissionMiddleware`。`ProviderSet` 里注册的是 `providePermissionMiddleware`，不是 `middleware.NewPermissionMiddleware` 本身（3.4 节体现这一点）。

### 3.2 新增 `provideMiddlewareBundle`：唯一还需要手写的 assembler

`MiddlewareBundle`（`internal/wire/middleware_bundle.go`）里的字段全部是 `rest.Middleware`（即 `func(http.HandlerFunc) http.HandlerFunc` 的具名类型），11 个字段类型相同——这是 Wire 做不到自动区分的地方（同一类型的多个值，Wire 不知道该把哪个 provider 的结果填进哪个字段）。所以仍然需要一个手写的 assembler 函数，区别在于它现在只依赖 11 个**互不相同的具体中间件指针类型**，而不是重新依赖 `svcCtx`：

```go
func provideMiddlewareBundle(
	auth *middleware.AuthMiddleware,
	apiEnabled *middleware.ApiEnabledMiddleware,
	permission *middleware.PermissionMiddleware,
	operationLog *middleware.OperationLogMiddleware,
	publicOperationLog *middleware.PublicOperationLogMiddleware,
	rateLimit *middleware.RateLimitMiddleware,
	performance *middleware.PerformanceMiddleware,
	cors *middleware.CorsMiddleware,
	sdkAuth *middleware.SDKAuthMiddleware,
	sdkRateLimit *middleware.SDKRateLimitMiddleware,
	sdkCallLog *middleware.SDKCallLogMiddleware,
) *MiddlewareBundle {
	return &MiddlewareBundle{
		Auth:               auth.Handle,
		ApiEnabled:         apiEnabled.Handle,
		Permission:         permission.Handle,
		OperationLog:       operationLog.Handle,
		PublicOperationLog: publicOperationLog.Handle,
		RateLimit:          rateLimit.Handle,
		Performance:        performance.Handle,
		Cors:               cors.Handle,
		SDKAuth:            sdkAuth.Handle,
		SDKRateLimit:       sdkRateLimit.Handle,
		SDKCallLog:         sdkCallLog.Handle,
	}
}
```

这一步和现状的关键区别：`buildMiddlewareBundle(svcCtx)` 是在 `provideServiceContext` 内部手工调用的普通函数；`provideMiddlewareBundle` 是注册进 `wire.NewSet` 的正式 provider，它的 11 个参数由 Wire 自动从各自的 `NewXxxMiddleware` provider 解析、构造、注入，`provideServiceContext` 不再需要知道任何一个中间件是怎么构造出来的，只需要声明自己依赖一个 `*MiddlewareBundle`。

### 3.3 `provideServiceContext` 简化

```go
func provideServiceContext(
	c config.Config,
	repo *repository.Repository,
	domain *registry.Domain,
	chatHub *hub.ChatHub,
	taskExecutors map[int]interfaces.TaskExecutor,
	taskScheduler *task.TaskScheduler,
	mw *MiddlewareBundle,
) (*svc.ServiceContext, func()) {
	svcCtx := &svc.ServiceContext{
		Config:                       c,
		Repository:                   repo,
		Domain:                       domain,
		ChatHub:                      chatHub,
		TaskExecutors:                taskExecutors,
		TaskScheduler:                taskScheduler,
		AuthMiddleware:               mw.Auth,
		ApiEnabledMiddleware:         mw.ApiEnabled,
		PermissionMiddleware:         mw.Permission,
		OperationLogMiddleware:       mw.OperationLog,
		PublicOperationLogMiddleware: mw.PublicOperationLog,
		RateLimitMiddleware:          mw.RateLimit,
		PerformanceMiddleware:        mw.Performance,
		CorsMiddleware:               mw.Cors,
		SDKAuthMiddleware:            mw.SDKAuth,
		SDKRateLimitMiddleware:       mw.SDKRateLimit,
		SDKCallLogMiddleware:         mw.SDKCallLog,
	}

	cleanup := func() {
		if taskScheduler != nil {
			taskScheduler.Stop()
			logx.Infof("任务调度器已停止")
		}
	}
	return svcCtx, cleanup
}
```

`buildMiddlewareBundle` 函数整体删除——不再需要，它存在的唯一理由（中间件依赖完整 `svcCtx`）已经被 2 节的收窄改动消除。

### 3.4 `ProviderSet` 最终形态

```go
var ProviderSet = wire.NewSet(
	provideRepository,
	provideDomain,
	provideChatHub,
	provideTaskExecutors,
	provideTaskScheduler,

	middleware.NewAuthMiddleware,
	middleware.NewApiEnabledMiddleware,
	providePermissionMiddleware, // 适配函数，不是 middleware.NewPermissionMiddleware 本身，见 2.2/3.1
	middleware.NewOperationLogMiddleware,
	middleware.NewPublicOperationLogMiddleware,
	middleware.NewRateLimitMiddleware,
	middleware.NewPerformanceMiddleware,
	middleware.NewCorsMiddleware,
	middleware.NewSDKAuthMiddleware,
	middleware.NewSDKRateLimitMiddleware,
	middleware.NewSDKCallLogMiddleware,
	provideMiddlewareBundle,

	provideServiceContext,
)
```

`internal/wire/wire.go`（`wireinject` build tag 文件）和 `InitializeApp` 签名都不需要改，`wire.Build(ProviderSet)` 自动感知新增的 provider。

---

## 4. `make wire` 重新生成

改完 `providers.go` 后跑：

```
cd admin-server && make wire
```

`Makefile` 里 `wire` target 的实际命令是 `go run github.com/google/wire/cmd/wire ./internal/wire`，重新生成 `internal/wire/wire_gen.go`，再跟一次 `go mod tidy`。这一步属于 `10-dev-execution-and-review-points.md` 里"开发期可以直接执行、事后 review、不阻塞"的第 3 类（`make wire`/`goctl rpc` 生成产物的提交，diff 留给用户日常 review）——AI 可以直接跑 `make wire` 并提交生成结果，不需要为这一步单独停下来问用户；如果 Wire 编译期依赖图解析失败（例如遗漏某个 provider、出现循环依赖），报错信息会直接指出缺哪个 provider，按提示修正即可。

## 完成的定义

1. 全部 11 个 `internal/middleware/*.go` 构造函数按第 2 节收窄完成，`Handle` 内部引用从 `m.svcCtx.X` 改成 `m.repo`/`m.jwtConfig`/`m.config` 等窄字段，`svcCtx` 字段/import 从结构体里移除（`CorsMiddleware` 不变）。
2. `internal/wire/providers.go` 按第 3 节改完：`buildMiddlewareBundle` 删除，新增 `provideMiddlewareBundle`，`ProviderSet` 注册全部 11 个中间件构造函数 + `provideMiddlewareBundle`。
3. `internal/wire/middleware_bundle.go` 不需要改（`MiddlewareBundle` 结构体形状不变）。
4. `svc.ServiceContext`（`internal/svc/servicecontext.go`）结构体形状不变——11 个 `rest.Middleware` 字段保持扁平，不引入嵌套。
5. `make wire` 重新生成 `internal/wire/wire_gen.go` 成功，无需手工修补生成结果。
6. `go build ./...` 全仓库编译通过。
7. 人工冒烟：启动服务，验证登录（`AuthMiddleware`）、一个受权限保护的接口（`PermissionMiddleware`）、一个限流触发场景（`RateLimitMiddleware`）、一个 SDK 调用（`SDKAuthMiddleware`）行为与改动前一致。
