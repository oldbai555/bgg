# Admin Nuxt - 公共内容展示项目

基于 Nuxt 3 的公共内容展示系统，用于展示博客文章和视频内容。

## 技术栈

- **框架**：Nuxt 3 (Vue 3 + SSR)
- **语言**：TypeScript
- **包管理**：pnpm
- **UI 组件库**：Element Plus
- **状态管理**：Pinia
- **样式**：SCSS

## 快速开始

### 1. 安装依赖

```bash
cd admin-nuxt
pnpm install
```

### 2. 配置环境变量

复制 `.env.example` 为 `.env` 并修改配置：

```bash
cp .env.example .env
```

编辑 `.env` 文件，设置 API 基础地址：

```env
API_BASE_URL=http://localhost:20000
```

**注意**：
- 开发环境：如果后端 API 运行在 `http://localhost:20000`，直接使用即可
- 如果需要通过代理访问，可以配置为 `/api`（需要配置 Nuxt 代理）

### 3. 启动开发服务器

```bash
pnpm dev
```

开发服务器将在 `http://localhost:3000` 启动。

### 4. 访问页面

- 博客列表：http://localhost:3000/blog
- 博客详情：http://localhost:3000/blog/:id
- 视频列表：http://localhost:3000/videos
- 视频详情：http://localhost:3000/videos/:id

## 其他命令

### 构建生产版本

```bash
pnpm build
```

### 预览生产构建

```bash
pnpm preview
```

### 生成静态站点

```bash
pnpm generate
```

## 项目结构

```
admin-nuxt/
├── api/                    # API 封装
│   ├── blog.ts            # 博客 API
│   ├── video.ts           # 视频 API
│   ├── metric.ts          # 埋点 API
│   └── generated/         # 自动生成的 API 代码
├── assets/                # 静态资源
│   └── styles/           # 样式文件
├── components/            # 组件
│   ├── blog/             # 博客相关组件
│   └── common/           # 通用组件
├── composables/          # Composables
├── pages/                # 页面（文件系统路由）
│   ├── blog/             # 博客页面
│   └── videos/           # 视频页面
├── plugins/              # 插件
├── utils/                # 工具函数
├── nuxt.config.ts        # Nuxt 配置
└── package.json          # 项目配置
```

## 开发规范

### Nuxt 3 开发规范

1. **自动导入**：`useRouter`、`useRoute`、`useRuntimeConfig` 等 composables 会自动导入，无需手动导入
2. **路由导航**：使用 `NuxtLink` 替代 `router-link`
3. **API 请求**：优先使用 `$fetch`，在 composable 中使用 `useRuntimeConfig()`
4. **客户端组件**：使用 `<ClientOnly>` 包裹只在客户端运行的组件
5. **页面元数据**：使用 `definePageMeta()` 定义页面配置

详细规范请参考 `docs/前端nuxt开发进度.md`。

## 常见问题

### 1. API 请求失败

- 检查 `.env` 文件中的 `API_BASE_URL` 配置是否正确
- 确认后端服务是否已启动
- 检查浏览器控制台的错误信息

### 2. 样式不生效

- 确认 `nuxt.config.ts` 中的 CSS 配置正确
- 检查 SCSS 文件路径是否正确

### 3. 组件未自动导入

- 确认组件文件在 `components/` 目录下
- 检查组件命名是否符合规范（PascalCase）

## 相关文档

- [Nuxt 3 文档](https://nuxt.com/docs)
- [Vue 3 文档](https://vuejs.org/)
- [Element Plus 文档](https://element-plus.org/)
