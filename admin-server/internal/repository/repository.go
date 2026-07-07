package repository

import (
	"postapocgame/admin-server/internal/config"
	businesscache "postapocgame/admin-server/pkg/cache"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"postapocgame/admin-server/internal/model/blog"
	"postapocgame/admin-server/internal/model/chat"
	"postapocgame/admin-server/internal/model/iam"
	"postapocgame/admin-server/internal/model/misc"
	"postapocgame/admin-server/internal/model/monitoring"
	"postapocgame/admin-server/internal/model/sdk"
	"postapocgame/admin-server/internal/model/system"
	"postapocgame/admin-server/internal/model/task"
	"postapocgame/admin-server/internal/model/video"
)

// Repository 聚合 goctl 生成的 Model，统一数据访问入口。
type Repository struct {
	DB            sqlx.SqlConn
	CacheConf     cache.CacheConf
	Redis         *redis.Redis                 // go-zero stores/redis 客户端
	BusinessCache *businesscache.BusinessCache // 业务层缓存工具

	AdminUserModel           iam.AdminUserModel
	AdminRoleModel           iam.AdminRoleModel
	AdminPermissionModel     iam.AdminPermissionModel
	AdminMenuModel           iam.AdminMenuModel
	AdminDepartmentModel     iam.AdminDepartmentModel
	AdminUserRoleModel       iam.AdminUserRoleModel
	AdminRolePermissionModel iam.AdminRolePermissionModel
	AdminApiModel            iam.AdminApiModel
	AdminPermissionMenuModel iam.AdminPermissionMenuModel
	AdminPermissionApiModel  iam.AdminPermissionApiModel
	AdminConfigModel         system.AdminConfigModel
	AdminDictTypeModel       system.AdminDictTypeModel
	AdminDictItemModel       system.AdminDictItemModel
	AdminFileModel           system.AdminFileModel
	DemoModel                misc.DemoModel
	ChatModel                chat.ChatModel
	ChatUserModel            chat.ChatUserModel
	ChatMessageModel         chat.ChatMessageModel
	AdminOperationLogModel   monitoring.AdminOperationLogModel
	AdminLoginLogModel       monitoring.AdminLoginLogModel
	AuditLogModel            monitoring.AuditLogModel
	AdminPerformanceLogModel monitoring.AdminPerformanceLogModel
	AdminNoticeModel         system.AdminNoticeModel
	AdminNotificationModel   system.AdminNotificationModel
	DailyShortSentenceModel  misc.DailyShortSentenceModel
	VideoModel               video.VideoModel
	SdkKeyModel              sdk.SdkKeyModel
	SdkInterfaceModel        sdk.SdkInterfaceModel
	SdkKeyApiModel           sdk.SdkKeyApiModel
	SdkCallLogModel          sdk.SdkCallLogModel
	AdminTaskModel           task.AdminTaskModel
	BlogTagModel             blog.BlogTagModel
	BlogArticleModel         blog.BlogArticleModel
	BlogArticleTagModel      blog.BlogArticleTagModel
	BlogArticleAuditModel    blog.BlogArticleAuditModel
	BlogFriendLinkModel      blog.BlogFriendLinkModel
	BlogSocialInfoModel      blog.BlogSocialInfoModel
}

func NewRepository(conn sqlx.SqlConn, cacheConf cache.CacheConf, rdb *redis.Redis) (*Repository, error) {
	if conn == nil {
		return nil, errors.New("repository requires sqlx conn")
	}
	if rdb == nil {
		return nil, errors.New("repository requires redis")
	}
	return &Repository{
		DB:                       conn,
		CacheConf:                cacheConf,
		Redis:                    rdb,
		BusinessCache:            businesscache.NewBusinessCache(rdb),
		AdminUserModel:           iam.NewAdminUserModel(conn, cacheConf),
		AdminRoleModel:           iam.NewAdminRoleModel(conn, cacheConf),
		AdminPermissionModel:     iam.NewAdminPermissionModel(conn, cacheConf),
		AdminMenuModel:           iam.NewAdminMenuModel(conn, cacheConf),
		AdminDepartmentModel:     iam.NewAdminDepartmentModel(conn, cacheConf),
		AdminUserRoleModel:       iam.NewAdminUserRoleModel(conn, cacheConf),
		AdminRolePermissionModel: iam.NewAdminRolePermissionModel(conn, cacheConf),
		AdminApiModel:            iam.NewAdminApiModel(conn, cacheConf),
		AdminPermissionMenuModel: iam.NewAdminPermissionMenuModel(conn, cacheConf),
		AdminPermissionApiModel:  iam.NewAdminPermissionApiModel(conn, cacheConf),
		AdminConfigModel:         system.NewAdminConfigModel(conn, cacheConf),
		AdminDictTypeModel:       system.NewAdminDictTypeModel(conn, cacheConf),
		AdminDictItemModel:       system.NewAdminDictItemModel(conn, cacheConf),
		AdminFileModel:           system.NewAdminFileModel(conn, cacheConf),
		DemoModel:                misc.NewDemoModel(conn, cacheConf),
		ChatModel:                chat.NewChatModel(conn, cacheConf),
		ChatUserModel:            chat.NewChatUserModel(conn, cacheConf),
		ChatMessageModel:         chat.NewChatMessageModel(conn, cacheConf),
		AdminOperationLogModel:   monitoring.NewAdminOperationLogModel(conn, cacheConf),
		AdminLoginLogModel:       monitoring.NewAdminLoginLogModel(conn, cacheConf),
		AuditLogModel:            monitoring.NewAuditLogModel(conn, cacheConf),
		AdminPerformanceLogModel: monitoring.NewAdminPerformanceLogModel(conn, cacheConf),
		AdminNoticeModel:         system.NewAdminNoticeModel(conn, cacheConf),
		AdminNotificationModel:   system.NewAdminNotificationModel(conn, cacheConf),
		DailyShortSentenceModel:  misc.NewDailyShortSentenceModel(conn, cacheConf),
		VideoModel:               video.NewVideoModel(conn, cacheConf),
		SdkKeyModel:              sdk.NewSdkKeyModel(conn, cacheConf),
		SdkInterfaceModel:        sdk.NewSdkInterfaceModel(conn, cacheConf),
		SdkKeyApiModel:           sdk.NewSdkKeyApiModel(conn, cacheConf),
		SdkCallLogModel:          sdk.NewSdkCallLogModel(conn, cacheConf),
		AdminTaskModel:           task.NewAdminTaskModel(conn, cacheConf),
		BlogTagModel:             blog.NewBlogTagModel(conn, cacheConf),
		BlogArticleModel:         blog.NewBlogArticleModel(conn, cacheConf),
		BlogArticleTagModel:      blog.NewBlogArticleTagModel(conn, cacheConf),
		BlogArticleAuditModel:    blog.NewBlogArticleAuditModel(conn, cacheConf),
		BlogFriendLinkModel:      blog.NewBlogFriendLinkModel(conn, cacheConf),
		BlogSocialInfoModel:      blog.NewBlogSocialInfoModel(conn, cacheConf),
	}, nil
}

// BuildSources 根据配置初始化数据源，供 ServiceContext 调用。
func BuildSources(cfg config.Config) (*Repository, error) {
	conn, err := NewSQLConn(cfg.Database)
	if err != nil {
		return nil, errors.Wrap(err, "init sqlx connection")
	}
	cacheConf := BuildCacheConf(cfg.Redis)
	// 创建 go-zero 的 Redis 客户端
	rdb, err := redis.NewRedis(redis.RedisConf{
		Host: cfg.Redis.Address,
		Pass: cfg.Redis.Password,
		Type: "node",
	})
	if err != nil {
		return nil, errors.Wrap(err, "init redis")
	}
	return NewRepository(conn, cacheConf, rdb)
}
