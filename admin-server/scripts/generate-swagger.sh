#!/bin/bash
# goctl-swagger API 文档生成脚本
# 从 admin.api 生成 OpenAPI (Swagger) JSON 文档

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ADMIN_SERVER_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
API_DIR="${ADMIN_SERVER_DIR}/api"
OPENAPI_DIR="${ADMIN_SERVER_DIR}/docs/openapi"
OUTPUT_FILE="admin-api.json"

# GOCTL_BIN 解析（与 generate-model.sh / generate-ts.sh 一致）
GOCTL_BIN="${GOCTL_BIN:-}"
[ -z "$GOCTL_BIN" ] && GOCTL_BIN="$(command -v goctl 2>/dev/null || true)"
if [ -z "$GOCTL_BIN" ]; then
    GOPATH_BIN="$(go env GOPATH 2>/dev/null)/bin/goctl"
    [ -x "$GOPATH_BIN" ] && GOCTL_BIN="$GOPATH_BIN"
fi
if [ -z "$GOCTL_BIN" ] || [ ! -x "$GOCTL_BIN" ]; then
    echo -e "${RED}错误: goctl 未安装或不可执行${NC}"
    echo "请运行: go install github.com/zeromicro/go-zero/tools/goctl@latest"
    exit 1
fi

# GOCTL_SWAGGER_BIN 解析（同一套模式）
GOCTL_SWAGGER_BIN="${GOCTL_SWAGGER_BIN:-}"
[ -z "$GOCTL_SWAGGER_BIN" ] && GOCTL_SWAGGER_BIN="$(command -v goctl-swagger 2>/dev/null || true)"
if [ -z "$GOCTL_SWAGGER_BIN" ]; then
    GOPATH_SWAGGER_BIN="$(go env GOPATH 2>/dev/null)/bin/goctl-swagger"
    [ -x "$GOPATH_SWAGGER_BIN" ] && GOCTL_SWAGGER_BIN="$GOPATH_SWAGGER_BIN"
fi
if [ -z "$GOCTL_SWAGGER_BIN" ] || [ ! -x "$GOCTL_SWAGGER_BIN" ]; then
    echo -e "${RED}错误: goctl-swagger 未安装或不可执行${NC}"
    echo "请运行: go install github.com/zeromicro/goctl-swagger@latest"
    exit 1
fi

usage() {
    echo -e "${GREEN}go-zero Swagger/OpenAPI 文档生成工具${NC}"
    echo ""
    echo "用法: $0 [api_file]"
    echo "  api_file  API 文件路径（可选，默认 api/admin.api）"
    echo ""
    echo "输出: docs/openapi/admin-api.json（固定路径）"
}

API_FILE="${1:-admin.api}"
[ "$API_FILE" = "-h" ] || [ "$API_FILE" = "--help" ] && { usage; exit 0; }

if [[ "$API_FILE" != /* ]]; then
    [ -f "${API_DIR}/${API_FILE}" ] && API_FILE="${API_DIR}/${API_FILE}"
fi
[ -f "$API_FILE" ] || { echo -e "${RED}错误: 找不到文件: ${API_FILE}${NC}"; exit 1; }

mkdir -p "$OPENAPI_DIR"

echo -e "${GREEN}=== goctl-swagger 文档生成 ===${NC}"
echo "API 文件:   $API_FILE"
echo "输出文件:   ${OPENAPI_DIR}/${OUTPUT_FILE}"
echo ""

read -p "确认生成 Swagger 文档? (y/N): " -n 1 -r
echo
[[ ! $REPLY =~ ^[Yy]$ ]] && { echo -e "${YELLOW}已取消${NC}"; exit 0; }

cd "$OPENAPI_DIR"
"$GOCTL_BIN" api plugin \
    -plugin "${GOCTL_SWAGGER_BIN}=\"swagger -filename ${OUTPUT_FILE}\"" \
    -api "$API_FILE" -dir .

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Swagger 文档生成成功: ${OPENAPI_DIR}/${OUTPUT_FILE}${NC}"
else
    echo -e "${RED}✗ Swagger 文档生成失败${NC}"
    exit 1
fi
