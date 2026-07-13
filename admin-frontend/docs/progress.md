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
