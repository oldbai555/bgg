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
├── api/admin.api      # 唯一 .api 定义文件，所有路由的真源
├── cmd/adminseed/      # 管理员账号初始化工具
├── db/                 # tables.sql/data.sql（首次部署）+ migrations/（增量SQL）
├── etc/                # admin-api.yaml、middleware.yaml
├── internal/
│   ├── handler/ logic/  # goctl 生成骨架，按模块分子目录
│   ├── repository/       # 手写数据访问层（squirrel）
│   ├── model/             # goctl 生成
│   ├── middleware/        # 手写中间件
│   ├── consts/ types/ svc/ config/
├── pkg/                  # errs/response/jwt/cache/audit/monitor/...
├── scripts/              # generate-sql.sh / generate-model.sh / generate-api.sh / generate-ts.sh
└── .template/            # goctl 自定义模板
```

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

- 首次部署：执行 `db/tables.sql`（建表）+ `db/data.sql`（初始数据）
- 后续新增模块/字段：增量 SQL 放 `db/migrations/`（字典SQL → 业务表SQL → 权限SQL 顺序执行）

## 本地开发

推荐直接用 IDE（GoLand 等）运行 `admin.go`。也可以用命令行脚本：

```bash
bash script/admin.sh dev start   # 启动（带健康检查）
bash script/admin.sh dev status  # 查看状态
bash script/admin.sh dev logs    # 查看日志
bash script/admin.sh dev stop    # 停止
```

健康检查接口：`GET /api/v1/ping`

## 新增功能模块

项目内置工程化脚手架：`scripts/generate-sql.sh -group <group> -name <name>` 一条命令即可生成建表 SQL、RBAC 初始化数据（菜单/权限/接口）、`.api` 草稿、前端列表页骨架。完整流程见根目录 [`AGENTS.md`](../AGENTS.md) 第 2、2.1 节 或 [`.cursor/rules/00-workflow.mdc`](../.cursor/rules/00-workflow.mdc)，这里不重复。

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
- [`docs/后端开发进度.md`](../docs/后端开发进度.md)：已完成功能、技术决策记录、关键代码位置索引
