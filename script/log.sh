#!/bin/bash

ERROR_COLOR="\033[31m"  # 错误消息
SUCCESS_COLOR="\033[32m"  # 成功消息
WARNING_COLOR="\033[33m"  # 告警消息
LOG_COLOR="\033[36m"  # 日志消息
PLAIN_TEXT_COLOR='\033[0m'

info() {
  echo -e "${SUCCESS_COLOR}[INFO]$1${PLAIN_TEXT_COLOR}"
}
error() {
  echo -e "${ERROR_COLOR}[ERROR]$1${PLAIN_TEXT_COLOR}"
}
warning() {
  echo -e "${WARNING_COLOR}[WARN]$1${PLAIN_TEXT_COLOR}"
}
log() {
  echo -e "${LOG_COLOR}[LOG]$1${PLAIN_TEXT_COLOR}"
}