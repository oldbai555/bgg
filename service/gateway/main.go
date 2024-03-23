package main

import (
	"context"
	"github.com/oldbai555/lbtool/log"
)

func main() {
	ctx := context.Background()
	err := Server(ctx)
	if err != nil {
		log.Warnf("err is %v", err)
		return
	}
}
