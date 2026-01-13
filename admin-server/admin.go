// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"postapocgame/admin-server/internal/consts"
	"strings"
	"syscall"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"postapocgame/admin-server/internal/config"
	"postapocgame/admin-server/internal/handler"
	"postapocgame/admin-server/internal/middleware"
	"postapocgame/admin-server/internal/model"
	"postapocgame/admin-server/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var (
	// GIT_COMMIT_VERSION 编译时通过 -ldflags 传入的 git 提交版本号
	GIT_COMMIT_VERSION string
)

var (
	configFile           = flag.String("f", "", "the config file")
	mysqlConfigFile      = flag.String("mysql-config", "", "MySQL config file path (default: /etc/work/mysql.json)")
	redisConfigFile      = flag.String("redis-config", "", "Redis config file path (default: /etc/work/redis.json)")
	middlewareConfigFile = flag.String("middleware-config", "", "Middleware config file path (default: etc/middleware.yaml)")
)

func main() {
	flag.Parse()

	// 打印版本信息
	if GIT_COMMIT_VERSION != "" {
		log.Printf("GIT_COMMIT_VERSION: %s", GIT_COMMIT_VERSION)
	} else {
		log.Printf("GIT_COMMIT_VERSION: dev (not set)")
	}

	var c config.Config
	conf.MustLoad(*configFile, &c)

	// 从外部文件加载 MySQL、Redis 和中间件配置（如果存在）
	if err := config.MergeExternalConfig(&c, *mysqlConfigFile, *redisConfigFile, *middlewareConfigFile); err != nil {
		log.Fatalf("加载外部配置失败: %v", err)
	}

	err := logx.SetUp(logx.LogConf{
		Encoding: "plain",
	})
	if err != nil {
		log.Fatalf("Failed to set up logging: %v", err)
	}

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx, err := svc.NewServiceContext(c)
	if err != nil {
		log.Fatalf("init service context: %v", err)
	}

	// 初始化中间件（避免循环依赖，在外部初始化）
	authMiddleware := middleware.NewAuthMiddleware(ctx)
	apiEnabledMiddleware := middleware.NewApiEnabledMiddleware(ctx)
	permissionMiddleware := middleware.NewPermissionMiddleware(ctx)
	operationLogMiddleware := middleware.NewOperationLogMiddleware(ctx)
	publicOperationLogMiddleware := middleware.NewPublicOperationLogMiddleware(ctx)
	rateLimitMiddleware := middleware.NewRateLimitMiddleware(ctx)
	performanceMiddleware := middleware.NewPerformanceMiddleware(ctx)
	sdkAuthMiddleware := middleware.NewSDKAuthMiddleware(ctx)
	sdkRateLimitMiddleware := middleware.NewSDKRateLimitMiddleware(ctx)
	sdkCallLogMiddleware := middleware.NewSDKCallLogMiddleware(ctx)
	ctx.AuthMiddleware = authMiddleware.Handle
	ctx.ApiEnabledMiddleware = apiEnabledMiddleware.Handle
	ctx.PermissionMiddleware = permissionMiddleware.Handle
	ctx.OperationLogMiddleware = operationLogMiddleware.Handle
	ctx.PublicOperationLogMiddleware = publicOperationLogMiddleware.Handle
	ctx.RateLimitMiddleware = rateLimitMiddleware.Handle
	ctx.PerformanceMiddleware = performanceMiddleware.Handle
	ctx.SDKAuthMiddleware = sdkAuthMiddleware.Handle
	ctx.SDKRateLimitMiddleware = sdkRateLimitMiddleware.Handle
	ctx.SDKCallLogMiddleware = sdkCallLogMiddleware.Handle

	handler.RegisterHandlers(server, ctx)
	// 注册自定义路由（WebSocket 等）
	handler.RegisterCustomRoutes(server, ctx)
	// 同步路由到 admin_api 表
	syncRoutesToAdminAPI(ctx, server)

	// 设置优雅关闭：监听系统信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 在 goroutine 中启动服务器
	go func() {
		logx.Infof("Starting server at %s:%d...", c.Host, c.Port)
		server.Start()
	}()

	// 等待关闭信号
	<-sigChan
	logx.Infof("收到关闭信号，开始优雅关闭...")

	// 停止任务调度器
	if ctx.TaskScheduler != nil {
		ctx.TaskScheduler.Stop()
		logx.Infof("任务调度器已停止")
	}

	logx.Infof("服务器已关闭")
}

// syncRoutesToAdminAPI 将已注册的路由同步到 admin_api 表（method+path 唯一）
func syncRoutesToAdminAPI(ctx *svc.ServiceContext, server *rest.Server) {
	logx.Infof("====== 同步路由到 admin_api 表开始 ======")
	routes := server.Routes()
	if len(routes) == 0 {
		return
	}

	now := time.Now().Unix()
	bg := context.Background()

	for _, r := range routes {
		method := strings.ToUpper(strings.TrimSpace(r.Method))
		path := strings.TrimSpace(r.Path)
		if method == "" || path == "" {
			continue
		}

		// 已存在则跳过
		_, err := ctx.Repository.AdminApiModel.FindOneByMethodPath(bg, method, path)
		if err == nil {
			continue
		}
		if err != model.ErrNotFound {
			logx.Errorf("检查接口 %s %s 失败: %v", method, path, err)
			continue
		}

		apiName := fmt.Sprintf("%s_%s", method, sanitizePathForName(path))
		data := &model.AdminApi{
			Name:        apiName,
			Method:      method,
			Path:        path,
			Description: sql.NullString{}, // 由管理员补充
			Status:      consts.Open,      // 默认启用
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		if _, err = ctx.Repository.AdminApiModel.Insert(bg, data); err != nil {
			logx.Errorf("写入接口 %s %s 失败: %v", method, path, err)
		}
	}
	logx.Infof("====== 同步路由到 admin_api 表结束 ======")
}

// sanitizePathForName 将路径转换为名称可用的格式
func sanitizePathForName(path string) string {
	path = strings.Trim(path, "/")
	if path == "" {
		return "ROOT"
	}
	path = strings.ReplaceAll(path, "/", "_")
	path = strings.ReplaceAll(path, ":", "_")
	return path
}
