# 20. API 文档生成：goctl-swagger

> 本文档是可直接执行的任务说明，对应总纲 Part C.2（Phase 3，约 Week 13-14+）。执行者在改动前应完整阅读一遍。改完要跑通「完成的定义」里列出的验证步骤。

## 0. 前置依赖

- `api/admin.api` 是 gateway 的 API 真源（Phase 2 拆分后，gateway 仍然是唯一持有 `.api` 文件的服务，五个 RPC 服务用 `.proto`，不用 `.api`）。
- 本文档只覆盖 gateway 的 HTTP API 文档生成；RPC 服务的文档策略见第 3 节，不需要新工具。

## 1. 为什么选 `goctl-swagger`

goctl 官方工具链本身不内置 `.api → OpenAPI` 的转换（`goctl api go`/`goctl api ts` 分别生成后端骨架和前端类型，但不生成 OpenAPI JSON），社区插件 `goctl-swagger`（`github.com/zeromicro/goctl-swagger`）通过 goctl 的插件机制（`goctl api plugin -plugin <plugin-name>`）读取同一份 `.api` AST，输出标准 OpenAPI 2.0（Swagger）JSON。

选它而不是手写/其他工具的原因：
- **单一真源不变**——文档直接从 `admin.api` 生成，`.api` 文件改了，重新跑一次生成脚本文档就同步，不需要手工维护一份平行的 OpenAPI 文件（那样几乎必然很快就会和真实接口对不上）。
- **和现有生成脚本是同一套机制**——`goctl api plugin` 只是 `goctl api go`/`goctl api ts` 用的同一个 `goctl api` 命令族的插件模式，接入成本低，不引入新的构建工具链。

## 2. `goctl-swagger` 工作流

### 2.1 安装

```bash
go install github.com/zeromicro/goctl-swagger@latest
```

产出的 `goctl-swagger` 二进制需要在 `PATH` 上（或通过环境变量指定路径，见第 2.3 节脚本设计），供 `goctl api plugin` 以子进程方式调用。

### 2.2 命令形态

```bash
goctl api plugin -plugin goctl-swagger="swagger -filename admin.json" -api admin.api -dir .
```

`-plugin` 参数的值是"插件二进制名 + 插件自己的子命令和参数"，goctl 把它当一整段字符串解析后转发给插件进程；`-api`/`-dir` 是 goctl 插件机制的标准参数（指定输入 `.api` 文件和输出目录）。

### 2.3 新增 `scripts/generate-swagger.sh`

设计上完全follow `scripts/generate-model.sh`/`scripts/generate-ts.sh` 已经确立的约定（三个脚本读起来应该是"同一个人写的"，不引入新的参数风格）：

- **`GOCTL_BIN` 解析三段式**：优先环境变量 `GOCTL_BIN`，其次 `PATH` 上的 `goctl`，再尝试 `$(go env GOPATH)/bin/goctl`；找不到就报错退出，提示安装命令——`generate-model.sh` 第 21-33 行、`generate-ts.sh` 第 21-33 行是完全相同的一段，`generate-swagger.sh` 原样复用。
- 新增同样风格的 `GOCTL_SWAGGER_BIN` 解析（`PATH` 上的 `goctl-swagger`，或 `$(go env GOPATH)/bin/goctl-swagger`），找不到时提示 `go install github.com/zeromicro/goctl-swagger@latest`。
- **颜色输出**：复用 `RED`/`GREEN`/`YELLOW`/`NC` 四个 ANSI 变量定义（两个现有脚本里逐字相同的一段）。
- **参数解析**：位置参数为可选的 `.api` 文件路径（不传则默认 `api/admin.api`，与 `generate-ts.sh` 默认使用 `admin.api` 的逻辑一致），`-h/--help` 打印用法。
- **确认执行**：`read -p "确认生成 Swagger 文档? (y/N): "` 交互确认，与两个现有脚本的确认交互逐字一致的风格——这是"用户亲自执行"政策的落地方式（见 `10-dev-execution-and-review-points.md`）：脚本本身允许开发期 AI 直接跑（`10-dev-execution-and-review-points.md` 第 1 条已经把 `generate-*.sh` 的实际执行列为"可以直接做，事后 review"），但脚本设计上保留交互确认，与其余 `generate-*.sh` 保持同一套用户预期。
- **输出目标固定为 `docs/openapi/admin-api.json`**（而不是像 model/ts 那样可通过 `-d`/`-dir` 自定义输出目录）——API 文档只有一份、不需要多目的地，固定路径让"文档在哪"这件事不需要每次确认。

脚本骨架（依照上述约定的完整实现，非伪代码）：

```bash
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
```

### 2.4 输出目标

`docs/openapi/admin-api.json`——新目录，只放这一个生成产物（未来如果给 `sdk-rpc` 单独做 swagger，见第 3 节，产物放同一目录下，如 `docs/openapi/sdk-api.json`）。这个 JSON 文件本身是生成产物，跟 `internal/wire/wire_gen.go`、前端 `src/api/generated/` 同一类性质：**不手改**，改动源头是 `admin.api`，改完重新跑脚本。是否提交进 git 由用户决定（生成产物提交策略与 `wire_gen.go` 一致，建议提交，方便不跑脚本的场景下也能看到最新文档），本文档不强制。

## 3. RPC 服务不做单独文档生成工具

`iam-rpc`/`content-rpc`/`chat-rpc`/`task-rpc`/`sdk-rpc` 五个服务**不配置任何 swagger/OpenAPI 类工具**。理由：
- 它们是纯内部服务，调用方只有 gateway 和彼此（B.3 的服务边界设计），不存在外部消费者需要一份可视化 API 文档。
- `.proto` 文件本身（配好字段注释、方法注释）就是文档——这是 gRPC/protobuf 生态的标准做法，`protoc` 生态本身有一堆基于 `.proto` 注释生成文档的工具（如需要可以后续接，本轮不做），不需要重新发明一套。
- 维护两套文档生成工具链（`.api → swagger` + `.proto → 某种 rpc 文档`）对一个内部服务而言是不成比例的投入。

`16-rpc-conventions.md` 里对 `.proto` 文件字段/方法注释的约定，直接就是 RPC 服务的"文档规范"，不需要在本文档重复。

## 4. 优先级提示：`sdk-rpc` 是第一个值得做真 Swagger UI 的候选

`sdk-rpc` 承载的是 API Key 鉴权的对外调用（`SDKAuthMiddleware`/`SDKRateLimitMiddleware`/`SDKCallLogMiddleware`，见 `15-service-boundaries.md` B.1 的拆分理由），信任边界和调用方性质与其余五个服务不同——它是唯一可能在未来有真正外部调用方（不是内部管理员，是第三方通过 API Key 调用）的服务。

一旦 `sdk-rpc` 有真实外部调用方，它是第一候选去接一个真正可交互的 Swagger UI（例如 `swagger-ui` 静态页面挂载生成的 JSON，或迁移到 gateway 侧为 `sdk` 相关路由单独出一份文档）——但这是**未来触发条件成立后才做**的事，不在本轮范围内，本轮只交付 `generate-swagger.sh` 覆盖 gateway 的 `admin.api` 全量文档生成能力。

## 5. 完成的定义

- `scripts/generate-swagger.sh` 落地，参数解析/颜色输出/确认交互风格与 `generate-model.sh`/`generate-ts.sh` 一致。
- 本地安装 `goctl-swagger` 后跑通一次生成，产出 `docs/openapi/admin-api.json`，是合法 JSON（`jq . docs/openapi/admin-api.json` 不报错）。
- 抽查 JSON 内容：至少能看到几个真实存在的路由（如 `/api/v1/auth/login`）出现在生成的文档里，字段/参数与 `admin.api` 定义一致。
- `docs/openapi/` 目录新增，`.gitignore` 视用户决定是否收录该目录（本文档不强制，默认建议提交）。
- 确认没有为任何 RPC 服务新增文档生成脚本或依赖。
