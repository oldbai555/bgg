package register

import (
	"context"
	"fmt"
	"github.com/oldbai555/bgg/pkg/grpc_tool"
	"github.com/oldbai555/lbtool/log"
	eclient "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
	"time"
)

func EndPointToEtcd(ctx context.Context, addr, svrName string) {
	log.Infof("register svr is %s, address: %s", svrName, addr)
	// 创建 etcd 客户端
	config := grpc_tool.GetConfig()
	etcdClient, _ := eclient.New(eclient.Config{
		Endpoints:   config.GetEndpointList(),
		DialTimeout: 5 * time.Second,
	})
	etcdManager, _ := endpoints.NewManager(etcdClient, svrName)

	// 创建一个租约，每隔 10s 需要向 etcd 汇报一次心跳，证明当前节点仍然存活
	var ttl int64 = 10
	lease, _ := etcdClient.Grant(ctx, ttl)

	// 添加注册节点到 etcd 中，并且携带上租约 id
	_ = etcdManager.AddEndpoint(ctx, fmt.Sprintf("%s/%s", svrName, addr), endpoints.Endpoint{Addr: addr}, eclient.WithLease(lease.ID))

	// 每隔 5 s进行一次延续租约的动作
	for {
		select {
		case <-time.After(5 * time.Second):
			// 续约操作
			resp, _ := etcdClient.KeepAliveOnce(ctx, lease.ID)
			log.Infof("keep alive resp: %+v", resp)
		case <-ctx.Done():
			return
		}
	}
}
