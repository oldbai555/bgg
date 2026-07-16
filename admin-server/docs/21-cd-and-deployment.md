# 21. CD 与部署：镜像化 + docker-compose 生产切换

> 本文档是可直接执行的任务说明，对应总纲 Part C.3（Phase 3，约 Week 13-14+）。执行者在改动前应完整阅读一遍。改完要跑通「完成的定义」里列出的验证步骤。

## 0. 前置依赖

- Part B 六个部署单元（`gateway` + 5 个 RPC 服务）已经拆分完成，各自能独立 `go build`。
- `09-ci-cd-and-deployability.md`（Phase 1 产出）里已经有单体阶段的 Dockerfile 设计，本文档直接复用同一份模板，不重新设计 Dockerfile 结构。

## 1. 现状：`script/admin.sh` 的真实部署机制

现状核查（`script/admin.sh`，工作区根目录 `script/`，不是 `admin-server/scripts/`）。这是一个约 260 行的 bash 脚本，`main()` 分发四组子命令：`build`/`package`/`frontend`/`supervisor`。当前后端的完整部署链路：

1. **`build_server()`**（第 63-79 行）：`cd admin-server && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w -X main.GIT_COMMIT_VERSION=$git_version" -o dist/admin-server admin.go`，构建产物 + `etc/admin-api.yaml` 一起放进 `dist/`。单一二进制，`git rev-parse --short HEAD` 作为版本号注入。
2. **`package_server()`**（第 108-119 行）：调用 `build_server`，再把 `dist/` 打包成 `package/admin-server_<git-sha>.tar.gz`。
3. **`supervisor_deploy(package_file, app_name)`**（第 190-212 行）：`tar -xzf` 解包到 `$SUPERVISOR_DIR/$app_name`（默认 `/home/work/service/admin-server`），生成/复制 Supervisor `.conf`（`supervisor_gen_conf`，第 168-187 行，`command=$service_dir/admin-server -f $service_dir/admin-api.yaml`），执行 `supervisorctl update && supervisorctl restart admin-server`。
4. **`supervisor_install()`**（第 214-221 行）是 1-3 的一键组合：build → package → deploy。

**这个链路里没有 scp 这个动作本身**——脚本本身假设在目标服务器上直接执行（`supervisor_deploy` 直接操作本地文件系统的 `$SUPERVISOR_DIR`），把打包产物"搬到服务器"这一步（无论是 `scp`/`rsync`/手工上传）不在脚本职责内，是使用者的外部动作。核心特征是：**单二进制 + Supervisor 进程管理 + 配置走 `/etc/work/*.json` 外部文件**（`admin.go` 里 `-mysql-config`/`-redis-config`/`-middleware-config` 三个 flag，指向 `/etc/work/mysql.json` 等，不提交仓库）。

这套机制在 Phase 2 之前（单体阶段）继续有效——**Phase 3 不是替换它，是新增一条平行路径**（`admin.sh compose ...`），供 Phase 2 拆分完成、准备切到多服务部署时使用；单体阶段的 `supervisor` 系列命令保留不动。

## 2. 演进目标：docker-compose，不上 Kubernetes

**推荐方案：生产环境从"单二进制 + Supervisor"演进到"六个容器 + docker-compose"，不引入 Kubernetes 或任何编排器/调度器。**

这是一个刻意的非目标，不是能力缺口：

- **场景是一台服务器、一个人**。Kubernetes 解决的核心问题——跨多台机器调度、自动故障转移、滚动升级的复杂编排——在单机场景下完全不存在，引入它只会带来 `kubectl`/`etcd`/`kubelet` 本身的运维负担和学习成本，没有对应的收益。
- **`docker-compose` 已经覆盖这个场景需要的全部能力**：多容器编排（六个服务）、声明式配置（`docker-compose.prod.yml`）、`docker compose pull && docker compose up -d` 一键滚动重启、日志采集（`docker compose logs`，对接 `19-observability.md` 的 JSON stdout 约定）。
- **给未来多机扩展留了后路，不是死路**：因为 Part B 已经把 DSN/配置按服务拆好（`15-service-boundaries.md` 的 5 个独立 schema），六个服务本来就是独立部署单元，真正需要多机时（如果那一天到来），是把 `docker-compose.prod.yml` 换成任何编排方案的问题，不需要回头重新拆分服务边界——现在的选择不构成对未来的锁定。

这一条已经写进 `11-descoped.md`（"不上 Kubernetes"），本文档是这个决策在部署机制设计上的具体落地。

## 3. 六个 Dockerfile：复用同一个模板

每个服务（`gateway`/`iam-rpc`/`content-rpc`/`chat-rpc`/`task-rpc`/`sdk-rpc`）各有一个 `Dockerfile`，但**不是六份独立设计**——全部遵循 `09-ci-cd-and-deployability.md` 里 Phase 1 已经为单体设计的同一份 Dockerfile 结构（该文档复用了仓库已有但当时未使用的 `.template/docker/docker.tpl`：多阶段构建，`golang:{version}-alpine` 编译 → 精简运行时基础镜像，`CGO_ENABLED=0` 静态编译，可选时区数据），六个服务的 Dockerfile 只有两处差异：

```dockerfile
# 差异点 1：构建产物名 / 入口 main 包路径
RUN go build -ldflags="-s -w" -o /app/iam-rpc ./services/iam
# gateway 对应 ./cmd/gateway，其余四个 RPC 服务同构替换服务名

# 差异点 2：是否 EXPOSE HTTP 端口
# 仅 gateway 需要 EXPOSE（对外唯一入口，见 15-service-boundaries.md B.3
# "只有 gateway 有公网端口"这条不变式）；五个 RPC 服务 EXPOSE 的是 zrpc 内部端口，
# 不做公网映射，docker-compose 内部网络即可互通
```

六份 Dockerfile 放在各自服务目录下（`services/iam/Dockerfile`、`cmd/gateway/Dockerfile` 等），与 `18-service-extraction-runbook.md` 里每次拆分一个服务时的产出物对齐——**不是 Phase 3 一次性新增六个，是 Phase 2 拆每个服务时就带一个 Dockerfile，Phase 3 只是把它们统一接进 CI 和 compose**。

## 4. CI：按服务的镜像构建+推送

在现有 GitHub Actions（`09-ci-cd-and-deployability.md` 里 Phase 1 落地的 lint+build/单元测试/集成测试三个 job 之外）新增一个 matrix job，按服务并行构建+推送镜像：

```yaml
build-images:
  needs: [lint-and-build, unit-test]   # 复用 Phase 1 已有 job 作为前置门槛
  strategy:
    matrix:
      service: [gateway, iam, content, chat, task, sdk]
  runs-on: ubuntu-latest
  steps:
    - uses: actions/checkout@v4
    - uses: docker/login-action@v3
      with:
        registry: ghcr.io
        username: ${{ github.actor }}
        password: ${{ secrets.GITHUB_TOKEN }}
    - uses: docker/build-push-action@v5
      with:
        context: ./admin-server
        file: ./admin-server/${{ matrix.service == 'gateway' && 'cmd/gateway' || format('services/{0}', matrix.service) }}/Dockerfile
        push: true
        tags: ghcr.io/<user>/admin-${{ matrix.service }}:${{ github.sha }}
```

镜像命名规则：`ghcr.io/<user>/admin-<service>:<git-sha>`——六个独立镜像仓库（同一 `ghcr.io/<user>` 命名空间下的六个 image name），tag 用完整 git sha（不用 `latest`，保证每次部署可追溯到具体提交，与 `build_server()` 现有用 `git rev-parse --short HEAD` 做版本标识的习惯一致，只是精度从"打进二进制的版本字符串"变成"镜像 tag"）。**这一部分全自动**——每次 push 到主分支触发，不需要人工介入。

## 5. `script/admin.sh` 新增 `compose` 模式

生产发布本身**仍是人工触发的一步**（下一节详述这一点为什么不能自动化），但发布动作要收敛成脚本里的一个命令，不是让使用者手打 `docker compose` 一长串参数。仿照现有 `supervisor` 子命令组的结构（`case "${2:-}" in gen-conf|install|deploy|status|start|stop|restart|logs`），新增 `compose` 子命令组，在 `script/admin.sh` 里新增以下函数（放在 Supervisor 管理那一段之后，主程序分发逻辑之前，和现有函数同一种命名习惯：`动词_名词`、`log_info`/`log_error` 输出、`command -v` 前置检查）：

```bash
# ============================================
# docker-compose 管理（Phase 3 新增，与 Supervisor 模式并存）
# ============================================

COMPOSE_REMOTE_HOST="${COMPOSE_REMOTE_HOST:-}"      # 生产服务器 SSH host（~/.ssh/config 别名或 user@ip）
COMPOSE_REMOTE_DIR="${COMPOSE_REMOTE_DIR:-/home/work/admin-compose}"
COMPOSE_FILE="docker-compose.prod.yml"

compose_pull() {
  local host="${1:-$COMPOSE_REMOTE_HOST}"
  [ -z "$host" ] && { log_error "请指定远程主机（COMPOSE_REMOTE_HOST 或参数）"; return 1; }

  log_info "在 $host 拉取最新镜像..."
  ssh "$host" "cd $COMPOSE_REMOTE_DIR && docker compose -f $COMPOSE_FILE pull"
}

compose_deploy() {
  local host="${1:-$COMPOSE_REMOTE_HOST}"
  [ -z "$host" ] && { log_error "请指定远程主机（COMPOSE_REMOTE_HOST 或参数）"; return 1; }

  compose_pull "$host" || return 1

  log_info "在 $host 应用新镜像..."
  ssh "$host" "cd $COMPOSE_REMOTE_DIR && docker compose -f $COMPOSE_FILE up -d"

  log_info "部署完成，检查状态..."
  compose_status "$host"
}

compose_status() {
  local host="${1:-$COMPOSE_REMOTE_HOST}"
  [ -z "$host" ] && { log_error "请指定远程主机"; return 1; }
  ssh "$host" "cd $COMPOSE_REMOTE_DIR && docker compose -f $COMPOSE_FILE ps"
}

compose_logs() {
  local host="${1:-$COMPOSE_REMOTE_HOST}"
  local service="${2:-}"
  [ -z "$host" ] && { log_error "请指定远程主机"; return 1; }
  ssh "$host" "cd $COMPOSE_REMOTE_DIR && docker compose -f $COMPOSE_FILE logs --tail=200 ${service}"
}
```

`main()` 里加一段和 `supervisor)` 平行的分发：

```bash
compose)
  case "${2:-}" in
    pull) compose_pull "${3:-}" ;;
    deploy) compose_deploy "${3:-}" ;;
    status) compose_status "${3:-}" ;;
    logs) compose_logs "${3:-}" "${4:-}" ;;
    *) log_error "未知命令: compose $2"; usage; exit 1 ;;
  esac
  ;;
```

调用形态：`bash script/admin.sh compose deploy myserver`（或设置 `COMPOSE_REMOTE_HOST` 环境变量后 `bash script/admin.sh compose deploy`）。本地开发同样可以用同一份 `docker-compose.yml`（非 `.prod.yml`，本地版本 mysql/redis 用容器而不是外部依赖，`09-ci-cd-and-deployability.md` 已有本地 compose 设计）跑六个服务，**本地和生产用同一种文件形状**（都是 `services: {gateway, iam-rpc, content-rpc, chat-rpc, task-rpc, sdk-rpc}` 的 compose 结构），不是两套完全不同的部署体系，差异只在 `image:` 字段指向本地 `build:` 还是 `ghcr.io/...:sha` 远程镜像、以及外部依赖（MySQL/Redis）是否容器化。

`docker-compose.prod.yml` 骨架：

```yaml
services:
  gateway:
    image: ghcr.io/<user>/admin-gateway:${TAG:-latest}
    ports: ["20000:20000"]
    volumes: ["./etc/gateway.yaml:/app/etc/gateway.yaml", "/etc/work:/etc/work:ro"]
  iam-rpc:
    image: ghcr.io/<user>/admin-iam:${TAG:-latest}
    volumes: ["./etc/iam.yaml:/app/etc/iam.yaml", "/etc/work:/etc/work:ro"]
  content-rpc: { ... }
  chat-rpc: { ... }
  task-rpc: { ... }
  sdk-rpc: { ... }
```

`/etc/work/*.json`（MySQL/Redis 外部配置文件，`05-密钥管理`/`A.5` 已确认的机制）继续挂载只读卷进容器，不打进镜像——这条约束延续单体阶段就有的"敏感配置不提交仓库、不打进构建产物"的原则。

## 6. 人工触发，不是自动 CD

**镜像构建+推送（第 4 节）全自动**，但**生产环境的 `compose deploy`（第 5 节）是人工触发的一步，不接入 CI 自动执行**。原因和 `10-dev-execution-and-review-points.md` 的"开发期执行策略"里明确写出的第 4 条停下来条件一致：**任何触及真实生产部署动作本身，需要用户拍板时机**——生产发布涉及真实的服务重启窗口、可能的短暂不可用，这是产品/运维判断，不是"技术上能自动化就应该自动化"的事，尤其是 Phase 2 每次服务拆分后是否真的切换到 compose 部署这类决策（见 `10-dev-execution-and-review-points.md` 停下来清单第 4 条）。

**执行主体是人**：使用者在确认好版本、确认好维护窗口之后，手动跑 `bash script/admin.sh compose deploy <host>`；CI 只负责保证"镜像已经构建好、可以随时被拉取"，不负责决定"现在要不要上线"。

## 7. 与上线部署清单的关系

这一步产生的所有部署动作继续记入 `14-production-deployment-checklist.md`（策略不变，只是条目变多），具体新增条目类型：

- 每个服务第一次切到 compose 部署时，对应的迁移顺序（建议：先切风险最低的 `task-rpc`，与 `18-service-extraction-runbook.md` B.6 的拆分顺序一致的保守思路，验证过一次机制再推其余服务）。
- 每个服务的生产配置文件（`etc/<service>.yaml`）与 `/etc/work/*.json` 挂载路径的对应关系。
- 每次发布对应的镜像 tag（git sha），便于回滚时明确"回滚到哪一个 tag"。
- `ghcr.io` 推送凭证（CI 用的 `GITHUB_TOKEN` 已经够用，不需要额外 PAT，除非镜像仓库权限模型另有要求）在生产服务器上拉取镜像时的登录方式（`docker login ghcr.io`，一次性配置，记入清单）。

## 8. 完成的定义

- 六个 `Dockerfile`（`cmd/gateway/Dockerfile` + 五个 `services/<name>/Dockerfile`）落地，均以 `.template/docker/docker.tpl` 的结构为模板，本地 `docker build` 每个都能成功出镜像。
- GitHub Actions 新增 `build-images` matrix job，六个服务并行构建+推送到 `ghcr.io/<user>/admin-<service>:<git-sha>`。
- `script/admin.sh` 新增 `compose pull`/`compose deploy`/`compose status`/`compose logs` 四个子命令，函数命名/日志输出风格与现有 `supervisor_*` 系列一致，原有 `build`/`package`/`frontend`/`supervisor` 命令组不变、不删除。
- `docker-compose.yml`（本地）+ `docker-compose.prod.yml`（生产，`image:` 指向 `ghcr.io`）落地，本地 `docker compose up -d` 能拉起六个服务并互通。
- `14-production-deployment-checklist.md` 补上本节列出的新增条目类型。
- 确认没有引入任何 Kubernetes 相关配置/依赖，`compose deploy` 保持人工触发，未被接入任何自动化 CD 流水线。
