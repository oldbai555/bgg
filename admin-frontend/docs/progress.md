# admin-frontend 重构过程记录（Phase 1-3）

> 本文档是本轮 `admin-frontend` 重构 Phase 1-3 期间的唯一过程记录（阶段/周次/关键决策），只追加不重写。与仓库根目录 `docs/前端开发进度.md`（跨越整个项目生命周期的功能级进度索引）是两份不同的记录，不要互相替代——分工原则见 `00-refactor-overview.md` 第 4 节。

---

## 2026-07-13：文档规划阶段完成

**What**：产出 `admin-frontend/docs/00~09.md` 全部 10 篇任务书文档，风格对齐 `admin-server/docs/00-refactor-overview.md` 的先例。`admin-frontend/src/` 下没有任何代码改动，本次交付是文档。

**核心决策**（详见各文档，此处只记摘要，不复制正文）：
- 域目录重组对齐后端 9 域（`iam/system/monitoring/misc/blog+video→content/chat/sdk/task`），而不是 5 个部署服务粒度。
- API 层：8 个域各建一个手写 wrapper，视图禁止直接 import `api/generated/` 里的请求函数。
- 引入 vitest，核心逻辑（stores/composables/request 拦截器/纯函数 utils）补测试，不设覆盖率门槛。
- 后台管理界面视觉：设计令牌化 + 精修（稳健路线，不做大幅重塑）。
- 公共页面视觉：完全重构为响应式优先的企业级方案，废弃"小程序风格"单一 768px 断点契约。
- 暗色模式：本轮做成全面支持（后台 + 公共页均需适配）。
- 规则/文档同步纳入每个 Phase 的收尾动作，不是收尾才做（见 `09-rules-and-docs-sync-checklist.md`），已额外核实脚手架模板 `list_page.vue.tpl` 存在与新 API 规范的直接冲突（第 47 行硬编码 import generated），需要在 Phase 1 尽早修模板，避免后续新生成代码持续违规。

**Why**：用户要求 admin-frontend 参照已完成三阶段重构的 admin-server，进行"架构 + 视觉/UX 一起重做"的大规模重构；项目未上线，无兼容性负担，用户明确表态"时间足够，可以放心大胆重构，要改得彻底"。

**下一步**：从下次会话开始按 `00-refactor-overview.md` 第 3 节的 Phase 1 Week 1 顺序实际动代码（域目录重组 + API wrapper 全覆盖），完成后回来本文档追加条目。

---

## 2026-07-13：Phase 1 Week 1 落地（域目录重组 + API wrapper 全覆盖 + 类型安全 + composables 合并）

**What**：

1. **域目录重组**（`git mv` 保留历史，按 `02-domain-reorg-and-api-layer.md` §1 逐文件迁移）：`views/system/` 拆分为 `views/iam/`（User/Role/Permission/Department/Menu/Api/Profile）+ `views/monitoring/`（AuditLog/LoginLog/OperationLog/PerformanceLog/Monitor/MetricStats）+ `views/task/`（TaskList），`views/system/` 只保留 Config/DictType/DictItem/File/Notice/Notification 6 个真正的 system 域页面；`views/blog/` + `views/video/` 合并进 `views/content/`；`views/chatroom/` 改名 `views/chat/`。`views/temp/*` 死代码留到 Week 2（`07-cleanup-and-tooling.md` 范围）未动。
2. **8 个域 API wrapper 全部建立**：新增 `src/api/iam.ts`、`system.ts`、`monitoring.ts`（吸收原 `metric.ts`）、`content.ts`（合并原 `blog.ts` + `video.ts`）、`chat.ts`、`sdk.ts`、`task.ts`、`misc.ts`（demo/dailyShortSentence 去留待 Week 2 按 07 号文档结论处理），删除 `blog.ts`/`video.ts`/`metric.ts`；`src/api/public.ts` 不变。全仓 `grep "from '@/api/generated"` 的函数调用形式引用已归零，只剩 `import type`。
3. **`request.ts` 类型安全改造**：新增 `src/types/envelope.ts`（`Envelope<T>` + `isEnvelope` 类型守卫），消灭原 4 处 `as any` 裸转换，运行时行为不变。
4. **`composables/`/`hooks/` 合并**：`usePermission.ts` 移入 `composables/`，删除 `hooks/` 目录，全仓 `from '@/hooks/'` 引用归零。
5. **`router/index.ts` 联动修复**：6 条手写 `/system/*` 路由的 `component` import 路径同步指向新目录；抽取 `generateUniqueRouteName()` 消灭 `buildRoutesFromMenus` 内部两处重复的 routeName 生成逻辑。
6. **菜单 `component` 字段修正 SQL 已生成、未执行**：`admin-server/db/services/iam/menu/migrations/fix_menu_component_path_20260713.sql`，基于本机开发库 `admin_menu` 表实际查询结果核对（23 条 `UPDATE ... WHERE component = '<旧值>'`），等待用户确认后执行；执行前旧菜单路径仍会 404，这是预期的中间态。

**验证**：`npm run build` 通过；`npm run lint` 无新增错误（现有 5 个 `no-extra-semi` 错误和若干 warning 经 diff 核实是 `BlogDetail.vue` 等文件的存量问题，与本次改动无关）。`npm run typecheck` 因环境问题无法运行（`vue-tsc@1.8.27` 与本机安装的 `typescript@5.9.3` 不兼容，`git stash` 验证过未改动的代码同样报错，是预置环境问题，不在本轮改动范围内，也不属于"允许直接决定的依赖升级"，已如实记录，未擅自升级 vue-tsc/typescript）。

**Why**：按 `00-refactor-overview.md` 第 3 节 Phase 1 Week 1 排期，域目录重组是后续所有 Phase 的地基（`01-architecture-target.md` 明确"Phase 1 是 Phase 2 的地基"）。

**已知问题 / 下一步**：
- `vue-tsc`/`typescript` 版本不兼容导致 `npm run typecheck` 不可用，需要用户决定是否在后续 Phase 里升级（`07-cleanup-and-tooling.md` 已声明本轮默认不做主版本升级，这是执行中新发现的阻塞点，按 `08-dev-execution-and-review-points.md` 第 3 条应停下问用户，此处如实记录待确认）。
- 菜单 `component` 字段修正 SQL 待用户审阅后在开发库执行，执行后需登录验证菜单可点击、无 `[Router] 无法解析组件` 报错。
- Week 1 剩余的死代码清理（`views/temp/*`、`api:gen`、Prettier 残留配置）、vitest 基建按排期属于 Week 2，未在本次一并处理。

---

## 2026-07-14：Phase 1 Week 2 落地（router 修复 + 死代码/工具链清理 + typecheck 环境修复 + D2Table 泛型化 + vitest 基建）

**What**：

1. **Week 1 提交**：域目录重组等改动正式 `git commit`（`a525d0f`）。提交过程中触发本仓库的 `gga`（Gentleman Guardian Angel）预提交 AI 审查两轮，第一轮修了真实存量 bug（`MessageNotification.vue` 的 `timeStr` 未定义引用、`readStatus` 语义颠倒——字典约定 1=未读/2=已读/0=全部，之前误用 0 当未读；`NoticeReader.vue` 公告类型硬编码改 `useDictOptions`；`user.ts` 消灭多余 `as any`；`NotificationList.vue` 移除无效 `have-add` prop）；第二轮 `gga` 继续拦在跨 Phase 的遗留问题上（公共博客页样式契约、多处下拉硬编码等，均属 Phase 2/3 范畴），与用户确认后用 `--no-verify` 跳过该轮。

2. **router `resolveComponent()` 修复**（对应 `01-architecture-target.md` A.5）：新增 `validateMenuComponents()`，登录后 `fetchMenus()` 成功时（仅 dev 环境）遍历一次菜单树，把 `component` 字段与 `import.meta.glob` 得到的 `viewModules` key 集合比对，不匹配项通过 `console.warn` + `ElMessage.warning` 提前暴露，不再是"点击菜单才发现 404"。

3. **死代码清理**：删除 `views/temp/BlogList.vue`/`DemoList.vue`/`MetricList.vue`（确认零引用）；`views/temp/DailyShortSentenceList.vue` 核实后不是死代码——它管理的"每日一言"数据正被 `Dashboard.vue` 实际消费，只是管理页面本身从未真正挪出临时目录，已 `git mv` 到 `views/misc/DailyShortSentenceList.vue` 并改用 `miscApi` 二次封装（新建 `views/misc/` 域目录），菜单 `component` 字段修正 SQL 见下；`api/misc.ts` 同步移除 `demo*` 死导出。`package.json` 删除失效的 `api:gen` 脚本，`README.md` 同步更新；确认 `video.js`/`@types/video.js` 全仓无引用（实际播放器是 `dplayer`），从依赖中移除。`.eslintrc.cjs` 删除"Vue 文件用 Prettier 格式化"的死配置注释和 `indent: 'off'`，但保留了 `overrides: [{files: ['*.vue']}]` 这个空壳——排查 `npm run lint` 一度报"No files matching the pattern .. were found" 才发现这个 `overrides.files` glob 本身是 `eslint .`（无 `--ext`）发现 `.vue` 文件的**唯一机制**，不是纯粹的死配置，删除 override 整体会连带删掉文件发现能力，只能删规则不能删结构。

4. **`vue-tsc`/`typescript` 环境修复升级为泛型组件改造的连锁工程**（范围超出原计划，过程见下"范围决策"）：
   - 升级 `typescript` → `^5.9.3`、`vue-tsc` → `^3.3.7`（1.8.27 起初升到 2.2.12 时能跑但把 `element-plus`/`@vueuse/core` 等三方库内部 `.d.ts` 也纳入检查，报出上百条与本项目代码无关的 `GlobalComponents` 类型错误；`tsconfig.app.json` 补上 `skipLibCheck: true` 后这类噪音清零，这是本项目此前遗漏的标准配置，不是本轮引入的新决策）。
   - `D2Table.vue` 改造成泛型组件（`<script setup lang="ts" generic="T">`）：`data: T[]`、`onclick-delete`/`onclick-update-row` 事件的 `row` 参数类型改为 `T`（原来统一是 `Record<string, unknown>`，导致 ~15 个业务页面各自的强类型 handler 如 `(row: UserItem) => void` 全部类型不匹配）；组件内部动态 `column.prop` 取值/抽屉表单仍用 `Record<string, any>` 处理，只在 `data`/两个事件的边界做类型收窄。踩坑记录：Props 若声明成独立 `interface`（无论加不加 `export`）在 vue-tsc 处理泛型 SFC 时都报错（`TS4025`/`TS1184`），最终改成内联在 `defineProps<{...}>()` 里才过；`Props` 泛型不加约束（`T extends Record<string, unknown>` 这个看起来自然的写法会导致所有具体 `XxxItem` 类型因为没有索引签名而报"不可赋值"，验证过 `interface Foo{a:number}` 确实不满足 `Record<string,unknown>`）。
   - 顺带修掉 `typecheck` 打开后暴露的全部预存在真实 bug（不是本轮改动引入，是环境修好后第一次被看见）：`D2TableElemType.EditSelect`/`EditInputNumber` 拼写错误（应为 `Select`/`Number`，`BlogTagList`/`VideoList`/`BlogFriendLinkList`/`BlogSocialInfoList`）；`NoticeList.vue` 的 `drawerColumns.options` 直接赋值成 `ComputedRef` 而不是 `.value`（真实运行时 bug——公告类型/状态下拉在编辑抽屉里此前实际渲染不出选项）；`MenuList.vue` 父级选择器的树过滤 `map` 缺少显式返回类型导致 `children` 可选/必填推断冲突；`Profile.vue`/`NotificationList.vue` 的可能 `undefined` 访问；`md-editor-v3` 动态 `import()` 的 `default` 兜底访问类型收紧；`SdkInterfaceList.vue` 创建接口缺 `apiCode` 字段——查了 `services/sdk/internal/logic/sdkinterfacecreatelogic.go` 确认后端**完全忽略**入参 `apiCode`（自己按 path+method 算），但 `admin.api` 的 `SdkInterfaceCreateReq.apiCode` 没打 `optional` 标签，前端传空串占位即可，**这是后端 `.api` 定义的遗留不一致，本轮未动后端，只记录**；`gocliRequest.ts` 的 `webapi` 缺 `options` 方法（`admin.api` 里 `metricReportOptions`/`videoCollectOptions` 两个探测接口调用 `webapi.options()`），照 `get/delete/put/post/patch` 的模式补上；`ChatGroupList.vue` 拼装 `availableUsers` 时假设 `UserItem` 有 `departmentName`/`roleNames` 字段，实际类型没有——**这两个字段一直是空的**（`|| ''`/`|| []` 兜底吞掉了），修类型没有改变行为，但记录下来这是"添加群成员"选人框里部门/角色一直显示空白的真实原因，后续要显示需要新增 join 查询，不在本轮范围；`ChatMessageList.vue` 更实质：`roomId`/`userId`/`toUserName` 都不是 `ChatMessageListReq`/`ChatMessageItem` 的真实字段（真实字段是 `chatId`，`ChatMessageItem` 只有 `fromUserId`/`fromUserName`，没有"接收人"概念），原代码靠 `as` 强转掩盖，这两个筛选字段+"接收用户"列此前对后端完全不生效（`toUserName` 在数据里永远是 `undefined`），已改成用真实的 `chatId` 筛选并删除"接收用户"这个从未真正工作过的列/筛选框。

5. **`npm run lint` 清零**：`BlogDetail.vue` 存量的 5 处 `no-extra-semi`（ASI 安全分号在 `semi: never` 规则下被判定为多余）用 scoped `eslint --fix` 修掉；移除 `.eslintrc.cjs` 的 Prettier 死配置后 `.vue` 文件首次被 `indent` 规则检查，新增 115 条 `indent` warning（不是 error，不阻塞 `npm run lint` 退出码）——**刻意没有跑全仓 `eslint --fix`**，因为会产出一次不可审查的大规模缩进重排 diff，留作后续独立清理项，不在本轮处理范围内标注在此处，避免遗忘。

6. **vitest 基建 + 首批核心逻辑测试**：`vite.config.ts` 的 `defineConfig` 改从 `vitest/config` 导入并扩展 `test` 字段（复用同一份配置文件，符合 `03-state-management-and-testing.md` 要求）；**环境选型偏离文档默认值**——`happy-dom` 在本机 `Node 26` 下 `window.localStorage` 取不到值（直接用 `new Window()` 单测验证过是 happy-dom 自身问题），换 `jsdom` 仍然复现，最终在 `src/test-setup.ts` 里用一个最小内存 `Storage` 实现兜底全局 `localStorage`，环境定为 `jsdom`（`03` 号文档"除非遇到兼容性问题优先用 happy-dom"的例外条款覆盖了这次选择）。为支持测试，顺手做了两处最小重构（不改变运行时行为）：`request.ts` 的响应/错误拦截器从内联匿名函数拆成具名导出的 `handleResponse`/`handleResponseError`/`isPublicPath`/`handleTokenExpired`；`router/index.ts` 的 `generateUniqueRouteName` 加 `export`。新增 41 个测试，覆盖 `stores/dict.ts`、`stores/user.ts`（含 TTL 缓存、login/logout 状态清理）、`composables/usePermission.ts`、`composables/useDictOptions.ts`、`utils/request.ts` 拦截器（Envelope 解包、10003 过期码分支、公共页面不跳登录）、`types/envelope.ts` 的 `isEnvelope`、`generateUniqueRouteName`，全部通过。

7. **提交前 `gga` 第二轮审查修复**（真问题 + 又一批跨文件预存在问题，用户确认后一并处理）：`gocliRequest.ts` 补充了"重新生成后需手动重新应用这个手工补丁"的注释，防止下次 `generate-ts.sh` 覆盖后忘记补回 `options` 方法；`ChatMessageList.vue` 的 `total.value = filteredList.length`（前端过滤后的页内数量）改成 `resp.total`（后端按 `chatId` 筛选后的真实总数），原写法会导致翻页错乱。以及一批状态下拉字典化：`SdkInterfaceList.vue`（`sdk_status`，与 DB `1启用2禁用` 一致）、`BlogTagList.vue`（`blog_tag_status`——建表注释里声明过但字典一直没建，新增 `dict_blog_tag_status_20260714.sql` 待执行，执行前走 `useDictOptions` 的 fallback 照常工作）、`BlogFriendLinkList.vue`/`BlogSocialInfoList.vue`（已有 `useDictOptions` 但 `#cell` 展示仍手写死三元表达式，改用 `getLabel`）改为真正读字典；`REQUIRED_DICT_CODES` 补上 `sdk_status`/`blog_tag_status`/`blog_friend_link_status`/`blog_social_info_status` 四个原来漏掉的 code。**`ApiList`/`RoleList`/`UserList`/`DictItemList`/`DictTypeList`/`FileList` 六个文件的 `status` 下拉刻意没有跟着字典化**——查了对应建表 SQL 确认 `admin_api`/`admin_role`/`admin_user`/`admin_dict_item`/`admin_dict_type`/`admin_file` 六张表的 `status` 列本身就是"1 启用/0 禁用"的原始 DB 布尔列，不是字典驱动的业务枚举，"字典 value 从 1 开始、0 表示不筛选"这条规则管的是 `admin_dict_item` 里的字典项数据，不覆盖实体自己的原始状态列；如果照 `gga` 的建议把这六处的 `0` 改掉，提交的其实是错误的 DB 值，等于把一个假问题改成真 bug——已在这六处代码旁加注释说明，不是遗漏。`ConfigList.vue` 的"刷新缓存"按钮删除：`admin.api` 里根本没有对应的刷新缓存接口，原代码只是弹一个假成功提示，属于误导用户且没有后端可对接，不是缺一行 API 调用能补齐的。

8. **`gga` 第三轮审查修复**（这轮抓到一个本轮之前就存在、真实会影响生产环境的 bug）：`utils/request.ts` 的 `isPublicPath()` 原来用 `window.location.pathname.startsWith('/blog')` 判断是否公共页面，但 `vite.config.ts` 的 `base` 是 `/admin/`，生产环境真实 URL 是 `/admin/blog/...`，这个判断在生产环境**永远返回 false**——意味着公共页游客遇到 10003 错误码时会被错误地强制跳转登录页，是一个此前从未被发现的真实缺陷（`04` 号 wrapper 层重构、`01` 号类型安全改造都没碰到这段逻辑，这次是 `gga` 第三轮审查连带指出的）。改成读 `router.currentRoute.value.path`（Vue Router 解析后天然不含 base 前缀，和 `router/index.ts` 里 `router.beforeEach` 自己判断 `isPublicPath` 用的 `to.path` 是同一套语义，两处现在一致了）。连带补了一个测试专门覆盖"路由 path 不含 /admin 前缀也能正确识别公共页"这个场景，防止回归。另外两处非阻塞死代码清理：`README.md` 目录结构小节仍写着 Week 1 已经删除的 `hooks/` 和已经不存在依赖的 `video.js`，同步更新；`DictTypeList.vue` 底部有一段"修复刷新缓存按钮样式"的 `:deep(.el-button--warning)` 死 CSS（该文件模板里根本没有这个按钮，是历史遗留的复制粘贴残留，与本次删除 `ConfigList.vue` 按钮无关），一并删除。

**范围决策**（对话中途和用户确认过，记录留痕）：vue-tsc 升级本身发现 D2Table 的 `Record<string, unknown>` 事件类型与业务页面强类型 handler 系统性不兼容（~15 处），这本是 `04-component-library-refactor.md`（Phase 2）的范畴；用户明确选择"现在就把 D2Table 改成泛型组件，一次性清零"而不是留到 Phase 2，所以本条目实际上提前完成了 Phase 2 的一部分工作，`04` 号文档执行时需要注意 D2Table 泛型化已经不是待办。另外 Week 1 的 gga 预提交审查跳过、菜单 SQL 执行方式（用户在结论确认阶段一度提出"委托你执行"，但 `mysql` MCP 配置了 `ALLOW_UPDATE_OPERATION=false` 的硬限制，且这类数据修正 SQL 按 `00-workflow.md`/`08-dev-execution-and-review-points.md` 明确不在开发期例外范围内，最终仍按"生成 SQL、由用户亲自执行"处理，两条 SQL 均未执行）。

**验证**：`npm run lint`（0 error，115 warning）、`npm run typecheck`（0 error）、`npm run build`（成功，既有的 chunk 过大警告与本轮无关）、`npm run test`（7 个测试文件、41 个用例全部通过）四项全绿。

**已知问题 / 下一步**：
- `fix_menu_daily_short_sentence_component_20260713.sql`、`dict_blog_tag_status_20260714.sql` 已由用户确认后在本机开发库（`root@127.0.0.1:3306/admin`，对应 `.vscode/launch.json` 的 `IAM_MYSQL_DSN`）执行并核实生效；执行方式是直接 `mysql` CLI，不是走 `mysql` MCP（MCP 的 `ALLOW_INSERT/UPDATE_OPERATION` 仍保持 `.mcp.json` 里团队共享的只读默认值未动，只在 `~/.config/bgg/mysql-mcp.env`（机器本地、未提交）补了连接信息，方便后续用 MCP 查数据）。
- `fix_menu_component_path_20260713.sql`（Week 1，域目录重组的 6 个域）仍待执行，执行前对应旧菜单路径仍是预期中的 404 中间态。
- `SdkInterfaceCreateReq.apiCode` 后端定义与实现不一致（后端忽略该字段但 `.api` 未标 `optional`）——本轮只在前端传空串绕过，未改后端 `.api`/重新生成，后续如需彻底修正需要走 `.api` 改动 + `generate-api.sh`（用户执行）流程。
- `ChatGroupList.vue` 选人框部门/角色列一直空白（`UserItem` 没有对应字段）：类型已修正为诚实反映现状，视觉/数据补全是独立的功能增强，不在本轮范围。
- `BlogDetail.vue` 仍用 `blog.scss` 的三栏布局（分类导航+正文+目录），未切到 `21-public-pages.mdc` 要求的 `public-detail.scss` 单栏卡片模板——核实过不是换 import 能解决的机械活，涉及布局取舍，按 `08-dev-execution-and-review-points.md` 的要求需要先出预览确认，留到 Phase 3（`06-responsive-and-public-pages-redesign.md`）一起做。
- `.vue` 文件新增的 115 条 ESLint `indent` warning 未做全仓格式化，留待后续单独批次处理（避免与本轮功能性改动混在一个大 diff 里）。
- D2Table 泛型化已提前完成，`04-component-library-refactor.md` 执行 Phase 2 时需要对照更新该文档的现状描述，避免重复规划已完成的工作。

---

## 2026-07-14：Phase 2 Week 3 落地（D2Table 复用收敛审计 + ChatList.vue 拆分）

**What**：

1. **D2Table 收敛审计**（对应 `04-component-library-refactor.md` §1）：逐一打开该文档列出的 6 处未使用 D2Table 的视图核实，结论与文档初筛一致，全部保留例外，在每个文件模板顶部补一行注释说明原因（此前均无说明，属于文档"完成的定义"里明确要求的一项）：`iam/DepartmentList.vue`/`iam/MenuList.vue`（树形数据，用 `el-tree`）、`public/BlogList.vue`/`public/VideoList.vue`（公共展示页，无权限/CRUD，不适用 D2Table；核实时发现 `BlogList.vue` 实际仍是 `blog.scss` 三栏布局、未切到 `21-public-pages.md` 的 `public-list.scss` 模板，且滚动位置恢复是死代码，与已知的 `BlogDetail.vue` 问题同源，已在该文件注释里补充说明并等 Phase 3 一起处理；`VideoList.vue` 已正确落地 `public-list-page` 模板）、`chat/ChatList.vue`（即时通讯 UI）；`monitoring/MonitorList.vue` 打开核实后确认是统计卡片 + 30 秒轮询的实时资源监控看板，不是分页列表，同样保留例外并补注释。未额外排查其余 11 个"未使用 D2Table"的文件（Login/Dashboard/Home/错误页/Profile/详情页/编辑器等）——这些本质上不是列表页，是否需要例外说明是自明的，不属于 `04` 号文档 §1 表格圈定的审计范围，加注释反而是噪音。

2. **`ChatList.vue` 拆分**（1033 行 → 238 行，对应 `04` 号文档 §2）：
   - `components/chat/ChatListItem.vue`：单条会话列表项渲染（头像/名称/群组标签/描述），含 `formatChatDesc` 纯格式化函数。
   - `components/chat/ChatMessageBubble.vue`：单条消息气泡渲染（头像/用户名/时间/文本或图片内容），含 `formatMessageContent`/`formatTime` 两个纯格式化函数。
   - `components/chat/ChatMessageInput.vue`：Emoji 选择器 + 图片上传 + 文本框 + 发送按钮整个输入区，内部自管理 emoji 分页状态和待发送图片暂存；发送动作通过 `onSendText`/`onSendImage` 两个回调 prop（而非普通 `emit`）交给调用方执行，因为需要 `await` 结果来决定何时清空本地输入框/暂存图片（与原实现的清空时机保持一致：文本框在发起请求前清空，暂存图片在请求成功后清空）。**顺带修复一个拆分前就存在的问题**：原 `ChatList.vue` 里 emoji 每行列数/行数写死为 `8`/`3`，注释却写"从字典获取"，但 `loadChatConfig` 只读了字典 `chat_config` 里的"聊天窗口消息数量"一项，"Emoji每行显示数量"/"Emoji显示行数"两个字典项（`db/services/iam/dict/init_dict.sql` 里已存在）从未被读取——运营在字典里改这两个值不会生效。本次拆分顺手把这两项也接进 `loadChatConfig`，通过 `emojiColsPerRow`/`emojiRows` 两个 prop（默认值 8/3 兜底）传给 `ChatMessageInput.vue`，使其名副其实地"从字典获取"，不算拆分范围外的行为改动（只是把注释已经声称的行为补齐）。
   - `composables/useChatList.ts`：会话/消息状态（`chats`/`messages`/`selectedChat` 等）、聊天配置加载、`chatList`/`chatMessageList`/`chatMessageSend` 三个 API 调用、WebSocket `lastMessage` 监听与去重合并逻辑，完全不碰 DOM；滚动到底部通过 `onMessagesChanged` 回调交还给页面组件（页面组件持有 `messageListRef`），保持"UI 渲染"与"数据/网络同步"关注点分离，这也是 `04` 号文档要求的拆分方向。
   - 拆分后 `ChatList.vue` 主文件只剩顶层布局 + 组件组合 + `onMounted` 启动流程，238 行，达成文档"≤300 行"的目标。
   - 拆分是纯重构，未改变任何可观察行为；`sendTextMessage`/`sendImageMessage` 在 composable 内部对"未选中聊天"补了一道防御性 `throw`（原实现里这个判断在 `handleSendMessage` 内联一次，拆分后表单组件已经在调用前判断过，这里是双保险，不会被正常路径触发）。

**验证**：`npm run typecheck`（0 error）、`npm run lint`（0 error，99 warning，比 Week 2 的 115 条更少——原 `ChatList.vue` 里未格式化的嵌套 SCSS 缩进问题随拆分自然消失）、`npm run build`（成功）、`npm run test`（7 个文件 42 个用例全部通过）四项全绿。**未能完成人工浏览器走查**：本环境没有可用的浏览器自动化工具（无 Playwright/DevTools 类 MCP），且本机 `admin-server` 需要 6 个进程（gateway + 5 个 RPC 服务）+ MySQL/Redis 才能跑通完整链路，未在本次会话拉起；只启动了 `vite` dev server 确认能正常编译服务（无运行时报错），未做真正的登录后收发消息/切换会话/未读角标人工验证。这是本条目区别于以往"验证"小节的地方，如实记录，不冒充已做浏览器验证。

**Why**：按 `00-refactor-overview.md` 第 3 节 Phase 2 Week 3 排期（`04-component-library-refactor.md`）推进；D2Table 泛型化已在 Week 2 提前完成，本轮只剩收敛审计和大文件拆分两项。

**已知问题 / 下一步**：
- **需要用户后续在本机拉起完整 admin-server 链路后，人工过一遍会话列表的核心交互**（收发文本/图片消息、未读消息、切换会话、WebSocket 断线重连提示），确认拆分后行为与拆分前一致，再视为这项任务真正完成——当前只有静态检查（类型/lint/构建/单测）通过，不构成运行时行为的证明。
- `layout/` 6 个组件（`04` 号文档 §3）本轮未动，按文档结论"组件边界合理，本轮不拆分"，样式层面的暗色模式/响应式适配留给 Phase 3。
- Week 4（`websocket.ts` store 拆分、其余 store 审计、组件层测试覆盖、Phase 1-2 规则文档阶段性同步）尚未开始。

---

## 2026-07-14：Phase 2 Week 4 落地（websocket store 拆分 + 组件层测试补齐 + 脚手架模板修复 + Phase 1-2 规则/文档同步）

**What**：

1. **`stores/websocket.ts` 拆分**（对应 `03-state-management-and-testing.md` §1）：读完全文 + `MessageNotification.vue`/`NotificationList.vue` 实际字段命名后确认走"合并成一个新 store"的方案——`MessageNotification.vue` 本就已经在组件内手写"合并 API 通知 + WebSocket 未读消息"的逻辑，说明这两个概念天然相关。拆分结果：`stores/websocket.ts`（429→244 行）只保留连接生命周期（`connect`/`disconnect`/重连/心跳）+ 原始消息广播（`lastMessage`）+ 任务浮球的 `recentTasks`/`scheduleRefreshRecentTasks`；新增 `stores/notification.ts`（190 行，setup store 写法）承接未读消息列表/已读状态（`unreadMessages`/`unreadCount`/`hasUnreadChat`/`markAsRead`/`markAllAsRead`/`clearReadMessages`）+ 聊天/任务/通知三类消息的处理逻辑（`handleChatMessage`/`handleTaskProgress`/`handleNotification`），通过 `watch(() => wsStore.lastMessage, ...)` 单向订阅前者广播的消息，不反向依赖，与 `01` 号文档"连接 store 广播、通知 store 消费"的设计一致。`MessageNotification.vue` 同步改用 `useNotificationStore()` 替换原来对 `wsStore.unreadMessages`/`markAsRead`/`markAllAsRead` 的直接访问，全仓核实无其它文件引用旧字段。`websocket.ts` 244 行略超文档"150-200 行"的目标区间，核实过剩余部分全是 `connect()` 内在的 WebSocket URL 拼装+重连逻辑本身的注释和 `onopen/onmessage/onerror/onclose` 四个回调，没有可再拆的独立职责，按文档"不要为了凑数字而强行拆分"的说明保留现状，不强拆。
2. **其余 store 审计**：复核 `app.ts`/`dict.ts`/`user.ts`，结论与 `03` 号文档既有审计一致（单一职责，不需要拆分），本轮未改动这三个文件本身。
3. **组件层测试覆盖补齐**（对应 `03` 号文档"完成的定义"：stores/composables 下所有文件都要有 spec）：新增 `stores/websocket.spec.ts`（9 用例，覆盖未登录跳过连接、连接幂等、disconnect 复位状态、sendMessage 两种分支、`handleMessage` 广播 + 任务浮球刷新防抖）、`stores/notification.spec.ts`（9 用例，覆盖未读消息增删/已读状态流转 + 订阅 `wsStore.lastMessage` 后聊天/任务/通知三类消息的入队规则）、`stores/app.spec.ts`（5 用例）、`composables/useAppConfig.spec.ts`（6 用例，字典配置缓存/失败兜底/清缓存）、`composables/useChatList.spec.ts`（10 用例，覆盖 `ChatList.vue` 拆分出来的会话加载/切换/发送文本图片/WebSocket 消息去重同步）。另给 `components/common/D2Table.vue` 补了 `D2Table.spec.ts`（3 用例，`@vue/test-utils` 浅层测试）——**过程中发现一个真实预存在 bug**：`el-pagination` 同时用 `v-model:current-page`/`v-model:page-size`（会触发我们自己的 computed setter 里的 `emit`）和显式 `@current-change`/`@size-change` 监听（又调用一次 `emit`），两条路径叠加导致 D2Table 每次翻页/改每页条数都会把 `current-change`/`size-change` 上抛两次，父页面因此每次分页操作都多发一次列表查询请求；已删除冗余的显式监听和对应的 `handleSizeChange`/`handleCurrentChange` 函数，只保留 `v-model` 一条路径，测试直接模拟 el-pagination 真实的双事件行为验证只上抛一次。这个 bug 一直存在（不是本轮改动引入），29/47 视图用了 D2Table，理论上都受影响，但实际后果只是多一次等效查询（同一页码/页大小），不是数据错误，用户侧基本不可感知。原计划测试"点击表格行内删除按钮"，实测中发现 el-table 在 jsdom 下 body 行渲染成占位空 `<tr>`（真实按钮位于 el-table 自己用于测量列宽的隐藏区域，不是真实交互行），是 Element Plus + jsdom 生态的已知限制，不是本项目代码问题；评估后放弃这条用例（按 `03` 号文档"组件测试按需，不强求"的定位，不为了凑测试硬堆一套自定义 el-table stub）。
4. **脚手架模板修复**（对应 `09-rules-and-docs-sync-checklist.md` 已核实的具体冲突，本应 Phase 1 域重组完成后立即改，本轮补做）：`admin-server/scripts/sqlgen/templates/list_page.vue.tpl` 第 51 行原来硬编码 `import {...} from '@/api/generated/admin'`，与新规则"视图禁止直接 import generated 函数"直接冲突。修复链路：`generate-sql.sh` 原来解析 `-group <domain>/<module>` 后只把 `<module>` 部分传给 Go 生成器（`GROUP="$MODULE"`），`<domain>` 只用于计算后端 `db/services/<service>/` 落地路径，从未传给模板；新增 `domain_to_frontend_api()` 映射（`blog`/`video`→`content`，其余域名不变，对应前端 8 个 API wrapper 的实际域名），新增 `-domain` 参数透传给 `scripts/sqlgen/main.go`（`Config`/`TemplateData` 各加一个 `Domain` 字段，`-domain` 与 `-group`/`-name` 一样是必填参数，缺失时报错退出，不做静默兜底）；模板改为 `import { {{.Domain}}Api } from '@/api/{{.Domain}}'`，四处 CRUD 调用点改为 `{{.Domain}}Api.{{.GroupFuncName}}List/Create/Update/Delete(...)`，与现有 8 个 wrapper 文件"具名导出 `<domain>Api` 对象，视图用 `iamApi.userList(...)` 这种命名空间调用"的实际约定完全一致；用独立的 `text/template` 脚本渲染验证过输出正确（未接入项目真实生成流程，`generate-sql.sh` 仍需用户亲自执行）。顺带核实了 `init_module.api.tpl` 的 `@server(group: {{.Group}})` 只输出裸模块名、不含 domain 前缀——确认这不是"隐含了旧前端目录分组假设"（`09` 号文档 §1 关注的问题），只是设计上把 `<domain>/<module>` 拼接留给人工在把 `.api.temp` 追加进 `admin.api` 时手动补，按文档"若无强绑定则不用改"的结论未动这个模板。
5. **Phase 1-2 规则/文档阶段性同步**（对应 `09-rules-and-docs-sync-checklist.md` §1"Phase 1 结束后先改一版"，本轮一次性补上 Phase 1+2 的量）：发现 `.claude/rules/20-frontend.md` 是 `.cursor/rules/20-frontend.mdc`（SSOT）经 `make sync-claude-rules` 生成的派生文件，直接改前者会在下次同步时被覆盖——先改了 SSOT（`.cursor/rules/20-frontend.mdc`）的目录结构/API 层规范/组件与状态管理/命名规范四节（域目录重组、composables/hooks 合并、"禁止直接 import generated 函数 + import type 例外"从建议改成强制规则、D2Table 泛型化说明、websocket/notification 两个 store 的职责边界），再跑 `make sync-claude-rules` 重新生成 `.claude/rules/20-frontend.md`，`go run ./script/sync_claude_rules.go --check` 确认两端零漂移。同步过程顺带发现 `.claude/rules/00-workflow.md`/`07-anthropic-skills.md`/`10-go-code-style.md` 三个文件此前有过不经过 SSOT 的手工改动（内部交叉引用被改成了 `.md` 后缀，但 `.cursor/rules/*.mdc` 源文件仍是 `.mdc` 后缀），本次全量重新生成顺带纠正了这个漂移，回到"引用其它规则文件一律用 `.mdc` 后缀"的 SSOT 原文状态，不属于本轮议题但如实记录。`AGENTS.md` §4 前端规范同步补了域目录重组、API 层强制规则、D2Table 泛型化、websocket/notification 拆分的对应描述，并删除了已经过时的"`api:gen` 脚本已失效"提法（脚本本体已在 Week 2 被删除，不只是失效）。`docs/前端开发进度.md`：§0 目录结构要点、§1 核心功能索引、§5 关键代码位置三节涉及的具体文件路径全部按重组后的新路径更新（含 `hooks/usePermission.ts`→`composables/usePermission.ts`、`views/blog+video`→`views/content`、`views/chatroom`→`views/chat`、`stores/websocket.ts`→拆分后的两个 store、技术栈里的 `video.js`→`dplayer`）；§2 "demo 功能页面" 标注为已废弃（不删除历史记录，只加状态标注）；§4 技术决策记录追加一条本轮重构中发现并修复的真实功能行为 bug 汇总（公共页生产环境错误跳转登录、`NoticeList` 编辑抽屉下拉不渲染、`ChatMessageList` 从未生效的接收人筛选、D2Table 分页双重上抛），只记功能行为发生实质变化的部分，不搬 `progress.md` 的全部内容，遵循两份文档的既有分工。
6. **提交前 `gga` 审查修复**（两轮）：第一轮——`useAppConfig.ts` 一直直接 `import {dictGet} from '@/api/generated/admin'`（预存在写法，不是本轮引入），与本轮刚写进规则的"禁止视图/业务代码直接 import generated 函数"冲突——`api/system.ts` 早已具名导出同一个 `dictGet`，顺手改成 `systemApi.dictGet(...)`，`useAppConfig.spec.ts` 同步改 mock `@/api/system`。`MessageNotification.vue` 的 `handleClearRead` 只清了后端已读消息，没有像 `handleMarkAllAsRead` 那样同步调 `notificationStore.clearReadMessages()`——这个 store 方法拆分前后都存在但从未被任何交互调用过（预存在遗漏，`websocket.ts` 拆分只是把这个死代码原样搬过去），现已补上调用，两个操作行为对称。同文件 SCSS 的 `&__header { &-title {} &-actions {} }` 嵌套编译出 `__header-title`/`__header-actions`，但模板用的类名是 `__title`/`__actions`，标题加粗和操作区 flex 布局从未真正生效（预存在样式 bug）——已把这两段样式提到 `&__header` 同级，改成 `&__title`/`&__actions`，与模板类名对齐。`notification.ts` 新写时沿用了旧 `websocket.ts` 的分号风格，被 `gga` 指出与项目 `semi: never` 约定不符，`npx eslint --fix` 清理。第二轮——`gga` 继续指出本轮touch 过的 `useAppConfig.ts`/`websocket.ts` 全文仍是旧的带分号写法（`websocket.ts` 本轮做了实质性拆分改写，不能算"没碰过"），两个文件补跑 `npx eslint --fix` 清理；`list_page.vue.tpl` 第 169 行 `onMounted(loadData);` 带分号，新生成的骨架从第一天就偏离规范，删掉分号；`MessageNotification.vue` 的 `onMounted` 里 `setInterval(loadNotifications, 30000)`（预存在代码，本轮未碰这段但因为改了同文件而被一并审查到）从未保存 timer id、也没有 `onUnmounted` 清理，组件卸载后仍会继续轮询——虽然这个组件实际挂在 `AppHeader.vue` 里、登录期间不会真正卸载，实践影响很小，但仍是真实的资源泄漏写法，已补 `onUnmounted` + `clearInterval`。

**关于 ESLint 覆盖范围的发现（未处理，记录留痕）**：排查"为什么 `notification.ts` 的分号从没在 `npm run lint` 里报出来"时发现 `npm run lint`（即 `eslint .`，无 `--ext`）实际上**只发现 `.vue` 文件**——`.eslintrc.cjs` 里 `overrides: [{files: ['*.vue']}]` 是唯一让 ESLint 在无 `--ext` 参数时发现 `.vue` 文件的机制（Week 2 已记录过这一点），但这也意味着所有裸 `.ts` 文件（`stores/`、`composables/`、`api/*.ts`、`utils/*.ts` 等）**从未被 `npm run lint` 实际检查过**，包括这些文件里大量的分号使用（`websocket.ts`/`user.ts`/`dict.ts`/`app.ts`/`useAppConfig.ts` 等全部有分号，直接 `npx eslint <文件>` 才会报出来）。此前所有 Week 报告里"`npm run lint` 0 error，N warning"的结论**只覆盖了 `.vue` 文件**，不是全仓结论，需要澄清但本轮不修——修复方式（给 `lint` 脚本加 `--ext .ts,.vue` 或调整文件发现机制）会一次性暴露全仓 `.ts` 文件的分号/其它规则违规，是一次独立的、需要专门 diff 审阅的清理批次，不适合夹在本轮改动里顺手做，留给后续 `07-cleanup-and-tooling.md` 相关工作或专门的一次性清理提交。

**验证**：`npm run typecheck`（0 error）、`npm run lint`（0 error，99 warning，与改动前持平；如上所述这个数字只覆盖 `.vue` 文件）、`npm run build`（成功）、`npm run test`（13 个测试文件、84 个用例全部通过，含本轮新增 42 个）四项全绿；`admin-server/scripts/sqlgen`（`go build -o sqlgen main.go`）编译通过；`go run ./script/sync_claude_rules.go --check` 确认规则文件两端同步；本地 `gga` pre-commit 审查两轮（第一轮拦下第 6 条列出的 4 个问题，修复后第二轮通过）。脚手架模板改动未通过真实 `generate-sql.sh` 走一遍完整生成流程验证（该脚本必须由用户亲自执行），只用独立 `text/template` 渲染验证过输出的 import/调用语句正确。

**Why**：按 `00-refactor-overview.md` 第 3 节 Phase 2 Week 4 排期收尾 Phase 2；`03`/`09` 号文档均明确"完成的定义"包含测试覆盖和规则文档同步这两项，不是可选项。

**已知问题 / 下一步**：
- **需要用户后续找机会实际执行一次 `generate-sql.sh` 验证脚手架模板改动**：确认生成出的 `<GroupUpper>List.vue` 确实 `import` 了正确的域 wrapper 且能正常编译/运行，本轮只做了离线模板渲染验证，未走完整生成流程。
- Phase 2（Week 3-4）全部完成，`04`/`03` 号文档规划的工作项已清零；下一步进入 Phase 3（Week 5-7，视觉与响应式重构），起点是 `05-design-system-and-tokens.md` 的设计令牌落地。
- D2Table 的分页双重上抛 bug 修复后，理论上所有用 D2Table 的视图分页请求次数会减半（从每次操作 2 次变 1 次），这是行为改进但不是本轮范围内需要额外验证的独立任务，测试已覆盖回归。
- `views/temp/` 目录本身在 Week 2 死代码清理后已清空但目录还在磁盘上（git 不追踪空目录，无需处理）。
- **`npm run lint` 实际只检查 `.vue` 文件，裸 `.ts` 文件全仓未被检查**（本条目"关于 ESLint 覆盖范围的发现"小节已详述），需要用户决定何时安排一次独立的 `.ts` 文件规则清理批次，不建议顺带处理。
