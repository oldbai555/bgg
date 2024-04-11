#!/bin/bash

# 设置环境变量
# source .env

# 取代码目录
sh_dir=$(cd "$(dirname "$0")" && pwd)
echo "sh dir:"$sh_dir
pro_dir=$(dirname ${sh_dir})
projectRootDir=$(dirname ${pro_dir})
echo "projectRootDir:"${projectRootDir}

# 定义全局变量
supervisorDir="/home/work/service"
supervisorLogDir="/home/work/supervisor/logs"
packOutputDir="${projectRootDir}/package"

# 检查并创建日志目录 打包输出目录
mkdir -p "$supervisorDir"
mkdir -p "$supervisorLogDir"
mkdir -p "$packOutputDir"
chmod +x -R "$supervisorLogDir"

# 函数：生成 Supervisor 配置文件
outputSupervisorConf() {
  local programName=$1
  local outputDir=$2
  local appName=$3

  cat >"$outputDir/$programName.ini" <<EOF
[program:$programName]
directory=$supervisorDir/$programName
command=$supervisorDir/$programName/$appName
autostart=true
autorestart=true
startsecs=10
startretries=3
user=root
redirect_stderr=true
stdout_logfile=$supervisorLogDir/${programName}_stdout.log
stdout_logfile_maxbytes=20MB
stdout_logfile_backups=20
EOF
}

# 函数：打包并部署网关服务
localPackGate() {
  local programName=$1
  local outputDir="$packOutputDir/$programName"
  mkdir -p "$outputDir"

  outputSupervisorConf "$programName" "$outputDir" "$programName"
  goBuild "$outputDir" "$projectRootDir/service/$programName"
  cp "$projectRootDir/service/$programName/application.yaml" "$outputDir"
  tar -cvf "$outputDir.tar" -C "$outputDir" .

  deploySupervisorService "$programName"
}

# 函数：打包并部署服务器服务
localPackServer() {
  local appName=$1
  local outputDir="$packOutputDir/$appName"
  mkdir -p "$outputDir"

  outputSupervisorConf "$appName" "$outputDir" "cmd"
  goBuild "$outputDir/cmd" "$projectRootDir/service/$appName/server/cmd"
  cp "$projectRootDir/service/$appName/server/cmd/application.yaml" "$outputDir"
  tar -cvf "$outputDir.tar" -C "$outputDir" .

  deploySupervisorService "$appName"
}

# 函数：部署 Supervisor 服务
deploySupervisorService() {
  local appName=$1

  cd "$supervisorDir" || exit
  rm -rf "$appName"
  mkdir -p "$appName"
  tar -xvf "$packOutputDir/$appName.tar" -C "$supervisorDir/$appName"
  chmod +x -R "$supervisorDir/$appName"
  cp "$supervisorDir/$appName/$appName.ini" /etc/supervisor/conf.d

  # 更新并重启 Supervisor 服务
  supervisorctl update
  supervisorctl restart "$appName"
  # 检查服务状态
  supervisorctl status "$appName"
}

# 函数：构建 Go 应用
goBuild() {
  local outputDir=$1
  local appPath=$2

  export GOOS=linux
  export GOARCH=amd64
  go build -o "$outputDir/app" "$appPath"
  unset GOOS
  unset GOARCH
}

# 使用说明
usage() {
  echo "Usage: sh supervisor.sh [OPTION]"
  echo "  lpg | localPackGate: 打包并部署网关服务"
  echo "  lps | localPackServer: 打包并部署服务器服务"
  # ... 可以添加更多选项
}

# 主程序逻辑
if [ "$#" -eq 0 ]; then
  usage
  exit 1
fi

# 解析命令行参数
case "$1" in
lpg)
  shift
  localPackGate "$1"
  ;;
lps)
  shift
  localPackServer "$1"
  ;;
*)
  usage
  exit 1
  ;;
esac

exit 0
