package main

import (
	"context"
	"github.com/oldbai555/bgg/pkg/syscfg"
	"github.com/oldbai555/bgg/service/lbstore"
	"github.com/oldbai555/bgg/service/lbstoreserver/impl"
	"github.com/oldbai555/bgg/service/lbstoreserver/impl/service"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/utils"
	"github.com/oldbai555/micro"
	"google.golang.org/grpc"
)

func main() {
	syscfg.InitGlobal("", utils.GetCurDir(), syscfg.OptionWithServer())

	// 初始化中间件
	err := impl.InitMiddlewareComponent()
	if err != nil {
		log.Errorf("err:%v", err)
		panic(err)
	}

	// 初始化客户端
	err = impl.InitClient()
	if err != nil {
		log.Errorf("err:%v", err)
		panic(err)
	}

	conf, err := syscfg.GetServerConf()
	if err != nil {
		log.Errorf("err:%v", err)
		panic(err)
	}

	// 启动微服务
	err = micro.NewGrpcWithGateSrv(lbstore.ServerName, conf.Ip, conf.Port,
		micro.WithRegisterFunc(func(server *grpc.Server) error {
			lbstore.RegisterLbstoreServer(server, service.OnceSvrImpl)
			return nil
		}),
		micro.WithCmdList(cmdList),
		micro.WithUseDefaultSrvReg(),
		//micro.WithUnaryServerInterceptors(lbuser.NewCheckUserInterceptor(cmdList)),
		//micro.WithCheckAuthFunc(lbuser.CheckLoginUser),
	).Start(context.Background())
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
}
