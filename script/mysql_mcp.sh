#!/usr/bin/env bash
# MySQL MCP 启动包装：加载 ~/.config/bgg/mysql-mcp.env，再启动 @benborla29/mcp-server-mysql。
#
# Cursor Shared MCP 进程常不带 HOME，内联 `source $HOME/.config/...` 会静默失败，
# 导致回退到 127.0.0.1:3306 → 进程启动后立刻退出 → MCP error -32000 Connection closed。
set -euo pipefail

if [[ -z "${HOME:-}" ]]; then
  if [[ -n "${USER:-}" ]]; then
    HOME="$(eval echo "~${USER}")"
  else
    HOME="$(cd ~ && pwd)"
  fi
  export HOME
fi

ENV_FILE="${MYSQL_MCP_ENV:-${HOME}/.config/bgg/mysql-mcp.env}"
if [[ -f "$ENV_FILE" ]]; then
  set -a
  # shellcheck disable=SC1090
  source "$ENV_FILE"
  set +a
fi

export ALLOW_INSERT_OPERATION="${ALLOW_INSERT_OPERATION:-false}"
export ALLOW_UPDATE_OPERATION="${ALLOW_UPDATE_OPERATION:-false}"
export ALLOW_DELETE_OPERATION="${ALLOW_DELETE_OPERATION:-false}"

exec npx -y @benborla29/mcp-server-mysql@2.0.9
