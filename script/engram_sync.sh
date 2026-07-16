#!/usr/bin/env bash
# Engram 跨设备记忆同步：导出/导入 .engram/，可选与 git 联动。
#
# 用法:
#   ./script/engram_sync.sh export          # 导出记忆到 .engram/
#   ./script/engram_sync.sh import          # 从 .engram/ 导入记忆
#   ./script/engram_sync.sh status          # 查看同步状态
#   ./script/engram_sync.sh push [-m msg]   # 导出 + git add/commit .engram/
#   ./script/engram_sync.sh pull            # git pull + 导入
#   ./script/engram_sync.sh pull --no-git   # 仅导入（不拉代码）
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
ENGRAM_DIR="$REPO_ROOT/.engram"
DEFAULT_COMMIT_MSG="chore: sync engram memories"

log_info() { echo "[INFO] $*"; }
log_warning() { echo "[WARN] $*" >&2; }
log_error() { echo "[ERROR] $*" >&2; }

usage() {
  cat <<'EOF'
Engram 跨设备记忆同步

子命令:
  export              导出本项目记忆到 .engram/
  import              从 .engram/ 导入记忆
  status              查看 engram sync 状态
  push [-m "message"] 导出后 git add/commit .engram/（有变更才提交）
  pull [--no-git]     git pull 后导入（--no-git 跳过 git pull）

示例:
  make engram-sync-push
  make engram-sync-pull
  ./script/engram_sync.sh push -m "chore: sync engram after refactor"
EOF
}

require_engram() {
  if ! command -v engram &>/dev/null; then
    log_error "未找到 engram 命令，请先安装: https://github.com/Gentleman-Programming/engram"
    exit 1
  fi
}

require_git_repo() {
  if ! git -C "$REPO_ROOT" rev-parse --is-inside-work-tree &>/dev/null; then
    log_error "当前目录不是 git 仓库: $REPO_ROOT"
    exit 1
  fi
}

run_export() {
  require_engram
  cd "$REPO_ROOT"
  log_info "导出 Engram 记忆到 $ENGRAM_DIR ..."
  engram sync
  log_info "导出完成"
  engram sync --status || true
}

run_import() {
  require_engram
  cd "$REPO_ROOT"
  if [ ! -d "$ENGRAM_DIR" ]; then
    log_warning ".engram/ 不存在，跳过导入（请先在本机执行 export/push，或 git pull）"
    exit 0
  fi
  log_info "从 $ENGRAM_DIR 导入 Engram 记忆 ..."
  engram sync --import
  log_info "导入完成"
  engram sync --status || true
}

run_status() {
  require_engram
  cd "$REPO_ROOT"
  engram sync --status
}

run_push() {
  local commit_msg="$DEFAULT_COMMIT_MSG"
  while [ $# -gt 0 ]; do
    case "$1" in
      -m|--message)
        shift
        commit_msg="${1:-}"
        if [ -z "$commit_msg" ]; then
          log_error "push: -m/--message 需要提交说明"
          exit 1
        fi
        ;;
      -h|--help)
        usage
        exit 0
        ;;
      *)
        log_error "未知参数: $1"
        usage
        exit 1
        ;;
    esac
    shift
  done

  require_engram
  require_git_repo
  run_export

  cd "$REPO_ROOT"
  if [ ! -d "$ENGRAM_DIR" ]; then
    log_warning ".engram/ 不存在，无内容可提交"
    exit 0
  fi

  git add "$ENGRAM_DIR"
  if git diff --cached --quiet -- "$ENGRAM_DIR"; then
    log_info ".engram/ 无变更，跳过 commit"
    exit 0
  fi

  git commit -m "$commit_msg"
  log_info "已提交 .engram/，请执行 git push 同步到远端"
}

run_pull() {
  local skip_git=0
  while [ $# -gt 0 ]; do
    case "$1" in
      --no-git)
        skip_git=1
        ;;
      -h|--help)
        usage
        exit 0
        ;;
      *)
        log_error "未知参数: $1"
        usage
        exit 1
        ;;
    esac
    shift
  done

  require_engram
  if [ "$skip_git" -eq 0 ]; then
    require_git_repo
    cd "$REPO_ROOT"
    log_info "拉取远端更新 ..."
    git pull --rebase --autostash
  fi
  run_import
}

main() {
  local cmd="${1:-export}"
  if [ $# -gt 0 ]; then
    shift
  fi

  case "$cmd" in
    export|sync)
      run_export "$@"
      ;;
    import)
      run_import "$@"
      ;;
    status)
      run_status "$@"
      ;;
    push)
      run_push "$@"
      ;;
    pull)
      run_pull "$@"
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
