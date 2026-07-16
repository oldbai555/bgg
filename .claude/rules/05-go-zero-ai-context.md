---
alwaysApply: false
paths: admin-server/**
---

# go-zero AI 工具链

本仓库通过子模块安装 [zeromicro/zero-skills](https://github.com/zeromicro/zero-skills)（`git submodule update --init .ai-context/zero-skills`）。

## 规则优先级（冲突时以此为准）

1. **本项目规则**：`.cursor/rules/*.mdc`、`AGENTS.md`（脚手架脚本、errs、squirrel、软删除等）
2. **go-zero 通用规则**：`.ai-context/zero-skills/`

通用 go-zero 指引不得覆盖本项目的 `generate-*.sh` 工作流、禁止手改生成目录、中间件顺序等约定。

## 工作流层（轻量，优先阅读）

| 文件 | 用途 |
|------|------|
| `.ai-context/zero-skills/SKILL.md` | 总览、何时加载各模块 |
| `.ai-context/zero-skills/getting-started/cursor-guide.md` | Cursor 环境 go-zero 指引 |
| `.ai-context/zero-skills/skill-patterns/generate-service.md` | 新建 API / RPC / Model 步骤 |
| `.ai-context/zero-skills/references/goctl-commands.md` | goctl 命令速查 |
| `.ai-context/zero-skills/references/rest-api-patterns.md` | 常用 REST 代码模式 |

## 知识层（需要细节时再读）

| 文件 | 用途 |
|------|------|
| `.ai-context/zero-skills/references/rest-api-patterns.md` | Handler → Logic → Model |
| `.ai-context/zero-skills/references/database-patterns.md` | SQL / Redis / 缓存 |
| `.ai-context/zero-skills/references/goctl-commands.md` | goctl 完整参考 |
| `.ai-context/zero-skills/references/resilience-patterns.md` | 熔断、限流 |
| `.ai-context/zero-skills/troubleshooting/common-issues.md` | 常见错误排查 |

## 本项目与通用 go-zero 的关键差异

- **代码生成**：用 `admin-server/scripts/generate-sql.sh` / `generate-model.sh` / `generate-api.sh`，不要直接 `goctl api go` 覆盖生成目录
- **脚本执行**：`generate-*.sh` 必须由用户亲自运行，AI 不得代执行
- **错误处理**：`pkg/errs`，不是 `errorx.NewCodeError`
- **SQL 构建**：Repository 层用 `squirrel`，不用字符串拼接
- **API 文件**：唯一入口 `admin-server/api/admin.api`，禁止路径参数 `:id`

## 环境

- `goctl` 已安装（`go install github.com/zeromicro/go-zero/tools/goctl@latest`）
- 生成后惯例：`go mod tidy` → 验证 import → `go build ./...`
