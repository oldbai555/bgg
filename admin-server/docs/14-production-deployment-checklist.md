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

**状态**：`TBD`（A.5 代码改动尚未落地，占位）。

### 2 · MySQL/Redis 配置的环境变量化现状确认

**触发条件**：核查发现 `admin.go` 已经支持 `-mysql-config`/`-redis-config` 命令行参数指定 `/etc/work/mysql.json`/`/etc/work/redis.json` 路径，说明 MySQL/Redis 连接信息**已经**走外部文件、不提交仓库——这部分不是本轮新增的改动，A.5 只处理 JWT 密钥这一项。

**部署时要做什么**：无新增动作，仅在此记录现状，避免后续误以为 MySQL/Redis 配置也需要迁移。如果 Phase 1 执行过程中发现 `admin.go` 现有的 `/etc/work/mysql.json`/`/etc/work/redis.json` 机制需要调整（比如要支持 `conf.UseEnv()` 统一），再单独补一条新条目，不要复用这一条。

**如何验证生效**：不适用（无动作）。

**状态**：`已就绪`（现状确认，非待办）。

### 3 · IAM 事务化改造涉及的 RBAC 种子数据（TBD，Phase 1 Week 2 补全）

**触发条件**：`04-domain-iam-chat.md` 任务 1（`user_create_logic.go` 修复）+ `01-architecture-target.md` A.1（`Repository.Transact`）落地后，"新用户加入默认群 + 批量私聊初始化"这部分逻辑会从同步改成异步（通过 `internal/domain/task` 现有调度器派发）。如果这个改造过程中涉及新增/调整任何 RBAC 权限种子数据（比如新的任务类型需要对应的操作权限，或者默认群组 `chat_id=1` 这类硬编码 ID 的初始化数据在生产环境是否已存在需要核实），这里要记一条。

**部署时要做什么**：**TBD——待 Phase 1 Week 2 实际执行 IAM+Chat 域改造时，根据当时产生的真实 SQL 差异填写这里的具体命令**，本篇现在不预先编造内容。填写时要包含：具体的 `INSERT`/`UPDATE` SQL 语句（或指向 `db/migrations/` 下具体文件名）、执行顺序（是否依赖条目 1 的密钥改动先生效）。

**如何验证生效**：TBD，随上面一起补全。

**状态**：`TBD`。

### 4 · Dockerfile / docker-compose 引入后的部署方式变化（如果 Phase 1 期间就切换）

**触发条件**：`09-ci-cd-and-deployability.md` 设计的 Dockerfile/docker-compose 是 Phase 1 的产出，但**是否在 Phase 1 就把生产部署从当前的 `script/admin.sh` 构建→scp→`supervisorctl restart` 流程切到 docker-compose，还是继续用 Supervisor、只把 Docker 化留到 Phase 3 统一切换**，属于 `10-dev-execution-and-review-points.md` 第 2 节"生产环境部署动作本身"的必停项，需要用户拍板。

**部署时要做什么**：TBD，取决于用户的决定。如果决定 Phase 1 就切：需要在此追加完整的"Supervisor → docker-compose 迁移步骤"（镜像构建、`.env` 生产值配置、`docker compose up -d`、旧 Supervisor 进程停止确认、回滚方案）。如果决定留到 Phase 3：这条先标记为"留待 Phase 3"，不需要现在填内容。

**如何验证生效**：TBD。

**状态**：`TBD`（待用户决定切换时机）。

---

## Phase 2/3 预留条目类型（占位，届时按真实改动追加，现在不编造内容）

以下只是"预计会出现哪类条目"的提示，**不是已确定的条目**，落地时按 Phase 2/3 实际改动的真实细节追加，不要现在就编内容占位：

- Phase 2 每次服务拆分（B.6 顺序：task-rpc → sdk-rpc → chat-rpc → content-rpc → iam-rpc）对应的：新数据库 schema 创建 SQL、旧单体到新服务的数据迁移步骤、镜像/compose/Supervisor 切换步骤、服务间 zrpc `Endpoints` 生产环境实际地址配置。
- Phase 2 服务间凭证（如果引入了服务间认证机制）的生产环境设置方式。
- Phase 3 CI 镜像发布凭证配置（`ghcr.io` 的 push 权限 token 等）。
- Phase 3 Telemetry/结构化日志接入后，如果需要调整生产环境日志采集方式（即使不上 ELK/Loki，Supervisor/docker-compose 的日志落盘路径可能需要调整）。

## 完成的定义

本篇没有"完成"这个状态——它的生命周期跟整个 Phase 1-3 重构一样长。**唯一的检查标准**：每次 Phase 收尾（Week 5 / Week 12 / Week 14+）回顾时，核对本文档里 `TBD` 状态的条目是否都已经在对应 Phase 结束前补全为"已就绪，待执行"，不允许某个 Phase 收尾时还有本该在该 Phase 产生却仍是 `TBD` 的条目遗留到下一个 Phase。
