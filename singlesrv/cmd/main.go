package main

import (
	"context"
	"fmt"
	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth_gin"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	"github.com/judwhite/go-svc"
	"github.com/oldbai555/bgg/pkg/syscfg"
	"github.com/oldbai555/bgg/pkg/tool"
	"github.com/oldbai555/bgg/singlesrv/server"
	"github.com/oldbai555/bgg/singlesrv/server/cache"
	"github.com/oldbai555/bgg/singlesrv/server/ctx"
	"github.com/oldbai555/bgg/singlesrv/server/mq"
	"github.com/oldbai555/bgg/singlesrv/server/mysql"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/jsonpb"
	"github.com/oldbai555/lbtool/pkg/lberr"
	"github.com/oldbai555/lbtool/pkg/routine"
	"github.com/oldbai555/lbtool/utils"
	"github.com/oldbai555/micro/bcmd"
	"github.com/oldbai555/micro/bgin"
	"github.com/oldbai555/micro/bgin/gate"
	"github.com/oldbai555/micro/blimiter"
	"github.com/oldbai555/micro/brpc/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"reflect"
	"sync"
	"syscall"
)

type program struct {
	once       sync.Once
	port       uint32
	srvName    string
	prometheus *http.Server
	ginSrv     *http.Server
}

func (p *program) Init(_ svc.Environment) error {
	syscfg.InitGlobal("", utils.GetCurDir(), syscfg.OptionWithServer())
	conf, err := syscfg.GetServerConf()
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	srvName := conf.Name
	p.srvName = srvName
	p.port = conf.Port
	log.SetModuleName(srvName)

	// 检查默认生成的路由是否有误
	gate.CheckCmdList(cmdList)

	// 初始化mysql
	err = mysql.Init()
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

func (p *program) Start() error {
	var err error
	routine.GoV2(func() error {
		err = p.startPrometheusMonitor()
		return nil
	})
	if err != nil {
		panic(err)
	}

	routine.GoV2(func() error {
		err = p.startGinHttpServer()
		return nil
	})
	if err != nil {
		panic(err)
	}

	err = server.InitTopic()
	if err != nil {
		panic(err)
	}

	err = mq.Start()
	if err != nil {
		panic(err)
	}

	err = server.SyncFileIndex()
	if err != nil {
		panic(err)
	}

	mysql.InitDefaultAccount()
	return nil
}

func (p *program) Stop() error {
	if p.prometheus != nil {
		log.Infof("stop prometheus")
		err := p.prometheus.Shutdown(context.Background())
		if err != nil {
			log.Errorf("stop prometheus err:%v", err)
		}
	}
	if p.ginSrv != nil {
		log.Infof("stop gin server")
		err := p.ginSrv.Shutdown(context.Background())
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
	}
	mq.Stop()
	return nil
}

func (p *program) registerCmd(r *gin.Engine, cmd *bcmd.Cmd) {
	r.POST(cmd.Path, func(c *gin.Context) {
		handler := bgin.NewHandler(c)
		// call
		h := cmd.GRpcFunc
		v := reflect.ValueOf(h)
		t := v.Type()

		// 拼装 request
		reqT := t.In(1).Elem()
		reqV := reflect.New(reqT)
		msg := reqV.Interface().(proto.Message)
		err := jsonpb.Unmarshal(c.Request.Body, msg)
		if err != nil {
			log.Errorf("err:%v", err)
			handler.Error(err)
			return
		}

		if validator, ok := reqV.Interface().(middleware.Validator); ok {
			err := validator.Validate()
			if err != nil {
				log.Errorf("err:%v", err)
				handler.Error(err)
				return
			}
		}

		nCtx := ctx.NewCtx(
			c,
			ctx.WithGinHeaderSid(c),
			ctx.WithGinHeaderAuthType(c, cmd),
		)

		// 需要校验
		if cmd.IsUserAuthType() {
			info, err := server.CheckAuth(nCtx)
			if err != nil {
				log.Errorf("err:%v", err)
				handler.Error(err)
				return
			}
			nCtx.SetExtInfo(info)
		}

		log.Infof("req:[%s]", msg.String())

		handlerRet := v.Call([]reflect.Value{reflect.ValueOf(nCtx), reqV})

		// 检查是否有误
		var callRes error
		if !handlerRet[1].IsNil() {
			callRes = handlerRet[1].Interface().(error)
		}

		if callRes != nil {
			log.Errorf("err:%v", callRes)
			handler.Error(callRes)
			return
		}

		// 检查返回值
		if handlerRet[0].IsValid() && !handlerRet[0].IsNil() {
			rspBody, ok := handlerRet[0].Interface().(proto.Message)
			if !ok {
				log.Errorf("proto convert failed")
				handler.Error(lberr.NewErr(http.StatusInternalServerError, "proto convert failed"))
				return
			}

			handler.Success(rspBody)
			return
		}

		// 走到这里说明走不动了
		handler.Error(lberr.NewErr(http.StatusInternalServerError, "Internal error"))
	})
}

func (p *program) startPrometheusMonitor() error {
	onePort := tool.GetOnePort()
	if onePort == 0 {
		log.Warnf("获取到无效端口,无法开启 Prometheus")
		return nil
	}
	srv := http.NewServeMux()
	srv.Handle("/metrics", promhttp.Handler())
	p.prometheus = &http.Server{
		Addr:    fmt.Sprintf(":%d", onePort),
		Handler: srv,
	}
	log.Infof("====== start prometheus monitor, port is %d ======", onePort)
	err := p.prometheus.ListenAndServe()
	if err != nil {
		log.Warnf("err:%v", err)
		return err
	}

	return nil
}

func (p *program) startGinHttpServer() error {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = log.GetWriter()
	gin.DefaultErrorWriter = log.GetWriter()
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Infof("%-6s %-25s --> %s (%d handlers)", httpMethod, absolutePath, handlerName, nuHandlers)
	}
	router := gin.New()

	limiter := tollbooth.NewLimiter(blimiter.Max, blimiter.DefaultExpiredAbleOptions())

	router.Use(
		gin.Recovery(),
		gin.LoggerWithConfig(gin.LoggerConfig{
			Formatter: tool.NewLogFormatter(p.srvName),
			Output:    log.GetWriter(),
		}),
		bgin.RegisterUuidTrace(),
		bgin.Cors(),
		tollbooth_gin.LimitHandler(limiter),
	)

	for _, cmd := range cmdList {
		p.registerCmd(router, cmd)
	}

	p.ginSrv = &http.Server{
		Addr:    fmt.Sprintf(":%d", p.port),
		Handler: router,
	}

	// 启动服务
	log.Infof("====== start gin %s server, port is %d ======", p.srvName, p.port)
	err := p.ginSrv.ListenAndServe()
	if err != nil {
		log.Warnf("err is %v", err)
		return err
	}
	return nil
}

func main() {
	prg := &program{}
	err := svc.Run(prg,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGKILL,
	)
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
}
