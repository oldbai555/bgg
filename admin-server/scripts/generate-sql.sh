#!/bin/bash

# SQL 脚本生成工具
# 用于快速生成新功能模块的初始化 SQL 脚本
# 支持在任何目录下运行，自动定位项目目录

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 获取脚本所在目录的绝对路径
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# 项目根目录（scripts的父目录）
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
SQLGEN_DIR="${PROJECT_ROOT}/scripts/sqlgen"
# 输出目录按 -group 的 <domain>/<module> 解析后计算，见下方「域→服务映射」
OUTPUT_DIR=""

# 域→服务映射（15-service-boundaries.md 第 4 节 SSOT）：
# iam/system/monitoring/misc 四个域合并进 iam-rpc；blog/video 合并进 content-rpc；
# chat/task/sdk 各自独立成服务。domain 也接受直接写服务名（iam/content/chat/task/sdk）。
domain_to_service() {
    case "$1" in
        iam|system|monitoring|misc) echo "iam" ;;
        blog|blog_extension|video|content) echo "content" ;;
        chat) echo "chat" ;;
        task) echo "task" ;;
        sdk) echo "sdk" ;;
        *) echo "" ;;
    esac
}

# 域→前端 API wrapper 映射（admin-frontend/docs/02-domain-reorg-and-api-layer.md SSOT）：
# 前端按 8 个域分 wrapper（iam/system/monitoring/misc/content/chat/sdk/task），blog/video
# 合并进同一个 content.ts，其余域名与后端 domain 一致。用于生成的 list_page.vue.tpl
# import '@/api/<前端 Domain>'，不是 import '@/api/generated/admin'。
domain_to_frontend_api() {
    case "$1" in
        blog|blog_extension|video) echo "content" ;;
        *) echo "$1" ;;
    esac
}

# 显示使用说明
usage() {
    echo -e "${GREEN}SQL 脚本生成工具${NC}"
    echo ""
    echo "用法:"
    echo "  $0 -group <domain>/<module> -name <name>"
    echo ""
    echo "参数:"
    echo "  -group <domain>/<module>  功能组名（必需，如 iam/user、blog/article、chat/chat）"
    echo "                            <domain> 决定落进哪个服务的 db/services/<service>/ 目录，"
    echo "                            取值：iam/system/monitoring/misc → iam；blog/video → content；"
    echo "                            chat/task/sdk 各自独立；也可以直接写服务名 iam/content/chat/task/sdk"
    echo "  -name <name>          功能名称（必需，如 用户管理, 文件管理）"
    echo ""
    echo "选项:"
    echo "  -parent-id <id>       父菜单 ID（可选，优先级最高）"
    echo "  -parent-path <path>   前端父目录路径（可选，如 /system，默认 /temp）"
    echo "  -h, --help            显示此帮助信息"
    echo ""
    echo "示例:"
    echo "  $0 -group iam/user -name 用户管理"
    echo "  $0 -group system/file -name 文件管理"
    echo "  $0 -group monitoring/operation_log -name 操作日志 -parent-path /system"
    echo "  $0 -group blog/article -name 文章管理"
    echo ""
    echo "注意:"
    echo "  - 生成的 SQL 文件在 admin-server/db/services/<service>/<module>/ 目录下"
    echo "  - 文件名格式: create_table_<module>.sql、init_<module>.sql"
    echo "  - 主键为自增，不需要手动赋值"
    echo "  - 默认菜单父目录为临时目录 /temp，如需挂到系统管理请使用 -parent-path /system"
    echo "  - 包含菜单、权限、接口及关联关系"
    echo ""
}

# 检查参数
if [ $# -eq 0 ]; then
    usage
    exit 0
fi

GROUP=""
NAME=""
PARENT_ID=""
PARENT_PATH=""

# 解析参数
while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            usage
            exit 0
            ;;
        -group)
            GROUP="$2"
            shift 2
            ;;
        -name)
            NAME="$2"
            shift 2
            ;;
        -parent-id)
            PARENT_ID="$2"
            shift 2
            ;;
        -parent-path)
            PARENT_PATH="$2"
            shift 2
            ;;
        *)
            echo -e "${RED}错误: 未知参数: $1${NC}"
            usage
            exit 1
            ;;
    esac
done

# 检查必需参数
if [ -z "$GROUP" ] || [ -z "$NAME" ]; then
    echo -e "${RED}错误: 必须提供 -group 和 -name 参数${NC}"
    usage
    exit 1
fi

# 解析 -group 的 <domain>/<module> 格式，映射到 db/services/<service>/<module>/
if [[ "$GROUP" != */* ]]; then
    echo -e "${RED}错误: -group 必须是 <domain>/<module> 格式（如 iam/user），收到: ${GROUP}${NC}"
    usage
    exit 1
fi
DOMAIN="${GROUP%%/*}"
MODULE="${GROUP#*/}"
if [[ -z "$DOMAIN" || -z "$MODULE" || "$MODULE" == */* ]]; then
    echo -e "${RED}错误: -group 格式不合法: ${GROUP}${NC}"
    exit 1
fi
SERVICE="$(domain_to_service "$DOMAIN")"
if [ -z "$SERVICE" ]; then
    echo -e "${RED}错误: 未知 domain: ${DOMAIN}（有效值: iam, system, monitoring, misc, blog, video, chat, task, sdk，或直接写服务名 iam/content/chat/task/sdk）${NC}"
    exit 1
fi
FRONTEND_DOMAIN="$(domain_to_frontend_api "$DOMAIN")"
OUTPUT_DIR="${PROJECT_ROOT}/db/services/${SERVICE}/${MODULE}"
GROUP="$MODULE"

# 检查 sqlgen 目录是否存在
if [ ! -d "$SQLGEN_DIR" ]; then
    echo -e "${RED}错误: sqlgen 目录不存在: ${SQLGEN_DIR}${NC}"
    exit 1
fi

# 检查模板文件是否存在
TEMPLATE_FILE="${SQLGEN_DIR}/templates/init_module.sql.tpl"
if [ ! -f "$TEMPLATE_FILE" ]; then
    echo -e "${RED}错误: 模板文件不存在: ${TEMPLATE_FILE}${NC}"
    exit 1
fi

# 显示配置信息
echo -e "${GREEN}=== SQL 脚本生成工具 ===${NC}"
echo "项目根目录:  $PROJECT_ROOT"
echo "功能组名:    $GROUP"
echo "功能名称:    $NAME"
echo "输出目录:    ${OUTPUT_DIR}"
echo "输出文件:    ${OUTPUT_DIR}/init_${GROUP}.sql"
echo ""

# 确认执行
read -p "确认生成 SQL 脚本? (y/N): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo -e "${YELLOW}已取消${NC}"
    exit 0
fi

# 确保输出目录存在
mkdir -p "$OUTPUT_DIR"

# 编译并运行 Go 程序
echo -e "${GREEN}正在生成 SQL 脚本...${NC}"
cd "$SQLGEN_DIR"

# 检查 Go 是否安装
if ! command -v go &> /dev/null; then
    echo -e "${RED}错误: Go 未安装${NC}"
    echo "请安装 Go: https://golang.org/dl/"
    exit 1
fi

# 编译 Go 程序
go build -o sqlgen main.go

# 运行程序
# 注意：在 Windows 环境下，如果遇到中文乱码问题，请使用 chcp 65001 设置代码页为 UTF-8
# 或者在 PowerShell 中设置：$OutputEncoding = [System.Text.Encoding]::UTF8
./sqlgen -group "$GROUP" -name "$NAME" -domain "$FRONTEND_DOMAIN" -output "$OUTPUT_DIR" -template "${SQLGEN_DIR}/templates" -parent-id "$PARENT_ID" -parent-path "$PARENT_PATH"
SQLGEN_EXIT_CODE=$?

# 清理编译产物
rm -f sqlgen

if [ $SQLGEN_EXIT_CODE -eq 0 ]; then
    echo -e "${GREEN}✓ SQL 脚本生成成功!${NC}"
    echo -e "${YELLOW}注意:${NC}"
    echo -e "  - 生成的 SQL 文件:"
    echo -e "      - 建表: ${OUTPUT_DIR}/create_table_${GROUP}.sql"
    echo -e "      - 初始化: ${OUTPUT_DIR}/init_${GROUP}.sql"
    echo -e "  - 请在数据库中按顺序执行建表和初始化 SQL"
    echo -e "  - 默认菜单父目录为临时目录 /temp，如需挂到系统管理请使用 -parent-path /system"
    echo -e "  - 可根据需要修改菜单、按钮、接口的启用状态"
else
    echo -e "${RED}✗ SQL 脚本生成失败${NC}"
    exit 1
fi

