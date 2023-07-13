package internal

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"

	"github.com/oldbai555/bgg/service/gateway/internal/conf"

	"github.com/oldbai555/bgg/pkg/gin_tool"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/lberr"
	"reflect"
	"strings"
)

func Server(_ context.Context) error {
	conf.InitWebTool()
	log.SetModuleName("gateway")

	gin.DefaultWriter = log.GetWriter()
	h := gin.New()

	// 配置跨域中间件
	// 链路追踪
	h.Use(Cors(), RegisterUuidTrace(), gin.LoggerWithFormatter(defaultLogFormatter), gin.Recovery(), RegisterJwt())

	h.POST("/gateway/*path", revProxy)

	// 启动服务
	err := h.Run(fmt.Sprintf(":%d", conf.Global.ServerConf.Port))
	if err != nil {
		log.Errorf("err is %v", err)
		return err
	}
	return nil
}

func revProxy(c *gin.Context) {
	// 找到对应的 rpc path
	var cm *gin_tool.Cmd
	handler := gin_tool.NewHandler(c)
	for _, cmd := range CmdList {
		if strings.HasSuffix(c.Request.RequestURI, cmd.Path) {
			cm = &cmd
			break
		}
	}

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

	handlerRet := v.Call([]reflect.Value{reflect.ValueOf(c), reqV})
	var callRes error
	if !handlerRet[1].IsNil() {
		callRes = handlerRet[1].Interface().(error)
	}
	if callRes != nil {
		log.Errorf("err:%v", callRes)
		handler.RespError(err)
		return
	}
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
	handler.RespError(lberr.NewInvalidArg("un ok"))
}
