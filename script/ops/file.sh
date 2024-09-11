#!/bin/bash

source ./log.sh

# 创建目录 如果目录不存在
CreateDirIfNotExists() {
  local dir_path="$1"

  # 使用stat命令检查目录是否存在
  if ! stat -t --dereference --format='%F' "$dir_path" &>/dev/null; then
    # 检查失败，可能是权限问题、路径不存在或其他错误
    local error_message=$(stat -t --dereference --format='%s\n%F' "$dir_path" 2>&1)
    echo "Error: Failed to check directory existence for '$dir_path': $error_message" >&2
    return 1
  fi

  # 如果目录不存在，则创建
  if [ "$(stat -t --dereference --format='%F' "$dir_path")" != "directory" ]; then
    mkdir -p "$dir_path"
  fi
}
