package internal

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/oldbai555/bgg/pkg/_const"
	"github.com/oldbai555/bgg/pkg/cmd"
	"github.com/oldbai555/lbtool/pkg/signal"
	"google.golang.org/grpc/metadata"
	"net/http"
	"os"
	"time"

	"github.com/oldbai555/bgg/service/gateway/internal/conf"

	"github.com/oldbai555/bgg/pkg/gin_tool"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/lberr"
	"reflect"
	"strings"
)

func Server(ctx context.Context) error {
	conf.InitWebTool()
	log.SetModuleName("gateway")

	gin.DefaultWriter = log.GetWriter()
	router := gin.Default()

	// 配置跨域中间件
	// 链路追踪
	router.Use(
		gin_tool.Cors(),
		gin_tool.RegisterUuidTrace(),
		RegisterSvr(),
		gin.LoggerWithFormatter(gin_tool.DefaultLogFormatter),
		gin.Recovery(),
		RegisterJwt(),
	)

	router.POST("/gateway/*path", revProxy)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", conf.Global.ServerConf.Port),
		Handler: router,
	}

	signal.Reg(func(signal os.Signal) error {
		log.Warnf("exit: close gatewway server connect , signal[%v]", signal)
		err := srv.Shutdown(ctx)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
		return nil
	})
	signal.Do()

	// 启动服务
	err := srv.ListenAndServe()
	if err != nil {
		log.Warnf("err is %v", err)
		return err
	}
	return nil
}

func revProxy(c *gin.Context) {
	// 找到对应的 rpc path
	var cm *cmd.Cmd
	handler := gin_tool.NewHandler(c)
	for _, cc := range CmdList {
		if strings.HasSuffix(c.Request.RequestURI, cc.Path) {
			cm = &cc
			break
		}
	}

	// 校验方法
	h := cm.GRpcFunc
	v := reflect.ValueOf(h)
	t := v.Type()
	if !t.In(0).Implements(reflect.TypeOf((*context.Context)(nil)).Elem()) {
		panic("XX(context.Context, proto.Message)(proto.Message, error): first in arg must be context.Context")
	}
	if !t.Out(0).Implements(reflect.TypeOf((*proto.Message)(nil)).Elem()) {
		panic("XX(context.Context, proto.Message)(proto.Message, error): first out arg must be proto.Message")
	}
	if t.Out(1).String() != "error" {
		panic("XX(context.Context, proto.Message)(proto.Message, error): second out arg must be error")
	}

	// 拼装 request
	reqT := t.In(1).Elem()
	reqV := reflect.New(reqT)
	msg := reqV.Interface().(proto.Message)
	unmarshaler := &jsonpb.Unmarshaler{AllowUnknownFields: true}
	err := unmarshaler.Unmarshal(c.Request.Body, msg)
	if err != nil {
		log.Errorf("err:%v", err)
		handler.RespError(err)
		return
	}
	log.Infof("req is %v", msg.String())

	// 携带一些参数过去
	header := metadata.New(map[string]string{
		_const.HeaderLBTraceId:  c.GetHeader(_const.HeaderLBTraceId),
		_const.HeaderLBDeviceId: c.GetHeader(_const.HeaderLBDeviceId),
		_const.HeaderLBSid:      c.GetHeader(_const.HeaderLBSid),
		_const.HeaderLBCallFrom: c.GetHeader(_const.HeaderLBCallFrom),
	})

	// 执行方法 - 得到结果 - 最长等待 5 s
	newCtx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()
	var ctx = metadata.NewOutgoingContext(newCtx, header)
	handlerRet := v.Call([]reflect.Value{reflect.ValueOf(ctx), reqV})

	// 检查是否有误
	var callRes error
	if !handlerRet[1].IsNil() {
		callRes = handlerRet[1].Interface().(error)
	}
	if callRes != nil {
		log.Errorf("err:%v", callRes)
		handler.RespError(callRes)
		return
	}

	// 检查返回值
	if handlerRet[0].IsValid() && !handlerRet[0].IsNil() {
		rspBody, ok := handlerRet[0].Interface().(proto.Message)
		if !ok {
			log.Errorf("proto.Marshal err %v", err)
			handler.RespError(lberr.NewErr(500, "not proto.Message"))
			return
		}
		handler.Success(rspBody)
		return
	}

	// 走到这里说明走不动了
	handler.RespError(lberr.NewInvalidArg("un ok"))
}
