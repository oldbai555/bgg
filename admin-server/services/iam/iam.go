package main

import (
	"flag"
	"fmt"

	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/config"
	"postapocgame/admin-server/services/iam/internal/consumer"
	"postapocgame/admin-server/services/iam/internal/server"
	"postapocgame/admin-server/services/iam/internal/svc"
	iamcallbackpb "postapocgame/admin-server/pkg/iamcallback/pb"
	taskcallbackpb "postapocgame/admin-server/pkg/taskcallback/pb"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/iam.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c, conf.UseEnv())
	ctx := svc.NewServiceContext(c)

	// 同一个 zrpc.Server 上注册三个 gRPC 服务：Iam（本服务的原生契约）+ TaskCallback/
	// IamCallback（从单体 internal/rpcserver/{taskcallback,iamcallback}/ 原样搬迁过来的
	// 服务端实现，契约不变，task-rpc/chat-rpc/content-rpc 侧代码零改动，只需要把回调
	// endpoint 配置改指向 iam-rpc，见 18-service-extraction-runbook.md 2.5 节）。
	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		iam.RegisterIamServer(grpcServer, server.NewIamServer(ctx))
		taskcallbackpb.RegisterTaskCallbackServer(grpcServer, server.NewTaskCallbackServer(ctx.Repository, ctx.SdkRPC))
		iamcallbackpb.RegisterIamCallbackServer(grpcServer, server.NewIamCallbackServer(ctx.Domain))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	// task 通知消费者：消费 task-rpc 发布的 stream:task.notification，写 admin_notification +
	// 回调 chat-rpc.PushToUser 推 WS（原单体 internal/consumer/task_notification_consumer.go
	// 原样搬迁，consumer group 从 "iam-chat-task-notify" 改名 "iam-rpc-task-notify"——搬迁前
	// 是"临时合并消费者"的命名，现在 iam-rpc 是真正的最终归属，见 progress.md 本轮条目）。
	taskNotificationConsumer := consumer.NewTaskNotificationConsumer(ctx.Repository.Redis, ctx.Repository, ctx.ChatRPC)
	taskNotificationConsumer.Start()
	defer taskNotificationConsumer.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
