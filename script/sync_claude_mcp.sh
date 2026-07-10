#!/usr/bin/env bash
# Claude Code 项目 MCP 管理。
#
# 团队 SSOT：仓库根 .mcp.json（已提交 git，第三人 clone 即可见项目需要哪些 MCP）
#
# 用法:
#   ./script/sync_claude_mcp.sh check           # 检查项目 .mcp.json 与 claude mcp 状态（默认）
#   ./script/sync_claude_mcp.sh approve         # 写入本机 .claude/settings.local.json 自动批准
#   ./script/sync_claude_mcp.sh import-cursor   # 维护者：从 ~/.cursor/mcp.json 导入并规范化路径后更新 .mcp.json
#   ./script/sync_claude_mcp.sh import-cursor --dry-run
#
# 环境变量（import-cursor 写入 .mcp.json 时用于替换本机绝对路径）:
#   CURSOR_MCP_JSON      源文件，默认 ~/.cursor/mcp.json
#   GO_ZERO_MCP_PATH     mcp-zero 可执行文件路径（第三人需在 shell 或 .env 中设置）
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
CURSOR_MCP_JSON="${CURSOR_MCP_JSON:-$HOME/.cursor/mcp.json}"
TARGET_MCP_JSON="$REPO_ROOT/.mcp.json"
SETTINGS_JSON="$REPO_ROOT/.claude/settings.json"
SETTINGS_LOCAL="$REPO_ROOT/.claude/settings.local.json"

log_info() { echo "[INFO] $*"; }
log_warn() { echo "[WARN] $*" >&2; }
log_error() { echo "[ERROR] $*" >&2; }

usage() {
  cat <<'EOF'
Claude Code 项目 MCP（.mcp.json 为团队 SSOT，已提交 git）

子命令:
  check (默认)     列出项目 .mcp.json 中的 server，并运行 claude mcp list
  approve          写入 settings.json + settings.local.json（对话插件 env/批准）
  import-cursor    维护者专用：从 ~/.cursor/mcp.json 导入并规范化路径，更新 .mcp.json
                   加 --dry-run 仅预览

第三人上手:
  1. git clone 后阅读 .mcp.json 与 docs/AI工具链上手.md
  2. make setup-ai
  3. 设置 GO_ZERO_MCP_PATH（若使用 mcp-zero）
  4. 在终端 REPL: claude → /mcp reconnect all

示例:
  make sync-claude-mcp-check
  ./script/sync_claude_mcp.sh approve
  ./script/sync_claude_mcp.sh import-cursor --dry-run
EOF
}

require_command() {
  local cmd="$1"
  local hint="$2"
  if ! command -v "$cmd" &>/dev/null; then
    log_error "未找到命令: $cmd"
    log_error "$hint"
    return 1
  fi
}

normalize_cursor_mcp_json() {
  local source="$1"
  python3 - "$source" <<'PY'
import json
import os
import re
import sys
from pathlib import Path

source = Path(sys.argv[1])
raw = source.read_text(encoding="utf-8")
converted = raw.replace("${workspaceFolder}", "${CLAUDE_PROJECT_DIR:-.}")

data = json.loads(converted)
servers = data.get("mcpServers")
if not isinstance(servers, dict) or not servers:
    raise SystemExit("源文件缺少 mcpServers 或为空")

home = str(Path.home())
brew_bin = "/opt/homebrew/bin"
local_bin = str(Path.home() / ".local/bin")

KNOWN_CMDS = {
    "engram", "codegraph", "npx", "uvx", "mcp-language-server",
    "vue-ts-lsp", "node",
}

def normalize_command(cmd: str) -> str:
    if not cmd:
        return cmd
    if cmd in KNOWN_CMDS:
        return cmd
    if cmd.startswith(f"{brew_bin}/"):
        return Path(cmd).name
    if cmd.startswith(f"{local_bin}/"):
        return Path(cmd).name
    if cmd.endswith("/go-zero-mcp") or cmd.endswith("/go-zero-mcp/go-zero-mcp"):
        return "${GO_ZERO_MCP_PATH}"
    return cmd

def normalize_env(env: dict) -> dict:
    if not isinstance(env, dict):
        return env
    out = {}
    for key, value in env.items():
        if key == "GOCTL_PATH" and isinstance(value, str) and value.startswith(home):
            out[key] = "${GOCTL_PATH:-goctl}"
        elif key in ("REDIS_HOST", "REDIS_PORT"):
            out[key] = f"${{{key}:-{value}}}"
        elif key in ("MYSQL_HOST", "MYSQL_PORT", "MYSQL_USER", "MYSQL_PASS", "MYSQL_DB"):
            defaults = {
                "MYSQL_HOST": "127.0.0.1",
                "MYSQL_PORT": "3306",
                "MYSQL_USER": "root",
                "MYSQL_PASS": "",
                "MYSQL_DB": "postapoc_admin",
            }
            default = defaults.get(key, value)
            out[key] = f"${{{key}:-{default}}}"
        elif key in (
            "ALLOW_INSERT_OPERATION",
            "ALLOW_UPDATE_OPERATION",
            "ALLOW_DELETE_OPERATION",
        ):
            out[key] = "false"
        else:
            out[key] = value
    return out

def normalize_args(args: list) -> list:
    out = []
    for arg in args:
        if isinstance(arg, str) and arg.startswith("mongodb://") and "localhost" in arg:
            out.append("${MONGODB_MCP_URI:-mongodb://localhost:27017/?replicaSet=rs0}")
        else:
            out.append(arg)
    return out

normalized = {}
for name, cfg in servers.items():
    entry = dict(cfg)
    if "command" in entry:
        entry["command"] = normalize_command(entry["command"])
    if "args" in entry:
        entry["args"] = normalize_args(entry["args"])
    if "env" in entry:
        entry["env"] = normalize_env(entry["env"])
    normalized[name] = entry

print(json.dumps({"mcpServers": normalized}, indent=2, ensure_ascii=False))
print("", end="")
PY
}

write_auto_approve_settings() {
  mkdir -p "$(dirname "$SETTINGS_JSON")"
  if [ -f "$SETTINGS_JSON" ]; then
    python3 - "$SETTINGS_JSON" <<'PY'
import json
import sys
from pathlib import Path

path = Path(sys.argv[1])
data = {}
if path.exists():
    data = json.loads(path.read_text(encoding="utf-8") or "{}")
data["enableAllProjectMcpServers"] = True
path.write_text(json.dumps(data, indent=2, ensure_ascii=False) + "\n", encoding="utf-8")
PY
  else
    cat >"$SETTINGS_JSON" <<'EOF'
{
  "enableAllProjectMcpServers": true
}
EOF
  fi
  log_info "已写入 ${SETTINGS_JSON} (enableAllProjectMcpServers: true)"

  # 对话插件不继承 shell 的 ~/.zshrc，需在本机 settings.local.json 注入 env
  python3 - "$SETTINGS_LOCAL" "$TARGET_MCP_JSON" <<'PY'
import json
import os
import sys
from pathlib import Path

local_path = Path(sys.argv[1])
mcp_path = Path(sys.argv[2])

data = {}
if local_path.exists():
    data = json.loads(local_path.read_text(encoding="utf-8") or "{}")

data["enableAllProjectMcpServers"] = True

servers = []
if mcp_path.exists():
    servers = sorted(json.loads(mcp_path.read_text(encoding="utf-8")).get("mcpServers", {}).keys())
if servers:
    data["enabledMcpjsonServers"] = servers

env = dict(data.get("env") or {})
for key in (
    "GO_ZERO_MCP_PATH",
    "GOCTL_PATH",
    "MONGODB_MCP_URI",
    "REDIS_HOST",
    "REDIS_PORT",
    "MYSQL_HOST",
    "MYSQL_PORT",
    "MYSQL_USER",
    "MYSQL_PASS",
    "MYSQL_DB",
):
    val = os.environ.get(key)
    if val:
        env[key] = val
if env:
    data["env"] = env

local_path.parent.mkdir(parents=True, exist_ok=True)
local_path.write_text(json.dumps(data, indent=2, ensure_ascii=False) + "\n", encoding="utf-8")
PY
  log_info "已写入 ${SETTINGS_LOCAL} (env + enabledMcpjsonServers，供对话插件使用)"
  if [ -z "${GO_ZERO_MCP_PATH:-}" ]; then
    log_warn "GO_ZERO_MCP_PATH 未设置，插件中 mcp-zero 可能无法连接"
    log_warn "请 export 后重跑: make sync-claude-mcp-approve"
  fi
  echo ""
  log_info "对话插件生效步骤:"
  echo "  1. 完全退出并重启 Cursor"
  echo "  2. 在集成终端执行: cd ${REPO_ROOT} && claude  （首次接受「信任工作区」）"
  echo "  3. 在插件新对话输入: /mcp"
  echo "  4. 若仍 not connected: Command Palette → Developer: Reload Window"
}

run_import_cursor() {
  local dry_run="${1:-0}"

  if [ ! -f "$CURSOR_MCP_JSON" ]; then
    log_error "未找到 Cursor MCP 配置: ${CURSOR_MCP_JSON}"
    return 1
  fi

  require_command python3 "需要 Python 3"

  local content
  content="$(normalize_cursor_mcp_json "$CURSOR_MCP_JSON")"
  local server_count
  server_count="$(python3 -c 'import json,sys; print(len(json.load(sys.stdin).get("mcpServers", {})))' <<<"$content")"

  if [ "$dry_run" = "1" ]; then
    echo "$content"
    log_info "dry-run: 将写入 ${server_count} 个 server 到 .mcp.json (未落盘)"
    return 0
  fi

  echo "$content" >"$TARGET_MCP_JSON"
  log_info "已更新 ${TARGET_MCP_JSON} (${server_count} 个 server, 源: ${CURSOR_MCP_JSON})"
  log_info "请检查 diff 后提交 git，供团队使用"
}

print_project_mcp() {
  local file="$1"
  local label="$2"
  if [ ! -f "$file" ]; then
    log_warn "未找到 ${label}: ${file}"
    return 0
  fi
  log_info "${label}: ${file}"
  python3 - "$file" <<'PY'
import json, sys
from pathlib import Path
data = json.loads(Path(sys.argv[1]).read_text(encoding="utf-8"))
servers = data.get("mcpServers", {})
print(f"  servers: {len(servers)}")
for name in sorted(servers):
    cfg = servers[name]
    cmd = cfg.get("command", cfg.get("url", "?"))
    env_keys = ", ".join(sorted(cfg.get("env", {}).keys()))
    extra = f" env=[{env_keys}]" if env_keys else ""
    print(f"  - {name}: {cmd}{extra}")
PY
}

run_check() {
  echo ""
  echo "========== 项目 MCP（团队 SSOT，已提交 git）=========="
  print_project_mcp "$TARGET_MCP_JSON" "项目 .mcp.json"

  echo ""
  echo "========== 维护者 Cursor 配置（参考）=========="
  print_project_mcp "$CURSOR_MCP_JSON" "Cursor ~/.cursor/mcp.json"

  echo ""
  echo "========== 环境变量（按需设置）=========="
  echo "  GO_ZERO_MCP_PATH=${GO_ZERO_MCP_PATH:-<未设置>}"
  echo "  GOCTL_PATH=${GOCTL_PATH:-goctl (默认)}"
  echo "  MONGODB_MCP_URI=${MONGODB_MCP_URI:-mongodb://localhost:27017/?replicaSet=rs0 (默认)}"
  echo "  REDIS_HOST=${REDIS_HOST:-127.0.0.1 (默认)}  REDIS_PORT=${REDIS_PORT:-6379 (默认)}"

  echo ""
  echo "========== 前端 MCP 依赖 =========="
  local fe_root="$REPO_ROOT/admin-frontend"
  if [ -f "$fe_root/node_modules/@vue/language-server/bin/vue-language-server.js" ]; then
    echo "  @vue/language-server: OK"
  else
    log_warn "@vue/language-server 未安装 → 在 admin-frontend 执行: pnpm install"
  fi
  if [ -f "$fe_root/mcp/dist/index.js" ]; then
    echo "  frontend-ui: OK ($fe_root/mcp/dist/index.js)"
  else
    log_warn "frontend-ui 未构建 → 在 admin-frontend 执行: pnpm mcp:build"
  fi

  echo ""
  echo "========== Claude CLI =========="
  if ! command -v claude &>/dev/null; then
    log_warn "未安装 claude CLI，跳过 mcp list"
    return 0
  fi

  log_info "claude 版本: $(claude --version 2>/dev/null || echo unknown)"
  echo ""
  (cd "$REPO_ROOT" && claude mcp list) || log_warn "claude mcp list 失败（可能需首次运行 claude 并信任工作区）"
}

main() {
  local cmd="${1:-check}"
  if [ $# -gt 0 ]; then
    shift
  fi

  case "$cmd" in
    check|status)
      run_check "$@"
      ;;
    approve|auto-approve)
      write_auto_approve_settings
      ;;
    import-cursor|import)
      local dry_run=0
      if [ "${1:-}" = "--dry-run" ] || [ "${1:-}" = "dry-run" ]; then
        dry_run=1
      fi
      run_import_cursor "$dry_run"
      ;;
    -h|--help|help)
      usage
      ;;
    *)
      log_error "未知子命令: $cmd"
      usage
      exit 1
      ;;
  esac
}

main "$@"
