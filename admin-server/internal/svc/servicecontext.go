package svc

import (
	"os"
	"path/filepath"

	"postapocgame/admin-server/internal/config"
	"postapocgame/admin-server/internal/consts"
	"postapocgame/admin-server/internal/domain/task"
	"postapocgame/admin-server/internal/domain/task/executors"
	"postapocgame/admin-server/internal/hub"
	"postapocgame/admin-server/internal/interfaces"
	"postapocgame/admin-server/internal/repository"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
)

type ServiceContext struct {
	Config                       config.Config
	Repository                   *repository.Repository
	ChatHub                      *hub.ChatHub
	TaskExecutors                map[int]interfaces.TaskExecutor
	TaskScheduler                *task.TaskScheduler
	AuthMiddleware               rest.Middleware
	ApiEnabledMiddleware         rest.Middleware
	PermissionMiddleware         rest.Middleware
	OperationLogMiddleware       rest.Middleware
	PublicOperationLogMiddleware rest.Middleware
	RateLimitMiddleware          rest.Middleware
	PerformanceMiddleware        rest.Middleware
	CorsMiddleware               rest.Middleware
	SDKAuthMiddleware            rest.Middleware
	SDKRateLimitMiddleware       rest.Middleware
	SDKCallLogMiddleware         rest.Middleware
}

func NewServiceContext(c config.Config) (*ServiceContext, error) {
	repo, err := repository.BuildSources(c)
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(consts.UploadDir, 0755); err != nil {
		logx.Errorf("创建 uploads 目录失败: %v", err)
		return nil, err
	}
	absPath, _ := filepath.Abs(consts.UploadDir)
	logx.Infof("uploads 目录已初始化: %s", absPath)

	chatHub := hub.NewChatHub()
	go chatHub.Run()

	taskExecutors := make(map[int]interfaces.TaskExecutor)
	taskExecutors[1] = executors.NewExcelExportExecutor(repo)

	svcCtx := &ServiceContext{
		Config:        c,
		Repository:    repo,
		ChatHub:       chatHub,
		TaskExecutors: taskExecutors,
	}

	taskScheduler := task.NewTaskScheduler(repo, chatHub, taskExecutors)
	taskScheduler.Start()
	svcCtx.TaskScheduler = taskScheduler

	return svcCtx, nil
}
