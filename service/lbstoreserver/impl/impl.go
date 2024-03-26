package impl

import (
	"context"
	"github.com/oldbai555/bgg/service/lbstore"
	"github.com/oldbai555/bgg/service/lbstoreserver/impl/cache"
	"github.com/oldbai555/bgg/service/lbstoreserver/impl/mysql"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/micro/brpc/bresolver"
	"google.golang.org/grpc/resolver"
)

func InitMiddlewareComponent() error {
	// 注册服务名
	log.SetModuleName(lbstore.ServerName)

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

	//err = discover.V2(ptuser.ServerName, ptuser.Init)
	//if err != nil {
	//	log.Errorf("err:%v", err)
	//	return err
	//}

	return nil
}
