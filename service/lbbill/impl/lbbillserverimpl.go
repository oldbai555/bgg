package impl

import (
	"context"
	"fmt"
	"github.com/oldbai555/bgg/client/lbbill"
	"github.com/oldbai555/bgg/pkg/register"
	"github.com/oldbai555/bgg/service/lbbill/impl/conf"
	"github.com/oldbai555/bgg/service/lbbill/impl/dao/impl/mysql"
	"github.com/oldbai555/bgg/service/lbbill/impl/service"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/routine"
	"github.com/oldbai555/lbtool/utils"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"net"
)

func Run(ctx *cli.Context) error {
	conf.InitWebTool()

	mysql.RegisterModel([]interface{}{}...)
	if err := mysql.RegisterOrm(conf.Global.MysqlConf.Dsn()); err != nil {
		panic(fmt.Sprintf("err is %v", err))
	}

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", "0.0.0.0", conf.Global.ServerConf.Port))
	if err != nil {
		log.Fatalf("net.Listen err: %v", err)
	}
	// 新建gRPC服务器实例
	grpcServer := grpc.NewServer()
	// 在gRPC服务器注册我们的服务
	lbbill.RegisterLbbillServer(grpcServer, &service.BillServer)

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

	// 启动 grpc 服务
	if err := grpcServer.Serve(listener); err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}
