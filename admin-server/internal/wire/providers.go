package wire

import (
	"os"
	"path/filepath"

	"postapocgame/admin-server/internal/config"
	"postapocgame/admin-server/internal/consts"
	"postapocgame/admin-server/internal/domain/task"
	"postapocgame/admin-server/internal/domain/task/executors"
	"postapocgame/admin-server/internal/hub"
	"postapocgame/admin-server/internal/interfaces"
	"postapocgame/admin-server/internal/middleware"
	"postapocgame/admin-server/internal/repository"
	"postapocgame/admin-server/internal/repository/registry"
	"postapocgame/admin-server/internal/svc"

	"github.com/google/wire"
	"github.com/zeromicro/go-zero/core/logx"
)

// ProviderSet 组合根依赖注入集合。
var ProviderSet = wire.NewSet(
	provideRepository,
	provideDomain,
	provideChatHub,
	provideTaskExecutors,
	provideTaskScheduler,

	middleware.NewAuthMiddleware,
	middleware.NewApiEnabledMiddleware,
	providePermissionMiddleware, // 适配函数，不是 middleware.NewPermissionMiddleware 本身，见下方注释
	middleware.NewOperationLogMiddleware,
	middleware.NewPublicOperationLogMiddleware,
	middleware.NewRateLimitMiddleware,
	middleware.NewPerformanceMiddleware,
	middleware.NewCorsMiddleware,
	middleware.NewSDKAuthMiddleware,
	middleware.NewSDKRateLimitMiddleware,
	middleware.NewSDKCallLogMiddleware,
	provideMiddlewareBundle,

	provideServiceContext,
)

func provideRepository(c config.Config) (*repository.Repository, error) {
	if err := initUploadDir(); err != nil {
		return nil, err
	}
	return repository.BuildSources(c)
}

func initUploadDir() error {
	if err := os.MkdirAll(consts.UploadDir, 0o755); err != nil {
		logx.Errorf("创建 uploads 目录失败: %v", err)
		return err
	}
	absPath, _ := filepath.Abs(consts.UploadDir)
	logx.Infof("uploads 目录已初始化: %s", absPath)
	return nil
}

func provideDomain(repo *repository.Repository) *registry.Domain {
	return registry.NewDomain(repo)
}

func provideChatHub() *hub.ChatHub {
	chatHub := hub.NewChatHub()
	go chatHub.Run()
	return chatHub
}

func provideTaskExecutors(repo *repository.Repository) map[int]interfaces.TaskExecutor {
	executorsMap := make(map[int]interfaces.TaskExecutor)
	executorsMap[1] = executors.NewExcelExportExecutor(repo)
	return executorsMap
}

func provideTaskScheduler(
	repo *repository.Repository,
	chatHub *hub.ChatHub,
	executorsMap map[int]interfaces.TaskExecutor,
) *task.TaskScheduler {
	scheduler := task.NewTaskScheduler(repo, chatHub, executorsMap)
	scheduler.Start()
	return scheduler
}

// providePermissionMiddleware 是 PermissionMiddleware 的 Wire 适配函数：PermissionMiddleware
// 构造函数吃的是 *iamdomain.PermissionResolver，但那不是一个独立的 Wire 节点，只是
// *registry.Domain 结构体里的一个字段（domain.IAM.PermissionResolver）。Wire 没法凭空
// "生产"出一个游离的 *iamdomain.PermissionResolver，所以需要这一层适配：用已有的
// *registry.Domain 节点取出字段，再调用 middleware.NewPermissionMiddleware。
func providePermissionMiddleware(domain *registry.Domain) *middleware.PermissionMiddleware {
	return middleware.NewPermissionMiddleware(domain.IAM.PermissionResolver)
}

// provideMiddlewareBundle 是唯一还需要手写的 assembler：MiddlewareBundle 的 11 个字段
// 类型全部是 rest.Middleware（同一个具名类型），Wire 无法自动区分该把哪个 provider
// 的结果填进哪个字段，所以仍然需要一个手写函数——区别在于它现在只依赖 11 个互不相同
// 的具体中间件指针类型，不再依赖尚未构造完成的 svcCtx 本身。
func provideMiddlewareBundle(
	auth *middleware.AuthMiddleware,
	apiEnabled *middleware.ApiEnabledMiddleware,
	permission *middleware.PermissionMiddleware,
	operationLog *middleware.OperationLogMiddleware,
	publicOperationLog *middleware.PublicOperationLogMiddleware,
	rateLimit *middleware.RateLimitMiddleware,
	performance *middleware.PerformanceMiddleware,
	cors *middleware.CorsMiddleware,
	sdkAuth *middleware.SDKAuthMiddleware,
	sdkRateLimit *middleware.SDKRateLimitMiddleware,
	sdkCallLog *middleware.SDKCallLogMiddleware,
) *MiddlewareBundle {
	return &MiddlewareBundle{
		Auth:               auth.Handle,
		ApiEnabled:         apiEnabled.Handle,
		Permission:         permission.Handle,
		OperationLog:       operationLog.Handle,
		PublicOperationLog: publicOperationLog.Handle,
		RateLimit:          rateLimit.Handle,
		Performance:        performance.Handle,
		Cors:               cors.Handle,
		SDKAuth:            sdkAuth.Handle,
		SDKRateLimit:       sdkRateLimit.Handle,
		SDKCallLog:         sdkCallLog.Handle,
	}
}

func provideServiceContext(
	c config.Config,
	repo *repository.Repository,
	domain *registry.Domain,
	chatHub *hub.ChatHub,
	taskExecutors map[int]interfaces.TaskExecutor,
	taskScheduler *task.TaskScheduler,
	mw *MiddlewareBundle,
) (*svc.ServiceContext, func()) {
	svcCtx := &svc.ServiceContext{
		Config:                       c,
		Repository:                   repo,
		Domain:                       domain,
		ChatHub:                      chatHub,
		TaskExecutors:                taskExecutors,
		TaskScheduler:                taskScheduler,
		AuthMiddleware:               mw.Auth,
		ApiEnabledMiddleware:         mw.ApiEnabled,
		PermissionMiddleware:         mw.Permission,
		OperationLogMiddleware:       mw.OperationLog,
		PublicOperationLogMiddleware: mw.PublicOperationLog,
		RateLimitMiddleware:          mw.RateLimit,
		PerformanceMiddleware:        mw.Performance,
		CorsMiddleware:               mw.Cors,
		SDKAuthMiddleware:            mw.SDKAuth,
		SDKRateLimitMiddleware:       mw.SDKRateLimit,
		SDKCallLogMiddleware:         mw.SDKCallLog,
	}

	cleanup := func() {
		if taskScheduler != nil {
			taskScheduler.Stop()
			logx.Infof("任务调度器已停止")
		}
	}
	return svcCtx, cleanup
}
