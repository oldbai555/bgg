#!/bin/bash

# go-zero RPC 服务代码生成脚本
# 从 .proto 文件生成 services/<name>/ 下的 pb/client/server/logic 骨架
# 支持在任何目录下运行，自动定位项目目录

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
ADMIN_SERVER_DIR="${PROJECT_ROOT}"
TEMPLATE_DIR="${ADMIN_SERVER_DIR}/.template"
SERVICES_DIR="${ADMIN_SERVER_DIR}/services"

# 解析 goctl 路径（优先环境变量 GOCTL_BIN，其次 PATH，再尝试 GOPATH/bin/goctl）
GOCTL_BIN="${GOCTL_BIN:-}"
if [ -z "$GOCTL_BIN" ]; then
    GOCTL_BIN="$(command -v goctl 2>/dev/null || true)"
fi
if [ -z "$GOCTL_BIN" ]; then
    GOPATH_BIN="$(go env GOPATH 2>/dev/null)/bin/goctl"
    if [ -x "$GOPATH_BIN" ]; then
        GOCTL_BIN="$GOPATH_BIN"
    fi
fi
if [ -z "$GOCTL_BIN" ] || [ ! -x "$GOCTL_BIN" ]; then
    echo -e "${RED}错误: goctl 未安装或不可执行${NC}"
    echo "请运行: go install github.com/zeromicro/go-zero/tools/goctl@latest"
    echo "或设置环境变量 GOCTL_BIN 指向 goctl 可执行文件"
    exit 1
fi

# 检查 protoc / 插件
if ! command -v protoc >/dev/null 2>&1; then
    echo -e "${RED}错误: protoc 未安装${NC}"
    echo "macOS: brew install protobuf；Linux: apt-get install -y protobuf-compiler"
    exit 1
fi
if [ -z "$(command -v protoc-gen-go 2>/dev/null)" ]; then
    echo -e "${RED}错误: protoc-gen-go 未安装${NC}"
    echo "请运行: go install google.golang.org/protobuf/cmd/protoc-gen-go@latest"
    exit 1
fi
if [ -z "$(command -v protoc-gen-go-grpc 2>/dev/null)" ]; then
    echo -e "${RED}错误: protoc-gen-go-grpc 未安装${NC}"
    echo "请运行: go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest"
    exit 1
fi

# 检查模板目录
if [ ! -d "$TEMPLATE_DIR/rpc" ]; then
    echo -e "${RED}错误: RPC 模板目录不存在: ${TEMPLATE_DIR}/rpc${NC}"
    exit 1
fi

usage() {
    echo -e "${GREEN}go-zero RPC 服务代码生成工具${NC}"
    echo ""
    echo "用法:"
    echo "  $0 <service_name> [options]"
    echo ""
    echo "参数:"
    echo "  service_name      服务名（如 iam / content / chat / task / sdk）"
    echo ""
    echo "选项:"
    echo "  -f, --file FILE   proto 文件路径（默认: services/<service_name>/rpc/<service_name>.proto）"
    echo "  -h, --help        显示此帮助信息"
    echo ""
    echo "示例:"
    echo "  $0 iam"
    echo "  $0 task -f services/task/rpc/task.proto"
    echo ""
}

SERVICE_NAME=""
PROTO_FILE=""

while [[ $# -gt 0 ]]; do
    case $1 in
        -f|--file)
            PROTO_FILE="$2"
            shift 2
            ;;
        -h|--help)
            usage
            exit 0
            ;;
        *)
            if [ -z "$SERVICE_NAME" ]; then
                SERVICE_NAME="$1"
            else
                echo -e "${RED}错误: 未知参数: $1${NC}"
                usage
                exit 1
            fi
            shift
            ;;
    esac
done

if [ -z "$SERVICE_NAME" ]; then
    echo -e "${RED}错误: 请指定服务名${NC}"
    usage
    exit 1
fi

SERVICE_DIR="${SERVICES_DIR}/${SERVICE_NAME}"
if [ -z "$PROTO_FILE" ]; then
    PROTO_FILE="${SERVICE_DIR}/rpc/${SERVICE_NAME}.proto"
fi

if [ ! -f "$PROTO_FILE" ]; then
    echo -e "${RED}错误: proto 文件不存在: ${PROTO_FILE}${NC}"
    echo "请先按 16-rpc-conventions.md 第 5 节的格式手写 .proto 文件"
    exit 1
fi

# protoc 默认 proto_path 是当前目录（下面会 cd 进 SERVICE_DIR），且要求输入文件路径是
# proto_path 的字面前缀——传绝对路径会被 protoc 判定为"不在任何 proto_path 下"而报错，
# 所以这里统一转换成相对 SERVICE_DIR 的相对路径。
PROTO_FILE_ABS="$(cd "$(dirname "$PROTO_FILE")" && pwd)/$(basename "$PROTO_FILE")"
case "$PROTO_FILE_ABS" in
    "$SERVICE_DIR"/*)
        PROTO_FILE_REL="${PROTO_FILE_ABS#"$SERVICE_DIR"/}"
        ;;
    *)
        echo -e "${RED}错误: proto 文件必须位于 ${SERVICE_DIR}/ 下: ${PROTO_FILE_ABS}${NC}"
        exit 1
        ;;
esac

echo -e "${GREEN}=== go-zero RPC 代码生成 ===${NC}"
echo "项目根目录:  $PROJECT_ROOT"
echo "服务名称:    $SERVICE_NAME"
echo "proto 文件:  $PROTO_FILE"
echo "输出目录:    $SERVICE_DIR"
echo "模板目录:    ${TEMPLATE_DIR}/rpc"
echo "goctl 路径:  $GOCTL_BIN"
echo ""

read -p "确认生成 RPC 代码? (y/N): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo -e "${YELLOW}已取消${NC}"
    exit 0
fi

echo -e "${GREEN}正在生成 RPC 代码...${NC}"
cd "$SERVICE_DIR"

"$GOCTL_BIN" rpc protoc "$PROTO_FILE_REL" \
    --go_out=. --go-grpc_out=. --zrpc_out=. \
    --home "$TEMPLATE_DIR"
GEN_EXIT=$?

if [ $GEN_EXIT -eq 0 ]; then
    echo -e "${GREEN}✓ RPC 代码生成成功!${NC}"
    echo -e "${GREEN}输出目录: ${SERVICE_DIR}${NC}"
else
    echo -e "${RED}✗ RPC 代码生成失败${NC}"
    exit 1
fi
