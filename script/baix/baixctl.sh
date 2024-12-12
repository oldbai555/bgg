#!/bin/bash

# 定义函数部分
function safe_process_env() {
  # 在处理.env文件之前，检查文件是否存在和可读
  if [ ! -f .env ]; then
    echo "Error: .env file does not exist." >&2
    exit 1
  fi
  if ! grep -q -E '^[A-Za-z0-9_]+=[A-Za-z0-9_]+$' .env; then
    echo "Error: .env file format is incorrect." >&2
    exit 1
  fi
  # 进行环境变量的导出
  export $(xargs <.env)
}

function sed_i_backup() {
  # 对.env文件进行备份并删除\r字符
  local file=".env"
  if ! cp "$file" "$file.bak"; then
    echo "Error: Failed to backup $file." >&2
    exit 1
  fi
  if ! sed -i 's/\r//' "$file"; then
    echo "Error: Failed to remove \r characters from $file." >&2
    exit 1
  fi
}

function safe_exec_baixctl() {
  # 安全执行baixctl命令，对输入参数进行验证
  local cmd="$1"
  shift
  local params=("$@")
  # 这里可以添加更多的参数验证逻辑
  if [ -z "$cmd" ]; then
    echo "Error: Missing baixctl command." >&2
    exit 1
  fi
  ./baixctl "$cmd" "${params[@]}"
}

function usage() {
  echo "Usage: sh baixctl.sh [OPTION]"
  echo "
        genAddCurdRpc 添加pb的Curd rpc方法
        genAddProto   初始化pb文件
        genAddRpc     添加pb的rpc方法
        genClient     生成客户端代码
        genServer     生成服务端代码
        gensinglepb   执行单体应用proto代码生成器"
  exit 1
}

# 命令执行函数
function genAddCurdRpc() {
  safe_exec_baixctl genAddCurdRpc -p "$1" -m "$2" -s "$3"
}

function genAddProto() {
  safe_exec_baixctl genAddProto -p "$1"
}

function genAddRpc() {
  safe_exec_baixctl genAddRpc -p "$1" -r "$2"
}

function genClient() {
  safe_exec_baixctl genClient -p "$1"
}

function genServer() {
  safe_exec_baixctl genServer -p "$1"
}

function gensinglepb() {
  safe_exec_baixctl gensinglepb -p "$1"
}

# 主逻辑部分
sed_i_backup
safe_process_env

case "$1" in
  "gacr" | "genAddCurdRpc")
    genAddCurdRpc "$2" "$3" "$4"
    ;;
  "gap" | "genAddProto")
    genAddProto "$2"
    ;;
  "gar" | "genAddRpc")
    genAddRpc "$2" "$3"
    ;;
  "gc" | "genClient")
    genClient "$2"
    ;;
  "gs" | "genServer")
    genServer "$2"
    ;;
  "gall")
    genClient "$2"
    genServer "$2"
    ;;
  "gb" | "gensinglepb")
    gensinglepb "$2"
    ;;
  *)
    usage
    ;;
esac
