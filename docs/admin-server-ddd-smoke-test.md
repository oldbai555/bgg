# DDD-lite 重构后冒烟测试与前端联调指南

> 配合 `feature/bggadmin` 分支 DDD-lite 重构（commit `ddf3e61`）使用。  
> 重构**只改后端内部目录与 group 命名**，HTTP 路由路径未变；前端现有 `src/api/generated/` 与二次封装层**无需因重构而改 import 路径**。

---

## 一、前置条件

### 1.1 环境

| 项 | 要求 |
|----|------|
| Go | 1.24+，`go build ./...` 已通过 |
| goctl | 建议 1.9.2+（与 `routes.go` 头注释一致） |
| MySQL / Redis | 已初始化（`db/tables.sql` + `data.sql` + 业务增量 SQL） |
| 后端配置 | `etc/admin-api.yaml`；MySQL/Redis 见 `config/mysql.json.example`、`config/redis.json.example` 复制到 `/etc/work/` 或通过启动参数指定 |
| 前端 | Node 18+，`cd admin-frontend && npm install` |

### 1.2 启动命令

```bash
# 终端 1：后端（admin-server 目录）
go run admin.go -f etc/admin-api.yaml \
  -mysql-config /path/to/mysql.json \
  -redis-config /path/to/redis.json

# 终端 2：前端
cd admin-frontend && npm run dev
# 开发地址：http://localhost:5173/admin/
# API 代理：/api → http://localhost:20000
```

### 1.3 测试账号（`db/data.sql` 初始化）

| 用户名 | 密码 | 角色 |
|--------|------|------|
| `oldbai` | （见部署文档 / 本地 seed） | 超级管理员 |
| `admin` | `admin` | 业务管理员 |

---

## 二、冒烟测试清单

按 Phase 验收项组织。每项勾选 `[ ]` → `[x]` 表示通过。

### 2.0 编译与路由基线

- [ ] `cd admin-server && go build ./...` 无错误
- [ ] `internal/handler/routes.go` 中路由均为 `/api/v1/...`，无 `:id` 路径参数
- [ ] 维护导航自测：按 [`admin-server-维护导航.md`](admin-server-维护导航.md)「我要改 X」决策树，能在 30 秒内定位到目标文件（如改 RBAC → `domain/iam/permission_resolver.go`）

### 2.1 公开接口（无需登录）

| # | 操作 | 预期 | 通过 |
|---|------|------|------|
| 1 | `GET /api/v1/ping` | `code=0`，返回 pong | [ ] |
| 2 | `GET /api/v1/public/dict?code=xxx` | 返回字典项 | [ ] |
| 3 | `GET /api/v1/public/blog/articles` | 返回已发布文章列表 | [ ] |
| 4 | `GET /api/v1/public/videos` | 返回公开视频列表 | [ ] |

**curl 示例：**

```bash
curl -s http://localhost:20000/api/v1/ping | jq .
curl -s 'http://localhost:20000/api/v1/public/blog/articles?page=1&page_size=10' | jq .
```

**前端页面：**

- [ ] 打开 `http://localhost:5173/blog/` — 博客列表正常
- [ ] 打开 `http://localhost:5173/videos/` — 视频列表正常
- [ ] 无需登录，不出现 401 跳转

### 2.2 IAM：登录 / 鉴权 / 权限

| # | 操作 | 预期 | 通过 |
|---|------|------|------|
| 5 | `POST /api/v1/login` body `{"username":"admin","password":"admin"}` | 返回 `access_token` + `refresh_token` | [ ] |
| 6 | 带 `Authorization: Bearer <token>` 调 `GET /api/v1/profile` | 返回当前用户信息 | [ ] |
| 7 | 带 token 调 `GET /api/v1/users` | 有 `user:list` 权限则 `code=0` | [ ] |
| 8 | 带 token 调**无权限**接口（如无 `role:create` 的用户调 `POST /api/v1/roles`） | `code` 为权限拒绝码，非 500 | [ ] |
| 9 | 不带 token 调 `GET /api/v1/users` | 401 / 10003 | [ ] |
| 10 | `POST /api/v1/refresh` | 返回新 token 对 | [ ] |
| 11 | `POST /api/v1/logout` | 成功，后续旧 token 失效 | [ ] |

**curl 示例：**

```bash
TOKEN=$(curl -s -X POST http://localhost:20000/api/v1/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"admin","password":"admin"}' | jq -r '.data.access_token')

curl -s http://localhost:20000/api/v1/profile \
  -H "Authorization: Bearer $TOKEN" | jq .
```

**前端页面：**

- [ ] `http://localhost:5173/admin/login` 登录成功，跳转 Dashboard
- [ ] 侧边栏按角色动态加载菜单
- [ ] 无权限按钮被 `v-permission` 隐藏

### 2.3 Blog 域（ServiceContext 精简后重点）

| # | 操作 | 预期 | 通过 |
|---|------|------|------|
| 12 | `GET /api/v1/blog/articles`（带 token） | 后台文章列表 | [ ] |
| 13 | `POST /api/v1/blog/articles` 新建草稿 | `code=0` | [ ] |
| 14 | `POST /api/v1/blog/articles/publish` | 发布成功 | [ ] |
| 15 | `GET /api/v1/public/blog/articles/detail?id=...` | 公开页可见该文章 | [ ] |
| 16 | `GET /api/v1/blog/tags` | 标签列表 | [ ] |
| 17 | `PUT /api/v1/blog/articles` 编辑 | 更新成功 | [ ] |
| 18 | `DELETE /api/v1/blog/articles` | 软删除成功 | [ ] |

**前端页面：**

- [ ] `/admin/blog/articles` — 列表 / 新增 / 编辑 / 发布
- [ ] `/admin/blog/tags` — 标签 CRUD
- [ ] 发布后 `/blog/` 公开页能刷到新文章

### 2.4 任务调度（domain/task 迁移后）

| # | 操作 | 预期 | 通过 |
|---|------|------|------|
| 19 | 后台触发异步导出（如审计日志导出 `GET /api/v1/audit-logs/export`） | 返回 `task_id`，`admin_task` 表新增 pending 记录 | [ ] |
| 20 | 等待调度器执行（默认轮询间隔内） | 任务状态变为 success / failed | [ ] |
| 21 | `GET /api/v1/tasks/recent` | 浮球/任务列表可见该任务 | [ ] |
| 22 | WebSocket 任务通知（如已配置） | 前端 `TaskFloatBall` 收到状态更新 | [ ] |

**说明：** 调度器在 `ServiceContext` 启动时注册，执行器在 `internal/domain/task/executors/`。

### 2.5 各域抽样（确保 handler → logic 链路未断）

每个域至少测 1 个受保护接口：

| 域 | 抽样接口 | 通过 |
|----|----------|------|
| iam | `GET /api/v1/roles` | [ ] |
| video | `GET /api/v1/videos` | [ ] |
| chat | `GET /api/v1/chats` | [ ] |
| sdk | `GET /api/v1/sdk/api-keys` | [ ] |
| monitoring | `GET /api/v1/operation-logs` | [ ] |
| system | `GET /api/v1/dict-types` | [ ] |
| misc | `GET /api/v1/demos` | [ ] |

### 2.6 回归结论

- [ ] 以上全部通过 → DDD-lite 重构冒烟验收完成
- [ ] 有失败项 → 记录接口路径、响应 body、对应 `internal/handler/<domain>/` 与 `logic/<domain>/` 路径，对照维护导航排查

---

## 三、前端 `generate-ts.sh` 联调步骤

### 3.1 是否需要重新生成？

重构后已对比 goctl 输出（`goctl api ts -api admin.api`）与仓库内现有 `admin.ts`：

| 对比项 | 结论 |
|--------|------|
| `export function` 函数名 | **169 个，完全一致** |
| `/api/v1/...` 路径 | **168 条，完全一致** |
| 差异 | 仅函数**排序**不同 + goctl 版本注释行 |

**结论：不重新生成也能联调**；若 `admin.api` 有新增接口或类型变更，再执行生成即可。

### 3.2 生成步骤（须用户亲自执行）

```bash
cd admin-server/scripts
./generate-ts.sh
# 或指定 api 文件：
./generate-ts.sh admin.api
```

脚本行为：

1. 读取 `admin-server/api/admin.api`
2. 调用 `goctl api ts` 输出到 `admin-frontend/src/api/generated/`
3. 产出：`admin.ts`、`adminComponents.ts`、`gocliRequest.ts`

**禁止手改 `generated/` 目录。**

### 3.3 生成后检查

```bash
cd admin-frontend

# 1. 类型检查
npm run typecheck

# 2. 确认二次封装层 import 仍有效（函数名未变则无需改）
#    src/api/blog.ts、video.ts、metric.ts、public.ts
#    以及 views 里直接 import generated 的页面

# 3. 启动开发服
npm run dev
```

若 `typecheck` 报错，常见原因：

| 现象 | 处理 |
|------|------|
| 某 `XxxReq` 字段新增/删除 | 更新对应 `.vue` 表单字段，或改 `src/api/*.ts` 封装 |
| 返回类型从 `Response` 变 `null` | 封装层泛型调整，不影响运行时 |
| 全新模块接口 | 在 `src/api/` 新建二次封装文件，页面从封装层导入 |

### 3.4 二次封装约定（保持不变）

```
views/*.vue
    ↓ import
src/api/blog.ts 等（手写封装：统一错误处理、类型导出）
    ↓ import
src/api/generated/admin.ts（goctl 生成，禁止手改）
    ↓
src/api/generated/gocliRequest.ts → src/utils/request.ts（Axios + JWT）
```

- 开发环境 baseURL：`/api`（Vite 代理到 `localhost:20000`）
- 生产环境 baseURL：`/gateway/api`
- 时间字段：后端 `int64` 秒级时间戳，前端展示层格式化

### 3.5 标准联调流程（新功能 / 重构后验证）

按 [`00-workflow.mdc`](../.cursor/rules/00-workflow.mdc) 步骤 10–13：

1. 后端接口在 Swagger/Postman/curl 验证通过
2. **用户执行** `generate-ts.sh`
3. `npm run typecheck` + `npm run dev`
4. 后台管理页走查（登录 → 目标模块 CRUD）
5. 若涉及公开页，再测 `/blog/*`、`/videos/*`
6. 更新 `docs/后端开发进度.md` / `docs/前端开发进度.md`

---

## 四、遗留清理（2026-07-07 已完成）

Phase 3 迁移后残留的孤儿文件已删除（`routes.go` 未引用的旧命名 handler/logic）：

| 清理项 | 数量 | 说明 |
|--------|------|------|
| `internal/handler/task/*.go`（域根目录旧命名） | 3 | 保留 `task/task/` 子目录 |
| `internal/handler/video/*.go`（域根目录旧命名） | 4 | 保留 `video/video/` 子目录 |
| `internal/logic/task/*.go`（域根目录旧命名） | 3 | 保留 `task/task/`、`task/public/` |
| `internal/logic/sdk/*.go`（域根目录旧命名） | 10 | 保留 `sdk/sdk/`、`sdk/public/` |
| `internal/logic/video/*.go`（域根目录旧命名） | 4 | 保留 `video/video/` 等子目录 |
| `internal/task/` | — | 此前已迁移至 `internal/domain/task/`，目录不存在 |

**保留**：`internal/handler/chat/chatwshandler.go`（`custom_routes.go` WebSocket 路由依赖，非孤儿文件）。

**一次性迁移脚本**（`migrate-phase*.py`、`merge-phase*.py`、`fix-*.py`）已于 2026-07-07 删除；通用工具 `migrate-menu.sh` 保留。

---

## 五、快速验收命令汇总

```bash
# 后端编译
cd admin-server && go build ./...

# 公开接口
curl -s http://localhost:20000/api/v1/ping

# 登录 + 受保护接口
TOKEN=$(curl -s -X POST http://localhost:20000/api/v1/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"admin","password":"admin"}' | jq -r '.data.access_token')
curl -s "http://localhost:20000/api/v1/blog/articles?page=1&page_size=10" \
  -H "Authorization: Bearer $TOKEN"

# 前端类型检查（需先 npm install）
cd admin-frontend && npm run typecheck
```
