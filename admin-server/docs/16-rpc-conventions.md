# 16 — RPC 约定（zrpc 服务间通信）

> 前置依赖：已读 `15-service-boundaries.md`（6 服务边界、5 个 MySQL schema、跨服务引用 ID 的三种模式）。本文档只处理"服务之间怎么通信"，不重复服务边界的决策理由。

## 0. 前置依赖

- Part A 已完成：`internal/wire/providers.go` 里的 Wire 组合根、`registry.Domain`、`repository.Repository` 已经是独立可注入的节点（本文档第 3 节直接在这个基础上扩展，不是重新设计 Wire）。
- 阅读过 `AGENTS.md` 第 3 节后端规范——RPC 服务内部的 Logic/Repository 分层、squirrel-only、错误处理约定（`pkg/errs`）**原样适用于每个新 RPC 服务**，不因为换了传输协议（HTTP → gRPC）而改变。唯一的区别是 RPC 服务的入口不再是 `httpx.Parse` 而是 pb 生成的请求结构体，错误也不再走 `pkg/errs` 转 HTTP 状态码，而是走 gRPC status code（细节见第 5 节）。

## 1. 仓库结构：单 Go module，多个 main 二进制

不做多 module monorepo（对独立维护者而言，多 module 的 `replace`/`go.work`/跨仓库版本管理成本远大于收益）。拆分后的目录结构：

```
admin-server/
  go.mod                    (module 路径不变：postapocgame/admin-server)
  cmd/gateway/               (原 admin.go 的继任者，纯 HTTP 入口，只剩 wire 装配 + zrpc.MustNewClient)
  services/
    iam/
      rpc/
        iam.proto
        pb/                 (protoc 生成的 *.pb.go，禁止手改)
        iamclient/          (goctl 生成的 client 包，gateway 侧 import 这个)
      internal/
        config/
        svc/
        logic/
        server/
      iam.go                (原 goctl rpc 生成的 main.go)
      etc/iam.yaml
    content/  chat/  task/  sdk/   (同构目录结构)
  internal/                 (拆分完成后只剩 gateway 自己的 handler/types/middleware/wire，不再有 9 个业务域的 logic/repository)
  pkg/                       (跨服务共享、不含业务逻辑：errs/jwt/cache/response/consts，不变)
  db/services/<service>/     (见 15 文档第 4 节)
```

`services/<name>/internal/` 内部的分层完全复刻现在单体的 `internal/{logic,repository,model,domain}` 结构，只是 scope 收窄到这一个服务自己的表和业务——例如 `services/task/internal/repository/` 只有 `task_repository.go` 一个文件，不会有 `iam/`、`chat/` 这些子目录,因为 task-rpc 不直接访问其他服务的表。

`internal/` 顶层收缩为 gateway 专属之后，`ServiceContext` 也要瘦身（见第 4 节）——这一步是"挪目录"而不是"重新设计"：现在 `internal/logic/<domain>/<module>/` 下的手写方法体，凡是被拆到某个 RPC 服务的，原样搬进 `services/<name>/internal/logic/<module>/`，包名和方法签名基本不变，只是 `svc.ServiceContext` 换成该服务自己的（不再是 gateway 的大 ServiceContext）。

## 2. 工具链前置依赖

拆分服务需要额外安装以下工具（Part A 的单体阶段不需要）：

```bash
# protoc（protobuf 编译器本体，非 Go 工具，走系统包管理器）
# macOS
brew install protobuf
# Linux (Debian/Ubuntu)
apt-get install -y protobuf-compiler
# 验证版本 >= 3.21
protoc --version

# protoc-gen-go / protoc-gen-go-grpc（Go 插件，goctl rpc 内部会调用）
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# 确认 $GOPATH/bin 在 PATH 里，否则 goctl rpc 找不到这两个插件
export PATH="$PATH:$(go env GOPATH)/bin"
```

`goctl` 本身已经是 Part A 阶段的依赖（`generate-model.sh`/`generate-api.sh` 都在用），版本要求与 `internal/handler/routes.go` 文件头注释标注的版本一致（当前是 goctl 1.10.1，与 `internal/logic/**/*.go` 顶部注释里看到的版本一致）；`goctl rpc` 是 `goctl` 自带的子命令，不需要单独安装。

## 3. `scripts/generate-rpc.sh`

`.template/rpc/*.tpl` 目录**已经存在**（`goctl template init --home .template` 生成过，但目前从未被任何脚本引用），包含 10 个模板文件：`call.tpl`、`config.tpl`、`etc.tpl`、`logic.tpl`、`logic-func.tpl`、`main.tpl`、`server.tpl`、`server-func.tpl`、`svc.tpl`、`template.tpl`（`template.tpl` 是 `goctl rpc new` 脚手架用的示例 proto，实际生成走的是 `goctl rpc protoc` 走已有 `.proto` 文件,不使用这个模板）。这些模板是标准 go-zero rpc 模板的直接产物，没有做过项目定制（不像 `.template/api`/`.template/model` 那样已经改过软删除/统一时间戳的部分），Phase 2 落地时如果需要在生成代码里统一加什么（比如所有 Logic 构造函数都注入 trace 相关字段），要在这批模板上改,而不是生成后手改。

新增 `scripts/generate-rpc.sh`，参数解析、GOCTL_BIN 解析、彩色输出、确认执行的风格与现有的 `scripts/generate-model.sh` 保持一致（同一套 "用户亲自执行" 的政策——AI 可以在开发期直接跑，但脚本本身仍然要有确认提示，供用户在非开发期/生产前审阅时使用）：

```bash
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

"$GOCTL_BIN" rpc protoc "$PROTO_FILE" \
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
```

与 `generate-model.sh` 的一致点：`set -e`、同一套颜色变量、`SCRIPT_DIR`/`PROJECT_ROOT` 解析方式、`GOCTL_BIN` 三级解析（环境变量 → PATH → GOPATH/bin）、生成前 `read -p` 二次确认、成功/失败用 ✓/✗ 打印。差异点：`generate-model.sh` 判断 `$?` 时直接测的是 `goctl` 命令本身的退出码（因为它是脚本里最后一条命令），`generate-rpc.sh` 同样紧跟在 `goctl rpc protoc` 后面立刻存到 `GEN_EXIT` 变量,不经过任何中间命令，避免重蹈计划里提到的 `generate-sql.sh`（`rm -f sqlgen` 之后判断 `$?`，实际取的是 `rm` 的退出码）那个坑。

## 4. 第一个要写的 proto：`iam-rpc`

所有其他服务都要靠 `iam-rpc` 做权限校验，因此它是第一个要落地的 proto。核心方法对应 `internal/middleware/permissionmiddleware.go`（Phase 1 提取出来的 `PermissionResolver.CanAccess`）的直接 RPC 化，加上用户资料查询和 token 黑名单查询（黑名单实际上**不需要**出现在这个 proto 里，见第 6 节说明,这里为了展示完整性仍给出注释掉的原因）：

```protobuf
syntax = "proto3";

package iam;

option go_package = "./iam";

// CheckPermissionRequest 对应现有 PermissionResolver.CanAccess(ctx, userID, method, path)
message CheckPermissionRequest {
  uint64 user_id = 1;
  string method = 2;
  string path = 3;
}

message CheckPermissionResponse {
  bool allowed = 1;
  // 拒绝时的原因，便于网关侧写审计日志，不用于用户可见文案
  string reason = 2;
}

message GetUserProfileRequest {
  uint64 user_id = 1;
}

message UserProfile {
  uint64 id = 1;
  string username = 2;
  string nickname = 3;
  string avatar = 4;
  uint64 department_id = 5;
  int32 status = 6;
}

message GetUserProfileResponse {
  UserProfile profile = 1;
}

message BatchGetUserProfilesRequest {
  repeated uint64 user_ids = 1;
}

message BatchGetUserProfilesResponse {
  // key: user_id（proto3 map 的 key 只能是标量类型，用 string 存 uint64 的十进制表示）
  map<string, UserProfile> profiles = 1;
}

// IsTokenBlacklisted 保留在 proto 里仅为接口完整性；
// 实际调用方（gateway 的 AuthMiddleware）不应该走这个 RPC，
// 应该直接查共享 Redis（见 16-rpc-conventions.md 第 6 节）。
// 只有当某个服务确实没有直连 Redis 的场景时才退回这个 RPC。
message IsTokenBlacklistedRequest {
  string token = 1;
}

message IsTokenBlacklistedResponse {
  bool blacklisted = 1;
}

service Iam {
  rpc CheckPermission(CheckPermissionRequest) returns (CheckPermissionResponse);
  rpc GetUserProfile(GetUserProfileRequest) returns (GetUserProfileResponse);
  rpc BatchGetUserProfiles(BatchGetUserProfilesRequest) returns (BatchGetUserProfilesResponse);
  rpc IsTokenBlacklisted(IsTokenBlacklistedRequest) returns (IsTokenBlacklistedResponse);
}
```

`BatchGetUserProfiles` 直接服务于 `15-service-boundaries.md` 第 5 节模式①的两个真实调用方：`chat-rpc` 展示聊天列表/群成员（`internal/logic/chat/{chat,group,message}/*.go` 的 8 个文件）、`content-rpc` 展示博客作者信息（`public_blog_author_info_logic.go`，目前写死查 `userID=1`，拆分后改成读 gateway 传下来的实际作者 ID）。`CheckPermission` 的请求参数直接对应 `PermissionResolver.CanAccess` 现有的三个参数（`userID`/`method`/`path`），迁移时字段名不用改。

其余 4 个服务（content/chat/task/sdk）的 proto 在各自 `services/<name>/rpc/<name>.proto` 里定义，方法集合直接来自各服务现有 logic 层对外暴露的能力（1:1 对应现有 `.api` 定义的接口),不在本文档展开,由 `18-service-extraction-runbook.md` 每次拆分时现场编写。

## 5. gateway 变薄机制

`internal/wire/providers.go` 现在的 `provideServiceContext` 把 `Repository`、`Domain`、`ChatHub`、`TaskExecutors`、`TaskScheduler` 一次性塞进 `svc.ServiceContext`。Phase 2 完成后，这个函数签名会大幅收缩——`Repository`/`Domain`/`ChatHub`/`TaskExecutors`/`TaskScheduler` 这些指向具体业务数据/领域服务的字段全部消失（它们连同对应的表一起搬进了各自的 RPC 服务），换成 5 个 zrpc client 字段：

```go
// internal/svc/servicecontext.go（gateway 侧，拆分完成后的形态）
type ServiceContext struct {
	Config config.Config

	IamRPC     iamclient.Iam
	ContentRPC contentclient.Content
	ChatRPC    chatclient.Chat
	TaskRPC    taskclient.Task
	SdkRPC     sdkclient.Sdk

	// 中间件字段维持不变的扁平结构（AGENTS.md 明令禁止嵌套），
	// 但中间件内部实现从"直接查 Repository/Domain"改成"调用对应的 zrpc client"
	AuthMiddleware         rest.Middleware
	PermissionMiddleware   rest.Middleware
	OperationLogMiddleware rest.Middleware
	// ... 其余中间件字段不变
}
```

对应的 Wire provider（在 `internal/wire/providers.go` 里新增,替换掉 `provideRepository`/`provideDomain`/`provideChatHub`/`provideTaskExecutors`/`provideTaskScheduler` 这 5 个 provider）：

```go
func provideIamRPC(c config.Config) iamclient.Iam {
	return iamclient.NewIam(zrpc.MustNewClient(c.IamRPCConf))
}

func provideContentRPC(c config.Config) contentclient.Content {
	return contentclient.NewContent(zrpc.MustNewClient(c.ContentRPCConf))
}

func provideChatRPC(c config.Config) chatclient.Chat {
	return chatclient.NewChat(zrpc.MustNewClient(c.ChatRPCConf))
}

func provideTaskRPC(c config.Config) taskclient.Task {
	return taskclient.NewTask(zrpc.MustNewClient(c.TaskRPCConf))
}

func provideSdkRPC(c config.Config) sdkclient.Sdk {
	return sdkclient.NewSdk(zrpc.MustNewClient(c.SdkRPCConf))
}
```

`config.Config` 新增 5 个 `zrpc.RpcClientConf` 字段（`IamRPCConf`/`ContentRPCConf`/`ChatRPCConf`/`TaskRPCConf`/`SdkRPCConf`），`etc/gateway.yaml` 对应新增 5 段配置。**初期用静态 `Endpoints`（host:port 列表），不引入 etcd 服务发现**：

```yaml
IamRPCConf:
  Endpoints:
    - 127.0.0.1:8081
  NonBlock: true
ContentRPCConf:
  Endpoints:
    - 127.0.0.1:8082
  NonBlock: true
# ... 其余 3 个同构
```

这与 `11-descoped.md` 里"不引入 etcd 服务发现——初期用静态 zrpc target，需要时可无痛升级"的决定一致：`zrpc.RpcClientConf` 本身同时支持 `Endpoints`（静态）和 `Etcd`（服务发现）两种配置方式，是否切换只是改配置文件，代码不用动。

**约 120 个纯 CRUD logic 文件是 1:1 机械替换**：把 `svcCtx.Domain.X.Y(...)` 或内联的 `xxxrepo.NewXxxRepository(svcCtx.Repository)` 调用替换成 `svcCtx.<Service>RPC.<Method>(ctx, &pb.XxxRequest{...})`，请求/响应之间做字段映射,不改变业务逻辑本身。Phase 1 已识别的 35-40 个编排类文件（跨表写、跨仓储读写、非平凡业务规则）不是"替换成 RPC 调用"，而是把整段业务逻辑物理搬到对应服务自己的 `logic/` 包里（因为它依赖的数据库也搬过去了），gateway 侧只剩"解析 HTTP 请求 → 拼一次 RPC 请求 → 映射 RPC 响应成 HTTP 响应"这三步薄胶水。

## 6. RBAC/权限校验跨服务方案

**推荐方案（不是选项罗列）：gateway 同步调用 `iam-rpc.CheckPermission`，走缓存，不把权限塞进 JWT。**

不用 JWT 内嵌权限的原因：管理后台场景下，管理员撤销权限应该"下一次请求就生效"，不应该等到 token 过期（当前 access token 过期时间较长）——这是当前系统没有的过期漏洞，不值得为了省一点延迟引入。

具体到缓存：`pkg/cache/business_cache.go` 里已经定义、但目前**从未被任何 logic 调用**的 `CacheKeyUserPermissions`（key 格式 `cache:user:permissions:<userID>`，30 分钟 TTL,即 `CacheExpireUserPermissions = 30*60`）,配套的 `BusinessCache.GetUserPermissions(ctx, userID)`/`SetUserPermissions(ctx, userID, permissions)`/`DeleteUserPermissions(ctx, userID)` 三个方法也已经写好，Phase 2 落地 `iam-rpc.CheckPermission` 内部实现时直接接上这套现成 API：

```go
// services/iam/internal/logic/permission/check_permission_logic.go（示意）
func (l *CheckPermissionLogic) CheckPermission(in *pb.CheckPermissionRequest) (*pb.CheckPermissionResponse, error) {
	permissions, err := l.svcCtx.BusinessCache.GetUserPermissions(l.ctx, in.UserId)
	if err == cache.ErrCacheMiss {
		// 缓存未命中：走原有 PermissionResolver 的权限集合计算逻辑，
		// 计算完写回缓存，一条缓存服务所有路由（不是按 method+path 分别缓存）
		permissions, err = l.resolver.ComputeUserPermissions(l.ctx, in.UserId)
		if err != nil {
			return nil, err
		}
		_ = l.svcCtx.BusinessCache.SetUserPermissions(l.ctx, in.UserId, permissions)
	} else if err != nil {
		return nil, err
	}
	allowed := matchPermission(permissions, in.Method, in.Path)
	return &pb.CheckPermissionResponse{Allowed: allowed}, nil
}
```

RBAC 变更（角色分配权限变化、用户角色变化）时要主动调用 `DeleteUserPermissions` 让缓存失效，不只靠 30 分钟 TTL 兜底——这一点在 `18-service-extraction-runbook.md` 的 iam-rpc 附录里要落实成"角色/用户角色的 Update/Delete logic 收尾时调一次失效"的具体改动点。

**Token 黑名单检查完全不需要走 RPC**：`internal/repository/iam/token_blacklist_repository.go` 的 `IsBlacklisted`/`Blacklist` 两个方法 100% 基于 `r.repo.Redis.Exists`/`r.repo.Redis.Setex`，没有触碰任何 MySQL 表。Redis 在所有服务间共享（见下），gateway 的 `AuthMiddleware` 直接构造一个 `tokenBlacklistRepository` 风格的薄封装直连共享 Redis 即可,热路径零 RPC——这也是第 4 节 proto 里 `IsTokenBlacklisted` 方法标注"不建议实际调用"的原因。

**Redis 保持全服务共享（不像 MySQL 那样拆分）**：Redis 里的东西都是缓存/锁/队列，不是系统记录数据，天然适合共享；已有的 key 前缀（`jwt:blacklist:`、`rate_limit:global`/`rate_limit:ip:`/`rate_limit:user:`/`rate_limit:api:`、`task:lock:`/`task:config:`）不需要重新设计,`pkg/cache`/`internal/consts` 里的这些常量原样复制到需要用到它们的每个服务（不是把 `internal/consts` 整体挪成共享包也不必要,量很小,直接复制维护成本可忽略）。

## 7. Chat 的 WebSocket ↔ gRPC 双向流桥接

**明确决定**：gateway 继续终结 WebSocket 连接（保持"只有 gateway 有公网端口"这条不变式，不用给 nginx 配第二个公网路由），每个 WS 连接通过一条 gRPC 双向流桥接到 `chat-rpc`，`ChatHub`（现在的 `internal/hub/chathub.go`）的实际权威搬到 chat-rpc 里,gateway 只做协议转换,不再持有在线连接的业务状态（谁在线、谁在哪个群）。

proto 契约：

```protobuf
// services/chat/rpc/chat.proto（节选，仅展示 Stream 部分）
message ClientFrame {
  oneof payload {
    JoinFrame join = 1;         // 连接建立后第一帧：携带已鉴权的 user_id
    SendMessageFrame send = 2;  // 发消息
    PingFrame ping = 3;         // 心跳
  }
}

message ServerFrame {
  oneof payload {
    MessageFrame message = 1;   // 有新消息推送给这个连接
    AckFrame ack = 2;           // 消息发送确认
    PongFrame pong = 3;         // 心跳回应
    ErrorFrame error = 4;       // 服务端错误（如发消息失败）
  }
}

message JoinFrame {
  uint64 user_id = 1;
}

message SendMessageFrame {
  uint64 chat_id = 1;
  string content = 2;
  int32 message_type = 3; // 对应 chat_message.message_type：1文本 2图片 3文件
}

message MessageFrame {
  uint64 chat_id = 1;
  uint64 from_user_id = 2;
  string content = 3;
  int32 message_type = 4;
  int64 created_at = 5;
}

service Chat {
  rpc Stream(stream ClientFrame) returns (stream ServerFrame);
  // 其余非流式方法（群管理、聊天列表等常规 CRUD）与其他服务一样走普通 unary RPC，不在这里展开
}
```

网关侧桥接 goroutine 的形状（`internal/handler/chat/ws_bridge_handler.go`，示意，不是最终代码）：

```go
func (h *ChatWSHandler) ServeWS(w http.ResponseWriter, r *http.Request) {
	// 1. 常规 WS upgrade（复用现有 gorilla/websocket 或标准库升级逻辑）
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	// 2. 从已通过 AuthMiddleware 鉴权的请求上下文取 user_id
	user, ok := jwthelper.FromContext(r.Context())
	if !ok {
		conn.WriteMessage(websocket.CloseMessage, nil)
		return
	}

	// 3. 建立到 chat-rpc 的 gRPC 双向流
	stream, err := h.svcCtx.ChatRPC.Stream(r.Context())
	if err != nil {
		logx.Errorf("建立 chat-rpc 流失败: %v", err)
		return
	}

	// 4. 第一帧发 JoinFrame 完成"登记"
	if err := stream.Send(&pb.ClientFrame{Payload: &pb.ClientFrame_Join{
		Join: &pb.JoinFrame{UserId: user.UserID},
	}}); err != nil {
		return
	}

	errCh := make(chan error, 2)

	// 5a. WS → gRPC：读 WS 帧，转成 ClientFrame 发给 chat-rpc
	go func() {
		for {
			_, data, err := conn.ReadMessage()
			if err != nil {
				errCh <- err
				return
			}
			frame, err := decodeClientFrame(data) // WS 传的 JSON，转成 pb.ClientFrame
			if err != nil {
				continue // 单帧解析失败不断连，记日志即可
			}
			if err := stream.Send(frame); err != nil {
				errCh <- err
				return
			}
		}
	}()

	// 5b. gRPC → WS：读 chat-rpc 推来的 ServerFrame，转成 WS 帧写回客户端
	go func() {
		for {
			frame, err := stream.Recv()
			if err != nil {
				errCh <- err
				return
			}
			data, err := encodeServerFrame(frame) // pb.ServerFrame 转 JSON
			if err != nil {
				continue
			}
			if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
				errCh <- err
				return
			}
		}
	}()

	// 6. 任一方向出错（含正常断连）就结束这条连接的桥接，
	// defer conn.Close() 和 stream 的取消都会级联清理
	<-errCh
}
```

两个方向各一个 goroutine、一个共享的 `errCh` 做退出协调，是双向流桥接的标准形状——工作量可控（计划估算约 150-250 行，含 `decodeClientFrame`/`encodeServerFrame` 的编解码和错误处理细节），换来的是"只有 gateway 有公网端口"这条架构不变式继续成立。`chat-rpc` 侧 `Stream` 方法的实现拿到 `stream.Context()` 后的第一件事是从 `JoinFrame` 拿 `user_id`，把这条 gRPC stream 注册进原来 `ChatHub` 的连接表（数据结构不变，只是"连接"从 `*websocket.Conn` 换成 `grpc.ServerStream`），后续推送消息时遍历连接表调 `stream.Send` 而不是 `conn.WriteMessage`。

## 8. 非目标

- 不引入 etcd/Consul 等服务发现——静态 `Endpoints`，见第 5 节。
- 不做 gRPC 层面的自动重试/熔断框架——`zrpc.RpcClientConf` 的默认超时/重试行为够用，出现真实问题再评估引入 `go-zero` 自带的 breaker（默认已开启，不需要额外配置）。
- 不做 REST/gRPC 双协议网关（gRPC-Gateway 之类）——服务间只用 gRPC，对外仍然只有 gateway 一个 HTTP 入口，没有必要让每个 RPC 服务自己再暴露一份 HTTP。
- 不把 `pkg/consts` 里的 Redis key 前缀常量拆成独立共享模块——直接复制到各服务自己的 `internal/consts` 包，维护成本可忽略（见第 6 节）。

## 9. 完成的定义

- `protoc`/`protoc-gen-go`/`protoc-gen-go-grpc` 三个工具链依赖已安装并可用，`goctl rpc protoc` 能跑通一次最小 proto（如 `template.tpl` 自带的 Ping/Pong 示例）验证模板确实生效。
- `scripts/generate-rpc.sh` 落地，风格审阅通过（对照 `generate-model.sh` 逐条核对：颜色输出、GOCTL_BIN 解析、确认提示、退出码处理）。
- `iam.proto` 编译通过（`protoc --go_out=... iam.proto` 无报错），字段命名与现有 `PermissionResolver.CanAccess` 参数、`iam.AdminUser` 字段保持语义一致。
- gateway 侧 `ServiceContext` 的 5 个 zrpc client 字段通过 Wire 装配成功，`go build ./...` 通过。
- Chat WS↔gRPC 桥接跑通一次真实的端到端连接（本地起 gateway + chat-rpc，浏览器/wscat 连接，发一条消息，另一个连接收到）。
