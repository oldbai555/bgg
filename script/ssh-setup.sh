#!/bin/bash
# 一键配置 SSH 免密登录 + 写入 ~/.ssh/config 别名
# 用法: ./ssh-setup.sh <user> <ip> [alias] [port]
#   alias 默认取 ip；port 默认 22
set -e

USER_NAME="$1"
HOST_IP="$2"
ALIAS_NAME="${3:-$HOST_IP}"
PORT="${4:-22}"
KEY_PATH="$HOME/.ssh/id_ed25519"
CONFIG_PATH="$HOME/.ssh/config"

if [ -z "$USER_NAME" ] || [ -z "$HOST_IP" ]; then
  echo "用法: $0 <user> <ip> [alias] [port]"
  exit 1
fi

# 1. 密钥不存在则生成
if [ ! -f "$KEY_PATH" ]; then
  echo "[1/3] 生成 SSH 密钥: $KEY_PATH"
  ssh-keygen -t ed25519 -N "" -f "$KEY_PATH"
else
  echo "[1/3] 已存在密钥，跳过生成: $KEY_PATH"
fi

# 2. 拷贝公钥到目标机器（需要输入一次目标机密码）
echo "[2/3] 拷贝公钥到 ${USER_NAME}@${HOST_IP}:${PORT}（需要输入一次密码）"
ssh-copy-id -p "$PORT" -i "${KEY_PATH}.pub" "${USER_NAME}@${HOST_IP}"

# 3. 写入 ~/.ssh/config 别名（若已存在同名 Host 则跳过，避免重复）
mkdir -p "$HOME/.ssh"
touch "$CONFIG_PATH"
if grep -qE "^Host[[:space:]]+${ALIAS_NAME}$" "$CONFIG_PATH" 2>/dev/null; then
  echo "[3/3] ~/.ssh/config 中已存在 Host ${ALIAS_NAME}，跳过写入"
else
  echo "[3/3] 写入 ~/.ssh/config 别名: ${ALIAS_NAME}"
  {
    echo ""
    echo "Host ${ALIAS_NAME}"
    echo "    HostName ${HOST_IP}"
    echo "    User ${USER_NAME}"
    echo "    Port ${PORT}"
    echo "    IdentityFile ${KEY_PATH}"
  } >> "$CONFIG_PATH"
fi

echo ""
echo "完成。验证: ssh ${ALIAS_NAME} 'echo ok'"
