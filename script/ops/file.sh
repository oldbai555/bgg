#!/bin/bash

source ./log.sh

# 创建目录 如果 目录不存在
MkdirIfFileUnExist() {
  stat $1
  if [ $? != 0 ]; then
    mkdir -p "$1"
  fi
}
