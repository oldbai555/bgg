# 04 — IAM + Chat 领域改造（Phase 1 Week 2，最高风险/最高价值）

> 本文档是可直接执行的任务说明。执行者（Cursor / Claude Code）在改动前应完整阅读一遍，按顺序推进任务，每完成一个任务跑一次 `go build ./...`，全部任务完成后跑通"完成的定义"里的验证步骤。

## 0. 前置依赖

在动手之前必须已经落地（本文档假设这些已经存在，只按文件名引用，不复述内容）：

- [`01-architecture-target.md`](./01-architecture-target.md) —— 领域服务放置位置、`registry.Domain` 分层原则、跨域窄接口约定。
- [`02-transactions-and-uow.md`](./02-transactions-and-uow.md) —— `repository.Repository.Transact(ctx, fn)` 方法 + `registry.Transact(ctx, repo, fn)` 包级函数 + 各 Model 的 `WithSession`。
- [`03-wire-and-middleware.md`](./03-wire-and-middleware.md) —— 中间件构造函数收窄到 `config.Config`/`*repository.Repository`/`*registry.Domain`，Wire provider 化。

## 1. 现状核查（已读真实代码，不是推测）

### 1.1 `user_create_logic.go`（`internal/logic/iam/user/user_create_logic.go`）—— 本轮的旗舰 bug

- `UserCreateLogic.UserCreate` 直接 `iamrepo.NewUserRepository(...).Create(...)`，成功后调用 `initChatForNewUser`，失败只 `logx.Errorf`，不回滚——这个"失败不影响主流程"的语义本身是对的（产品期望），但承载它的写法有两个真实问题：
  1. `initChatForNewUser` 内部 **直接 `import chatrepo "postapocgame/admin-server/internal/repository/chat"`**——IAM 域越界直接持有 Chat 域仓储，这是需要在 Phase 2 拆分前修掉的跨域违规。
  2. 为存量用户批量建私聊那段用 `userRepo.FindPage(l.ctx, 1, 10000, "")` 一次拉全表（`internal/repository/iam/user_repository.go:122` 对应调用点），用户表一旦过万这里就是一次性全表扫描+ 内存驻留。`UserRepository` 接口已经声明了 `FindChunk(ctx, limit, lastId uint64) ([]iammodel.AdminUser, uint64, error)`（`internal/repository/iam/user_repository.go:17,95-97`），底层 `r.model.FindChunk` 已经实现，**但目前全仓库没有任何调用点用它**——`internal/logic/system/notice/notice_create_logic.go` 和 `notice_update_logic.go` 已经在用同样的"分批 + 异步 goroutine + 失败只记日志"模式（`FindChunk` + for 循环 + `lastID` 游标），可以直接照抄这个已验证过的模式，不用重新发明。
  3. 群组加入 + N 次私聊创建（每次私聊又是 `chat` 行 + 2 条 `chat_user` 行）中间没有任何事务保护，任何一步失败都可能留下"群已建、成员未加全"或"私聊已建、只有一方在会话里"的脏数据。

### 1.2 `PermissionResolver`（`internal/domain/iam/permission_resolver.go`）—— 现有领域服务，绕过了已装配的 `registry.Domain`

`PermissionResolver.CanAccess` 内部手动 `iamrepo.NewApiRepository(r.repo)` / `NewUserRoleRepository` / `NewRolePermissionRepository` / `NewPermissionApiRepository` 四次构造仓储（`permission_resolver.go:26,36,45,61`），完全没用上已经在 `registry.Domain.IAM` 里现成的 `Api`/`UserRole`/`RolePermission`/`PermissionApi` 字段。更严重的是调用方 `internal/middleware/permissionmiddleware.go` 自己 `iamdomain.NewPermissionResolver(m.svcCtx.Repository)` 现场构造，每次请求 new 一次——这条路径完全没有走 Wire 装配好的 `registry.Domain`，是"领域服务已经写了但没被正确接线"的典型案例。

### 1.3 `AuthMiddleware`（`internal/middleware/authmiddleware.go`）

现状：吃整个 `*svc.ServiceContext`，方法体里 `iamrepo.NewTokenBlacklistRepository(m.svcCtx.Repository)` 现场构造。属于 03 文档要收窄的 11 个中间件之一，本文档负责的是 IAM 侧改完后如何对接收窄后的构造函数（见任务 4）。

### 1.4 IAM 域内其余"多仓储调用"文件全量核查

对 `internal/logic/iam/**` 做了 `grep -o "[A-Za-z]*repo\.New[A-Za-z]*Repository("` 统计，出现 ≥2 个不同仓储构造调用的文件如下（按用途分类，**不是所有多仓储调用都需要领域服务**——只有满足"跨 ≥2 表写 / 跨域 / 非平凡业务规则（RBAC 授权、密码认证）"才需要）：

| 文件 | 仓储数 | 写操作 | 判断 |
|---|---|---|---|
| `user/user_create_logic.go` | 3 | 建用户 + 建群关系 + 建私聊（跨域写） | **需要领域服务**（本文档任务 1） |
| `role_permission/role_permission_update_logic.go` | 3 | `RolePermissionRepository.UpdateRolePermissions`（先物理删后批量插，单表两步无事务保护） | **需要领域服务**（RBAC 授权规则，任务 5） |
| `permission_menu/permission_menu_update_logic.go` | 3 | 同上模式，写 `admin_permission_menu` | **需要领域服务**（任务 5） |
| `user_role/user_role_update_logic.go` | 3 | `UserRoleRepository.UpdateUserRoles`（先删后循环插） | **需要领域服务**（任务 5） |
| `permission_api/permission_api_update_logic.go` | 3 | 写 `admin_permission_api` | **需要领域服务**（任务 5） |
| `auth/login_logic.go` | 4 | 登录日志 + 公告未读通知均为 `go func(){...}` 异步、失败只记日志 | **不需要 Transact**（已经是正确的异步尽力而为写法），但登录本身是"密码/认证"非平凡业务规则，建议后续下沉为薄 `AuthDomainService.Login`，优先级低于上面几项（见任务 6，可选） |
| `role_permission/role_permission_list_logic.go`、`permission_menu/permission_menu_list_logic.go`、`permission_menu/menu_tree_logic.go`、`menu/menu_my_tree_logic.go`、`permission_api/permission_api_list_logic.go`、`user_role/user_role_list_logic.go`、`auth/profile_logic.go` | 2-3 | 全部只读聚合（查树/查列表） | **保持薄 logic 直调**，不需要领域服务 |

结论：IAM 域除 `user_create_logic.go` 外，另有 4 个 RBAC 授权分配类文件（role_permission/permission_menu/user_role/permission_api 的 update）是真实的领域服务候选，这一点和计划文档的猜测一致，现已核实确认。

### 1.5 关键架构决策：领域服务用哪一层的 Transact，避免循环 import

`01`/`02` 的设计是"领域服务挂在 `registry.Domain` 各域子结构体上"（如 `IAMDomain.UserService`），这意味着 **`internal/repository/registry` 包要 import `internal/domain/iam`**。如果领域服务反过来调用 `registry.Transact(...)`，`internal/domain/iam` 就要 import `internal/repository/registry`——两个包互相 import，编译直接失败。

**本文档的领域服务一律只依赖 `internal/repository` 包的 `Repository.Transact(ctx, fn)` 方法（02 文档定义的底层原语），不 import `internal/repository/registry`。** 回调里按需 `iamrepo.NewXxxRepository(txRepo)` / `chatrepo.NewXxxRepository(txRepo)` 自己构造子仓储，和现有 `PermissionResolver`、`internal/domain/task` 的 import 方式完全一致（它们目前也只 import `internal/repository`，不 import `registry`）。`registry.Transact` 包级函数保留给**没有专属领域服务的 logic 文件**直接使用（拿到一个事务内的 `*registry.Domain`），两条路径互不冲突。

## 2. 目标目录结构（本文档新增/修改的文件）

```
internal/domain/
├── iam/
│   ├── permission_resolver.go     # 已存在，任务 3 改造
│   ├── user_service.go            # 新增，任务 1
│   ├── rbac_service.go            # 新增，任务 5
│   └── auth_service.go            # 新增（可选，任务 6）
└── chat/
    └── onboarding.go              # 新增，任务 2（Onboarding 接口 + ChatOnboardingService 实现）

internal/repository/registry/domain.go   # 任务 3/1/5：IAMDomain/ChatDomain 加字段，NewDomain 按依赖顺序构造
internal/middleware/authmiddleware.go       # 任务 4：改用 03 收窄后的构造函数
internal/middleware/permissionmiddleware.go # 任务 4：改为持有 *iamdomain.PermissionResolver（由 providePermissionMiddleware 从 registry.Domain.IAM.PermissionResolver 取出注入）
internal/logic/iam/user/user_create_logic.go              # 任务 1：改为薄委托
internal/logic/iam/role_permission/role_permission_update_logic.go  # 任务 5：改为薄委托
internal/logic/iam/permission_menu/permission_menu_update_logic.go  # 任务 5：改为薄委托
internal/logic/iam/user_role/user_role_update_logic.go              # 任务 5：改为薄委托
internal/logic/iam/permission_api/permission_api_update_logic.go    # 任务 5：改为薄委托
```

## 3. 非目标

- 不改变 `initChatForNewUser` 现有的产品语义（新用户加默认群 + 为存量用户建私聊，失败不影响建用户主流程）——只改实现方式，不改行为。
- 不在本轮给 Chat 域引入除 `Onboarding` 之外的领域服务（`chat_group_create_logic.go`、`chat_message_send_logic.go` 等 Chat 自身的多表写方法虽然也有事务缺口，但不在"IAM+Chat 联合改造"这个切口内，留给 Chat 单独排期，见第 6 节备注）。
- 不改 RBAC 权限判定算法本身（`CanAccess` 的业务逻辑不变，只改它怎么被构造和注入）。
- 不做权限缓存接入（`pkg/cache/business_cache.go` 的 `CacheKeyUserPermissions` 真正用起来是 Part B.5 的事，本轮不做）。

## 4. 任务清单

### 任务 1：`UserDomainService`（`internal/domain/iam/user_service.go`）

```go
package iam

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	"postapocgame/admin-server/internal/repository"
	iammodel "postapocgame/admin-server/internal/model/iam"
	iamrepo "postapocgame/admin-server/internal/repository/iam"
	"postapocgame/admin-server/pkg/errs"

	"github.com/zeromicro/go-zero/core/logx"
	chatdomain "postapocgame/admin-server/internal/domain/chat"
)

// UserDomainService 承载"建用户"这类跨越 IAM 自身表 + 需要触发 Chat 域初始化的编排逻辑。
type UserDomainService struct {
	repo       *repository.Repository
	onboarding chatdomain.Onboarding // 窄接口，不 import internal/repository/chat
}

func NewUserDomainService(repo *repository.Repository, onboarding chatdomain.Onboarding) *UserDomainService {
	return &UserDomainService{repo: repo, onboarding: onboarding}
}

type CreateUserInput struct {
	Username, Nickname, Password, Avatar, Signature string
	DepartmentId                                    uint64
	Status                                           int64
}

// CreateUser 建用户（用户名唯一性校验 + 密码加密 + 落库包在事务里，
// 为后续可能追加的"建默认角色关联"等同表内多步操作预留原子性；
// Chat 初始化异步尽力而为，失败不回滚用户创建，这是产品既定语义）。
func (s *UserDomainService) CreateUser(ctx context.Context, in CreateUserInput) (*iammodel.AdminUser, error) {
	if in.Username == "" || in.Password == "" {
		return nil, errs.New(errs.CodeBadRequest, "用户名和密码不能为空")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, "密码加密失败", err)
	}

	user := &iammodel.AdminUser{
		Username:     in.Username,
		Nickname:     in.Nickname,
		PasswordHash: string(hash),
		Avatar:       in.Avatar,
		Signature:    in.Signature,
		DepartmentId: in.DepartmentId,
		Status:       in.Status,
	}

	err = s.repo.Transact(ctx, func(ctx context.Context, txRepo *repository.Repository) error {
		userRepo := iamrepo.NewUserRepository(txRepo)
		if _, err := userRepo.FindByUsername(ctx, in.Username); err == nil {
			return errs.New(errs.CodeBadRequest, "用户名已存在")
		}
		if err := userRepo.Create(ctx, user); err != nil {
			return errs.Wrap(errs.CodeInternalError, "创建用户失败", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// Chat 初始化：异步、尽力而为，不阻塞建用户请求、失败不回滚。
	// ⚠️ 触发方式（进程内 goroutine 直调 vs 通过 internal/domain/task 调度器派发）见下方"待确认事项"。
	go func(newUserID uint64) {
		defer func() {
			if r := recover(); r != nil {
				logx.Errorf("Chat onboarding 发生 panic: userId=%d, err=%v", newUserID, r)
			}
		}()
		if err := s.onboarding.InitNewUser(context.Background(), newUserID); err != nil {
			logx.Errorf("初始化新用户聊天数据失败: userId=%d, err=%v", newUserID, err)
		}
	}(user.Id)

	return user, nil
}
```

**⚠️ 待确认事项（必须在实现前停下来问用户，不是技术细节，是产品/延迟取舍）**：

计划文档给出的方向是"Phase 1 先通过 `internal/domain/task` 现有调度器实现异步派发"，但现有调度器（`internal/domain/task/scheduler.go`）的执行单元是 `interfaces.TaskExecutor`，注册在 `internal/wire/providers.go:59-62` 的 `executorsMap[int]interfaces.TaskExecutor`（当前只有 `executorsMap[1] = ExcelExportExecutor`）；把 Chat onboarding 做成 `executorsMap[2]` 意味着：① 需要真的往 `admin_task` 表插一行任务记录（用户会在"任务中心"里看到"新用户初始化"这类任务，之前这个动作是完全隐式的，现在变成可见的）；② 从"请求处理完立刻触发"变成"最长等一个调度周期（默认 5 秒，`consts.TaskDefaultScanInterval`）才执行"。上面代码里给出的 `go func(){...}` 直接异步调用是最小改动、维持现有的"立即触发、纯后台、用户不可见"语义。**这两种方案哪个符合产品预期，需要用户在实现前拍板**；本文档先按"进程内 goroutine 直调"实现，如果用户选择调度器派发方案，把 `go func(){...}` 那段替换成 `taskRepo.Create(ctx, &taskmodel.AdminTask{Type: <新分配的 task_type 字典值>, ExecutionType: consts.TaskExecutionTypeAsync, Params: ...})`，并新增一个 `ChatOnboardingExecutor implements interfaces.TaskExecutor`，在 `provideTaskExecutors` 里注册。

### 任务 2：`chatdomain.Onboarding` 接口 + `ChatOnboardingService`（`internal/domain/chat/onboarding.go`）

**跨域方向对称，Chat 依赖 IAM 的用户数据同样不能直接 `import internal/repository/iam`**：`user_create_logic.go` 的原罪之一就是 IAM 反过来 import `internal/repository/chat`（1.1 节），如果修复方案在 Chat 侧又 import `internal/repository/iam`，只是把越界方向倒过来，同一个问题原地打转，而且 Phase 2 拆分后 `chat-rpc`/`iam-rpc` 是两个进程、两个数据库，`chatdomain` 包到时候里如果还留着 `iamrepo` 的 import，会直接编译不过——这不是风格问题，是 Phase 2 能否顺利抽取的硬约束。做法和 IAM→Chat 方向完全对称：`chatdomain` 包自己声明一个只依赖最小数据形状的窄接口，实现和适配都放在 `registry.NewDomain`（组合根，允许同时 import 两个域的包）里：

```go
package chatdomain

import (
	"context"
	"time"

	"postapocgame/admin-server/internal/consts"
	chatmodel "postapocgame/admin-server/internal/model/chat"
	"postapocgame/admin-server/internal/repository"
	chatrepo "postapocgame/admin-server/internal/repository/chat"
	"postapocgame/admin-server/pkg/errs"

	"github.com/zeromicro/go-zero/core/logx"
)

// Onboarding 供其它领域（当前只有 IAM）触发的窄接口：
// 只暴露"新用户上线要做什么"，不暴露 Chat 域仓储/模型细节，IAM 不需要 import internal/repository/chat。
type Onboarding interface {
	InitNewUser(ctx context.Context, newUserID uint64) error
}

// UserRef 是 Chat 域视角下对"一个用户"的最小引用，不依赖 iammodel.AdminUser，
// 这样 chatdomain 包不需要 import internal/model/iam，跟不 import internal/repository/iam 是同一个理由：
// Phase 2 拆分后 chat-rpc 进程里根本不会有 IAM 的 model/repository 代码可 import。
type UserRef struct {
	ID uint64
}

// UserLister 是 Chat 域向外部要求提供的窄接口：分批遍历"活跃、未删除"的用户引用，
// 用于批量建私聊。实现方（IAM 的 UserRepository.FindChunk）在 registry.NewDomain 里适配注入，
// chatdomain 包本身不知道实现方是谁。
type UserLister interface {
	FindChunk(ctx context.Context, limit int, lastID uint64) ([]UserRef, uint64, error)
}

// defaultGroupChatID 引用 internal/consts 的常量，不是本地魔法数（AGENTS.md 第 3 节：
// 系统级枚举/常量放 internal/consts，禁止业务代码硬编码）。internal/consts 里目前没有这个值，
// 属于本任务需要新增的常量，加在 internal/consts/chat.go（如果这个文件不存在就新建），
// 命名建议 consts.DefaultGroupChatID，不要在 chatdomain 包内部私有声明。
const defaultGroupChatID = consts.DefaultGroupChatID

type ChatOnboardingService struct {
	repo       *repository.Repository
	userLister UserLister
}

func NewChatOnboardingService(repo *repository.Repository, userLister UserLister) *ChatOnboardingService {
	return &ChatOnboardingService{repo: repo, userLister: userLister}
}

func (s *ChatOnboardingService) InitNewUser(ctx context.Context, newUserID uint64) error {
	s.joinDefaultGroup(ctx, newUserID)
	return s.createPrivateChatsForExistingUsers(ctx, newUserID)
}

// joinDefaultGroup 加入默认企业群组；失败只记日志，不影响后续私聊初始化。
func (s *ChatOnboardingService) joinDefaultGroup(ctx context.Context, newUserID uint64) {
	chatRepo := chatrepo.NewChatRepository(s.repo)
	chatUserRepo := chatrepo.NewChatUserRepository(s.repo)

	groupChat, err := chatRepo.FindByID(ctx, defaultGroupChatID)
	if err != nil || groupChat.DeletedAt != 0 {
		logx.Infof("默认企业群组不存在或已删除，跳过加入群组操作: userId=%d", newUserID)
		return
	}

	chatUsers, _ := chatUserRepo.FindUsersByChatID(ctx, defaultGroupChatID)
	for _, cu := range chatUsers {
		if cu.UserId == newUserID {
			return // 已在群组中
		}
	}

	now := time.Now().Unix()
	if err := chatUserRepo.Create(ctx, &chatmodel.ChatUser{
		ChatId: defaultGroupChatID, UserId: newUserID, JoinedAt: now, CreatedAt: now, UpdatedAt: now,
	}); err != nil {
		logx.Errorf("将新用户加入默认企业群组失败: userId=%d, err=%v", newUserID, err)
	}
}

// createPrivateChatsForExistingUsers 用 s.userLister.FindChunk 分批遍历存量用户，
// 每个用户的"建私聊 + 拉两人入会"包一层 s.repo.Transact，批内一个用户失败不影响其他用户，
// 整个方法失败也不影响 IAM 那边已经提交的用户创建（由调用方 UserDomainService 决定"失败只记日志"）。
// 注意：UserLister 实现（registry.NewDomain 里的适配器）负责只返回未删除、状态正常的用户，
// chatdomain 包本身不感知 IAM 的 deleted_at/status 字段语义。
func (s *ChatOnboardingService) createPrivateChatsForExistingUsers(ctx context.Context, newUserID uint64) error {
	chatRepo := chatrepo.NewChatRepository(s.repo)

	const limit = 100
	lastID := uint64(0)
	for {
		users, newLastID, err := s.userLister.FindChunk(ctx, limit, lastID)
		if err != nil {
			return errs.Wrap(errs.CodeInternalError, "分批查询用户失败", err)
		}
		if len(users) == 0 {
			break
		}

		for _, existingUser := range users {
			if existingUser.ID == newUserID {
				continue
			}
			if existing, err := chatRepo.FindPrivateChatByUserIDs(ctx, newUserID, existingUser.ID); err == nil && existing != nil {
				continue // 私聊已存在
			}
			if err := s.createPrivateChat(ctx, newUserID, existingUser.ID); err != nil {
				logx.Errorf("创建私聊失败: newUserId=%d, existingUserId=%d, err=%v", newUserID, existingUser.ID, err)
				continue
			}
		}

		if len(users) < limit {
			break
		}
		lastID = newLastID
	}
	return nil
}

// createPrivateChat 建一条私聊 + 两条 chat_user 关联，三条写在一个事务里，
// 修复原来"私聊建了、一方或两方没加进去"的部分失败风险。
func (s *ChatOnboardingService) createPrivateChat(ctx context.Context, userA, userB uint64) error {
	return s.repo.Transact(ctx, func(ctx context.Context, txRepo *repository.Repository) error {
		chatRepo := chatrepo.NewChatRepository(txRepo)
		chatUserRepo := chatrepo.NewChatUserRepository(txRepo)

		now := time.Now().Unix()
		privateChat := &chatmodel.Chat{Type: consts.ChatTypePrivate, CreatedBy: 0, CreatedAt: now, UpdatedAt: now} // consts.ChatTypePrivate 需要新增，见下方说明
		if err := chatRepo.Create(ctx, privateChat); err != nil {
			return err
		}
		if err := chatUserRepo.Create(ctx, &chatmodel.ChatUser{
			ChatId: privateChat.Id, UserId: userA, JoinedAt: now, CreatedAt: now, UpdatedAt: now,
		}); err != nil {
			return err
		}
		return chatUserRepo.Create(ctx, &chatmodel.ChatUser{
			ChatId: privateChat.Id, UserId: userB, JoinedAt: now, CreatedAt: now, UpdatedAt: now,
		})
	})
}
```

`FindUsersByChatID` 对应现有 `ChatRepository` 接口方法名（不是原代码里误写的 `chatUserRepo.FindByChatID`——原 `user_create_logic.go` 调的是 `chatUserRepo.FindByChatID`，但 `internal/repository/chat/chat_repository.go` 里这个方法实际挂在 `ChatRepository.FindUsersByChatID` 上；改造时以 `internal/repository/chat/*.go` 当前真实接口签名为准，如果 `ChatUserRepository` 另有同名方法就用那个，写代码前先 `grep -n "func (r \*chatUserRepository)"` 确认一遍，不要照抄本文档的方法名）。

### 任务 3：`PermissionResolver` 移入 `registry.NewDomain`

`internal/domain/iam/permission_resolver.go` 本身不用大改（构造参数仍是 `*repository.Repository`），改的是"谁来构造它"：

`internal/repository/registry/domain.go`：

```go
import (
	"context"

	iamdomain "postapocgame/admin-server/internal/domain/iam"
	chatdomain "postapocgame/admin-server/internal/domain/chat"
	iamrepo "postapocgame/admin-server/internal/repository/iam"
	"postapocgame/admin-server/internal/consts"
	// ...原有 import 不变
)

type IAMDomain struct {
	// ...原有 11 个仓储字段不变
	PermissionResolver *iamdomain.PermissionResolver
	UserService         *iamdomain.UserDomainService
	RBAC                 *iamdomain.RBACService // 任务 5
}

type ChatDomain struct {
	// ...原有 3 个仓储字段不变
	Onboarding chatdomain.Onboarding
}

// iamUserListerAdapter 实现 chatdomain.UserLister，是 registry 包（组合根，允许同时
// import iam 和 chat 两个域）承担的"翻译"职责：chatdomain 包本身永远不 import
// internal/repository/iam 或 internal/model/iam。只返回未删除、状态正常的用户，
// 过滤逻辑放在这里而不是 chatdomain 里，因为"状态正常"是 IAM 的业务语义。
type iamUserListerAdapter struct {
	userRepo iamrepo.UserRepository
}

func (a *iamUserListerAdapter) FindChunk(ctx context.Context, limit int, lastID uint64) ([]chatdomain.UserRef, uint64, error) {
	users, newLastID, err := a.userRepo.FindChunk(ctx, int64(limit), lastID)
	if err != nil {
		return nil, 0, err
	}
	refs := make([]chatdomain.UserRef, 0, len(users))
	for _, u := range users {
		// consts.UserStatusEnabled 目前 internal/consts 里没有，需要新增（AGENTS.md 第 3 节：
		// 系统级枚举/常量放 internal/consts，禁止硬编码）；deleted_at 沿用现有软删除判断惯例（!=0 即已删除）。
		if u.DeletedAt != 0 || u.Status != consts.UserStatusEnabled {
			continue
		}
		refs = append(refs, chatdomain.UserRef{ID: u.Id})
	}
	return refs, newLastID, nil
}

func NewDomain(repo *repository.Repository) *Domain {
	// ...
	userRepo := iamrepo.NewUserRepository(repo)
	chatOnboarding := chatdomain.NewChatOnboardingService(repo, &iamUserListerAdapter{userRepo: userRepo})
	return &Domain{
		// ...
		Chat: ChatDomain{
			Chat: chatrepo.NewChatRepository(repo),
			ChatUser: chatrepo.NewChatUserRepository(repo),
			ChatMessage: chatrepo.NewChatMessageRepository(repo),
			Onboarding: chatOnboarding,
		},
		IAM: IAMDomain{
			// ...原有 11 个字段
			PermissionResolver: iamdomain.NewPermissionResolver(repo),
			UserService:        iamdomain.NewUserDomainService(repo, chatOnboarding),
			RBAC:               iamdomain.NewRBACService(repo),
		},
		// ...
	}
}
```

**构造顺序要求**：`chatOnboarding` 必须先于 `IAM.UserService` 构造好（`UserService` 依赖它），这也是为什么 `NewDomain` 里把 `chatOnboarding` 提到局部变量而不是内联在两处字面量里各构造一次——内联会构造出两个不同的 `*ChatOnboardingService` 实例，IAM 手上那个和 `registry.Domain.Chat.Onboarding` 对不上，虽然功能上问题不大（都是无状态服务），但会造成理解成本，禁止这样写。

### 任务 4：`AuthMiddleware`/`PermissionMiddleware` 接入收窄后的构造函数

`AuthMiddleware` 按 03 文档的通用模式收窄为只依赖 `*repository.Repository`/`config.Config`。`PermissionMiddleware` 是 03 文档"全部中间件只依赖 Repository/config"这条通用规则的**唯一例外**：它需要的是已经构造好的 `*iamdomain.PermissionResolver`（任务 3 已经把它的构造挪进了 `registry.NewDomain`），而不是整个 `*registry.Domain`——传整个 `Domain` 会让一个只做权限判定的中间件拿到"读写任意域任意仓储"的能力，超出它实际需要的范围，也和其他中间件"依赖面收到刚好够用"的原则不一致。正确的收窄方式是：**Wire provider 函数**的入参是 `*registry.Domain`（因为这是 Wire 图里已有的节点），但**中间件构造函数**本身只吃从中取出的窄类型：

```go
// internal/middleware/permissionmiddleware.go
type PermissionMiddleware struct {
	resolver *iamdomain.PermissionResolver
}
func NewPermissionMiddleware(resolver *iamdomain.PermissionResolver) *PermissionMiddleware {
	return &PermissionMiddleware{resolver: resolver}
}
func (m *PermissionMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := jwthelper.FromContext(r.Context())
		if !ok {
			response.ErrorCtx(r.Context(), w, errs.New(errs.CodeUnauthorized, "未登录或登录已过期"))
			return
		}
		allowed, err := m.resolver.CanAccess(r.Context(), user.UserID, r.Method, r.URL.Path)
		// ...其余逻辑不变
	}
}

// internal/wire/providers.go —— provider 函数负责"从 Domain 里取出窄类型"这一步，
// 中间件本身不知道 registry.Domain 的存在。
func providePermissionMiddleware(domain *registry.Domain) *middleware.PermissionMiddleware {
	return middleware.NewPermissionMiddleware(domain.IAM.PermissionResolver)
}

// internal/middleware/authmiddleware.go —— 字段/构造函数签名以 03-wire-and-middleware.md §2.1 为准，
// 这里不重复定义一份不一致的版本，直接照抄：
type AuthMiddleware struct {
	repo      *repository.Repository
	jwtConfig config.JWTConf
}
func NewAuthMiddleware(cfg config.Config, repo *repository.Repository) *AuthMiddleware {
	return &AuthMiddleware{repo: repo, jwtConfig: cfg.JWT}
}
// Handle 内部 iamrepo.NewTokenBlacklistRepository(m.svcCtx.Repository) 改成 iamrepo.NewTokenBlacklistRepository(m.repo)，
// m.svcCtx.Config.JWT.AccessSecret 改成 m.jwtConfig.AccessSecret（03 文档已经写清楚，这里不重复）。
```

`providers.go` 里 `AuthMiddleware` 的 provider 直接传 `*repository.Repository`/`config.Config`；`PermissionMiddleware` 的 provider 传 `*registry.Domain`（仅用于取出 `.IAM.PermissionResolver`，不是把整个 Domain 交给中间件持有）。`svcCtx.AuthMiddleware`/`svcCtx.PermissionMiddleware` 两个扁平字段本身不变（03 文档已定的规则：中间件字段保持扁平，不引入嵌套 struct）。这一处是 03 文档"全部 11 个中间件只依赖 Repository/config"表格的唯一例外，03 文档需要同步补一条脚注说明。

### 任务 5：RBAC 授权分配类文件改造（`internal/domain/iam/rbac_service.go`）

覆盖 1.4 节列出的 4 个文件。以 `role_permission_update_logic.go` 为例，当前 `RolePermissionRepository.UpdateRolePermissions` 内部是"先物理删除该角色全部权限关联，再批量插入新的"两步 SQL、无事务保护（`internal/repository/iam/role_permission_repository.go:47-73`），删除成功后插入失败会导致该角色权限被清空且无法自动恢复。`user_role_update_logic.go` 对应的 `UpdateUserRoles` 也是同样的"先删后插"结构（`internal/repository/iam/user_role_repository.go:39-61`）。

**这里的"物理删除"不违反 `AGENTS.md` 的软删除规则，本任务不改这一点，只加事务**：`admin_role_permission`/`admin_user_role`/`admin_permission_menu`/`admin_permission_api` 是纯关联/中间表（role↔permission、user↔role 这类多对多映射），不是承载业务实体审计历史的表（不是 `admin_user`/`admin_role` 这类需要软删除保留恢复能力的表）——"整批替换某个角色当前生效的权限集合"这个操作本身的语义就是"重新计算当前状态"，不是"删除一条业务记录"，物理删除关联行不丢失任何需要审计/恢复的信息。这是仓库既有约定，`role_permission_repository.go` 本身就是 `AGENTS.md`/`.cursor/rules/10-go-code-style.mdc` 明确点名的 squirrel 参考实现，其"先删后插"写法不是本次要修的问题；本任务的范围严格限定在"给这两步 SQL 包一层事务，避免删除成功、插入失败导致权限被清空且无法恢复"，不改动删除策略本身，也不要顺手把关联表改成软删除。

```go
package iam

import (
	"context"

	"postapocgame/admin-server/internal/repository"
	iamrepo "postapocgame/admin-server/internal/repository/iam"
	"postapocgame/admin-server/pkg/errs"
)

// RBACService 承载角色/权限/菜单/接口的授权分配类写操作——
// 这些方法之所以需要领域服务，不是因为跨了很多表（大多数只写一张关联表），
// 而是因为它们是"非平凡业务规则"（RBAC 授权变更），且现有的"先删后插"两步写法
// 本身就需要事务保护，属于 01 文档判分标准里明确列出的场景。
type RBACService struct {
	repo *repository.Repository
}

func NewRBACService(repo *repository.Repository) *RBACService {
	return &RBACService{repo: repo}
}

// UpdateRolePermissions 校验角色/权限存在性 + 事务内先删后插。
func (s *RBACService) UpdateRolePermissions(ctx context.Context, roleID uint64, permissionIDs []uint64) error {
	return s.repo.Transact(ctx, func(ctx context.Context, txRepo *repository.Repository) error {
		roleRepo := iamrepo.NewRoleRepository(txRepo)
		if _, err := roleRepo.FindByID(ctx, roleID); err != nil {
			return errs.Wrap(errs.CodeBadRequest, "角色不存在", err)
		}
		permRepo := iamrepo.NewPermissionRepository(txRepo)
		permissions, err := permRepo.ListByIds(ctx, permissionIDs)
		if err != nil || len(permissions) != len(permissionIDs) {
			return errs.New(errs.CodeBadRequest, "权限不存在")
		}
		rpRepo := iamrepo.NewRolePermissionRepository(txRepo)
		if err := rpRepo.UpdateRolePermissions(ctx, roleID, permissionIDs); err != nil {
			return errs.Wrap(errs.CodeInternalError, "更新角色权限失败", err)
		}
		return nil
	})
	// 缓存清理（DeleteMenuTree）保持原来的 go func() 尽力而为写法，放在调用方 logic 里，
	// 不要塞进领域服务——缓存失效不是这个方法的核心职责，也不需要事务保护。
}

// UpdateUserRoles / UpdatePermissionMenus / UpdatePermissionApis 结构与上面完全一致，
// 分别对应 user_role_update_logic.go、permission_menu_update_logic.go、permission_api_update_logic.go
// 当前的校验+更新逻辑，照此模式实现，此处不重复贴代码。
```

对应的 4 个 logic 文件改为薄委托：

```go
func (l *RolePermissionUpdateLogic) RolePermissionUpdate(req *types.RolePermissionUpdateReq) error {
	if req.RoleId == 0 {
		return errs.New(errs.CodeBadRequest, "角色ID不能为空")
	}
	if err := l.svcCtx.Domain.IAM.RBAC.UpdateRolePermissions(l.ctx, req.RoleId, req.PermissionIds); err != nil {
		return err
	}
	cache := l.svcCtx.Repository.BusinessCache
	go func() {
		if err := cache.DeleteMenuTree(context.Background()); err != nil {
			l.Errorf("清除菜单树缓存失败: %v", err)
		}
	}()
	return nil
}
```

### 任务 6（可选，优先级低于 1-5）：`AuthDomainService.Login`

`login_logic.go` 目前的写操作（登录日志、公告未读通知）已经是正确的异步尽力而为写法，**不需要 Transact**。把 `Login` 方法体（用户名查询 → 状态校验 → bcrypt 密码比对 → token 生成）下沉成 `AuthDomainService.Login(ctx, username, password) (*TokenPair, error)` 的价值主要是可测试性（sqlmock 覆盖"密码错误"/"账号禁用"/"用户不存在"分支不需要绕开整个 HTTP logic 层）。如果时间不够，这一项可以推迟到 Phase 1 Week 4-5 的"测试覆盖补漏"阶段，不阻塞本文档其余任务的完成判定。

## 5. sqlmock 测试要求

按 [`02-transactions-and-uow.md`](./02-transactions-and-uow.md) 定的testing 约定，每个开 `Transact` 的方法都要有 happy-path + rollback-path 两个测试。至少覆盖：

- `UserDomainService.CreateUser`
  - happy path：`sqlmock.ExpectBegin()` → 用户名查询 `ExpectQuery` 返回空结果集 → `ExpectExec`（INSERT `admin_user`）成功 → `ExpectCommit()`；断言返回的 `user.Id` 等于 mock 里 `LastInsertId`。
  - rollback path：`ExpectBegin()` → 用户名查询返回一行（模拟用户名已存在）→ 直接 `ExpectRollback()`（不应该走到 INSERT）；断言返回 `errs.CodeBadRequest`。
  - 另一个 rollback 变体：用户名查询正常返回空 → INSERT 失败（`ExpectExec(...).WillReturnError(...)`）→ `ExpectRollback()`；断言返回 `errs.CodeInternalError`。
- `ChatOnboardingService.createPrivateChat`
  - happy path：`ExpectBegin()` → 3 个 `ExpectExec`（chat 插入、chat_user×2）均成功 → `ExpectCommit()`。
  - rollback path：第二个 `ExpectExec`（第一条 chat_user）失败 → `ExpectRollback()`，断言第三个 exec（第二条 chat_user）**不会**被调用（sqlmock 的 `ExpectationsWereMet()` 天然保证这一点，不用额外断言逻辑）。
- `RBACService.UpdateRolePermissions`
  - happy path：角色/权限校验查询均返回合法数据 → `ExpectBegin` → 先 DELETE 后批量 INSERT 两个 `ExpectExec` 均成功 → `ExpectCommit`。
  - rollback path：批量 INSERT 失败 → `ExpectRollback`，断言角色的旧权限关联"看起来"仍然是 DELETE 前的状态（用 mock 层面验证 DELETE 和 INSERT 都在同一个事务里，没有部分提交）。
- `ChatOnboardingService.createPrivateChatsForExistingUsers` 的分页边界：不需要 sqlmock，用 mockery 生成的 `chatdomain.UserLister` mock（按 `08-testing-strategy.md` §1.1 的约定，不手写 fake struct）配两页 `.On("FindChunk", ...)` 返回值（第一页刚好 `limit` 条触发第二次调用，第二页不足 `limit` 条触发循环退出），断言 mock 的 `FindChunk` 被调用两次且 `lastID` 参数正确传递——这是窄接口带来的直接好处：测这个方法不需要关心 IAM 那边的 SQL 长什么样。`iamUserListerAdapter`（`registry.NewDomain` 里的适配器）单独用 sqlmock 测一次"正确过滤了 deleted_at/status"即可，不需要和 `ChatOnboardingService` 的测试耦合在一起。

跳过测试：goctl 生成的透传方法（`FindByID`/`FindPage` 等）、`PermissionResolver.CanAccess` 里超级管理员短路分支之外的已有测试如果之前就有可以保留，本次不强求补全没有分支复杂度的纯查询。

## 6. 备注：Chat 域自身的事务缺口不在本文档范围

`internal/logic/chat/group/chat_group_create_logic.go` 的"建群 + 建群主关联 + 批量拉初始成员"、`chat_message_send_logic.go` 的多仓储写操作也有和 `user_create_logic.go` 类似的事务缺口，且同样有 `iamrepo` 交叉引用（方向相反：Chat 读 IAM 的用户信息）。**这个方向不是例外，一样要修**——跨域越界不因为方向是"读"就不算数，本文档任务 2 给 `ChatOnboardingService` 引入的 `UserLister` 窄接口 + `registry.NewDomain` 适配器模式，就是这两个文件后续改造时应该照抄的同一套解法，不要在这两个文件里重新发明。这两个文件本身的改造留给 Chat 域后续排期处理，不在"IAM+Chat 联合改造"这个 Phase 1 Week 2 切口内，执行时不要顺手改掉，避免这次改动范围失控——但如果顺手改了，必须按 `UserLister` 模式改，不能延续直接 `import iamrepo` 的旧写法。

## 7. 完成的定义

1. `go build ./...` 通过。
2. `internal/domain/iam/user_service.go`、`internal/domain/iam/rbac_service.go`、`internal/domain/chat/onboarding.go` 的 sqlmock 测试全部通过（`go test ./internal/domain/... -run TestUserDomainService -v`、`-run TestChatOnboarding`、`-run TestRBACService`）。
3. `grep -rn "internal/repository/chat" internal/logic/iam` 无匹配（越界 import 已清除）。
4. `grep -n "iamdomain.NewPermissionResolver" internal/middleware/permissionmiddleware.go` 无匹配（确认 `PermissionMiddleware` 不再每次请求现场构造 `PermissionResolver`，而是持有任务 4 里 `providePermissionMiddleware` 注入的 `m.resolver *iamdomain.PermissionResolver`）。
5. 人工冒烟测试（本地起服务）：
   - 用 Postman/curl 建一个新用户（`POST /api/v1/iam/user/create`），确认接口在几十毫秒内返回（不被 Chat 初始化阻塞），几秒后查 `chat_user` 表确认新用户已加入默认群组（`chat_id=1`）。
   - 用一个只有普通角色（无超级管理员）的账号登录，访问一个受权限保护的接口，确认 `PermissionMiddleware` 通过持有的 `m.resolver`（`providePermissionMiddleware` 从 `registry.Domain.IAM.PermissionResolver` 取出后注入）正常放行/拒绝（分别测试有权限和无权限两种情况）。
   - 在角色管理页面改一次角色的权限勾选并保存，确认保存成功且刷新后权限确实变了；再故意制造一次保存失败（例如传一个不存在的 `permissionId`），确认角色原有权限没有被清空（验证事务回滚生效）。
