# 响应式断点体系 + 公共页面完全重构（Phase 3 Week 6-7）

> 前置依赖：`05-design-system-and-tokens.md`（令牌化 + 暗色模式基建）已完成，本文档在此基础上做公共页面的结构性重构。用户明确要求："我期望是响应式布局，这样我在 web 端、移动端可以按企业级的处理方案进行处理，我支持你完全重构"——范围是**完全重构**，不是在现有"小程序风格"上小修小补。

## 1. 现状问题

- `.claude/rules/21-public-pages.md` 约束的"小程序风格"：暖色渐变背景（`linear-gradient(135deg, #fff7e6 0%, #ffe9d9 45%, #ffd1a4 100%)`）、卡片式列表/详情、固定 DOM 契约（`public-list-page`/`public-detail-page` 根类名 + `.container`/`.hero`/`.list-grid`/`.list-card` 等）。
- 断点处理是**单一 `@media (max-width: 768px)` hack**：移动端下卡片从网格切换成左图右文的行布局，分页组件精简为 `prev, pager, next`——这是"给桌面布局打个移动端补丁"的做法，不是从响应式原则出发的设计。
- `variables.scss` 已经定义了三级断点（`$screen-sm: 768px`/`$screen-md: 1024px`/`$screen-lg: 1280px`），但公共页面样式（`public-list.scss`/`public-detail.scss`，264/171 行）完全没有使用，是独立于全局令牌系统之外的一套自成一体的样式。
- 硬约束但仍然有效、需要保留的**交互契约**（与视觉风格无关，是产品逻辑）：`MetricReporter` 统一打点、`IcpFooter` 合规展示、列表→详情跳转的分页/搜索/滚动位置 `sessionStorage` 持久化与恢复。这些不因视觉重构而废弃。

## 2. 目标：企业级响应式方案

"企业级响应式"在本项目语境下具体指：

1. **断点体系统一**：公共页面改为使用 `variables.scss` 的 `$screen-sm/md/lg` 三级断点，而不是自造一个 768px 硬编码值；后台管理界面如果本身也有响应式需求（目前主要是桌面场景，但侧边栏折叠等已有基础响应式行为），一并对齐同一套断点变量，做到全站断点定义只有一处来源。
2. **布局策略按设备类型分层设计，不是简单媒体查询缩放**：
   - **桌面/Web 端（≥ `$screen-md`）**：多栏网格布局（列表页），详情页可以引入侧边栏（目录/相关推荐/作者信息，复用现有 `BlogTOC.vue`/`BlogAuthorCard.vue`/`BlogCategoryNav.vue`/`BlogSocialLinks.vue` 组件——这些组件此前是为"标准博客风格"方案设计但未启用，见 `docs/changelog/archive-frontend.md` §4 决策记录，本轮重构是复用它们的机会，不需要重新设计）。
   - **平板（`$screen-sm` ~ `$screen-md`）**：单栏但保留卡片式网格，密度介于桌面和移动之间。
   - **移动端（< `$screen-sm`）**：单栏堆叠，触控友好的点击区域尺寸，分页/搜索交互简化——这部分延续现有"移动端下卡片改行布局、分页精简"的**产品行为**（这个决策本身是对的，只是实现方式要从硬编码 media query 改成基于断点变量的系统化实现）。
3. **不是"小程序风格"专属**：新的响应式方案是公共页面统一的布局系统，视觉皮肤（配色、圆角、阴影）走 `05` 号文档的设计令牌，不再有公共页面自己的一套独立配色（渐变暖色系是否保留、替换成什么，属于视觉方向的具体取舍，按 `08-dev-execution-and-review-points.md` 的原则，**这类具体视觉细节必须先出预览/截图给用户确认**，不能在文档里直接拍板一个新配色方案）。

## 3. 与既有"标准博客风格"方案的关系

`docs/changelog/archive-frontend.md` §4 记录过一次"标准博客风格"（顶部导航栏+左侧分类导航+右侧 TOC 目录+阅读时间/字数统计+相邻文章导航）的详细设计评审，当时因为方向未定而搁置，对应组件 `BlogHeader.vue`/`BlogTOC.vue`/`BlogCategoryNav.vue`/`BlogAuthorCard.vue`/`BlogSocialLinks.vue` 已经存在于 `components/blog/` 但未被 `views/public/` 使用。

本轮"企业级响应式完全重构"与该方案的目标高度重合（顶部导航+侧边栏/TOC 正是典型的企业级/桌面优先布局形态），**执行时应该先确认这批闲置组件的实际完成度和可用性**（读一遍组件源码，不要假设文档里的描述就是当前代码状态），能复用的直接复用，避免重新发明。这不等于直接照搬旧方案的每一个细节——旧方案是纯桌面导航设计，本轮还需要补上移动端的对应形态，且是否采纳"标准博客风格"的具体视觉语言仍需按 §2 第 3 条走预览确认流程。

## 4. 保留不变的契约

以下来自 `.claude/rules/21-public-pages.md` 的产品逻辑约束继续有效，视觉重构不影响：

- `MetricReporter` 组件的打点调用方式（`<module>_list`/`<module>_detail` 命名约定），不允许改回 ad-hoc `metricApi.report` 调用。
- `IcpFooter` 必须在所有公共页面挂载。
- 列表→详情跳转时分页/搜索状态/滚动位置的 `sessionStorage` 持久化与恢复逻辑保留（`BlogList.vue`/`VideoList.vue` 现有实现模式），详情页返回优先 `router.back()`，回退到存储状态兜底的逻辑不变。

## 5. 执行顺序

1. 先在 `variables.scss` 基础上确定响应式断点的具体使用规则（哪些布局用 flex、哪些用 grid、断点切换点），产出一份轻量的"响应式设计规范"作为本次视觉实现的依据（可以直接作为代码注释/一个新的 `styles/responsive.scss` 存在，不需要单独再写一篇文档）。
2. 评估 `components/blog/` 下闲置组件的复用可行性（§3）。
3. 用 `frontend-design` skill 出视觉方向预览，按 `08` 号文档的流程找用户确认具体配色/布局细节。
4. 确认后重写 `views/public/**` 四个页面 + `public-list.scss`/`public-detail.scss`（大概率整体重写而不是修补，因为断点体系和布局策略都变了）。
5. 更新 `.claude/rules/21-public-pages.md`（详见 `09-rules-and-docs-sync-checklist.md`）反映新契约，废弃"小程序风格"专属描述。
6. 在 `05` 号文档暗色模式适配步骤里回归验证公共页面的暗色展示。

## 完成的定义

- 公共页面不再依赖硬编码 768px 断点，改用 `variables.scss` 断点变量。
- Web 端与移动端各自有明确设计的布局形态（不是同一套桌面布局简单缩放），且已经过用户对预览效果的确认。
- `MetricReporter`/`IcpFooter`/滚动位置恢复等产品级契约验证无回归。
- 亮色/暗色模式下公共页面均实测通过。
- `.claude/rules/21-public-pages.md` 已同步更新（`09` 号文档跟踪）。
