# 归档：admin-server 历史开发记录（2025-01 ~ 2026-07-07，changelog 起点之前）

> **本文件是历史归档，不再维护，不要往这里追加新内容。** 原文件 `docs/后端开发进度.md` 已于 2026-07-17 按「文档分层与生命周期」规则（见 `.cursor/rules/00-workflow.mdc`）退场：2026-07-10（`docs/changelog/` 起点）之前的内容与 changelog 没有重复，整篇原样搬到这里；2026-07-10 之后的内容（Phase 1-3 重构记录）与 `docs/changelog/2026-07-10.md` 起的各篇日期文件完全重复，直接删除未保留。新的后端变更一律写进 `docs/changelog/`，日常"系统现在是什么样"的参考见 `docs/admin-server-维护导航.md`。
>
> 以下正文原样保留（含内部编号、"见第 X 节"式交叉引用），**不做重新整理**；部分引用指向的原文档（`docs/admin-server-ddd-refactor-prompt.md`/`admin-server-ddd-smoke-test.md`/`admin-server-phase0-goctl-spike.md`）已随 DDD-lite 重构专项结束一并删除，相关链接已在下方就地失效处理为纯文本。原文档第 16-22、24 节（2026-07-10 起的 Phase 1-3 记录）与 `docs/changelog/` 重复，已删除不在本文件中；**唯一例外是第 23 节**（一处未被任何 changelog 日期文件覆盖、且仍被多处代码注释引用的 SQL bug 修复记录），原样保留在本文件末尾。

---

## 0. 项目概述与技术栈

基于 go-zero 框架的后端管理系统，分层架构（Handler → Logic → Repository → Model），支持 RBAC 权限管理、系统配置、数据字典、文件存储、日志监控、聊天、博客、视频、SDK 开放接口等功能模块。

**核心技术栈：**
- go-zero 1.9.3+ / Go 1.24+
- MySQL 8.0+ + Redis 6.0+，数据访问层统一使用 go-zero sqlx + cache（goctl model 生成），动态 SQL 使用 squirrel
- JWT（双令牌：Access + Refresh）+ bcrypt
- go-zero logx 日志，gopsutil 系统监控
- gorilla/websocket（聊天/任务通知）

**目录结构（`admin-server/`，本节描述的是归档时点的现状，最新现状见 `docs/admin-server-维护导航.md`）：**
```
api/            # admin.api —— 统一 API 定义（路由真源）
cmd/adminseed/  # 管理员初始化工具，现纯 gRPC 客户端调 IamRPC.UserCreate，不再直连数据库
db/services/    # 按 <service>/<module>/{create_table,init,migrations/} 拆分的建表+初始化 SQL
etc/            # admin-api.yaml、middleware.yaml
internal/
  ├─ config/ consts/ middleware/  # 中间件已全部改调 zrpc client / 共享 Redis 直连，不再直连数据库
  ├─ handler/ logic/        # goctl 生成骨架 + 薄胶水（解析请求 -> 调对应 XxxRPC -> 映射响应）
  ├─ redisconn/              # gateway 唯一保留的存储直连：共享 Redis（token 黑名单、限流滑动窗口）
  ├─ types/ svc/             # svc.ServiceContext 只持有 Redis + 5 个 XxxRPC client 字段
                              # （internal/repository/、internal/model/、internal/domain/ 三个目录已整体删除）
pkg/            # errs/ response/ jwt/ cache/ audit/ monitor/ useragent/ taskcallback/ iamcallback/
scripts/        # generate-sql.sh / generate-model.sh / generate-api.sh / generate-ts.sh / generate-rpc.sh
services/iam/, services/task/, services/sdk/, services/chat/, services/content/
  # Phase 2 拆出的 iam-rpc/task-rpc/sdk-rpc/chat-rpc/content-rpc 五个独立服务，全部领域代码/
  # repository/model/domain 都已搬出单体；iam-rpc 额外承载 RBAC 权限校验 + pkg/taskcallback、
  # pkg/iamcallback 两个跨服务回调契约的服务端实现
```

**中间件体系（执行顺序）**：Performance → RateLimit → Auth → Permission → OperationLog；此外有 ApiEnabledMiddleware（业务开关，与 Permission 互斥）、SDKAuth / SDKRateLimit / SDKCallLog（SDK 对外接口专用）。详细规则见 `10-go-code-style.mdc`。

**前后端协同**：后端定义 `.api` → `goctl api go` 生成后端 → `goctl api ts` 生成前端类型 → 前端二次封装。详细生成流程见 `00-workflow.mdc`。

---

## 1. 开发规范

开发规范已迁移至 `.cursor/rules/00-workflow.mdc` 与 `.cursor/rules/10-go-code-style.mdc`，请查阅规则文件。

---

## 2. 核心功能模块一览

- **认证与授权**：用户名/密码 + bcrypt 比对，JWT 双令牌（Access/Refresh），登出黑名单（Redis `jwt:blacklist:*`），RBAC 权限验证。
- **RBAC 权限管理**：用户/角色/权限/部门/菜单（树形，含按钮级 type=3）/接口 六大对象的 CRUD 及关联关系管理。
- **系统支撑**：系统配置（热更新+缓存）、数据字典（类型/项管理+缓存）、文件管理（上传/下载）。
- **日志与监控**：操作日志、登录日志、审计日志、性能日志、系统监控（CPU/内存/磁盘/网络）。
- **通用打点统计（PV/UV/VV/IP）**：`POST /api/v1/metrics/report` 统一上报入口（`module` 白名单校验），Redis 做实时计数，`metric_daily_stats` 表做日汇总（`MetricRepository.UpsertDailyStats` 增量落库，异步落库必须用独立 `context.Background()` + 超时，避免 HTTP 请求 context 被取消），后台查询走 `GET /api/v1/metrics/stats`。
- **缓存策略**：go-zero 自动缓存单条记录查询；业务层缓存菜单树/权限列表/字典项/配置项，过期时间：用户权限/菜单树 30 分钟，字典 1 小时，配置 10 分钟，变更时主动清除。
- **SDK 管理**：API Key（key+secret）管理、接口管理、key-接口授权与自定义限频、调用记录与导出，专用中间件鉴权/限流/记录。

---

## 3. （已合并至第 1 节）

---

## 4. 已完成功能

- **基础工程**：go-zero Rest 骨架、健康检查、分层结构、统一响应/错误码（`pkg/errs`、`pkg/response`）、配置与 ServiceContext（DB/Redis/JWT/bcrypt）。
- **认证模块**：登录/刷新/登出、JWT 鉴权中间件、`/auth/profile`。
- **RBAC 完整实现**：
  - 角色/权限/部门（树形）/菜单（树形，按钮级权限 type=3）/用户/接口 六个对象的 CRUD 及前端页面。
  - 关联关系管理：用户-角色、角色-权限、权限-菜单、权限-接口。
  - `PermissionMiddleware`：基于用户权限验证 API 访问。
- **系统支撑**：系统配置（CRUD+按 key 查询+刷新缓存）、字典类型/项管理、公共字典查询接口 `/api/v1/dict`、文件管理（上传/下载）、缓存刷新接口 `/api/v1/cache/refresh`。
- **日志与监控**：
  - 操作日志：中间件自动记录所有增删改操作，异步写入，支持分页查询/筛选/导出。
  - 登录日志：记录登录/登出、IP、设备信息，支持统计（总次数/失败次数/今日统计/在线用户数）。
  - 审计日志：记录权限分配、角色变更、配置修改等敏感操作，异步写入。
  - 系统监控：健康检查（`/ping`）、资源使用情况（`/monitor/status`）、系统统计（`/monitor/stats`，基于 gopsutil）。
  - 限流中间件：按 IP/用户/接口限流，429 状态码返回。
- **性能与缓存优化**：
  - 业务层缓存工具 `pkg/cache/business_cache.go`（Get/Set/Delete/GetOrSet，含防雪崩随机过期）。
  - 缓存落地位置：用户权限列表、用户菜单树、完整菜单树、字典项列表、配置项，均在对应 create/update/delete 时清除。
  - 慢查询监控 `pkg/monitor/slow_query.go`；DB 连接池调优（反射配置 MaxOpen/MaxIdle/ConnMaxLifetime）。
- **demo 功能**：开发流程示例（sqlgen → Model → API → 业务实现），供新功能开发参考。
- **聊天模块**：会话/群组/消息 CRUD，WebSocket 实时通信（`ChatWSHandler` + `ChatHub`），JWT 认证、房间管理、消息广播。
- **公告与通知模块**：公告 CRUD（类型/状态/定时发布），消息通知（多来源类型、已读管理）。
- **性能日志模块**：中间件自动记录接口响应时间/慢接口标记，支持查询筛选。
- **视频管理模块**：视频 CRUD、视频代理接口（解决跨域播放，支持 m3u8/mp4），前端智能播放策略（直连失败自动切代理）。详见第 8、10 节。
- **异步任务模块**：任务管理（列表/详情/取消/最近任务）、任务调度器（5 秒轮询+Redis 分布式锁）、Excel/CSV 导出执行器、任务通知（复用通知系统）；操作/登录/审计/SDK调用/性能日志的导出统一改造为异步任务模式。
- **SDK 管理模块**：`sdk_key`/`sdk_interface`/`sdk_key_api`/`sdk_call_log` 数据层；SDKAuth/SDKRateLimit/SDKCallLog 中间件；管理接口（Key/接口/授权/调用记录）；对外文件上传接口 `/sdk/file/upload`。
- **视频与公共接口后端实施**：中间件梳理（ApiEnabledMiddleware 覆盖公共接口组）、公共字典接口 `/api/v1/public/dict`（白名单仅 `video_proxy_url`）、视频采集接口 CORS 支持、M3U8 代理 Range 请求与边下边播（详见第 10 节）。
- **博客功能模块**（含扩展功能）：标签/文章/审核三大子模块，友情链接、社交信息、文章置顶、公共接口补充（作者信息/文章统计/相邻文章）。详见第 12 节。

---

## 5. 待实现 / 待完善功能（归档时点的存量，不代表当前状态）

- **内部代码清理（2026-01-19 审查结论，仍是开放项）**：
  - [x] ~~`performance_log_repository.go`、`chat_repository.go` 尚未改用 squirrel 构建 SQL~~——**订正（2026-07-12，sdk-rpc 拆分会话核实）**：`internal/repository/monitoring/performance_log_repository.go`（DDD-lite 域重组后的真实路径）已经完整迁移到 squirrel，不再是遗留项；`internal/repository/chat/chat_repository.go` 仍有部分方法用参数化的静态多行 SQL（无拼接风险，但不是 squirrel），继续是开放项，`10-go-code-style.mdc` 已同步更正例外清单。**再订正（2026-07-12，chat-rpc 拆分会话）**：chat 域已拆分成独立服务，这个文件的真实路径变成 `services/chat/internal/repository/chat/chat_repository.go`，问题本身（部分方法未用 squirrel）仍未解决，继续是开放项。
  - [ ] GET 详情请求的 form 标签补充 `optional`（audit/login/operation/task 等）
  - [ ] 抽取公共分页默认值 helper，替换日志/视频等列表逻辑的重复代码
  - [ ] 聊天/任务列表的字典配置读取增加缓存或 svcCtx 预加载
  - [ ] 清理 `VideoCollectLogic` 的调试日志与敏感输出
- **性能与缓存优化（阶段六，部分未完成）**：
  - [ ] 缓存预热机制（服务启动时预加载热点数据）
  - [ ] 防止缓存穿透（对不存在的 key 也缓存较短过期时间）
  - [ ] 防止缓存击穿（分布式锁/单机锁保护热点 key）
  - [ ] 数据库索引优化（用户表 username/email、日志表 user_id/created_at 等）
  - [ ] 查询字段优化（避免 `SELECT *`）、关联查询优化（避免 N+1）
  - [ ] 接口缓存（菜单树/字典列表等查询类接口）、批量操作优化、异步处理耗时操作
  - [ ] 接口文档完善（go-zero 自带文档生成、参数说明、在线访问地址、版本管理）
  - [ ] 单元测试与集成测试（Repository/Logic/Handler/中间件覆盖，核心业务测试覆盖率 ≥ 80%）
- **博客模块收尾**：
  - [ ] 12.7.9 关键 Logic 单元测试（状态流转、字典长度校验、标签约束、PV/UV/VV/IP 上报）
  - [ ] 13.6.5 社交信息管理 Logic/Handler 的分页/校验/软删除细节待补
  - [ ] 13.6.8 博客扩展功能的 SQL 执行与联调测试
  - [ ] 第 14 章（作者信息/文章统计/相邻文章接口）：API 定义、Repository 扩展方法（`FindPrevArticle`/`FindNextArticle`/`CountPublishedArticles`）、Logic/Handler 实现、TypeScript 生成均待完成

其余各阶段（一～十一）均已完成，详见第 4 节「已完成功能」。

---

## 6. 技术决策记录

- 2025-12-23：后台按 go-zero + sqlx + cache 重建，不保留旧兼容路径；目录遵循 Handler → Logic → Repository → Model（不使用 Service 层）。
- 2025-12-26：数据访问层统一使用 goctl 生成的 sqlx + cache Model，不再使用 GORM。
- 2025-12-30：Redis 客户端统一使用 go-zero `stores/redis` 组件，不再直接依赖 go-redis/v9；系统级固定枚举/常量统一放入 `internal/consts` 包，禁止硬编码字符串（优先常量与数据字典）。
- 2026-01-13：**表格导出统一使用异步任务机制**：所有导出类接口一律通过 `admin_task` 记录 + `ExcelExportExecutor` 异步生成 CSV，再由文件管理模块创建 `AdminFile` 记录并返回下载 URL；禁止新增同步向 HTTP 响应流写出 CSV 的接口；前端导出按钮只触发任务创建，通过任务列表/浮动任务球查看结果并下载。
- 2026-01-15：**通用打点模块与统计规范**：`module` 参数走白名单校验；实时统计用 Redis，日汇总用 MySQL `metric_daily_stats`（`UpsertDailyStats` 做 `INSERT ... ON DUPLICATE KEY UPDATE`）；异步落库必须用独立 `context.Background()` + 超时，禁止复用 HTTP 请求 context；后台统计查询统一走 `GET /api/v1/metrics/stats`，不直接读 Redis，便于长期留存与离线分析。
- 2026-07-13：**`LAST_INSERT_ID()` 多行插入后的返回值是批次第一行的自增 ID，不是最后一行**（MySQL 官方文档行为，已用临时表实测验证）。写初始化 SQL 时如果一条 `INSERT` 插入多行、后续要按行反推各自的 ID，必须用 `LAST_INSERT_ID() + N`（以第一行为基准正向递增），不能用 `LAST_INSERT_ID() - N`（那是假设返回的是最后一行）。详见第 23 节。

---

## 7. API 清单

接口清单以 `admin-server/api/admin.api` 为准，不在此文档中维护（历史版本记录的路由格式已过时，例如 `:id` 路径参数写法在当前项目中已不再使用）。

---

## 8. 关键代码位置（归档时点的现状，最新现状见 `docs/admin-server-维护导航.md`）

- 入口与路由：`admin-server/admin.go`、`internal/handler/routes.go`
- API 定义：`admin-server/api/admin.api`
- 配置：`internal/config/config.go`、`etc/admin-api.yaml`
- **RBAC（用户/角色/权限/部门/菜单/接口 + 关联关系，2026-07-13 起随 iam+system+monitoring+misc 四个域整体拆分为独立服务 iam-rpc）**：
  - gateway 侧薄胶水：`internal/handler/iam/{user,role,permission,department,menu,api,user_role,role_permission,permission_menu,permission_api,auth}/`（goctl 生成）+ `internal/logic/iam/{user,role,permission,department,menu,api,user_role,role_permission,permission_menu,permission_api,auth}/`（解析 HTTP 请求 → 调 `svcCtx.IamRPC` → 映射响应）
  - iam-rpc 本体：`admin-server/services/iam/`（`internal/domain/iam/{permission_resolver,rbac_service,user_service}.go`、`internal/repository/iam/`、`internal/model/iam/`、`internal/logic/`（约 90 个文件，`Xxx*Logic` 命名，未按域再分子目录）、`internal/server/iamserver.go`）
  - 权限中间件：`internal/middleware/permissionmiddleware.go`（本身仍在 gateway，内部实现已改调 `IamRPC.CheckPermission`，不再直连数据库）；`internal/middleware/authmiddleware.go` 的 JWT 黑名单校验改成直连共享 Redis（`internal/redisconn/`），不走 RPC
- **系统支撑（配置/字典/文件，随 RBAC 一起拆到 iam-rpc）**：
  - gateway 侧薄胶水：`internal/handler/system/{config,dict_type,dict_item,dict,file}/` + `internal/logic/system/{config,dict_type,dict_item,dict,file}/`（`file` 例外：文件字节仍由 gateway 直接读写共享 `uploads` 卷，只有元数据登记/查询走 `IamRPC.FileRegister`/`FileGetMeta`，见 `internal/logic/system/file/file_upload_logic.go`/`file_download_logic.go`）
  - iam-rpc 本体：`services/iam/internal/repository/system/`、`services/iam/internal/model/system/`
- **日志与监控（随 RBAC 一起拆到 iam-rpc）**：
  - gateway 侧薄胶水：`internal/handler/monitoring/{operation_log,login_log,monitor,audit_log}/` + `internal/logic/monitoring/{operation_log,login_log,monitor,audit_log}/`
  - iam-rpc 本体：`services/iam/internal/repository/monitoring/`、`services/iam/internal/model/monitoring/`
  - 中间件：`internal/middleware/operationlogmiddleware.go`（改批量调 `IamRPC.BatchRecordOperationLog`）、`ratelimitmiddleware.go`（改直连共享 Redis）；审计工具：`pkg/audit/audit.go`（改调 `IamCallbackRPC.RecordAuditLog`）
- **性能与缓存优化**：
  - 性能监控：`pkg/monitor/slow_query.go`、`pkg/monitor/performance.go`、`internal/middleware/performancemiddleware.go`（改异步调 `IamRPC.RecordPerformanceLog`）
  - 原 `internal/repository/sql_conn.go`（反射设置 MaxOpen/MaxIdle）随 gateway `internal/repository/` 整体删除，未在 iam-rpc 侧找到等价文件，连接池调优现状待确认（不确定是否需要在 `services/iam/internal/repository/` 补一份，还是现在依赖 go-zero sqlx 默认值，标记为待核实项，不是本次编造）
  - 业务层缓存工具 `pkg/cache/business_cache.go` 使用位置随对应 logic 一起搬到了 iam-rpc：`services/iam/internal/logic/profilelogic.go`（权限）、`menumytreelogic.go`/`menutreelogic.go`（菜单树）、`dictgetlogic.go`（字典）、`configgetlogic.go`（配置）
- **demo 示例**（开发流程参考，随 misc 域拆到 iam-rpc）：gateway 侧 `internal/handler/misc/demo/` + `internal/logic/misc/demo/`；iam-rpc 本体 `services/iam/internal/repository/misc/demo_repository.go`、`services/iam/internal/model/misc/demomodel*.go`；`db/services/iam/demo/create_table_demo.sql`、`init_demo.sql`
- **聊天模块（2026-07-12 起拆分为独立服务 chat-rpc）**：
  - gateway 侧薄胶水：`admin-server/internal/logic/chat/{chat,group,message}/`（解析 HTTP 请求→调 ChatRPC→映射响应）；WS↔gRPC 桥接：`admin-server/internal/handler/chat/chatwshandler.go`
  - chat-rpc 本体：`admin-server/services/chat/`（`internal/domain/chat/onboarding.go`、`internal/repository/chat/`（`chat_repository.go` 部分方法未使用 squirrel，见第 5 节）、`internal/hub/chathub.go`（连接表）、`internal/consumer/chat_user_created_consumer.go`（消费 `stream:chat.user.created`）、`internal/model/chat/`、`internal/logic/`）
  - 跨服务契约：`admin-server/pkg/iamcallback/`（`FindActiveUserChunk`/`GetUserProfile`，服务端实现现为 `services/iam/internal/server/iamcallbackserver.go`；原单体内嵌在 `internal/rpcserver/iamcallback/`，2026-07-13 起随 iam-rpc 拆分搬迁）
- **公告与通知（随 system 域一起拆到 iam-rpc）**：gateway 侧薄胶水 `internal/handler/system/{notice,notification}/` + `internal/logic/system/{notice,notification}/`；iam-rpc 本体 `services/iam/internal/repository/system/{notice,notification}_repository.go`
- **性能日志（随 monitoring 域一起拆到 iam-rpc）**：gateway 侧薄胶水 `internal/handler/monitoring/performance_log/` + `internal/logic/monitoring/performance_log/`；iam-rpc 本体 `services/iam/internal/repository/monitoring/performance_log_repository.go`（已使用 squirrel）
- **视频管理（2026-07-12 起 CRUD/采集/公开列表拆分为独立服务 content-rpc；M3U8 代理与 CORS 预检不涉及域数据，继续留在 gateway）**：
  - gateway 侧薄胶水：`admin-server/internal/logic/video/{video,video_collect,public}/`（`M3u8ProxyLogic`、`VideoCollectOptionsLogic` 例外，原样留在 gateway 不变）
  - content-rpc 本体：`admin-server/services/content/`（`internal/repository/video/video_repository.go`、`internal/model/video/`）
- **异步任务（2026-07-11 起拆分为独立服务 task-rpc）**：
  - gateway 侧薄胶水：`admin-server/internal/logic/task/{task,public}/`（解析 HTTP 请求→调 TaskRPC→映射响应，不再持有业务逻辑）
  - task-rpc 本体：`admin-server/services/task/`（`internal/domain/task/{scheduler.go,notifier.go,executors/}`、`internal/repository/task_repository.go`、`internal/model/task/`、`internal/logic/`、`internal/consts/`）
  - 跨服务契约：`admin-server/pkg/taskcallback/`（`FetchExportData`/`RegisterExportFile`，服务端实现现为 `services/iam/internal/server/taskcallbackserver.go`；原单体内嵌在 `internal/rpcserver/taskcallback/`）、`services/iam/internal/consumer/task_notification_consumer.go`（消费 `stream:task.notification`；原 `admin-server/internal/consumer/`，均于 2026-07-13 随 iam-rpc 拆分搬迁）
- **SDK 管理**（Phase 2 已拆分成独立服务 sdk-rpc，gateway 侧只剩薄胶水）：
  - gateway 侧：Handler：`internal/handler/sdk/`（goctl 生成）；Logic：`internal/logic/sdk/sdk/`（薄胶水，调 `svcCtx.SdkRPC`）、`internal/logic/sdk/public/sdk_file_upload_logic.go`（委托 `internal/logic/system/file`，不经过 sdk-rpc）；中间件：`internal/middleware/sdkauthmiddleware.go`、`sdkratelimitmiddleware.go`、`sdkcalllogmiddleware.go`（本身仍在 gateway，内部实现调 `svcCtx.SdkRPC` 的 `VerifyApiKey`/`GetEffectiveRateLimit`/`RecordCallLog`）
  - sdk-rpc 本体：`admin-server/services/sdk/`（`internal/domain/sdk/sdk_service.go`、`internal/repository/{store.go,sdk/}`、`internal/model/sdk/`、`internal/logic/`、`internal/consts/`）
  - 跨服务契约：sdk-rpc 新增的 `SdkCallLogExport` 方法（`services/sdk/rpc/sdk.proto`）供 `services/iam/internal/server/taskcallbackserver.go`（原单体内嵌的 `internal/rpcserver/taskcallback/server.go`，2026-07-13 起随 iam-rpc 拆分搬迁）回调取 SDK 调用日志导出数据
- **M3U8 代理与公共接口**：`internal/handler/video/m3u8/m3u8_proxy_handler.go`、`internal/logic/video/m3u8/m3u8_proxy_logic.go`（不涉及域数据，content-rpc 拆分时原样留在 gateway）、`internal/handler/misc/public/public_dict_get_handler.go`（公共字典查询，调 `IamRPC.DictGet`）、`internal/handler/video/video_collect/`；Nginx 配置：`config/nginxconfig.txt`
- **博客模块（含扩展，2026-07-12 起拆分为独立服务 content-rpc）**：
  - gateway 侧薄胶水：`admin-server/internal/logic/blog/{tag,article,article_audit,friend_link,social_info,public}/`（解析 HTTP 请求→调 ContentRPC→映射响应）
  - content-rpc 本体：`admin-server/services/content/`（`internal/domain/content/blog_service.go`、`internal/repository/blog/`、`internal/model/blog/`、`internal/consts/`、`internal/logic/`）
  - 跨服务契约：`admin-server/pkg/iamcallback/`（`GetUserProfile` 新增 `signature` 字段供 `PublicBlogAuthorInfo` 用、新增 `RecordAuditLog` 供 `BlogArticleAudit`/`BlogArticleAuditUnpublish` 写审计日志）
- **打点统计（随 monitoring 域一起拆到 iam-rpc）**：gateway 侧薄胶水 `internal/logic/monitoring/metric/metric_report_logic.go`、`internal/logic/monitoring/metric_admin/metric_stats_logic.go`；iam-rpc 本体 `services/iam/internal/repository/monitoring/metric_repository.go`

---

## 9. 数据库变更记录

- **2025-12-24**：统一初始化 SQL `db/tables.sql`，包含 `admin_user`/`admin_role`/`admin_permission`/`admin_department`/`admin_user_role`/`admin_role_permission`/`admin_menu`/`admin_api`/`admin_permission_menu`/`admin_permission_api`；所有表统一 `created_at`/`updated_at`/`deleted_at` BIGINT 秒级字段，软删除通过 `deleted_at` 标识（关联关系表不含 `deleted_at`）；菜单/用户管理菜单项与权限 SQL 见 `db/permissions_menu_user.sql`。
- **2025-12-26**：按钮级菜单（type=3）及权限-菜单关联补充（共 18 个按钮菜单）；阶段四新增表 `admin_config`（系统配置）、`admin_dict_type`/`admin_dict_item`（数据字典）、`admin_file`（文件）及对应权限/菜单/接口初始化数据；补充通用权限 `common:profile`/`common:logout`/`common:dict`/`common:cache_refresh`/`menu:my_tree`；`MenuMyTreeLogic` 菜单过滤逻辑修复；demo 功能示例增量 SQL（`create_table_demo.sql`/`init_demo.sql`）。
- **阶段七～九**（聊天/公告通知/性能日志）：新增表 `admin_chat`/`admin_chat_group`/`admin_chat_message`/`admin_chat_group_member`（聊天）、`admin_notice`/`admin_notification`（公告与通知）、`admin_performance_log`（性能日志）。
- **2025-01-01（阶段十，视频模块）**：新增表 `video`（id/name/cover/duration/play_url/description + 时间戳字段），字典 `video_proxy_url`（视频代理地址配置）；增量 SQL：`db/migrations/dict_video_20250101.sql`、`create_table_video.sql`、`init_video.sql`、`init_video_parent_menu.sql`；后续 M3U8 代理重构支持 Range 请求、流式传输与连接池优化。
- **2026-01-13（阶段十一，异步任务模块）**：新增表 `admin_task`（id/name/type/execution_type/status/params/result/error_message/user_id/scheduled_at/started_at/finished_at + 时间戳字段），字典 `task_type`/`task_execution_type`/`task_status`/`task_config`，扩展 `notification_source_type` 增加 "task" 选项；增量 SQL：`db/migrations/dict_task_20250115.sql`、`create_table_task.sql`、`init_task.sql`；操作/登录/审计/SDK调用/性能日志导出统一改造为异步任务模式。
- **2026-01-15（通用打点模块）**：新增表 `metric_daily_stats`（module/biz_id/day/pv/uv/vv/ip + 时间戳字段），对 `module+biz_id+day` 建唯一索引 `uk_metric_stats_module_biz_day`；增量 SQL：`db/migrations/create_table_metric.sql`（表结构）、`init_metric_stats.sql`（数据统计菜单、`metric:stats` 权限、接口映射）；Redis 仍作为实时计数，该表用于日级汇总与后台查询。
- **博客模块（阶段十二）**：新增表 `blog_tag`/`blog_article`/`blog_article_tag`/`blog_article_audit`，及状态/审核状态/长度限制等字典与权限初始化；增量 SQL：`db/migrations/create_table_blog.sql`、`dict_blog_*.sql`、`init_blog.sql`。
- **博客模块扩展（阶段十三）**：新增表 `blog_friend_link`（友情链接）、`blog_social_info`（社交信息），`blog_article` 表新增 `is_top` 置顶字段；增量 SQL：`db/migrations/create_table_blog_extension.sql`、`dict_blog_extension_20260116.sql`、`init_blog_extension.sql`。

---

## 10. HTTP Range 请求与边下边播（M3U8 代理）

**功能**：M3U8/视频代理接口透传客户端 `Range` 请求头到源服务器，支持 `206 Partial Content` 响应的流式转发，实现视频边下边播（低延迟、低内存占用、支持拖拽播放），同时解决跨域播放限制。

**关键设计决策**：
- 代理透传 `Range` 请求头与 `Content-Range`/`Accept-Ranges`/`Content-Length` 响应头，使用 `io.Copy` 流式转发而非缓存整个文件；HTTP 客户端使用连接池（`MaxIdleConnsPerHost` 等）提升并发性能。
- **CORS 统一由 Nginx 处理**：`ApiEnabledMiddleware` 返回错误时 Handler 不会执行，导致后端设置的 CORS 头失效；因此最终方案是在 `/gateway/` location 用 `add_header ... always` 兜底设置 CORS 头（包括错误响应），OPTIONS 预检直接由 Nginx 处理返回 204，不转发到后端；后端 Handler/Logic 不再处理 CORS，逻辑简化。
- 评估过 HTTP/2 Server Push、HTTP/3、WebSocket 流式传输、CDN、智能预加载、ABR、P2P 等更高级方案，结论：当前规模下 HTTP Range 已足够（简单、低成本、效果好），CDN+Range 是未来扩展方向，其余方案暂不需要。

**代码位置**：`internal/handler/m3u8/`、`internal/logic/m3u8/`（Handler/Logic 均已简化，不处理 CORS）、Nginx 配置 `config/nginxconfig.txt`（`/gateway/` location）。

---

## 11. 视频与公共接口实施小结

已完成内容（与第 4、10 节部分重复，此处仅作索引）：
- 中间件梳理：公共接口组统一接入 `ApiEnabledMiddleware`（业务启停控制）。
- 公共字典接口 `GET /api/v1/public/dict`：白名单仅暴露 `video_proxy_url`，供未登录页面获取视频代理地址。代码：`internal/handler/public/publicdictgethandler.go`、`internal/logic/public/publicdictgetlogic.go`。
- 视频采集接口 CORS 支持（`/api/v1/videos/collect` + OPTIONS 预检），允许第三方页面调用。代码：`internal/handler/video_collect/`。
- M3U8 代理接口详见第 10 节。

---

## 12. 博客功能模块

**模块构成**：标签管理（`blog_tag`）、文章管理（`blog_article`，状态机：草稿→待审核→审核通过→上架/下架）、文章审核（`blog_article_audit`，审核记录+审计日志联动）、公共文章接口（仅展示审核通过+上架的文章）。

**关键设计**：
- 所有业务可变枚举/长度限制走数据字典（如 `blog_article_status`、`blog_article_title_max_length`、`blog_article_summary_length`），枚举值从 1 开始；状态流转在 Logic 层严格校验，禁止越级跳转。
- 编辑权限：上架状态文章不可编辑（`403`）；其余状态编辑后自动流转回相应审核状态。
- 封面为空时返回 `coverFallback`（标题首字），由前端生成占位封面。
- Markdown 内容原样存储与传输，不做 HTML 转义（避免破坏原文），由前端渲染时负责安全过滤。
- 打点统计（PV/UV/VV/IP）复用通用 `metric` 模块，见第 2、6 节。

**代码位置**：
- Handler/Logic：`internal/{handler,logic}/blog_tag/`、`blog_article/`、`blog_article_audit/`、`public_blog/`
- Repository：`internal/repository/blog_tag_repository.go`、`blog_article_repository.go`、`blog_article_tag_repository.go`、`blog_article_audit_repository.go`
- 常量：`internal/consts/blog.go`；长度校验工具：`pkg/dict/dict.go`
- SQL：`db/migrations/create_table_blog.sql`、`dict_blog_*.sql`、`init_blog.sql`

单元测试与接口联调（12.7.9）尚未完成，见第 5 节。

> 以上「代码位置」是拆分前（Phase 2 之前）的原始单体实现路径。2026-07-12 起博客域整体拆分为独立服务 content-rpc，现状见第 8 节「博客模块」条目。

---

## 13. 博客模块扩展（友情链接、社交信息、文章置顶）

**模块构成**：友情链接管理（`blog_friend_link`）、社交信息管理（`blog_social_info`，均含后台 CRUD + 公共列表接口）、文章置顶（`blog_article.is_top` 字段，置顶数量走字典 `blog_article_top_max_count` 限制，默认 1 篇，超限时自动取消最早置顶）。

**代码位置**：
- Handler/Logic：`internal/{handler,logic}/blog_friend_link/`、`blog_social_info/`
- Repository：`internal/repository/blog_friend_link_repository.go`、`blog_social_info_repository.go`、`blog_article_repository.go`（扩展 `UpdateTopStatus`/`FindPublicPage`）
- 文章置顶逻辑：`internal/logic/blog_article/blogarticletoplogic.go`、`blogarticleuntoplogic.go`
- 公共接口扩展：`internal/logic/public_blog/`
- 常量：`internal/consts/blog.go`（扩展）
- SQL：`db/migrations/create_table_blog_extension.sql`、`dict_blog_extension_20260116.sql`、`init_blog_extension.sql`

社交信息管理 Logic/Handler 细节（13.6.5）与联调测试（13.6.8）尚未完全完成，见第 5 节。

> 以上「代码位置」是拆分前（Phase 2 之前）的原始单体实现路径，现状见第 8 节「博客模块」条目。

---

## 14. 博客页面改造后端接口补充（待实现，归档时已完成，见下方订正）

配合前端博客页面改造（详见 `archive-frontend.md`），需在 `public_blog` 分组下补充三个只读接口：
- `GET /api/v1/public/blog/author-info`：超级管理员（`admin_user.id=1`）公开信息（昵称/头像/签名）。
- `GET /api/v1/public/blog/article-stats`：已发布文章总数统计。
- `GET /api/v1/public/blog/articles/prev` / `.../next`：按 `publish_time` 查询相邻已发布文章，供详情页翻页导航。

以上接口的 API 定义、Repository 扩展（`FindPrevArticle`/`FindNextArticle`/`CountPublishedArticles`）、Logic/Handler 实现均尚未完成，落地位置计划为 `internal/logic/blog/public/` 下对应 logic 文件，详见第 5 节待办。

> **本节描述的三个接口已实现**（作者信息/文章统计/上一篇下一篇），随博客域一起拆分进了 content-rpc，2026-07-12 起真实环境验证通过；现状位置为 `services/content/internal/logic/`（`PublicBlogAuthorInfoLogic` 等）与 gateway 薄胶水 `internal/logic/blog/public/`，本节按"只追加"原则保留原始待办记录不重写。

---

## 15. DDD-lite 领域分层重构（2026-07-07，已完成）

**做了什么**：在单体进程内按 9 个业务域（iam/blog/video/chat/sdk/task/monitoring/system/misc）重新组织 handler/logic/repository/model 目录；提取 `internal/domain/iam/permission_resolver.go`（RBAC）、迁移 `internal/domain/task/`（调度器）；精简 `ServiceContext`（删除 7 个具名 repository 字段，统一内联构造）。

**为什么**：个人维护 5 万行代码需要领域边界；找代码、改权限、加模块时有明确路径。

**关键决策**：
- goctl 嵌套 group `<domain>/<module>` 经验证可行（Spike 结论：可行，验证方法/过程原记录于已删除的 `docs/admin-server-phase0-goctl-spike.md`，本节保留结论）
- repository 包与 model 包同名时用 `xxxrepo` 别名（如 `iamrepo`）避免冲突
- 仅 IAM/Task 引入 `internal/domain/`，其余域保持 repository+logic

**维护入口**：[`docs/admin-server-维护导航.md`](../admin-server-维护导航.md)

**冒烟测试与前端联调**：原记录于已删除的 `docs/admin-server-ddd-smoke-test.md`（内容是构建/启动/联调步骤清单，随重构完成失去时效性，未保留）。

---

## 23. 本机新装 MySQL 全新初始化时发现并修复两处 SQL bug（2026-07-13）

> **例外说明**：本节原属于原文档第 16 节起的 Phase 1-3 记录（本文件其余部分已略去，见上方说明），按理不该出现在这里；但复核时发现这处 SQL bug 修复**没有被任何 `docs/changelog/` 日期文件覆盖**，而 `admin-server/services/iam/internal/repository/iam/department_repository.go`、`db/services/content/blog/init_blog.sql`、`db/services/content/blog_extension/init_blog_extension.sql`、`db/services/iam/department/migrations/add_feishu_pending_department_20260716.sql` 四处代码注释仍指向"docs/后端开发进度.md 第 23 节"，所以整节原样保留在此，编号维持原文档的"23"不变，代码注释里的路径改指向本文件即可，节号不用改。

**背景**：本机通过 Homebrew 新装 `mysql@8.0` 用于本地开发，第一次在完全空库上从头跑一遍 `db/services/init-dev-db.sh`（此前所有 dev/CI 库都是逐步演进积累出来的，从未有人真正从零重放过全部初始化 SQL），暴露了两处此前从未触发过的 SQL bug。

**1. `init_blog.sql` / `init_blog_extension.sql`：`LAST_INSERT_ID()` 用反了，多行批量插入的场景全部错位**

`LAST_INSERT_ID()` 在一条多行 `INSERT ... VALUES (...),(...),(...)` 语句后返回的是**第一行**的自增 ID（MySQL 官方行为，不是最后一行），但这两个文件里所有多行批量插入后续都用 `LAST_INSERT_ID() - N` 反推每一行 ID，等于假设返回的是最后一行。在全新空库上跑 content/blog 模块初始化时，这个错位刚好让某一行算出的 `permission_id` 撞上了 iam 域已存在的 `daily_short_sentence:update` 权限（真实 ID 恰好等于错误算出的值），触发 `admin_permission_menu.uk_admin_permission_menu` 唯一键冲突而报错中止；如果没有恰好撞上已存在的键，这类错位会**静默产生错误的权限-菜单/权限-接口关联**，不报错但数据是错的，是比"报错中止"更值得警惕的风险点。

修复：把两个文件里全部 `LAST_INSERT_ID() - N` 改成 `LAST_INSERT_ID() + N`（以批次第一行为基准正向递增）。已在全新库上完整重跑验证，`blog_article:list` 等相关权限-菜单/权限-接口关联人工核对无误。

**2. `db/services/task/task/migrations/dict_task_20250115.sql`：多余逗号导致语法错误**

单行 `INSERT ... VALUES (...)` 后多写了一个逗号，紧跟 `ON DUPLICATE KEY UPDATE`，导致 `ERROR 1064` 语法错误。已删除多余逗号。

**排查方式**：先用临时表验证了 `LAST_INSERT_ID()` 多行插入后的真实返回值（确认返回第一行 ID），排除环境问题，定位为系统性写法错误；再全项目搜索 `LAST_INSERT_ID()` 用法逐一核对：
- `db/services/` 下用到 `LAST_INSERT_ID()` 的文件中，只有 `init_blog.sql`/`init_blog_extension.sql` 是"多行插入 + 反推"模式，其余（`init_task.sql`/`init_demo.sql`/`init_metric.sql`/`init_metric_stats.sql` 等）都是单行插入紧跟 `SET`，不受影响
- `scripts/sqlgen/templates/init_module.sql.tpl`（`generate-sql.sh` 脚手架模板）本身只生成单行插入模式，不会把这个 bug 带进新模块，**确认不需要同步修改**

**影响范围确认**：仅 `db/services/content/blog/init_blog.sql`、`db/services/content/blog_extension/init_blog_extension.sql`、`db/services/task/task/migrations/dict_task_20250115.sql` 三个文件；`scripts/sqlgen/templates/*.tpl` 无需修改。

**遗留提醒**：如果生产/团队共享的 dev 库是逐步演进出来的（不是从这次修复后的脚本重新跑出来的），这两个文件里原来错位关联的权限-菜单/权限-接口数据可能已经以错误状态存在于那些库里，需要人工核对（或等下次该库整体重建时自然修复）；本次修复只改了脚本本身，不含针对已有存量库的迁移脚本。

**3. 提交前 pre-commit 代码审查暴露的第三个问题：`init_blog.sql`/`init_blog_extension.sql` 全文缺少幂等性，且发现 `admin_menu` 表没有唯一键导致项目里"约定俗成"的 `ON DUPLICATE KEY UPDATE` 写法对它实际不生效**

修完上面两处 bug 后提交时，仓库的 pre-commit AI 代码审查（Gentleman Guardian Angel）拦下了这次提交：`init_blog.sql`/`init_blog_extension.sql` 全文所有 `admin_menu`/`admin_permission`/`admin_api` 的 `INSERT` 都是裸插入，没有任何幂等保护，重复执行会产生重复数据——这是这两个文件从最初写的时候就有的问题，跟 `LAST_INSERT_ID()` bug 无关。

参照同仓库 `init_task.sql`/`init_video.sql`/`init_chat.sql` 的写法补幂等性时，**用真实数据库实测发现这三个"参照对象"本身也有问题**：它们对 `admin_menu` 的插入同样套用了 `INSERT ... ON DUPLICATE KEY UPDATE`，但 `admin_menu` 表（`db/services/iam/menu/create_table_menu.sql`）从建表起就只有主键 `id`（自增），没有任何唯一键——`ON DUPLICATE KEY UPDATE` 依赖唯一键冲突才会触发，对没有唯一键的表等于什么都没做，纯粹是"看起来幂等"的装饰写法。实测验证：把 `init_blog.sql`/`init_blog_extension.sql` 按这个"约定写法"改完后，在本机 dev 库上连续跑两次，`admin_menu` 行数从 123 涨到 173（每跑一次多插入一遍全部菜单/按钮行），证实了这个猜测。

**最终修复**：`admin_menu` 类插入全部改为 `INSERT ... SELECT ... FROM DUAL WHERE NOT EXISTS (...)`，按 `path`（有路由的菜单）或 `parent_id + name`（无路由的按钮）判重，这是唯一能在没有唯一键的表上做到真正幂等的写法；`admin_permission`（`code` 唯一键）、`admin_api`（`method+path` 唯一键）、`admin_permission_menu`/`admin_permission_api`（联合唯一键）继续用 `ON DUPLICATE KEY UPDATE`，这几个是真的有效。已在本机 dev 库连续重跑 3 次验证：`admin_menu`/`admin_permission`/`admin_api`/`admin_permission_menu`/`admin_permission_api` 五张表行数三次都完全一致，且逐条核对了全部 `blog_*`/`blog_article:*` 权限-菜单关联，ID 指向均正确。

**遗留提醒（新增）**：`init_task.sql`/`init_video.sql`/`init_chat.sql` 里同样存在这个"`admin_menu` 用 `ON DUPLICATE KEY UPDATE` 但实际不生效"的问题，本次未修改这三个文件（超出本次任务范围，且这三个模块的菜单结构此前从未被验证过重复执行的实际效果），如果后续这些模块的库需要重新初始化或有人真的重复跑过这些脚本，需要单独核查是否已产生重复菜单行。

---

> 本文档到此结束（对应原文档第 15、23 节，中间第 16-22、24 节与 `docs/changelog/` 重复已删除）。2026-07-10 起的后续内容（Phase 1-3 重构记录）见 `docs/changelog/2026-07-10.md` 起的各篇日期文件，索引见 `docs/changelog/README.md`。
