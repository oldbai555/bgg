package main

import (
	"gmicro/common"
	"gmicro/pkg/log"
	"gmicro/pkg/routine"
	"gmicro/pkg/rpc"
	"gmicro/pkg/syscfg"
	"gmicro/service/dbproxy"
	"gmicro/service/dbproxyserver/impl"
)

var (
	IsProd bool // 是否生产环境
)

func main() {
	// 加载日志
	log.InitLogger(log.WithAppName(dbproxy.ServerName), log.WithScreen(!IsProd), log.WithLevel(log.DebugLevel))
	defer func() {
		log.Flush()
	}()

	// 加载配制
	syscfg.InitGlobal(common.GetCurDir(), syscfg.OptionWithServer())
	srvAddr, err := syscfg.GetSrvAddr()
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}

	// 初始化状态
	err = impl.InitState()
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}

	// 启动服务
	routine.Run(func() {
		err = rpc.ServerRun(dbproxy.ServerName, srvAddr, cmdList)
		if err != nil {
			log.Errorf("err:%v", err)
			return
		}
	})
}
