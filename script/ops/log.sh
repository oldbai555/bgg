#!/bin/bash

RD="\033[31m" # 错误消息
GR="\033[32m" # 成功消息
YL="\033[33m" # 告警消息
BL="\033[36m" # 日志消息
PL='\033[0m'

INFO() {
  echo -e "${GR}[INFO]$1${PL}"
}
ERROR() {
  echo -e "${RD}[ERROR]$1${PL}"
}
WANF() {
  echo -e "${YL}[WANF]$1${PL}"
}
LOG() {
  echo -e "${BL}[LOG]$1${PL}"
}