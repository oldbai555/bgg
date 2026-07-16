#!/bin/bash
# bgg-dev（8.138.33.161）一键部署脚本：拉最新代码 → 重新构建改动的服务镜像 → 重启容器。
# 只负责 admin-server 的 6 个 docker 服务（gateway/iam/task/sdk/chat/content）。
# 前端（admin-frontend）走独立的构建+上传流程，见 script/deploy-frontend.sh。
#
# 用法（在服务器上，仓库 clone 目录内执行，即 /home/work/src-bgg）：
#   bash script/deploy-dev.sh              # 全量重新 build + up 六个服务
#   bash script/deploy-dev.sh gateway iam  # 只重新 build + up 指定服务
#
# 前置条件：
#   - /home/work/src-bgg/admin-server/.env 已按 .env.dev-mixed.example 配好真实值
#   - Docker 已安装，且 /etc/docker/daemon.json 配置了国内镜像加速（见 docs/changelog/）
#   - 宿主机 /etc/redis/redis.conf 的 bind 已加上 docker 网桥地址（172.17.0.1）
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(dirname "$SCRIPT_DIR")"
ADMIN_SERVER_DIR="$REPO_ROOT/admin-server"
COMPOSE_FILE="docker-compose.dev-mixed.yml"

log() { echo -e "\n\033[36m[DEPLOY]\033[0m $1"; }

if [ ! -f "$ADMIN_SERVER_DIR/.env" ]; then
  echo "未找到 $ADMIN_SERVER_DIR/.env，请先参考 .env.dev-mixed.example 创建"
  exit 1
fi

log "拉取最新代码（当前分支: $(git -C "$REPO_ROOT" branch --show-current)）"
git -C "$REPO_ROOT" pull

cd "$ADMIN_SERVER_DIR"

# DOCKER_CONFIG 指向可写目录，避免 ubuntu 用户 $HOME（/home/work，root 属主）导致
# docker CLI 建 ~/.docker 配置目录时权限不足
export DOCKER_CONFIG="${DOCKER_CONFIG:-/tmp/docker-config}"
mkdir -p "$DOCKER_CONFIG"

SERVICES="$*"
if [ -z "$SERVICES" ]; then
  log "构建全部六个服务镜像"
else
  log "构建指定服务镜像: $SERVICES"
fi
docker compose -f "$COMPOSE_FILE" --env-file .env build $SERVICES

log "重启服务"
docker compose -f "$COMPOSE_FILE" --env-file .env up -d $SERVICES

log "当前状态"
docker compose -f "$COMPOSE_FILE" ps

log "健康检查"
sleep 3
curl -s -o /dev/null -w "gateway /api/v1/ping -> %{http_code}\n" http://127.0.0.1:20000/api/v1/ping || echo "gateway 健康检查失败，请查看日志: docker logs admin-server-gateway-1"

log "部署完成"
