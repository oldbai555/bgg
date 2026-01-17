# 前端 Nuxt 开发进度

## 1. 项目概述

### 1.1 项目目标
将 `admin-frontend/src/views/public` 下的公共功能迁移到新的 Nuxt 3 项目 `admin-nuxt`，实现前后端分离的架构：
- **admin-frontend**：后台管理系统（Vue 3 + Vite）
- **admin-nuxt**：公共内容展示（Nuxt 3 + SSR）
- **admin-server**：统一后端 API（Go-zero）

### 1.2 技术栈
- **框架**：Nuxt 3（Vue 3 + SSR）
- **语言**：TypeScript
- **包管理**：pnpm
- **UI 组件库**：Element Plus
- **状态管理**：Pinia（如需要）
- **样式**：SCSS
- **Markdown 渲染**：md-editor-v3
- **视频播放器**：DPlayer

---

## 2. 迁移内容清单

### 2.1 页面迁移（4 个页面）
- [x] `BlogListPage.vue` → `pages/blog/index.vue`
- [x] `BlogDetailPage.vue` → `pages/blog/[id].vue`
- [x] `VideoListPage.vue` → `pages/videos/index.vue`
- [x] `VideoDetailPage.vue` → `pages/videos/[id].vue`

### 2.2 组件迁移（5 个组件）
- [x] `BlogHeader.vue` → `components/blog/BlogHeader.vue`
- [x] `BlogCategoryNav.vue` → `components/blog/BlogCategoryNav.vue`
- [x] `BlogAuthorCard.vue` → `components/blog/BlogAuthorCard.vue`
- [x] `BlogSocialLinks.vue` → `components/blog/BlogSocialLinks.vue`
- [x] `BlogTOC.vue` → `components/blog/BlogTOC.vue`

### 2.3 通用组件（2 个）
- [x] `MetricReporter.vue` → `components/common/MetricReporter.vue`
- [x] `IcpFooter.vue` → `components/common/IcpFooter.vue`

### 2.4 样式文件（3 个）
- [x] `public-list.scss` → `assets/styles/public-list.scss`
- [x] `public-detail.scss` → `assets/styles/public-detail.scss`
- [x] `blog.scss` → `assets/styles/blog.scss`

### 2.5 工具函数（2 个）
- [x] `date.ts` → `utils/date.ts`
- [x] `clipboard.ts` → `utils/clipboard.ts`

### 2.6 API 封装
- [x] `api/blog.ts`：博客相关 API
- [x] `api/video.ts`：视频相关 API
- [x] `api/metric.ts`：埋点上报 API

### 2.7 路由映射
- [x] `/public/blog` → `/blog`
- [x] `/public/blog/:id` → `/blog/:id`
- [x] `/public/videos` → `/videos`
- [x] `/public/videos/:id` → `/videos/:id`

---

## 3. 实施步骤

### 阶段一：项目初始化（✅ 已完成）

#### 1.1 创建 Nuxt 3 项目
- [x] 使用 `pnpm create nuxt-app admin-nuxt` 创建项目
- [x] 选择 TypeScript、SCSS、Element Plus 等配置
- [x] 配置 `package.json` 和依赖安装

#### 1.2 配置项目
- [x] 配置 `nuxt.config.ts`（Element Plus、Pinia、运行时配置等）
- [x] 创建项目目录结构（`pages/`、`components/`、`assets/`、`utils/`、`api/`）
- [x] 配置 TypeScript 和 ESLint

### 阶段二：修改 generate-ts.sh 脚本（✅ 已完成）

#### 2.1 脚本修改
- [x] 修改 `admin-server/scripts/generate-ts.sh`，支持 `--frontend admin-nuxt` 参数
- [x] 配置输出目录为 `admin-nuxt/api/generated`
- [x] 测试脚本生成功能

### 阶段三：公共组件和工具迁移（✅ 已完成）

#### 2.6 迁移通用组件
- [x] 迁移 `MetricReporter.vue`（埋点上报，使用 `<ClientOnly>` 包裹）
- [x] 迁移 `IcpFooter.vue`（备案信息）
- [x] 适配 Nuxt 3 的组件使用方式

#### 2.7 迁移工具函数
- [x] 迁移日期格式化工具（`utils/date.ts`）
- [x] 迁移剪贴板工具（`utils/clipboard.ts`，添加客户端环境检查）
- [x] 创建 HTTP 请求封装（`utils/request.ts`，提供工厂函数，符合 Nuxt 3 规范）
- [x] 创建 `composables/useApiRequest.ts`（使用 `$fetch` 和 `useRuntimeConfig`，符合 Nuxt 3 规范）

#### 2.8 迁移样式文件
- [x] 迁移 `public-list.scss` → `assets/styles/public-list.scss`
- [x] 迁移 `public-detail.scss` → `assets/styles/public-detail.scss`
- [x] 迁移 `blog.scss` → `assets/styles/blog.scss`
- [x] 在 `nuxt.config.ts` 中配置全局样式导入

### 阶段四：博客功能迁移（✅ 已完成）

#### 2.9 迁移博客组件
- [x] 迁移 `BlogHeader.vue`（适配路由路径：/public/blog → /blog，移除手动导入 useRouter）
- [x] 迁移 `BlogCategoryNav.vue`
- [x] 迁移 `BlogAuthorCard.vue`
- [x] 迁移 `BlogSocialLinks.vue`
- [x] 迁移 `BlogTOC.vue`（添加客户端环境检查）
- [x] 适配 Nuxt 3 的组件使用方式（使用自动导入）

#### 2.10 迁移博客页面
- [x] 迁移 `BlogListPage.vue` → `pages/blog/index.vue`
- [x] 迁移 `BlogDetailPage.vue` → `pages/blog/[id].vue`
- [x] 适配 Nuxt 3 的路由和页面生命周期（使用 NuxtLink、ClientOnly）
- [x] 使用 `definePageMeta()` 定义页面元数据
- [x] 移除手动导入 composables（useRouter、useRoute）
- [x] 适配样式导入路径（@/styles → @/assets/styles）

#### 2.11 博客 API 封装
- [x] 创建 `api/blog.ts`，封装博客相关 API 调用
- [x] 配置错误处理和拦截器（在 gocliRequest.ts 中使用 $fetch）

### 阶段五：视频功能迁移（✅ 已完成）

#### 2.12 迁移视频页面
- [x] 迁移 `VideoListPage.vue` → `pages/videos/index.vue`（适配路由路径：/public/videos → /videos）
- [x] 迁移 `VideoDetailPage.vue` → `pages/videos/[id].vue`（适配 DPlayer、useRuntimeConfig）
- [x] 适配 Nuxt 3 的路由和页面生命周期（使用 ClientOnly、NuxtLink）
- [x] 使用 `definePageMeta()` 定义页面元数据
- [x] 移除手动导入 composables（useRouter、useRoute）
- [x] 适配客户端环境检查（window、document、sessionStorage）

#### 2.13 视频 API 封装
- [x] 创建 `api/video.ts`，封装视频相关 API 调用
- [x] 配置错误处理和拦截器（在 gocliRequest.ts 中使用 $fetch）

### 阶段六：功能测试和优化（🔄 进行中）

#### 3.1 功能测试
- [x] 测试博客列表页（分页、搜索、标签筛选）
- [x] 测试博客详情页（内容渲染、目录导航、相邻文章）
- [x] 修复博客详情页 Markdown 内容显示问题（md-editor-v3 预览区域）
- [x] 优化滚动行为（中间内容区域可滚动，左右侧边栏固定）
- [x] 修复布局问题（右侧边栏宽度、对齐、裁剪问题）
- [x] 移除最外层滚动条（创建全局样式文件）
- [ ] 测试视频列表页（分页、搜索、预览播放）
- [ ] 测试视频详情页（播放器、磁力链接复制）
- [ ] 测试移动端适配
- [ ] 测试 SSR 渲染

#### 3.2 性能优化
- [ ] 优化图片加载（懒加载、WebP）
- [ ] 优化 API 请求（缓存、防抖）
- [ ] 优化首屏加载时间
- [ ] 优化 SEO（meta 标签、结构化数据）

#### 3.3 错误处理
- [ ] 统一错误处理机制
- [ ] 添加错误边界
- [ ] 优化错误提示

#### 3.4 样式和布局优化（✅ 已完成）
- [x] 修复博客详情页 Markdown 内容显示问题
  - 问题：md-editor-v3 预览区域宽度为 0，内容无法显示
  - 解决：修复 md-editor-content 的显示逻辑，确保预览区域可见且有宽度
  - 文件：`admin-nuxt/pages/blog/[id].vue`、`admin-nuxt/assets/styles/blog.scss`
- [x] 优化滚动行为
  - 中间文章内容区域（`.blog-main`）可滚动
  - 左右侧边栏使用 `position: sticky` 固定，滚动时保持固定
  - 移除预览元素的滚动条，由中间内容区域控制滚动
  - 文件：`admin-nuxt/assets/styles/blog.scss`、`admin-nuxt/pages/blog/[id].vue`
- [x] 修复布局问题
  - 修复右侧边栏宽度问题（添加 `min-width`、`max-width`、`flex-shrink: 0`）
  - 修复右侧边栏被裁剪问题（移除容器的 `overflow-x: hidden`）
  - 修复 blog-toc 组件被裁剪问题（调整宽度和 overflow 设置）
  - 修复左右侧边栏对齐问题（移除 blog-toc 的双重定位）
  - 文件：`admin-nuxt/assets/styles/blog.scss`、`admin-nuxt/components/blog/BlogTOC.vue`
- [x] 移除最外层滚动条
  - 创建全局样式文件 `admin-nuxt/assets/styles/global.scss`
  - 设置 `html`、`body`、`#__nuxt` 的 `overflow: hidden`
  - 所有页面使用 `height: 100vh` 替代 `min-height: 100vh`
  - 文件：`admin-nuxt/assets/styles/global.scss`、`admin-nuxt/nuxt.config.ts`、各页面文件
- [x] 修复博客详情页布局和滚动问题（2025-01-17）
  - 问题：`.detail-content` 内容过长时不会撑开高度，`.detail-navigation` 不会跟随在内容后方
  - 解决：
    1. 将 `.blog-main` 设置为滚动容器（`overflow-y: auto`），移除 flex 布局
    2. 移除 `.blog-detail-container` 的 flex 布局，让内容自然流动
    3. 移除 `.detail-content` 的滚动和 flex 设置，添加 `height: auto` 让内容自然撑开
    4. 修复 md-editor 相关容器的高度设置：
       - `.md-editor-content` 改为 `display: block`，添加 `height: auto`
       - `.md-editor-preview-wrapper` 移除 flex 设置，添加 `height: auto`
       - `.md-editor-preview` 添加 `height: auto`
  - 效果：`.blog-main` 作为滚动容器，内容过长时出现滚动条；`.detail-content` 内容可以自然撑开高度；`.detail-navigation` 跟随在内容后面
  - 文件：`admin-nuxt/assets/styles/blog.scss`、`admin-nuxt/pages/blog/[id].vue`

### 阶段七：清理 admin-frontend（待开始）

#### 4.1 删除已迁移代码
- [ ] 删除 `admin-frontend/src/views/public` 目录
- [ ] 删除 `admin-frontend/src/components/blog` 目录（如需要）
- [ ] 删除 `admin-frontend/src/styles/public-*.scss` 和 `blog.scss`
- [ ] 更新 `admin-frontend/src/router/index.ts`，移除公共路由

#### 4.2 更新文档
- [ ] 更新 `docs/前端开发进度.md`，标记迁移完成
- [ ] 更新项目 README，说明项目结构

---

## 4. 开发规范

### 4.1 Nuxt 3 开发规范（必须遵守）

#### 4.1.1 自动导入
- ✅ **不要手动导入 composables**：`useRouter`、`useRoute`、`useRuntimeConfig` 等会自动导入
- ✅ **不要手动导入 Vue API**：`ref`、`reactive`、`computed` 等会自动导入（但显式导入也可以）
- ✅ **组件自动导入**：`components/` 目录下的组件会自动导入

#### 4.1.2 路由和导航
- ✅ **使用 `NuxtLink`**：替代 `router-link`，Nuxt 3 会自动优化
- ✅ **使用 `useRouter()` 和 `useRoute()`**：自动导入，无需手动导入
- ✅ **使用 `definePageMeta()`**：定义页面元数据（layout、middleware 等）

#### 4.1.3 API 请求
- ✅ **优先使用 `$fetch`**：Nuxt 3 内置的 fetch 工具，支持 SSR
- ✅ **使用 `useFetch()` 或 `useAsyncData()`**：用于页面数据获取，自动处理加载状态和错误
- ✅ **使用 `useRuntimeConfig()`**：在 composable 或 setup 中获取配置，不能在模块顶层调用

#### 4.1.4 客户端/服务端区分
- ✅ **使用 `<ClientOnly>`**：包裹只在客户端运行的组件（如 MetricReporter、DPlayer）
- ✅ **检查浏览器环境**：使用 `typeof window !== 'undefined'` 检查客户端环境
- ✅ **使用 `onMounted()`**：确保在客户端执行的操作

#### 4.1.5 状态管理
- ✅ **使用 `useState()`**：用于跨组件共享状态（替代 `ref`）
- ✅ **使用 Pinia**：用于复杂的状态管理（已配置 `@pinia/nuxt`）

#### 4.1.6 代码规范
- 使用 TypeScript 严格模式
- 遵循 Vue 3 Composition API 最佳实践
- 使用 ESLint + Prettier（与 admin-frontend 保持一致）

---

## 5. 技术细节

### 5.1 Nuxt 3 配置要点

#### nuxt.config.ts 示例
```typescript
export default defineNuxtConfig({
  devtools: { enabled: true },
  typescript: {
    strict: true,
    typeCheck: true
  },
  css: [
    'element-plus/dist/index.css',
    '@/assets/styles/public-list.scss',
    '@/assets/styles/public-detail.scss',
    '@/assets/styles/blog.scss'
  ],
  modules: [
    '@pinia/nuxt',
    '@element-plus/nuxt'
  ],
  runtimeConfig: {
    public: {
      apiBase: process.env.API_BASE_URL || 'http://localhost:8888'
    }
  },
  app: {
    head: {
      title: '博客与视频',
      meta: [
        { charset: 'utf-8' },
        { name: 'viewport', content: 'width=device-width, initial-scale=1' }
      ]
    }
  }
})
```

### 5.2 API 请求封装

#### composables/useApiRequest.ts
使用 `$fetch` 和 `useRuntimeConfig`，符合 Nuxt 3 规范。

#### api/generated/gocliRequest.ts
已修改为使用 `$fetch`，并在 `getBaseURL()` 函数中安全地调用 `useRuntimeConfig()`。

### 5.3 路由配置

Nuxt 3 使用文件系统路由，无需手动配置：
- `pages/blog/index.vue` → `/blog`
- `pages/blog/[id].vue` → `/blog/:id`
- `pages/videos/index.vue` → `/videos`
- `pages/videos/[id].vue` → `/videos/:id`

### 5.4 状态管理

使用 Pinia 管理全局状态（如需要）：
- 用户信息（如果需要）
- 字典数据缓存
- 其他共享状态

---

## 6. 注意事项

### 6.1 与 admin-frontend 的差异

1. **路由路径**：`/public/blog` → `/blog`，`/public/videos` → `/videos`
2. **API 调用**：使用 `$fetch` 替代 `axios`
3. **组件导入**：使用自动导入，无需手动导入
4. **客户端组件**：使用 `<ClientOnly>` 包裹
5. **配置获取**：使用 `useRuntimeConfig()` 在 setup 中获取

### 6.2 开发注意事项

1. **SSR 兼容性**：确保所有代码在服务端和客户端都能正常运行
2. **环境变量**：使用 `runtimeConfig` 管理配置
3. **性能优化**：利用 Nuxt 3 的自动代码分割和预加载
4. **SEO 优化**：使用 `useHead()` 或 `definePageMeta()` 设置页面元数据

---

## 7. 测试清单

### 7.1 功能测试
- [ ] 博客列表页：分页、搜索、标签筛选、滚动位置恢复
- [ ] 博客详情页：内容渲染、目录导航、相邻文章、字数统计
- [ ] 视频列表页：分页、搜索、预览播放、滚动位置恢复
- [ ] 视频详情页：播放器、磁力链接复制、相邻视频
- [ ] 移动端适配：响应式布局、触摸交互

### 7.2 性能测试
- [ ] 首屏加载时间
- [ ] 页面切换速度
- [ ] API 请求响应时间
- [ ] 图片加载优化

### 7.3 SEO 测试
- [ ] 页面标题和描述
- [ ] 结构化数据
- [ ] 链接可访问性

---

## 8. 后续优化

### 8.1 功能增强
- [ ] 添加搜索功能（全文搜索）
- [ ] 添加评论功能
- [ ] 添加分享功能
- [ ] 添加收藏功能

### 8.2 性能优化
- [ ] 实现服务端缓存
- [ ] 实现客户端缓存
- [ ] 优化图片加载（WebP、懒加载）
- [ ] 优化代码分割

### 8.3 用户体验
- [ ] 添加加载动画
- [ ] 优化错误提示
- [ ] 添加骨架屏
- [ ] 优化移动端体验

---

## 9. 完成情况总结

### ✅ 已完成
- 项目初始化和配置
- TypeScript 代码生成脚本修改
- 公共组件和工具迁移
- 博客功能完整迁移（符合 Nuxt 3 规范）
- 视频功能完整迁移（符合 Nuxt 3 规范）
- 符合 Nuxt 3 开发规范的代码调整
- 博客详情页 Markdown 内容显示修复
- 滚动行为和布局优化
- 最外层滚动条移除

### 🔄 进行中
- 功能测试和优化（视频功能测试、移动端适配、SSR 测试）

### ⏳ 待开始
- 性能优化（图片加载、API 请求优化、SEO）
- 错误处理优化
- 清理 admin-frontend 中的已迁移代码
- 文档更新

---

**最后更新时间**：2025-01-17
**当前进度**：核心功能迁移完成，博客功能测试和优化完成 ✅（包括布局和滚动问题修复），视频功能测试进行中 🔄
