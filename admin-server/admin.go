// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/zeromicro/go-zero/core/logx"

	"postapocgame/admin-server/internal/config"
	"postapocgame/admin-server/internal/handler"
	"postapocgame/admin-server/internal/svc"
	appwire "postapocgame/admin-server/internal/wire"
	"postapocgame/admin-server/services/iam/iamclient"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var (
	// GIT_COMMIT_VERSION 编译时通过 -ldflags 传入的 git 提交版本号
	GIT_COMMIT_VERSION string
)

var (
	configFile           = flag.String("f", "", "the config file")
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
	conf.MustLoad(*configFile, &c, conf.UseEnv())

	if c.JWT.AccessSecret == "" || c.JWT.RefreshSecret == "" {
		log.Fatalf("JWT_ACCESS_SECRET / JWT_REFRESH_SECRET 未设置，拒绝以空密钥启动")
	}
	if len(c.IamRpc.Endpoints) == 0 || c.IamRpc.Endpoints[0] == "" {
		log.Fatalf("IAM_RPC_ENDPOINT 未设置，拒绝以空 iam-rpc 地址启动")
	}
	if len(c.IamCallbackRpc.Endpoints) == 0 || c.IamCallbackRpc.Endpoints[0] == "" {
		log.Fatalf("IAM_CALLBACK_RPC_ENDPOINT 未设置，拒绝以空 iam-rpc 地址启动")
	}
	if len(c.TaskRPCConf.Endpoints) == 0 || c.TaskRPCConf.Endpoints[0] == "" {
		log.Fatalf("TASK_RPC_ENDPOINT 未设置，拒绝以空 task-rpc 地址启动")
	}
	if len(c.SdkRPCConf.Endpoints) == 0 || c.SdkRPCConf.Endpoints[0] == "" {
		log.Fatalf("SDK_RPC_ENDPOINT 未设置，拒绝以空 sdk-rpc 地址启动")
	}
	if len(c.ChatRPCConf.Endpoints) == 0 || c.ChatRPCConf.Endpoints[0] == "" {
		log.Fatalf("CHAT_RPC_ENDPOINT 未设置，拒绝以空 chat-rpc 地址启动")
	}
	if len(c.ContentRPCConf.Endpoints) == 0 || c.ContentRPCConf.Endpoints[0] == "" {
		log.Fatalf("CONTENT_RPC_ENDPOINT 未设置，拒绝以空 content-rpc 地址启动")
	}

	// 从外部文件加载 Redis 和中间件配置（如果存在）；iam+system+monitoring+misc 域拆分
	// 成 iam-rpc 后，gateway 不再直连任何 MySQL，不再需要 mysql-config
	if err := config.MergeExternalConfig(&c, *redisConfigFile, *middlewareConfigFile); err != nil {
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

	ctx, cleanup, err := appwire.InitializeApp(c)
	if err != nil {
		log.Fatalf("init service context: %v", err)
	}
	defer cleanup()

	handler.RegisterHandlers(server, ctx)
	// 注册自定义路由（WebSocket 等）
	handler.RegisterCustomRoutes(server, ctx)
	// 同步路由到 admin_api 表（admin_api 表物理属于 iam-rpc，改成一次性批量 RPC）
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

	logx.Infof("服务器已关闭")
}

// syncRoutesToAdminAPI 将已注册的路由同步到 admin_api 表（method+path 唯一）
func syncRoutesToAdminAPI(ctx *svc.ServiceContext, server *rest.Server) {
	logx.Infof("====== 同步路由到 admin_api 表开始 ======")
	routes := server.Routes()
	if len(routes) == 0 {
		return
	}

	routeRefs := make([]*iamclient.ApiRouteRef, 0, len(routes))
	for _, r := range routes {
		method := strings.ToUpper(strings.TrimSpace(r.Method))
		path := strings.TrimSpace(r.Path)
		if method == "" || path == "" {
			continue
		}
		routeRefs = append(routeRefs, &iamclient.ApiRouteRef{Method: method, Path: path})
	}

	if _, err := ctx.IamRPC.SyncApiRoutes(context.Background(), &iamclient.SyncApiRoutesRequest{Routes: routeRefs}); err != nil {
		logx.Errorf("同步路由到 admin_api 表失败: %v", err)
	}
	logx.Infof("====== 同步路由到 admin_api 表结束 ======")
}
