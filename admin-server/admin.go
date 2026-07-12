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
	"postapocgame/admin-server/internal/consumer"
	"postapocgame/admin-server/internal/handler"
	iamcallbacksrv "postapocgame/admin-server/internal/rpcserver/iamcallback"
	taskcallbacksrv "postapocgame/admin-server/internal/rpcserver/taskcallback"
	"postapocgame/admin-server/internal/svc"
	appwire "postapocgame/admin-server/internal/wire"
	iamcallbackpb "postapocgame/admin-server/pkg/iamcallback/pb"
	taskcallbackpb "postapocgame/admin-server/pkg/taskcallback/pb"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"postapocgame/admin-server/internal/model"
	"postapocgame/admin-server/internal/model/iam"
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
	conf.MustLoad(*configFile, &c, conf.UseEnv())

	if c.JWT.AccessSecret == "" || c.JWT.RefreshSecret == "" {
		log.Fatalf("JWT_ACCESS_SECRET / JWT_REFRESH_SECRET 未设置，拒绝以空密钥启动")
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

	ctx, cleanup, err := appwire.InitializeApp(c)
	if err != nil {
		log.Fatalf("init service context: %v", err)
	}
	defer cleanup()

	handler.RegisterHandlers(server, ctx)
	// 注册自定义路由（WebSocket 等）
	handler.RegisterCustomRoutes(server, ctx)
	// 同步路由到 admin_api 表
	syncRoutesToAdminAPI(ctx, server)

	// TaskCallback zrpc server：与 REST server 并存，供 services/task/（task-rpc）回调取导出
	// 数据/登记导出文件。见 16-rpc-conventions.md、17-async-eventing.md 第 1 节；sdk_call_log
	// 分支已经改成回调 ctx.SdkRPC（sdk-rpc 拆分完成），admin_file 登记（iam 域）仍在这里；
	// Phase 2 iam-rpc 真正拆分后这部分原样搬过去。
	taskCallbackServer := zrpc.MustNewServer(c.TaskCallbackRPCConf, func(grpcServer *grpc.Server) {
		taskcallbackpb.RegisterTaskCallbackServer(grpcServer, taskcallbacksrv.NewServer(ctx.Repository, ctx.SdkRPC))
	})
	defer taskCallbackServer.Stop()

	// IamCallback zrpc server：与 REST server 并存，供 services/chat/（chat-rpc）回调枚举
	// 存量用户 / 取用户展示信息。见 pkg/iamcallback 包注释、internal/rpcserver/iamcallback/。
	// iam 域还没拆分成独立服务前的临时方案，和 TaskCallback 同一个模式。
	iamCallbackServer := zrpc.MustNewServer(c.IamCallbackRPCConf, func(grpcServer *grpc.Server) {
		iamcallbackpb.RegisterIamCallbackServer(grpcServer, iamcallbacksrv.NewServer(ctx.Domain))
	})
	defer iamCallbackServer.Stop()

	// task 通知消费者：消费 task-rpc 发布的 stream:task.notification，写 admin_notification +
	// 推 WS（推 WS 现在通过 ctx.ChatRPC.PushToUser 回调 chat-rpc，见
	// internal/consumer/task_notification_consumer.go 包注释、17-async-eventing.md）。
	taskNotificationConsumer := consumer.NewTaskNotificationConsumer(ctx.Repository.Redis, ctx.Repository, ctx.ChatRPC)
	taskNotificationConsumer.Start()
	defer taskNotificationConsumer.Stop()

	// 设置优雅关闭：监听系统信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 在 goroutine 中启动服务器
	go func() {
		logx.Infof("Starting server at %s:%d...", c.Host, c.Port)
		server.Start()
	}()
	go func() {
		logx.Infof("Starting TaskCallback rpc server at %s...", c.TaskCallbackRPCConf.ListenOn)
		taskCallbackServer.Start()
	}()
	go func() {
		logx.Infof("Starting IamCallback rpc server at %s...", c.IamCallbackRPCConf.ListenOn)
		iamCallbackServer.Start()
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
		data := &iam.AdminApi{
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
