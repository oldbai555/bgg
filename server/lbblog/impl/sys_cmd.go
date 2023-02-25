package impl

import (
	"fmt"
	"github.com/oldbai555/bgg/client/lbblog"
	webtool2 "github.com/oldbai555/bgg/pkg/webtool"
	"github.com/oldbai555/lbtool/log"
	"google.golang.org/grpc"
	"net"
)

var lb *Tool

type Tool struct {
	*webtool2.WebTool

	// 可以向下扩展其他的rpc服务
}

func StartServer() {
	var err error
	lb = &Tool{}
	lb.WebTool, err = webtool2.NewWebTool(webtool2.OptionWithOrm(), webtool2.OptionWithRdb(), webtool2.OptionWithStorage())

	if err != nil {
		log.Errorf("err:%v", err)
		return
	}

	// 初始化数据库工具类
	InitDbOrm()

	// 地址
	addr := fmt.Sprintf(":%d", lb.Sc.Port)
	// 1.监听
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Errorf("err is %v", err)
		return
	}
	log.Infof("监听端口：%s", addr)
	// 2.实例化gRPC
	s := grpc.NewServer()
	// 3.在gRPC上注册微服务
	lbblog.RegisterLbblogServer(s, &lbblogServer)
	// 4.启动服务端
	if err = s.Serve(listener); err != nil {
		log.Errorf("err is %v", err)
	}
}
