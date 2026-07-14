# admin-frontend

基于 Vue 3 + Vite 的前端项目，承担双重角色：`/admin/*` 后台管理界面（需登录）+ `/blog/*`、`/videos/*` 公共展示页面（免登录）。统一调用 `admin-server` 提供的接口。

## 技术栈

Vite 5 + Vue 3.4（Composition API）+ TypeScript 5.9 + Element Plus + Pinia + Axios + vue-i18n + dplayer（视频播放）+ vitest。**不是 Nuxt**——曾有过 Nuxt SSR 迁移实验，已完全回滚。

## 目录结构

```
src/
├── api/
│   ├── generated/   # goctl 从后端 .api 生成的接口代码，禁止手改
│   └── *.ts          # 对 generated 的二次封装（iam/system/monitoring/content/chat/sdk/task/misc/public），业务代码统一从这里导入
├── components/common/ # 通用组件（D2Table 等，见 components/common/README.md）
├── composables/        # useDictOptions、usePermission、useAppConfig 等
├── directives/permission.ts # v-permission 指令
├── stores/            # Pinia：app/dict/user/websocket
├── styles/             # 含公共页样式模板 public-list.scss / public-detail.scss
└── views/               # 按后端域分目录：iam/system/monitoring/misc/content/chat/sdk/task 后台页面 + public/ 公共页面
```

详细规范（API 分层、字典/权限用法、命名规则）见根目录 [`AGENTS.md`](../AGENTS.md) 与 [`.cursor/rules/20-frontend.mdc`](../.cursor/rules/20-frontend.mdc)、[`.cursor/rules/21-public-pages.mdc`](../.cursor/rules/21-public-pages.mdc)。

## 环境准备

- Node 18+
- 建议使用 pnpm（也可用 npm）

## 本地开发

```bash
pnpm install
pnpm dev
```

开发服务器需要 `admin-server` 已在 `localhost:20000` 运行（Vite 会把 `/api` 代理过去）。访问路径固定为 `/admin/`（见 `vite.config.ts` 的 `base` 配置）。

## 常用脚本

| 命令 | 说明 |
|---|---|
| `pnpm dev` | 本地开发服务器 |
| `pnpm build` | 生产构建，产出 `dist/` |
| `pnpm preview` | 预览构建产物 |
| `pnpm typecheck` | `vue-tsc --noEmit` 类型检查 |
| `pnpm lint` | ESLint 检查 |
| `pnpm test` | vitest 单测（`vitest run`） |
| `pnpm test:watch` | vitest watch 模式 |

**注意**：接口代码生成入口是 `admin-server/scripts/generate-ts.sh`（产物在 `src/api/generated/`，禁止手改）；`package.json` 曾有的 `api:gen` 死脚本（对应 `scripts/api-gen.mjs` 从未存在）已在 Phase 1 Week 2 清理移除。

## 新增页面/模块

标准 CRUD 页面不需要手写骨架：后端脚手架（`generate-sql.sh`）会连带生成基于 `D2Table` 的列表页骨架到 `src/views/temp/`，直接在骨架上补充业务字段即可，详见根目录 [`AGENTS.md`](../AGENTS.md) 第 2.1 节。表格组件用法见 [`src/components/common/README.md`](src/components/common/README.md)。

## 构建与部署

```bash
pnpm build                              # 产出 dist/
bash script/admin.sh package frontend   # 或走统一打包脚本
```

生产环境部署在 Nginx 的 `/admin/` 路径下（参考 [`config/nginxconfig.txt`](../config/nginxconfig.txt)）；静态资源的 tmpfs 缓存技巧见 [`docs/使用tmpfs（内存文件系统）缓存静态文件.md`](../docs/使用tmpfs（内存文件系统）缓存静态文件.md)。

## AI / MCP

本项目为前端开发配置了两个 MCP server（清单见仓库根 [`.mcp.json`](../.mcp.json)）：

| Server | 用途 |
|--------|------|
| `vue-lsp` | Vue/TS 语言服务：`definition` / `references` / `diagnostics` / `hover`（workspace 为本目录） |
| `frontend-ui` | 项目 UI 组件文档与约定：`ui_get_component`、`ui_list_components`、`ui_get_patterns` |

第三人上手：

```bash
pnpm install                    # 安装 vue-ts-lsp（vue-lsp 依赖）
# go install mcp-language-server  # 与 go-lsp 相同，若已装可跳过
make sync-claude-mcp-check      # 在仓库根目录验证 MCP 连接
```

维护 `frontend-ui` 源码后：`pnpm mcp:build`，并 commit `mcp/dist/`。

## 更多文档

- 根目录 [`AGENTS.md`](../AGENTS.md)、[`.cursor/rules/20-frontend.mdc`](../.cursor/rules/20-frontend.mdc)、[`.cursor/rules/21-public-pages.mdc`](../.cursor/rules/21-public-pages.mdc)：开发规范
- [`docs/前端开发进度.md`](../docs/前端开发进度.md)：已完成功能、技术决策记录、关键代码位置索引
