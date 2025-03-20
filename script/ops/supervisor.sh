#!/bin/bash

# 设置环境变量
# source .env

sh_dir=$(cd "$(dirname "$0")" && pwd)
echo "脚本目录 shDir:"$sh_dir

# 项目目录
pro_dir=$(dirname ${sh_dir})
echo "脚本根目录 proDir:"${pro_dir}

# 项目根目录
proRootDir=$(dirname ${pro_dir})
echo "项目根目录 proRootDir:"${proRootDir}

# 定义全局变量
supervisorDir="/home/work/service"
supervisorLogDir="/home/work/supervisor/logs"
packOutputDir="${proRootDir}/package"

MkdirIfFileUnExist() {
  stat $1 >/dev/null 2>&1
  if [ $? != 0 ]; then
    mkdir -p "$1"
  fi
}

# 检查并创建日志目录 打包输出目录
MkdirIfFileUnExist "$supervisorDir"
MkdirIfFileUnExist "$supervisorLogDir"
MkdirIfFileUnExist "$packOutputDir"
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

# 函数：根据文件名前缀获取最新的文件
getLatestFileByPrefix() {
  local dir=$1
  local prefix=$2

  # 使用 find 命令查找所有以 prefix 开头的文件，并按修改时间排序，取最新的一个
  local latestFile=$(find "$dir" -type f -name "${prefix}*" -printf '%T+ %p\n' | sort -r | head -n 1 | cut -d' ' -f2-)
  echo "$latestFile"
}

# 函数：部署 Supervisor 服务
deploySupervisorService() {
  local packName
  local appName
  packName=$1
  appName=$(echo "$packName" | cut -d'_' -f1)
  echo "====== 包体 $packName ======"
  echo "====== 开始部署 $appName ======"

  cd "$supervisorDir" || exit

  MkdirIfFileUnExist "$appName"

  tar -xvf "$sh_dir/$packName" -C "$supervisorDir/$appName"

  chmod +x -R "$supervisorDir/$appName"

  cp "$supervisorDir/$appName/$appName.conf" /etc/supervisor/conf.d

  # 更新并重启 Supervisor 服务
  supervisorctl update

  supervisorctl restart "$appName"

  # 检查服务状态
  supervisorctl status "$appName"
  echo "====== 完成部署 $appName ======"
}

# 获取改动时间部署最新的包体
deployV2() {
  appName=$1
  packName=$(getLatestFileByPrefix $sh_dir $appName)
  echo "====== 包体 $packName ======"
  echo "====== 开始部署 $appName ======"

  cd "$supervisorDir" || exit

  MkdirIfFileUnExist "$appName"

  tar -xvf "$packName" -C "$supervisorDir/$appName"

  chmod +x -R "$supervisorDir/$appName"

  cp "$supervisorDir/$appName/$appName.conf" /etc/supervisor/conf.d

  # 更新并重启 Supervisor 服务
  supervisorctl update

  supervisorctl restart "$appName"

  # 检查服务状态
  supervisorctl status "$appName"
  echo "====== 完成部署 $appName ======"
}

# 函数：构建 Go 应用
goBuild() {
  local outputDir=$1
  local appPath=$2
  local appName=$3
  echo "====== 开始编译 $appName ======"
  echo "====== appPath: $appPath ======"
  echo "====== outputDir: $outputDir ======"
  go env -w CGO_ENABLED=0
  export GOOS=linux
  export GOARCH=amd64
  go build -o "$outputDir/$appName" "$appPath"
  unset GOOS
  unset GOARCH
  go env -w CGO_ENABLED=1
  echo "====== 完成编译 $appName ======"
}

# 函数：打包并部署服务
localPackSrv() {
  # 服务名称
  local appName=$1
  echo "====== 开始打包 $appName ======"

  hash=$(git rev-parse --short HEAD)
  echo "版本号："${hash}

  # 打包路径
  local proSrvPath=$2

  # 输出目录
  local outputDir=${packOutputDir}"/"${appName}

  MkdirIfFileUnExist ${outputDir}

  outputSupervisorConf ${appName} ${outputDir}

  goBuild ${outputDir} ${proSrvPath} ${appName}

  #  cp "${proSrvPath}/application.yaml" ${outputDir}

  local srvOutPut=${packOutputDir}"/${appName}_${hash}.tar"
  tar -cvf ${srvOutPut} -C ${outputDir} . --remove-files
  echo "====== 输出文件 ${srvOutPut}"
  echo "====== 完成打包 $appName ======"

  start_upload ${srvOutPut}
}

start_upload() {
  ./uploader.exe $1
}

deployWeb() {
  file=$1
  mv ${file} /home/work/web && cd /home/work/web && rm -rf dist && tar -zxvf ${file}
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
  #  deploySupervisorService "$1"
  ;;
lpcmd)
  shift
  localPackSrv "$1" ${proRootDir}"/service/"$1"server/cmd"
  #  deploySupervisorService "$1"
  ;;
deploy)
  shift
  deploySupervisorService "$1"
  ;;
dall)
  shift
  deployV2 gateway
  deployV2 lbsingle
  deployV2 lboss
  deployV2 lbwxmp
  ;;
web)
  shift
  deployWeb "$1"
  ;;
upload)
  shift
  start_upload "$1"
  ;;
*)
  usage
  exit 1
  ;;
esac

exit 0
