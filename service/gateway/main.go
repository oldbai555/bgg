package main

/*
	desc:   很简单的网关
	author: bgg
	time:	2024/09/03
*/
import (
	"fmt"
	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth_gin"
	"github.com/gin-gonic/gin"
	"github.com/oldbai555/bgg/pkg/syscfg"
	"github.com/oldbai555/bgg/pkg/tool"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/utils"
	"github.com/oldbai555/micro/bconst"
	"github.com/oldbai555/micro/bgin"
	"github.com/oldbai555/micro/blimiter"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func main() {
	syscfg.InitGlobal("", utils.GetCurDir(), syscfg.OptionWithServer(), syscfg.OptionWithProxyConf())
	srvName := syscfg.Global.ServerConf.Name
	log.SetModuleName(srvName)
	gin.SetMode(gin.DebugMode)
	gin.DefaultWriter = log.GetWriter()
	gin.DefaultErrorWriter = log.GetWriter()
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Infof("%-6s %-25s --> %s (%d handlers)", httpMethod, absolutePath, handlerName, nuHandlers)
	}
	limiter := tollbooth.NewLimiter(blimiter.Max, blimiter.DefaultExpiredAbleOptions())
	engine := gin.New()
	engine.Use(
		gin.Recovery(),
		gin.LoggerWithConfig(gin.LoggerConfig{
			Formatter: tool.NewLogFormatter(srvName),
			Output:    log.GetWriter(),
		}),
		bgin.RegisterUuidTrace(),
		bgin.Cors(),
		tollbooth_gin.LimitHandler(limiter),
	)
	engine.GET("/gateway/*path", handleRevProxy)
	engine.POST("/gateway/*path", handleRevProxy)
	err := engine.Run(fmt.Sprintf(":%d", syscfg.Global.ServerConf.Port))
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
}

func handleRevProxy(ctx *gin.Context) {
	handler := bgin.NewHandler(ctx)
	param := ctx.Param("path")
	var srv = param
	paths := strings.Split(strings.TrimPrefix(param, "/"), "/")
	if len(paths) > 1 {
		srv = paths[0]
	}

	if syscfg.Global.Proxys == nil || syscfg.Global.Proxys.Map == nil {
		handler.Error(fmt.Errorf("未配置代理服务"))
		return
	}

	ipAddr, ok := syscfg.Global.Proxys.Map[srv]
	if !ok || len(ipAddr) == 0 {
		handler.Error(fmt.Errorf("未配置代理服务"))
		return
	}

	var target = ipAddr
	proxyUrl, err := url.Parse(target)
	if err != nil {
		log.Errorf("err:%v", err)
		handler.Error(err)
		return
	}

	// 重置 path
	ctx.Request.URL.Path = strings.Join(paths, "/")
	ctx.Request.Header.Set(bconst.ProtocolType, bconst.PROTO_TYPE_API_JSON)

	// todo 过滤一下请求
	proxy := httputil.NewSingleHostReverseProxy(proxyUrl)

	// todo 过滤一下响应
	proxy.ModifyResponse = func(resp *http.Response) error {
		resp.Header.Del(bconst.HeaderAccessControlAllowOrigin)
		resp.Header.Del(bconst.HeaderAccessControlAllowCredentials)
		return nil
	}

	proxy.ServeHTTP(ctx.Writer, ctx.Request)
}
