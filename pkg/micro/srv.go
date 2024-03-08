package micro

import (
	"context"
	"fmt"
	"github.com/oldbai555/bgg/internal/_cmd"
	"github.com/oldbai555/bgg/internal/bgin/grpc_gate"
	"github.com/oldbai555/bgg/internal/bgrpc/register"
	"github.com/oldbai555/bgg/internal/bgrpc/srv"
	"github.com/oldbai555/bgg/pkg/prometheus_tool"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/routine"
	"google.golang.org/grpc"
)

// 聚合 grpc server 和 网关功能
// 可以单独拎出去 自己组装

type GrpcWithGateSrv struct {
	ip   string
	name string
	port uint32

	rf            srv.RegisterFunc
	checkAuthFunc grpc_gate.CheckAuthFunc
	cmdList       []*_cmd.Cmd
	interceptors  []grpc.UnaryServerInterceptor

	useDefaultSrvReg bool
}

func NewGrpcWithGateSrv(name, ip string, port uint32, opts ...Option) *GrpcWithGateSrv {
	s := &GrpcWithGateSrv{name: name, ip: ip, port: port}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

type Option func(*GrpcWithGateSrv)

func WithCheckAuthFunc(checkAuthFunc grpc_gate.CheckAuthFunc) Option {
	return func(gateSrv *GrpcWithGateSrv) {
		gateSrv.checkAuthFunc = checkAuthFunc
	}
}

func WithCmdList(cmdList []*_cmd.Cmd) Option {
	return func(gateSrv *GrpcWithGateSrv) {
		gateSrv.cmdList = cmdList
	}
}

func WithRegisterFunc(rf srv.RegisterFunc) Option {
	return func(gateSrv *GrpcWithGateSrv) {
		gateSrv.rf = rf
	}
}

func WithUnaryServerInterceptors(list ...grpc.UnaryServerInterceptor) Option {
	return func(gateSrv *GrpcWithGateSrv) {
		gateSrv.interceptors = list
	}
}

func WithUseDefaultSrvReg() Option {
	return func(gateSrv *GrpcWithGateSrv) {
		gateSrv.useDefaultSrvReg = true
	}
}

func (s *GrpcWithGateSrv) Start(ctx context.Context) error {
	grpcSrv := srv.NewSvr(s.name, s.port, s.rf, s.interceptors...)
	gateSrv := grpc_gate.NewSvr(s.name, s.genGatePort(), s.cmdList, s.checkAuthFunc)
	defer func() {
		grpcSrv.Stop()
		gateSrv.Stop()
	}()

	// 启动 grpc
	routine.GoV2(func() error {
		err := grpcSrv.StartGrpcSrv(ctx)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
		return nil
	})

	if s.useDefaultSrvReg {
		// 服务注册
		err := register.V2(ctx, s.ip, s.name, int(s.port), fmt.Sprintf("%d", s.genGatePort()))
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
	}

	// 启动监控
	routine.GoV2(func() error {
		err := prometheus_tool.StartPrometheusMonitor("", s.genPrometheusPort())
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
		return nil
	})

	// 启动 网关
	err := gateSrv.StartSrv(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}

func (s *GrpcWithGateSrv) genGatePort() uint32 {
	return s.port + 100
}

func (s *GrpcWithGateSrv) genPrometheusPort() uint32 {
	return s.port + 1000
}
