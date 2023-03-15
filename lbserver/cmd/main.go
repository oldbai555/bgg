package main

import (
	"github.com/oldbai555/bgg/lbserver/impl"
	"github.com/oldbai555/bgg/lbserver/impl/conf"
	"github.com/oldbai555/lbtool/log"
)

func main() {
	conf.InitWebTool()
	err := impl.StartServer()
	if err != nil {
		log.Errorf("err is %v", err)
		return
	}
}
