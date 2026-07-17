# 15 — 服务边界与数据归属（Phase 2 总纲）

> 本文档是 Part B（微服务拆分）的前置文档，`16-rpc-conventions.md`、`17-async-eventing.md`、`18-service-extraction-runbook.md` 都以本文档定义的服务边界、schema 归属、目录结构为前提，改动这三份文档时如果与本文档冲突，以本文档为准，并回来同步修正。

## 0. 前置依赖

- Part A（`01`~`14`）已经落地：`internal/domain/<domain>` 的包边界已经按本文档第 1 节的 5 服务分组组织（不是按 9 个业务域各自一个包），`registry.Transact`/`Wire` 组合根/密钥管理已经就绪。
- 读者已了解 DDD-lite 分层的历史背景（原任务书已删除，决策摘要见 `docs/changelog/archive-backend.md` §15）——本文档假定 9 个业务域（iam/blog/video/chat/sdk/task/monitoring/system/misc）在代码里已经是 `internal/{handler,logic,repository,model}/<domain>/` 的既成事实。
- `AGENTS.md` 里对 `admin-server/**` 的后端规范（squirrel-only SQL、`deleted_at` 软删除、snake_case 命名、Group 格式 `<domain>/<module>`、错误处理走 `pkg/errs`、中间件声明顺序等）在拆分成微服务之后**继续对每一个新 RPC 服务生效**——这些规范描述的是"如何写 Go 代码/如何设计表结构"，与代码运行在单体进程还是独立进程无关。拆分只改变部署单元和数据库边界，不改变编码规范，这一点在每个服务的 `services/<name>/README` 或包注释里都要复述一次，避免执行 Phase 2 的人误以为拆分服务是"重新设计"的机会。

## 1. 服务边界：6 个独立部署单元，不是机械的 9 个

```
gateway        HTTP 唯一入口，无状态，不持有自己的数据库
iam-rpc        承载：iam + system + monitoring + misc  → DB: admin_platform
content-rpc    承载：blog + video                       → DB: admin_content
chat-rpc       承载：chat（含 hub）                       → DB: admin_chat
task-rpc       承载：task（通用异步任务/队列基础设施）        → DB: admin_task
sdk-rpc        承载：sdk（admin+public+调用日志）           → DB: admin_sdk
```

9 个业务域 → 5 个数据服务 + 1 个无状态网关，一共 6 个独立部署单元。**拆分粒度按"是否有独立扩缩容/独立生命周期/独立信任边界的真实理由"决定，不是"一个业务域一个服务"的机械映射**——9 个服务对一个独立维护者是不必要的运维负担（要多盯 9 套日志、9 套健康检查、9 个 docker-compose 条目），这与用户明确要求的"做简单"直接冲突；反过来把所有域塞进一个单体又放弃了微服务拆分本该带来的独立部署/独立扩缩容收益。下面逐条给出每个合并/拆分决定的具体理由，理由本身来自对现有代码行为的观察，不是预先设定好答案再找理由。

### 1.1 iam + system + monitoring + misc → iam-rpc

**核心理由：monitoring 的写入路径和请求路径是同构的，拆开只会制造无意义的同步 RPC。**

`admin_operation_log`、`admin_performance_log` 的写入分别挂在 `OperationLogMiddleware`、`PerformanceMiddleware` 上——这两个中间件在**每一次 HTTP 请求**上都会触发一次写入（前者记录谁在什么时间调用了哪个接口、参数是什么，后者记录耗时和是否慢查询）。如果把 monitoring 拆成独立服务，网关处理的每一个请求，无论业务本身是否需要跨服务调用，都必须额外发起一次同步 RPC 只是为了把日志写进 monitoring 自己的库——这是纯负担、没有对应收益，因为日志写入本身在语义上就该是"尽力而为、最终一致"的操作（丢一条操作日志不影响任何业务正确性）。把 monitoring 和 iam 放进同一个服务、同一个数据库连接后，这条写入路径可以走进程内异步（Redis Streams，见 `17-async-eventing.md`），完全不需要跨进程网络跳转。

`system`、`monitoring` 体量也支持这个决定：`system` 域约 30 个 logic 文件、`monitoring` 域约 20 个 logic 文件（本文档第 2 节有精确的表清单），都是低频后台配置/日志查询场景，QPS 特征和 iam 高度一致（都是登录态管理员的后台管理操作），没有独立扩缩容的理由。而且 `system/notice` 的通知创建逻辑（`internal/logic/system/notice/notice_create_logic.go` 的 `createNotificationsForAllUsers`）本来就直接 `import "postapocgame/admin-server/internal/repository/iam"` 去分页拉取全量用户 ID 来批量建通知——这是一处真实存在的跨域直接耦合（详见第 3 节的核查结果），拆进同一个服务之后，这类耦合从"跨进程 RPC"降级为"进程内函数调用"，不需要重新设计。

`misc`（`demo`、`daily_short_sentence`，共 2 张表、约 8 个 logic 文件）本身没有明确的领域归属，体量太小不值得单独占用一个部署单元，并入 iam-rpc 作为兜底工具端点。

### 1.2 blog + video → content-rpc

两个域的公开面（`blog/public`、`video/public`、`video/m3u8`）都是免登录的公共只读接口，读 QPS 天然比后台管理接口高、且随外部流量（爬虫、真实访客、突发热点内容）波动更剧烈——这是一个真实的独立扩缩容理由：content-rpc 需要能单独加副本应对读流量峰值，而不必连带扩容整个后台系统。但 blog 和 video 彼此之间架构非常简单：都是"单表/几张关联表的 CRUD + 一个公开只读视图"，没有 chat 那种服务端状态、没有 sdk 那种独立信任边界，合并成一个服务不会制造耦合，反而避免了两个几乎同构的微小服务各自维护一套 etc.yaml/Dockerfile/健康检查的重复成本。

### 1.3 chat → chat-rpc

chat 是 9 个域里**唯一有真正服务端状态**的域：`internal/hub/chathub.go` 维护的是进程内存里的在线连接表和消息广播逻辑，重启这个进程会打断所有当前建立的 WebSocket 连接。这种"长连接 + 内存态"的失败模式和扩缩容特征与其余所有域（无状态的请求-响应式 CRUD）完全不同——chat-rpc 挂掉或需要重启升级时，只应该影响正在聊天的用户，不该拖累后台管理员登录、博客发布这些完全不相关的操作。这是一个独立信任边界+独立生命周期的双重理由，拆分收益明确。

### 1.4 task → task-rpc

`internal/interfaces.AsyncTaskBackend` 从设计之初就是面向可插拔中间件的接口（注释里明确写了"对接内置调度器/xxl-job/quartz/k8s CronJob 等"），调度器现有的 `acquireLock`/`releaseLock` 用 Redis `SETEX` 实现，已经是多副本安全的分布式锁语义，不需要为了拆分重新设计并发模型。也就是说代码里已经有一条现成的拆分缝——把 `AsyncTaskBackend.Submit` 从进程内接口调用升级为 RPC 方法，架构不需要变,只是调用方式变了。task 域本身只有 1 张表（`admin_task`）、约 4 个 logic 文件，是天然的最小拆分单元（这也是它被选为 Part B.6 第一个拆分对象的直接原因）。

### 1.5 sdk → sdk-rpc

sdk 是理由最充分的一个拆分：
- **已有独立的中间件链**：`SDKAuthMiddleware`/`SDKRateLimitMiddleware`/`SDKCallLogMiddleware`，与后台管理走的 `AuthMiddleware`/`PermissionMiddleware`/`OperationLogMiddleware` 完全是两套体系，没有共享状态。
- **信任边界不同**：sdk 的调用方是持有 API Key/Secret 的外部系统，不是 JWT 登录态的后台管理员，两者的认证机制、限流策略、审计要求天然应该在代码和部署层面都分开。
- **未来最可能需要独立的网络暴露方式**：sdk-rpc 一旦有真实外部调用方接入，很可能需要独立域名/独立网关限流策略/独立 SLA，现在就用独立服务边界可以让这条演进路径无痛发生，不需要事后从单体里再挖一次。

## 2. 表清单核实（对照 `db/tables.sql` + `db/migrations/*.sql` + `db/demo/*.sql` 实测结果）

对计划草案里给出的"iam 10、system 6、monitoring 4、chat 3、sdk 4、video 1、blog 6、task 1、misc 2"逐条核实，方法是数 `internal/model/<domain>/` 下除 `vars.go` 外的每个 `*model.go` 文件（每个文件对应一张表的 goctl Model），再用 `internal/repository/<domain>/` 下的手写 repository 文件交叉核对。

| 域 | 草案表数 | 实测表数 | 实测表清单 | 差异说明 |
|---|---|---|---|---|
| iam | 10 | **10**（核实无误） | admin_user, admin_role, admin_permission, admin_department, admin_user_role, admin_role_permission, admin_menu, admin_permission_menu, admin_api, admin_permission_api | `repository/iam/` 还有第 11 个文件 `token_blacklist_repository.go`，但它没有对应的表/Model——100% 基于 Redis（见 16/17 文档），不计入表数 |
| system | 6 | **6**（核实无误） | admin_config, admin_dict_type, admin_dict_item, admin_file, admin_notice, admin_notification | `internal/model/system/filemodel.go` 是一个孤立的、未被任何代码引用的 `FileModel` 定义（与 `adminfilemodel.go` 的 `AdminFileModel` 重复，疑似历史遗留的生成残留），拆分时直接删除，不要当成第 7 张表 |
| **monitoring** | 4 | **5**（草案少算 1 张） | admin_operation_log, admin_login_log, audit_log, admin_performance_log, **metric_daily_stats** | `metric_daily_stats` 由 `internal/repository/monitoring/metric_repository.go` 直接用 squirrel 拼 `INSERT ... ON DUPLICATE KEY UPDATE` 操作，**没有走 goctl Model**，因此在只数 `model/monitoring/*.go` 文件时会被漏掉；建表语句在 `db/migrations/create_table_metric.sql`。这张表记录的是博客/视频公开页的 PV/UV/VV 日度统计，语义上更贴近 monitoring 而不是 content，随 monitoring 一起并入 iam-rpc/admin_platform |
| chat | 3 | **3**（核实无误） | chat, chat_user, chat_message | — |
| sdk | 4 | **4**（核实无误） | sdk_key, sdk_interface, sdk_key_api, sdk_call_log | — |
| video | 1 | **1**（核实无误） | video | — |
| blog | 6 | **6**（核实无误） | blog_tag, blog_article, blog_article_tag, blog_article_audit, blog_friend_link, blog_social_info | 前 4 张建表语句在 `db/migrations/create_table_blog.sql`，后 2 张在 `create_table_blog_extension.sql`（历史上分两批加的） |
| task | 1 | **1**（核实无误） | admin_task | 建表语句在 `db/migrations/create_table_task.sql` |
| misc | 2 | **2**（核实无误） | daily_short_sentence, demo | `demo` 表的建表/初始化 SQL 单独放在 `db/demo/`（不在 `db/migrations/` 下），是脚手架自带的演示模块，拆分时按普通表处理即可 |

**结论**：全库真实表数是 **38 张**（10+6+5+3+4+1+6+1+2），比草案统计的 37 张多 1 张（`metric_daily_stats` 漏计）。`db/tables.sql` 本身只收录了 29 张核心表（iam 全部 10 张、system 6 张、monitoring 的 4 张日志表、chat 3 张、misc 的 `daily_short_sentence`、sdk 4 张、video 1 张），blog 的 6 张、`metric_daily_stats`、`admin_task`、`demo` 分别放在 `db/migrations/*.sql` 和 `db/demo/*.sql` 里——这是历史遗留的文件组织方式（先有 `tables.sql` 做的首批建表，后续新增域改成了增量迁移文件），第 4 节的 `db/services/` 目录重组会把这些分散的建表文件统一收拢，届时不再有"核心表在 tables.sql、新表在 migrations/"这种区分。

再次确认：仓库里没有任何 `FOREIGN KEY` 约束（`db/tables.sql` + `db/migrations/*.sql` + `db/demo/*.sql` 全零），"拆库后外键失效"这个问题在 schema 层面不存在。

## 3. 跨域越界代码清单（实测，比计划草案更完整）

用 `grep -rl "postapocgame/admin-server/internal/repository/<other>\"" internal/logic/<domain>` 逐个域对逐个域核查（只认真实 import 路径，不是简单字符串匹配），完整结果：

| 越界方向 | 文件 | 说明 |
|---|---|---|
| iam → system, monitoring | `internal/logic/iam/auth/login_logic.go` | 登录成功后异步创建"未读公告通知"（读 `system/notice`、写 `system/notification`），以及异步写登录日志（写 `monitoring/login_log`）。两处都已经是 `go func(){...}()` + 内部 `recover`，失败只 `logx.Errorf`，不影响登录主流程 |
| iam → chat | `internal/logic/iam/user/user_create_logic.go` | `initChatForNewUser`：建用户成功后同步调用（非 goroutine，但错误只记日志），把新用户拉进默认群 + 为存量用户逐个建私聊 |
| system → iam | `internal/logic/system/notice/notice_create_logic.go`, `notice_update_logic.go` | 公告发布后异步给全量用户创建通知，需要 `iamrepo.NewUserRepository(...).FindChunk` 分批拉用户 |
| monitoring → task | `internal/logic/monitoring/{operation_log,login_log,audit_log,performance_log}/*_export_logic.go`（4 个文件） | 导出改成了"创建异步任务记录"模式：调 `taskrepo.NewTaskRepository(...).Create` 建一条 `admin_task` 记录，真正的导出由 `internal/domain/task/executors/excel_export_executor.go` 的 `ExcelExportExecutor` 异步执行——而这个 executor 目前直接持有完整的 `*repository.Repository`，直连 `monitoringrepo`/`systemrepo` 查询要导出的数据 |
| sdk → task | `internal/logic/sdk/sdk/sdk_call_log_export_logic.go` | 同上模式，导出 `sdk_call_log` |
| chat → iam | `internal/logic/chat/{chat,group,message}/*.go`（8 个文件：`chat_list_logic.go`、`chat_message_list_logic.go`、`chat_group_member_list_logic.go`、`chat_group_detail_logic.go`、`chat_group_member_add_logic.go`、`chat_message_list_admin_logic.go`、`chat_group_create_logic.go`、以及创建流程里的部门/角色校验） | 展示聊天列表/群成员时需要拉用户昵称、头像（`iamrepo.NewUserRepository`），群创建时校验部门/角色是否存在（`iamrepo.NewDepartmentRepository`、`NewRoleRepository`、`NewUserRoleRepository`） |
| blog → iam | `internal/logic/blog/public/public_blog_author_info_logic.go` | 公开博客页展示"作者信息"，固定读 `iamrepo.NewUserRepository(...).FindByID(ctx, 1)`（写死 userID=1，即站长账号） |
| task → system | `internal/logic/task/public/task_recent_logic.go` | 读取"最近任务展示条数"配置，先查 `l.svcCtx.Repository.BusinessCache`（本地 Redis 缓存），未命中再读 `systemrepo` 的字典表兜底，字典也查不到再用硬编码默认值 10 |

**与计划草案的差异**：草案原文只列了 `chat/*`、`blog/public/public_blog_author_info_logic.go`、`system/notice/*`、`sdk/sdk_call_log_export_logic.go`、4 个 `monitoring/*_export_logic.go`，本次核查在此基础上补充确认了 `iam → system/monitoring/chat`（login_logic.go 和 user_create_logic.go 本身发起的跨域调用）以及 `task → system`（task_recent_logic.go）。这两组遗漏不影响拆分决策本身，但会影响 Phase 2 每个服务的 RPC 契约设计——尤其是 `iam → chat`（登录/建用户流程内发起，最终变成第一个要落地的 Redis Streams 生产者）和 `task → system`（这是本文档第 5 节要单独讨论的一个可以规避掉的伪跨服务调用）。

一个重要的观察：`iam ↔ system`、`iam ↔ monitoring` 这两组耦合在 Phase 2 落地后**会自动消失**——因为 B.1 已经决定 iam、system、monitoring 三个域合并进同一个 `iam-rpc` 服务，`login_logic.go`、`notice_create_logic.go` 里的这些"跨域" import 拆分后仍然是同一个 Go 二进制内的包引用，不需要改成 RPC 调用，也不需要经过 `17-async-eventing.md` 的 Streams 机制——它们本来就该继续走进程内调用（异步与否取决于原来的 goroutine 用法是否保留，与是否跨服务无关）。真正需要在 `16-rpc-conventions.md`/`17-async-eventing.md` 里设计成 RPC 或 Streams 的，只有跨越 B.1 服务边界的那几条：`iam-rpc → chat-rpc`（异步 Streams）、`iam-rpc/sdk-rpc → task-rpc`（同步 RPC 提交任务）、`task-rpc → iam-rpc/sdk-rpc`（TaskCallback RPC 取导出数据）、`chat-rpc → iam-rpc`（同步 RPC 查用户资料）、`content-rpc → iam-rpc`（同步 RPC 查作者信息）。

## 4. 数据归属设计

**每个 RPC 服务从第一天起就有自己独立的 MySQL schema**（逻辑隔离现在就做对，物理隔离——独立主机——可以之后再做，不影响现在）：5 个 schema `admin_platform`/`admin_content`/`admin_chat`/`admin_task`/`admin_sdk`，初期共用同一台 MySQL 物理实例，`gateway` 不持有数据库、不直连 MySQL。

`db/tables.sql`/`db/migrations/`/`db/demo/` 按服务拆分，目录结构确定为三层——公共目录 `db/services/` → 服务（`iam`/`content`/`chat`/`task`/`sdk`）→ 业务功能/模块（与现有 `scripts/generate-sql.sh -group <domain>/<module>` 的模块粒度一致），每个模块目录下放 `create_table_<module>.sql`、`init_<module>.sql`、`migrations/`。这个结构已经和用户确认过，**逐字保留，不做调整**：

```
db/
└── services/                          # 公共目录
    ├── iam/                           # 服务
    │   ├── user/                      # 业务功能/模块
    │   │   ├── create_table_user.sql
    │   │   ├── init_user.sql
    │   │   └── migrations/
    │   ├── role/
    │   ├── permission/
    │   ├── menu/
    │   ├── department/
    │   ├── dict/
    │   ├── config/
    │   ├── file/
    │   ├── notice/
    │   ├── notification/
    │   ├── operation_log/
    │   ├── login_log/
    │   ├── audit_log/
    │   ├── performance_log/
    │   ├── demo/
    │   └── daily_short_sentence/
    ├── content/
    │   ├── blog/         (tag/article/article_tag/article_audit/friend_link/social_info 各自子目录或合并，按现有表粒度定)
    │   └── video/
    ├── chat/
    │   └── chat/          (chat/chat_user/chat_message)
    ├── task/
    │   └── task/
    └── sdk/
        └── sdk/           (sdk_key/sdk_interface/sdk_key_api/sdk_call_log)
```

对照第 2 节的实测表清单，`db/services/iam/` 下要补一个第 2 节新发现的模块目录 `metric/`（对应 `metric_daily_stats`，随 monitoring 一起进 `admin_platform`）——这是本文档相对计划草案原始目录树的唯一必要补充,其余目录逐字照抄计划草案，不再调整。

`scripts/generate-sql.sh` 加一张固定的 域→服务 映射表（6 条：iam/system/monitoring/misc → iam，content 对应 blog/video，chat/task/sdk 各自对应自己），让新表按 `-group <domain>/<module>` 的 `<domain>` 自动落进 `db/services/<service>/<module>/`，服务归属只在这一张映射表里维护一次，不散落在别处。

**这一层目录拆分在 Phase 1 就可以先做**（不用等 Phase 2 才动），因为它只是把现有 SQL 文件重新归档，不影响单体运行——单体阶段所有服务的库仍然连的是同一个 MySQL 实例/同一套连接配置，只是 SQL 源文件已经按未来的服务边界分好目录，Phase 2 拆库时直接对应搬迁，不用重新梳理归属。

## 5. 跨服务引用 ID 处理：三种模式

拆分之后，一个服务的数据行引用另一个服务拥有的实体 ID（例如 `chat_message.from_user_id` 引用 `iam-rpc` 的 `admin_user.id`）不能再指望同库 JOIN 校验，必须显式选择以下三种模式之一：

**① 按需 RPC 查询 + 本地缓存**——用于"展示型"引用：拿到 ID 之后需要展示对方服务里的可读信息（用户昵称、头像、作者信息），但不需要实时强一致。典型例子：`chat-rpc` 展示聊天列表/群成员时调 `iam-rpc.BatchGetUserProfiles` 批量拉用户资料（对应第 3 节 `chat → iam` 的 8 个文件）；`content-rpc` 的公开博客页展示作者信息调 `iam-rpc.GetUserProfile`（对应 `public_blog_author_info_logic.go`）。查询结果按 TTL 做本地缓存（复用 `pkg/cache` 的模式），减少热路径的 RPC 次数。

**② 异步事件（Redis Streams）**——用于"建立型"引用：一个服务的写操作需要在另一个服务里连带建立关联数据，但这个连带操作允许最终一致、允许失败后只记日志。典型例子：`iam-rpc` 的用户创建事务提交后，发布 `stream:chat.user.created` 事件，`chat-rpc` 消费后异步执行"拉入默认群 + 建私聊"（对应第 3 节 `iam → chat` 的 `initChatForNewUser`，这个方法当前的实现本身已经是"失败只记日志、不回滚建用户"的语义，天然适合搬进 Streams，不需要改变行为，只是改变触发机制）。完整设计见 `17-async-eventing.md`。

**③ Task 回调 RPC**——用于"task-rpc 执行任务时需要读别的服务的数据"这种特殊情况。`task-rpc` 本身不持有 `admin_platform`/`admin_sdk` 的数据（它只有 `admin_task` 一张表记录任务元数据和结果），但导出类任务（操作日志导出、SDK 调用日志导出）必须读到实际要导出的业务数据。现状是 `ExcelExportExecutor` 直接持有 `*repository.Repository` 全量句柄、直连 `monitoringrepo`/`systemrepo`（详见第 3 节的 `monitoring → task`、`sdk → task` 两组耦合），拆分后 task-rpc 拿不到这些库,必须反过来在执行任务时回调发起方服务的一个 `TaskCallback` RPC 取数据。完整契约设计见 `17-async-eventing.md` 第 2 节。

**一处可以规避掉的伪跨服务调用**：第 3 节提到的 `task → system`（`task_recent_logic.go` 读"最近任务展示条数"配置）表面上是模式①的候选，但代码本身已经证明这个值几乎不变、有本地缓存兜底、缓存/字典都查不到还有硬编码默认值——真正跨服务查一次字典表换来的收益极小。建议 Phase 2 执行 task-rpc 拆分时评估直接把这个阈值做成 `task-rpc` 自己的静态配置项（`etc/task.yaml` 里的一个字段），不做 RPC，省掉一次不必要的服务间依赖；如果确认要保留可配置性，再退回模式①。这个决定留到 `18-service-extraction-runbook.md` 的 task-rpc 附录里由执行者根据当时的实际需求判断，不在本文档里强行定论。

## 6. 非目标

- 不做物理隔离（独立数据库主机/独立 MySQL 实例）——5 个 schema 初期共用一台 MySQL，Phase 2 只做逻辑隔离。
- 不引入分布式事务/Saga——跨服务一致性统一走本文档第 5 节的三种模式，不引入 2PC/TCC 之类的重量级方案。
- 不为了拆分而拆出第 6、7 个数据服务——9 个业务域到此为止收敛为 5 个数据服务 + 1 个网关，没有找到独立扩缩容/独立生命周期/独立信任边界理由的域不再继续细分。
- 不在本轮改变任何表结构本身（字段、索引）——本文档只处理"表归哪个服务"，不重新设计 schema。

## 7. 完成的定义

- `db/services/` 目录树按第 4 节结构建好（含 `iam/metric/` 补充目录），`db/tables.sql`/`db/migrations/*.sql`/`db/demo/*.sql` 里的每张表的建表+初始化 SQL 都能在新目录树里找到对应文件，一一核对不遗漏（38 张表，逐张勾选）。
- `scripts/generate-sql.sh` 的 域→服务 映射表已经加上，新建一个模块跑一次 `-group iam/xxx` 能验证落进 `db/services/iam/xxx/`。
- 第 3 节列出的每一处跨域越界 import，在 `16-rpc-conventions.md`/`17-async-eventing.md` 里都能找到对应的处理方式（RPC/Streams/TaskCallback/进程内合并），没有遗漏项。
- 团队（即用户本人）过一遍第 1 节的 5 条合并/拆分理由，确认认可，再进入 `18-service-extraction-runbook.md` 的实际执行阶段。
