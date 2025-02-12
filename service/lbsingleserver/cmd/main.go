package main

import (
	"context"
	"embed"
	"fmt"
	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth_gin"
	"github.com/gin-gonic/gin"
	"github.com/judwhite/go-svc"
	"github.com/oldbai555/bgg/pkg/bctx"
	scan "github.com/oldbai555/bgg/pkg/osargs"
	"github.com/oldbai555/bgg/pkg/syscfg"
	"github.com/oldbai555/bgg/pkg/tool"
	"github.com/oldbai555/bgg/service/lbsingleserver"
	"github.com/oldbai555/bgg/service/lbsingleserver/cache"
	"github.com/oldbai555/bgg/service/lbsingleserver/mq"
	"github.com/oldbai555/bgg/service/lbsingleserver/wsmgr"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/lberr"
	"github.com/oldbai555/lbtool/pkg/routine"
	"github.com/oldbai555/lbtool/utils"
	"github.com/oldbai555/micro/bcmd"
	"github.com/oldbai555/micro/bgin"
	"github.com/oldbai555/micro/bgin/gate"
	"github.com/oldbai555/micro/blimiter"
	"github.com/oldbai555/micro/uctx"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/protobuf/proto"
	"net/http"
	"os"
	"path"
	"sync"
	"syscall"
)

//go:embed application.yaml
var configFile embed.FS

type program struct {
	once       sync.Once
	port       uint32
	srvName    string
	prometheus *http.Server
	ginSrv     *http.Server
}

func (p *program) Init(_ svc.Environment) error {
	syscfg.InitGlobal("", utils.GetCurDir(), syscfg.OptionWithServer(), syscfg.OptionWithWxMiniProgram(), syscfg.OptionWithDeepSeek())
	conf, err := syscfg.GetServerConf()
	if err != nil {
		return lberr.Wrap(err)
	}
	srvName := conf.Name
	p.srvName = srvName
	p.port = conf.Port
	log.SetModuleName(srvName)

	// 检查默认生成的路由是否有误
	gate.CheckCmdList(cmdList)

	// 初始化mysql
	err = lbsingleserver.Init()
	if err != nil {
		return lberr.Wrap(err)
	}

	// 初始化redis
	err = cache.InitCache()
	if err != nil {
		return lberr.Wrap(err)
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
		return lberr.Wrap(err)
	}

	routine.GoV2(func() error {
		err = p.startGinHttpServer()
		return nil
	})
	if err != nil {
		return lberr.Wrap(err)
	}

	err = lbsingleserver.InitTopic()
	if err != nil {
		return lberr.Wrap(err)
	}

	err = mq.Start()
	if err != nil {
		return lberr.Wrap(err)
	}

	err = lbsingleserver.SyncFileIndex(context.Background(), true)
	if err != nil {
		return lberr.Wrap(err)
	}
	return nil
}

func (p *program) Stop() error {
	wsmgr.CloseAllWsConn()
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
			return lberr.Wrap(err)
		}
	}
	mq.Stop()
	cache.Stop()
	return nil
}

func (p *program) registerCmd(r *gin.Engine, cmd *bcmd.Cmd) {
	cmd.WithGenIUCtx(func(ctx *gin.Context) uctx.IUCtx {
		return bctx.NewCtx(ctx, bctx.WithGinHeaderAuthorization(ctx), bctx.WithGinHeaderAuthType(ctx, cmd), bctx.WithClientIp(ctx))
	}).WithCheckAuthF(func(nCtx uctx.IUCtx) (extInfo interface{}, err error) {
		return lbsingleserver.CheckAuth(nCtx)
	}).WithHandleError(func(ctx *gin.Context, err error) {
		handler := bgin.NewHandler(ctx)
		handler.Error(err)
	}).WithHandleResult(func(ctx *gin.Context, result proto.Message) {
		handler := bgin.NewHandler(ctx)
		handler.Success(result)
	})
	r.POST(cmd.Path, cmd.GinPost)
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
		return lberr.Wrap(err)
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

	// proto 生成的路由
	for _, cmd := range cmdList {
		p.registerCmd(router, cmd)
	}

	// 注册自定义路由
	lbsingleserver.RegisterCustomRouter(router)

	p.ginSrv = &http.Server{
		Addr:    fmt.Sprintf(":%d", p.port),
		Handler: router,
	}

	// 启动服务
	log.Infof("====== start gin %s server, port is %d ======", p.srvName, p.port)
	err := p.ginSrv.ListenAndServe()
	if err != nil {
		return lberr.Wrap(err)
	}
	return nil
}

func initConfig() {
	configPath := scan.OptStrDefault("configPath", path.Join(utils.GetCurDir(), "application.yaml"))
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		configData, err := configFile.ReadFile("application.yaml")
		if err != nil {
			log.Errorf("err:%v", err)
			return
		}

		err = os.WriteFile(configPath, configData, 0644)
		if err != nil {
			log.Errorf("err:%v", err)
			return
		}
		log.Errorf("config file not found, created a new one, please restart")
		return
	}
}

func main() {
	initConfig()
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
