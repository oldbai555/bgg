package svc

import (
	"os"
	"path/filepath"
	"postapocgame/admin-server/internal/task"
	"postapocgame/admin-server/internal/task/executors"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
	"postapocgame/admin-server/internal/config"
	"postapocgame/admin-server/internal/consts"
	"postapocgame/admin-server/internal/hub"
	"postapocgame/admin-server/internal/interfaces"
	"postapocgame/admin-server/internal/repository"
)

type ServiceContext struct {
	Config                       config.Config
	Repository                   *repository.Repository
	ChatHub                      *hub.ChatHub
	TaskExecutors                map[int]interfaces.TaskExecutor // 任务执行器映射
	TaskScheduler                *task.TaskScheduler             // 任务调度器
	AuthMiddleware               rest.Middleware
	ApiEnabledMiddleware         rest.Middleware
	PermissionMiddleware         rest.Middleware
	OperationLogMiddleware       rest.Middleware
	PublicOperationLogMiddleware rest.Middleware
	RateLimitMiddleware          rest.Middleware
	PerformanceMiddleware        rest.Middleware
	SDKAuthMiddleware            rest.Middleware
	SDKRateLimitMiddleware       rest.Middleware
	SDKCallLogMiddleware         rest.Middleware
}

func NewServiceContext(c config.Config) (*ServiceContext, error) {
	repo, err := repository.BuildSources(c)
	if err != nil {
		return nil, err
	}

	// 初始化 uploads 目录（如果不存在）
	if err := os.MkdirAll(consts.UploadDir, 0755); err != nil {
		logx.Errorf("创建 uploads 目录失败: %v", err)
		return nil, err
	}
	absPath, _ := filepath.Abs(consts.UploadDir)
	logx.Infof("uploads 目录已初始化: %s", absPath)

	// 初始化 ChatHub
	chatHub := hub.NewChatHub()
	go chatHub.Run()

	// 初始化任务执行器映射
	taskExecutors := make(map[int]interfaces.TaskExecutor)
	// 注册Excel导出执行器（task_type=1）
	taskExecutors[1] = executors.NewExcelExportExecutor(repo)

	// 创建ServiceContext（先不包含TaskScheduler，避免循环依赖）
	svcCtx := &ServiceContext{
		Config:        c,
		Repository:    repo,
		ChatHub:       chatHub,
		TaskExecutors: taskExecutors,
		// AuthMiddleware 和 PermissionMiddleware 需要在外部初始化，避免循环依赖
	}

	// 初始化任务调度器（传入Repository和ChatHub，避免循环依赖）
	taskScheduler := task.NewTaskScheduler(repo, chatHub, taskExecutors)

	// 启动任务调度器
	taskScheduler.Start()

	// 将TaskScheduler赋值给ServiceContext
	svcCtx.TaskScheduler = taskScheduler

	return svcCtx, nil
}
