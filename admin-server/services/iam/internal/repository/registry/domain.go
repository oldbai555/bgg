package registry

import (
	"postapocgame/admin-server/services/iam/internal/repository"
	iamrepo "postapocgame/admin-server/services/iam/internal/repository/iam"
	miscrepo "postapocgame/admin-server/services/iam/internal/repository/misc"
	monitoringrepo "postapocgame/admin-server/services/iam/internal/repository/monitoring"
	systemrepo "postapocgame/admin-server/services/iam/internal/repository/system"

	iamdomain "postapocgame/admin-server/services/iam/internal/domain/iam"
)

// Domain 聚合各领域 Repository，启动时构造一次，Logic 通过 svcCtx.Domain 访问。
// 不再有 Task/SDK/Chat/Blog/Video 字段：task 域已拆成独立的 task-rpc（services/task/），
// sdk 域已拆成独立的 sdk-rpc（services/sdk/），chat 域已拆成独立的 chat-rpc（services/chat/），
// blog+video 域已拆成独立的 content-rpc（services/content/），gateway 侧改成走
// svcCtx.TaskRPC/svcCtx.SdkRPC/svcCtx.ChatRPC/svcCtx.ContentRPC 这四个 zrpc client，不通过
// Domain 聚合根访问。
type Domain struct {
	IAM        IAMDomain
	Monitoring MonitoringDomain
	System     SystemDomain
	Misc       MiscDomain
}

type IAMDomain struct {
	User               iamrepo.UserRepository
	UserThirdParty     iamrepo.UserThirdPartyRepository
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
			UserThirdParty:     iamrepo.NewUserThirdPartyRepository(repo),
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
