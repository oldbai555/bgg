package main

import (
	"context"
	"github.com/oldbai555/bgg/pkg/syscfg"
	"github.com/oldbai555/bgg/service/lbbill"
	"github.com/oldbai555/bgg/service/lbbillserver/impl"
	"github.com/oldbai555/bgg/service/lbbillserver/impl/service"
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
	err = micro.NewGrpcWithGateSrv(lbbill.ServerName, conf.Ip, conf.Port,
		micro.WithRegisterFunc(func(server *grpc.Server) error {
			lbbill.RegisterLbbillServer(server, service.OnceSvrImpl)
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
