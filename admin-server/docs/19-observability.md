# 19. 可观测性：Telemetry + 结构化日志

> 本文档是可直接执行的任务说明，对应总纲 Part C.1（Phase 3，约 Week 13-14+）。执行者在改动前应完整阅读一遍。改完要跑通「完成的定义」里列出的验证步骤，不要跨文档并行改动（先做完本篇，再做 `20-api-docs-generation.md`/`21-cd-and-deployment.md`）。

## 0. 前置依赖

- Part B（`15-service-boundaries.md` ~ `18-service-extraction-runbook.md`）已落地：`gateway` + `iam-rpc`/`content-rpc`/`chat-rpc`/`task-rpc`/`sdk-rpc` 五个 RPC 服务 + gateway 共六个独立部署单元已经存在，各自有自己的 `etc/*.yaml`。
- 本文档**不需要**等六个服务全部拆完才能开始——Telemetry 配置块和 `pkg/logging` 包可以在最后一个服务（`iam-rpc`）拆分完成后统一接入一次，也可以在每个服务拆分时顺手接入；两种顺序都不影响正确性，按 Phase 2 实际进度选择。

## 1. 为什么这一步不能省（这不是可选项）

单体时代排查一个线上问题，是"看一个日志流、跟一条调用链"——`admin.go` 一个进程、一份 stdout，出错时 `grep` 一下请求时间点附近的日志基本就能看到从 handler 到 repository 的完整轨迹，人肉对照时间戳就能拼出因果关系。

微服务拆分后，这个前提彻底不成立了。一个用户可见的错误现在可能横跨 `gateway → iam-rpc → chat-rpc → task-rpc` 四个独立进程、四份完全独立的日志目的地（不同的 stdout、不同的 Supervisor/容器日志文件）。没有一个跨服务共享的 trace id 把这四份日志串起来，排查方式退化成：在四个地方分别搜索"大概那个时间点、大概那个用户 ID"的日志行，再靠人工核对时间戳、用户 ID、请求参数去猜哪几行属于同一次请求。

这个退化不是线性的——**不是"慢一点"，是过了 1-2 跳调用就基本查不动**。两个服务之间的调用还能靠时间窗口勉强对上；三跳以上，加上并发请求导致的时间戳交错、加上服务重启/多副本导致的日志文件轮转，人工关联的成本呈指数增长，实际后果是"这个 bug 没法排查，只能改完重现",这在生产环境是不可接受的。

这正是拆分本身制造出来的问题——单体时没有这个问题，拆分之后才有。所以它必须在**第一次真实的多服务故障之前**就位，不能等出了事故再补，那时候补是"debug 一个无法 debug 的系统"的悖论。

## 2. Telemetry 接入：go-zero 原生能力，零新增依赖

先核实过依赖现状：`admin-server/go.mod` 里已经通过 go-zero 1.9.3 间接带了完整的 `go.opentelemetry.io/otel*` 系列（`otel`、`otel/sdk`、`otel/trace`、`otel/metric`、`otel/exporters/otlp/otlptrace{,grpc,http}`、`otel/exporters/jaeger`、`otel/exporters/zipkin`、`otel/exporters/stdout/stdouttrace`），当前全部标记为 `// indirect`——go-zero 内部用，业务代码目前完全没有引用（`grep -rn "Telemetry\|opentelemetry\|otel" admin-server --include="*.go" --include="*.yaml"` 零命中）。这意味着接入 Telemetry **不需要 `go get` 任何新包**，只是把 go-zero 已经实现好的能力通过配置打开。

go-zero 的 `rest.RestConf`/`zrpc.RpcServerConf` 都内嵌了 `service.ServiceConf`，其中包含 `Telemetry trace.Config`，`conf.MustLoad` 解析 YAML 时会自动识别这个字段，框架层在 `rest.MustNewServer`/`zrpc.MustNewServer` 内部自动完成 span 创建、HTTP `traceparent` 头和 gRPC metadata 的传播拼接——业务代码不需要手写任何 span 相关代码就能获得跨进程 trace 传播。

### 2.1 每个服务 `etc/*.yaml` 加 `Telemetry` 块

以 `etc/iam.yaml`（zrpc 服务）为例：

```yaml
Name: iam.rpc
ListenOn: 0.0.0.0:8081
Etcd:
  Hosts: []          # 静态 target，未接入服务发现（见 16-rpc-conventions.md）

Telemetry:
  Name: iam-rpc
  Endpoint: http://127.0.0.1:4318/v1/traces   # 本地/自建 OTLP collector，未部署时留空即可关闭
  Sampler: 1.0        # 开发环境全采样；生产按流量调低（如 0.1）
  Batcher: otlpgrpc    # 可选 otlpgrpc / otlphttp / jaeger / zipkin / stdout
```

`gateway` 的 `etc/gateway.yaml`（HTTP 服务）同样加：

```yaml
Name: gateway-api
Host: 0.0.0.0
Port: 20000

Telemetry:
  Name: gateway
  Endpoint: http://127.0.0.1:4318/v1/traces
  Sampler: 1.0
  Batcher: otlpgrpc
```

六个服务（`gateway`/`iam-rpc`/`content-rpc`/`chat-rpc`/`task-rpc`/`sdk-rpc`）各自的 `Telemetry.Name` 必须不同（用于在 trace 后端区分服务来源），`Endpoint` 指向同一个 collector。

**本轮不部署 OTLP collector/Jaeger/Zipkin 后端**——`Telemetry` 配置块本身会自动生成并传播 trace-id，这个能力在没有后端收集时依然生效（span 数据会被静默丢弃，但 trace-id 已经写进日志字段，见第 3 节），本轮的重点是把 trace-id 打进结构化日志，让"跨服务用同一个 trace-id 搜四份日志"这件事可行；接一个真正的 trace 后端做可视化调用链，是后续可选的基础设施投入，不阻塞本轮。

### 2.2 gRPC/HTTP 传播的落地位置

不需要手写传播代码——`rest.MustNewServer`（gateway）和 `zrpc.MustNewServer`（五个 RPC 服务）在检测到 `Telemetry` 配置非空时，会自动挂载 otel 的 HTTP/gRPC 拦截器/中间件，完成：
- HTTP 请求头 `traceparent` 的读取（如果上游已经带了）或生成（如果是新请求的入口，即 gateway）；
- gateway 调用各 RPC 服务时，trace context 自动注入 gRPC metadata；
- 各 RPC 服务收到请求时自动从 metadata 提取 trace context，延续同一条 trace。

`Sampler` 建议开发环境设 `1.0`（全量采样，方便本地调试全链路），生产环境按实际流量调低（如 `0.1`），避免 trace 数据量过大；这是纯配置改动，不需要代码配合。

## 3. `pkg/logging`：统一结构化日志

### 3.1 现状核查

`admin-server/admin.go` 第 62-67 行：

```go
err := logx.SetUp(logx.LogConf{
    Encoding: "plain",
})
```

确认现状是 `plain` 编码（纯文本行日志），不是 JSON。六个服务（`cmd/gateway/main.go` 继任 `admin.go`，加上五个 `services/<name>/<name>.go` 的 RPC `main`）目前各自有一份类似的启动代码，日志格式没有统一约定，也没有 `trace_id`/`service` 等字段。

### 3.2 新增 `pkg/logging` 包

```
pkg/logging/
├── setup.go      # 统一的 logx.SetUp 封装，替换六处分散的 logx.SetUp 调用
└── fields.go      # 标准字段常量 + 从 ctx 提取字段的 helper
```

`pkg/logging/setup.go`：

```go
package logging

import (
	"github.com/zeromicro/go-zero/core/logx"
)

// Setup 统一的日志初始化，六个服务的 main 函数都调用这一个函数，
// 不再各自内联 logx.SetUp(logx.LogConf{...})。
// serviceName 写入每条日志的 service 字段，用于在聚合后按服务过滤。
func Setup(serviceName string) error {
	return logx.SetUp(logx.LogConf{
		ServiceName: serviceName,
		Encoding:    "json", // 从 plain 改为 json，可被日志采集/检索工具结构化解析
		Level:       "info",
	})
}
```

`pkg/logging/fields.go`：

```go
package logging

// 标准字段集：所有服务的结构化日志最少要包含这几个字段，
// 便于跨服务用同一个 trace_id 检索、按 service/user_id 过滤。
const (
	FieldTraceID = "trace_id"
	FieldSpanID  = "span_id"
	FieldService = "service"
	FieldUserID  = "user_id"
)
```

`trace_id`/`span_id` 由 go-zero 的 `logx.WithContext(ctx)` **自动注入**——只要 `ctx` 里已经有 Telemetry 中间件放进去的 span context（第 2 节的自动传播），`logx.WithContext(ctx)` 返回的 logger 打印每一条日志时就会自带 `trace` 字段（go-zero 内部字段名，含 trace-id），不需要业务代码手动从 ctx 取 trace-id 再拼进日志——这是选用 go-zero 原生 Telemetry + logx 组合而不是自己接 OTel SDK 手写日志字段的关键原因，两者是同一套框架、天然打通。

`service` 字段由 `Setup(serviceName)` 的 `ServiceName` 参数在初始化时全局设置一次，每条日志自动带上。`user_id` 目前日志里没有统一字段，需要在能拿到用户身份的地方（鉴权中间件之后的 handler/logic）显式补：约定通过 `logx.WithContext(ctx).WithFields(logx.Field("user_id", uid))` 附加，而不是塞进 ctx 让框架自动带出——`user_id` 是业务字段，不是 go-zero 框架层认识的东西。

### 3.3 现有代码里 `logx.WithContext(ctx)` 的实际用法（用于对照，不用改）

现状核查（`grep -rn "logx.WithContext" internal`，命中约 130+ 处），典型模式来自 goctl 生成的骨架，例如 `internal/logic/misc/ping/ping_logic.go`：

```go
type PingLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PingLogic {
	return &PingLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}
```

每个 Logic 构造函数已经把 `logx.WithContext(ctx)` 嵌进 `logx.Logger` 匿名字段，所有 `l.Errorf(...)`/`l.Infof(...)` 调用天然带上下文——这个模式**保持不变**，Telemetry 接入后不需要改一行 Logic 代码；`ctx` 里的 span context 是 Telemetry 中间件自动放进去的，`logx.WithContext(ctx)` 自动读取。这是选择"配置层接入 + logx 原生集成"而不是"侵入式改造 130+ 个 Logic 文件"的直接原因。

唯一需要改的是六处 `main` 函数里的 `logx.SetUp(...)` 调用，替换成 `logging.Setup("iam-rpc")` 等。

### 3.4 唯一需要新写代码的地方：`user_id` 字段

鉴权中间件（`AuthMiddleware`，拆分后在 gateway 侧）已经把解析出的用户 ID 放进 ctx（通过 `pkg/jwt/context.go` 现有机制）。约定：写操作类 Logic（涉及审计/操作日志的路径，已经在调用 `pkg/audit`）在记录关键日志时，用 `logx.WithContext(ctx).WithFields(logx.Field(logging.FieldUserID, userID))` 附加一次，不要求全部 130+ 个 Logic 文件都补——这是"有真实排查价值的地方才加"，跟 `08-testing-strategy.md` 的"测有 bug 风险的地方，不是测生成代码"是同一个判断标准。

## 4. 明确不做的事

- **不引入日志聚合系统（ELK/Loki/Grafana Loki 等）**。每个服务继续往 stdout 打 JSON（而不是文件），复用现有的 Supervisor（Phase 1-2 过渡期）/`docker compose logs`（Phase 3 之后，见 `21-cd-and-deployment.md`）日志采集方式。JSON 格式已经是"可被任何聚合系统摄入"的标准形态，以后要接 Loki/ELK 只是加一个 Promtail/Filebeat 配置指向已有的 JSON stdout，不需要回头重写日志代码——这就是本轮只做 JSON 化、不做聚合系统选型的原因。
- **不部署 trace 后端（Jaeger/Zipkin/Tempo）**。`Telemetry` 配置块本身在没有 collector 时不报错、不影响服务运行（span 静默丢弃），trace-id 依然会打进日志字段，本轮的排查手段是"用同一个 trace-id grep 四份 JSON 日志"，不依赖可视化调用链界面；接一个真正的 trace 后端是后续可选投入。
- **不做 metrics（Prometheus 指标）体系化接入**。`go.mod` 里已经有 `prometheus/client_golang`（go-zero 间接依赖，`/metrics` 端点 go-zero 框架自带），这部分本轮不额外设计，不在本文档范围内。
- **不要求六个服务一次性全部接入**。可以随 Phase 2 每个服务拆分时顺手做（推荐），也可以在 Phase 2 全部拆完后统一做一次；两种顺序都不影响正确性，按实际进度选。

## 5. 完成的定义

- 六个服务的 `etc/*.yaml` 都有 `Telemetry` 配置块，`Telemetry.Name` 互不相同。
- `pkg/logging/setup.go`、`pkg/logging/fields.go` 落地，六个服务的 `main` 函数（`cmd/gateway/main.go` + 五个 `services/<name>/<name>.go`）改用 `logging.Setup("<service-name>")` 替换原来内联的 `logx.SetUp(logx.LogConf{Encoding: "plain"})`。
- `go build ./...` 通过。
- 冒烟验证：本地起 gateway + 至少一个 RPC 服务（如 `iam-rpc`），发一个会跨两个进程的请求（如登录后调一个需要鉴权的接口），确认 gateway 和 iam-rpc 各自的 stdout 输出都是合法 JSON，且能找到一个共同的 trace 字段值把两边日志关联起来。
- 确认改动没有触碰任何业务逻辑分支——这一步应该是纯配置 + 一个新增小包 + 六处 `main` 函数的机械替换。
