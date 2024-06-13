package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/oldbai555/bgg/pkg/syscfg"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/dispatch"
	"github.com/oldbai555/lbtool/utils"
	"github.com/oldbai555/micro/bconst"
	"github.com/oldbai555/micro/bgin"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func Server(ctx context.Context) error {
	syscfg.InitGlobal("", utils.GetCurDir(), syscfg.OptionWithServer())
	srvName := syscfg.Global.ServerConf.Name
	log.SetModuleName(srvName)
	//d, err := dispatchimpl.New()
	//if err != nil {
	//	log.Errorf("err:%v", err)
	//	return err
	//}
	//return ginhelper.QuickStart(ctx, srvName, syscfg.Global.ServerConf.Port, func(router *gin.Engine) {
	//	router.POST("/gateway/*path", handleRevProxy(d))
	//})
	return nil
}

func handleRevProxy(d dispatch.IDispatch) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		handler := bgin.NewHandler(ctx)
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
}
