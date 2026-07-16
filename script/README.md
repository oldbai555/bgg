# Admin System 脚本工具

本目录包含 admin-system 项目的统一管理脚本。

## 📁 脚本文件

| 脚本 | 功能 | 说明 |
|------|------|------|
| `setup_ai_toolchain.sh` | AI 工具链初始化 | 新设备一键配置 Gentle-AI + CodeGraph + Engram，见 `docs/AI工具链上手.md` |
| `sync_claude_mcp.sh` | Claude Code MCP | `.mcp.json` 团队 SSOT、检查连接、从 Cursor 导入 |
| `engram_sync.sh` | Engram 记忆同步 | 跨设备导出/导入 `.engram/`，配合 `make engram-sync-*` |
| `admin.sh` | 统一管理脚本 | 开发、构建、打包、Supervisor管理 |
| `admin-completion.bash` | Bash 自动补全 | 启用后支持 Tab 键自动补全命令 |
| `utils.sh` | 工具函数库 | 通用函数（日志、路径、检查等），被 admin.sh 引用 |
| `ssh-setup.sh` | 一键配置 SSH 免密登录 | `./script/ssh-setup.sh <user> <ip> [alias] [port]`，生成密钥（如不存在）+ `ssh-copy-id` + 写入 `~/.ssh/config` 别名，用于快速接入部署服务器 |
| `deploy-dev.sh` | bgg-dev 后端一键部署 | **在服务器上**执行：`bash script/deploy-dev.sh [service...]`，git pull + docker compose build/up，见 `docs/changelog/` |
| `deploy-frontend.sh` | bgg-dev 前端一键部署 | **在本机**执行：`bash script/deploy-frontend.sh [ssh-host]`，本机构建打包 + 上传 + 远端解压替换 dist |

## 🚀 快速开始

### AI 工具链（Gentle-AI + CodeGraph）

换设备或新同事接入时（默认同时配置 **Cursor + Claude Code 插件**），在仓库根目录执行：

```bash
make setup-ai
```

详细说明见 [docs/AI工具链上手.md](../docs/AI工具链上手.md)。

### 启用自动补全（可选）

为了在使用脚本时能够按 Tab 键自动补全命令，可以启用 bash 自动补全：

```bash
# 临时启用（当前会话有效）
source script/admin-completion.bash

# 永久启用（添加到 ~/.bashrc）
echo "source $(pwd)/script/admin-completion.bash" >> ~/.bashrc
source ~/.bashrc
```

启用后，支持以下自动补全：

- `sh script/admin.sh bu` + Tab → `sh script/admin.sh build`
- `bash script/admin.sh dev s` + Tab → `bash script/admin.sh dev start`
- `./script/admin.sh pack` + Tab → `./script/admin.sh package`

### 开发环境

**注意：开发环境建议使用 GoLand 直接运行，脚本主要用于生产环境部署。**

如果需要在命令行启动服务：

```bash
# 启动后端服务（带健康检查）
bash script/admin.sh dev start

# 查看服务状态
bash script/admin.sh dev status

# 查看日志
bash script/admin.sh dev logs

# 停止服务
bash script/admin.sh dev stop
```

### 构建打包

```bash
# 构建后端
bash script/admin.sh build server

# 构建前端
bash script/admin.sh build frontend

# 打包后端（构建+打包）
bash script/admin.sh package server

# 打包前端（构建+打包）
bash script/admin.sh package frontend
```

### Supervisor 管理

```bash
# 生成 Supervisor 配置文件
bash script/admin.sh supervisor gen-conf

# 安装服务（构建+部署）
bash script/admin.sh supervisor install

# 部署服务
bash script/admin.sh supervisor deploy package.tar.gz

# 管理服务
bash script/admin.sh supervisor status    # 查看状态
bash script/admin.sh supervisor start     # 启动服务
bash script/admin.sh supervisor stop      # 停止服务
bash script/admin.sh supervisor restart   # 重启服务
bash script/admin.sh supervisor logs      # 查看日志
```

## 📖 详细说明

### admin.sh - 统一管理脚本

**开发命令：**
- `dev start` - 启动后端服务（带健康检查，确保服务完全启动）
- `dev stop` - 停止后端服务
- `dev status` - 查看服务状态
- `dev logs [行数]` - 查看日志（默认100行）

**构建命令：**
- `build server` - 构建后端服务
- `build frontend` - 构建前端项目
- `package server` - 打包后端（构建+打包）
- `package frontend` - 打包前端（构建+打包）

**Supervisor 命令：**
- `supervisor gen-conf` - 生成 Supervisor 配置文件
- `supervisor install` - 安装服务（构建+部署）
- `supervisor deploy <file>` - 部署打包文件
- `supervisor status` - 查看服务状态
- `supervisor start` - 启动服务
- `supervisor stop` - 停止服务
- `supervisor restart` - 重启服务
- `supervisor logs` - 查看服务日志

## ⚙️ 配置说明

### 服务启动检测

服务启动时会进行严格的健康检查：
1. 检查进程是否存在
2. 检查端口是否监听（默认示例为 20000，实际以配置为准）
3. 检查健康检查接口（`/api/v1/ping`）是否返回 200

只有所有检查通过，才认为服务启动成功。

### 配置文件

**后端配置：**
- 主配置文件：`admin-server/etc/admin-api.yaml`（业务配置：JWT、Bcrypt、BaseURL等）
- MySQL配置（必须存在，包含所有MySQL相关参数：连接信息、连接池参数等）：
  - Linux: `/etc/work/mysql.json`
  - Windows: `/c/work/mysql.json`
- Redis配置（必须存在，包含所有Redis相关参数：连接信息、超时参数等）：
  - Linux: `/etc/work/redis.json`
  - Windows: `/c/work/redis.json`
- 中间件配置：`admin-server/etc/middleware.yaml`（可选，不存在则使用主配置文件中的配置）

**Supervisor 配置：**
- 服务目录：`/home/work/service`（可通过环境变量 `SUPERVISOR_DIR` 配置）
- 日志目录：`/home/work/supervisor/logs`（可通过环境变量 `SUPERVISOR_LOG_DIR` 配置）
- 配置目录：`/etc/supervisor/conf.d`（可通过环境变量 `SUPERVISOR_CONF_DIR` 配置）

### 文件上传下载

文件上传下载通过 nginx 反向代理处理：
- 上传接口：`POST /files/upload` → nginx 代理到 `POST /api/v1/files/upload`
- 下载接口：`GET /files/download?id=xxx` → nginx 代理到 `GET /api/v1/files/download?id=xxx`
- 文件访问：`GET /files/uploads/xxx` → nginx 代理到 `GET /api/v1/uploads/xxx`

后端返回的路径格式：`/files/uploads/xxx` 或 `/files/download?id=xxx`

## 📝 使用示例

### 开发流程

```bash
# 1. 启动服务
bash script/admin.sh dev start

# 2. 查看日志
bash script/admin.sh dev logs

# 3. 停止服务
bash script/admin.sh dev stop
```

### 部署流程

```bash
# 1. 打包后端
bash script/admin.sh package server
# 输出: package/admin-server_abc123.tar.gz

# 2. 打包前端
bash script/admin.sh package frontend
# 输出: package/admin-frontend_abc123.tar.gz

# 3. 部署后端到 Supervisor
bash script/admin.sh supervisor deploy package/admin-server_abc123.tar.gz

# 4. 前端部署到 nginx（手动操作）
# 解压到 /home/work/web/dist 目录
```

### 一键安装

```bash
# 构建并部署后端服务
bash script/admin.sh supervisor install
```

## 🔧 前置要求

### 必需工具

- **Go** (1.19+) - 后端开发
- **Node.js** (18+) - 前端开发
- **pnpm** 或 **npm** - 前端包管理
- **Supervisor** - 进程管理（生产环境）

### 可选工具

- **curl** 或 **wget** - 健康检查（用于服务启动检测）
- **netstat** 或 **ss** - 端口检查（用于服务启动检测）

## 📝 注意事项

1. **配置文件路径**：
   - MySQL配置：`/etc/work/mysql.json`（必须存在）
   - Redis配置：`/etc/work/redis.json`（必须存在）
   - 如果文件不存在，服务启动会失败

2. **服务启动检测**：
   - 使用健康检查接口确保服务完全启动
   - 最多等待30秒，超时则启动失败

3. **打包文件**：
   - 后端打包：包含可执行文件和配置文件
   - 前端打包：仅包含构建产物（dist目录）

4. **Windows 环境**：
   - 需要使用 Git Bash 或 WSL 运行脚本
   - 或使用 PowerShell 调用 Git Bash：`& "C:\Program Files\Git\bin\bash.exe" script/admin.sh dev start`

## 🐛 故障排查

### 服务启动失败

```bash
# 查看详细日志
bash script/admin.sh dev logs 50

# 检查配置文件
cat admin-server/etc/admin-api.yaml
cat /etc/work/mysql.json
cat /etc/work/redis.json
```

### 健康检查失败

```bash
# 手动测试健康检查接口（请将端口替换为当前后端监听端口，示例使用 20000）
curl http://localhost:20000/api/v1/ping

# 检查端口是否监听
netstat -tuln | grep 20000
# 或
ss -tuln | grep 20000
```

### 配置文件不存在

```bash
# Linux 系统：检查配置文件
ls -l /etc/work/mysql.json
ls -l /etc/work/redis.json

# Windows 系统：检查配置文件
ls -l /c/work/mysql.json
ls -l /c/work/redis.json

# 如果不存在，需要创建配置文件
# 可以参考 config/mysql.json.example 和 config/redis.json.example
```

## 📚 相关文档

- [后端开发进度](../docs/后端开发进度.md)（已整合 go-zero 实现方案内容）
- [前端开发进度](../docs/前端开发进度.md)（已整合 Vue3 实现方案内容）
