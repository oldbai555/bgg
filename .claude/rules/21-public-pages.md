---
alwaysApply: false
paths: admin-frontend/src/views/public/**,admin-frontend/src/components/blog/**,admin-frontend/src/components/common/PublicHeader.vue
---

# 适用范围

`admin-frontend/src/views/public/*`、`admin-frontend/src/components/blog/*`、`admin-frontend/src/components/common/PublicHeader.vue` 下的所有公共页面（博客列表/详情、视频列表/详情等）。

Phase 3 Week 7（`admin-frontend/docs/06-responsive-and-public-pages-redesign.md`）起，博客与视频统一成同一套"企业级响应式"视觉方案，不再各写一套配色/布局；之前"小程序风格"（暖色渐变 + 固定 768px hack）的描述已废弃，不要再参照。

# 统一布局（禁止各写一套，必须复用模板）

- 页面根元素必须同时包含业务类名（如 `blog-list-page`）和契约类名 `public-list-page`/`public-detail-page`
- 顶部统一挂载 `<PublicHeader />`（`components/common/PublicHeader.vue`，sticky 定位，含 logo + 博客/视频导航 + 社交图标），不要各页面自己写 header
- 列表页：`@import '@/styles/public-list.scss';`，HTML 结构固定为：
  - `.page-shell`（容器）→ `.page-intro`（`__title`/`__desc`/`__search` 可选）→ `.page-layout`（grid 容器）
  - `.page-layout` 内部 `.card-grid > .list-card`（`.cover` 可选 `.cover-fallback` / `.card-content` 内 `.card-title`/`.card-summary` 可选/`.card-meta` 可选/`.card-tags` 可选）+ `.pagination-bar`
  - 有分类导航等侧栏内容的页面（目前只有博客）把侧栏作为 `.page-layout` 的其它 grid 子元素，`grid-template-columns` 由页面自己在 `scoped style` 里定义（例：博客三栏 `200px 1fr 240px`，视频单栏 `1fr`），不要在共享文件里为不存在的数据（如视频没有分类）硬套侧栏
- 详情页：`@import '@/styles/public-detail.scss';`，HTML 结构固定为：
  - `.page-shell` → `.page-layout`（grid 容器）→ `.detail-card`（`.back-link` 可选 / `.title` / `.meta` / `.cover` 可选 / `.content`）
  - 侧栏内容（如博客的分类导航 + 目录）同样作为 `.page-layout` 的其它 grid 子元素，列数页面自定义
- 配色统一走 `src/styles/theme.scss` 的设计令牌（`--color-primary`/`--color-bg-card`/`--color-text-*` 等），**不允许**暖色渐变背景或页面私有配色；暗色模式通过 `[data-theme='dark']` 自动响应，无需页面自己适配
- 页面自定义样式只能在自己的业务类名（如 `.blog-list-page`/`.video-detail-page`）下做"小范围覆盖"（列数、页面独有的卡片细节如视频悬停播放），不得改写 `.list-card`/`.detail-card` 等共享基础样式的整体结构

# 响应式规范

- 断点统一使用 `src/styles/responsive.scss` 的 `mobile`/`tablet`/`tablet-up`/`desktop`/`wide` mixin，不要写字面量 `768px`/`1024px` 媒体查询；JS 侧断点判断统一读 `src/constants/breakpoints.ts` 的 `MOBILE_BREAKPOINT`
- 桌面/平板：`.page-layout` 多栏 grid（`.card-grid` 三列卡片网格 + 可选侧栏）；移动端：`.page-layout` 收窄为单栏，侧边分类导航横向滚动条化（`BlogCategoryNav.vue` 自身的 `@include mobile` 已内置该行为，不需要共享文件提供额外的侧栏包装类），`.card-grid` 单列，`.list-card` 切换为左图右文信息流（封面固定宽度，高度约 84px，摘要隐藏）
- 分页组件移动端简化为 `prev, pager, next`，隐藏总条数/每页条数/跳转输入框（`public-list.scss` 已内置）

# 交互与体验

- 列表页进入详情前，必须把分页参数、搜索条件、滚动位置写入 `sessionStorage`，返回时恢复（参考 `BlogList.vue`/`VideoList.vue` 实现）
- 详情页返回优先 `router.back()`，无历史记录时再按存储状态跳转带 query 的列表页
- 所有 Public 页面必须通过 `MetricReporter` 统一接入埋点上报（列表 `<module>_list`、详情 `<module>_detail`），不要各写一套 `metricApi.report`
- 所有 `views/public/*` 页面都必须在底部挂载 `IcpFooter` 组件，确保公网访问时备案信息统一展示
- 搜索框放在页面主体 `.page-intro__search` 里（随 `PublicHeader` 统一后，header 本身不再内置搜索框），列表页各自维护自己的搜索状态与路由同步
