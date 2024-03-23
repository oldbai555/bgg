package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/oldbai555/bgg/internal/_const"
	tool2 "github.com/oldbai555/bgg/internal/bgin/tool"
	grpctool "github.com/oldbai555/bgg/internal/bgrpc/tool"
	limiter2 "github.com/oldbai555/bgg/pkg/limiter"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth_gin"
	"github.com/oldbai555/bgg/pkg/prometheus_tool"
	"github.com/oldbai555/bgg/pkg/syscfg"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/dispatch"
	"github.com/oldbai555/lbtool/pkg/routine"
	"github.com/oldbai555/lbtool/pkg/signal"
	"github.com/oldbai555/lbtool/utils"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

func Server(ctx context.Context) error {
	syscfg.InitGlobal("", utils.GetCurDir(), syscfg.OptionWithServer(), syscfg.OptionWithPrometheus())

	log.SetModuleName(syscfg.Global.ServerConf.Name)

	d, err := grpctool.New()
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	gin.DefaultWriter = log.GetWriter()
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Infof("%-6s %-25s --> %s (%d handlers)", httpMethod, absolutePath, handlerName, nuHandlers)
	}
	router := gin.Default()

	// Create a limiter struct.
	limiter := tollbooth.NewLimiter(limiter2.Max, limiter2.DefaultExpiredAbleOptions())

	router.Use(
		gin.Recovery(),
		gin.LoggerWithFormatter(tool2.NewLogFormatter(syscfg.Global.ServerConf.Name)),
		tool2.Cors(),
		tool2.RegisterUuidTrace(),
		tollbooth_gin.LimitHandler(limiter),
	)

	router.POST("/gateway/*path", handleRevProxy(d))

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", syscfg.Global.ServerConf.Port),
		Handler: router,
	}

	routine.GoV2(func() error {
		err := prometheus_tool.StartPrometheusMonitor(syscfg.Global.PrometheusConf.Ip, syscfg.Global.PrometheusConf.Port)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
		return nil
	})

	signal.RegV2(func(signal os.Signal) error {
		log.Warnf("exit: close gateway server connect , signal[%v]", signal)
		err := srv.Shutdown(ctx)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
		return nil
	})

	// 启动服务
	err = srv.ListenAndServe()
	if err != nil {
		log.Warnf("err is %v", err)
		return err
	}
	return nil
}

func handleRevProxy(d dispatch.IDispatch) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		handler := tool2.NewHandler(ctx)
		param := ctx.Param("path")
		var srv = param
		paths := strings.Split(strings.TrimPrefix(param, "/"), "/")
		if len(paths) == 2 || len(paths) == 1 {
			srv = paths[0]
		}

		node, err := dispatch.Route(ctx, d, srv)
		if err != nil {
			log.Errorf("err:%v", err)
			handler.Error(err)
			return
		}

		var target = fmt.Sprintf("%s://%s:%s", "http", node.Host, node.Extra)
		proxyUrl, err := url.Parse(target)
		if err != nil {
			log.Errorf("err:%v", err)
			handler.Error(err)
			return
		}

		// 重置 path
		ctx.Request.URL.Path = strings.Join(paths, "/")
		ctx.Request.Header.Set(_const.ProtocolType, _const.PROTO_TYPE_API_JSON)

		// todo 过滤一下请求
		proxy := httputil.NewSingleHostReverseProxy(proxyUrl)

		// todo 过滤一下响应
		proxy.ModifyResponse = func(resp *http.Response) error {
			resp.Header.Del(_const.HeaderAccessControlAllowOrigin)
			resp.Header.Del(_const.HeaderAccessControlAllowCredentials)
			return nil
		}

		proxy.ServeHTTP(ctx.Writer, ctx.Request)

	}
}
