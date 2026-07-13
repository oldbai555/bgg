# 项目总览

本仓库是前后端分离的管理系统：
- `admin-server`：go-zero 后端，统一提供所有 API（后台管理接口 + 公共接口）
- `admin-frontend`：Vue3 + TS 前端，同时承载后台管理页面（`/admin/*`，需登录）和公共展示页面（`/blog/*`、`/videos/*`，无需登录）

规则文件地图：
- 本文件（`00-workflow.md`）：全局工作流 + 绝对禁止事项
- `05-go-zero-ai-context.md`：go-zero AI 上下文（zero-skills 子模块，编辑 `admin-server/**` 时）
- `06-mcp-toolchain.md`：MCP 工具链使用契约（CodeGraph / Engram / context7 等，**会话探索代码前必读**）
- `07-anthropic-skills.md`：Anthropic 官方 Skills 市场接入说明 + 何时用哪个 skill
- `10-go-code-style.md`：`admin-server/**` 后端规范
- `20-frontend.md`：`admin-frontend/**` 前端规范
- `21-public-pages.md`：公共展示页专属样式/交互契约
- 根目录 `AGENTS.md`：面向任意 AI 工具的整合版操作手册
- `docs/AI工具链上手.md`：Gentle-AI / CodeGraph / MCP / Engram 换设备与第三人上手（**新维护者必读**）
- `docs/后端开发进度.md` / `docs/前端开发进度.md`：已完成功能、技术决策记录、关键代码位置索引（历史/背景，不是规则）

# AI 工具链（新设备 / 新维护者）

日常开发：**Cursor + Claude Code 插件**，生态为 Gentle-AI + CodeGraph + Engram。

1. clone 后执行 `make setup-ai`，完全重启 Cursor
2. **必须手动完成** `docs/AI工具链上手.md` 中「初始化后必做清单」（`GO_ZERO_MCP_PATH`、`go-lsp` 等 `setup-ai` 不会自动安装）
3. 修改 `.cursor/rules/` 后执行 `make sync-claude-rules`
4. 维护者在 Cursor 增删 MCP 后执行 `make sync-claude-mcp-import` 并 **commit `.mcp.json`**
5. 跨设备收工：`make engram-sync-push` → `git push`；换设备：`make engram-sync-pull`

# 标准开发流程（严格按顺序，新增功能模块必走）

步骤1：明确功能需求，确定模块名称（snake_case，如 `order`、`blog_tag`）

步骤2：评估是否需要数据字典
- 需要：创建增量字典 SQL 文件，路径 `admin-server/db/services/<service>/<module>/migrations/dict_{module}_YYYYMMDD.sql`

步骤3：**用户**执行 `admin-server/scripts/generate-sql.sh -group <domain>/<module> -name <name>`
- 一次性生成整个模块的骨架，AI 不得代替用户执行此脚本，只能在用户确认生成结果后继续
- 具体产出见下方"新增模块脚手架"一节

步骤4：补齐 SQL 字段（`created_at`/`updated_at`/`deleted_at`，均为 BIGINT 秒级时间戳，默认 0）

步骤5：补齐 `.api` 接口参数和中间件声明
- 中间件按执行顺序声明（见 `10-go-code-style.md`）
- 所有 Query 参数和可选字段必须加 `optional` 标签（详见 `10-go-code-style.md`）

步骤6：**用户**执行 `generate-model.sh <sql_file>`，确认生成 Model 代码

步骤7：**用户**执行 `generate-api.sh <api_file>`，确认生成 Handler/Logic 骨架

步骤8：实现 Repository/Logic 业务逻辑

步骤9：执行 SQL（顺序：字典SQL → 业务表SQL → 权限SQL）

步骤10：启动后端服务测试接口

步骤11：**用户**执行 `generate-ts.sh`，确认生成前端 TS 代码

步骤12：完善前端页面（基于生成的 `.vue` 骨架）

步骤13：前后端联调测试通过

步骤14：更新进度文档（`docs/后端开发进度.md` / `docs/前端开发进度.md`：已完成功能、关键代码位置、技术决策记录）

**完成的定义**：一个功能只有同时满足「联调测试通过」+「进度文档已更新」才算完成，不要仅凭代码写完就声称任务结束。

# 新增模块脚手架（工程化能力，新模块默认走这条路）

admin-server 和 admin-frontend 从设计之初就是"工程化"的：新增一个标准 CRUD 业务模块（如"订单管理""优惠券管理"）不需要从零手写，`generate-sql.sh -group <domain>/<module> -name <name>` 一条命令会同时生成整个模块的四块骨架，且相互匹配、可以直接跑通：

1. **建表 SQL**：`admin-server/db/services/<service>/<module>/create_table_<module>.sql`（`id` 自增主键 + `created_at`/`updated_at`/`deleted_at` + 主键/`deleted_at` 索引，业务字段需手动补充）
2. **初始化 SQL**：`admin-server/db/services/<service>/<module>/init_<module>.sql`（菜单数据「主菜单 + 新增/编辑/删除按钮」、权限数据「list/create/update/delete」、接口数据「GET/POST/PUT/DELETE」、权限-菜单关联、权限-接口关联，全部主键自增、幂等可重复执行；菜单默认挂在"临时目录"下，路径 `/temp/<group>`，后续在菜单管理里手动挪到正式分类）
3. **`.api` 草稿**：`admin-server/api/<group>.api.temp`（`Item`/`ListReq`/`ListResp`/`CreateReq`/`UpdateReq` 类型定义 + 含 List/Create/Update/Delete 四个接口的 `@server` 服务块），需人工复制追加进 `admin-server/api/admin.api`（类型追加进 `type (...)` 块内，服务块追加到文件末尾），追加完删除 `.api.temp`
4. **前端页面骨架**：`admin-frontend/src/views/temp/<GroupUpper>List.vue`（基于 `D2Table`，已包含搜索/列表/新增/编辑/删除，自动调用 `generate-ts.sh` 之后生成的 `<group>Api.list/create/update/delete`），开箱即可用，后续按需完善业务字段和交互

四块产物的生成模板是可定制的 Go 模板文件，都在 `admin-server/scripts/sqlgen/templates/*.tpl`（`create_table.sql.tpl`、`init_module.sql.tpl`、`init_module.api.tpl`、`list_page.vue.tpl`）——如果发现生成骨架的默认写法不符合项目最新约定（例如要统一加表前缀、调整默认字段），应该改模板而不是每次生成后手工修一遍。

**新增标准 CRUD 模块时，默认走这条脚手架路径，不要绕过它直接手写四块骨架**；只有当模块的数据结构/交互明显偏离「单表 + 列表 CRUD」这个标准形态时（如聊天、任务调度、WebSocket 类模块），才不套用脚手架，转为手写。

# 何时必须停下来问用户（不得替用户跳过）

- 任何 `admin-server/scripts/generate-*.sh` 脚本的执行——这些脚本会生成/覆盖文件，必须由用户亲自运行并确认结果
- 任何数据库 SQL 的实际执行（建表、字典、权限初始化）
- 涉及 `script/admin.sh` 的构建/打包/Supervisor 部署操作
- 修改 `etc/admin-api.yaml`、`etc/middleware.yaml` 之外的生产环境配置文件（如 `/etc/work/mysql.json`、`/etc/work/redis.json`）
- 不确定某个字段/接口是否应该新增字典、或某个改动是否会破坏既有数据语义时

## 开发期执行策略的例外（仅限用户明确批准的大型重构项目）

上面的默认规则是为**日常新增业务模块**场景设计的。对于用户已经**明确批准、体量大、跨多周的整体性重构项目**（先例：admin-server 2026-07 单体加固→微服务拆分→可观测性/CD 三阶段重构，规划记录见 `admin-server/docs/10-dev-execution-and-review-points.md`），用户可以针对该项目显式放开以下几类操作，AI 直接执行、事后随 diff 交给用户日常 review，不需要逐次停下确认：

1. `generate-*.sh` 系列脚本的实际执行
2. 本地/开发环境数据库的真实 SQL 执行（建表、字典、权限初始化），前提是目标库明确是本地/团队约定的 dev 库，不是生产 `/etc/work/*.json` 指向的库
3. `make wire`/`goctl rpc` 等生成产物的提交
4. 本地 `docker compose up` 验证（执行环境里确实有 Docker 时；没有则如实告知用户需要本地验证，不得假装验证过）
5. `golangci-lint` 存量问题里能安全自动修复的部分（`gofmt`/`goimports`、明确的 `unused`/`ineffassign`）；改变行为、判断不了是否安全的问题不擅自改，记录下来汇总给用户

**这条例外不会自动生效**：必须是用户对某个具体项目明确授权过（授权范围写在该项目自己的规划文档里），且只在该项目的开发阶段内有效——不延伸到日常业务模块开发，也不延伸到被授权项目之外的其他改动，不能因为一个项目放开过就类推到另一个项目。

即使在被授权的项目范围内，以下情况仍然必须停下来问用户，这条例外不覆盖：
- 任何触及生产配置/密钥本身（`/etc/work/mysql.json`、`/etc/work/redis.json` 等）或生产部署动作本身
- 生产环境真实密钥/凭证的生成或取值（哪怕只是"先给个占位值"也不行）
- 产品/体验取舍（技术方向已定，但具体交互/时机等细节仍是产品判断，不能因为方向定了就自行拍板全部细节）
- 任何 AI 自己判断不准的架构/命名/取舍点——拿不准就问，不是逢清单必停，也不是没在清单里就可以自行决定

# 绝对禁止事项（常见错误速查表）

| 禁止行为 | 为什么错 | 正确做法 |
|---|---|---|
| 跳过脚本执行步骤，手写本该由 goctl/generate-*.sh 生成的文件 | 破坏代码生成一致性，后续再生成会冲突 | 能用工具生成的必须用工具生成，AI 只负责补齐业务逻辑 |
| 修改 `internal/handler`、`internal/model`、`admin-frontend/src/api/generated/*` 等生成目录 | 下次重新生成会覆盖手改内容 | 只在允许手改的文件（Logic、Repository、二次封装层）中改动 |
| Group 使用驼峰命名 | `.api` 文件的 group 名会影响生成路径和路由前缀 | Group 必须 snake_case，如 `user_role`、`dict_type` |
| API 定义使用路径参数 `:id` | 当前 `admin.api` 全部路由都不使用路径参数，混用会破坏一致性 | 一律用 Query/Body 参数，或 `/xxx/detail`、`/xxx/get` 子路径 |
| Query 参数或可选字段缺少 `optional` 标签 | go-zero 的 `httpx.Parse` 会因此报 400 参数解析失败 | 所有 Query 参数、所有可选 Body 字段都要加 `,optional`；必填校验放在 Logic 层，不依赖解析层 |
| Repository 层用 `fmt.Sprintf`/字符串拼接构建 SQL | 有 SQL 注入风险，且与项目风格不一致 | 必须用 `squirrel`（`sq.And/Or/Eq/Like` 等），参考 `services/iam/internal/repository/iam/role_permission_repository.go` |
| 业务表使用物理删除 | 破坏审计/恢复能力，且与 `deleted_at` 软删除体系冲突 | 一律软删除（`deleted_at` 置为当前时间戳） |
| 中间件声明顺序错误 | 部分中间件依赖前序中间件写入的上下文（如权限依赖鉴权结果） | 严格按 Performance → RateLimit → Auth → Permission → OperationLog 顺序（SDK 系列同理） |
| 字典SQL插入到模块自己的 `init_<module>.sql` | 破坏首次部署与增量更新的边界 | 必须新建独立的 `db/services/<service>/<module>/migrations/dict_{module}_YYYYMMDD.sql` |
| 在业务代码中硬编码字符串常量 | 难以维护、容易多处不一致 | 系统级枚举/常量统一放 `internal/consts` |
| 保留旧代码路径和兼容层 | 增加认知负担，掩盖真实调用路径 | 确认无引用后直接删除，不留"以防万一"的兼容代码 |
| 字典枚举 value 从 0 开始 | 0 在筛选场景中被约定为"不筛选/全部"的占位值，会和真实业务含义冲突 | 字典 value 从 1 开始；前端筛选参数用 0 表示不加该筛选条件 |

核心原则：能用工具生成的绝不手写，严格遵循分层架构，前后端协同开发，文档与代码同步更新。
