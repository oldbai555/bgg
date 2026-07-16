package svc

import (
	"log"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"

	"postapocgame/admin-server/pkg/taskcallback"
	"postapocgame/admin-server/services/task/internal/config"
	"postapocgame/admin-server/services/task/internal/consts"
	taskdomain "postapocgame/admin-server/services/task/internal/domain/task"
	"postapocgame/admin-server/services/task/internal/domain/task/executors"
	"postapocgame/admin-server/services/task/internal/interfaces"
	taskmodel "postapocgame/admin-server/services/task/internal/model/task"
	"postapocgame/admin-server/services/task/internal/repository"
)

type ServiceContext struct {
	Config    config.Config
	TaskRepo  repository.TaskRepository
	Scheduler *taskdomain.TaskScheduler
}

func NewServiceContext(c config.Config) *ServiceContext {
	if c.Mysql.DSN == "" {
		log.Fatalf("task-rpc: Mysql.DSN 未配置")
	}
	conn := sqlx.NewMysql(c.Mysql.DSN)

	cacheConf := cache.CacheConf{
		{
			RedisConf: redis.RedisConf{Host: c.TaskRedis.Address, Pass: c.TaskRedis.Password, Type: "node"},
			Weight:    100,
		},
	}
	model := taskmodel.NewAdminTaskModel(conn, cacheConf)
	taskRepo := repository.NewTaskRepository(model, conn)

	rds, err := redis.NewRedis(redis.RedisConf{Host: c.TaskRedis.Address, Pass: c.TaskRedis.Password, Type: "node"})
	if err != nil {
		log.Fatalf("task-rpc: 初始化 Redis 失败: %v", err)
	}

	callbackClient, err := taskcallback.NewClient(c.TaskCallbackRpc)
	if err != nil {
		log.Fatalf("task-rpc: 连接 TaskCallback 失败: %v", err)
	}

	// moduleRoutes：当前阶段全部 module 指向同一个单体内嵌 TaskCallback server，
	// 见 executors.ModuleServiceRoute 注释。Phase 2 iam-rpc/sdk-rpc 真正拆分后，
	// 这里改成指向各自服务的 client，其余代码不用动。
	moduleRoutes := executors.ModuleServiceRoute{
		consts.TaskModuleOperationLog:   callbackClient,
		consts.TaskModuleAuditLog:       callbackClient,
		consts.TaskModuleLoginLog:       callbackClient,
		consts.TaskModulePerformanceLog: callbackClient,
		consts.TaskModuleSdkCallLog:     callbackClient,
	}
	exportExecutor := executors.NewGenericExportExecutor(moduleRoutes, callbackClient)

	executorsMap := map[int]interfaces.TaskExecutor{
		exportExecutor.GetType(): exportExecutor,
	}

	notifier := taskdomain.NewTaskNotifier(rds)
	scheduler := taskdomain.NewTaskScheduler(taskRepo, rds, notifier, executorsMap)

	return &ServiceContext{
		Config:    c,
		TaskRepo:  taskRepo,
		Scheduler: scheduler,
	}
}
