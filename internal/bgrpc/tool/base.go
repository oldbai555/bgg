package tool

import (
	"fmt"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/dispatch"
	"github.com/oldbai555/lbtool/pkg/dispatch/impl/etcd"
	"github.com/oldbai555/lbtool/pkg/etcdcfg"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"sync"
	"time"
)

var onceD dispatch.IDispatch
var once sync.Once

var RoundRobinDialOpts = []grpc.DialOption{
	grpc.WithInsecure(),
	grpc.WithBlock(),
	grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"loadBalancingConfig": [{"%s":{}}]}`, roundrobin.Name)),
}

func New() (dispatch.IDispatch, error) {
	if onceD == nil {
		once.Do(func() {
			var err error
			onceD, err = etcd.NewDispatch(time.Second*5, clientv3.Config{
				Endpoints:   etcdcfg.GetConfig().GetEndpointList(),
				DialTimeout: time.Duration(etcdcfg.GetConfig().ConnectTimeoutMs) * time.Millisecond,
			})
			if err != nil {
				log.Errorf("err:%v", err)
				return
			}
		})
	}
	return onceD, nil
}
