package main

import (
	"flag"
	"fmt"

	"postapocgame/admin-server/services/chat/chat"
	"postapocgame/admin-server/services/chat/internal/config"
	"postapocgame/admin-server/services/chat/internal/server"
	"postapocgame/admin-server/services/chat/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/chat.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c, conf.UseEnv())
	ctx := svc.NewServiceContext(c)

	// 启动 chat onboarding 消费者（消费 stream:chat.user.created，进程内 goroutine），与
	// zrpc server 生命周期绑定，和 services/task/task.go 的 ctx.Scheduler.Start() 同一个模式。
	ctx.OnboardingConsumer.Start()
	defer ctx.OnboardingConsumer.Stop()

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		chat.RegisterChatServer(grpcServer, server.NewChatServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
