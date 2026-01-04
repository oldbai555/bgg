#!/bin/bash

# admin.sh bash 自动补全脚本
# 使用方法：
#   source script/admin-completion.bash
# 或者添加到 ~/.bashrc：
#   source /path/to/admin-completion.bash

_admin_completion() {
  local cur prev words cword
  COMPREPLY=()
  cur="${COMP_WORDS[COMP_CWORD]}"
  prev="${COMP_WORDS[COMP_CWORD-1]}"
  words=("${COMP_WORDS[@]}")

  # 检测是否是 sh script/admin.sh 或 bash script/admin.sh 的调用方式
  # 如果是，则第一个参数是脚本路径，第二个参数才是命令
  local cmd_index=1
  if [[ "${COMP_WORDS[0]}" =~ ^(sh|bash|zsh)$ ]] && [[ "${prev}" =~ admin\.sh$ ]]; then
    # sh script/admin.sh <command> 的情况
    cmd_index=2
  elif [[ "${COMP_WORDS[0]}" =~ ^(sh|bash|zsh)$ ]] && [[ "${cur}" != "" ]] && [ ${COMP_CWORD} -eq 2 ]; then
    # sh script/admin.sh <正在输入命令> 的情况
    cmd_index=2
  fi

  # 第一个参数：主命令
  if [ ${COMP_CWORD} -eq ${cmd_index} ]; then
    COMPREPLY=($(compgen -W "dev build package supervisor" -- "${cur}"))
    return 0
  fi

  # 第二个参数：子命令
  local main_cmd="${COMP_WORDS[${cmd_index}]}"
  case "${main_cmd}" in
    dev)
      if [ ${COMP_CWORD} -eq $((cmd_index + 1)) ]; then
        COMPREPLY=($(compgen -W "start stop status logs" -- "${cur}"))
      fi
      return 0
      ;;
    build)
      if [ ${COMP_CWORD} -eq $((cmd_index + 1)) ]; then
        COMPREPLY=($(compgen -W "server frontend" -- "${cur}"))
      fi
      return 0
      ;;
    package)
      if [ ${COMP_CWORD} -eq $((cmd_index + 1)) ]; then
        COMPREPLY=($(compgen -W "server frontend" -- "${cur}"))
      fi
      return 0
      ;;
    supervisor)
      if [ ${COMP_CWORD} -eq $((cmd_index + 1)) ]; then
        COMPREPLY=($(compgen -W "gen-conf install deploy status start stop restart logs" -- "${cur}"))
      fi
      return 0
      ;;
  esac

  return 0
}

# 注册补全函数（支持多种调用方式）
complete -F _admin_completion admin.sh
complete -F _admin_completion ./admin.sh
complete -F _admin_completion script/admin.sh
complete -F _admin_completion ./script/admin.sh

# 对于 sh/bash 命令，需要特殊处理
# 当输入 sh script/admin.sh 时，第二个参数是脚本路径，第三个参数才是命令
_admin_sh_completion() {
  local cur="${COMP_WORDS[COMP_CWORD]}"
  local prev="${COMP_WORDS[COMP_CWORD-1]}"
  
  # 如果上一个词是 admin.sh 相关的脚本路径，则补全命令
  if [[ "${prev}" =~ admin\.sh$ ]]; then
    _admin_completion
    return 0
  fi
  
  # 否则使用默认的文件补全
  COMPREPLY=($(compgen -f -- "${cur}"))
  return 0
}

complete -F _admin_sh_completion sh
complete -F _admin_sh_completion bash

