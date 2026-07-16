# 06 — Blog / Video / SDK 三域改造（Phase 1 Week 3）

> 本文档是可直接执行的任务说明。三个域放在一份文档里是因为在 IAM 那一轮把模式跑通之后（见 04），这里剩下的是结构相似的机械工作：找出真正的多表写方法、套用同一个"领域服务 + `repo.Transact`"模板、其余保持薄 logic 直调。执行者按域顺序（Blog → SDK → Video）推进，每个域改完跑一次 `go build ./...`。

## 0. 前置依赖

- [`01-architecture-target.md`](./01-architecture-target.md) —— 分层判断标准（跨 ≥2 表写才需要领域服务）、`internal/domain/<domain>` 按 Part B 服务边界分组的约定（Blog/Video 归入 `content` 分组，SDK 独立分组）。
- [`02-transactions-and-uow.md`](./02-transactions-and-uow.md) —— `repository.Repository.Transact(ctx, fn)`。
- [`03-wire-and-middleware.md`](./03-wire-and-middleware.md)。

## 1. 包边界（按 01 的分组约定）

```
internal/domain/
├── content/                    # 对应未来 content-rpc，Blog + Video 共用一个包
│   └── blog_service.go         # package content，只有 Blog 有真正的领域服务，Video 没有
└── sdk/
    └── sdk_service.go          # package sdk，对应未来 sdk-rpc
```

不新建 `internal/domain/video/`——审计结论（见第 3 节）是 Video 域没有任何方法满足"跨 ≥2 表写"的门槛，不需要领域服务这一层。

## 2. Blog 域

### 2.1 现状核查（已读真实代码）

`internal/logic/blog/**` 全部 32 个文件里，用 `grep -o "[A-Za-z]*repo\.New[A-Za-z]*Repository("` 统计出现 ≥2 个不同仓储构造的文件如下：

| 文件 | 仓储数 | 实际操作 | 判断 |
|---|---|---|---|
| `article/blog_article_list_logic.go`、`article/blog_article_detail_logic.go`、`public/public_blog_article_list_logic.go`、`public/public_blog_article_detail_logic.go` | 2 | 查文章 + 查标签，只读聚合 | 保持薄 logic |
| `article_audit/blog_article_audit_logic.go`、`article_audit/blog_article_audit_unpublish_logic.go` | 2 | 写审核记录（`blog_article_audit`）+ 更新文章状态（`blog_article`），两张表两次独立写，无事务保护（`blog_article_audit_logic.go:70,82`） | **需要领域服务** |

`blog_article_create_logic.go`/`blog_article_update_logic.go` 本身只调了 1 个仓储（`blogrepo.NewBlogArticleRepository`），grep 按"仓储构造次数"算不会命中它们，但要害不在 logic 层，**而在 Repository 层内部**：

- `internal/repository/blog/blog_article_repository.go` 的 `CreateWithTags`（第 187-238 行）：先 `INSERT blog_article`，拿到自增 ID 后再循环 `INSERT blog_article_tag`，两张表跨多条 SQL，**代码注释自己写着"使用 squirrel 手动插入，避免依赖事务 session API"**（第 196 行）——这是明确承认过、故意绕开事务的写法。失败处理是手动补偿删除文章行（`_ = r.articleModel.Delete(ctx, article.Id)`，第 228、232 行），属于"没有事务能力时的权宜之计"。
- `UpdateWithTags`（第 240-273 行）：更新文章 + 软删旧标签关联 + 插入新标签关联，三步跨表写，**完全没有任何失败补偿**，比 `CreateWithTags` 更脆弱。

这两个方法是 Blog 域**真正需要事务保护的地方**，是计划文档"blog 文章+标签创建"猜测的具体落点，现已核实确认（且发现 Update 比 Create 更需要修）。

`BlogArticleRepository` 接口（`internal/repository/blog/blog_article_repository.go:17-30`）目前没有一个不touch 标签的纯 `Update(ctx, article)` 方法——只有 `UpdateWithTags`（连带处理标签关联，语义比"只更新审核状态"重）。2.2 节的 `AuditArticle` 只需要更新 `blog_article` 一张表的两个字段，不应该套用 `UpdateWithTags`（会带上不相关的标签处理），也不应该绕过 Repository 直接碰 `BlogArticleModel`（违反"Repository 承载 SQL、Domain/Logic 不直连 Model"的分层）。所以本任务顺带给接口加一个方法：

```go
// internal/repository/blog/blog_article_repository.go 接口定义处新增：
Update(ctx context.Context, article *blogmodel.BlogArticle) error

// 实现（薄封装，直接透传给 Model，不涉及标签，不需要 squirrel）：
func (r *blogArticleRepository) Update(ctx context.Context, article *blogmodel.BlogArticle) error {
	return r.articleModel.Update(ctx, article)
}
```

### 2.2 领域服务：`internal/domain/content/blog_service.go`

关键设计点：`blogArticleRepository`（`internal/repository/blog/blog_article_repository.go:32-44`）构造时只是把 `repo.BlogArticleModel`/`repo.BlogArticleTagModel`/`repo.DB` 拆出来存到自己的字段里，`CreateWithTags`/`UpdateWithTags` 内部所有 SQL 都走同一个 `r.conn`。这意味着**只要把这个 Repository 构造在一个已经绑定了事务 session 的 `*repository.Repository` 之上，`CreateWithTags`/`UpdateWithTags` 内部原有的多条 SQL 自动就在同一个事务里**——不需要改 `blog_article_repository.go` 一行代码，只需要在领域服务里把调用包一层 `Transact`：

```go
package content

import (
	"context"

	"postapocgame/admin-server/internal/consts"
	blogmodel "postapocgame/admin-server/internal/model/blog"
	"postapocgame/admin-server/internal/repository"
	blogrepo "postapocgame/admin-server/internal/repository/blog"
	"postapocgame/admin-server/pkg/errs"
)

// BlogArticleService 承载文章创建/更新（跨 blog_article + blog_article_tag 两表写）
// 和审核（跨 blog_article_audit + blog_article 两表写）。
type BlogArticleService struct {
	repo *repository.Repository
}

func NewBlogArticleService(repo *repository.Repository) *BlogArticleService {
	return &BlogArticleService{repo: repo}
}

func (s *BlogArticleService) CreateArticle(ctx context.Context, article *blogmodel.BlogArticle, tagIDs []uint64) error {
	return s.repo.Transact(ctx, func(ctx context.Context, txRepo *repository.Repository) error {
		return blogrepo.NewBlogArticleRepository(txRepo).CreateWithTags(ctx, article, tagIDs)
	})
}

func (s *BlogArticleService) UpdateArticle(ctx context.Context, article *blogmodel.BlogArticle, tagIDs []uint64) error {
	return s.repo.Transact(ctx, func(ctx context.Context, txRepo *repository.Repository) error {
		return blogrepo.NewBlogArticleRepository(txRepo).UpdateWithTags(ctx, article, tagIDs)
	})
}

// AuditArticle 对应 blog_article_audit_logic.go 现有逻辑：写审核记录 + 更新文章审核状态，
// 原来是两次独立写，这里包进同一个事务。
func (s *BlogArticleService) AuditArticle(ctx context.Context, articleID uint64, result int64, remark string, auditorID uint64, auditorName string) (*blogmodel.BlogArticle, error) {
	var article *blogmodel.BlogArticle
	err := s.repo.Transact(ctx, func(ctx context.Context, txRepo *repository.Repository) error {
		articleRepo := blogrepo.NewBlogArticleRepository(txRepo)
		a, err := articleRepo.FindByID(ctx, articleID)
		if err != nil {
			return errs.Wrap(errs.CodeBadDB, "查询文章失败", err)
		}
		if a == nil || a.DeletedAt != 0 {
			return errs.New(errs.CodeNotFound, "文章不存在")
		}
		if a.AuditStatus != consts.BlogArticleAuditStatusPending { // internal/consts/blog.go:15，已存在，直接用
			return errs.New(errs.CodeForbidden, "当前状态不允许审核")
		}

		if err := blogrepo.NewBlogArticleAuditRepository(txRepo).Create(ctx, &blogmodel.BlogArticleAudit{
			ArticleId: a.Id, AuditStatus: result, AuditRemark: remark,
			AuditorId: auditorID, AuditorName: auditorName,
		}); err != nil {
			return err
		}

		a.AuditStatus = result
		if result == consts.BlogArticleAuditStatusPassed { // internal/consts/blog.go:16
			a.Status = consts.BlogArticleStatusAuditPassed // internal/consts/blog.go:7
		}
		if err := articleRepo.Update(ctx, a); err != nil {
			return errs.Wrap(errs.CodeBadDB, "更新文章审核状态失败", err)
		}
		article = a
		return nil
	})
	return article, err
}
```

上面代码已经直接用 `internal/consts/blog.go` 里真实存在的 `BlogArticleAuditStatusPending`（值 2）/`BlogArticleAuditStatusPassed`（值 3）/`BlogArticleStatusAuditPassed`（值 3）常量，不是占位符，照抄即可，不需要再去查值替换。

### 2.3 logic 文件改为薄委托

`blog_article_create_logic.go`：

```go
if err = l.svcCtx.Domain.Blog.ArticleService.CreateArticle(l.ctx, article, req.TagIds); err != nil {
	return nil, err
}
```

`blog_article_audit_logic.go` 里的 `audit.RecordAuditLog(...)` 调用保持在 logic 层原地不动（审计日志记录本身不属于这次事务保护的范围，也不应该失败就回滚审核结果）。

`CreateWithTags`/`UpdateWithTags` 内部原有的手动补偿删除（`_ = r.articleModel.Delete(ctx, article.Id)`）在真正走了事务之后是死代码（回滚会自动撤销未提交的 INSERT，不需要再手动删）——可以顺手删掉，但不是阻塞项，留着也无害（事务回滚后这行代码根本不会被执行到，因为方法返回 error 后整个函数栈都在事务边界内被回滚）。

### 2.4 保持薄 logic 直调的部分

`tag/*`、`friend_link/*`、`social_info/*` 全部是单表 CRUD，不动。`public/*` 全部只读，不动，但有一个需要单独说明的例外：

`public/public_blog_author_info_logic.go` 直接 `iamrepo.NewUserRepository(l.svcCtx.Repository).FindByID(l.ctx, 1)`——**跨域读取 IAM 域仓储**。这是计划文档 Part B.2 明确列出的"已定位到具体文件"之一，但它是**只读**、单次调用、没有事务/一致性风险，修复方式是 Phase 2 拆 `content-rpc`/`iam-rpc` 时改成调 `iam-rpc.GetUserProfile`。**本文档不处理这一处**——它不满足"跨 ≥2 表写"的领域服务门槛，强行现在就引入一个窄接口只是为了消除一次跨域 import，收益低于 04 文档里 IAM→Chat 那处（那处是写路径、且已经造成了真实的 bug）。原样保留，留个注释标记方便 Phase 2 时搜索：

```go
// TODO(phase2-content-rpc): 跨域读取 IAM 用户信息，Phase 2 拆分后改为调用 iam-rpc.GetUserProfile
```

## 3. Video 域

### 3.1 现状核查

`internal/logic/video/**` 8 个文件（`m3u8/m3u8_proxy_logic.go`、`public/*` 2 个、`video/*` 4 个、`video_collect/*` 2 个）逐个 grep 多仓储调用，**结果为空**——没有任何文件出现 2 个及以上不同仓储构造调用。`video_create_logic.go` 是典型代表：校验参数 → 构造 `video.Video` → `videorepo.NewVideoRepository(l.svcCtx.Repository).Create(...)`，单表单仓储。`m3u8_proxy_logic.go` 是纯 HTTP 反向代理（拉取 m3u8 清单/媒体分片转发），完全不碰数据库。

**确认计划文档的猜测：Video 域目前没有需要事务保护的写路径，是纯 CRUD + 网络代理。** 不需要 `internal/domain/content` 里额外加 Video 相关方法，`video/*` 全部保持 `svcCtx.Domain.Video.Video` 直调。

### 3.2 唯一要做的事

无代码改动。如果 04/05/06 三份文档执行完之后要给 Video 域补 sqlmock 测试，按 [`08-testing-strategy.md`](./08-testing-strategy.md) 的"手写 squirrel 查询的 repository 方法用 sqlmock 断言 SQL 和参数"补 `internal/repository/video/video_repository.go` 的测试即可，不需要新建领域服务测试。

## 4. SDK 域

### 4.1 现状核查

`internal/logic/sdk/**` 12 个文件里，`sdk_interface_update_logic.go`/`sdk_interface_create_logic.go` 各出现 2 次仓储构造（`sdkrepo.NewSdkAdminRepository` + `sdkrepo.NewSdkRepository`），**但两者不是同一次写操作跨了两张表**——`SdkRepository.BuildInterfaceCode(method, path)` 只是一个纯函数式的 apiCode 生成方法（不查库），真正的写只有 `repo.CreateInterface`/`repo.UpdateInterface` 一次，单表 `sdk_interface`。**计划文档"SDK 的 Admin/Public 仓储拆分需要事务保护"这个猜测不成立**，Admin/Public 只是"面向管理端的完整 CRUD 仓储"和"面向公开调用鉴权的轻量只读仓储"两个不同职责的仓储类，不是同一次写操作被拆成两半。

真正的多表写问题在 `internal/repository/sdk/sdk_admin_repository.go` 的 `SaveBindings`（第 231-263 行）：先软删除该 Key 名下全部旧的 `sdk_key_api` 绑定关系，再循环插入新的绑定，**"先删后插"两步无事务保护**，和 04 文档里 IAM 的 `UpdateRolePermissions`/`UpdateUserRoles` 是完全相同的模式（同一张关联表内的 replace-all 操作）。调用方是 `sdk_api_key_bind_save_logic.go`（`internal/logic/sdk/sdk/sdk_api_key_bind_save_logic.go:51`）。

另外，`SdkAdminRepository` 的构造方式（`internal/repository/sdk/sdk_admin_repository.go:16-21`）比 Blog 的仓储更直接——它本身就持有完整的 `repo *repository.Repository` 字段（不是拆出来的窄字段），意味着同样"把 Repository 构造在事务绑定过的 `*repository.Repository` 之上即可让内部所有 SQL 自动共享事务"这条规律照样适用。

`sdk_call_log_export_logic.go` 有跨域 import（`internal/repository/task`），用于创建导出任务——单表写（`admin_task`），且是"发起一个任务"这种性质的调用，不是需要原子性保护的组合写，计划文档 Part B.2 已经把它列为 Phase 2 拆分前需要处理的跨域 import，本文档不处理（原因同 2.4 节的 `public_blog_author_info_logic.go`）。

### 4.2 领域服务：`internal/domain/sdk/sdk_service.go`

```go
package sdk

import (
	"context"

	"postapocgame/admin-server/internal/repository"
	sdkrepo "postapocgame/admin-server/internal/repository/sdk"
	sdkmodel "postapocgame/admin-server/internal/model/sdk"
)

type SDKService struct {
	repo *repository.Repository
}

func NewSDKService(repo *repository.Repository) *SDKService {
	return &SDKService{repo: repo}
}

// SaveApiKeyBindings 把"软删旧绑定 + 插入新绑定"包进事务，
// SaveBindings 方法本身不用改，只需要把它构造在事务绑定过的 Repository 之上。
func (s *SDKService) SaveApiKeyBindings(ctx context.Context, sdkKeyID uint64, bindings []sdkmodel.SdkKeyApi) error {
	return s.repo.Transact(ctx, func(ctx context.Context, txRepo *repository.Repository) error {
		return sdkrepo.NewSdkAdminRepository(txRepo).SaveBindings(ctx, sdkKeyID, bindings)
	})
}
```

`sdk_api_key_bind_save_logic.go` 改为：

```go
if err := l.svcCtx.Domain.SDK.Service.SaveApiKeyBindings(l.ctx, req.SdkKeyId, bindings); err != nil {
	return errs.Wrap(errs.CodeInternalError, "保存授权失败", err)
}
```

### 4.3 保持薄 logic 直调的部分

`sdk_api_key_create/delete/list/update`、`sdk_interface_create/delete/list/update`、`sdk_call_log_list`、`sdk_api_key_bind_list`、`public/sdk_file_upload_logic.go` 全部单表写或纯读，不动。`sdk_call_log_export_logic.go` 按 4.1 节说明原样保留（加 TODO 注释标记 Phase 2 待办，格式同 2.4 节）。

## 5. `registry.Domain` 接线

**决定：`registry.Domain` 顶层保留现有的 `Domain.Blog`/`Domain.Video` 两个独立字段，不新增 `ContentDomain` 聚合体。** 01 文档"包边界按 Phase 2 服务分组组织"这条要求，管的是 `internal/domain/content/` 这个 **Go 包目录**（`blog_service.go`/`video_service.go` 放在同一个目录、同一个 `package content` 下，方便 Phase 2 整体搬走），不要求 `registry.Domain` 的字段结构也在 Phase 1 就跟着合并——Phase 1 阶段 blog/video 的表和 repository 仍然是分开的两套（`BlogDomain`/`VideoDomain` 两个独立子结构体，各自的 repository 字段不变），跟着改字段结构只是徒增这一轮的改动面，不改变任何实际行为，等 Phase 2 真正拆 `content-rpc`、数据库也物理合并时再考虑要不要把 `Domain.Blog`/`Domain.Video` 合并成一个字段（那时候是给 `content-rpc` 自己的组合根设计，不是 admin-server 单体的 `registry.Domain`）。所以新增服务字段直接挂在各自原有的子结构体上：

```go
// internal/repository/registry/domain.go
import (
	contentdomain "postapocgame/admin-server/internal/domain/content"
	sdkdomain "postapocgame/admin-server/internal/domain/sdk"
)

type BlogDomain struct {
	// ...原有仓储字段不变（Tag/Article/ArticleTag/ArticleAudit/FriendLink/SocialInfo）
	ArticleService *contentdomain.BlogArticleService // 新增，来自 internal/domain/content 包
}

// VideoDomain 不新增字段——第 3 节已确认 Video 域本轮没有需要领域服务的方法。

type SDKDomain struct {
	Admin  *sdkrepo.SdkAdminRepository // 原有字段不变
	Public *sdkrepo.SdkRepository      // 原有字段不变
	Service *sdkdomain.SDKService      // 新增
}
```

`internal/domain/content/blog_service.go` 里的 `package content` 和 `registry.Domain.Blog.ArticleService` 这个字段路径不矛盾——包名描述"这段代码逻辑上属于哪个未来服务"，字段路径描述"现在这个单体里数据实际怎么组织"，两者本轮解耦，不强求一致。

## 6. 非目标

- 不改 Video 域任何代码（第 3 节已确认不需要）。
- 不修 `public_blog_author_info_logic.go`/`sdk_call_log_export_logic.go` 的跨域 import，只加 TODO 注释标记，留给 Phase 2。
- 不改 `SdkRepository`（public 只读仓储）/`SdkAdminRepository` 的接口签名，只新增领域服务包一层。
- 不给 Blog/SDK 域的纯读接口（列表、详情、下拉选项）引入领域服务。

## 7. 完成的定义

1. `go build ./...` 通过。
2. `go test ./internal/domain/content/... ./internal/domain/sdk/... -v` 通过，覆盖：
   - `BlogArticleService.CreateArticle`：happy path（`ExpectBegin` → INSERT article 成功 → INSERT tag ×N 成功 → `ExpectCommit`）+ rollback path（第二条 tag INSERT 失败 → `ExpectRollback`，断言 article 那条 INSERT 不会被"手动补偿删除"逻辑触碰到，因为整个事务都没提交）。
   - `BlogArticleService.AuditArticle`：happy path + "状态不是待审核"分支（不触发任何写）+ 审核记录写入成功但文章状态更新失败的 rollback path。
   - `SDKService.SaveApiKeyBindings`：happy path（软删 + 批量插入成功）+ rollback path（批量插入中途失败，断言旧绑定"看起来"没有被清空，即 DELETE 和 INSERT 在同一事务内一起回滚）。
3. 人工冒烟测试：
   - 后台创建一篇带 2 个标签的文章，确认文章和标签关联都成功写入；再故意传一个不存在的 `tagId` 触发插入失败（如果标签表有外键式的业务校验的话，没有就跳过这一步，改成直接在测试环境临时改代码模拟插入失败），确认失败时文章行本身也没有残留。
   - 提交一篇文章审核（通过/驳回各测一次），确认 `blog_article_audit` 记录和 `blog_article.audit_status` 同时更新。
   - 给一个 SDK Key 保存接口绑定关系，确认旧绑定被替换为新绑定；查一下 `sdk_key_api` 表确认没有残留软删除失败的脏数据。
   - 创建/更新一条视频记录，确认功能行为和改造前一致（因为本次未改动 Video 域代码，这一步主要是回归验证没有被误改）。
