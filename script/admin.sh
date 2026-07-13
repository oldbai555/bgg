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
  # 打包 dist 目录下的所有文件（不包含 dist 这一层），方便在服务器上直接解压到 dist 目录中
  tar -czf "$package_file" -C dist .
  log_info "打包完成: $package_file"
  echo "$package_file"
}

# ============================================
# 前端部署（静态文件）
# ============================================

deploy_frontend() {
  local package_file="$1"

  [ -z "$package_file" ] && { log_error "请指定前端打包文件"; return 1; }
  [ ! -f "$package_file" ] && { log_error "前端打包文件不存在: $package_file"; return 1; }

  # 使用当前目录作为前端部署目录
  local target_dir="$(pwd)"
  log_info "部署前端到目录: $target_dir"

  cd "$target_dir" || return 1

  # 备份旧的 dist 目录
  if [ -d "dist" ]; then
    local backup_dir="dist_$(date +%Y%m%d%H%M%S).bak"
    mv dist "$backup_dir"
    log_info "已备份旧的 dist 目录为: $backup_dir"
  fi

  # 创建新的 dist 目录并解压到其中
  mkdir -p dist
  tar -xzf "$package_file" -C dist

  log_info "前端部署完成，访问目录: $target_dir/dist"
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
# docker-compose 管理（Phase 3 新增，与 Supervisor 模式并存，见
# admin-server/docs/21-cd-and-deployment.md 第 5 节）
# ============================================

COMPOSE_REMOTE_HOST="${COMPOSE_REMOTE_HOST:-}"      # 生产服务器 SSH host（~/.ssh/config 别名或 user@ip）
COMPOSE_REMOTE_DIR="${COMPOSE_REMOTE_DIR:-/home/work/admin-compose}"
COMPOSE_FILE="docker-compose.prod.yml"

compose_pull() {
  local host="${1:-$COMPOSE_REMOTE_HOST}"
  [ -z "$host" ] && { log_error "请指定远程主机（COMPOSE_REMOTE_HOST 或参数）"; return 1; }

  log_info "在 $host 拉取最新镜像..."
  ssh "$host" "cd $COMPOSE_REMOTE_DIR && docker compose -f $COMPOSE_FILE pull"
}

compose_deploy() {
  local host="${1:-$COMPOSE_REMOTE_HOST}"
  [ -z "$host" ] && { log_error "请指定远程主机（COMPOSE_REMOTE_HOST 或参数）"; return 1; }

  compose_pull "$host" || return 1

  log_info "在 $host 应用新镜像..."
  ssh "$host" "cd $COMPOSE_REMOTE_DIR && docker compose -f $COMPOSE_FILE up -d"

  log_info "部署完成，检查状态..."
  compose_status "$host"
}

compose_status() {
  local host="${1:-$COMPOSE_REMOTE_HOST}"
  [ -z "$host" ] && { log_error "请指定远程主机"; return 1; }
  ssh "$host" "cd $COMPOSE_REMOTE_DIR && docker compose -f $COMPOSE_FILE ps"
}

compose_logs() {
  local host="${1:-$COMPOSE_REMOTE_HOST}"
  local service="${2:-}"
  [ -z "$host" ] && { log_error "请指定远程主机"; return 1; }
  ssh "$host" "cd $COMPOSE_REMOTE_DIR && docker compose -f $COMPOSE_FILE logs --tail=200 ${service}"
}

# ============================================
# 主程序
# ============================================

usage() {
  cat <<EOF
用法: ./admin.sh <command> [subcommand] [options]

基础说明:
  建议使用 bash 执行本脚本（避免 sh 在某些系统上的兼容问题）：
    bash script/admin.sh <command> ...

构建命令:
  build server              构建后端二进制到 admin-server/dist
  build frontend            构建前端到 admin-frontend/dist

打包命令:
  package server            构建并打包后端，输出到 package/admin-server_<git>.tar.gz
  package frontend          构建并打包前端，输出到 package/admin-frontend_<git>.tar.gz

前端部署命令（静态文件）:
  frontend deploy <file>   在当前目录下部署前端，将前端包内容解压到 ./dist
                          注意：请在 Nginx 对应的 Web 根目录下执行该命令

Supervisor 命令:
  supervisor gen-conf        仅生成 Supervisor 配置文件（默认使用 \$APP_NAME）
  supervisor install         一键构建 + 打包 + 部署到 Supervisor 目录
  supervisor deploy <file>   使用已有的后端包部署到 Supervisor 目录
  supervisor status          查看服务状态
  supervisor start           启动服务
  supervisor stop            停止服务
  supervisor restart         重启服务
  supervisor logs            查看最近 100 行 stdout 日志

docker-compose 命令（Phase 3 新增，六个容器镜像部署，与 Supervisor 模式并存）:
  compose pull [host]          在远程主机拉取最新镜像（docker-compose.prod.yml）
  compose deploy [host]        拉取镜像 + docker compose up -d + 检查状态
  compose status [host]        查看远程主机上六个服务的容器状态
  compose logs [host] [svc]    查看远程主机指定服务的最近 200 行日志（svc 留空查全部）
  远程主机可通过参数传入，或设置环境变量 COMPOSE_REMOTE_HOST；
  远程 compose 目录默认 /home/work/admin-compose，可通过 COMPOSE_REMOTE_DIR 覆盖。

示例:
  # 本地构建
  bash script/admin.sh build server
  bash script/admin.sh build frontend

  # 本地打包
  bash script/admin.sh package server
  bash script/admin.sh package frontend

  # 服务器上部署后端（基于打包文件）
  bash script/admin.sh supervisor deploy package/admin-server_xxx.tar.gz

  # 服务器上部署前端（在 Nginx Web 根目录下执行）
  bash script/admin.sh frontend deploy /path/to/package/admin-frontend_xxx.tar.gz

  # 使用 Supervisor 管理服务
  bash script/admin.sh supervisor status
  bash script/admin.sh supervisor restart

  # 使用 docker-compose 部署（生产服务器已配置 docker-compose.prod.yml + ghcr 登录）
  bash script/admin.sh compose deploy myserver
  COMPOSE_REMOTE_HOST=myserver bash script/admin.sh compose status
  bash script/admin.sh compose logs myserver iam
EOF
}

main() {
  # 确保目录存在（忽略错误，避免在 Windows 环境下失败）
  ensure_dir "$PROJECT_ROOT/logs" 2>/dev/null || true
  ensure_dir "$PACKAGE_OUTPUT_DIR" 2>/dev/null || true
  ensure_dir "$SUPERVISOR_LOG_DIR" 2>/dev/null || true
  
  case "${1:-}" in
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
    frontend)
      case "${2:-}" in
        deploy) deploy_frontend "${3:-}" ;;
        *) log_error "未知命令: frontend $2"; usage; exit 1 ;;
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
    compose)
      case "${2:-}" in
        pull) compose_pull "${3:-}" ;;
        deploy) compose_deploy "${3:-}" ;;
        status) compose_status "${3:-}" ;;
        logs) compose_logs "${3:-}" "${4:-}" ;;
        *) log_error "未知命令: compose $2"; usage; exit 1 ;;
      esac
      ;;
    *)
      usage
      exit 1
      ;;
  esac
}

main "$@"

