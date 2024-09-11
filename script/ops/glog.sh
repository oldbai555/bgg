#!/bin/bash

# 检查是否提供了至少一个参数
if [[ "$#" -lt 1 ]]; then
  echo "用法: glog.sh <选项> <关键字>"
  echo "选项:"
  echo "-t 默认当前年月日时"
  echo "-t [前x小时]"
  echo "-t [指定日小时 例如 0114]"
  echo "-t [指定月日小时 例如 120114]"
  exit 1
fi

# 初始化变量
D=""
LOGDIR="/home/work/log/"

# 解析选项和关键字
while [[ $# -gt 0 ]]; do
  case "$1" in
  -t)
    # 检查是否有提供时间值
    if [[ $# -lt 2 ]]; then
      echo "错误: -t 选项需要一个参数"
      exit 1
    fi
    D="$2"
    shift # 跳过下一个参数
    ;;
  *)
    # 所有非选项参数都被视为关键字
    args_filter+=("$1")
    ;;
  esac
  shift # 跳过当前参数
done

# 如果没有提供时间选项，则默认为当前时间
if [[ "$D" == "" ]]; then
  D="0"
fi

# 根据提供的时间选项生成日志文件名
LOGFILE=""
case $D in
0)
  LOGFILE=$(date "+%Y%m%d%H").log
  ;;
[1-9]) # 注意这里改为 [1-9]，因为我们不期望单个数字代表小时数
  LOGFILE=$(date -d "-${D} hour" +%Y%m%d%H).log
  ;;
[0-9][0-9]) # 14
  LOGFILE=$(date "+%Y%m%d")${D:0:2}${D:2:2}.log
  ;;
[0-9][0-9][0-9][0-9]) #  0114
  LOGFILE=$(date "+%Y%m")${D:0:2}${D:2:2}${D:4:2}.log
  ;;
[0-9][0-9][0-9][0-9][0-9][0-9]) # 120114
  LOGFILE=$(date "+%Y")${D:0:2}${D:2:2}${D:4:2}.log
  ;;
*)
  echo "错误: 不支持的时间格式。请使用如下格式之一: [0], [1-9], [0-9][0-9], [0-9][0-9][0-9][0-9], [0-9][0-9][0-9][0-9][0-9][0-9]"
  exit 1
  ;;
esac

# 构造日志文件路径
LOG="$LOGDIR"*"$LOGFILE"
LOG=${LOG//\\/\//}

# 检查日志文件是否存在
NLOG=$(find "$LOGDIR" -name "*$LOGFILE" | wc -l)
if [ $NLOG -eq 0 ]; then
  echo "文件 $LOG 不存在，停止执行"
  exit 1
fi

# 输出执行的命令
echo "执行命令: cat $LOG | grep -a ${args_filter[*]}"

# 执行搜索并排序结果
cat $LOG | grep -a "${args_filter[@]}" | sort -k3
