package impl

import (
	"fmt"
	"github.com/oldbai555/bgg/client/lbuser"
	"github.com/oldbai555/bgg/service/lbuser/impl/conf"
	"github.com/oldbai555/bgg/service/lbuser/impl/dao/impl/mysql"
	"github.com/oldbai555/bgg/service/lbuser/impl/service"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"log"
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
	lbuser.RegisterLbuserServer(grpcServer, &service.UserServer)

	err = grpcServer.Serve(listener)
	if err != nil {
		panic(err)
	}
	return nil
}
