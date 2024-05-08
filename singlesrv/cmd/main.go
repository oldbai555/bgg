package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/judwhite/go-svc"
	"github.com/oldbai555/bgg/pkg/ginhelper"
	"github.com/oldbai555/bgg/pkg/syscfg"
	"github.com/oldbai555/bgg/singlesrv/server/cache"
	"github.com/oldbai555/bgg/singlesrv/server/mysql"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/routine"
	"github.com/oldbai555/lbtool/utils"
	"os"
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
