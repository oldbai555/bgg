# AI 工具链上手（Gentle-AI + CodeGraph）

本文面向**换设备继续开发**或**第三位维护者**快速接入本项目的 AI 辅助开发环境。

## 本项目的日常开发方式

团队默认使用 **Cursor IDE + Cursor 内的 Claude Code 插件** 双通道开发：

| 通道 | 用途 | 规则来源 | MCP 配置 |
|------|------|----------|----------|
| **Cursor**（Chat / Agent / Composer） | 日常对话、改代码、子 Agent | `.cursor/rules/*.mdc`（SSOT） | `~/.cursor/mcp.json` |
| **Claude Code 插件**（Cursor 侧边栏） | SDD 工作流、Task 子代理、长任务 | `.claude/rules/*.md`（由 Cursor 规则同步） | 仓库根 **`.mcp.json`**（已提交 git） |

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
| `make sync-claude-mcp-import` | 维护者：从 `~/.cursor/mcp.json` 更新团队 `.mcp.json` |
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
| `.mcp.json` | — | ✅ 项目 MCP 清单 | ✅ 团队 SSOT，改后 commit |

---

## 项目 MCP 清单（`.mcp.json`）

仓库根 `.mcp.json` **已提交 git**，第三人 clone 即可知道本项目依赖哪些 MCP。Claude Code 插件启动时会读取该文件（首次需在工作区信任对话框中确认）。

| Server | 用途 | 安装方式 | 必需 |
|--------|------|----------|------|
| `engram` | 跨会话持久记忆 | `make setup-ai`（随 gentle-ai） | ✅ |
| `codegraph` | 代码知识图谱，精准上下文 | `make setup-ai` | ✅ |
| `context7` | 查第三方库文档 | 需 Node/`npx`（`make setup-ai` 不单独装） | ✅ |
| `go-lsp` | Go 语言服务（gopls） | `go install .../mcp-language-server@latest` + `gopls` | ✅ 后端开发 |
| `mcp-zero` | go-zero / goctl 脚手架辅助 | 本机编译 `go-zero-mcp` 并设置 `GO_ZERO_MCP_PATH` | ✅ 后端开发 |
| `mongodb` | 本地 Mongo 查询 | `npx` + 本机 MongoDB 在跑 | 可选 |
| `redis` | 本地 Redis 查询 | `uvx` + 本机 Redis 在跑 | 可选 |

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

在 `~/.cursor/mcp.json` 增删 MCP 后，规范化并写回仓库：

```bash
make sync-claude-mcp-import    # 或加 --dry-run 先预览
git diff .mcp.json
git commit -m "chore: update project MCP servers"
```

`import-cursor` 会自动把本机绝对路径（如 `/opt/homebrew/bin/engram`）转为可移植的 `engram` 或 `${GO_ZERO_MCP_PATH}` 等形式。

**重要**：Cursor 的 MCP 配置在 `~/.cursor/mcp.json`（本机全局），Claude Code 插件读仓库 `.mcp.json`（团队 git）。两边**不会自动同步**——维护者改 Cursor 后必须走上面的 `import-cursor` 流程，否则双通道会漂移。

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

### 3. `context7` / `npx`

需 Node.js 18+。一般开发机已有；无则安装 Node 后 `make sync-claude-mcp-check` 中 `context7` 应显示 Connected。

### 4. `mongodb` / `redis` MCP（可选）

| Server | 前提 | 不用时 |
|--------|------|--------|
| `mongodb` | 本机 MongoDB 在跑，连接串可用 | `check` 失败可忽略；团队不用可从 `.mcp.json` 删除后 `import` 同步 |
| `redis` | 本机 Redis 在跑；需 `uv`（`uvx`） | 同上 |

仅做 admin-server MySQL 开发、不查 Mongo/Redis 时，可不启这两个服务。

### 5. Claude Code MCP 重连

`setup-ai` 后在 **Cursor 集成终端**（非对话面板）：

```bash
make sync-claude-mcp-check
claude
/mcp reconnect all
```

期望：必需 server 为 Connected；`mcp-zero` 在设置 `GO_ZERO_MCP_PATH` 后应 Connected。

### 6. 维护者：双通道 MCP 同步纪律

| 你改了什么 | 必须做什么 |
|------------|------------|
| `~/.cursor/mcp.json` 增删 server | `make sync-claude-mcp-import` → `git diff .mcp.json` → commit |
| `.cursor/rules/*.mdc` | `make sync-claude-rules` |
| AI 会话记忆要带给团队 | `make engram-sync-push` → `git push` |

**禁止**只改 Cursor 全局 MCP 而不更新 `.mcp.json`——第三人 clone 后 Claude Code 插件会对不上。

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
| Cursor 增删 MCP 后 | `make sync-claude-mcp-import` → commit `.mcp.json` |
| Claude Code MCP 连不上 | `make sync-claude-mcp-check`，REPL 内 `/mcp reconnect all` |
| Claude Code 里项目上下文过期 | 在插件中执行 `/sdd-init` |

CodeGraph 默认**自动同步**文件变更，一般无需手动 `codegraph sync`。

---

## 常见问题

### `make setup-ai` 后看不到 CodeGraph / Engram MCP

1. **完全重启 Cursor**（不是只重载窗口）
2. Cursor 通道：检查 `~/.cursor/mcp.json`
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

### `mongodb` / `redis` 连接失败

本机未启动对应服务，或团队不使用——可忽略。若全员不用，维护者从 `.mcp.json` 删除后 commit。

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
