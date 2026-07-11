package repository

import (
	"context"

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

// Transact 在单个 MySQL 事务内执行 fn。
// fn 收到的 txRepo 是 r 的克隆：DB 字段与全部 37 个 *Model 字段都已经换绑到本次事务的 session 上。
// fn 内部必须只通过 txRepo 访问数据，不能继续闭包引用外层 r —— 否则读写会跳出事务边界。
// 不做嵌套事务检测：调用方需保证不在已经开启的 Transact 内部再次调用 Transact。
func (r *Repository) Transact(ctx context.Context, fn func(ctx context.Context, txRepo *Repository) error) error {
	return r.DB.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		return fn(ctx, r.withSession(session))
	})
}

// withSession 返回一个新的 *Repository，DB 字段与全部 37 个 *Model 字段都换绑到给定的事务 session。
// CacheConf/Redis/BusinessCache 保持不变——事务内不使用查询缓存。
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
		BlogTagModel:             r.BlogTagModel.WithSession(session),
		BlogArticleModel:         r.BlogArticleModel.WithSession(session),
		BlogArticleTagModel:      r.BlogArticleTagModel.WithSession(session),
		BlogArticleAuditModel:    r.BlogArticleAuditModel.WithSession(session),
		BlogFriendLinkModel:      r.BlogFriendLinkModel.WithSession(session),
		BlogSocialInfoModel:      r.BlogSocialInfoModel.WithSession(session),
	}
}
