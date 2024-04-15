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
  echo "gc | genclient 根据 proto 生成客户端代码"
  echo "gs | genserver 根据 proto 生成服务端代码"
  echo "gsc | gensc 根据 proto 生成客户端 服务端代码"
  echo "gg | gengingateway 根据 proto 生成网关代码"
  echo "gt | gents 根据 proto 生成 ts 代码"
  echo "ap | addproto 新增 proto 文件"
  echo "ar | addrpc 新增 rpc 方法"
  echo "ac | addCurdRpc 新增 curd rpc 方法以及 message. 示例 sh bin.sh addCurdRpc lbbill.proto Bill,BillCategory"
  echo "asc | addCurdSysRpc 新增系统 curd rpc 方法以及 message. 示例 sh bin.sh addCurdRpc lbbill.proto Bill,BillCategory"
  echo "gts 生成 ts 代码"
  exit 1
}

# 命令执行函数
function gen_client() {
  safe_exec_baixctl genclient -p "$1"
}

function gen_server() {
  safe_exec_baixctl genserver -p "$1"
}

function gen_client_and_server() {
  safe_exec_baixctl genclient -p "$1"
  safe_exec_baixctl genserver -p "$1"
}

function gen_gateway() {
  safe_exec_baixctl gengingateway -p "$1"
}

function gen_ts() {
  safe_exec_baixctl gen_ts_vue -p "$1" -o "$2"
}

function add_rpc() {
  safe_exec_baixctl genAddRpc -p "$1" -r "$2"
}

function add_curd_rpc() {
  safe_exec_baixctl genAddCurdRpc -p "$1" -m "$2"
}

function add_curd_sys_rpc() {
  safe_exec_baixctl genAddCurdRpc -p "$1" -m "$2" -s true
}

function add_proto() {
  safe_exec_baixctl genAddProto -p "$1"
}

function gen_ts_admin() {
  safe_exec_baixctl gen_ts_vue  -p "$1" -o "/e/bgg/github.com/oldbai555/bgg/webv2/admin"
}

# 主逻辑部分
sed_i_backup
safe_process_env

case "$1" in
  "gc" | "genclient")
    gen_client "$2"
    ;;
  "gs" | "genserver")
    gen_server "$2"
    ;;
  "gsc" | "gensc")
    gen_client_and_server "$2"
    ;;
  "gg" | "gengingateway")
    gen_gateway "$2"
    ;;
  "gt" | "gents")
    gen_ts "$2" "$3"
    ;;
  "ar" | "addrpc")
    add_rpc "$2" "$3"
    ;;
  "ap" | "addproto")
    add_proto "$2"
    ;;
  "ac" | "addCurdRpc")
    add_curd_rpc "$2" "$3"
    ;;
  "asc" | "addCurdSysRpc")
    add_curd_sys_rpc "$2" "$3"
    ;;
  "gts")
    gen_ts_admin "$2"
    ;;
  *)
    usage
    ;;
esac
