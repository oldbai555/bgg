# 09 · CI/CD 与可部署性（Phase 1 · A.7）

## 前置依赖

- `08-testing-strategy.md`（unit-test / integration-test job 要跑的目标已确定）。
- 本篇只覆盖 **Phase 1（单体）版本**的 Dockerfile / docker-compose / CI。Phase 3 把它扩展成"6 个服务各自出镜像、CI 按服务 build+push"的版本，写在 `21-cd-and-deployment.md`（另一位执行者负责），**本篇不重复、不预写那部分内容**。

## 0. 现状核查

已实地确认：

- `admin-server/.template/docker/docker.tpl` **存在**，是 go-zero `goctl docker` 命令用的模板（`{{.Version}}`/`{{.HasTimezone}}`/`{{.ExeFile}}`/`{{.GoMainFrom}}`/`{{.BaseImage}}`/`{{.Argument}}`/`{{.GoRelPath}}`/`{{.HasPort}}`/`{{.Port}}` 等占位符），当前**没有被使用**——仓库里没有任何 `Dockerfile`。
- 仓库根目录、`admin-server/` 下均**没有 `Dockerfile`**（已用 `find` 确认零匹配）。
- 仓库根目录、`admin-server/` 下均**没有 `docker-compose*.yml`**（已用 `find` 确认零匹配）。
- **没有 `.github/workflows/` 目录**（`.github` 目录本身都不存在）。
- **没有任何 `.golangci.yml`/`.golangci.yaml`**（全仓库搜索零匹配）；已有 `admin-server/.staticcheck.conf`（`checks = ["all", "-SA5008"]`，`SA5008` 关闭是因为 go-zero 的 `optional` 标签会被 staticcheck 误报），这是不同的工具，不能替代 golangci-lint 配置。

以上四项都是"从零搭建"，不是"修复既有配置"。

## 1. Dockerfile

不手写，用 `goctl docker` 从 `.template/docker/docker.tpl` 生成，和仓库"能用 goctl 生成的必须用 goctl 生成"的原则一致（`AGENTS.md` 第 3 节）。生成命令（在 `admin-server/` 目录下执行）：

```bash
goctl docker -go admin.go
```

生成时需要确认/传入的关键参数（对应模板占位符）：
- `-go admin.go`：main 入口，对应模板里的 `{{.GoMainFrom}}`。
- Go 版本：跟 `go.mod` 的 `go 1.24.0` 对齐，生成的 `FROM golang:1.24-alpine AS builder`（`goctl` 默认会读本地 Go 版本，需要人工核对一致，不一致要手改）。
- 时区：`{{.HasTimezone}}` 分支建议开启（`Asia/Shanghai`），避免容器内时间戳与业务预期的秒级时间戳错位。
- `{{.HasPort}}`/`{{.Port}}`：对齐 `etc/admin-api.yaml` 的 `Port`（当前配置需要核实具体值，写 Dockerfile 时从该文件读取，不要猜测硬编码）。
- `{{.Argument}}`（是否带 `-f etc/xxx.yaml` 启动参数）：需要，容器内启动必须显式指定配置文件路径，对应 `COPY {{.GoRelPath}}/etc /app/etc` 这一段会被启用。
- `{{.BaseImage}}`：运行时基础镜像用 `alpine:latest`（配合 `CGO_ENABLED=0` 静态编译，镜像体积小）。

生成后人工核对清单：
1. `ADD go.mod .` / `ADD go.sum .` / `RUN go mod download` 这几行在 builder 阶段之前执行，确认没有把 `scripts/sqlgen/` 下的产物（尤其是待清理的 `sqlgen.exe`，见 `12-scripts-standardization.md`）意外 `COPY . .` 进镜像——`.dockerignore` 需要新增（当前不存在），至少排除 `scripts/sqlgen/sqlgen.exe`、`admin-server.exe`、`*.log`、`.git`。
2. `CMD ["./admin", "-f", "etc/admin-api.yaml"]`（或等价参数）能对上实际的可执行文件名和配置路径。
3. JWT 密钥走环境变量注入（见 A.5 / `01-architecture-target.md`），Dockerfile 本身不 `ENV` 硬编码任何密钥，只在 `docker-compose.yml`/生产 compose 文件里通过 `environment:`/`.env` 注入。

## 2. docker-compose.yml（本地开发）

新建 `admin-server/docker-compose.yml`（仓库当前没有），目标：一条命令拉起 MySQL + Redis + app，供本地开发/集成测试使用。

设计要点：

```yaml
services:
  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: ${MYSQL_ROOT_PASSWORD:-devroot}
      MYSQL_DATABASE: ${MYSQL_DATABASE:-admin}
    ports:
      - "3306:3306"
    volumes:
      - mysql_data:/var/lib/mysql
      - ./db:/db:ro                                          # 整个 db/ 只读挂载，供下面的入口脚本引用
      - ./db/docker-init.sh:/docker-entrypoint-initdb.d/00-init.sh:ro   # 唯一入口，委托给 db/services/init-dev-db.sh 按显式依赖顺序执行，不依赖 MySQL 官方镜像对 initdb.d 下多文件的字典序
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "localhost", "-u", "root", "-p${MYSQL_ROOT_PASSWORD:-devroot}"]
      interval: 5s
      timeout: 5s
      retries: 10

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 5s
      timeout: 5s
      retries: 10

  app:
    build:
      context: .
      dockerfile: Dockerfile
    depends_on:
      mysql:
        condition: service_healthy
      redis:
        condition: service_healthy
    environment:
      JWT_ACCESS_SECRET: ${JWT_ACCESS_SECRET:-local-dev-access-secret-not-for-prod}
      JWT_REFRESH_SECRET: ${JWT_REFRESH_SECRET:-local-dev-refresh-secret-not-for-prod}
    ports:
      - "20000:20000"   # 与 etc/admin-api.yaml 的 Port 一致

volumes:
  mysql_data:
```

要点说明：
- **healthcheck 是硬要求**：`app` 必须等 MySQL/Redis 真正可用（不是容器起来就行，是端口能连、能查询）才启动，否则 `admin.go` 里 `conf.MustLoad`/连接初始化阶段会直接崩溃退出，在 CI 里表现为间歇性失败，排查成本高。
- `db/` 目录整个挂载到 MySQL 的 `docker-entrypoint-initdb.d`：`db/services/` 拆分后已经不裸挂多个散落 `.sql` 文件靠字典序执行，而是 `db/docker-init.sh` 作为唯一入口委托给 `db/services/init-dev-db.sh`（显式依赖顺序：iam 各表 create→init，再按 content/chat/task/sdk 顺序 create→init），字典序问题已解决，不再是需要写 compose 时验证的开放问题。
- 本地开发用的 JWT 密钥给了默认值（`local-dev-*-not-for-prod`），按 `10-dev-execution-and-review-points.md` 的口径可以直接生成使用，注释里必须写清楚"仅供本地开发，生产环境必须通过环境变量覆盖"。

## 3. .golangci.yml（新增，最小化起步）

当前仓库没有任何 golangci-lint 配置。不要一上来对全仓库 512 个文件开满全部 linter（噪音太大，会挡住真正该关注的问题）。Phase 1 起步范围：**只扫 `internal/repository/**` 和 `internal/domain/**`**（本轮改动最集中、最需要保证质量的两层），其余目录留到 Phase 1 Week 4-5"扩大范围"阶段（这一步已经写进总纲的时间线）。

```yaml
run:
  timeout: 5m

linters:
  disable-all: true
  enable:
    - govet
    - staticcheck
    - errcheck
    - ineffassign
    - unused
    - gosimple

issues:
  exclude-dirs-use-default: true

linters-settings:
  staticcheck:
    checks: ["all", "-SA5008"]   # 与 admin-server/.staticcheck.conf 保持一致的例外
```

**范围控制不写在这份 YAML 里**：golangci-lint 的 `issues.include` 是"把默认被排除的某些 issue 代码重新纳入"，不是路径过滤器，写错了不会报配置错误但也不会起到"只扫两个目录"的效果，容易埋一个看起来生效实则没生效的坑。真正的路径范围通过命令行参数传，`.golangci.yml` 本身保持全仓库通用，Phase 1 起步只是**调用方式**上加了路径限制，不是配置文件本身的限制：

```bash
golangci-lint run ./internal/repository/... ./internal/domain/...
```

Phase 1 Week 4-5"扩大范围"、Phase 2/3 随服务拆分继续扩大，都是改这条命令行的参数列表（或者干脆去掉参数变成 `golangci-lint run ./...` 扫全仓库），`.golangci.yml` 本身不用跟着改。

`staticcheck` 的 `SA5008` 例外要和 `admin-server/.staticcheck.conf` 现有配置对齐（第 3 节已确认该文件存在且原因是 go-zero `optional` 标签误报），不要在两处配置里出现矛盾的规则。

## 4. .github/workflows/ci.yml

新增 `.github/workflows/ci.yml`（`.github/workflows/` 目录当前不存在，需要新建）。三个 job，对应 A.7：

```yaml
name: CI

on:
  push:
    branches: [main, feature/**]
  pull_request:
    branches: [main]

jobs:
  lint-build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: "1.24"
      - name: golangci-lint
        uses: golangci-lint-action@v6
        with:
          working-directory: admin-server
          # Phase 1 起步范围，见上方"范围控制不写在 YAML 里"的说明；扩大范围时只改这一行的参数列表
          args: --timeout=5m ./internal/repository/... ./internal/domain/...
      - name: go build
        run: cd admin-server && go build ./...

  unit-test:
    runs-on: ubuntu-latest
    needs: lint-build
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: "1.24"
      - name: go test (unit)
        # 范围严格对齐 08-testing-strategy.md 的口径：只跑领域服务 + repository 层，
        # 不是 go test ./...——全仓库跑会把 handler/logic 里大量"明确不测"的透传代码也跑一遍
        # （本身没坏处但拉长 CI 时间，且给人"这些也被测试覆盖了"的错觉），
        # 集成测试（//go:build integration）默认不会被这条命令带上，因为没传 -tags=integration。
        run: cd admin-server && go test ./internal/domain/... ./internal/repository/... -race -count=1

  integration-test:
    runs-on: ubuntu-latest
    needs: lint-build
    services:
      mysql:
        image: mysql:8.0
        env:
          MYSQL_ROOT_PASSWORD: root
          MYSQL_DATABASE: admin_test
        ports: ["3306:3306"]
        options: >-
          --health-cmd="mysqladmin ping" --health-interval=5s --health-timeout=5s --health-retries=10
      redis:
        image: redis:7-alpine
        ports: ["6379:6379"]
        options: >-
          --health-cmd="redis-cli ping" --health-interval=5s --health-timeout=5s --health-retries=10
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: "1.24"
      - name: 初始化测试库
        run: MYSQL_ROOT_PASSWORD=root MYSQL_DATABASE=admin_test admin-server/db/services/init-dev-db.sh -h127.0.0.1
      - name: go test (integration)
        run: cd admin-server && go test -tags=integration ./... -count=1
        env:
          TEST_MYSQL_DSN: root:root@tcp(127.0.0.1:3306)/admin_test
          TEST_REDIS_ADDR: 127.0.0.1:6379
```

要点：
- `unit-test`/`integration-test` 都 `needs: lint-build`，编译不过或 lint 不过直接短路，不浪费 CI 时间跑测试。
- `integration-test` 用 GitHub Actions 的 `services:` 字段起 MySQL/Redis（这是 CI 环境的标准做法，不同于本地开发用 docker-compose——两者殊途同归，都是"给测试提供真实 MySQL/Redis"，不需要统一成同一份 yaml）。
- 环境变量名（`TEST_MYSQL_DSN`/`TEST_REDIS_ADDR`）是占位建议，落地时要和 `08-testing-strategy.md` 里集成测试实际读取配置的方式对齐，不要出现 CI 传一套名字、代码读另一套名字的错配。
- 测试库初始化已改为调用 `db/services/init-dev-db.sh`（`15-service-boundaries.md` 第 4 节的 `db/services/` 目录拆分落地后的统一入口，替代了本篇最初写的 `db/tables.sql`），与 `.github/workflows/ci.yml` 保持同步。

## 5. 与 Phase 3 的衔接

Phase 3（`21-cd-and-deployment.md`）会把这里的 Dockerfile 从"1 个单体二进制"扩展成"6 个服务各一个"，CI 新增按服务的 `build+push` job（推到 `ghcr.io/<user>/admin-<service>:<git-sha>`），docker-compose 也会有一份"生产切换"版本。**本篇的 Dockerfile/compose/CI 设计是 Phase 3 那份的直接基础，不是要推翻重做的过渡产物**——Phase 3 落地时优先复用这里已验证的 healthcheck/环境变量注入方式，只扩展"单服务→多服务"这一维度。

## 完成的定义

- `admin-server/Dockerfile` 存在，`docker build -f admin-server/Dockerfile admin-server` 能成功构建镜像。
- `admin-server/docker-compose.yml` 存在，`docker compose up` 能拉起 MySQL/Redis/app 三个健康的容器，app 能连上另外两个。
- `admin-server/.golangci.yml` 存在，`golangci-lint run` 能在 `internal/repository`、`internal/domain` 范围内跑通（允许有 lint 问题待修，但配置本身要能正常执行不报配置错误）。
- `.github/workflows/ci.yml` 存在，push 到分支后三个 job 都能触发并至少跑到"能执行"（不要求 Phase 1 一落地就全绿，允许后续几周逐步把 lint/test 问题修完）。
