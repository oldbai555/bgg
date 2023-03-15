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
	"github.com/storyicon/grbac"
	"time"
)

func StartServer() error {
	// 初始化数据库
	err := service.InitDao(context.Background(), conf.Global.MysqlConf.Dsn())
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
	h.Use(Cors(), RegisterUuidTrace(), gin.LoggerWithFormatter(defaultLogFormatter), gin.Recovery(), RegisterShowReq())

	// 权限管理
	// 在这里，我们通过“grbac.WithLoader”接口使用自定义Loader功能
	// 并指定应每分钟调用一次LoadAuthorizationRules函数以获取最新的身份验证规则。
	if Rbac, err = grbac.New(grbac.WithLoader(LoadAuthorizationRules, time.Minute)); err != nil {
		log.Errorf("err is : %v", err)
		return err
	}

	// 初始化接口
	InitSysApi(h)

	// 启动服务
	err = h.Run(fmt.Sprintf(":%d", conf.Global.ServerConf.Port))
	if err != nil {
		log.Errorf("err is %v", err)
		return err
	}
	return nil
}

func InitSysApi(h *gin.Engine) {
	registerLbuserApi(h)
	registerLbblogApi(h)
	registerStoreApi(h)
	registerPublicApi(h)
	registerWechatGzhApi(h)
}
