package registry

import (
	"context"

	"postapocgame/admin-server/internal/repository"
	blogrepo "postapocgame/admin-server/internal/repository/blog"
	iamrepo "postapocgame/admin-server/internal/repository/iam"
	miscrepo "postapocgame/admin-server/internal/repository/misc"
	monitoringrepo "postapocgame/admin-server/internal/repository/monitoring"
	systemrepo "postapocgame/admin-server/internal/repository/system"
	videorepo "postapocgame/admin-server/internal/repository/video"

	contentdomain "postapocgame/admin-server/internal/domain/content"
	iamdomain "postapocgame/admin-server/internal/domain/iam"
)

// Domain 聚合各领域 Repository，启动时构造一次，Logic 通过 svcCtx.Domain 访问。
// 不再有 Task/SDK/Chat 字段：task 域已拆成独立的 task-rpc（services/task/），sdk 域已拆成
// 独立的 sdk-rpc（services/sdk/），chat 域已拆成独立的 chat-rpc（services/chat/），gateway
// 侧改成走 svcCtx.TaskRPC/svcCtx.SdkRPC/svcCtx.ChatRPC 这三个 zrpc client，不通过 Domain
// 聚合根访问。
type Domain struct {
	IAM        IAMDomain
	Blog       BlogDomain
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

	return &Domain{
		IAM: IAMDomain{
			User:               iamrepo.NewUserRepository(repo),
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
			UserService:        iamdomain.NewUserDomainService(repo, repo.Redis),
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
