#!/bin/bash

# ============================================
# admin-system 项目通用工具函数库
# ============================================

# 颜色定义
ERROR_COLOR="\033[31m"      # 错误消息（红色）
SUCCESS_COLOR="\033[32m"    # 成功消息（绿色）
WARNING_COLOR="\033[33m"    # 告警消息（黄色）
INFO_COLOR="\033[36m"       # 信息消息（青色）
PLAIN_TEXT_COLOR='\033[0m'  # 重置颜色

# 日志输出函数
log_info() {
  echo -e "${SUCCESS_COLOR}[INFO]${PLAIN_TEXT_COLOR} $1"
}

log_error() {
  echo -e "${ERROR_COLOR}[ERROR]${PLAIN_TEXT_COLOR} $1" >&2
}

log_warning() {
  echo -e "${WARNING_COLOR}[WARN]${PLAIN_TEXT_COLOR} $1"
}

log_debug() {
  echo -e "${INFO_COLOR}[DEBUG]${PLAIN_TEXT_COLOR} $1"
}

# 获取脚本所在目录
get_script_dir() {
  local script_path="${BASH_SOURCE[0]}"
  if [ -z "$script_path" ]; then
    script_path="$0"
  fi
  local script_dir=$(cd "$(dirname "$script_path")" 2>/dev/null && pwd)
  if [ -z "$script_dir" ]; then
    # 如果 cd 失败，尝试使用 dirname
    script_dir="$(dirname "$script_path")"
    # 转换为绝对路径
    if [[ "$script_dir" != /* ]]; then
      script_dir="$(cd "$script_dir" 2>/dev/null && pwd || echo "$script_dir")"
    fi
  fi
  echo "$script_dir"
}

# 获取项目根目录
get_project_root() {
  local script_dir=$(get_script_dir)
  local project_root=$(dirname "$script_dir" 2>/dev/null)
  # 确保路径存在
  if [ ! -d "$project_root" ]; then
    # 尝试使用当前工作目录
    project_root="$(pwd)"
  fi
  echo "$project_root"
}

# 检查命令是否存在
check_command() {
  if ! command -v "$1" &> /dev/null; then
    log_error "命令 '$1' 未找到，请先安装"
    return 1
  fi
  return 0
}

# 检查目录是否存在，不存在则创建
ensure_dir() {
  local dir_path="$1"
  if [ ! -d "$dir_path" ]; then
    mkdir -p "$dir_path"
    log_info "创建目录: $dir_path"
  fi
}

# 检查文件是否存在
check_file() {
  local file_path="$1"
  if [ ! -f "$file_path" ]; then
    log_error "文件不存在: $file_path"
    return 1
  fi
  return 0
}

# 获取 Go 进程 PID
get_go_pid() {
  local app_name="${1:-admin-server}"
  pgrep -f "$app_name" | head -n 1
}

# 检查服务是否运行
is_service_running() {
  local pid=$(get_go_pid "$1")
  if [ -n "$pid" ]; then
    return 0
  else
    return 1
  fi
}

# 项目路径配置
PROJECT_ROOT=$(get_project_root)
ADMIN_SERVER_DIR="$PROJECT_ROOT/admin-server"
ADMIN_FRONTEND_DIR="$PROJECT_ROOT/admin-frontend"
SCRIPT_DIR=$(get_script_dir)

# 日志目录配置
LOG_DIR="$PROJECT_ROOT/logs"
SERVER_LOG_DIR="$LOG_DIR/server"
FRONTEND_LOG_DIR="$LOG_DIR/frontend"

# 确保日志目录存在
ensure_dir "$SERVER_LOG_DIR"
ensure_dir "$FRONTEND_LOG_DIR"

