package svc

import (
	"os"
	"path/filepath"

	"postapocgame/admin-server/internal/config"
	"postapocgame/admin-server/internal/consts"
	"postapocgame/admin-server/internal/hub"
	"postapocgame/admin-server/internal/repository"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
)

type ServiceContext struct {
	Config                       config.Config
	Repository                   *repository.Repository
	ChatHub                      *hub.ChatHub
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

	return &ServiceContext{
		Config:     c,
		Repository: repo,
		ChatHub:    chatHub,
		// AuthMiddleware 和 PermissionMiddleware 需要在外部初始化，避免循环依赖
	}, nil
}
