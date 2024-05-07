package main

import (
	"context"
	"github.com/oldbai555/bgg/pkg/ginhelper"
	"github.com/oldbai555/bgg/pkg/syscfg"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/utils"
)

func main() {
	ctx := context.Background()
	err := server(ctx)
	if err != nil {
		log.Warnf("err is %v", err)
		return
	}
}

func server(ctx context.Context) error {
	syscfg.InitGlobal("", utils.GetCurDir(), syscfg.OptionWithServer())
	srvName := syscfg.Global.ServerConf.Name
	port := syscfg.Global.ServerConf.Port
	log.SetModuleName(srvName)
	return ginhelper.QuickStart(ctx, srvName, port, registerRouter)
}
