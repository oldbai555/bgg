package svc

import (
	"log"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"

	iamcallbackpb "postapocgame/admin-server/pkg/iamcallback/pb"
	"postapocgame/admin-server/services/content/internal/config"
	contentdomain "postapocgame/admin-server/services/content/internal/domain/content"
	"postapocgame/admin-server/services/content/internal/repository"
	blogrepo "postapocgame/admin-server/services/content/internal/repository/blog"
	videorepo "postapocgame/admin-server/services/content/internal/repository/video"
)

type ServiceContext struct {
	Config config.Config
	// Store 是聚合 blog 六张表 + video 一张表 Model 的句柄，下面 7 个字段是绑定在它上面的
	// repository；BlogArticleService 需要事务时直接调 Store.Transact，和 services/sdk 的
	// SDKService 持有 *repository.Store 是同一个模式。
	Store          *repository.Store
	BlogArticle    blogrepo.BlogArticleRepository
	BlogArticleTag blogrepo.BlogArticleTagRepository
	ArticleAudit   blogrepo.BlogArticleAuditRepository
	FriendLink     blogrepo.BlogFriendLinkRepository
	SocialInfo     blogrepo.BlogSocialInfoRepository
	Tag            blogrepo.BlogTagRepository
	Video          videorepo.VideoRepository
	ArticleService *contentdomain.BlogArticleService

	// IamCallback 回调单体内嵌的 pkg/iamcallback.IamCallback server：PublicBlogAuthorInfo
	// 需要展示的用户信息 + BlogArticleAudit/BlogArticleAuditUnpublish 的审计日志写入，
	// 见 pkg/iamcallback 包注释。
	IamCallback iamcallbackpb.IamCallbackClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	if c.Mysql.DSN == "" {
		log.Fatalf("content-rpc: Mysql.DSN 未配置")
	}
	conn := sqlx.NewMysql(c.Mysql.DSN)

	cacheConf := cache.CacheConf{
		{
			RedisConf: redis.RedisConf{Host: c.ContentRedis.Address, Pass: c.ContentRedis.Password, Type: "node"},
			Weight:    100,
		},
	}
	store := repository.NewStore(conn, cacheConf)

	iamCallbackClient := iamcallbackpb.NewIamCallbackClient(zrpc.MustNewClient(c.IamCallbackRpc).Conn())

	return &ServiceContext{
		Config:         c,
		Store:          store,
		BlogArticle:    blogrepo.NewBlogArticleRepository(store),
		BlogArticleTag: blogrepo.NewBlogArticleTagRepository(store),
		ArticleAudit:   blogrepo.NewBlogArticleAuditRepository(store),
		FriendLink:     blogrepo.NewBlogFriendLinkRepository(store),
		SocialInfo:     blogrepo.NewBlogSocialInfoRepository(store),
		Tag:            blogrepo.NewBlogTagRepository(store),
		Video:          videorepo.NewVideoRepository(store),
		ArticleService: contentdomain.NewBlogArticleService(store),
		IamCallback:    iamCallbackClient,
	}
}
