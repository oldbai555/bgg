#!/bin/bash

# ============================================
# admin-system 统一管理脚本
# 功能：开发、构建、部署、Supervisor管理
# ============================================

# set -e  # 注释掉，避免在 Windows 环境下因路径问题导致脚本静默退出

# 颜色定义
ERROR_COLOR="\033[31m"
SUCCESS_COLOR="\033[32m"
WARNING_COLOR="\033[33m"
INFO_COLOR="\033[36m"
PLAIN_TEXT_COLOR='\033[0m'

# 日志函数
log_info() { echo -e "${SUCCESS_COLOR}[INFO]${PLAIN_TEXT_COLOR} $1"; }
log_error() { echo -e "${ERROR_COLOR}[ERROR]${PLAIN_TEXT_COLOR} $1" >&2; }
log_warning() { echo -e "${WARNING_COLOR}[WARN]${PLAIN_TEXT_COLOR} $1"; }
log_debug() { echo -e "${INFO_COLOR}[DEBUG]${PLAIN_TEXT_COLOR} $1"; }

# 获取脚本目录和项目根目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
ADMIN_SERVER_DIR="$PROJECT_ROOT/admin-server"
ADMIN_FRONTEND_DIR="$PROJECT_ROOT/admin-frontend"

# Supervisor 配置
SUPERVISOR_DIR="${SUPERVISOR_DIR:-/home/work/service}"
SUPERVISOR_LOG_DIR="${SUPERVISOR_LOG_DIR:-/home/work/supervisor/logs}"
SUPERVISOR_CONF_DIR="${SUPERVISOR_CONF_DIR:-/etc/supervisor/conf.d}"
PACKAGE_OUTPUT_DIR="$PROJECT_ROOT/package"

# 应用配置
APP_NAME="admin-server"
APP_BINARY="admin-server"
APP_CONFIG_FILE="admin-api.yaml"

# 工具函数
ensure_dir() {
  [ ! -d "$1" ] && mkdir -p "$1" && log_info "创建目录: $1"
}

check_command() {
  command -v "$1" &> /dev/null || { log_error "命令 '$1' 未找到"; return 1; }
}

get_go_pid() {
  pgrep -f "$APP_NAME" | head -n 1
}

is_service_running() {
  [ -n "$(get_go_pid)" ]
}

# ============================================
# 开发工具
# ============================================

# 检查服务健康状态
check_service_health() {
  local max_attempts=30
  local attempt=0
  local health_url="http://localhost:8888/api/v1/ping"
  
  log_info "等待服务启动..."
  
  while [ $attempt -lt $max_attempts ]; do
    # 检查进程是否存在
    if ! is_service_running; then
      log_warning "进程未运行，等待中... ($attempt/$max_attempts)"
      sleep 1
      attempt=$((attempt + 1))
      continue
    fi
    
    # 检查端口是否监听
    if command -v netstat &> /dev/null; then
      if ! netstat -tuln 2>/dev/null | grep -q ":8888.*LISTEN"; then
        log_warning "端口未监听，等待中... ($attempt/$max_attempts)"
        sleep 1
        attempt=$((attempt + 1))
        continue
      fi
    elif command -v ss &> /dev/null; then
      if ! ss -tuln 2>/dev/null | grep -q ":8888"; then
        log_warning "端口未监听，等待中... ($attempt/$max_attempts)"
        sleep 1
        attempt=$((attempt + 1))
        continue
      fi
    fi
    
    # 检查健康检查接口
    if command -v curl &> /dev/null; then
      local response=$(curl -s -o /dev/null -w "%{http_code}" --max-time 2 "$health_url" 2>/dev/null || echo "000")
      if [ "$response" = "200" ]; then
        log_info "服务健康检查通过"
        return 0
      fi
    elif command -v wget &> /dev/null; then
      local response=$(wget -q -O /dev/null -S --timeout=2 "$health_url" 2>&1 | grep -o "HTTP/[0-9.]* [0-9]*" | awk '{print $2}' || echo "000")
      if [ "$response" = "200" ]; then
        log_info "服务健康检查通过"
        return 0
      fi
    else
      # 如果没有 curl 或 wget，至少检查进程和端口
      if is_service_running; then
        log_info "服务进程运行中（无法进行健康检查，请手动验证）"
        return 0
      fi
    fi
    
    log_warning "健康检查未通过，等待中... ($attempt/$max_attempts)"
    sleep 1
    attempt=$((attempt + 1))
  done
  
  log_error "服务启动超时（${max_attempts}秒）"
  return 1
}

dev_start() {
  log_info "启动后端服务..."
  
  # 检查服务是否已在运行
  if is_service_running; then
    if check_service_health; then
      log_warning "服务已在运行中"
      return 0
    else
      log_warning "服务进程存在但未正常响应，尝试重启..."
      dev_stop
      sleep 1
    fi
  fi
  
  cd "$ADMIN_SERVER_DIR" || exit 1
  [ ! -f "etc/$APP_CONFIG_FILE" ] && { log_error "配置文件不存在: etc/$APP_CONFIG_FILE"; return 1; }
  
  # 确保日志目录存在
  ensure_dir "$PROJECT_ROOT/logs"
  
  # 启动服务
  log_info "正在启动服务..."
  nohup go run admin.go > "$PROJECT_ROOT/logs/server.log" 2>&1 &
  local pid=$!
  
  # 等待进程启动
  sleep 1
  
  # 检查进程是否启动
  if ! kill -0 "$pid" 2>/dev/null; then
    log_error "服务启动失败，请查看日志: $PROJECT_ROOT/logs/server.log"
    tail -n 20 "$PROJECT_ROOT/logs/server.log" 2>/dev/null || true
    return 1
  fi
  
  # 等待服务完全启动（健康检查）
  if check_service_health; then
    log_info "服务启动成功，PID: $pid"
    log_info "健康检查接口: http://localhost:8888/api/v1/ping"
    log_info "日志文件: $PROJECT_ROOT/logs/server.log"
    return 0
  else
    log_error "服务启动失败，请查看日志: $PROJECT_ROOT/logs/server.log"
    tail -n 30 "$PROJECT_ROOT/logs/server.log" 2>/dev/null || true
    # 尝试停止失败的进程
    kill "$pid" 2>/dev/null || true
    return 1
  fi
}

dev_stop() {
  log_info "停止后端服务..."
  local pid=$(get_go_pid)
  [ -z "$pid" ] && { log_warning "服务未运行"; return 0; }
  kill "$pid" 2>/dev/null || kill -9 "$pid" 2>/dev/null
  log_info "服务已停止"
}

dev_status() {
  local pid=$(get_go_pid)
  [ -n "$pid" ] && { log_info "服务运行中，PID: $pid"; ps -p "$pid" -o pid,cmd,etime; } || log_warning "服务未运行"
}

dev_logs() {
  local lines="${1:-100}"
  [ -f "$PROJECT_ROOT/logs/server.log" ] && tail -n "$lines" "$PROJECT_ROOT/logs/server.log" || log_error "日志文件不存在"
}

# ============================================
# 构建工具
# ============================================

build_server() {
  log_info "构建后端服务..."
  cd "$ADMIN_SERVER_DIR" || exit 1
  [ ! -f "go.mod" ] && { log_error "go.mod 不存在"; return 1; }
  
  # 获取 git 提交版本号
  local git_version=$(git rev-parse --short HEAD 2>/dev/null || echo "dev")
  local output_dir="dist"
  ensure_dir "$output_dir"
  
  log_info "Git版本号: $git_version"
  
  # 构建时传递版本号
  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags "-s -w -X main.GIT_COMMIT_VERSION=$git_version" \
    -o "$output_dir/$APP_BINARY" \
    admin.go
  
  [ $? -eq 0 ] && {
    log_info "后端构建成功"
    [ -f "etc/$APP_CONFIG_FILE" ] && cp "etc/$APP_CONFIG_FILE" "$output_dir/" && log_info "已复制配置文件"
    chmod +x "$output_dir/$APP_BINARY"
    
    # 验证版本号是否成功注入
    log_info "构建完成，版本号已注入: $git_version"
  } || { log_error "后端构建失败"; return 1; }
}

build_frontend() {
  log_info "构建前端项目..."
  cd "$ADMIN_FRONTEND_DIR" || exit 1
  [ ! -f "package.json" ] && { log_error "package.json 不存在"; return 1; }
  
  [ ! -d "node_modules" ] && {
    log_info "安装依赖..."
    command -v pnpm &> /dev/null && pnpm install || npm install
  }
  
  log_info "开始构建..."
  command -v pnpm &> /dev/null && pnpm build || npm run build
  [ $? -eq 0 ] && log_info "前端构建成功" || { log_error "前端构建失败"; return 1; }
}

package_server() {
  log_info "打包后端服务..."
  local version=$(git rev-parse --short HEAD 2>/dev/null || echo "dev")
  ensure_dir "$PACKAGE_OUTPUT_DIR"
  
  build_server || return 1
  
  local package_file="$PACKAGE_OUTPUT_DIR/${APP_NAME}_${version}.tar.gz"
  cd "$ADMIN_SERVER_DIR/dist" || return 1
  tar -czf "$package_file" ./*
  log_info "打包完成: $package_file"
  echo "$package_file"
}

package_frontend() {
  log_info "打包前端项目..."
  local version=$(git rev-parse --short HEAD 2>/dev/null || echo "dev")
  ensure_dir "$PACKAGE_OUTPUT_DIR"
  
  build_frontend || return 1
  
  local package_file="$PACKAGE_OUTPUT_DIR/admin-frontend_${version}.tar.gz"
  cd "$ADMIN_FRONTEND_DIR" || return 1
  tar -czf "$package_file" -C dist .
  log_info "打包完成: $package_file"
  echo "$package_file"
}

# ============================================
# Supervisor 管理
# ============================================

supervisor_gen_conf() {
  local app_name="${1:-$APP_NAME}"
  local output_dir="${2:-$PACKAGE_OUTPUT_DIR/$app_name}"
  local service_dir="$SUPERVISOR_DIR/$app_name"
  
  ensure_dir "$output_dir"
  log_info "生成 Supervisor 配置: $app_name.conf"
  
  cat > "$output_dir/$app_name.conf" <<EOF
[program:$app_name]
directory=$service_dir
command=$service_dir/$APP_BINARY -f $service_dir/$APP_CONFIG_FILE
autostart=true
autorestart=true
startsecs=10
startretries=3
user=root
redirect_stderr=true
stdout_logfile=$SUPERVISOR_LOG_DIR/${app_name}_stdout.log
stderr_logfile=$SUPERVISOR_LOG_DIR/${app_name}_stderr.log
stdout_logfile_maxbytes=20MB
stdout_logfile_backups=20
environment=GOMAXPROCS=2
EOF
  
  log_info "配置文件已生成: $output_dir/$app_name.conf"
}

supervisor_deploy() {
  local package_file="$1"
  local app_name="${2:-$APP_NAME}"
  
  [ -z "$package_file" ] && { log_error "请指定打包文件"; return 1; }
  [ ! -f "$package_file" ] && { log_error "打包文件不存在: $package_file"; return 1; }
  
  log_info "部署服务: $app_name"
  ensure_dir "$SUPERVISOR_DIR/$app_name"
  
  tar -xzf "$package_file" -C "$SUPERVISOR_DIR/$app_name"
  chmod +x "$SUPERVISOR_DIR/$app_name/$APP_BINARY"
  
  [ -f "$SUPERVISOR_DIR/$app_name/$app_name.conf" ] || {
    supervisor_gen_conf "$app_name" "$SUPERVISOR_DIR/$app_name"
  }
  cp "$SUPERVISOR_DIR/$app_name/$app_name.conf" "$SUPERVISOR_CONF_DIR/"
  
  command -v supervisorctl &> /dev/null && {
    supervisorctl update
    supervisorctl restart "$app_name"
    sleep 2
    supervisorctl status "$app_name"
  } || log_warning "未找到 supervisorctl，请手动执行"
  
  log_info "部署完成"
}

supervisor_install() {
  local app_name="${1:-$APP_NAME}"
  local version=$(git rev-parse --short HEAD 2>/dev/null || echo "dev")
  
  log_info "安装服务: $app_name"
  local package_file=$(package_server)
  [ -z "$package_file" ] && return 1
  supervisor_deploy "$package_file" "$app_name"
}

supervisor_manage() {
  local action="$1"
  local app_name="${2:-$APP_NAME}"
  
  command -v supervisorctl &> /dev/null || { log_error "未找到 supervisorctl"; return 1; }
  
  case "$action" in
    status) supervisorctl status "$app_name" ;;
    start) supervisorctl start "$app_name" && supervisorctl status "$app_name" ;;
    stop) supervisorctl stop "$app_name" && supervisorctl status "$app_name" ;;
    restart) supervisorctl restart "$app_name" && supervisorctl status "$app_name" ;;
    logs) [ -f "$SUPERVISOR_LOG_DIR/${app_name}_stdout.log" ] && tail -n 100 "$SUPERVISOR_LOG_DIR/${app_name}_stdout.log" || log_warning "日志文件不存在" ;;
    *) log_error "未知操作: $action"; return 1 ;;
  esac
}

# ============================================
# 主程序
# ============================================

usage() {
  cat <<EOF
用法: ./admin.sh <command> [options]

开发命令（开发环境建议使用 GoLand，此部分可选）:
  dev start         启动后端服务（带健康检查）
  dev stop          停止后端服务
  dev status        查看服务状态
  dev logs [行数]   查看日志

构建命令:
  build server      构建后端
  build frontend    构建前端
  package server    打包后端（构建+打包）
  package frontend  打包前端（构建+打包）

Supervisor 命令:
  supervisor gen-conf          生成配置文件
  supervisor install           安装服务（构建+部署）
  supervisor deploy <file>     部署打包文件
  supervisor status            查看服务状态
  supervisor start             启动服务
  supervisor stop              停止服务
  supervisor restart           重启服务
  supervisor logs              查看日志

示例:
  ./admin.sh dev start
  ./admin.sh build server
  ./admin.sh package server
  ./admin.sh supervisor install
EOF
}

main() {
  # 确保目录存在（忽略错误，避免在 Windows 环境下失败）
  ensure_dir "$PROJECT_ROOT/logs" 2>/dev/null || true
  ensure_dir "$PACKAGE_OUTPUT_DIR" 2>/dev/null || true
  ensure_dir "$SUPERVISOR_LOG_DIR" 2>/dev/null || true
  
  case "${1:-}" in
    dev)
      case "${2:-}" in
        start) dev_start ;;
        stop) dev_stop ;;
        status) dev_status ;;
        logs) dev_logs "${3:-100}" ;;
        *) log_error "未知命令: dev $2"; usage; exit 1 ;;
      esac
      ;;
    build)
      case "${2:-}" in
        server) build_server ;;
        frontend) build_frontend ;;
        *) log_error "未知命令: build $2"; usage; exit 1 ;;
      esac
      ;;
    package)
      case "${2:-}" in
        server) package_server ;;
        frontend) package_frontend ;;
        *) log_error "未知命令: package $2"; usage; exit 1 ;;
      esac
      ;;
    supervisor)
      case "${2:-}" in
        gen-conf) supervisor_gen_conf ;;
        install) supervisor_install ;;
        deploy) supervisor_deploy "${3:-}" ;;
        status|start|stop|restart|logs) supervisor_manage "$2" ;;
        *) log_error "未知命令: supervisor $2"; usage; exit 1 ;;
      esac
      ;;
    *)
      usage
      exit 1
      ;;
  esac
}

main "$@"

