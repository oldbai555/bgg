package impl

import (
	"context"
	"github.com/oldbai555/bgg/internal/bgrpc/bresolver"
	"github.com/oldbai555/bgg/internal/bgrpc/discover"
	"github.com/oldbai555/bgg/service/lbblog"
	"github.com/oldbai555/bgg/service/lbblogserver/impl/cache"
	"github.com/oldbai555/bgg/service/lbblogserver/impl/mysql"
	"github.com/oldbai555/bgg/service/ptuser"
	"github.com/oldbai555/lbtool/log"
	"google.golang.org/grpc/resolver"
)

func InitMiddlewareComponent() error {
	// 注册服务名
	log.SetModuleName(lbblog.ServerName)

	// etcd 路径
	// etcdcfg.SetConfigPath("")

	// 初始化mysql
	err := mysql.RegisterOrm()
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	// 初始化redis
	err = cache.InitCache()
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

func InitClient() error {
	builder, err := bresolver.NewBuilder(context.Background())
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	resolver.Register(builder)

	err = discover.V2(ptuser.ServerName, ptuser.Init)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}
