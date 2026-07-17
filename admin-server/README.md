# admin-server

基于 [go-zero](https://go-zero.dev/) 的后台管理系统后端，统一提供后台管理接口与公共（免登录）接口，供 `admin-frontend` 及其他前端消费。

## 技术栈

- Go 1.24 + go-zero
- MySQL 8.0+ / Redis 6.0+（数据访问层为 go-zero `sqlx` + cache，动态 SQL 统一用 [squirrel](https://github.com/Masterminds/squirrel)）
- JWT 双令牌（Access + Refresh）+ bcrypt
- gorilla/websocket（聊天、任务通知）

## 目录结构

```
admin-server/
├── api/admin.api      # 唯一 .api 定义文件，group 格式 <domain>/<module>
├── internal/
│   ├── handler/<domain>/<module>/  # goctl 生成骨架
│   ├── logic/<domain>/<module>/    # goctl 骨架 + 薄胶水（解析请求 -> 调对应 XxxRPC -> 映射响应）
│   ├── middleware/                 # 全部通过 zrpc client 调用各服务，不再直连数据库
│   ├── redisconn/                  # gateway 唯一保留的存储直连：共享 Redis（token 黑名单、限流滑动窗口）
│   ├── consts/ config/ types/ svc/ # svc.ServiceContext 只持有 Redis + 各服务 XxxRPC client 字段
├── pkg/ scripts/ .template/ db/
└── services/iam/, services/task/, services/sdk/, services/chat/, services/content/
    # Phase 2 拆出的 iam-rpc/task-rpc/sdk-rpc/chat-rpc/content-rpc 五个独立服务，全部领域代码/
    # repository/model/domain 都已搬出单体；gateway 的 internal/repository/、internal/model/、
    # internal/domain/ 三个目录已整体删除，不再直连任何 MySQL
```

详细维护导航见 [`docs/admin-server-维护导航.md`](../docs/admin-server-维护导航.md)。

详细代码规范（squirrel 用法、中间件顺序、命名规则等）见根目录 [`AGENTS.md`](../AGENTS.md) 与 [`.cursor/rules/10-go-code-style.mdc`](../.cursor/rules/10-go-code-style.mdc)。

## 环境准备

- Go 1.24+
- MySQL 8.0+、Redis 6.0+
- 安装 `goctl`：
  ```bash
  go install github.com/zeromicro/go-zero/tools/goctl@latest
  export PATH=$PATH:$(go env GOPATH)/bin
  ```

## 配置

- `etc/admin-api.yaml`：业务配置（监听端口 `20000`、JWT、Bcrypt、限流阈值）
- 外部 MySQL/Redis 配置（**必须存在**，否则服务无法启动）：
  - Linux：`/etc/work/mysql.json`、`/etc/work/redis.json`
  - Windows：`/c/work/mysql.json`、`/c/work/redis.json`
  - 本地可参考仓库根目录的 [`config/mysql.json.example`](../config/mysql.json.example)、[`config/redis.json.example`](../config/redis.json.example) 复制一份改成自己的连接信息；**真实配置文件不要提交到仓库**
- `etc/middleware.yaml`：限流等中间件配置，可选，不存在则回退使用 `admin-api.yaml` 中的配置

## 数据库初始化

- 首次部署：执行 `db/services/init-dev-db.sh`（按 `db/services/<service>/<module>/` 目录树的固定顺序建表 + 初始化数据，见 `docs/15-service-boundaries.md` 第 4 节）
- 后续新增模块/字段：增量 SQL 放对应模块的 `db/services/<service>/<module>/migrations/`（字典SQL → 业务表SQL → 权限SQL 顺序执行）

## 本地开发

推荐直接用 IDE（GoLand 等）运行 `admin.go`。也可以用命令行脚本（`script/admin.sh` 在仓库根目录，不是 `admin-server/scripts/`，以下命令需在**仓库根目录**执行）：

```bash
bash script/admin.sh dev start   # 启动（带健康检查）
bash script/admin.sh dev status  # 查看状态
bash script/admin.sh dev logs    # 查看日志
bash script/admin.sh dev stop    # 停止
```

健康检查接口：`GET /api/v1/ping`

## 新增功能模块

项目内置工程化脚手架：`scripts/generate-sql.sh -group <domain>/<module> -name <name>`（如 `-group iam/user -name 用户管理`）一条命令即可生成建表 SQL、RBAC 初始化数据（菜单/权限/接口）、`.api` 草稿、前端列表页骨架。完整流程见根目录 [`AGENTS.md`](../AGENTS.md) 第 2、2.1 节 或 [`.cursor/rules/00-workflow.mdc`](../.cursor/rules/00-workflow.mdc)，这里不重复。

## 构建与部署

```bash
bash script/admin.sh build server     # 构建
bash script/admin.sh package server   # 构建+打包
```

Supervisor 部署、生产配置细节见 [`script/README.md`](../script/README.md)。

## 管理员初始化工具

`cmd/adminseed`：独立命令行工具，用于初始化/重置管理员账号，按需运行。

## 更多文档

- 根目录 [`AGENTS.md`](../AGENTS.md)、[`.cursor/rules/*.mdc`](../.cursor/rules/)：开发规范与工作流
- [`scripts/README.md`](scripts/README.md)：代码生成脚本详细说明
- [`docs/changelog/`](../docs/changelog/README.md)：开发交接记录（历史/背景）；早期历史存档于 `docs/changelog/archive-backend.md`
