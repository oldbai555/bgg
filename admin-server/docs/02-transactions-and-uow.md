# 事务与 `Transact`/`withSession` 落地清单

## 前置依赖

- 先读 `01-architecture-target.md` 的 A.1 节——本文档是它的可执行拆解，不重新解释设计动机。
- 代码库状态：`internal/repository/repository.go`、`internal/repository/registry/domain.go` 处于当前 HEAD（`36fbda9`）未改动状态。
- 本文档只改 `internal/model/**` 下**手改 sibling 文件**（非 `_gen.go`）+ `internal/repository/repository.go` + `internal/repository/registry/domain.go`。**不改** `internal/model/**/*_gen.go`——那是 goctl 生成产物，`AGENTS.md` 明令禁止手改，重新生成会覆盖任何手改内容。

---

## 1. `internal/repository/repository.go` 改动

在文件末尾（`BuildSources` 之后）新增：

```go
// Transact 在单个 MySQL 事务内执行 fn。
// fn 收到的 txRepo 是 r 的克隆：DB 字段与全部 *Model 字段都已经换绑到本次事务的 session 上。
// fn 内部必须只通过 txRepo 访问数据，不能继续闭包引用外层 r。
func (r *Repository) Transact(ctx context.Context, fn func(ctx context.Context, txRepo *Repository) error) error {
	return r.DB.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		return fn(ctx, r.withSession(session))
	})
}

// withSession 返回一个新的 *Repository，DB 与全部 37 个 *Model 字段都换绑到给定的事务 session。
func (r *Repository) withSession(session sqlx.Session) *Repository {
	return &Repository{
		DB:            sqlx.NewSqlConnFromSession(session),
		CacheConf:     r.CacheConf,
		Redis:         r.Redis,
		BusinessCache: r.BusinessCache,

		AdminUserModel:           r.AdminUserModel.WithSession(session),
		AdminRoleModel:           r.AdminRoleModel.WithSession(session),
		AdminPermissionModel:     r.AdminPermissionModel.WithSession(session),
		AdminMenuModel:           r.AdminMenuModel.WithSession(session),
		AdminDepartmentModel:     r.AdminDepartmentModel.WithSession(session),
		AdminUserRoleModel:       r.AdminUserRoleModel.WithSession(session),
		AdminRolePermissionModel: r.AdminRolePermissionModel.WithSession(session),
		AdminApiModel:            r.AdminApiModel.WithSession(session),
		AdminPermissionMenuModel: r.AdminPermissionMenuModel.WithSession(session),
		AdminPermissionApiModel:  r.AdminPermissionApiModel.WithSession(session),
		AdminConfigModel:         r.AdminConfigModel.WithSession(session),
		AdminDictTypeModel:       r.AdminDictTypeModel.WithSession(session),
		AdminDictItemModel:       r.AdminDictItemModel.WithSession(session),
		AdminFileModel:           r.AdminFileModel.WithSession(session),
		DemoModel:                r.DemoModel.WithSession(session),
		ChatModel:                r.ChatModel.WithSession(session),
		ChatUserModel:            r.ChatUserModel.WithSession(session),
		ChatMessageModel:         r.ChatMessageModel.WithSession(session),
		AdminOperationLogModel:   r.AdminOperationLogModel.WithSession(session),
		AdminLoginLogModel:       r.AdminLoginLogModel.WithSession(session),
		AuditLogModel:            r.AuditLogModel.WithSession(session),
		AdminPerformanceLogModel: r.AdminPerformanceLogModel.WithSession(session),
		AdminNoticeModel:         r.AdminNoticeModel.WithSession(session),
		AdminNotificationModel:   r.AdminNotificationModel.WithSession(session),
		DailyShortSentenceModel:  r.DailyShortSentenceModel.WithSession(session),
		VideoModel:               r.VideoModel.WithSession(session),
		SdkKeyModel:              r.SdkKeyModel.WithSession(session),
		SdkInterfaceModel:        r.SdkInterfaceModel.WithSession(session),
		SdkKeyApiModel:           r.SdkKeyApiModel.WithSession(session),
		SdkCallLogModel:          r.SdkCallLogModel.WithSession(session),
		AdminTaskModel:           r.AdminTaskModel.WithSession(session),
		BlogTagModel:             r.BlogTagModel.WithSession(session),
		BlogArticleModel:         r.BlogArticleModel.WithSession(session),
		BlogArticleTagModel:      r.BlogArticleTagModel.WithSession(session),
		BlogArticleAuditModel:    r.BlogArticleAuditModel.WithSession(session),
		BlogFriendLinkModel:      r.BlogFriendLinkModel.WithSession(session),
		BlogSocialInfoModel:      r.BlogSocialInfoModel.WithSession(session),
	}
}
```

这一步在全部 37 个 `.WithSession(session)` 方法存在之前无法编译通过——先做第 3 节的逐 Model 改造，再回来加这段代码，或者两者在同一个 PR/commit 里一次性做完都可以，顺序不影响正确性。

## 2. `internal/repository/registry/domain.go` 改动

在文件末尾新增：

```go
// Transact 在事务内执行 fn。fn 收到的 txDomain 是用换绑过事务 session 的 Repository
// 重新调用 NewDomain 构造出来的——每个 <domain>repo.NewXxxRepository(repo) 在构造时
// 就把 repo.DB / repo.XxxModel 捕获进内部字段（例如 iam.userRepository.conn = repo.DB），
// 事务场景下必须重新构造一遍，不能只换 Repository 本身而不重建 Domain。
func Transact(ctx context.Context, repo *repository.Repository, fn func(ctx context.Context, txDomain *Domain) error) error {
	return repo.Transact(ctx, func(ctx context.Context, txRepo *repository.Repository) error {
		return fn(ctx, NewDomain(txRepo))
	})
}
```

`context` 包需要新增 import。

## 3. 逐 Model 改造：工作示例（`AdminUserModel`）

以 `internal/model/iam/adminusermodel.go`（手改 sibling 文件，**不是** `adminusermodel_gen.go`）为例，改造前：

```go
package iam

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ AdminUserModel = (*customAdminUserModel)(nil)

type (
	AdminUserModel interface {
		adminUserModel
	}

	customAdminUserModel struct {
		*defaultAdminUserModel
	}
)

func NewAdminUserModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) AdminUserModel {
	return &customAdminUserModel{
		defaultAdminUserModel: newAdminUserModel(conn, c, opts...),
	}
}
```

改造后（新增 `WithSession` 到接口定义 + 实现）：

```go
package iam

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ AdminUserModel = (*customAdminUserModel)(nil)

type (
	AdminUserModel interface {
		adminUserModel
		// WithSession 返回一个绑定到事务 session 的新 AdminUserModel，供 Repository.withSession 调用。
		WithSession(session sqlx.Session) AdminUserModel
	}

	customAdminUserModel struct {
		*defaultAdminUserModel
	}
)

func NewAdminUserModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) AdminUserModel {
	return &customAdminUserModel{
		defaultAdminUserModel: newAdminUserModel(conn, c, opts...),
	}
}

// WithSession 见接口注释。table 字段直接复用，CachedConn 通过 sqlc.CachedConn.WithSession 换绑。
func (m *customAdminUserModel) WithSession(session sqlx.Session) AdminUserModel {
	return &customAdminUserModel{
		defaultAdminUserModel: &defaultAdminUserModel{
			CachedConn: m.CachedConn.WithSession(session),
			table:      m.table,
		},
	}
}
```

`defaultAdminUserModel`、`m.CachedConn`、`m.table` 都是 `adminusermodel_gen.go` 里定义的未导出标识符（`type defaultAdminUserModel struct { sqlc.CachedConn; table string }`），因为手改文件和生成文件在同一个包（`package iam`）内，可以直接引用，不需要额外导出。

**每个 Model 文件改动量约 10 行**：1 行接口方法声明 + 1 个 `WithSession` 方法实现（结构体名字和 `default` 前缀按各自模型替换）。改动模式完全一致，可以批量套用，但仍然要求逐文件替换 `defaultXxxModel`/`XxxModel` 里的具体类型名，不要用无脑字符串替换脚本导致误伤同名子串。

## 4. 全部需要改造的 Model 文件清单

以下 37 个文件通过 `find internal/model -name "*.go" ! -name "*_gen.go" ! -name "vars.go"` 实地核查得出，对应 `Repository` 结构体里的 37 个 `*Model` 字段（一一对应，顺序与 `repository.go` 源码一致）：

| # | `Repository` 字段 | 手改 sibling 文件（相对 `internal/model/`） |
|---|---|---|
| 1 | `AdminUserModel` | `iam/adminusermodel.go` |
| 2 | `AdminRoleModel` | `iam/adminrolemodel.go` |
| 3 | `AdminPermissionModel` | `iam/adminpermissionmodel.go` |
| 4 | `AdminMenuModel` | `iam/adminmenumodel.go` |
| 5 | `AdminDepartmentModel` | `iam/admindepartmentmodel.go` |
| 6 | `AdminUserRoleModel` | `iam/adminuserrolemodel.go` |
| 7 | `AdminRolePermissionModel` | `iam/adminrolepermissionmodel.go` |
| 8 | `AdminApiModel` | `iam/adminapimodel.go` |
| 9 | `AdminPermissionMenuModel` | `iam/adminpermissionmenumodel.go` |
| 10 | `AdminPermissionApiModel` | `iam/adminpermissionapimodel.go` |
| 11 | `AdminConfigModel` | `system/adminconfigmodel.go` |
| 12 | `AdminDictTypeModel` | `system/admindicttypemodel.go` |
| 13 | `AdminDictItemModel` | `system/admindictitemmodel.go` |
| 14 | `AdminFileModel` | `system/adminfilemodel.go` |
| 15 | `DemoModel` | `misc/demomodel.go` |
| 16 | `ChatModel` | `chat/chatmodel.go` |
| 17 | `ChatUserModel` | `chat/chatusermodel.go` |
| 18 | `ChatMessageModel` | `chat/chatmessagemodel.go` |
| 19 | `AdminOperationLogModel` | `monitoring/adminoperationlogmodel.go` |
| 20 | `AdminLoginLogModel` | `monitoring/adminloginlogmodel.go` |
| 21 | `AuditLogModel` | `monitoring/auditlogmodel.go` |
| 22 | `AdminPerformanceLogModel` | `monitoring/adminperformancelogmodel.go` |
| 23 | `AdminNoticeModel` | `system/adminnoticemodel.go` |
| 24 | `AdminNotificationModel` | `system/adminnotificationmodel.go` |
| 25 | `DailyShortSentenceModel` | `misc/dailyshortsentencemodel.go` |
| 26 | `VideoModel` | `video/videomodel.go` |
| 27 | `SdkKeyModel` | `sdk/sdkkeymodel.go` |
| 28 | `SdkInterfaceModel` | `sdk/sdkinterfacemodel.go` |
| 29 | `SdkKeyApiModel` | `sdk/sdkkeyapimodel.go` |
| 30 | `SdkCallLogModel` | `sdk/sdkcalllogmodel.go` |
| 31 | `AdminTaskModel` | `task/admintaskmodel.go` |
| 32 | `BlogTagModel` | `blog/blogtagmodel.go` |
| 33 | `BlogArticleModel` | `blog/blogarticlemodel.go` |
| 34 | `BlogArticleTagModel` | `blog/blogarticletagmodel.go` |
| 35 | `BlogArticleAuditModel` | `blog/blogarticleauditmodel.go` |
| 36 | `BlogFriendLinkModel` | `blog/blogfriendlinkmodel.go` |
| 37 | `BlogSocialInfoModel` | `blog/blogsocialinfomodel.go` |

**不在清单内、不要改动**：`internal/model/system/filemodel.go`（`FileModel`/`NewFileModel`）——实地核查确认它**没有**被 `Repository` 结构体引用（`Repository` 里注册的是 `AdminFileModel`，不是 `FileModel`），全仓库搜索 `system.NewFileModel`/`system.FileModel` 也没有调用点，是历史遗留的孤儿文件。本轮不处理它（不加 `WithSession`，也不删除）；是否清理记入 `11-descoped.md` 或作为独立的小 cleanup 项，不要在事务改造这个任务里顺带动它，避免把两件不相关的事混在一次 diff 里。

## 5. 执行顺序建议

不要求严格串行，但建议：

1. 先改 1-2 个 Model 文件（如 `AdminUserModel`），确认 `WithSession` 模式跑通、`go build ./internal/model/...` 通过。
2. 批量套用到剩余 35 个文件。
3. 加 `repository.go` 的 `Transact`/`withSession`，此时 37 个字段应该都能编译通过。
4. 加 `registry.go` 的 `Transact`。
5. 全仓库 `go build ./...`。

## 完成的定义

1. 37 个 Model 手改 sibling 文件全部新增 `WithSession(session sqlx.Session) XxxModel` 方法，且只改了手改文件，未触碰任何 `*_gen.go`。
2. `internal/repository/repository.go` 新增 `Transact`/`withSession`，`withSession` 覆盖全部 37 个字段（含 `DB`/`CacheConf`/`Redis`/`BusinessCache` 四个非 Model 字段的正确处理：`DB` 换绑、其余三个保持不变）。
3. `internal/repository/registry/domain.go` 新增 `Transact`，内部正确调用 `repo.Transact` + `NewDomain(txRepo)`。
4. `go build ./...` 全仓库编译通过。
5. 至少写一个 `sqlmock` 单测验证 `Repository.Transact` 的 happy-path（fn 返回 nil，数据确实提交）和 rollback-path（fn 返回 error，数据确实回滚）——这个测试本身作为 A.6/`08-testing-strategy.md` 范畴的一部分提前写出来，确认 `Transact` 机制本身是对的，不用等到 `04` 才第一次验证。
