# 设计令牌落地 + 暗色模式全面适配（Phase 3 Week 5-6）

> 前置依赖：Phase 1-2（架构地基 + 组件整改）已完成，视图目录结构和组件边界已经稳定，本文档的改动才不会和还在迁移中的文件冲突。视觉改动按 `.claude/rules/07-anthropic-skills.md` 约定使用 `frontend-design` skill；按 `00` 号文档 §2 第 4 条，后台管理界面本轮走"设计令牌化 + 精修"路线，不做大幅视觉重塑。

## 1. 现状

- `src/styles/variables.scss`（43 行）：完整的间距（`$spacing-xs/sm/md/lg/xl`：4/8/16/24/32px）、圆角（sm/md/lg：4/8/12px）、阴影（sm/md/lg）、断点（`$screen-sm: 768px`/`$screen-md: 1024px`/`$screen-lg: 1280px`）令牌，**0/47 视图文件 `import` 它**。
- `src/styles/theme.scss`（64 行）：`:root` + `[data-theme='dark']` 的 CSS 自定义属性系统，覆盖 primary/success/warning/danger/info/背景/文字/卡片等颜色语义，切换逻辑在 `stores/app.ts`；目前只有 `.el-card`、`.el-input__wrapper` 等极少数选择器真正消费。
- 17/47 视图在 `<style>` 块内硬编码十六进制色值，高频出现：`#909399`、`#667eea`、`#606266`、`#f5f7fa`、`#409eff`、`#303133`、`#764ba2`，以及公共页面渐变色 `#fff7e6`/`#ffe9d9`/`#ffd1a4`。
- 全仓 46 处内联 `style="..."` 属性，未逐一核实是否为魔法数字（执行时需要抽查分类：哪些是动态计算值必须内联，哪些是应该挪进 scoped style/令牌的静态值）。

## 2. 设计令牌落地规范

1. **强制规则**：所有新增/改动的 `<style>` 块，颜色一律引用 `theme.scss` 的 CSS 自定义属性（`var(--el-color-primary)` 或项目自定义的 `var(--xxx)`，视 `theme.scss` 实际变量名而定，执行时先读一遍该文件全部变量名再写规则文档，不要在这里编造变量名），间距/圆角/阴影一律引用 `variables.scss` 的 SCSS 变量，禁止新写硬编码色值/魔法数字。
2. **存量清理**：17 处硬编码色值逐一替换为对应令牌——替换前先建立"硬编码色值 → 语义令牌"的映射表（比如 `#409eff` 大概率就是 Element Plus 的 primary 色，`#f5f7fa` 大概率是背景色），不能盲目替换成视觉不对应的令牌，替换后每个文件都要过一遍浏览器视觉核对（`webapp-testing` skill），避免"类型对了颜色错了"。
3. **落地方式**：不要求每个视图文件顶部显式 `@import 'variables.scss'`（SCSS 变量在 Vite + sass-embedded 配置下可以通过全局注入方式提供，检查 `vite.config.ts` 的 `css.preprocessorOptions.scss.additionalData` 配置是否已经全局注入 `variables.scss`；如果没有，本轮顺带加上，这样所有 `.vue` 文件的 `<style>` 块都能直接用 `$spacing-md` 而不用逐个 `@import`）。CSS 自定义属性（`theme.scss` 的变量）本身是全局的，不需要 import。

## 3. 暗色模式全面适配方案

范围：后台管理界面全部视图 + 公共页面（`views/public/**`）——用户已明确要求"全面支持"，不是后台独占。

### 执行步骤

1. **审计**：`grep -rn "data-theme" src/` 结合 `theme.scss` 现有覆盖的选择器列表，确认目前真正响应暗色切换的 CSS 规则清单（预计只有 `.el-card`/`.el-input__wrapper` 等个位数）。
2. **按视图批量核查**：完成 §2 的令牌替换后，硬编码色值已经归零，此时大部分视图应该"顺带"就能在暗色下正确展示（因为颜色都走了会响应 `data-theme` 的 CSS 变量）——这是为什么令牌落地要排在暗色适配之前的原因，两者不是独立工作量,而是令牌化直接消解了暗色适配的大部分工作。
3. **剩余人工核查点**：Element Plus 组件本身的暗色主题需要确认是否已经通过 `theme.scss` 正确覆盖（Element Plus 官方有自己的暗色模式变量体系 `--el-*`，需要核实 `theme.scss` 的自定义变量和 Element Plus 官方暗色变量之间的映射关系是否完整，而不是只覆盖了两三个组件）；图片/图标类资源（如有明确设计成浅色背景专用的图标）需要暗色下的替代方案或滤镜处理。
4. **公共页面暗色适配**：`06` 号文档会把公共页面从"小程序风格"重构为响应式方案，暗色适配应该在那次重构里一并设计（新样式体系从一开始就用令牌，而不是先给旧的小程序风格暖色渐变做暗色适配再推倒重来）——**执行顺序：先做 `06` 号文档的公共页面重构，再回来确认暗色模式在新样式下工作正常**，本文档列出这个依赖关系但具体重构方案见 `06`。
5. **验收**：每个视图在亮色/暗色下各截一次图（或用 `webapp-testing` skill 实际切换 `data-theme` 走查），确认文字对比度、边框、卡片背景等无"亮色的白底白字"级别的可用性问题。

## 完成的定义

- 全仓硬编码十六进制色值归零（`grep -rn "#[0-9a-fA-F]\{3,6\}" src/views src/components` 只剩必要的品牌色定义源头，即 `theme.scss` 本身）。
- `variables.scss` 的间距/圆角/阴影令牌全局可用，新代码不再需要手写 `@import`。
- 后台全部视图 + 公共页面在 `data-theme="dark"` 下人工走查通过，无明显可用性问题。
- `npm run typecheck` + `npm run build` 通过；视觉改动额外要求浏览器实测（不能仅凭构建通过声称完成）。
