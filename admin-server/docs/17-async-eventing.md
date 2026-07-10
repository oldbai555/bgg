# 17 — 跨服务一致性：TaskCallback + Redis Streams

> 前置依赖：已读 `15-service-boundaries.md`（尤其第 5 节"跨服务引用 ID 处理三种模式"）和 `16-rpc-conventions.md`。本文档展开三种模式里的②异步事件、③Task 回调 RPC 的具体契约，①按需 RPC 查询已经在 16 文档第 4/6 节讲完，不重复。

## 0. 前置依赖与总体思路

两种机制，都基于已有的 Redis（`github.com/redis/go-redis/v9`，`go.mod` 里已经是依赖，零新增中间件）：

- **task-rpc 变成通用任务队列**，`internal/interfaces.AsyncTaskBackend` 从进程内接口升级为 RPC 方法；需要跨服务数据的任务走 `TaskCallback` RPC 契约（第 1 节）。
- **Redis Streams 作为"尽力而为副作用"的事件总线**（第 2 节），不是 `admin_task` 表——`admin_task` 是用户可见的、有状态跟踪的任务，高频内部事件（用户创建、审计日志写入）不该塞进那张表,会污染任务列表语义、也没必要持久化到 MySQL。

一条贯穿全文的判断规则（第 3 节展开）：**现在代码里子操作失败只是 `logx.Errorf` 记录、不影响主流程的，转成异步 Streams；现在代码里子操作失败会让整个请求报错的，转成同步 RPC。** 这个规则不是主观判断，是直接读现有代码的错误处理方式推导出来的——迁移不应该在拆分的同时顺带改变一个操作"失败了要不要紧"的语义，语义在 Phase 1 单体阶段就该是对的（如果发现现在的语义本身不合理，那是 Phase 1 的修复范围，不是 Phase 2 拆分时顺带改)。

## 1. TaskCallback RPC 契约

### 1.1 现状：`AsyncTaskBackend` 接口（`internal/interfaces/task.go`）

```go
type TaskExecutor interface {
	GetType() int
	Execute(ctx context.Context, task *task.AdminTask, paramsJSON string) (string, error)
}

type AsyncTaskBackend interface {
	Submit(ctx context.Context, task *task.AdminTask) (backendTaskID string, err error)
	Cancel(ctx context.Context, task *task.AdminTask) error
	SyncStatus(ctx context.Context, task *task.AdminTask) error
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}
```

现状实现（`internal/domain/task/executors/excel_export_executor.go` 的 `ExcelExportExecutor`）直接持有 `*repository.Repository` 全量句柄——`Execute` 方法内部按 `params.Module`（`operation_log`/`audit_log`/`login_log`/`performance_log`/`sdk_call_log`,对应 `internal/consts` 里的 `TaskModule*` 常量）分支，直连 `monitoringrepo`/`systemrepo`/`sdk` 的 repository 查要导出的数据、生成 Excel/CSV、写文件。这是 `15-service-boundaries.md` 第 3 节标出的 `monitoring → task`（4 个 `*_export_logic.go`）、`sdk → task`（`sdk_call_log_export_logic.go`）两组耦合背后的真正执行者。

### 1.2 拆分后的问题

`task-rpc` 只有 `admin_task` 一张表（`admin_task` schema）。拆分之后 `ExcelExportExecutor` 里 `monitoringrepo.NewOperationLogRepository(...)` 这类调用没法再直接构造——那些 repository 依赖的是 `admin_platform`/`admin_sdk` 的连接，task-rpc 的进程里根本没有这些库的连接配置。`executors map[int]interfaces.TaskExecutor` 因此要改造成 `executors map[int]taskcallback.Client`：每个 module 对应一个"回调哪个服务"的客户端，而不是直接持有 repository。

### 1.3 `TaskCallback` proto 契约

每个拥有导出数据的服务（当前是 `iam-rpc` 承载 monitoring 的 4 个导出模块、`sdk-rpc` 承载 `sdk_call_log` 导出）实现同一份 `TaskCallback` 接口,task-rpc 通过一张静态路由表（`task_type` 或更细粒度的 `module` → 服务地址）决定回调谁：

```protobuf
// pkg/taskcallback/taskcallback.proto —— 每个"承载可导出数据"的服务都实现这份接口
syntax = "proto3";

package taskcallback;

option go_package = "./taskcallback";

message FetchExportDataRequest {
  string module = 1;        // 对应 consts.TaskModuleOperationLog 等取值
  uint64 task_id = 2;       // admin_task.id，用于回调方自己记录/幂等
  uint64 requested_by = 3;  // 发起导出的用户 ID（对应现有 ExcelExportParams 里携带的 user 信息）
  string filters_json = 4;  // 原样透传现有 ExcelExportParams.Filters（JSON 编码后的筛选条件）
}

message FetchExportDataResponse {
  // 导出数据以行记录方式返回，每行是字段名到字段值的 map（JSON 编码字符串，
  // 避免为每个 module 单独定义一份 proto message，保持这个通用回调接口的稳定）
  repeated string rows_json = 1;
  repeated string headers = 2; // 列名，用于 task-rpc 侧统一生成 Excel/CSV 表头
  int64 total_count = 3;
}

service TaskCallback {
  rpc FetchExportData(FetchExportDataRequest) returns (FetchExportDataResponse);
}
```

`task-rpc` 里的路由表（静态 map，不需要动态注册,新增导出 module 时手动加一行,这与 `16-rpc-conventions.md` 第 5 节"不引入 etcd 服务发现"的简单化取向一致）：

```go
// services/task/internal/logic/exec/module_route.go
var moduleServiceRoute = map[string]taskcallback.Client{
	consts.TaskModuleOperationLog:   iamCallbackClient,
	consts.TaskModuleAuditLog:       iamCallbackClient,
	consts.TaskModuleLoginLog:       iamCallbackClient,
	consts.TaskModulePerformanceLog: iamCallbackClient,
	consts.TaskModuleSdkCallLog:     sdkCallbackClient,
}
```

`ExcelExportExecutor.Execute` 的新形态：不再直连 repository，而是查路由表拿到对应服务的 `taskcallback.Client`，调 `FetchExportData`，拿到 `rows_json`/`headers` 后按现有逻辑生成 Excel/CSV 文件——生成文件、写 `admin_task.result` 这部分逻辑完全不变，只是"怎么拿到要导出的数据"从直连 repository 换成一次 RPC。`iam-rpc`/`sdk-rpc` 侧新增一个 `TaskCallback` 的 server 实现，内部就是把原来 4 个 `*_export_logic.go`（现在只负责"创建任务记录"）里散落的查询逻辑集中实现一遍。

### 1.4 提交路径（`monitoring/sdk → task` 的同步 RPC 部分）

`*_export_logic.go` 现在的行为是"创建 `admin_task` 记录，创建失败直接 `return nil, errs.Wrap(...)`（错误会一路返回给前端）"——按第 3 节的判断规则,这是同步 RPC,不是 Streams。拆分后 `iam-rpc`/`sdk-rpc` 的导出 logic 改成调 `task-rpc.SubmitTask`（对应 `AsyncTaskBackend.Submit` 的 RPC 化），task-rpc 建好 `admin_task` 记录后同步返回 `task_id`，调用方原样返回给前端（前端继续用现有的任务列表机制轮询状态,不需要改前端）。

## 2. Redis Streams：两个具体的流

### 2.1 `stream:chat.user.created`

**生产者**：`iam-rpc` 的用户创建事务提交后发布。对应 `internal/logic/iam/user/user_create_logic.go` 的 `initChatForNewUser`——现有实现里这段代码本身已经是"失败只记日志、不回滚建用户"的语义（`logx.Errorf(...)` 后继续执行，方法整体返回的 error 在调用方 `UserCreate` 里也只是 `logx.Errorf` 记录，不影响 `UserCreate` 本身返回成功），完全匹配第 3 节的异步判断规则,不需要改变行为，只改变触发机制：

```go
// 拆分前（Phase 1 单体内）：
if err := l.initChatForNewUser(user.Id); err != nil {
    logx.Errorf("初始化新用户聊天数据失败: %v", err)
}

// 拆分后（iam-rpc 内）：
event := ChatUserCreatedEvent{UserID: user.Id, CreatedAt: time.Now().Unix()}
payload, _ := json.Marshal(event)
if _, err := l.svcCtx.Redis.XAdd(l.ctx, &redis.XAddArgs{
    Stream: "stream:chat.user.created",
    Values: map[string]interface{}{"payload": string(payload)},
}).Result(); err != nil {
    // 发布失败同样只记日志，不回滚用户创建——与现有语义一致
    logx.Errorf("发布 chat.user.created 事件失败: %v", err)
}
```

**消费者**：`chat-rpc` 起一个消费者组（`XGROUP CREATE stream:chat.user.created chat-rpc-init $ MKSTREAM`，消费者组名固定为 `<service>-<用途>` 格式，这里是 `chat-rpc-init`），消费到事件后执行"拉入默认群 + 为存量用户逐个建私聊"——这就是 `initChatForNewUser` 原来的方法体，原样搬进 chat-rpc 自己的领域服务（`services/chat/internal/domain/onboarding.go` 或类似位置），包括第 1 步查默认群、第 2 步 `FindChunk` 分批拉用户建私聊的逻辑都不变。**唯一需要补的是 Phase 1 里已经规划好的 `FindChunk` 替换全量 `FindPage(1, 10000, "")`**（`04-domain-iam-chat.md` 任务 1/2，`ChatOnboardingService.createPrivateChatsForExistingUsers` 用 `chatdomain.UserLister.FindChunk` 分批）——这个修复在 Phase 1 单体阶段就该做完，Phase 2 只是原样搬迁已经修好的版本，不是在拆分时顺带修。

### 2.2 `stream:audit.request`

**生产者**：gateway 的日志中间件（`OperationLogMiddleware`/`PerformanceMiddleware`）非阻塞 `XADD`，每次请求处理完成后发布一条事件,携带请求方法/路径/耗时/状态码/用户信息，等同于现在这两个中间件直接写 `admin_operation_log`/`admin_performance_log` 两张表时携带的字段。

**消费者**：`iam-rpc` 起消费者组（`iam-rpc-audit-writer`）批量写入自己的 `admin_operation_log`/`admin_performance_log` 表——**这就是让 `15-service-boundaries.md` 第 1.1 节"monitoring 合并进 iam-rpc"这个决定在没有同步 RPC 开销的前提下依然成立的机制**：如果 monitoring 是独立服务,每个 HTTP 请求都要多一次同步 RPC 只为了写日志；现在 gateway 只需要一次非阻塞 `XADD`（不等待任何服务确认），iam-rpc 按自己的节奏批量消费写入,网关侧的请求延迟不受这条日志写入路径影响。

批量写入的具体方式：消费者用 `XREADGROUP` 一次拉一批（如 `COUNT 100`），攒够一批或超时（如 200ms）后用一次事务批量 `INSERT`，减少 MySQL 往返次数——这是对现状"每次中间件调用各自 `Create` 一条记录"的一个自然优化，不是本文档强制要求的实现细节，具体 batch size/超时阈值留给实现时按实际负载调整。

### 2.3 幂等性要求

Streams 是至少一次投递（consumer 崩溃重启、`XACK` 之前进程被杀等场景都可能导致同一条消息被消费两次），消费者必须幂等。复用现有代码里已经在用的模式——`initChatForNewUser` 里"创建私聊前先查是否已存在"（`chatRepo.FindPrivateChatByUserIDs` 查到就跳过）、"加入群组前先查是否已在群里"（遍历 `chatUserRepo.FindByChatID` 结果比对 `UserId`）,这两处"插入前查是否已存在"的写法直接照搬到 Streams 消费者里,不需要引入新的幂等框架（如给每条 Stream 消息生成全局唯一 ID 再建一张去重表）——现有查询本身自带的唯一性约束（`chat_user` 表的 `uk_chat_user (chat_id, user_id)` 唯一键）已经能兜底真正的并发重复,消费者代码里的"先查后插"只是减少不必要的失败重试。

### 2.4 死信处理

消费者处理失败超过 N 次（建议 N=3，具体次数留给实现时按事件重要性调整,不做成可配置项，硬编码一个常量即可）后，把这条消息移到 `stream:<name>.deadletter`（如 `stream:chat.user.created.deadletter`），记一条 `logx.Errorf`，**不做更复杂的 DLQ**——不建告警、不建自动重放机制,人工发现问题时去查 deadletter 流里的内容手动处理即可。这与 `11-descoped.md` "不做日志聚合系统"的取向一致：先把最简单能工作的机制落地,复杂度按真实运维经验需要再加。

## 3. 判断规则与草案 disposition 表

### 3.1 判断规则（复述，供本节表格套用）

现在代码里子操作失败只是 `logx.Errorf` 记录、不影响主流程的 → **异步 Streams**；现在代码里子操作失败会让整个请求报错的 → **同步 RPC**。这条规则直接从现有代码行为读出来，不是主观判断。

### 3.2 草案 disposition 表

**这是草案，不是最终结论**——Phase 1 完成后才会有确定的"跨 ≥2 表写、跨仓储读写、跨域、或有非平凡业务规则"的 35-40 个域服务文件清单（`01-architecture-target.md` A.2 节"分层判断标准"一节的估算）。本表用一个可以在 Phase 1 之前跑的代理方法生成：对当前（Phase 1 尚未开始的）代码库跑 `grep -oE '[a-zA-Z_]*repo\.New[A-Za-z]*Repository' <file> | sort -u` 找出"调用了 ≥2 个不同 repository 构造函数"的 logic 文件，命中约 40 个文件（数量级与 Phase 1 估算的 35-40 吻合，但具体文件集合不完全相同——代理方法找到的是"多 repository 调用"，Phase 1 真正的标准还包括跨表写/非平凡业务规则等本文档无法用 grep 判断的维度）。下表只挑出其中**真正跨越 `15-service-boundaries.md` 第 1 节服务边界**的条目（同域内的多 repository 调用，如 `chat_group_create_logic.go` 调用 3 个 chat 域自己的 repository，属于 Phase 1 的域服务分层范畴，不属于 Phase 2 的跨服务同步/异步问题，不列入本表）：

| 方法 / 文件 | 目标服务方向 | 同步 RPC / 异步 Streams | 判断依据 |
|---|---|---|---|
| `UserCreate.initChatForNewUser`（`iam/user/user_create_logic.go`） | iam-rpc → chat-rpc | **异步 Streams**（`stream:chat.user.created`） | 现有代码：`initChatForNewUser` 失败只 `logx.Errorf`，不影响 `UserCreate` 返回成功 |
| `Login.recordLoginLog`（`iam/auth/login_logic.go`） | iam-rpc 内部（monitoring 已合并进 iam-rpc，不跨服务） | 不适用——保持进程内 `go func(){...}()` 异步写，不用 Streams | 拆分后 monitoring 与 iam 同进程，不需要跨进程机制；现状本身已经是 goroutine 异步 |
| `Login.createUnreadNoticeNotifications`（`iam/auth/login_logic.go`） | iam-rpc 内部（system 已合并进 iam-rpc，不跨服务） | 不适用，同上 | 同上 |
| `NoticeCreate/Update.createNotificationsForAllUsers`（`system/notice/*.go`） | iam-rpc 内部（system → iam 同进程） | 不适用，同上 | 同上；现状已经是 `go func(){...}()`，失败只记日志 |
| `OperationLogExport`/`AuditLogExport`/`LoginLogExport`/`PerformanceLogExport`（`monitoring/*_export_logic.go`，4 个文件） | iam-rpc → task-rpc（提交） | **同步 RPC**（`task-rpc.SubmitTask`） | 现有代码：`taskRepo.Create` 失败直接 `return nil, errs.Wrap(...)`，错误会返回给前端 |
| 同上 4 个方法对应的实际导出执行 | task-rpc → iam-rpc（取数） | **同步 RPC**（`TaskCallback.FetchExportData`，见第 1 节） | 任务执行必须拿到数据才能生成文件，取数失败应使任务状态置为失败（`TaskStatusFailed`），不能静默跳过 |
| `SdkCallLogExport`（`sdk/sdk/sdk_call_log_export_logic.go`） | sdk-rpc → task-rpc（提交）+ task-rpc → sdk-rpc（取数） | **同步 RPC**（两段都是） | 与 monitoring 导出同构，理由相同 |
| `ChatList`/`ChatMessageList`/群相关 8 个文件的用户资料查询 | chat-rpc → iam-rpc | **同步 RPC**（`BatchGetUserProfiles`，非 Streams） | 这不是"失败可忽略的副作用"，是请求路径上必须拿到的展示数据，属于 `15-service-boundaries.md` 模式①（RPC 查询+缓存），不适用本文档的同步/异步二分规则（那条规则针对的是"要不要做"的副作用，不针对"请求本身需要的数据"） |
| `PublicBlogAuthorInfo`（`blog/public/public_blog_author_info_logic.go`） | content-rpc → iam-rpc | **同步 RPC**（`GetUserProfile`），模式① | 同上，展示数据不是可选副作用 |
| `TaskRecent.getRecentTaskLimit`（`task/public/task_recent_logic.go`） | task-rpc → iam-rpc（可选） | **建议规避（见 `15-service-boundaries.md` 第 5 节末尾）**，退而求其次是同步 RPC 但容忍失败降级 | 现状：缓存/字典查不到时静默 fallback 到硬编码默认值 10，不返回错误——如果保留 RPC，应实现成"RPC 失败/超时也直接 fallback，不让 `TaskRecent` 整体失败"，即同步调用但错误处理策略是"降级不是报错" |

上表覆盖了当前能确认的跨服务方法；Phase 1 的域服务清单最终定稿后，如果发现新的跨服务候选（Phase 1 主要在同域内拆领域服务，理论上不应该新增跨域耦合，但如果拆分过程中发现遗漏，按第 3.1 节规则补充到本表，不需要重新走一遍这里的分析流程）。

## 4. 非目标

- 不引入 Kafka/RabbitMQ——跨服务异步场景全部用 Redis Streams 承载。
- 不做通用的分布式事件总线抽象（发布订阅框架、事件版本管理）——只有本文档列出的 2 个具体 Stream，新增 Stream 时直接照本文档的格式写清楚生产者/消费者/幂等/死信策略，不需要先设计一层通用抽象。
- 不为 Streams 引入监控告警（消费延迟、死信数量的 Dashboard）——`19-observability.md` 的结构化日志是唯一的可观测性手段，Streams 本身的运维状态通过 `XINFO GROUPS`/`XLEN` 手动查询即可。
- 不做 `TaskCallback` 的动态路由注册——`moduleServiceRoute` 是硬编码的 Go map，新增导出 module 时改代码重新部署，不做配置热更新。

## 5. 完成的定义

- `pkg/taskcallback/taskcallback.proto` 编译通过，`iam-rpc`/`sdk-rpc` 各自实现了 `TaskCallback` server，`task-rpc` 的 `moduleServiceRoute` 路由表覆盖第 1.3 节列出的全部 5 个 module。
- `stream:chat.user.created`、`stream:audit.request` 两个 Stream 的生产者/消费者都跑通一次端到端验证：建一个新用户能在 `chat-rpc` 侧观察到默认群加入 + 私聊建立；发一次 HTTP 请求能在 `iam-rpc` 的 `admin_operation_log`/`admin_performance_log` 表里查到对应记录。
- 幂等验证：故意让同一条 `stream:chat.user.created` 消息被消费两次（如手动 `XCLAIM` 重放），确认不会产生重复的群成员/私聊记录。
- 死信路径验证：故意让消费者对某条消息连续失败 N 次，确认消息出现在对应的 `.deadletter` 流里，且不会无限重试阻塞后续消息处理。
- 第 3.2 节的草案 disposition 表已经和 Phase 1 实际产出的域服务清单核对过一次，差异（如有）已经补充记录。
