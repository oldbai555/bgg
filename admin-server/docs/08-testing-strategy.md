# 08 · 测试策略（Phase 1 · A.6）

## 前置依赖

- `01-architecture-target.md`（事务方案/领域服务分层原则，含"两条使用路径"：领域服务内部只用 `Repository.Transact`，无专属领域服务的 logic 文件用 `registry.Transact`）、`02-transactions-and-uow.md`（两者具体实现）已落地，或至少签名已确定——本篇大量测试用例围绕这两个方法包裹的方法展开，没有签名测试没法写。**本篇 §2 的领域服务测试全部针对 `Repository.Transact`**（领域服务不 import `registry` 包，见 01），不要照抄成 `registry.Transact`。
- 已读 `internal/repository/repository.go`（`Repository` 聚合 37 个 `*Model` 字段 + `DB sqlx.SqlConn`）、`internal/repository/registry/domain.go`（`registry.Domain` 聚合）、`internal/repository/iam/role_permission_repository.go`（squirrel 查询范式）。

## 0. 现状

`admin-server` 当前 **零测试**（全仓库无 `_test.go` 文件），零 CI。目标不是把覆盖率从 0 拉到某个百分比，而是**覆盖真正可能出 bug、出了 bug 现在也没人能发现的地方**。全仓库唯一已知的真实事故现场就是 `internal/logic/iam/user/user_create_logic.go`：建用户 → 建群关系 → 建 N 个私聊，没有事务保护、没有测试、也没有人验证过它在部分失败时的行为——这正是本篇要优先补的洞。

**明确原则**：测试投入按"业务复杂度 + 出错代价"分配，不是雨露均沾。goctl 生成的透传代码、纯 CRUD 委托不测；跨表写、跨仓储读写、有非平凡规则的领域服务方法优先测。

## 1. 依赖确认

已核查 `admin-server/go.mod`、`go.sum`、`go list -m all`：

- `github.com/DATA-DOG/go-sqlmock v1.5.2` —— **已存在于 `go.sum` 且带完整 `h1:` 哈希**（当前是间接依赖，未出现在 `go.mod` 的 `require` 块里），`go list -m all` 能正常解析出该版本，说明模块已被下载过、可离线使用。
- `github.com/stretchr/testify v1.11.1` —— 同样已在 `go.sum` 里带 `h1:` 哈希，可用。
- **需要做的**：写第一个测试文件时执行 `go get github.com/DATA-DOG/go-sqlmock@v1.5.2 && go mod tidy`，把这两个包从"间接、仅因传递依赖出现在 go.sum"转成 `go.mod` 里的直接 `require`（`go mod tidy` 会自动处理，因为已有哈希，预期不需要联网重新下载）。这是本篇任务书里唯一需要执行的依赖变更命令，按 `10-dev-execution-and-review-points.md` 的口径，属于开发期可直接执行、事后 review 的操作，不用停下来问。
- go-zero 侧的关键 API：`sqlx.NewSqlConnFromDB(db *sql.DB, opts ...SqlOption) SqlConn`（`go-zero@v1.9.3/core/stores/sqlx/sqlconn.go:123`）——sqlmock 打桩的标准入口，把 `sqlmock.New()` 返回的 `*sql.DB` 包装成 `sqlx.SqlConn` 注入到 `repository.Repository{DB: ...}` 或直接构造 `iam.NewAdminXModel(conn, cacheConf)`。
- `github.com/vektra/mockery/v2` —— **仓库里目前完全没有引入**（`grep -rn "mockery" . --include=*.go --include=go.mod --include=go.sum` 无匹配），是本篇新增的依赖，用途见下方"接口边界 mock 用 mockery，不手写"。安装：`go install github.com/vektra/mockery/v2@latest`（本地/CI 都需要，属于开发工具链而不是项目依赖，不进 `go.mod`）。

### 1.1 两种 mock 各自的职责边界，不要混用

`sqlmock` 和 `mockery` 测的是两类完全不同的东西，本篇同时用两者，职责严格分开：

- **`sqlmock` 只用于验证"SQL/事务语义对不对"**：手写的 squirrel 查询拼没拼对、`Repository.Transact` 的 happy-path 是否 `Commit`、rollback-path 是否 `Rollback`（§2、§3）。这类测试关心"打到数据库的 SQL 长什么样"，必须走 sqlmock。
- **`mockery` 只用于隔离"跨接口边界的协作对象"**：领域服务依赖的窄接口（`chatdomain.Onboarding`、`chatdomain.UserLister`，见 `04-domain-iam-chat.md` 任务 1/2）、以及领域服务测试里不是这次测试重点的旁路 repository 依赖，用 mockery 生成的 mock 断言"调用没调用、传参对不对"，不关心 SQL。**禁止在 `*_test.go` 文件里手写 mock struct 代替 mockery**——手写 fake 没有生成器保证接口变更时同步报错，接口加一个方法，手写 fake 要靠人肉记得去改，mockery 生成物会在 `go build` 时因为漏实现方法直接报错。**禁止在 `*_test.go` 里声明 interface**——接口是生产代码的边界，属于 `internal/domain/<domain>/*.go`，不属于测试文件。
- 落地方式：在定义接口的生产代码文件顶部加 `//go:generate mockery --name=Onboarding --output=../../mocks/chatdomain --outpkg=chatdomain_mocks`（`chatdomain.Onboarding`、`chatdomain.UserLister` 各一条），生成物统一放 `internal/mocks/<domain>/`，执行 `go generate ./internal/domain/...` 批量生成，生成物本身按"goctl 生成代码"的同一条规则处理——不手改，重新生成会覆盖。
- 依赖多的领域服务测试用 fixture 构造函数收口（如 `newUserDomainServiceFixture(t, opts...)`），不要在每个 `Test*` 里重复写四五个参数的构造函数调用——构造函数签名一旦变，只用改一处 fixture，不用改几十个测试。

## 2. 领域服务测试（主战场）

`internal/domain/**` 下的领域服务是测试投入的核心，理由：这是唯一会出现"跨表写 + 需要回滚"这种真实故障模式的地方。

**强制要求**：每一个被 `Repository.Transact` 包裹的领域服务方法，必须同时有：

1. **happy-path 测试**：所有子操作都成功，断言最终状态正确、`mock.ExpectCommit()` 被调用。
2. **rollback-path 测试**：让其中一个子操作（不是第一个，选中间或最后一个更能验证"部分执行后回滚"的真实场景）返回 error，断言：
   - 方法整体返回 error；
   - `mock.ExpectRollback()` 被触发，其余子操作的 SQL 副作用不应该被提交；
   - 用 `mock.ExpectationsWereMet()` 收尾，确认没有遗漏或多余的 SQL 调用。

这一条要写清楚：**这是全仓库当前 100% 未被验证的行为**。`user_create_logic.go` 的 `initChatForNewUser` 现在连事务都没有，谈不上回滚测试；一旦按 04 的方案改造成 `UserDomainService.CreateUser` 并包上 `Repository.Transact`（用户名唯一性校验 + 建用户这部分，不含异步 chat 初始化——那部分是尽力而为、不回滚的），必须补上这一对测试才算完成，不是"加了事务就算做完"。

**测试范围**：`internal/domain/iam/`（含 `permission_resolver.go` 的 `CanAccess`——虽然它当前不涉及事务，但有真实的多分支权限判定逻辑，值得测）、`internal/domain/task/`（`scheduler.go` 的锁获取/释放、执行结果记录）、Phase 1 Week 2-3 陆续新增的约 35-40 个领域服务方法（凡是开了 `Repository.Transact` 的，无一例外都要有这两个测试）。没有专属领域服务、直接在 logic 文件里用 `registry.Transact` 的方法同样要有这一对测试，测试写法一致，区别只是构造事务期望时目标是 `*registry.Domain` 而不是某个具体 repository。

## 3. Repository 层测试（手写查询）

- **只测手写的 squirrel 查询方法**（`internal/repository/<domain>/*.go` 里非透传的方法，如 `role_permission_repository.go` 的 `ListPermissionIDsByRoleID`/`UpdateRolePermissions`）。
- 用 `sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))`，断言生成的 SQL 字符串（用正则容忍空格/反引号差异，不要求逐字节相等）和参数值，而不是只断言返回值——squirrel 拼错 WHERE 条件、拼错表名这类 bug，只测返回值的 mock 测试完全测不出来（因为 mock 本身就是按 query 匹配返回预设结果，SQL 拼错了但只要能匹配上任意查询依然会"测试通过"）。
- **明确跳过**：goctl 生成的 `*_gen.go` 里的 CRUD 透传方法（`FindOne`/`Insert`/`Update`/`Delete`/`FindPage`/`FindChunk`）——这些是生成代码，逻辑在 go-zero/goctl 模板里已经验证过，测它们是在测框架而不是测本项目的业务逻辑。

## 4. Logic 层测试（只测有真实分支的地方）

不是每个 logic 文件都要测，纯委托的 CRUD logic（`l.svcCtx.Domain.X.Y(...)` 一行转发）不单独测——领域服务测试已经覆盖了底层行为，logic 层再测一遍是重复劳动。

需要补测的 logic（有真实分支/状态机的）：
- 登录（`internal/logic/iam/auth/*login*.go`）：密码错误、用户不存在、用户被禁用、登录成功四个分支。
- 权限校验相关（依赖 `PermissionResolver.CanAccess`）：可以在领域服务层测过之后，logic 层只加一个"权限不足时返回正确错误码"的煙测。
- Refresh Token：token 过期、token 在黑名单、正常刷新三个分支。
- Task 调度触发路径（如果 logic 层有非平凡的状态转换）。

## 5. 集成测试套件（小而精）

- 构建标签 `//go:build integration`，独立于单元测试运行，跑真实 MySQL + Redis（本地或 CI 里的 docker-compose 服务）。
- 目标数量：**约 5-6 个**，不追多。已知必须覆盖的场景：
  1. 登录 e2e（真实建用户 → 真实登录 → 拿到 token）。
  2. 一个 RBAC 允许的请求（有权限的用户访问受保护接口，返回 200）。
  3. 一个 RBAC 拒绝的请求（无权限用户访问受保护接口，返回 403）。
  4. IAM 用户创建 → chat onboarding 全链路（验证 04 文档任务 1 修复后的异步/尽力而为语义：建用户请求本身不因 chat 初始化失败而失败，且 chat 数据最终正确写入）。
  5. Task 调度器跑一个完整周期（提交任务 → 调度器拾取 → 执行 → 状态更新为完成）。
  6.（可选，视 Phase 1 进度）一个领域服务方法的 `Repository.Transact` 回滚在真实 MySQL 上确实生效（unit test 用 sqlmock 验证了"调用了 Rollback"，集成测试验证"数据真的没落库"，两者互补，不重复）。

## 6. 明确的非目标（不做，写清楚为什么）

- **不做 HTTP handler 层的端到端测试套件**：handler 是 goctl 生成的薄胶水层（解析请求 → 调 logic → 序列化响应），业务分支全在 logic/domain 层，已经被 2/3/4 节覆盖；handler 本身随 `.api` 重新生成就要重写，维护一套跟着生成代码变的 E2E 测试是负收益。
- **不做 WebSocket 聊天的并发/压力测试**：`internal/hub/chathub.go` 的并发模型验证成本高、收益主要体现在生产真实负载下，用手写并发测试模拟意义有限，且不是当前阶段的风险重灾区（当前风险重灾区是无事务的多表写，不是 WS 并发）。
- **不设置全局 CI 覆盖率百分比门槛**。原因需要写清楚：覆盖率门槛是**结果指标**，会激励"测什么最容易拉高数字"而不是"测什么真的有风险"——最快拉高覆盖率的方式是给 goctl 生成的透传方法、简单 getter/setter 写测试，而这些恰恰是本文档明确排除的低价值目标。一旦设了百分比门槛，第 3、4 节里"明确跳过"的部分就会被倒逼着补测试，只为了凑数字，挤占本该花在第 2 节（真正有 bug 风险的领域服务）上的时间。CI 只检查"该测的测了"（见 09 篇 unit-test job：跑 `go test ./internal/domain/... ./internal/repository/...` 且不能有编译错误/测试失败），不检查"测了百分之多少"。

## 7. 领域服务测试骨架示例

以下是一个 `Repository.Transact` 包裹方法的测试骨架，对应 `04-domain-iam-chat.md` 已经落地的 `UserDomainService.CreateUser` 真实签名（领域服务内部只用 `Repository.Transact`，不 import `registry`，见 01）：

```go
type CreateUserInput struct {
    Username, Nickname, Password, Avatar, Signature string
    DepartmentId                                    uint64
    Status                                           int64
}

func (s *UserDomainService) CreateUser(ctx context.Context, in CreateUserInput) (*iammodel.AdminUser, error) {
    // ... 密码加密等前置校验 ...
    user := &iammodel.AdminUser{Username: in.Username /* ... */}
    err := s.repo.Transact(ctx, func(ctx context.Context, txRepo *repository.Repository) error {
        userRepo := iamrepo.NewUserRepository(txRepo)
        if _, err := userRepo.FindByUsername(ctx, in.Username); err == nil {
            return errs.New(errs.CodeBadRequest, "用户名已存在")
        }
        return userRepo.Create(ctx, user)
    })
    if err != nil {
        return nil, err
    }
    // Chat 初始化是事务外的异步尽力而为操作，不在本方法的事务测试范围内，
    // 单独用 mockery 生成的 chatdomain.Onboarding mock 验证"是否被调用"即可，不需要 sqlmock（见 §1.1）。
    go func(newUserID uint64) { /* ... */ }(user.Id)
    return user, nil
}
```

测试文件 `internal/domain/iam/user_service_test.go`：

```go
package iam_test

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	tmock "github.com/stretchr/testify/mock" // 别名 tmock：本文件里 sqlmock.Sqlmock 变量名就叫 mock，两者会撞名
	"github.com/stretchr/testify/require"

	iamdomain "postapocgame/admin-server/internal/domain/iam"
	"postapocgame/admin-server/internal/repository"
	chatdomainmocks "postapocgame/admin-server/internal/mocks/chatdomain" // mockery 生成，见 §1.1
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// newTestRepo 用 sqlmock 构造一个可注入事务期望的 *repository.Repository。
// 具体字段按 02-transactions-and-uow.md 确认的 Repository/Model 构造方式补全。
func newTestRepo(t *testing.T) (*repository.Repository, sqlmock.Sqlmock, func()) {
	t.Helper()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	require.NoError(t, err)

	conn := sqlx.NewSqlConnFromDB(db)
	repo, err := repository.NewRepository(conn, cache.CacheConf{}, nil /* 或注入一个可用的 miniredis/伪 Redis 客户端 */)
	require.NoError(t, err)

	return repo, mock, func() { _ = db.Close() }
}

func TestUserDomainService_CreateUser_HappyPath(t *testing.T) {
	repo, mock, cleanup := newTestRepo(t)
	defer cleanup()

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT")). // FindByUsername 唯一性校验，无匹配行
							WillReturnError(sql.ErrNoRows)
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `admin_user`")).
		WillReturnResult(sqlmock.NewResult(42, 1))
	mock.ExpectCommit()

	// chatdomainmocks.NewMockOnboarding(t) 是 mockery v2 生成的构造函数，自动注册
	// t.Cleanup 做 AssertExpectations，不需要手写 fake struct（见 §1.1 的强制要求）。
	onboarding := chatdomainmocks.NewMockOnboarding(t)
	done := make(chan struct{})
	onboarding.On("InitNewUser", tmock.Anything, uint64(42)).
		Run(func(tmock.Arguments) { close(done) }).
		Return(nil)

	svc := iamdomain.NewUserDomainService(repo, onboarding)
	user, err := svc.CreateUser(context.Background(), iamdomain.CreateUserInput{
		Username: "alice", Password: "s3cr3t",
	})

	require.NoError(t, err)
	assert.Equal(t, uint64(42), user.Id)
	assert.NoError(t, mock.ExpectationsWereMet())

	// CreateUser 内部是 go func(){ onboarding.InitNewUser(...) }(user.Id) 异步调用（04 文档已定），
	// 事务本身的 happy-path 断言在上面已经完成；这里额外等一小段时间断言异步调用确实发生，
	// 超时说明 CreateUser 没有触发 onboarding，是真实 bug，不是 flaky。
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("CreateUser 没有触发 onboarding.InitNewUser")
	}
}

func TestUserDomainService_CreateUser_RollbackOnDuplicateUsername(t *testing.T) {
	repo, mock, cleanup := newTestRepo(t)
	defer cleanup()

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT")). // FindByUsername 命中已存在用户
							WillReturnRows(sqlmock.NewRows([]string{"id", "username"}).AddRow(1, "bob"))
	mock.ExpectRollback()

	// 用户名冲突分支在 Transact 内部就 return 了，走不到 onboarding 调用，
	// mockery mock 默认零调用即满足期望，不需要 .On(...)。
	onboarding := chatdomainmocks.NewMockOnboarding(t)

	svc := iamdomain.NewUserDomainService(repo, onboarding)
	_, err := svc.CreateUser(context.Background(), iamdomain.CreateUserInput{
		Username: "bob", Password: "s3cr3t",
	})

	require.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
```

`newTestRepo` 里 `redis.Redis` 参数目前留空/待定——具体怎么在 sqlmock 场景下满足 `NewRepository` 对非空 Redis 客户端的要求（当前 `NewRepository` 在 `rdb == nil` 时会返回 error），由落地 `02-transactions-and-uow.md` 时一并确定（可选：引入一个轻量内存 Redis 模拟，或者把 `NewRepository` 的 Redis 校验放宽为可选依赖），此处不预先决定,写测试时按实际签名调整。上面两个测试用真实存在的"用户名唯一性校验"分支演示 happy-path/rollback-path，而不是虚构一个 04 里不存在的"建角色关联"步骤——测试骨架必须跟着领域服务的真实实现走，不能自己发明额外的 SQL 语句。

## 完成的定义

- `go test ./internal/domain/... ./internal/repository/...` 全绿。
- 每个 Phase 1 新增、被 `Repository.Transact`（领域服务）或 `registry.Transact`（无专属领域服务的 logic 文件）包裹的方法，happy-path + rollback-path 测试都存在且通过——这是本篇最硬性的验收标准，其余测试范围允许按时间预算酌情取舍。
- 集成测试套件（`//go:build integration`）能在本地/CI 跑通，不要求默认随 `go test ./...` 触发（避免开发者本地没有 MySQL/Redis 时跑不过）。
- `go.mod` 里 `github.com/DATA-DOG/go-sqlmock`、`github.com/stretchr/testify` 已从 indirect 转为 direct require。
