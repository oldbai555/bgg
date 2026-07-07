package registry

import (
	"postapocgame/admin-server/internal/repository"
	blogrepo "postapocgame/admin-server/internal/repository/blog"
	chatrepo "postapocgame/admin-server/internal/repository/chat"
	iamrepo "postapocgame/admin-server/internal/repository/iam"
	miscrepo "postapocgame/admin-server/internal/repository/misc"
	monitoringrepo "postapocgame/admin-server/internal/repository/monitoring"
	sdkrepo "postapocgame/admin-server/internal/repository/sdk"
	systemrepo "postapocgame/admin-server/internal/repository/system"
	taskrepo "postapocgame/admin-server/internal/repository/task"
	videorepo "postapocgame/admin-server/internal/repository/video"
)

// Domain 聚合各领域 Repository，启动时构造一次，Logic 通过 svcCtx.Domain 访问。
type Domain struct {
	IAM        IAMDomain
	Blog       BlogDomain
	Chat       ChatDomain
	SDK        SDKDomain
	Video      VideoDomain
	Task       TaskDomain
	Monitoring MonitoringDomain
	System     SystemDomain
	Misc       MiscDomain
}

type IAMDomain struct {
	User           iamrepo.UserRepository
	Role           iamrepo.RoleRepository
	Permission     iamrepo.PermissionRepository
	Menu           iamrepo.MenuRepository
	Department     iamrepo.DepartmentRepository
	UserRole       iamrepo.UserRoleRepository
	RolePermission iamrepo.RolePermissionRepository
	Api            iamrepo.ApiRepository
	PermissionMenu iamrepo.PermissionMenuRepository
	PermissionApi  iamrepo.PermissionApiRepository
	TokenBlacklist iamrepo.TokenBlacklistRepository
}

type BlogDomain struct {
	Article      blogrepo.BlogArticleRepository
	ArticleTag   blogrepo.BlogArticleTagRepository
	ArticleAudit blogrepo.BlogArticleAuditRepository
	FriendLink   blogrepo.BlogFriendLinkRepository
	SocialInfo   blogrepo.BlogSocialInfoRepository
	Tag          blogrepo.BlogTagRepository
}

type ChatDomain struct {
	Chat        chatrepo.ChatRepository
	ChatUser    chatrepo.ChatUserRepository
	ChatMessage chatrepo.ChatMessageRepository
}

type SDKDomain struct {
	Admin  *sdkrepo.SdkAdminRepository
	Public *sdkrepo.SdkRepository
}

type VideoDomain struct {
	Video videorepo.VideoRepository
}

type TaskDomain struct {
	Task taskrepo.TaskRepository
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
			User:           iamrepo.NewUserRepository(repo),
			Role:           iamrepo.NewRoleRepository(repo),
			Permission:     iamrepo.NewPermissionRepository(repo),
			Menu:           iamrepo.NewMenuRepository(repo),
			Department:     iamrepo.NewDepartmentRepository(repo),
			UserRole:       iamrepo.NewUserRoleRepository(repo),
			RolePermission: iamrepo.NewRolePermissionRepository(repo),
			Api:            iamrepo.NewApiRepository(repo),
			PermissionMenu: iamrepo.NewPermissionMenuRepository(repo),
			PermissionApi:  iamrepo.NewPermissionApiRepository(repo),
			TokenBlacklist: iamrepo.NewTokenBlacklistRepository(repo),
		},
		Blog: BlogDomain{
			Article:      blogrepo.NewBlogArticleRepository(repo),
			ArticleTag:   blogrepo.NewBlogArticleTagRepository(repo),
			ArticleAudit: blogrepo.NewBlogArticleAuditRepository(repo),
			FriendLink:   blogrepo.NewBlogFriendLinkRepository(repo),
			SocialInfo:   blogrepo.NewBlogSocialInfoRepository(repo),
			Tag:          blogrepo.NewBlogTagRepository(repo),
		},
		Chat: ChatDomain{
			Chat:        chatrepo.NewChatRepository(repo),
			ChatUser:    chatrepo.NewChatUserRepository(repo),
			ChatMessage: chatrepo.NewChatMessageRepository(repo),
		},
		SDK: SDKDomain{
			Admin:  sdkrepo.NewSdkAdminRepository(repo),
			Public: sdkrepo.NewSdkRepository(repo),
		},
		Video: VideoDomain{
			Video: videorepo.NewVideoRepository(repo),
		},
		Task: TaskDomain{
			Task: taskrepo.NewTaskRepository(repo),
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
