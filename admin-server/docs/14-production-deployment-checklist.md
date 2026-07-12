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

## Phase 2/3 预留条目类型（占位，届时按真实改动追加，现在不编造内容）

以下只是"预计会出现哪类条目"的提示，**不是已确定的条目**，落地时按 Phase 2/3 实际改动的真实细节追加，不要现在就编内容占位：

- Phase 2 每次服务拆分（B.6 顺序：task-rpc → sdk-rpc → chat-rpc → content-rpc → iam-rpc）对应的：新数据库 schema 创建 SQL、旧单体到新服务的数据迁移步骤、镜像/compose/Supervisor 切换步骤、服务间 zrpc `Endpoints` 生产环境实际地址配置。
- Phase 2 服务间凭证（如果引入了服务间认证机制）的生产环境设置方式。
- Phase 3 CI 镜像发布凭证配置（`ghcr.io` 的 push 权限 token 等）。
- Phase 3 Telemetry/结构化日志接入后，如果需要调整生产环境日志采集方式（即使不上 ELK/Loki，Supervisor/docker-compose 的日志落盘路径可能需要调整）。

## 完成的定义

本篇没有"完成"这个状态——它的生命周期跟整个 Phase 1-3 重构一样长。**唯一的检查标准**：每次 Phase 收尾（Week 5 / Week 12 / Week 14+）回顾时，核对本文档里 `TBD` 状态的条目是否都已经在对应 Phase 结束前补全为"已就绪，待执行"，不允许某个 Phase 收尾时还有本该在该 Phase 产生却仍是 `TBD` 的条目遗留到下一个 Phase。
