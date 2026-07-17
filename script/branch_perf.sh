#!/bin/sh
# 从 main 切 perf（性能优化）分支并开 PR（封装 branch_from_main.sh）。
# 用法: ./script/branch_perf.sh <slug> [--remote 远端名 ...]
# 示例: ./script/branch_perf.sh order-export
set -e
SCRIPT_DIR=$(CDPATH= cd -- "$(dirname "$0")" && pwd)
if [ "$#" -lt 1 ]; then
	echo "用法: $0 <slug> [其它参数见 branch_from_main.sh --help]" >&2
	exit 1
fi
SLUG=$1
shift
exec bash "$SCRIPT_DIR/branch_from_main.sh" --branch-type perf --slug "$SLUG" "$@"
