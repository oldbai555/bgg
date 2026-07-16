# 域目录重组 + API wrapper 可执行迁移清单（Phase 1 Week 1）

> 前置依赖：`01-architecture-target.md` A.1/A.2/A.5。本文档是可执行版清单，逐文件给出迁移目标，不重复 `01` 里的规则推导。

## 1. views/ 迁移清单

`git mv` 保留历史，禁止用删除+新建代替（会丢 blame）。

| 现路径 | 目标路径 | 备注 |
|---|---|---|
| `views/Dashboard.vue` | 不动 | 跨域首页，不属于任何单一业务域 |
| `views/Home.vue` | 不动 | 同上 |
| `views/Login.vue` | 不动 | 同上 |
| `views/error/NoAccess.vue`、`views/error/NotFound.vue` | 不动 | 跨域错误页 |
| `views/system/UserList.vue` | `views/iam/UserList.vue` | |
| `views/system/RoleList.vue` | `views/iam/RoleList.vue` | |
| `views/system/PermissionList.vue` | `views/iam/PermissionList.vue` | |
| `views/system/DepartmentList.vue` | `views/iam/DepartmentList.vue` | |
| `views/system/MenuList.vue` | `views/iam/MenuList.vue` | |
| `views/system/ApiList.vue` | `views/iam/ApiList.vue` | 对应 `admin.api` 的 `iam/api` group |
| `views/system/Profile.vue` | `views/iam/Profile.vue` | 自助资料页，数据归属仍是 iam.user，与其他 iam 页面放一起而不是单独跨域 |
| `views/system/ConfigList.vue` | `views/system/ConfigList.vue` | 留在 system，域名不变 |
| `views/system/DictTypeList.vue`、`DictItemList.vue` | 不动（留在 `views/system/`） | |
| `views/system/FileList.vue` | 不动 | |
| `views/system/NoticeList.vue`、`NotificationList.vue` | 不动 | |
| `views/system/AuditLogList.vue` | `views/monitoring/AuditLogList.vue` | |
| `views/system/LoginLogList.vue` | `views/monitoring/LoginLogList.vue` | |
| `views/system/OperationLogList.vue` | `views/monitoring/OperationLogList.vue` | |
| `views/system/PerformanceLogList.vue` | `views/monitoring/PerformanceLogList.vue` | |
| `views/system/MonitorList.vue` | `views/monitoring/MonitorList.vue` | |
| `views/system/MetricStats.vue` | `views/monitoring/MetricStats.vue` | |
| `views/system/TaskList.vue` | `views/task/TaskList.vue` | 后端已是独立 `task-rpc` |
| `views/blog/*.vue`（6 个文件） | `views/content/Blog*.vue` | 文件名不变，只搬目录 |
| `views/video/VideoList.vue`、`VideoPlayer.vue` | `views/content/Video*.vue` | 与 blog 合并进 `content/`，对应后端 `content-rpc` |
| `views/chatroom/ChatList.vue`、`ChatGroupList.vue`、`ChatMessageList.vue` | `views/chat/*.vue` | 目录改名 `chatroom` → `chat`，与后端 `chat-rpc` 命名对齐 |
| `views/sdk/*.vue`（3 个文件） | 不动 | 已经是独立域目录 |
| `views/public/*.vue`（4 个文件） | 不动（本轮不改目录结构，视觉/响应式重构见 `06`） | |
| `views/temp/*.vue`（4 个文件） | 删除 | 死代码，见 `07-cleanup-and-tooling.md` §1 |

迁移后 `views/system/` 只剩 `ConfigList`/`DictTypeList`/`DictItemList`/`FileList`/`NoticeList`/`NotificationList` 6 个真正的 system 域页面，不再是"大杂烩"。

## 2. 路由 meta / 菜单 `component` 字段联动迁移步骤

domains 目录搬迁后，`admin_menu` 表里存量数据的 `component` 字段值（如 `system/UserList`）会指向已经不存在的路径，导致 `01` A.5 描述的 `resolveComponent` 静默失败。**这是数据库数据，不是纯前端改动**，必须按以下顺序执行，不能只改前端：

1. 前端先完成 `views/` 目录迁移 + `router/index.ts` 里手写路由（第 64-105 行的 `system/role`、`system/permission` 等 6 条 `/system/*` 静态路由）里的 `component` import 路径同步更新。
2. 盘点 `admin_menu` 表里所有 `component` 字段引用了本次迁移涉及路径的记录（`system/UserList`、`system/RoleList`、`blog/*`、`video/*`、`chatroom/*` 等），生成一份对照 SQL（`UPDATE admin_menu SET component = 'iam/UserList' WHERE component = 'system/UserList'` 这种一一对应），作为增量 SQL 文件（依 `00-workflow.md` 的字典/业务 SQL 分离原则，这类"数据修正"SQL 建议放 `admin-server/db/services/iam/menu/migrations/` 下，文件名 `fix_menu_component_path_YYYYMMDD.sql`）。**注意这是数据修正 SQL，不是字典新增 SQL**：不要套用字典迁移文件的 `INSERT ... ON DUPLICATE KEY UPDATE` 写法（那是插入语义，这里没有新增记录）——正确的幂等方式是 `UPDATE ... WHERE component = '<旧值>'`，每条 `WHERE` 条件天然只匹配"还没被改过"的旧值，重复执行时旧值已不存在、`UPDATE` 影响 0 行，本身就是幂等的，不需要额外包装。
3. 该 SQL **必须由用户确认执行**（`00-workflow.md` "何时必须停下来问用户"——任何数据库 SQL 的实际执行都要停下确认），不在"开发期执行策略例外"范围内，因为这次不是 dev 库的脚手架生成数据，而是修正现网菜单数据。
4. SQL 执行后本地登录验证：菜单可点击、页面可达、无 `[Router] 无法解析组件` 控制台报错。

## 3. API wrapper 分域清单

依 `01` A.2 规则，8 个域各建一个 `src/api/<domain>.ts`。下表按 `admin.api` 已确认的 `group:` 声明分组（第一轮探索已核实的完整列表），每个 wrapper 覆盖对应 group 的全部接口：

| wrapper 文件 | 覆盖的 `admin.api` group | 现状 |
|---|---|---|
| `src/api/iam.ts`（新增） | `iam/auth`、`iam/user`、`iam/role`、`iam/permission`、`iam/department`、`iam/menu`、`iam/api`、`iam/user_role`、`iam/role_permission`、`iam/permission_menu`、`iam/permission_api` | 目前无 wrapper，视图直接调 `generated/`，是迁移重点 |
| `src/api/system.ts`（新增） | `system/file`、`system/config`、`system/dict_type`、`system/dict_item`、`system/dict`、`system/notice`、`system/notification` | 同上；注意 `system/dict` 的字典批量查询目前走 `src/api/public.ts`（登录前也要用），迁移时**保留 `public.ts` 对该接口的引用**，不要把公共可访问的字典查询也收进需要鉴权语境的 `system.ts`，两者可以都导出同一个底层调用，但入口保留独立文件是为了不引入循环依赖 |
| `src/api/monitoring.ts`（新增，吸收现有 `metric.ts`） | `monitoring/monitor`、`monitoring/metric`、`monitoring/metric_admin`、`monitoring/operation_log`、`monitoring/login_log`、`monitoring/audit_log`、`monitoring/performance_log` | 现有 `metric.ts` 内容原样迁入，文件删除、内容合并，调用处 import 路径同步改 |
| `src/api/misc.ts`（新增，视 07 号文档清理结论决定是否需要） | `misc/ping`、`misc/public`、`misc/demo`、`misc/daily_short_sentence` | 若 `demo`/`daily_short_sentence` 按 `07` 结论删除，则此 wrapper 只保留 `ping` 相关或直接不建 |
| `src/api/content.ts`（新增，合并现有 `blog.ts` + `video.ts`） | `blog/tag`、`blog/article`、`blog/article_audit`、`blog/public`、`blog/friend_link`、`blog/social_info`、`video/video`、`video/m3u8`、`video/video_collect`、`video/public` | 两个文件函数原样合并，按原 `blog.ts`/`video.ts` 内部分组用注释分节，不强行拆分文件（"是同一个后端服务"这件事只需要文件合并体现，不需要过度设计成子模块） |
| `src/api/chat.ts`（新增） | `chat/chat`、`chat/group`、`chat/message` | 目前无 wrapper |
| `src/api/sdk.ts`（新增） | `sdk/sdk`、`sdk/public` | 目前无 wrapper |
| `src/api/task.ts`（新增） | `task/task`、`task/public` | 目前无 wrapper |
| `src/api/public.ts`（保留不变） | 跨域公共/未登录可访问接口（字典查询等技术性入口） | 不动 |

### 迁移方式（每个域重复执行）

1. 新建 `src/api/<domain>.ts`，参考现有 `src/api/blog.ts` 的写法风格（错误处理、返回类型、path 前缀修正逻辑）。
2. 全仓 `grep -rn "from '@/api/generated" src/views/<相关目录> src/components` 找出该域下所有直接调用点。
3. 逐个调用点改成从新 wrapper 导入对应函数；类型 import（`import type {...}`）保留直接从 `generated/` 引入，不搬。
4. 该域全部视图迁移完后，`npm run typecheck` 确认无残留类型错误，`grep` 确认该域视图目录下 `from '@/api/generated'` 只剩类型 import。

### 完成的定义

- 8 个域 wrapper 全部建立（或按 `07` 结论明确 `misc.ts` 不需要建）。
- `grep -rn "from '@/api/generated" src/views src/components` 命中的行里，只有 `import type` 形式，函数调用形式的引用归零。
- 菜单 `component` 字段修正 SQL 已提交用户确认并在开发库执行，登录后动态路由可正常访问所有迁移后的页面。
- `npm run typecheck` + `npm run build` 通过。
