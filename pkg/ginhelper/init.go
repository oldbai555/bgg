/**
 * @Author: zjj
 * @Date: 2024/5/7
 * @Desc:
**/

package ginhelper

import (
	"context"
	"fmt"
	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth_gin"
	"github.com/gin-gonic/gin"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/routine"
	"github.com/oldbai555/lbtool/pkg/signal"
	"github.com/oldbai555/micro/bgin"
	"github.com/oldbai555/micro/blimiter"
	"github.com/oldbai555/micro/bprometheus"
	"net"
	"net/http"
	"os"
)

func QuickStart(ctx context.Context, srvName string, port uint32, registerRouter func(router *gin.Engine)) error {

	gin.DefaultWriter = log.GetWriter()
	gin.DefaultErrorWriter = log.GetWriter()
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Infof("%-6s %-25s --> %s (%d handlers)", httpMethod, absolutePath, handlerName, nuHandlers)
	}
	router := gin.Default()

	// Create a limiter struct.
	limiter := tollbooth.NewLimiter(blimiter.Max, blimiter.DefaultExpiredAbleOptions())

	router.Use(
		gin.Recovery(),
		gin.LoggerWithFormatter(bgin.NewLogFormatter(srvName)),
		bgin.Cors(),
		bgin.RegisterUuidTrace(),
		tollbooth_gin.LimitHandler(limiter),
	)

	registerRouter(router)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}

	routine.GoV2(func() error {
		onePort := getOnePort()
		if onePort == 0 {
			log.Warnf("获取到无效端口,无法开启Prometheus")
			return nil
		}
		err := bprometheus.StartPrometheusMonitor("", onePort)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
		return nil
	})

	signal.RegV2(func(signal os.Signal) error {
		log.Warnf("exit: close %s server connect , signal[%v]", srvName, signal)
		err := srv.Shutdown(ctx)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
		return nil
	})

	// 启动服务
	log.Infof("====== start %s, port is %d ======", srvName, port)
	err := srv.ListenAndServe()
	if err != nil {
		log.Warnf("err is %v", err)
		return err
	}
	return nil
}

func getOnePort() uint32 {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		log.Errorf("获取端口失败:%s", err)
		return 0
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Errorf("监听端口失败:%s", err)
		return 0
	}
	err = l.Close()
	if err != nil {
		log.Errorf("结束监听端口失败:%s", err)
		return 0
	}
	onePort := l.Addr().(*net.TCPAddr).Port
	log.Infof("获取端口成功:%v", onePort)
	return uint32(onePort)
}
