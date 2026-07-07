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

func provideServiceContext(
	c config.Config,
	repo *repository.Repository,
	domain *registry.Domain,
	chatHub *hub.ChatHub,
	taskExecutors map[int]interfaces.TaskExecutor,
	taskScheduler *task.TaskScheduler,
) (*svc.ServiceContext, func()) {
	svcCtx := &svc.ServiceContext{
		Config:        c,
		Repository:    repo,
		Domain:        domain,
		ChatHub:       chatHub,
		TaskExecutors: taskExecutors,
		TaskScheduler: taskScheduler,
	}
	mw := buildMiddlewareBundle(svcCtx)
	svcCtx.AuthMiddleware = mw.Auth
	svcCtx.ApiEnabledMiddleware = mw.ApiEnabled
	svcCtx.PermissionMiddleware = mw.Permission
	svcCtx.OperationLogMiddleware = mw.OperationLog
	svcCtx.PublicOperationLogMiddleware = mw.PublicOperationLog
	svcCtx.RateLimitMiddleware = mw.RateLimit
	svcCtx.PerformanceMiddleware = mw.Performance
	svcCtx.CorsMiddleware = mw.Cors
	svcCtx.SDKAuthMiddleware = mw.SDKAuth
	svcCtx.SDKRateLimitMiddleware = mw.SDKRateLimit
	svcCtx.SDKCallLogMiddleware = mw.SDKCallLog

	cleanup := func() {
		if taskScheduler != nil {
			taskScheduler.Stop()
			logx.Infof("任务调度器已停止")
		}
	}
	return svcCtx, cleanup
}

func buildMiddlewareBundle(svcCtx *svc.ServiceContext) *MiddlewareBundle {
	return &MiddlewareBundle{
		Auth:               middleware.NewAuthMiddleware(svcCtx).Handle,
		ApiEnabled:         middleware.NewApiEnabledMiddleware(svcCtx).Handle,
		Permission:         middleware.NewPermissionMiddleware(svcCtx).Handle,
		OperationLog:       middleware.NewOperationLogMiddleware(svcCtx).Handle,
		PublicOperationLog: middleware.NewPublicOperationLogMiddleware(svcCtx).Handle,
		RateLimit:          middleware.NewRateLimitMiddleware(svcCtx).Handle,
		Performance:        middleware.NewPerformanceMiddleware(svcCtx).Handle,
		Cors:               middleware.NewCorsMiddleware().Handle,
		SDKAuth:            middleware.NewSDKAuthMiddleware(svcCtx).Handle,
		SDKRateLimit:       middleware.NewSDKRateLimitMiddleware(svcCtx).Handle,
		SDKCallLog:         middleware.NewSDKCallLogMiddleware(svcCtx).Handle,
	}
}
