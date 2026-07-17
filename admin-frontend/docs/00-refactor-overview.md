# admin-frontend 重构任务书总纲（Phase 1-3）

> 本文档是 `admin-frontend/docs/` 文档集（`00`~`09` + `progress.md`）的入口和索引，风格对齐 `admin-server/docs/00-refactor-overview.md`：可直接执行的任务说明，不是设计讨论稿。执行者（Cursor / Claude Code）在开始任何一个 Phase 之前，应先完整读一遍本文档，再按需跳转到具体的分篇文档。**本轮（文档产出会话）只产出这 11 篇文档，`admin-frontend/src/` 下没有任何代码改动**；从下一次会话开始按 Phase 1 Week 1 顺序实际动代码。

## 0. 这一批文档要解决什么

`admin-frontend`（Vue 3.4 + TS 5.3 + Vite 5 + Element Plus + Pinia，`postapocgame` 项目前端）在多轮迭代式开发之后，工程质量与已经完成三阶段重构的 `admin-server`（见 `admin-server/docs/00-refactor-overview.md`）明显不对称。`admin-server` 重构前有"零事务、DI 未落地、密钥硬编码、旗舰级 N+1 bug"这类结构性硬伤；`admin-frontend` 的问题量级不同——**没有同等的结构性硬伤，但域边界、API 层规范、类型安全、测试、视觉体系都有实打实的欠账**，且用户明确要求这次"改得彻底"，architecture 与视觉/UX 一起重做。

## 1. 现状核查（为什么要做这一轮）

以下问题是本轮规划前对代码库的实地核查结果（含文件路径/命中次数），不是推测：

| 问题 | 证据 |
|---|---|
| 规模 | 23,398 行 `.vue`/`.ts`，47 个视图，17 个组件——约为 `admin-server`（4.26 万行）的一半量级，views 数量不多但域内混杂严重 |
| 域边界错配 | `src/views/system/` 一个文件夹塞进了 iam（`UserList`/`RoleList`/`PermissionList`/`DepartmentList`/`MenuList`/`ApiList`）+ system 本域（`ConfigList`/`DictTypeList`/`DictItemList`/`FileList`/`NoticeList`/`NotificationList`）+ monitoring（`AuditLogList`/`LoginLogList`/`MonitorList`/`OperationLogList`/`PerformanceLogList`/`MetricStats`）+ task（`TaskList`，后端已是独立 `task-rpc`）四个后端域；`src/views/blog` 与 `src/views/video` 两个文件夹对应后端已合并的同一个 `content-rpc`（见 `admin-server/docs/15-service-boundaries.md`）。后端 9 个业务域（`iam/system/monitoring/misc/blog/video/chat/sdk/task`，与 `admin.api` 的 `group:` 声明一一对应）是前端应对齐的粒度 |
| API 层规范名存实亡 | `src/api/` 下只有 4 个手写 wrapper（`blog.ts`/`video.ts`/`metric.ts`/`public.ts`），iam/system/chat/sdk/task 五个域完全没有 wrapper；`grep "from '@/api/generated"` 命中 85 处，遍布 views/components，与 `.claude/rules/20-frontend.md`"业务代码只从二次封装层导入"的规则直接冲突（其中相当一部分是类型 import，需要在 Phase 1 逐一甄别） |
| 类型安全 | `src/utils/request.ts` 响应拦截器里 4 处 `as any` 裸转换 Envelope，全局请求层类型没有真正收敛，业务代码依赖拦截器"裸转换"后的 `data` |
| 状态管理 | `src/stores/websocket.ts` 428 行，身兼连接生命周期 + 未读消息列表 + 多个 getters，是"胖 store"候选 |
| 组件复用不均 | D2Table 命中 29/47 视图；`src/views/chatroom/ChatList.vue`（1037 行，全仓最大单文件组件）等仍手搓表格/列表逻辑 |
| composables/ 与 hooks/ 偶然分裂 | `composables/` 2 个文件（`useAppConfig.ts`/`useDictOptions.ts`），`hooks/` 仅 1 个文件（`usePermission.ts`），无文档化分工约定 |
| 路由脆弱点 | `src/router/index.ts` 的 `resolveComponent()` 用字符串比对 `import.meta.glob` 的 key 去匹配后端下发的菜单 `component` 字段，没有编译期校验，出错只 `console.error` 静默失败——`docs/changelog/archive-frontend.md` §3 "菜单组件路径与后端约定的自动映射策略（当前手工 map）"待办的具体代码位置 |
| `buildRoutesFromMenus` 内部重复 | `router/index.ts` 约 145-155 行与 175-185 行，"目录类型菜单"和"页面类型菜单"各自重复了一遍生成唯一 routeName 的 while 循环 |
| 死代码 | `src/views/temp/` 下 `BlogList.vue`/`DailyShortSentenceList.vue`/`DemoList.vue`/`MetricList.vue` 共 700 行，`router/index.ts` 对 `temp` 零命中，确认未接入路由的孤儿页面 |
| 工具链卫生 | `package.json` 的 `api:gen` 脚本引用的 `scripts/api-gen.mjs` 不存在（已确认失效）；`.eslintrc.js` 里"Vue 文件用 Prettier 格式化"的注释和相关 `overrides` 是死配置——`devDependencies` 根本没装 Prettier；`video.js`/`@types/video.js` 在正文代码搜不到直接引用，疑似死依赖（需排除 `md-editor-v3` 等间接依赖后再定论） |
| 测试 | 零测试：无 `*.spec.ts`/`*.test.ts`，`package.json` 无 vitest/jest/cypress/playwright |
| 设计令牌未落地 | `src/styles/variables.scss` 定义了完整的间距（`$spacing-xs/sm/md/lg/xl`）、圆角、阴影、断点（`$screen-sm/md/lg`）令牌，但 **0/47** 视图文件 `import` 它；**17/47** 视图在 `<style>` 块里硬编码十六进制色值（`#909399`/`#667eea`/`#606266`/`#f5f7fa`/`#409eff`/`#303133`/`#764ba2` 等高频出现） |
| 暗色模式形同虚设 | `src/styles/theme.scss` 有一套可用的亮/暗 CSS 自定义属性系统（`:root` + `[data-theme='dark']`），切换逻辑在 `src/stores/app.ts`，但只有 `.el-card`/`.el-input__wrapper` 等极少数选择器真正消费这些变量，绝大多数视图在暗色模式下观感不可控 |
| 公共页面非响应式 | `.claude/rules/21-public-pages.md` 约束的"小程序风格"（暖色渐变 + 卡片列表）只有单一 `@media (max-width: 768px)` 断点 hack，`variables.scss` 已定义的三级断点体系（`$screen-sm/md/lg`）完全没用上，不是真正的响应式设计 |
| 文档与配置不一致 | `docs/changelog/archive-frontend.md`（原 `docs/前端开发进度.md`）仍写"ESLint + Prettier"，但 Prettier 未安装；`api:gen` 脚本已死但文档/`package.json` 仍保留引用 |

这是用户自己的项目，**尚未上线，没有外部用户，不需要考虑兼容性**。目标是把前端做成"对标开源项目"的工程质量 + 与已完成重构的后端匹配的视觉/响应式体验，用户已明确表态"时间足够，可以放心大胆重构，要改得彻底"。

## 2. 已确认的范围与决策（执行时不用再问）

1. **域目录重组对齐后端 9 域**（`iam/system/monitoring/misc/blog+video→content/chat/sdk/task`），而不是 5 个部署服务粒度——部署粒度是后端的关注点，前端只是普通 SPA，按业务域组织代码即可。
2. **API 层**：为每个域补齐/新增手写 wrapper（`src/api/<domain>.ts`），视图一律通过 wrapper 调用后端接口，禁止直接 `import` `src/api/generated/` 里的请求函数；类型（`import type`）允许直接复用 generated 里的定义，不强制重新导出。
3. **引入 vitest**，核心逻辑（stores、composables、`request.ts` 拦截器、纯函数 utils）补测试；组件级测试按需，不设覆盖率门槛。
4. **后台管理界面视觉方向：设计令牌化 + 精修**（稳健路线）——消灭硬编码色值，`variables.scss`/`theme.scss` 强制落地，视觉上仍是 Element Plus 的自然演进，不做大幅重塑。
5. **公共页面视觉方向：完全重构**为响应式优先的企业级方案——统一断点体系（复用 `variables.scss` 已定义但未用的 `$screen-sm/md/lg`），替换现有单一 768px hack 与"小程序风格"DOM 契约；Web 端与移动端各自按企业级标准处理，不是简单媒体查询套壳。
6. **暗色模式：本轮做成全面支持**——后台管理界面 + 公共页面均需在暗色下正确展示，不是后台独占。
7. **项目未上线，无兼容性负担**，允许自由重命名目录/文件、调整路由结构、废弃现有视觉契约、重写 `.claude/rules`/`AGENTS.md` 对应章节。
8. **本轮全部纳入，无 descope 保留项**：域目录重组 + API wrapper 全覆盖、类型安全、composables/hooks 合并、路由映射修复、死代码/工具链卫生清理、vitest 测试基建、D2Table 复用收敛、大文件拆分、状态管理拆分、设计令牌落地、暗色模式全面适配、响应式断点体系、公共页面完全重构。
9. **规则/文档同步不是收尾才做的事**——每个 Phase 结束都要过一遍 `09-rules-and-docs-sync-checklist.md`，最后一个 Phase 做一次实质性重写收尾，与 `admin-server` 的 `13-rules-sync-checklist.md` 模式一致。
10. **时间线不设硬性 deadline**（用户明确"时间足够"），但仍按 Phase/Week 给出粗粒度节奏，便于分批执行和追踪进度，参考约 7 周。
11. **开发期执行策略比照 `admin-server/docs/10-dev-execution-and-review-points.md` 的先例**，本项目自己的边界见 `08-dev-execution-and-review-points.md`——域目录重组/API wrapper 新增/vitest 配置等可以直接执行、事后随 diff 走查；但视觉重塑的具体效果必须先出预览/截图，不能因为大方向定了就自行拍板全部视觉细节，`generate-*.sh` 系列脚本仍必须由用户亲自执行。

## 3. 三阶段结构与时间线

```
Phase 1  架构地基         Week 1-2   域目录重组、API wrapper 全覆盖、类型安全、composables/hooks 合并、
                                     router 组件映射修复、死代码/工具链卫生清理、vitest 基建 + 首批核心逻辑测试
Phase 2  组件与状态整改   Week 3-4   D2Table 复用收敛、大文件拆分（ChatList 等）、websocket store 拆分、
                                     测试覆盖补充、规则文档阶段性同步
Phase 3  视觉与响应式重构 Week 5-7   设计令牌落地、暗色模式全面适配、响应式断点体系、公共页面完全重构、
                                     规则文档最终同步重写（AGENTS.md/.cursor/rules 实质性更新）+ progress 收尾
```

**关键前提**：Phase 3 排在架构之后是有意为之——视觉/响应式重构依赖 Phase 1 已经清理好的域目录和组件边界，避免边做架构迁移边做视觉重做导致的双重返工。这与 `admin-server` "Phase 1 是 Phase 2 的地基，不是过渡阶段"的原则一致。

### Phase 1（Week 1-2，架构地基）

- **Week 1**：域目录重组（`views/system` 拆分 + `views/blog`+`views/video` 合并为 `views/content`）、8 个域 API wrapper 全覆盖、`request.ts` 类型安全改造（消灭 4 处 `as any`）、composables/hooks 合并。对应 `01`、`02`。
- **Week 2**：router `resolveComponent()` 修复、死代码清理（`views/temp/*`、`api:gen`、疑似死依赖）、ESLint/Prettier 配置卫生、vitest 基建搭建 + stores/composables/request 拦截器首批测试。对应 `02`、`03`（测试基建部分）、`07`。

### Phase 2（Week 3-4，组件与状态整改）

- **Week 3**：D2Table 复用收敛（逐视图核实该收敛还是合理例外）、`ChatList.vue` 拆分。对应 `04`。
- **Week 4**：`websocket.ts` store 拆分、其余 store 审计、测试覆盖补充到组件层、Phase 1-2 的规则文档阶段性同步。对应 `03`、`09`。

### Phase 3（Week 5-7，视觉与响应式重构）

- **Week 5**：设计令牌落地（`variables.scss`/`theme.scss` 强制引用、17 处硬编码色值清理）。对应 `05`。
- **Week 6**：暗色模式全面适配（后台 + 公共页）、响应式断点体系落地。对应 `05`、`06`。
- **Week 7**：公共页面完全重构（替换"小程序风格"契约）、`AGENTS.md`/`.cursor/rules` 实质性重写、`progress.md` 收尾。对应 `06`、`09`。

## 4. 如何使用这套文档

- 每份 Phase 文档（`01` 起）开头写"前置依赖"，结尾写"完成的定义"（`npm run typecheck` + `npm run lint` + `npm run build` 通过，涉及运行时/视觉行为的改动额外做人工冒烟或截图确认）。
- `01-architecture-target.md` 是 Phase 1 的技术总纲，`02`~`07` 是它的可执行拆解，遇到需要复核设计决策的地方回指 `01`，不要在分篇文档里重新推导一遍。
- 不要跨 Phase 并行改动；每个 Phase / 每个子任务结束都要跑上述验证命令。
- 任何不确定是否属于"可以直接做"还是"必须停下来问用户"的判断，先查 `08-dev-execution-and-review-points.md`，拿不准就问，不要逢清单必停也不要逢清单都不停。
- **`admin-frontend/docs/progress.md` 是本轮重构 Phase 1-3 期间的过程记录（阶段/周次/关键决策），只追加不重写**。仓库根目录曾有一份跨项目生命周期的功能级进度索引 `docs/前端开发进度.md`，2026-07-17 起已按「文档分层与生命周期」规则（`.cursor/rules/00-workflow.mdc`）退场并归档为 `docs/changelog/archive-frontend.md`；此后功能行为发生实质变化（不是纯内部重构）改为在 `docs/changelog/` 补一篇交接记录，不再更新已归档文件。
- 视觉重塑相关工作按 `.claude/rules/07-anthropic-skills.md` 的约定使用 `frontend-design` skill，避免"一看就是模板默认值"；前端改动完成后按项目规范启动 dev server 实际验证（`webapp-testing` skill），不能只凭类型检查/构建通过就声称完成。

## 5. 全部文档索引

| 文件 | 一句话用途 | 状态 |
|---|---|---|
| `00-refactor-overview.md` | 总纲，覆盖 Phase 1-3 全貌，本文档 | 已产出 |
| `01-architecture-target.md` | 技术决策正文：域目录规范、API wrapper 规范、类型安全、composables 合并、路由映射修复 | 已产出 |
| `02-domain-reorg-and-api-layer.md` | 域目录重组 + API wrapper 的可执行迁移清单 | 已产出 |
| `03-state-management-and-testing.md` | store 拆分方案 + vitest 引入与测试覆盖范围 | 已产出 |
| `04-component-library-refactor.md` | D2Table 复用收敛 + 大文件拆分方案 | 已产出 |
| `05-design-system-and-tokens.md` | 设计令牌落地规范 + 暗色模式全面适配方案 | 已产出 |
| `06-responsive-and-public-pages-redesign.md` | 响应式断点体系 + 公共页面完全重构方案 | 已产出 |
| `07-cleanup-and-tooling.md` | 死代码/依赖清理 + ESLint/Prettier/TS 严格性复查 | 已产出 |
| `08-dev-execution-and-review-points.md` | 开发期直接执行 vs 停下确认的边界 | 已产出 |
| `09-rules-and-docs-sync-checklist.md` | 规则/文档同步清单（`.cursor/rules`、`AGENTS.md`、脚手架模板等） | 已产出 |
| `progress.md` | 贯穿 Phase 1-3 的唯一进度记录，不分叉 | 已产出（种子条目） |

## 6. 明确不做的事

- 不引入 Nuxt/SSR——`docs/changelog/archive-frontend.md` §4 已有明确的失败先例（曾迁移到 admin-nuxt 后完全回滚），不重复尝试。
- 不做后台管理界面的大幅视觉重塑——已确认走"设计令牌化 + 精修"路线，仍在 Element Plus 体系内（见 §2 第 4 条）。
- 不引入除 vitest 外的额外测试框架（不上 Cypress/Playwright E2E，浏览器端验证走 `webapp-testing` skill 人工/半自动走查即可）。
- 不追求测试覆盖率百分比门槛，测试范围按 `03-state-management-and-testing.md` 的优先级来，不是"全面覆盖"。

## 7. 当前进度

本次交付：`00`~`09` 全部 10 篇任务书文档 + `progress.md` 种子条目，一次性完成，`admin-frontend/src/` 下没有任何代码改动。从下一次会话开始，按 Phase 1 Week 1（见第 3 节）实际动代码，每个子任务完成后回来更新 `progress.md`（追加条目，不要重写）。
