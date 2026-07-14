# AI 工具链上手（Gentle-AI + CodeGraph）

本文面向**换设备继续开发**或**第三位维护者**快速接入本项目的 AI 辅助开发环境。

## 本项目的日常开发方式

团队默认使用 **Cursor IDE + Cursor 内的 Claude Code 插件** 双通道开发：

| 通道 | 用途 | 规则来源 | MCP 配置 |
|------|------|----------|----------|
| **Cursor**（Chat / Agent / Composer） | 日常对话、改代码、子 Agent | `.cursor/rules/*.mdc`（SSOT） | 仓库内 **`.cursor/mcp.json`**（已提交 git，团队 SSOT，唯一手改入口） |
| **Claude Code 插件**（Cursor 侧边栏） | SDD 工作流、Task 子代理、长任务 | `.claude/rules/*.md`（由 Cursor 规则同步） | 仓库根 **`.mcp.json`**（已提交 git，由 `.cursor/mcp.json` 全量生成，禁止手改） |

`.cursor/mcp.json` 是唯一权威来源，按当前项目实际需要精简（本项目未启用 `mongodb`；`redis` 已启用，见下文）；`.mcp.json` 跑 `make sync-claude-mcp-import` 全量生成，两者 server 列表**始终一致**。个人跨项目才用的 server 放 `~/.cursor/mcp.json`（本机全局，不提交，不影响本项目）。

两者共享同一套 Gentle-AI 生态（Engram 记忆、CodeGraph、GGA、Skills），但配置入口不同。`make setup-ai` 会**同时初始化两者**。

规则同步链路：

```
.cursor/rules/*.mdc  ──make sync-claude-rules──►  .claude/rules/*.md
.cursor/skills/      ──软链──────────────────►  .claude/skills/
```

修改规范时只改 `.cursor/rules/`，然后执行 `make sync-claude-rules`。

---

本仓库的核心 AI 栈：

| 工具 | 作用 | 官方仓库 |
|------|------|----------|
| [Gentle-AI](https://github.com/Gentleman-Programming/gentle-ai) | 配置 Cursor + Claude Code：SDD、Skills、Engram、GGA、MCP | Gentleman-Programming/gentle-ai |
| [CodeGraph](https://github.com/colbymchenry/codegraph) | 预建代码知识图谱，减少 grep/读文件，通过 MCP 给 Agent 精准上下文 | colbymchenry/codegraph |
| Engram（Gentle-AI 组件） | 跨会话持久记忆（决策、踩坑、约定）；两通道共享 | 随 gentle-ai 安装 |
| GGA（Gentle-AI 组件） | Git pre-commit AI 代码审查 | 随 gentle-ai 安装 |

项目内规则（与工具链配合，**已在 git 中**）：

- `AGENTS.md` / `.cursor/rules/*.mdc` — 开发规范与工作流（SSOT 在 `.cursor/rules/`）
- `.gga` — GGA 审查配置
- `.atl/skill-registry.md` — 技能索引（`/sdd-init` 或 `gentle-ai skill-registry refresh` 生成）
- `.engram/` — Engram 记忆块（跨设备通过 git 同步，见下文）

---

## 前置条件

| 项 | 说明 |
|----|------|
| 操作系统 | macOS / Linux（推荐）；Windows 请用 WSL2 或 Git Bash |
| Git | 已 clone 本仓库 |
| Cursor | 已安装，并启用 **Claude Code 插件** |
| Go 1.24+ | 后端开发（见 `admin-server`） |
| Node 18+ / pnpm | 前端开发（见 `admin-frontend`） |

---

## 一键初始化（新设备 / 新同事）

在仓库根目录执行：

```bash
make setup-ai
```

等价于：

```bash
bash script/setup_ai_toolchain.sh
```

脚本会依次：

1. 安装 `gentle-ai`、`codegraph` CLI（若本机缺失）
2. `gentle-ai install --agent cursor,claude-code` — 同时配置 Cursor 与 Claude Code 插件
3. `make sync-claude-rules` — 将 `.cursor/rules` 同步到 `.claude/rules`（并维护 skills 软链）
4. `gentle-ai skill-registry refresh` — 刷新 `.atl/skill-registry.md`
5. `codegraph install --target=cursor,claude` — 将 CodeGraph 接入两个通道的 MCP
6. `codegraph init` — 构建本仓库代码索引（仅首次；已有索引则跳过）
7. `engram sync --import` — 从 `.engram/` 导入团队记忆（若目录存在）

**Claude Code MCP**：项目需要的 MCP 清单在仓库根 **`.mcp.json`**（已提交 git）。第三人 clone 后即可看到要装哪些 server；执行 `make setup-ai` 会安装大部分 CLI 依赖。若使用 `mcp-zero`，需在本机设置 `GO_ZERO_MCP_PATH`（见下文「项目 MCP 清单」）。

**完成后请完全重启 Cursor**（关闭再打开），让 Cursor 与 Claude Code 插件都重新加载 MCP。

首次在 Claude Code 插件中打开本项目时，可执行 `/sdd-init` 注册项目上下文（SDD 编排器也会在无上下文时自动触发）。

验证：

```bash
make setup-ai-check
```

---

## Makefile 命令速查

| 命令 | 说明 |
|------|------|
| `make setup-ai` | 完整 AI 工具链初始化 |
| `make setup-ai-check` | 健康检查（gentle-ai doctor、codegraph status、engram status） |
| `make sync-claude-rules` | 同步 Cursor 规则到 Claude Code 格式 |
| `make sync-claude-mcp-check` | 检查 `.mcp.json` 与 `claude mcp list` 连接状态 |
| `make sync-claude-mcp-approve` | 确保 `.claude/settings.json` 自动批准项目 MCP |
| `make sync-claude-mcp-import` | 维护者：从项目 `.cursor/mcp.json` 更新团队 `.mcp.json` |
| `make engram-sync-push` | 离开设备前：导出记忆并 commit `.engram/` |
| `make engram-sync-pull` | 换设备后：`git pull` + 导入记忆 |
| `make engram-sync-status` | 查看 Engram 同步状态 |

脚本细粒度子命令：

```bash
./script/setup_ai_toolchain.sh gentle-ai   # 仅 Gentle-AI
./script/setup_ai_toolchain.sh codegraph   # 仅 CodeGraph
./script/setup_ai_toolchain.sh engram      # 仅导入记忆
./script/setup_ai_toolchain.sh check       # 仅检查
```

---

## 跨设备工作流

### 设备 A（收工）

```bash
make engram-sync-push
git push
```

### 设备 B（开工）

```bash
git clone <repo> && cd bgg
make setup-ai          # 首次：装工具链 + 导入记忆
# 或已装过工具链时：
git pull
make engram-sync-pull
```

记忆在 `.engram/`（压缩块，可安全提交）；CodeGraph 索引在 `.codegraph/`（**不提交**，每台机器本地构建）。

---

## 目录说明（维护者必知）

```
bgg/
├── AGENTS.md                 # AI 操作手册（规则整合版）
├── .cursor/rules/*.mdc       # Cursor 规则 SSOT
├── .cursor/mcp.json          # Cursor 项目 MCP 清单（提交 git，团队 SSOT，按本项目按需加载）
├── .claude/rules/*.md        # 由 sync-claude-rules 生成，勿手改
├── .claude/settings.json     # Claude Code 项目设置（自动批准 .mcp.json）
├── .mcp.json                 # Claude Code 项目 MCP 清单（提交 git，团队 SSOT）
├── .gga                      # GGA pre-commit 审查配置
├── .atl/skill-registry.md    # 技能索引
├── .engram/                  # Engram 记忆（提交 git，跨设备共享）
│   ├── config.json
│   └── chunks/
├── .codegraph/               # CodeGraph 本地索引（gitignore，每台机器独立）
├── script/
│   ├── setup_ai_toolchain.sh # 工具链一键初始化
│   ├── sync_claude_mcp.sh    # Claude Code MCP 检查 / 维护者导入
│   └── engram_sync.sh        # Engram 导入/导出
└── docs/AI工具链上手.md       # 本文
```

业务开发规范仍以 `AGENTS.md` 与 `.cursor/rules/` 为准；本文只覆盖 AI 工具链。

### Cursor vs Claude Code 插件：各自读什么

| 路径 | Cursor | Claude Code 插件 | 能否手改 |
|------|--------|------------------|----------|
| `.cursor/rules/*.mdc` | ✅ 自动挂载 | — | ✅ SSOT，改这里 |
| `.claude/rules/*.md` | — | ✅ 读取 | ❌ 运行 `make sync-claude-rules` 生成 |
| `.cursor/skills/` | ✅ | — | ✅ |
| `.claude/skills/` | — | ✅（软链到 `.cursor/skills`） | ❌ 自动生成 |
| `AGENTS.md` | 部分 Agent 读取 | ✅ CLAUDE.md 指向此处 | ✅ |
| `.engram/` | 共享记忆 | 共享记忆 | 用 `make engram-sync-*` |
| `.codegraph/` | 共享本地索引 | 共享本地索引 | 每台机器 `codegraph init` |
| `.cursor/mcp.json` | ✅ 项目 MCP 清单 | — | ✅ 团队 SSOT，唯一权威来源，改后 commit + `make sync-claude-mcp-import` |
| `.mcp.json` | — | ✅ 项目 MCP 清单 | ❌ 由 `sync-claude-mcp-import` 从 `.cursor/mcp.json` **全量生成**，不要手改，server 列表须与 `.cursor/mcp.json` 保持一致 |
| `~/.cursor/mcp.json` | ✅ 个人跨项目 server（不影响本项目） | — | 本机全局，不提交 |

---

## 项目 MCP 清单（`.mcp.json`）

仓库内 `.cursor/mcp.json` 是唯一权威来源（已提交 git，按本项目实际需要精简加载，含 `redis`，不含 `mongodb`）。仓库根 `.mcp.json` 由它跑 `make sync-claude-mcp-import` **全量生成**，第三人 clone 即可知道本项目依赖哪些 MCP，Claude Code 插件启动时会读取该文件（首次需在工作区信任对话框中确认）。两份清单的 server 列表**始终一致**——新增/删除 server 一律改 `.cursor/mcp.json`，改完跑 import 同步进 `.mcp.json`，不要单独手改 `.mcp.json`。

| Server | 用途 | 安装方式 | 必需 |
|--------|------|----------|------|
| `engram` | 跨会话持久记忆 | `make setup-ai`（随 gentle-ai） | ✅ |
| `codegraph` | 代码知识图谱，精准上下文 | `make setup-ai` | ✅ |
| `context7` | 查第三方库文档 | 需 Node/`npx`（`make setup-ai` 不单独装） | ✅ |
| `go-lsp` | Go 语言服务（gopls） | `go install .../mcp-language-server@latest` + `gopls` | ✅ 后端开发 |
| `vue-lsp` | Vue/TS 语言服务（@vue/language-server） | 与 go-lsp 共用 `mcp-language-server` + `cd admin-frontend && pnpm install` | ✅ 前端开发 |
| `frontend-ui` | 项目 UI 组件与前端约定 | 已提交 `admin-frontend/mcp/dist/`（改源码后 `pnpm mcp:build`） | ✅ 前端开发 |
| `mcp-zero` | go-zero / goctl 脚手架辅助 | 本机编译 `go-zero-mcp` 并设置 `GO_ZERO_MCP_PATH` | ✅ 后端开发 |
| `admin-mcp` | 本仓库自建脚手架/进度查询工具 | 本机构建 `admin-server/tool/admin-mcp`（见 `admin-server/docs/22-admin-mcp-tool.md`） | ✅ 后端开发 |
| `mysql` | 本地 MySQL 查询（默认只读） | `npx` + 本机 MySQL 在跑；连接信息见 `MYSQL_*` 环境变量 | 可选 |
| `redis` | 本地 Redis 查询（调试 `sqlc.CachedConn` 缓存与 MySQL 不一致问题） | `uvx --from git+https://github.com/redis/mcp-redis.git redis-mcp-server` + 本机 Redis 在跑（默认 127.0.0.1:6379，见 `REDIS_HOST`/`REDIS_PORT`） | 可选 |

`mongodb` 本项目当前**未注册为 MCP server**（不是仅本机服务未启动）；团队需要时先在 `.cursor/mcp.json` 里加回对应条目（参考 Cursor `${env:VAR}` 语法配置连接参数），再跑 `make sync-claude-mcp-import` 同步进 `.mcp.json`。

### 第三人必做（Claude Code 对话插件）

终端 `claude mcp list` 正常 **≠** 对话插件已连通——插件是独立进程，且不继承 `~/.zshrc` 里的环境变量。

```bash
git clone <repo> && cd bgg
make setup-ai

# 若使用 mcp-zero（后端开发）
export GO_ZERO_MCP_PATH=/path/to/go-zero-mcp   # 写入 ~/.zshrc

# 为对话插件写入 env + MCP 批准（本机 settings.local.json）
make sync-claude-mcp-approve

# 1) 完全退出并重启 Cursor
# 2) 集成终端首次信任工作区:
claude    # 弹出「信任此工作区」时选接受，然后 exit

# 3) 在对话插件新开对话，输入:
/mcp      # 应看到 engram/codegraph 等 Connected

# CLI 验证（可选）
make sync-claude-mcp-check
```

若插件里 `/mcp reconnect` 报 *MCP controls aren't available*，用 **Command Palette → Developer: Reload Window** 后新开对话，不要在面板里强行走 reconnect。

若 Google Drive/Calendar/Gmail 干扰，在终端 REPL 执行 `/mcp disable claude.ai Gmail` 等逐个禁用，或到 claude.ai 账户设置里移除连接器。

### 维护者：从 Cursor 更新团队清单

Cursor 侧的团队 SSOT 是仓库内 **`.cursor/mcp.json`**（已提交 git，按本项目实际需要按需加载，不是每个 server 都要有）。改动步骤：

1. 直接编辑仓库内 `.cursor/mcp.json`（新增/删除 server），本机可执行路径一律用 `${workspaceFolder}`（仓库内路径）或 `${env:VAR_NAME}`（需要本机环境变量的场景，如 `GO_ZERO_MCP_PATH`）表达，**不要写死本机绝对路径**（如 `/opt/homebrew/bin/xxx`、`/Users/<you>/...`）——这些路径只在你自己机器上有效，写进 git 会让其他人打不开
2. 完全重启 Cursor 确认新 server 能连上
3. 跑 `import-cursor` 规范化并同步进 Claude Code 侧的 `.mcp.json`（`${workspaceFolder}` → `${CLAUDE_PROJECT_DIR:-.}`，`${env:VAR}` → `${VAR}`，两边语义等价，只是各自 client 的变量引用语法不同）：

```bash
make sync-claude-mcp-import    # 或加 --dry-run 先预览
git diff .mcp.json
git add .cursor/mcp.json .mcp.json
git commit -m "chore: update project MCP servers"
```

若源文件里仍有本机绝对路径遗漏，`import-cursor` 也会尽量把已知命令（如 `engram`、`codegraph`）规范化为可移植的裸命令名，或把 `go-zero-mcp` 可执行文件路径转成 `${GO_ZERO_MCP_PATH}`，但**新增 server 时优先自己在源头写成可移植形式**，不要依赖脚本兜底。

**重要**：Cursor 项目内 `.cursor/mcp.json` 与 Claude Code 插件读的 `.mcp.json` 是两份独立文件，**不会自动同步**——维护者改了 `.cursor/mcp.json` 后必须走上面的 `import-cursor` 流程并一起 commit，否则双通道会漂移。个人只在其他项目用得到、与本项目无关的 server，放你自己的 `~/.cursor/mcp.json`（本机全局，不提交），不要混进仓库内的 `.cursor/mcp.json`。

---

## 初始化后必做清单（`make setup-ai` 不会自动完成）

`make setup-ai` 安装 Gentle-AI、CodeGraph、同步规则与记忆导入，但以下项**必须每人本机手动完成**：

### 1. `GO_ZERO_MCP_PATH`（后端开发必需）

`.mcp.json` 中 `mcp-zero` 使用环境变量占位，未设置时 `make sync-claude-mcp-check` 会报 `mcp-zero` 连接失败。

```bash
# 示例：按你本机 go-zero-mcp 编译路径修改
export GO_ZERO_MCP_PATH=/path/to/go-zero-mcp/go-zero-mcp

# 持久化（zsh）
echo 'export GO_ZERO_MCP_PATH=/path/to/go-zero-mcp/go-zero-mcp' >> ~/.zshrc
source ~/.zshrc
```

维护者本机若路径不同，只改 shell 环境即可，**不要**把个人绝对路径写进 `.mcp.json`。

### 2. `go-lsp`（后端开发必需）

`setup-ai` **不会**安装 `mcp-language-server` 与 `gopls`，需手动：

```bash
go install github.com/isaacphi/mcp-language-server@latest
go install golang.org/x/tools/gopls@latest
```

验证：`command -v mcp-language-server gopls`

### 3. 前端 MCP（`vue-lsp` + `frontend-ui`，前端开发必需）

`vue-lsp` 与 `go-lsp` 共用 `mcp-language-server`（见上一节）。另需安装前端依赖：

```bash
cd admin-frontend && pnpm install
```

- **`vue-lsp`**：workspace 指向 `admin-frontend/`，通过 `node` 启动 `@vue/language-server`（`node_modules` 内路径相对 workspace）
- **`frontend-ui`**：Node 直启 `admin-frontend/mcp/dist/index.js`，无需额外安装（dist 已提交 git）

验证：`test -f admin-frontend/node_modules/@vue/language-server/bin/vue-language-server.js && node admin-frontend/mcp/dist/index.js`（后者会挂起等待 stdio，Ctrl+C 退出即可）

### 4. `context7` / `npx`

需 Node.js 18+。一般开发机已有；无则安装 Node 后 `make sync-claude-mcp-check` 中 `context7` 应显示 Connected。

### 5. `mysql` MCP（可选）

| Server | 前提 | 不用时 |
|--------|------|--------|
| `mysql` | 本机 MySQL 在跑；设置 `MYSQL_HOST`/`MYSQL_PORT`/`MYSQL_USER`/`MYSQL_PASS`/`MYSQL_DB`（可与 `config/mysql.json.example` 或 `/etc/work/mysql.json` 对齐） | `check` 失败可忽略；团队不用可从 `.cursor/mcp.json` 删除后 `import` 同步 |

`mongodb` 本项目当前未注册（见上文「项目 MCP 清单」）；需要时先在 `.cursor/mcp.json` 加回对应 server 再 `import`。

`mysql` MCP 默认 **只读**（`ALLOW_*_OPERATION=false`），用于调试联调查表/跑 SELECT，不替代业务代码里的 Repository 层。

示例（写入 `~/.zshrc`，字段名与 `config/mysql.json.example` 对应）：

```bash
export MYSQL_HOST=127.0.0.1
export MYSQL_PORT=3306
export MYSQL_USER=root
export MYSQL_PASS=your-password
export MYSQL_DB=postapoc_admin
```

Claude Code 插件需 `make sync-claude-mcp-approve` 把上述 env 写入 `.claude/settings.local.json`（或在本机创建 `~/.config/bgg/mysql-mcp.env` 后 source 再执行 approve）。

推荐把凭据放在 **`~/.config/bgg/mysql-mcp.env`**（`chmod 600`），`~/.zshrc` 里 source 该文件。Cursor / Claude Code 均通过仓库内 **`script/mysql_mcp.sh`** 启动 `mysql` MCP（自动补全 `HOME` 并 source 上述 env 文件），避免密码写进 mcp.json。

### 6. `redis` MCP（可选）

| Server | 前提 | 不用时 |
|--------|------|--------|
| `redis` | `uvx` 可用（`brew install uv` 或参考 [astral-sh/uv](https://github.com/astral-sh/uv)）；本机 Redis 在跑；Cursor 侧 `.cursor/mcp.json` 里 `env.REDIS_HOST`/`env.REDIS_PORT` 硬编码 `127.0.0.1`/`6379`（开箱可用），连别的 host/port 需直接改这两个字段再 `import` | `check` 失败可忽略；团队不用可从 `.cursor/mcp.json` 删除后 `import` 同步 |

Claude Code 侧 `.mcp.json` 由 `import-cursor` 自动把上述字段转成 `${REDIS_HOST:-127.0.0.1}`/`${REDIS_PORT:-6379}`，支持用 shell 环境变量临时覆盖；Cursor 侧无此机制，改连接目标必须直接编辑 `.cursor/mcp.json`。

主要用于排查 go-zero `sqlc.CachedConn` 缓存与 MySQL 数据不一致的问题——例如手动 `TRUNCATE`/`DELETE` 重置开发库表时绕过了 Model 层的 `Insert`/`Update`/`Delete`，导致 Redis 里按主键缓存的旧记录（`cache:<model>:id:<id>`）未被失效，后续 `FindOne`/`FindByID` 读到脏数据而走 squirrel 直查的 `FindPage`/列表接口正常。排查方法：`GET`/`KEYS`/`TTL` 确认缓存内容与 MySQL 是否一致，确认脏读后 `DEL` 对应 key 即可（cache-aside 模式下下次读取会自动回源重建缓存）。

### 7. Claude Code MCP 重连

`setup-ai` 后在 **Cursor 集成终端**（非对话面板）：

```bash
make sync-claude-mcp-check
claude
/mcp reconnect all
```

期望：必需 server 为 Connected；`mcp-zero` 在设置 `GO_ZERO_MCP_PATH` 后应 Connected。

### 8. 维护者：双通道 MCP 同步纪律

| 你改了什么 | 必须做什么 |
|------------|------------|
| 仓库内 `.cursor/mcp.json` 增删 server | `make sync-claude-mcp-import` → `git diff .mcp.json` → 两个文件一起 commit |
| `.cursor/rules/*.mdc` | `make sync-claude-rules` |
| AI 会话记忆要带给团队 | `make engram-sync-push` → `git push` |

**禁止**只改仓库内 `.cursor/mcp.json` 而不更新 `.mcp.json`——第三人 clone 后 Claude Code 插件会对不上。个人 `~/.cursor/mcp.json`（本机全局）里的改动不影响团队，也不需要走这套同步流程。

---

## 可选配置

### 仅当前仓库安装 Gentle-AI（不影响其他项目）

```bash
GENTLE_AI_INSTALL_SCOPE=workspace make setup-ai
```

### 只使用单一通道

```bash
# 仅 Cursor
GENTLE_AI_AGENT=cursor CODEGRAPH_TARGET=cursor make setup-ai

# 仅 Claude Code 插件（gentle-ai 用 claude-code，codegraph 用 claude）
GENTLE_AI_AGENT=claude-code CODEGRAPH_TARGET=claude make setup-ai
```

支持列表见 [gentle-ai agents 文档](https://github.com/Gentleman-Programming/gentle-ai/blob/main/docs/agents.md)。

### 跳过部分步骤

```bash
SKIP_CODEGRAPH=1 make setup-ai    # 暂不建 CodeGraph 索引
SKIP_ENGRAM=1 make setup-ai       # 不导入记忆
```

---

## 日常维护

| 场景 | 操作 |
|------|------|
| 改了 `.cursor/rules/` | `make sync-claude-rules`（Claude Code 插件才会跟上） |
| 升级 Gentle-AI 生态 | `gentle-ai upgrade` 或 `gentle-ai sync` |
| 升级 CodeGraph | `codegraph upgrade` |
| 新增/删除 Skills 后 | `gentle-ai skill-registry refresh` |
| 大重构后重建索引 | `codegraph index --force` |
| 查看 Engram 记忆 | `engram tui` 或 `engram search "关键词"` |
| 生态健康检查 | `gentle-ai doctor` |
| Cursor 增删 MCP 后 | 改 `.cursor/mcp.json` → `make sync-claude-mcp-import` → commit `.cursor/mcp.json` + `.mcp.json` |
| Claude Code MCP 连不上 | `make sync-claude-mcp-check`，REPL 内 `/mcp reconnect all` |
| Claude Code 里项目上下文过期 | 在插件中执行 `/sdd-init` |

CodeGraph 默认**自动同步**文件变更，一般无需手动 `codegraph sync`。

---

## 常见问题

### `make setup-ai` 后看不到 CodeGraph / Engram MCP

1. **完全重启 Cursor**（不是只重载窗口）
2. Cursor 通道：检查仓库内 **`.cursor/mcp.json`**（团队 SSOT）和个人 `~/.cursor/mcp.json`（本机全局）
3. Claude Code 插件：检查仓库根 **`.mcp.json`**，运行 `make sync-claude-mcp-check`
4. 重装 MCP 接线：
   ```bash
   codegraph install --target=cursor,claude --yes
   gentle-ai sync
   ```

### CodeGraph 报 "not initialized"

在仓库根目录执行 `codegraph init` 或 `make setup-ai`。

### `mcp-zero` 连接失败

`GO_ZERO_MCP_PATH` 未设置或路径错误。见上文「初始化后必做清单」第 1 节。

### `go-lsp` 连接失败

未安装 `mcp-language-server` 或 `gopls`。见上文「初始化后必做清单」第 2 节。

### `vue-lsp` 连接失败

1. 未安装 `mcp-language-server`（与 go-lsp 相同）
2. 未执行 `cd admin-frontend && pnpm install`（缺少 `@vue/language-server`）
3. `admin-frontend/node_modules/@vue/language-server/bin/vue-language-server.js` 不存在 → 在 `admin-frontend` 重新 `pnpm install`

### `frontend-ui` 连接失败

检查 `admin-frontend/mcp/dist/index.js` 是否存在；若改过 `mcp/src/` 需 `cd admin-frontend && pnpm mcp:build` 后 commit dist。

### `mysql` 连接失败

本机未启动 MySQL 服务或连接信息未配置——可忽略。若全员不用，维护者从 `.cursor/mcp.json` 删除后 `import` 同步进 `.mcp.json` 一起 commit。`mongodb` 本项目当前未注册，不适用本节。

### `redis` 连接失败

本机未启动 Redis 服务、`uvx` 不存在，或 `REDIS_HOST`/`REDIS_PORT` 配置不对——可忽略（属可选 server）。首次连接需要 `uvx` 从 GitHub 拉取并构建 `redis/mcp-redis`，耗时较久属正常现象。若全员不用，维护者从 `.cursor/mcp.json` 删除后 `import` 同步进 `.mcp.json` 一起 commit。

**`mysql` 日志出现 `Connection closed`（连上立刻断开）** 常见两类原因：

1. **`mysql-mcp.env` 未加载**（Cursor Shared MCP 常不带 `HOME`，旧版内联 `source $HOME/.config/...` 会静默失败）→ 进程回退连 `127.0.0.1:3306` 后退出。确认 `.cursor/mcp.json` 使用 `${workspaceFolder}/script/mysql_mcp.sh`，并执行 `chmod +x script/mysql_mcp.sh`。
2. **凭据/库名错误** 或目标库不可达。本地调试：`ENABLE_LOGGING=1 bash script/mysql_mcp.sh`（会打印连接目标；Ctrl+C 退出）。远程库（如 SQLPub）一般无需 `MYSQL_SSL=true`；若报 SSL 相关错误再按需配置。

其他：`MYSQL_PASS`/`MYSQL_DB` 与真实库不一致时也会启动失败。

### GGA pre-commit 导致 commit 失败

GGA 在 `git commit` 时做 AI 审查。若 Provider 未配置或额度用尽：

- 配置 `~/.config/gga/config` 中的 `PROVIDER`
- 或临时：`GGA_PROVIDER=claude git commit ...`
- 排查：`/opt/homebrew/bin/gga config`（注意 shell 里 `gga` 若有 alias 会冲突）

### Engram 项目名不一致

```bash
engram projects list
engram projects consolidate
```

### `.codegraph/` 要不要提交？

**不要。** `.codegraph/.gitignore` 已忽略数据库；每台机器 `codegraph init` 本地构建即可。

---

## 相关文档

- [后端维护导航](admin-server-维护导航.md)
- [后端开发进度](后端开发进度.md)
- [前端开发进度](前端开发进度.md)
- [Gentle-AI 官方文档](https://github.com/Gentleman-Programming/gentle-ai)
- [CodeGraph 官方文档](https://github.com/colbymchenry/codegraph)
