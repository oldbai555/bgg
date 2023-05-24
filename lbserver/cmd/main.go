package main

import (
	"context"
	"github.com/oldbai555/bgg/lbserver/impl"
	"github.com/oldbai555/bgg/lbserver/impl/conf"
	"github.com/oldbai555/lbtool/log"
)

func main() {
	ctx := context.Background()
	conf.InitWebTool()
	err := impl.Server(ctx)
	if err != nil {
		log.Errorf("err is %v", err)
		return
	}
}
