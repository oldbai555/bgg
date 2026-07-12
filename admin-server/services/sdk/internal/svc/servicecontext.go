package svc

import (
	"log"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"

	"postapocgame/admin-server/services/sdk/internal/config"
	sdkdomain "postapocgame/admin-server/services/sdk/internal/domain/sdk"
	"postapocgame/admin-server/services/sdk/internal/repository"
	sdkrepo "postapocgame/admin-server/services/sdk/internal/repository/sdk"
)

type ServiceContext struct {
	Config           config.Config
	Admin            *sdkrepo.SdkAdminRepository
	Public           *sdkrepo.SdkRepository
	Service          *sdkdomain.SDKService
	RateLimitDefault int64
}

func NewServiceContext(c config.Config) *ServiceContext {
	if c.Mysql.DSN == "" {
		log.Fatalf("sdk-rpc: Mysql.DSN 未配置")
	}
	conn := sqlx.NewMysql(c.Mysql.DSN)

	cacheConf := cache.CacheConf{
		{
			RedisConf: redis.RedisConf{Host: c.SdkRedis.Address, Pass: c.SdkRedis.Password, Type: "node"},
			Weight:    100,
		},
	}
	store := repository.NewStore(conn, cacheConf)

	return &ServiceContext{
		Config:           c,
		Admin:            sdkrepo.NewSdkAdminRepository(store),
		Public:           sdkrepo.NewSdkRepository(store),
		Service:          sdkdomain.NewSDKService(store),
		RateLimitDefault: c.RateLimitDefault,
	}
}
