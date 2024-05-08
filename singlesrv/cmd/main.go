package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	"github.com/judwhite/go-svc"
	"github.com/oldbai555/bgg/pkg/ginhelper"
	"github.com/oldbai555/bgg/pkg/syscfg"
	"github.com/oldbai555/bgg/singlesrv/server/cache"
	"github.com/oldbai555/bgg/singlesrv/server/mysql"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/jsonpb"
	"github.com/oldbai555/lbtool/pkg/lberr"
	"github.com/oldbai555/lbtool/pkg/routine"
	"github.com/oldbai555/lbtool/utils"
	"github.com/oldbai555/micro/bcmd"
	"github.com/oldbai555/micro/bgin"
	"github.com/oldbai555/micro/bgin/gate"
	"net/http"
	"os"
	"reflect"
	"sync"
	"syscall"
)

type program struct {
	once    sync.Once
	port    uint32
	srvName string
}

func (p *program) Init(env svc.Environment) error {
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
	routine.GoV2(func() error {
		err := ginhelper.QuickStart(context.Background(), p.srvName, p.port, func(router *gin.Engine) {
			for _, cmd := range cmdList {
				p.registerCmd(router, cmd)
			}
		})
		if err != nil {
			log.Errorf("err:%v", err)
			err = p.Stop()
			os.Exit(1)
			return err
		}
		return nil
	})
	return nil
}

func (p *program) Stop() error {
	log.Infof("stop srv %s", p.srvName)
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
		log.Infof("req:[%s]", msg.String())

		handlerRet := v.Call([]reflect.Value{reflect.ValueOf(c), reqV})

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
	}
}
