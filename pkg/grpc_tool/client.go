package grpc_tool

import (
	"errors"
	"fmt"
	"github.com/oldbai555/lbtool/log"
	eclient "go.etcd.io/etcd/client/v3"
	eresolver "go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type ConnMgr struct {
	conn *grpc.ClientConn
	Srv  string
	Err  error
}

func NewConnMgr(srv string) *ConnMgr {
	// 创建 etcd 客户端
	config := GetConfig()
	etcdClient, _ := eclient.New(eclient.Config{
		Endpoints:   config.GetEndpointList(),
		DialTimeout: 5 * time.Second,
	})

	// 创建 etcd 实现的 grpc 服务注册发现模块 resolver
	etcdResolverBuilder, _ := eresolver.NewBuilder(etcdClient)

	// 拼接服务名称，需要固定义 etcd:/// 作为前缀
	etcdTarget := fmt.Sprintf("etcd:///%s", srv)

	// 创建 grpc 连接代理
	conn, err := grpc.Dial(
		// 服务名称
		etcdTarget,
		// 注入 etcd resolver
		grpc.WithResolvers(etcdResolverBuilder),
		// 声明使用的负载均衡策略为 round robin
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Errorf("err:%v", err)
	}

	return &ConnMgr{
		conn: conn,
		Srv:  srv,
		Err:  errors.New("init conn failed"),
	}
}

func (m *ConnMgr) GetConn() (*grpc.ClientConn, bool) {
	if m.conn == nil {
		return nil, false
	}
	return m.conn, true
}

func (m *ConnMgr) Close() error {
	if m.conn == nil {
		return nil
	}
	return m.conn.Close()
}
