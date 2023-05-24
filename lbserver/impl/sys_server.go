package impl

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/oldbai555/bgg/lbserver/impl/cache"
	"github.com/oldbai555/bgg/lbserver/impl/conf"
	"github.com/oldbai555/bgg/lbserver/impl/service"
	"github.com/oldbai555/bgg/lbserver/impl/storage"
	"github.com/oldbai555/lbtool/log"
)

func Server(ctx context.Context) error {
	// 初始化数据库
	err := service.InitDao(ctx, conf.Global.MysqlConf.Dsn())
	if err != nil {
		log.Errorf("err:%v", err)
		panic(err)
	}

	// 初始化缓存
	cache.InitCacheHelper(conf.Global.RedisConf)

	// 初始化存储桶
	storage.InitStorage(conf.Global.StorageConf)

	gin.DefaultWriter = log.GetWriter()
	h := gin.New()

	// 配置跨域中间件
	// 链路追踪
	h.Use(Cors(), RegisterUuidTrace(), gin.LoggerWithFormatter(defaultLogFormatter), gin.Recovery(), RegisterJwt())

	// 注册api
	register(h)

	// 启动服务
	err = h.Run(fmt.Sprintf(":%d", conf.Global.ServerConf.Port))
	if err != nil {
		log.Errorf("err is %v", err)
		return err
	}
	return nil
}
