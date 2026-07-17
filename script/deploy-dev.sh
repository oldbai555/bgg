#!/bin/bash
# bgg-dev 一键部署脚本：拉最新代码（仅同步 compose/配置文件，不再编译源码）
# → 拉取 ghcr.io 镜像 → 重启容器。
# 只负责 admin-server 的 6 个 docker 服务（gateway/iam/task/sdk/chat/content）。
# 前端（admin-frontend）走独立的构建+上传流程，见 script/deploy-frontend.sh。
#
# 用法（在服务器上，仓库 clone 目录内执行，即 /home/work/src-bgg）：
#   TAG=<git-sha> bash script/deploy-dev.sh              # 全量 pull + up 六个服务
#   TAG=<git-sha> bash script/deploy-dev.sh gateway iam  # 只 pull + up 指定服务
#   不传 TAG 时默认取 main 分支 pull 之后的最新 commit sha（见下方逻辑），
#   前提是该 sha 对应的 CI build-images 已经跑成功、镜像已推送到 ghcr.io
#   （可以去 https://github.com/oldbai555?tab=packages 确认 tag 是否存在）。
#
# 前置条件：
#   - /home/work/src-bgg/admin-server/.env 已按 .env.dev-mixed.example 配好真实值
#   - Docker 已安装，且 /etc/docker/daemon.json 配置了国内镜像加速（见 docs/changelog/）
#   - 宿主机 /etc/redis/redis.conf 的 bind 已加上 docker 网桥地址（172.17.0.1）
#   - ghcr.io 六个 package 可见性为 public（当前已确认），pull 不需要 docker login；
#     若之后改成 private，需要在这里补一步 `docker login ghcr.io`
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

CURRENT_BRANCH="$(git -C "$REPO_ROOT" branch --show-current)"
if [ "$CURRENT_BRANCH" != "main" ]; then
  echo "当前分支是 $CURRENT_BRANCH，不是 main——ghcr 镜像只在 push 到 main 时由 CI 构建，"
  echo "该分支的 commit sha 在 ghcr 上不会有对应 tag。请先 git checkout main。"
  exit 1
fi

log "拉取最新代码（main 分支）"
git -C "$REPO_ROOT" pull origin main

export TAG="${TAG:-$(git -C "$REPO_ROOT" rev-parse HEAD)}"
log "使用镜像 tag: $TAG（可去 https://github.com/oldbai555?tab=packages 确认该 tag 已推送）"

cd "$ADMIN_SERVER_DIR"

# DOCKER_CONFIG 指向可写目录，避免 ubuntu 用户 $HOME（/home/work，root 属主）导致
# docker CLI 建 ~/.docker 配置目录时权限不足
export DOCKER_CONFIG="${DOCKER_CONFIG:-/tmp/docker-config}"
mkdir -p "$DOCKER_CONFIG"

SERVICES="$*"
if [ -z "$SERVICES" ]; then
  log "拉取全部六个服务镜像"
else
  log "拉取指定服务镜像: $SERVICES"
fi
docker compose -f "$COMPOSE_FILE" --env-file .env pull $SERVICES

log "重启服务"
docker compose -f "$COMPOSE_FILE" --env-file .env up -d $SERVICES

log "当前状态"
docker compose -f "$COMPOSE_FILE" ps

log "健康检查"
sleep 3
curl -s -o /dev/null -w "gateway /api/v1/ping -> %{http_code}\n" http://127.0.0.1:20000/api/v1/ping || echo "gateway 健康检查失败，请查看日志: docker logs admin-server-gateway-1"

log "部署完成"
