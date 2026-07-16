# 路由命名空间彻底分离：/bgg/admin/* 与 /bgg/front/*

## 背景（这个问题怎么被发现的）

用户反馈：在后台「博客管理」相关页面 F5 硬刷新后，经常会莫名其妙跳到公共博客列表页；反馈视频管理也有类似现象。

排查过程（用 codegraph + mysql MCP 交叉核对源码和真实数据，不是靠猜）：

- `admin-frontend` 用同一个 `createWebHistory('/admin/')` 同时承载后台管理页面和公共展示页面。两者在路由 **path** 层面共用同一套字面量前缀——公共页（`/blog`、`/blog/:id`、`/videos`、`/videos/:id`）是启动时就注册好的**静态路由**；后台管理页是登录后异步拉取菜单（`userStore.fetchMenus()`）才 `router.addRoute()` 注册的**动态路由**。
- `admin_menu` 表里「博客管理」目录菜单的 `path` 字段就是字面量 `/blog`——和公共博客列表页的路由 path 完全一样。
- F5 硬刷新会让 SPA 重新启动，`initialized`（模块级变量）复位。如果这时地址栏恰好是 `/admin/blog`，浏览器发起的这次导航会在动态路由注册完成**之前**被解析——此时路由表里只有静态路由，`/blog` 精确匹配到公共 `BlogList.vue`，而不是等到菜单异步加载完、动态路由注册后再匹配后台页面。
- `router/index.ts` 里原有的"重新解析"补救逻辑（`if (to.name === 'NotFound') { 重新 resolve }`）只处理"首次匹配到 NotFound"这一种情况，没有处理"首次匹配到了一个存在但错误的路由"这种情况，于是错误地停留在公共页面，用户看到的就是"莫名其妙跳到了公共页"。
- 已排除的次要怀疑：`BlogDetail.vue` 请求失败时没有自动跳转到 `/blog` 的逻辑（读代码确认过），不是这条路径导致的。

用户明确：项目未上线，不用考虑现有公共 URL 的兼容性/SEO，要求"做彻底"、可以联动后端一起调整——因此没有停留在"修补时序竞态"，而是从路由 **path** 设计上把后台和公共彻底分成两个不共享任何前缀段的命名空间，让这类冲突在结构上不可能发生。

## 最终方案

`createWebHistory` base 从 `/admin/` 改为 **`/bgg/`**，内部两个分支：

- **`/front/*`**（公共，免登录）：`/front`（Home 首页）、`/front/blog`、`/front/blog/:id`、`/front/videos`、`/front/videos/:id`
- **`/admin/*`**（后台，需登录）：现有全部后台子路径原样保留最后一段，只在最前面加 `/admin` 前缀，例如 `/blog/article` → `/admin/blog/article`、`/system/role` → `/admin/system/role`、`/video/list` → `/admin/video/list`、`/chatroom/chat` → `/admin/chatroom/chat`、`/sdk/api-key` → `/admin/sdk/api-key`
- 裸 `/`（即 `/bgg/`）→ `redirect: '/front'`，匿名访客默认落在公共首页
- `/:pathMatch(.*)*` → NotFound，全局唯一，两个分支共用

选 `/bgg/` 这层外层前缀（而不是直接裸 `/admin/*` + `/front/*`）除了沿用最初提议外，还避开了 nginx 现有 `location ^~ /video/`（`config/nginxconfig.txt`，一个完全不相关的视频流后端代理，转发到 `127.0.0.1:8889`）——套一层 `/bgg/` 外壳后，前端的 `/admin/video/*` 不会和这个既有后端路径产生任何组合歧义。

这个方案是结构性的：逐条核对现有全部后台子路径（`system` 18 个、`temp` 4 个、`chatroom` 4 个、`sdk` 4 个、blog/video 系列）后确认，新方案下不存在任何"段数相同又字面量相同"的交叉，之前那种"动态路由还没注册时被静态路由误接"的可能性在结构上被消除，不再依赖猜时序。

## 涉及改动

- **`admin-frontend/src/router/index.ts`**：base 改 `/bgg/`；路由表按上述方案重新分支；`isPublicPath` 简化为 `to.path.startsWith('/front')`；顺手加固了"重新解析"逻辑——不再只在 `to.name === 'NotFound'` 时才重新 resolve，只要本次新注册了动态路由，就用 `to.fullPath` 重新解析一次，解析结果和当前不一致就纠正过去（防御性补强，覆盖"匹配到错误页面而非 NotFound"这种理论上还可能出现的情况）
- **`admin-frontend/src/utils/request.ts`**：里面重复了一份 `isPublicPath()`，同步改成 `path.startsWith('/front')`
- 全项目所有 `router.push`/`router.replace`/`:to`/`redirect` 硬编码字面量路径统一按"后台页面加 `/admin`，公共页面加 `/front`"改过一遍，代表性文件：`Login.vue`、`UserMenu.vue`、`MessageNotification.vue`、`TaskFloatBall.vue`、`DefaultLayout.vue`、`AppHeader.vue`、`utils/breadcrumb.ts`、`BlogArticleEdit.vue`、`BlogArticleList.vue`、`content/VideoList.vue`、`PublicHeader.vue`、`Home.vue`、`views/public/*`（Blog/Video 的 List/Detail 四个页面）
- **`stores/notification.ts`**：聊天页检测的字面量改成 `/admin/chatroom/chat`；顺带修了 code review 抓出来的两处预存在问题——① 删掉从未真正用到的 `window.location.hash` 兼容分支（应用全程 `createWebHistory`，不存在 hash 路由场景）；② 不再用裸 `.includes()` 匹配 `window.location.pathname`，改成路径段边界正则 `/(^|\/)admin\/chatroom\/chat(\/|$)/`——`admin_menu` 里还有 `/admin/chatroom/chat-message`、`/admin/chatroom/chat-group` 两个同前缀但语义不同的管理页面，裸 `includes` 会把它们也误判成"在聊天页"，导致真实聊天消息的未读提醒被错误抑制；也不再判断不存在的 `/admin/temp/chat`（核实 `admin_menu` 后确认这条路径从未真实存在过）
- **`admin-frontend/vite.config.ts`**：`base: '/admin/'` → `base: '/bgg/'`，dev 重定向中间件、`open` 路径同步改
- **`config/nginxconfig.txt`**：`location /admin` → `location /bgg`，相关正则和根路径 302 目标（`/admin/` → `/bgg/front/`）同步改；`alias /tmp/admin-cache` 这个**文件系统路径**故意没改——它牵涉部署脚本/tmpfs 挂载，不在本次前端路由改动范围内
- **后端数据迁移**：新增 `admin-server/db/services/iam/menu/migrations/add_admin_prefix_to_menu_path_20260714.sql`，给 `admin_menu.path` 统一加 `/admin` 前缀，顺带修正字典项 `chat_config`/「在线聊天页面路径」的值。已通过 codegraph 确认后台权限/接口校验按菜单 **id** 关联 `admin_permission_menu`，从不按 `path` 字符串匹配，这条 UPDATE 不影响任何鉴权逻辑。**这条 SQL 需要用户在本地/团队 dev 库亲自执行**，AI 不代为执行数据库写操作
- **文档同步**：`AGENTS.md`、`.cursor/rules/00-workflow.mdc`（SSOT，改完跑过 `make sync-claude-rules` 同步到 `.claude/rules/00-workflow.md`）里陈述的 URL 约定改成新方案；`request.spec.ts`、`notification.spec.ts` 的路径断言同步更新

## 明确排除

- 不处理生产部署本身（`script/admin.sh` 构建/打包/Supervisor 动作）——只改了仓库里的 nginx 配置文件，没有执行任何实际部署
- 不新增对旧 `/blog`、`/videos`、`/admin/blog` 等 URL 的重定向兼容层——未上线，无需兼容
- `.cursor/rules/21-public-pages.mdc` 未改——只按目录/组件名限定范围，不含字面量路径，不受这次改动影响

## 验证

1. `npm run typecheck`、`npm run test`
2. Dev server 手动复现原 bug 场景：登录后台，地址栏直接访问 `/bgg/admin/blog/article` 后 F5 硬刷新，确认停留在后台文章管理页而不是跳到公共页；`/bgg/admin/video/list` 同样验证
3. 公共页面免登录独立可访问：`/bgg/front/blog`、`/bgg/front/blog/:id`、`/bgg/front/videos`、`/bgg/front/videos/:id`
4. 侧边栏高亮、面包屑在新路径下正确（`AppSidebar.vue`/`utils/breadcrumb.ts` 是精确字符串匹配，DB 迁移 SQL 执行后菜单数据和前端路由前缀就能对上）
5. 登录/登出、403、404 兜底跳转、聊天未读消息点击跳转在新前缀下正常工作
6. DB 迁移 SQL 由用户执行后，联调确认菜单驱动的动态路由和数据库 `path` 一致
