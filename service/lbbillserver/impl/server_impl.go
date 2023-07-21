package impl

import (
	"context"
	"fmt"
	"github.com/oldbai555/bgg/service/lbbill"

	"github.com/oldbai555/bgg/pkg/grpc_tool"
	"github.com/oldbai555/bgg/pkg/register"

	"github.com/oldbai555/bgg/service/lbbillserver/impl/conf"
	"github.com/oldbai555/bgg/service/lbbillserver/impl/dao/impl/mysql"
	"github.com/oldbai555/bgg/service/lbbillserver/impl/service"

	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/routine"
	"github.com/oldbai555/lbtool/utils"

	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func Run(ctx *cli.Context) error {
	// 初始化配置
	conf.InitWebTool()

	// 注册表
	mysql.RegisterModel([]interface{}{}...)

	// 初始化mysql
	if err := mysql.RegisterOrm(conf.Global.MysqlConf.Dsn()); err != nil {
		panic(fmt.Sprintf("err is %v", err))
	}

	// 监听端口
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", conf.Global.ServerConf.Port))
	if err != nil {
		log.Errorf("net.Listen err: %v", err)
		return err
	}

	// 新建gRPC服务器实例
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(grpc_tool.Recover(), grpc_tool.AutoValidate()),
	)

	// 在gRPC服务器注册我们的服务
	lbbill.RegisterLbbillServer(grpcServer, &service.ServerImpl)

	newCtx, cancel := context.WithCancel(ctx.Context)
	defer cancel()

	// 注册 grpc 服务节点到 etcd 中
	boundIP, err := utils.GetOutBoundIP()
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	routine.Go(newCtx, func(ctx context.Context) error {
		register.EndPointToEtcd(ctx, fmt.Sprintf("%s:%d", boundIP, conf.Global.ServerConf.Port), lbbill.ServerName)
		return nil
	})

	routine.Go(newCtx, func(ctx context.Context) error {
		signalCh := make(chan os.Signal, 1)
		signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM, os.Kill)
		log.Infof("listening on %s exited notify", lbbill.ServerName)
		<-signalCh
		log.Warnf("exit: close %s server connect", lbbill.ServerName)
		grpcServer.Stop()
		return nil
	})

	// 启动 grpc 服务
	if err := grpcServer.Serve(listener); err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}
