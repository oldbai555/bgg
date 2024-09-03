#!/bin/bash

# 设置环境变量
# source .env

sh_dir=$(cd "$(dirname "$0")" && pwd)
echo "脚本目录 shDir:"$sh_dir

# 项目目录
pro_dir=$(dirname ${sh_dir})
echo "项目目录 proDir:"${pro_dir}

# 项目根目录
proRootDir=$(dirname ${pro_dir})
echo "项目根目录 proRootDir:"${proRootDir}

# 定义全局变量
supervisorDir="/home/work/service"
supervisorLogDir="/home/work/supervisor/logs"
packOutputDir="${proRootDir}/package"

# 检查并创建日志目录 打包输出目录
mkdir -p "$supervisorDir"
mkdir -p "$supervisorLogDir"
mkdir -p "$packOutputDir"
chmod +x -R "$supervisorLogDir"

# 函数：生成 Supervisor 配置文件
outputSupervisorConf() {
  local programName=$1
  local outputDir=$2

  cat >"$outputDir/$programName.conf" <<EOF
[program:$programName]
directory=$supervisorDir/$programName
command=$supervisorDir/$programName/$programName
autostart=true
autorestart=true
startsecs=10
startretries=3
user=root
redirect_stderr=true
stdout_logfile=$supervisorLogDir/${programName}_stdout.log
stderr_logfile=$supervisorLogDir/${programName}_stderr.log
stdout_logfile_maxbytes=20MB
stdout_logfile_backups=20
EOF
}

# 函数：部署 Supervisor 服务
deploySupervisorService() {
  local appName=$1

  cd "$supervisorDir" || exit

  rm -rf "$appName"

  mkdir -p "$appName"

  tar -xvf "$packOutputDir/$appName.tar" -C "$supervisorDir/$appName"

  chmod +x -R "$supervisorDir/$appName"

  cp "$supervisorDir/$appName/$appName.conf" /etc/supervisor/conf.d

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
  local appName=$3

  export GOOS=linux
  export GOARCH=amd64
  go build -o "$outputDir/$appName" "$appPath"
  unset GOOS
  unset GOARCH
}

# 函数：打包并部署服务
localPackSrv() {
  # 服务名称
  local appName=$1

  # 打包路径
  local proSrvPath=$2

  # 输出目录
  local outputDir="$packOutputDir/$appName"

  mkdir -p ${outputDir}

  outputSupervisorConf ${appName} ${outputDir}

  goBuild ${outputDir} ${proSrvPath}

  cp "${proSrvPath}/application.yaml" ${outputDir}

  tar -cvf "${appName}.tar" -C "$outputDir" . --remove-files
}

# 使用说明
usage() {
  echo "Usage: sh supervisor.sh [OPTION]"
  echo "  lpg: 打包网关"
  echo "  lpcmd: 打包服务"
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
  localPackSrv "$1" ${proRootDir}"/service/gateway"
  deploySupervisorService "$1"
  ;;
lpcmd)
  shift
  localPackSrv "$1" ${proRootDir}"/singlesrv/cmd"
  deploySupervisorService "$1"
  ;;
*)
  usage
  exit 1
  ;;
esac

exit 0
