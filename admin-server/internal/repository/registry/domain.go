package registry

import (
	"context"

	"postapocgame/admin-server/internal/consts"
	chatdomain "postapocgame/admin-server/internal/domain/chat"
	contentdomain "postapocgame/admin-server/internal/domain/content"
	iamdomain "postapocgame/admin-server/internal/domain/iam"
	sdkdomain "postapocgame/admin-server/internal/domain/sdk"
	"postapocgame/admin-server/internal/repository"
	blogrepo "postapocgame/admin-server/internal/repository/blog"
	chatrepo "postapocgame/admin-server/internal/repository/chat"
	iamrepo "postapocgame/admin-server/internal/repository/iam"
	miscrepo "postapocgame/admin-server/internal/repository/misc"
	monitoringrepo "postapocgame/admin-server/internal/repository/monitoring"
	sdkrepo "postapocgame/admin-server/internal/repository/sdk"
	systemrepo "postapocgame/admin-server/internal/repository/system"
	videorepo "postapocgame/admin-server/internal/repository/video"
)

// Domain 聚合各领域 Repository，启动时构造一次，Logic 通过 svcCtx.Domain 访问。
// 不再有 Task 字段：task 域已拆成独立的 task-rpc（services/task/），gateway 侧改成走
// svcCtx.TaskRPC 这个 zrpc client，不通过 Domain 聚合根访问。
type Domain struct {
	IAM        IAMDomain
	Blog       BlogDomain
	Chat       ChatDomain
	SDK        SDKDomain
	Video      VideoDomain
	Monitoring MonitoringDomain
	System     SystemDomain
	Misc       MiscDomain
}

type IAMDomain struct {
	User               iamrepo.UserRepository
	Role               iamrepo.RoleRepository
	Permission         iamrepo.PermissionRepository
	Menu               iamrepo.MenuRepository
	Department         iamrepo.DepartmentRepository
	UserRole           iamrepo.UserRoleRepository
	RolePermission     iamrepo.RolePermissionRepository
	Api                iamrepo.ApiRepository
	PermissionMenu     iamrepo.PermissionMenuRepository
	PermissionApi      iamrepo.PermissionApiRepository
	TokenBlacklist     iamrepo.TokenBlacklistRepository
	PermissionResolver *iamdomain.PermissionResolver
	UserService        *iamdomain.UserDomainService
	RBAC               *iamdomain.RBACService
}

type BlogDomain struct {
	Article        blogrepo.BlogArticleRepository
	ArticleTag     blogrepo.BlogArticleTagRepository
	ArticleAudit   blogrepo.BlogArticleAuditRepository
	FriendLink     blogrepo.BlogFriendLinkRepository
	SocialInfo     blogrepo.BlogSocialInfoRepository
	Tag            blogrepo.BlogTagRepository
	ArticleService *contentdomain.BlogArticleService
}

type ChatDomain struct {
	Chat        chatrepo.ChatRepository
	ChatUser    chatrepo.ChatUserRepository
	ChatMessage chatrepo.ChatMessageRepository
	Onboarding  chatdomain.Onboarding
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
		if u.DeletedAt != 0 || u.Status != consts.UserStatusEnabled {
			continue
		}
		refs = append(refs, chatdomain.UserRef{ID: u.Id})
	}
	return refs, newLastID, nil
}

type SDKDomain struct {
	Admin   *sdkrepo.SdkAdminRepository
	Public  *sdkrepo.SdkRepository
	Service *sdkdomain.SDKService
}

type VideoDomain struct {
	Video videorepo.VideoRepository
}

type MonitoringDomain struct {
	OperationLog   monitoringrepo.OperationLogRepository
	Metric         monitoringrepo.MetricRepository
	LoginLog       monitoringrepo.LoginLogRepository
	PerformanceLog monitoringrepo.PerformanceLogRepository
	AuditLog       monitoringrepo.AuditLogRepository
}

type SystemDomain struct {
	DictItem     systemrepo.DictItemRepository
	DictType     systemrepo.DictTypeRepository
	Config       systemrepo.ConfigRepository
	File         systemrepo.FileRepository
	Notification systemrepo.NotificationRepository
	Notice       systemrepo.NoticeRepository
}

type MiscDomain struct {
	Demo               miscrepo.DemoRepository
	DailyShortSentence miscrepo.DailyShortSentenceRepository
}

// NewDomain 从聚合 Repository 一次性构造全部领域 Repo。
func NewDomain(repo *repository.Repository) *Domain {
	if repo == nil {
		return nil
	}

	// chatOnboarding 必须先于 IAM.UserService 构造好（UserService 依赖它）。提到局部变量
	// 而不是在两处字面量里各构造一次，避免出现两个不同的 *ChatOnboardingService 实例
	// （IAM 手上那个和 Domain.Chat.Onboarding 对不上，造成理解成本）。
	userRepo := iamrepo.NewUserRepository(repo)
	chatOnboarding := chatdomain.NewChatOnboardingService(repo, &iamUserListerAdapter{userRepo: userRepo})

	return &Domain{
		IAM: IAMDomain{
			User:               userRepo,
			Role:               iamrepo.NewRoleRepository(repo),
			Permission:         iamrepo.NewPermissionRepository(repo),
			Menu:               iamrepo.NewMenuRepository(repo),
			Department:         iamrepo.NewDepartmentRepository(repo),
			UserRole:           iamrepo.NewUserRoleRepository(repo),
			RolePermission:     iamrepo.NewRolePermissionRepository(repo),
			Api:                iamrepo.NewApiRepository(repo),
			PermissionMenu:     iamrepo.NewPermissionMenuRepository(repo),
			PermissionApi:      iamrepo.NewPermissionApiRepository(repo),
			TokenBlacklist:     iamrepo.NewTokenBlacklistRepository(repo),
			PermissionResolver: iamdomain.NewPermissionResolver(repo),
			UserService:        iamdomain.NewUserDomainService(repo, chatOnboarding),
			RBAC:               iamdomain.NewRBACService(repo),
		},
		Blog: BlogDomain{
			Article:        blogrepo.NewBlogArticleRepository(repo),
			ArticleTag:     blogrepo.NewBlogArticleTagRepository(repo),
			ArticleAudit:   blogrepo.NewBlogArticleAuditRepository(repo),
			FriendLink:     blogrepo.NewBlogFriendLinkRepository(repo),
			SocialInfo:     blogrepo.NewBlogSocialInfoRepository(repo),
			Tag:            blogrepo.NewBlogTagRepository(repo),
			ArticleService: contentdomain.NewBlogArticleService(repo),
		},
		Chat: ChatDomain{
			Chat:        chatrepo.NewChatRepository(repo),
			ChatUser:    chatrepo.NewChatUserRepository(repo),
			ChatMessage: chatrepo.NewChatMessageRepository(repo),
			Onboarding:  chatOnboarding,
		},
		SDK: SDKDomain{
			Admin:   sdkrepo.NewSdkAdminRepository(repo),
			Public:  sdkrepo.NewSdkRepository(repo),
			Service: sdkdomain.NewSDKService(repo),
		},
		Video: VideoDomain{
			Video: videorepo.NewVideoRepository(repo),
		},
		Monitoring: MonitoringDomain{
			OperationLog:   monitoringrepo.NewOperationLogRepository(repo),
			Metric:         monitoringrepo.NewMetricRepository(repo),
			LoginLog:       monitoringrepo.NewLoginLogRepository(repo),
			PerformanceLog: monitoringrepo.NewPerformanceLogRepository(repo),
			AuditLog:       monitoringrepo.NewAuditLogRepository(repo),
		},
		System: SystemDomain{
			DictItem:     systemrepo.NewDictItemRepository(repo),
			DictType:     systemrepo.NewDictTypeRepository(repo),
			Config:       systemrepo.NewConfigRepository(repo),
			File:         systemrepo.NewFileRepository(repo),
			Notification: systemrepo.NewNotificationRepository(repo),
			Notice:       systemrepo.NewNoticeRepository(repo),
		},
		Misc: MiscDomain{
			Demo:               miscrepo.NewDemoRepository(repo),
			DailyShortSentence: miscrepo.NewDailyShortSentenceRepository(repo),
		},
	}
}

// Transact 在事务内执行 fn。fn 收到的 txDomain 是用换绑过事务 session 的 Repository
// 重新调用 NewDomain 构造出来的——每个 <domain>repo.NewXxxRepository(repo) 在构造时
// 就把 repo.DB / repo.XxxModel 捕获进内部字段，事务场景下必须重新构造一遍，不能只换
// Repository 本身而不重建 Domain。
func Transact(ctx context.Context, repo *repository.Repository, fn func(ctx context.Context, txDomain *Domain) error) error {
	return repo.Transact(ctx, func(ctx context.Context, txRepo *repository.Repository) error {
		return fn(ctx, NewDomain(txRepo))
	})
}
