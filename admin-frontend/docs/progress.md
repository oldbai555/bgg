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
