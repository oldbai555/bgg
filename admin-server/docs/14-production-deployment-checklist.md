# 14 · 生产部署清单（持续追加，不是一次性文档）

## 这篇文档的性质

**这是一份随 Phase 1-3 推进持续追加的活文档，不是写完就结束的静态清单。** 按 `10-dev-execution-and-review-points.md` 的口径，开发阶段的本地/开发环境操作（SQL 执行、脚本运行、`docker compose up`）AI 可以直接做，但**生产环境部署动作本身、生产配置文件、生产密钥**永远是"必须停下来问用户"的事项——AI 不会替用户执行这里的任何一条，这份文档的作用是保证等用户真的要上线/发布时，有一份完整、按时间顺序整理好的操作清单可以照着做，不用临时回忆"这次改动到底触发了哪些部署侧的动作"。

**给后续执行者（AI 或人）的强制要求**：Phase 1-3 期间，任何一次改动如果满足"这个改动会影响生产环境的配置/密钥/部署步骤/需要执行的 SQL"，**必须在完成对应代码改动的同一个工作批次里，在本文档追加一条**，不能拖到"回头补"——拖到最后补大概率会漏项。追加时不要删除或改写已有条目的历史记录，除非该条目本身描述错误需要订正（订正也要注明"此前描述有误，已更正"，不要静默改写造成误解）。

## 条目格式（新增条目必须遵循这个结构）

```
### <序号> · <简短标题>

**触发条件**：什么改动/什么时间点会用到这一条（对应哪个 Phase/Week/PR）。

**部署时要做什么**：具体命令 / SQL / 配置文件改动，直接可执行，不要写"参照 xxx 自行判断"这种模糊描述。

**如何验证生效**：执行完之后怎么确认这一步真的生效了（检查什么日志、请求什么接口、查什么表）。

**状态**：`TBD`（占位，等对应改动真正落地时补全）/ `已就绪，待执行`（改动已完成，清单已写好，等用户执行部署）/ `已执行`（用户已在生产环境执行过，附执行日期）。
```

---

## Phase 1 已知条目

### 1 · JWT 密钥改用环境变量注入

**触发条件**：A.5（密钥管理）——`etc/admin-api.yaml` 里硬编码的 JWT 占位符（当前 `AccessSecret: "replace-with-secure-access-secret"`、`RefreshSecret: "replace-with-secure-refresh-secret"`，已核查确认是当前仓库真实内容）会改成 `${JWT_ACCESS_SECRET}`/`${JWT_REFRESH_SECRET}` 环境变量占位符，`admin.go` 里 `conf.MustLoad` 会加 `conf.UseEnv()` 做展开。

**部署时要做什么**：
1. 在生产环境的进程管理配置里（当前是 Supervisor，`script/admin.sh` 管理；Phase 3 之后可能切到 docker-compose，见条目 4）新增两个环境变量：`JWT_ACCESS_SECRET`、`JWT_REFRESH_SECRET`。
2. 两个值必须是用户自己生成的高强度随机字符串（例如 `openssl rand -base64 48`），**不能沿用本地开发 `docker-compose.yml` 里的默认值**（`09-ci-cd-and-deployability.md` 里那两个 `local-dev-*-not-for-prod` 占位值只用于本地）。
3. Supervisor 配置文件（具体路径待用户确认，属于"触及生产配置"，AI 不直接改）需要在 `environment=` 字段里加上这两个变量，或者通过 `/etc/work/` 下的独立环境变量文件加载（具体机制以用户现有部署方式为准，AI 不擅自决定新增文件位置）。
4. 确认 `etc/admin-api.yaml` 生产环境实际部署的那一份也已经同步改成 `${JWT_ACCESS_SECRET}`/`${JWT_REFRESH_SECRET}` 占位符（不是只改了仓库里的模板，生产服务器上实际生效的配置文件也要改）。

**如何验证生效**：
- 启动日志里不应出现 A.5 提到的"空值 fail-fast 检查"报错（说明环境变量已正确注入且非空）。
- 用一个测试账号走一次登录 → 拿到 token → 用 token 访问一个受保护接口，确认 token 签发/校验正常（间接验证密钥被正确读取且前后端一致，不是"改了但没生效，实际还在用默认值"）。

**状态**：`已就绪，待执行`（A.5 代码改动已落地：`etc/admin-api.yaml` 的 `JWT.AccessSecret`/`RefreshSecret` 已改成 `${JWT_ACCESS_SECRET}`/`${JWT_REFRESH_SECRET}`，`admin.go` 的 `conf.MustLoad` 已加 `conf.UseEnv()` 并在其后新增空值 fail-fast 检查；本地开发环境变量由用户自行在 shell/`.env` 里设置任意非空值即可跑通，生产环境的真实高强度密钥值仍需用户按下方步骤在部署环境里生成和设置，AI 不代为生成）。

### 2 · MySQL/Redis 配置的环境变量化现状确认

**触发条件**：核查发现 `admin.go` 已经支持 `-mysql-config`/`-redis-config` 命令行参数指定 `/etc/work/mysql.json`/`/etc/work/redis.json` 路径，说明 MySQL/Redis 连接信息**已经**走外部文件、不提交仓库——这部分不是本轮新增的改动，A.5 只处理 JWT 密钥这一项。

**部署时要做什么**：无新增动作，仅在此记录现状，避免后续误以为 MySQL/Redis 配置也需要迁移。如果 Phase 1 执行过程中发现 `admin.go` 现有的 `/etc/work/mysql.json`/`/etc/work/redis.json` 机制需要调整（比如要支持 `conf.UseEnv()` 统一），再单独补一条新条目，不要复用这一条。

**如何验证生效**：不适用（无动作）。

**状态**：`已就绪`（现状确认，非待办）。

### 3 · IAM 事务化改造涉及的 RBAC 种子数据（TBD，Phase 1 Week 2 补全）

**触发条件**：`04-domain-iam-chat.md` 任务 1（`user_create_logic.go` 修复）+ `01-architecture-target.md` A.1（`Repository.Transact`）落地后，"新用户加入默认群 + 批量私聊初始化"这部分逻辑会从同步改成异步（通过 `internal/domain/task` 现有调度器派发）。如果这个改造过程中涉及新增/调整任何 RBAC 权限种子数据（比如新的任务类型需要对应的操作权限，或者默认群组 `chat_id=1` 这类硬编码 ID 的初始化数据在生产环境是否已存在需要核实），这里要记一条。

**部署时要做什么**：**TBD——待 Phase 1 Week 2 实际执行 IAM+Chat 域改造时，根据当时产生的真实 SQL 差异填写这里的具体命令**，本篇现在不预先编造内容。填写时要包含：具体的 `INSERT`/`UPDATE` SQL 语句（或指向 `db/services/<service>/<module>/migrations/` 下具体文件名）、执行顺序（是否依赖条目 1 的密钥改动先生效）。

**如何验证生效**：TBD，随上面一起补全。

**状态**：`TBD`。

### 4 · Dockerfile / docker-compose 引入后的部署方式变化（如果 Phase 1 期间就切换）

**触发条件**：`09-ci-cd-and-deployability.md` 设计的 Dockerfile/docker-compose 是 Phase 1 的产出，但**是否在 Phase 1 就把生产部署从当前的 `script/admin.sh` 构建→scp→`supervisorctl restart` 流程切到 docker-compose，还是继续用 Supervisor、只把 Docker 化留到 Phase 3 统一切换**，属于 `10-dev-execution-and-review-points.md` 第 2 节"生产环境部署动作本身"的必停项，需要用户拍板。

**部署时要做什么**：TBD，取决于用户的决定。如果决定 Phase 1 就切：需要在此追加完整的"Supervisor → docker-compose 迁移步骤"（镜像构建、`.env` 生产值配置、`docker compose up -d`、旧 Supervisor 进程停止确认、回滚方案）。如果决定留到 Phase 3：这条先标记为"留待 Phase 3"，不需要现在填内容。

**如何验证生效**：TBD。

**状态**：`TBD`（待用户决定切换时机）。

---

### 5 · task-rpc 拆分（Phase 2 第一个落地的服务）

**触发条件**：`services/task/` 完整拆分完成（`docs/progress.md` 对应条目），task-rpc 成为独立部署单元。

**部署时要做什么**：
1. 建 `admin_task` schema（逻辑隔离，物理上和主库同一个 MySQL 实例）：`mysql -uroot -p<pass> -e "CREATE DATABASE IF NOT EXISTS admin_task CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"`，然后跑 `db/services/task/task/create_table_task.sql`（**只需要这一个文件**——`init_task.sql`/`dict_task_20250115.sql` 写的是 `admin_menu`/`admin_permission`/`admin_api`，物理属于主库 iam 部分，已经在主库初始化时跑过，不属于 `admin_task` schema）。docker-compose 场景下这一步已经自动化，见 `db/services/init-task-db.sh` + `docker-compose.yml` 的 mysql 服务第二个 init 脚本挂载。
2. 生产环境实际地址配置（三处，均通过环境变量，`conf.UseEnv()` 已接入，fail-fast 不给静默默认值，和 `JWT_ACCESS_SECRET` 同一套约定）：
   - `TASK_RPC_ENDPOINT`（gateway 侧 `etc/admin-api.yaml` 的 `TaskRpc.Endpoints`，指向 task-rpc 实际监听地址）
   - `TASK_CALLBACK_ENDPOINT`（task-rpc 侧 `services/task/etc/task.yaml` 的 `TaskCallbackRpc.Endpoints`，指向承载 `pkg/taskcallback` server 的进程——当前阶段是单体 gateway 自己）
   - `TASK_MYSQL_DSN`、`TASK_REDIS_ADDRESS`、`TASK_REDIS_PASSWORD`（task-rpc 侧 `services/task/etc/task.yaml`）
3. `uploads` 共享卷：task-rpc 生成导出文件、gateway 的 `/api/v1/files/uploads/*` 对外提供下载，两个进程必须能读到同一份物理文件。docker-compose 已经配好一个命名卷 `uploads` 同时挂载给 `app`/`task` 两个服务；非容器化部署（Supervisor 方式）需要保证两个进程的 `consts.UploadDir`（`./uploads`，两边 `internal/consts`/`services/task/internal/consts` 各自定义，值必须一致）指向同一个物理目录，比如都软链到同一个共享路径。

**如何验证生效**：`docs/progress.md` 本轮条目记录了一次真实端到端验证（借用远程 MySQL，起 task-rpc + 单体真实连上，提交一个操作日志导出任务，确认状态流转、文件生成、`admin_notification` 记录三方面都符合预期，验证完已清理测试数据）——生产部署时按同样的路径（提交一个真实导出任务、检查任务详情接口的 `result` 字段、检查下载链接可访问）复核一遍即可，不需要重新设计验证方法。

**状态**：`已就绪，待生产环境实际部署时执行`（本机无 Docker，`docker compose up` 本身未做容器化实测,只做过非容器化的直接进程验证）。

### 6 · sdk-rpc 拆分（Phase 2 第二个落地的服务）

**触发条件**：`services/sdk/` 完整拆分完成（`docs/progress.md` 对应条目），sdk-rpc 成为独立部署单元。sdk 域没有调度器/Redis 锁，比 task-rpc 更简单，但多了两处 task-rpc 没有的复杂点：① `SDKAuthMiddleware`/`SDKRateLimitMiddleware`/`SDKCallLogMiddleware` 三个中间件仍在 gateway，内部实现改成调 sdk-rpc；② 单体内嵌的 `pkg/taskcallback` server（`internal/rpcserver/taskcallback/server.go`）的 `sdk_call_log` 导出分支改成回调 sdk-rpc 新增的 `SdkCallLogExport` 方法。

**部署时要做什么**：
1. 建 `admin_sdk` schema（逻辑隔离，物理上和主库同一个 MySQL 实例）：`mysql -uroot -p<pass> -e "CREATE DATABASE IF NOT EXISTS admin_sdk CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"`，然后跑 `db/services/sdk/sdk/create_table_sdk.sql`（**只需要这一个文件**——`init_sdk.sql` 写的是 `admin_menu`/`admin_permission`/`admin_api`，物理属于主库 iam 部分，已经在主库初始化时跑过，不属于 `admin_sdk` schema）。docker-compose 场景下这一步已经自动化，见 `db/services/init-sdk-db.sh` + `docker-compose.yml` 的 mysql 服务第三个 init 脚本挂载。
2. 生产环境实际地址配置（均通过环境变量，`conf.UseEnv()` 已接入，fail-fast 不给静默默认值，和 `JWT_ACCESS_SECRET`/task-rpc 那批同一套约定）：
   - `SDK_RPC_ENDPOINT`（gateway 侧 `etc/admin-api.yaml` 的 `SdkRpc.Endpoints`，指向 sdk-rpc 实际监听地址）
   - `SDK_MYSQL_DSN`（sdk-rpc 侧 `services/sdk/etc/sdk.yaml`）
   - `SDK_REDIS_ADDRESS`、`SDK_REDIS_PASSWORD`（sdk-rpc 侧 `services/sdk/etc/sdk.yaml`；sdk-rpc 业务本身不用 Redis，纯粹是满足 goctl 生成 Model 的 `CachedConn` 强制要求非空缓存节点，和 gateway 共享同一个 Redis 实例即可，不需要单独的 Redis 部署）
   - `RateLimitDefault`（`services/sdk/etc/sdk.yaml` 里的静态配置项，不是环境变量，默认 60——取代原来读字典 `sdk_rate_limit_default` 的做法，字典数据本身仍留在 `admin_menu` 所在的主库，但已经不再被任何代码读取，生产环境如果要调整这个默认值，改这个配置项然后重启 sdk-rpc，不要指望改字典生效）。
3. `SdkCallLogExport` 是本次拆分新增、17-async-eventing.md 原文档没有的一个 `pkg/taskcallback`-类似回调（但不是复用 `pkg/taskcallback` 契约本身，是 sdk.proto 里单独定义的方法），单体内嵌的 TaskCallback server 需要重启才能加载新代码；如果部署时先重启了 sdk-rpc、还没重启单体，SDK 调用日志导出任务会失败（`fetchSdkCallLog` 报错），不是数据丢失，任务失败后可以重新提交导出任务，但建议按顺序：先部署 sdk-rpc → 确认它监听正常 → 再重启单体 gateway（一次性完成 TaskCallback server 和其余 SdkRPC 调用点的切换）。

**如何验证生效**：`docs/progress.md` 本轮条目记录了一次真实端到端验证（借用远程 MySQL + 本地 Redis，起 sdk-rpc + 单体真实连上）：直接调 sdk-rpc 的全部 14 个 RPC 方法（API Key/接口 CRUD、绑定、调用记录列表/导出、VerifyApiKey/GetEffectiveRateLimit/RecordCallLog）+ 通过真实 HTTP 请求走完整 `SDKAuthMiddleware→SDKRateLimitMiddleware→handler→SDKCallLogMiddleware` 链路（含限流生效、调用日志真实落库），验证完已清理测试数据。生产部署时按同样的路径（建一个真实 API Key、绑定一个接口、发几次请求确认鉴权/限流/调用记录都符合预期）复核一遍即可。

**已知遗留（本次验证中发现，不属于本轮拆分引入，但值得部署前留意）**：`sdk_key_api` 表的唯一索引 `uk_sdk_key_api (sdk_key_id, sdk_interface_id)` 不包含 `deleted_at`，`SaveApiKeyBindings` 的"软删旧绑定 + 插入新绑定"模式在给同一个 Key 重复绑定同一个接口时会触发 `Duplicate entry` 报错（软删除的旧行仍然占用唯一索引位）。这是拆分前就存在的原始代码行为（本轮只是原样搬迁，未修改这段逻辑），真实生产使用中大概率只在管理员反复调整同一个 Key 的接口授权时才会撞到，建议后续单独排期修复（唯一索引加上 `deleted_at` 或改成物理删除旧绑定），不在本轮 sdk-rpc 拆分范围内处理。

**状态**：`已就绪，待生产环境实际部署时执行`（本机无 Docker，`docker compose up` 本身未做容器化实测，只做过非容器化的直接进程验证）。

### 7 · chat-rpc 拆分（Phase 2 第三个落地的服务）

**触发条件**：`services/chat/` 完整拆分完成（`docs/progress.md` 对应条目），chat-rpc 成为独立部署单元。这是 5 个服务里第一个涉及 WS↔gRPC 双向流桥接、也是 `stream:chat.user.created`（Redis Streams）第一次真正投入生产路径的服务，比 task-rpc/sdk-rpc 多了两处新增复杂点：① gateway 侧 `internal/handler/chat/chatwshandler.go` 从终结 WebSocket 连接改成"终结 WS + 桥接一条到 chat-rpc 的 gRPC 双向流"；② 单体新增一个临时的 `IamCallback` zrpc server（`internal/rpcserver/iamcallback/`，和已有的 `TaskCallback` 同一个"iam-rpc 真正拆分前的过渡方案"模式），供 chat-rpc 回调枚举存量用户 / 取用户展示信息。

**部署时要做什么**：
1. 建 `admin_chat` schema：`mysql -uroot -p<pass> -e "CREATE DATABASE IF NOT EXISTS admin_chat CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"`，然后跑 `db/services/chat/chat/create_table_chat.sql`（`chat`/`chat_user`/`chat_message` 三张表；`init_chat.sql` 写的是 `admin_menu`/`admin_permission`/`admin_api`，物理属于主库 iam 部分，已经在主库初始化时跑过）。docker-compose 场景下已自动化，见 `db/services/init-chat-db.sh` + `docker-compose.yml` 的 mysql 服务第四个 init 脚本挂载。
2. 生产环境实际地址配置（均通过环境变量，`conf.UseEnv()` 已接入，fail-fast 不给静默默认值，和 `JWT_ACCESS_SECRET`/task-rpc/sdk-rpc 那批同一套约定）：
   - `CHAT_RPC_ENDPOINT`（gateway 侧 `etc/admin-api.yaml` 的 `ChatRpc.Endpoints`，指向 chat-rpc 实际监听地址）
   - `CHAT_MYSQL_DSN`（chat-rpc 侧 `services/chat/etc/chat.yaml`）
   - `CHAT_REDIS_ADDRESS`、`CHAT_REDIS_PASSWORD`（chat-rpc 侧 `services/chat/etc/chat.yaml`；这个 Redis 既满足 goctl 生成 Model 的 `CachedConn` 缓存节点要求，也真正用于消费 `stream:chat.user.created`，和 gateway 共享同一个 Redis 实例）
   - `IAM_CALLBACK_ENDPOINT`（chat-rpc 侧 `services/chat/etc/chat.yaml` 的 `IamCallbackRpc.Endpoints`，指向单体内嵌 `IamCallback` server 的地址，对应 gateway 侧 `etc/admin-api.yaml` 新增的 `IamCallbackRpc.ListenOn`，默认 `0.0.0.0:9002`，和 `TaskCallbackRpc`（9001）同一个模式，docker-compose 场景下容器间要能互相连到所以监听 `0.0.0.0` 不是 `127.0.0.1`）
3. 部署顺序建议：先部署 chat-rpc → 确认它监听正常且能连上 `IAM_CALLBACK_ENDPOINT` → 再重启单体 gateway（一次性完成 `IamCallback` server 上线和 gateway `ChatRPC`/WS 桥接切换）。如果先重启了 gateway、chat-rpc 还没起来，`ChatRPC` 走 `NonBlock: true` 不会阻塞 gateway 启动，但 `/api/v1/chats*` 接口和 WS 端点在 chat-rpc 真正连上之前会持续报错，不是数据丢失。
4. 单元测试范围（`.github/workflows/ci.yml` 的 `unit-test` job）已加上 `./services/chat/...`。

**如何验证生效**：`docs/progress.md` 本轮条目记录了一次真实端到端验证（借用远程 MySQL + 本地 Redis，起 gateway + chat-rpc 真实连上）：① IAM 建用户（真实调用 `UserDomainService.CreateUser`）触发 `stream:chat.user.created`，chat-rpc 消费者在秒级内完成默认群加入 + 与全部存量用户建私聊，全程真实 DB 落库验证；② 全部 11 个 unary CRUD 接口（`ChatList`/`ChatMessageList`/`ChatMessageSend`/`ChatMessageListAdmin`/`ChatMessageDelete`/`ChatGroupList`/`ChatGroupCreate`/`ChatGroupUpdate`/`ChatGroupDelete`/`ChatGroupDetail`/`ChatGroupMemberList`/`ChatGroupMemberAdd`/`ChatGroupMemberRemove`）通过真实 HTTP 请求走完整链路验证，`ChatList`/`ChatGroupDetail` 返回的部门名/角色名确认是真实回调 `IamCallback` 拿到的数据；③ 用一个最小 Go WS 客户端连 `/api/v1/chats/ws`，验证 `GetOnlineUserCount` 从 0 变 1，再发一条 `ChatMessageSend` 确认 WS 客户端实时收到广播（wire JSON 格式和拆分前逐字段一致）。验证完已清理全部测试数据（用户、群组、私聊、消息），行数复核回到基线。生产部署时可以按同样的路径（建群、发消息、开一个 WS 连接看能不能收到推送）复核一遍。

**已知遗留**：
- `docker compose up` 本身未做容器化实测（本机无 Docker，和 task-rpc/sdk-rpc 拆分时同样的限制），已用非容器化真实进程验证过完整链路（含 WS），建议用户在有 Docker 的环境跑一遍确认 `chat` 容器能通过服务名连接。
- WS 客户端→服务端发消息（`SendMessageFrame`，对应 `chat.proto` 里 `ClientFrame` 的 `send` 分支）按文档骨架实现，但当前真实前端只用 WS 做服务端推送、发消息走 REST `ChatMessageSend`，这条路径没有真实前端联调，只做过最小 Go 客户端级别的手工验证（Join/心跳/服务端推送路径），如果后续前端真的要切到 WS 发消息，需要单独联调一次。

**状态**：`已就绪，待生产环境实际部署时执行`（本机无 Docker，`docker compose up` 本身未做容器化实测，已做过非容器化的直接进程 + 真实 WS 连接验证）。

### 8 · content-rpc 拆分（Phase 2 第四个落地的服务）

**触发条件**：`services/content/` 完整拆分完成（`docs/progress.md` 对应条目），content-rpc 成为独立部署单元。这是 blog（标签/文章/审核/友情链接/社交信息/公共展示）+ video（管理/采集/公共展示）合并拆出的服务，`18-service-extraction-runbook.md` 2.4 节定性为"文件数最多但架构最简单"，没有新增跨服务机制（`M3u8Proxy`/`VideoCollectOptions` 不涉及域数据，继续留在 gateway；`PublicBlogAuthorInfo`/`BlogArticleAudit`/`BlogArticleAuditUnpublish` 复用已有的 `IamCallback` 回调模式，只是新增了 `signature` 字段和 `RecordAuditLog` 方法）。

**部署时要做什么**：
1. 建 `admin_content` schema：`mysql -uroot -p<pass> -e "CREATE DATABASE IF NOT EXISTS admin_content CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"`，然后依次跑 `db/services/content/blog/create_table_blog.sql`、`db/services/content/blog_extension/create_table_blog_extension.sql`、`db/services/content/video/create_table_video.sql`（共 7 张表；三个模块各自的 `init_*.sql` 写的是 `admin_menu`/`admin_permission`/`admin_api`，物理属于主库 iam 部分，已经在主库初始化时跑过）。docker-compose 场景下已自动化，见 `db/services/init-content-db.sh` + `docker-compose.yml` 的 mysql 服务第五个 init 脚本挂载。
2. 生产环境实际地址配置（均通过环境变量，`conf.UseEnv()` 已接入，fail-fast 不给静默默认值，和 `JWT_ACCESS_SECRET`/task-rpc/sdk-rpc/chat-rpc 那批同一套约定）：
   - `CONTENT_RPC_ENDPOINT`（gateway 侧 `etc/admin-api.yaml` 的 `ContentRpc.Endpoints`，指向 content-rpc 实际监听地址）
   - `CONTENT_MYSQL_DSN`（content-rpc 侧 `services/content/etc/content.yaml`）
   - `CONTENT_REDIS_ADDRESS`、`CONTENT_REDIS_PASSWORD`（content-rpc 侧 `services/content/etc/content.yaml`；纯粹满足 goctl 生成 Model 的 `CachedConn` 缓存节点要求，业务本身不用 Redis，和 gateway 共享同一个 Redis 实例）
   - `IAM_CALLBACK_ENDPOINT`（content-rpc 侧 `services/content/etc/content.yaml` 的 `IamCallbackRpc.Endpoints`，指向单体内嵌 `IamCallback` server 的地址，和 task-rpc/chat-rpc 用的是同一个 `0.0.0.0:9002`）
3. **`services/content/etc/content.yaml` 的 `Limits` 段需要人工核对一次**：原来读字典（`blog_article_title_max_length`/`blog_article_summary_length`/`blog_article_top_max_count`/`blog_tag_name_max_length`/`blog_friend_link_*_max_length`/`blog_social_info_*_max_length`，共 10 个）决定的长度/数量上限，改成了这个静态配置段，默认值已对齐字典种子数据；如果生产环境之前通过后台字典管理改过这些值，需要手动同步改这个配置项，否则会变回默认值。
4. 部署顺序建议：先部署 content-rpc → 确认它监听正常且能连上 `IAM_CALLBACK_ENDPOINT` → 再重启单体 gateway。`ContentRPC` 走 `NonBlock: true` 不会阻塞 gateway 启动，但 `/api/v1/blog/*`、`/api/v1/videos*` 系列接口在 content-rpc 真正连上之前会持续报错。
5. 单元测试范围（`.github/workflows/ci.yml` 的 `unit-test` job）已加上 `./services/content/...`。

**如何验证生效**：`docs/progress.md` 本轮条目记录了一次真实端到端验证（借用远程 MySQL + 本地 Redis，起 gateway + content-rpc 真实连上，对着真实存量数据——27 篇已发布文章、5 个标签等——做只读验证，写操作全部用新建的测试数据并在验证后精确物理清理）：① 全部公共只读接口（文章列表/详情/上一篇/下一篇/统计、标签/友情链接/社交信息列表、作者信息）对着真实数据验证通过；② 用临时管理员账号 + 临时角色（绑定"全部权限"）验证了标签/文章/友情链接/社交信息的完整 CRUD，文章的提交→审核→上架→置顶→取消置顶→审核下架全流程状态流转正确，审核/下架的审计日志确认真实写入了 `audit_log` 表（回调 `IamCallback.RecordAuditLog` 全链路打通）；③ 视频采集（`VideoCollect`）、管理端视频新增/列表/删除验证通过。

**已知遗留**：
- `docker compose up` 本身未做容器化实测（本机无 Docker，和前三次服务拆分同样的限制），已用非容器化真实进程验证过完整链路，建议用户在有 Docker 的环境跑一遍确认 `content` 容器能通过服务名连接。
- ~~**发现一个预先存在、和本轮拆分无关的真实 bug，未修复**：`services/content/internal/model/video/videomodel_gen.go` 的自定义 `Update` 方法……~~——**已修复（2026-07-13，见 `docs/progress.md` 对应条目）**：根因是 `video` 表 DDL 后来加了 `uuid`/`god_num`/`xlzz_urls`/`type` 四个字段，但 `videomodel_gen.go` 一直没有跟着重新生成。已按项目规则重新执行 `generate-model.sh` 生成正确的 `Insert`/`Update`（`git log -S` 定位到的引入提交 `9580866` 之后首次真正修复），并顺手处理了 regenerate 带来的连带 break（`video_repository.go` 里对已被新模板移除的 `model.FindPage` 的调用，改成统一走仓库层已有的 `findPageWithFilter`）。`go build`/`go vet`/`go test`/`golangci-lint` 全绿，`PUT /api/v1/videos` 现在参数绑定正确。

**状态**：`已就绪，待生产环境实际部署时执行`（本机无 Docker，`docker compose up` 本身未做容器化实测，已做过非容器化的直接进程 + 真实数据验证；`VideoUpdate` 已知 bug 已于 2026-07-13 修复，不再是部署前需要额外决策的项）。

### 9 · iam-rpc 拆分（Phase 2 第五个、也是最后一个落地的服务）

**触发条件**：`services/iam/` 完整拆分完成（`docs/progress.md` 对应条目），iam-rpc 成为独立部署单元，Phase 2 五个服务拆分至此全部落地。这是体量最大、安全敏感度最高的一次拆分：iam+system+monitoring+misc 四个域（约 94 个 RPC 方法）整体搬出单体，`AuthMiddleware`/`PermissionMiddleware`/`ApiEnabledMiddleware`/`OperationLogMiddleware`/`PerformanceMiddleware` 五个中间件第一次真正从直连数据库切换成调 zrpc client，`pkg/taskcallback`/`pkg/iamcallback` 两个"iam 域还没拆分前的临时回调契约"服务端实现也整体搬到这里，成为它们的永久归宿。

**部署时要做什么**：
1. **不需要新建 schema**：iam-rpc 复用现有的 `"admin"` schema（用户明确选择，见 `docs/progress.md` 本轮条目"MySQL schema 决策"）——iam+system+monitoring+misc 从单体拆分前就一直存在这个库里，task/sdk/chat/content 早已各自独立成 `admin_task`/`admin_sdk`/`admin_chat`/`admin_content`，`"admin"` 库现在事实上就是 iam 专属库，零数据迁移、零改名操作。`db/services/init-dev-db.sh` 第一步（`[1/4] iam 建表+初始化数据`）本来就在维护这部分表，不需要新增 `init-iam-db.sh`。
2. 生产环境实际地址配置（均通过环境变量，`conf.UseEnv()` 已接入，fail-fast 不给静默默认值）：
   - `IAM_RPC_ENDPOINT`、`IAM_CALLBACK_RPC_ENDPOINT`（gateway 侧 `etc/admin-api.yaml` 的 `IamRpc.Endpoints`/`IamCallbackRpc.Endpoints`，两者指向同一个 iam-rpc 地址，只是调用不同的 gRPC service）
   - `IAM_MYSQL_DSN`（iam-rpc 侧 `services/iam/etc/iam.yaml`，**必须指向现有 `"admin"` schema，不是新库**）
   - `IAM_REDIS_ADDRESS`、`IAM_REDIS_PASSWORD`（iam-rpc 侧，和 gateway/其余四个服务共享同一个物理 Redis 实例——JWT 黑名单、限流滑动窗口、业务缓存都在这个实例上，不能指错）
   - `SDK_RPC_ENDPOINT`、`CHAT_RPC_ENDPOINT`（iam-rpc 侧，前者供原 `TaskCallback.FetchExportData` 的 `fetchSdkCallLog` 分支回调 sdk-rpc 用，后者供 task 通知消费者推送 WS 通知时回调 chat-rpc 用）
   - `TASK_CALLBACK_ENDPOINT`（task-rpc 侧 `services/task/etc/task.yaml`）、`IAM_CALLBACK_ENDPOINT`（chat-rpc/content-rpc 侧对应 yaml）三者**全部从原来的 `app:9001`/`app:9002`（gateway 容器）改指向 iam-rpc 地址**——这是本轮和前四次拆分的关键差异：TaskCallback/IamCallback 两个 server 不再由 gateway 内嵌承载。
3. 部署顺序建议：先部署 iam-rpc → 确认它监听正常且能连上 `"admin"` schema → 再依次重启/部署 task-rpc、chat-rpc、content-rpc（它们的回调 endpoint 都改指向了 iam-rpc）→ 最后重启 gateway。gateway 的 `IamRpc`/`IamCallbackRpc` 走 `NonBlock: true` 不会阻塞启动，但在 iam-rpc 真正连上之前，几乎全部接口（除了少数不依赖 iam-rpc 的纯代理接口）都会报错——这和前四次拆分"只影响对应域接口"的影响面完全不同，iam-rpc 是唯一一个"挂了就全站不可用"的服务，上线/回滚窗口需要格外谨慎。
4. 单元测试范围（`.github/workflows/ci.yml` 的 `unit-test` job）已加上 `./services/iam/...`，同时移除了已删除的 `./internal/domain/... ./internal/repository/... ./internal/rpcserver/...` 三个路径。

**如何验证生效**：`docs/progress.md` 本轮条目记录了一次真实端到端验证（借用远程 MySQL + 本地 Redis，起 iam-rpc + gateway 真实连上，对着真实存量数据验证）：① `GET /api/v1/ping` 返回 `database:ok,redis:ok`；② 超级管理员（`user_id=1`）请求 `GET /api/v1/users` 返回真实数据且和直接查库结果一致，验证了 `PermissionMiddleware` 超级管理员 bypass 分支 + `UserList` RPC；③ `GET /api/v1/monitor/stats` 返回真实聚合统计（用户/角色/权限/菜单/操作日志/登录日志计数），验证了 `MonitorStats` 多表聚合查询；④ `GET /api/v1/public/dict`（无需登录）返回真实字典值，验证了 `ApiEnabledMiddleware.CheckApiEnabled` + `DictGet`；⑤ 用无任何角色的真实用户 id 签发 JWT 请求受保护接口，正确返回 `403`，验证了 `CheckPermission` 拒绝分支。全程只做只读+权限判定验证，未触发任何写操作，验证后核对 `admin_operation_log` 表行数与验证前完全一致。

**已知遗留**：
- `docker compose up` 本身未做容器化实测（本机无 Docker，和前四次服务拆分同样的限制），建议用户在有 Docker 的环境跑一遍确认新增的 `iam` 容器能通过服务名被 `app`/`task`/`chat`/`content` 四个容器连接。
- ~~本轮真实环境验证**没有覆盖** `Login`/`Refresh`（需要真实管理员密码）、`FileUpload`/`FileDownload`（涉及本地磁盘写入）、`Notice`/`Notification` 批量通知创建（涉及真实批量写入）这几类路径……~~——**已补验证（2026-07-13，见 `docs/progress.md` 对应条目）**：起真实 iam-rpc + gateway 进程连上 oldbai，新建一次性测试账号+角色（不使用/不触碰真实管理员账号）验证了 `Login`（含内部的未读公告通知创建）、`Refresh`、`FileUpload`/`FileDownload`（磁盘字节 + `admin_file` 元数据双重核对）、`Notice`/`Notification` 批量创建（真实 3 个用户全部收到通知）全部路径，验证后精确清理，8 张相关表行数核对回到验证前基线。`Login` 迁移前后行为一致性已通过真实数据确认，不再是"理论上应该一致"。

**状态**：`已就绪，待生产环境实际部署时执行`（本机无 Docker，`docker compose up` 本身未做容器化实测；`Login`/`Refresh`/`FileUpload`/`FileDownload`/`Notice`/`Notification` 几类写路径已于 2026-07-13 补验证通过，不再是部署前的待办项）。

### 10 · Phase 3：六容器镜像 + docker-compose 生产切换（`21-cd-and-deployment.md`）

**触发条件**：Phase 3 第三篇（`21-cd-and-deployment.md`）完成——CI 新增 `build-images` matrix job（六个服务并行构建+推送 `ghcr.io/<owner>/admin-<service>:<git-sha>`）、`script/admin.sh` 新增 `compose pull/deploy/status/logs` 四个子命令、新增 `docker-compose.prod.yml`。这是**前置工作**：本次会话只落地代码/配置本身，未做任何真实部署动作，用户后续实际上线时再具体调整（按用户 2026-07-13 明确要求，见 `docs/progress.md` 对应条目）。

**部署时要做什么**（用户真正要切换到 compose 部署时，按顺序做）：
1. **生产服务器一次性准备**：安装 Docker + docker compose plugin；`docker login ghcr.io`（用户自己的 PAT 或 deploy token，不是 CI 用的 `GITHUB_TOKEN`——那个只在 GitHub Actions runner 里有效）；建 `COMPOSE_REMOTE_DIR`（默认 `/home/work/admin-compose`，可通过环境变量覆盖）并把 `admin-server/docker-compose.prod.yml`、六个服务的 `etc/*.yaml` 模板文件放过去。
2. **`.env` 或 shell 环境变量**（生产服务器上的 compose 目录，不提交仓库）：`GHCR_NAMESPACE`（默认 `oldbai555`，如果镜像推到别的 GitHub 账号/组织下要改）、`TAG`（建议每次发布显式传 git sha，不用 `latest`）、`JWT_ACCESS_SECRET`/`JWT_REFRESH_SECRET`（同条目 1，不能沿用本地开发默认值）、五个服务各自的 `<SERVICE>_MYSQL_DSN`/`<SERVICE>_REDIS_ADDRESS`/`<SERVICE>_REDIS_PASSWORD`（同条目 5-9 已经确认的生产地址，此前是 Supervisor 场景下设进进程环境，切到 compose 后改成同一批变量值放进 `.env`）。
3. gateway 容器仍然依赖 `/etc/work/redis.json`（`-redis-config` 机制未变，见 `admin.go`），`docker-compose.prod.yml` 已把宿主机 `/etc/work` 只读挂载进容器，生产服务器上这个文件需要已经存在（和 Supervisor 部署时的要求一致，不是新增依赖）。
4. 首次切换建议顺序：先在生产服务器上验证六个镜像都能被拉取（`bash script/admin.sh compose pull <host>`）→ 确认 `.env` 齐全 → `bash script/admin.sh compose deploy <host>` → `compose status` 确认六个容器都 `Up` → 用条目 5-9 各自记录的验证路径逐个服务复核一遍 → 确认无误后再考虑停掉旧的 Supervisor 进程（不建议同一时刻双跑，端口/Redis 键会冲突）。
5. `GHCR_NAMESPACE` 默认值取自当前仓库实际 owner（`github.com/oldbai555/bgg`），`.github/workflows/ci.yml` 的 `build-images` job 用 `${{ github.repository_owner }}` 自动取值，两处不需要手动保持同步，但如果仓库迁移到别的 GitHub 账号/组织，两处都要注意（CI 侧自动跟着变，`docker-compose.prod.yml` 的 `GHCR_NAMESPACE` 默认值需要手动改或在 `.env` 里覆盖）。

**如何验证生效**：`docker compose -f docker-compose.prod.yml ps` 六个服务全部 `Up`；`GET /api/v1/ping` 走网关公网端口返回 `database:ok,redis:ok`；按条目 5-9 各自记录的真实验证路径（不是重新设计新的验证方法）逐个复核一遍。

**状态**：`已就绪，待生产环境实际部署时执行`（六个 Dockerfile 在 Phase 2 各服务拆分时已落地并可本地 `go build`/镜像构建成功；`build-images` CI job、`compose` 系列命令、`docker-compose.prod.yml` 三块本次新增，均未在真实 Docker/生产环境验证过——本机无 Docker，`docker compose config` 用 Ruby YAML 解析器做了语法校验，未做 `docker build`/`docker compose up` 实测；未引入 Kubernetes 或任何编排器）。

## Phase 2/3 预留条目类型（占位，届时按真实改动追加，现在不编造内容）

以下只是"预计会出现哪类条目"的提示，**不是已确定的条目**，落地时按 Phase 2/3 实际改动的真实细节追加，不要现在就编内容占位：

- Phase 2 服务间凭证（如果引入了服务间认证机制）的生产环境设置方式。
- Phase 3 Telemetry/结构化日志接入后，如果需要调整生产环境日志采集方式（即使不上 ELK/Loki，Supervisor/docker-compose 的日志落盘路径可能需要调整）。

## 完成的定义

本篇没有"完成"这个状态——它的生命周期跟整个 Phase 1-3 重构一样长。**唯一的检查标准**：每次 Phase 收尾（Week 5 / Week 12 / Week 14+）回顾时，核对本文档里 `TBD` 状态的条目是否都已经在对应 Phase 结束前补全为"已就绪，待执行"，不允许某个 Phase 收尾时还有本该在该 Phase 产生却仍是 `TBD` 的条目遗留到下一个 Phase。
