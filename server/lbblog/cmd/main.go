package main

import (
	"github.com/oldbai555/bgg/server/lbblog/impl"
	"github.com/oldbai555/lbtool/log"
)

func main() {
	err := impl.StartServer()
	if err != nil {
		log.Errorf("err is %v", err)
		return
	}
}
