package impl

import (
	"fmt"
	"github.com/oldbai555/bgg/client/lbuser"
	"github.com/oldbai555/bgg/pkg/webtool"
	"github.com/oldbai555/lbtool/log"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"net"
)

var lb *Tool

type Tool struct {
	*webtool.WebTool

	// 可以向下扩展其他的rpc服务
}

func StartServer() error {
	v, err := initViper()
	if err != nil {
		log.Errorf("err is %v", err)
		return err
	}

	lb = &Tool{}
	lb.WebTool, err = webtool.NewWebTool(v, webtool.OptionWithOrm(), webtool.OptionWithRdb(), webtool.OptionWithStorage())

	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	// 初始化数据库工具类
	InitDbOrm()

	// 地址
	addr := fmt.Sprintf(":%d", lb.Sc.Port)
	// 1.监听
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Errorf("err is %v", err)
		return err
	}
	log.Infof("监听端口：%s", addr)
	// 2.实例化gRPC
	s := grpc.NewServer()
	// 3.在gRPC上注册微服务
	lbuser.RegisterLbuserServer(s, &lbuserServer)
	// 4.启动服务端
	if err = s.Serve(listener); err != nil {
		log.Errorf("err is %v", err)
		return err
	}
	return nil
}

func initViper() (*viper.Viper, error) {
	viper.SetConfigName("application")          // name of config file (without extension)
	viper.SetConfigType("yaml")                 // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("/etc/work/")           // path to look for the config file in
	viper.AddConfigPath("./")                   // optionally look for config in the working directory
	viper.AddConfigPath("./server/lbuser/cmd/") // optionally look for config in the working directory
	err := viper.ReadInConfig()                 // Find and read the config file
	if err != nil {                             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	return viper.GetViper(), nil
}
