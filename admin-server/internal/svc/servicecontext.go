package svc

import (
	"postapocgame/admin-server/internal/config"
	"postapocgame/admin-server/internal/domain/task"
	"postapocgame/admin-server/internal/hub"
	"postapocgame/admin-server/internal/interfaces"
	"postapocgame/admin-server/internal/repository"
	"postapocgame/admin-server/internal/repository/registry"

	"github.com/zeromicro/go-zero/rest"
)

type ServiceContext struct {
	Config                       config.Config
	Repository                   *repository.Repository
	Domain                       *registry.Domain
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
