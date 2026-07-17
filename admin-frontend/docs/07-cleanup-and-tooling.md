# 死代码清理 + 工具链卫生（Phase 1 Week 2）

> 前置依赖：无（可独立执行，且大部分改动风险低，建议在 Phase 1 Week 2 与 `01`/`02` 号文档的改动一起提交，减少 PR 数量）。

## 1. 死代码清理清单

| 项目 | 结论 | 操作 |
|---|---|---|
| `views/temp/BlogList.vue` | 未接入路由，与 `views/content/BlogList.vue`（原 `views/blog/`，见 `02` 号文档）功能重复的早期脚手架产物 | 删除 |
| `views/temp/DailyShortSentenceList.vue` | 未接入路由 | 执行时确认 `daily_short_sentence` 功能是否还有产品需求；若无人在用，删除；若仍需要，从 `temp/` 移入正式域目录（`misc/` 或就近放 `content/`，视后端归属而定）而不是继续留在 `temp/` |
| `views/temp/DemoList.vue` | 明确是开发流程示例脚手架 | 删除，`docs/changelog/archive-frontend.md` §2 提到它"作为开发流程示例"的描述同步更新（见 `09` 号文档） |
| `views/temp/MetricList.vue` | 未接入路由，功能已被 `views/monitoring/MetricStats.vue`（原 `system/MetricStats.vue`）取代 | 删除前 diff 一下两个文件确认没有 `MetricStats.vue` 遗漏的独有功能，确认后删除 |
| `package.json` 的 `api:gen` script | 引用的 `scripts/api-gen.mjs` 不存在，已确认失效 | 删除该 script 条目；真正的生成入口是 `admin-server/scripts/generate-ts.sh`，`.claude/rules/20-frontend.md` 已经写清楚，不需要在前端仓库保留误导性的死脚本 |
| `.eslintrc.js` 里 "Vue 文件用 Prettier 格式化" 的注释及相关 `overrides`（`indent: 'off'` 等） | 死配置，`devDependencies` 未装 Prettier | 二选一，执行时定案：(a) 删除死配置注释和相关 `overrides`，明确本项目现状就是"ESLint 管一切，没有 Prettier"；(b) 补装 Prettier 让配置名副其实。**建议选 (a)**——现有 ESLint 规则（单引号/无分号/2 空格缩进）已经覆盖了 Prettier 通常负责的格式化职责，引入 Prettier 反而要处理两者规则冲突，价值不大，属于不必要的新增依赖 |
| `video.js`/`@types/video.js` 依赖 | 正文代码搜不到直接引用，疑似死依赖 | 执行时先 `grep -rn "video.js" src/` 排除 `md-editor-v3` 等间接依赖场景（`md-editor-v3` 不太可能依赖 `video.js`，但需要实际确认，不能凭猜测删除依赖），确认无引用后从 `package.json` 移除，同时确认 `VideoPlayer.vue` 实际使用的播放库是否是 `dplayer`（4 处引用）而 `video.js` 从一开始就是未采纳的备选方案 |

## 2. ESLint/TS 严格性复查

- `tsconfig.app.json` 已开启 `strict: true`，本轮不降低此设置；`01` 号文档的类型安全改造（消灭 `as any`）应该在严格模式下顺利通过类型检查，作为交叉验证。
- `@typescript-eslint/no-explicit-any` 目前是 `warn`，本轮全仓 `any` 使用量本就很低（`: any` 2 处、`as any` 5 处、`<any>` 1 处，其中 4 处 `as any` 在 `request.ts` 会被 `01` 号文档消灭）。清理完 `request.ts` 后回头检查剩余 `any` 使用点，若确认都已消灭或都是有正当理由的例外，**可以把该规则从 `warn` 升级为 `error`**，防止未来再引入不受控的 `any`；这是本轮清理工作的自然收尾动作，不需要单独立项。
- `vue-tsc` 版本目前是 `1.8.27`（对应 TS `5.3.3`），本轮不做主版本升级（升级 `vue-tsc`/`typescript`/`vite`/`element-plus` 等主版本是另一类风险更高的维护工作，不在"架构+视觉重构"范围内，除非升级是解决某个具体阻塞问题的必需前提）。

## 3. `package.json` 脚本卫生

清理完 `api:gen` 后，`scripts` 应保留：`dev`/`build`/`preview`/`typecheck`/`lint`/`mcp:build`（不变），新增 `test`/`test:watch`（见 `03` 号文档 vitest 部分）。执行时顺带确认 `mcp:build` 脚本对应的 `admin-frontend/mcp/` 子项目是否仍在使用，若也是历史遗留可以记录但不属于本轮清理范围（不要顺手清理未经核实的模块，扩大变更面）。

## 完成的定义

- `views/temp/` 目录不存在（或确认迁出的 `DailyShortSentenceList.vue` 已在正式域目录下）。
- `package.json` 无 `api:gen` 死脚本，`.eslintrc.js` 无 Prettier 死配置残留。
- `video.js`/`@types/video.js` 依赖去留已有明确结论并落地（保留则说明理由，删除则确认 `npm run build` 后无缺失依赖报错）。
- `npm run lint` + `npm run typecheck` + `npm run build` 全部通过。
