#!/bin/bash
# admin-frontend 一键部署到 bgg-dev：本机构建 → 打包 → 上传 → 远端解压替换 dist。
# 必须在本机（有 Node/pnpm 的开发机）执行，服务器上没有装 Node 工具链。
# CI 同源逻辑见 .github/workflows/ci.yml 的 deploy-frontend job（Runner 上构建，无 sudo）。
#
# 前提：远端目录属主为 SSH 登录用户（推荐 ubuntu）：
#   sudo chown -R ubuntu:ubuntu /home/work/service/admin-frontend
#
# 用法: bash script/deploy-frontend.sh [ssh-host]
#   ssh-host 默认 bgg-dev（见 ~/.ssh/config 别名，script/ssh-setup.sh 生成）
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(dirname "$SCRIPT_DIR")"
SSH_HOST="${1:-bgg-dev}"
REMOTE_FRONTEND_DIR="/home/work/service/admin-frontend"

log() { echo -e "\n\033[36m[DEPLOY]\033[0m $1"; }

log "本机构建 + 打包前端"
PACKAGE_FILE=$(bash "$SCRIPT_DIR/admin.sh" package frontend | tail -1)
if [ ! -f "$PACKAGE_FILE" ]; then
  echo "打包失败，未找到产物: $PACKAGE_FILE"
  exit 1
fi
log "打包产物: $PACKAGE_FILE"

log "检查远端目录可写（无需 sudo）"
ssh "$SSH_HOST" "bash -s" <<EOF
set -euo pipefail
if [ ! -d "$REMOTE_FRONTEND_DIR" ]; then
  echo "远端目录不存在: $REMOTE_FRONTEND_DIR"
  exit 1
fi
if [ ! -w "$REMOTE_FRONTEND_DIR" ]; then
  echo "远端目录不可写: $REMOTE_FRONTEND_DIR"
  echo "请先在 bgg-dev 执行: sudo chown -R ubuntu:ubuntu $REMOTE_FRONTEND_DIR"
  exit 1
fi
EOF

log "上传到 $SSH_HOST:/tmp/"
scp "$PACKAGE_FILE" "$SSH_HOST:/tmp/"
REMOTE_TAR="/tmp/$(basename "$PACKAGE_FILE")"

log "远端解压替换 dist"
ssh "$SSH_HOST" "bash -s" <<EOF
set -euo pipefail
cd $REMOTE_FRONTEND_DIR
if [ -d dist ]; then
  mv dist dist_\$(date +%Y%m%d%H%M%S).bak
fi
mkdir -p dist
tar -xzf $REMOTE_TAR -C dist
rm -f $REMOTE_TAR
echo "前端部署完成: $REMOTE_FRONTEND_DIR/dist"
EOF

log "验证 https://oldbai.top/bgg/ 是否可访问"
curl -sk -o /dev/null -w "GET /bgg/ -> %{http_code}\n" https://oldbai.top/bgg/

log "完成"
