#!/usr/bin/env bash
# 新设备 / 第三位维护者：初始化 Gentle-AI + CodeGraph 工具链，并同步本项目 AI 上下文。
#
# 用法:
#   ./script/setup_ai_toolchain.sh              # 完整初始化（默认）
#   ./script/setup_ai_toolchain.sh check        # 仅健康检查
#   ./script/setup_ai_toolchain.sh gentle-ai    # 仅 Gentle-AI 相关
#   ./script/setup_ai_toolchain.sh codegraph    # 仅 CodeGraph 相关
#   ./script/setup_ai_toolchain.sh engram       # 仅导入 Engram 记忆
#
# 环境变量:
#   GENTLE_AI_AGENT=cursor,claude-code  # 默认：Cursor + Claude Code 插件
#   GENTLE_AI_INSTALL_SCOPE=global      # global | workspace（默认 global）
#   CODEGRAPH_TARGET=cursor,claude        # CodeGraph 侧 Claude Code 插件 id 为 claude
#   SKIP_GENTLE_AI=1                # 跳过 Gentle-AI 安装/同步
#   SKIP_CODEGRAPH=1                # 跳过 CodeGraph 安装/索引
#   SKIP_ENGRAM=1                   # 跳过 Engram 记忆导入
#   SKIP_CLAUDE_RULES=1             # 跳过 make sync-claude-rules
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

GENTLE_AI_AGENT="${GENTLE_AI_AGENT:-cursor,claude-code}"
GENTLE_AI_INSTALL_SCOPE="${GENTLE_AI_INSTALL_SCOPE:-global}"
CODEGRAPH_TARGET="${CODEGRAPH_TARGET:-cursor,claude}"

log_info() { echo "[INFO] $*"; }
log_warn() { echo "[WARN] $*" >&2; }
log_error() { echo "[ERROR] $*" >&2; }

usage() {
  cat <<'EOF'
AI 工具链初始化（Gentle-AI + CodeGraph + Engram）

子命令:
  init (默认)     完整初始化：CLI → Agent 配置 → 项目索引 → 记忆导入
  check           健康检查（gentle-ai doctor、codegraph status、engram status）
  gentle-ai       仅安装/同步 Gentle-AI（含 sync-claude-rules、skill-registry）
  codegraph       仅安装 CodeGraph CLI、接入 Cursor + Claude Code、构建项目索引
  engram          仅从 .engram/ 导入跨设备记忆

环境变量见脚本头部注释。

示例:
  make setup-ai
  GENTLE_AI_INSTALL_SCOPE=workspace ./script/setup_ai_toolchain.sh
  ./script/setup_ai_toolchain.sh check
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

install_gentle_ai_cli() {
  if command -v gentle-ai &>/dev/null; then
    log_info "gentle-ai 已安装: $(gentle-ai version 2>/dev/null || true)"
    return 0
  fi

  log_info "安装 gentle-ai CLI ..."
  if command -v brew &>/dev/null; then
    if ! brew tap Gentleman-Programming/homebrew-tap &>/dev/null; then
      log_warn "brew tap 失败，改用官方安装脚本"
    else
      brew install gentle-ai && return 0
    fi
  fi

  curl -fsSL https://raw.githubusercontent.com/Gentleman-Programming/gentle-ai/main/scripts/install.sh | bash
  hash -r 2>/dev/null || true
  require_command gentle-ai "请重新打开终端后重试，或手动安装: https://github.com/Gentleman-Programming/gentle-ai"
}

install_codegraph_cli() {
  if command -v codegraph &>/dev/null; then
    log_info "codegraph 已安装: $(codegraph version 2>/dev/null || true)"
    return 0
  fi

  log_info "安装 codegraph CLI ..."
  curl -fsSL https://raw.githubusercontent.com/colbymchenry/codegraph/main/install.sh | sh
  hash -r 2>/dev/null || true

  if ! command -v codegraph &>/dev/null; then
    log_warn "codegraph 可能安装在 ~/.local/bin，尝试加入 PATH"
    export PATH="${HOME}/.local/bin:${PATH}"
  fi
  require_command codegraph "请重新打开终端，或执行: export PATH=\"\$HOME/.local/bin:\$PATH\""
}

setup_gentle_ai() {
  install_gentle_ai_cli

  log_info "配置 Gentle-AI（agent=${GENTLE_AI_AGENT}, scope=${GENTLE_AI_INSTALL_SCOPE}）..."
  gentle-ai install \
    --agent "$GENTLE_AI_AGENT" \
    --scope "$GENTLE_AI_INSTALL_SCOPE"

  log_info "同步 Gentle-AI 组件到最新版本 ..."
  gentle-ai sync || log_warn "gentle-ai sync 未完成，可稍后手动执行"

  if [ "${SKIP_CLAUDE_RULES:-0}" != "1" ]; then
    log_info "同步 .cursor/rules → .claude/rules ..."
    make -C "$REPO_ROOT" sync-claude-rules
  fi

  if gentle-ai skill-registry refresh &>/dev/null; then
    log_info "已刷新 skill-registry（.atl/skill-registry.md）"
  else
    log_warn "skill-registry refresh 跳过（可在 Cursor 中执行 /sdd-init）"
  fi
}

setup_claude_mcp() {
  if [ ! -f "$REPO_ROOT/.mcp.json" ]; then
    log_warn "缺失 .mcp.json — Claude Code 插件无法加载项目 MCP（维护者执行 make sync-claude-mcp-import）"
    return 0
  fi

  if [ ! -f "$REPO_ROOT/.claude/settings.json" ]; then
    log_info "写入 .claude/settings.json（自动批准项目 MCP）..."
    bash "$SCRIPT_DIR/sync_claude_mcp.sh" approve
  fi

  if [ -z "${GO_ZERO_MCP_PATH:-}" ]; then
    log_warn "GO_ZERO_MCP_PATH 未设置 → mcp-zero 不可用（后端开发请 export 并写入 ~/.zshrc）"
  fi
}

setup_codegraph() {
  install_codegraph_cli

  log_info "将 CodeGraph 接入 Agent（target=${CODEGRAPH_TARGET}）..."
  codegraph install --target="$CODEGRAPH_TARGET" --yes

  cd "$REPO_ROOT"
  if codegraph status &>/dev/null; then
    log_info "CodeGraph 索引已存在，跳过 init（可用 codegraph index --force 重建）"
    codegraph status | head -20
  else
    log_info "构建 CodeGraph 项目索引（首次可能需数分钟）..."
    codegraph init
    codegraph status | head -20
  fi
}

setup_engram() {
  if ! command -v engram &>/dev/null; then
    log_warn "未找到 engram，跳过记忆导入（通常随 gentle-ai 安装）"
    return 0
  fi

  if [ ! -d "$REPO_ROOT/.engram" ]; then
    log_warn ".engram/ 不存在，跳过导入（团队记忆尚未提交到仓库）"
    return 0
  fi

  log_info "导入 Engram 跨设备记忆 ..."
  bash "$SCRIPT_DIR/engram_sync.sh" import
}

run_check() {
  local failed=0

  echo ""
  echo "========== Gentle-AI =========="
  if command -v gentle-ai &>/dev/null; then
    gentle-ai doctor || failed=1
  else
    log_error "gentle-ai 未安装"
    failed=1
  fi

  echo ""
  echo "========== CodeGraph =========="
  if command -v codegraph &>/dev/null; then
    (cd "$REPO_ROOT" && codegraph status) || failed=1
  else
    log_error "codegraph 未安装"
    failed=1
  fi

  echo ""
  echo "========== Engram =========="
  if command -v engram &>/dev/null; then
    (cd "$REPO_ROOT" && engram sync --status) || failed=1
  else
    log_warn "engram 未安装（可选）"
  fi

  echo ""
  echo "========== 项目 AI 规则文件 =========="
  local files=(
    "AGENTS.md"
    ".mcp.json"
    ".cursor/rules/00-workflow.mdc"
    ".gga"
    ".atl/skill-registry.md"
    ".engram/config.json"
  )
  for f in "${files[@]}"; do
    if [ -e "$REPO_ROOT/$f" ]; then
      log_info "OK  $f"
    else
      log_warn "缺失 $f"
    fi
  done

  echo ""
  echo "========== Claude Code 项目 MCP =========="
  if [ -f "$REPO_ROOT/.mcp.json" ]; then
    bash "$SCRIPT_DIR/sync_claude_mcp.sh" check || true
  else
    log_warn "缺失 .mcp.json"
    failed=1
  fi

  if [ "$failed" -ne 0 ]; then
    log_error "健康检查未全部通过，请根据上方输出修复"
    return 1
  fi
  log_info "健康检查通过"
}

run_init() {
  log_info "仓库: $REPO_ROOT"
  log_info "开始 AI 工具链初始化 ..."

  if [ "${SKIP_GENTLE_AI:-0}" != "1" ]; then
    setup_gentle_ai
  else
    log_info "跳过 Gentle-AI（SKIP_GENTLE_AI=1）"
  fi

  if [ "${SKIP_CODEGRAPH:-0}" != "1" ]; then
    setup_codegraph
  else
    log_info "跳过 CodeGraph（SKIP_CODEGRAPH=1）"
  fi

  if [ "${SKIP_ENGRAM:-0}" != "1" ]; then
    setup_engram
  else
    log_info "跳过 Engram（SKIP_ENGRAM=1）"
  fi

  setup_claude_mcp

  echo ""
  log_info "初始化完成。请完全重启 Cursor（Cursor + Claude Code 插件均需重载 MCP）。"
  log_info "Claude Code 插件：终端执行 make sync-claude-mcp-check，REPL 内 /mcp reconnect all"
  log_info "Claude Code 插件首次打开本项目可执行 /sdd-init"
  log_info "验证: make setup-ai-check"
  log_info "文档: docs/AI工具链上手.md"
}

main() {
  local cmd="${1:-init}"
  if [ $# -gt 0 ]; then
    shift
  fi

  case "$cmd" in
    init|setup|all)
      run_init "$@"
      ;;
    check|doctor|status)
      run_check "$@"
      ;;
    gentle-ai|gentleai)
      setup_gentle_ai "$@"
      ;;
    codegraph|cg)
      setup_codegraph "$@"
      ;;
    engram|memory)
      setup_engram "$@"
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
