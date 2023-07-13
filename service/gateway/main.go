package main

import (
	"context"
	"github.com/oldbai555/bgg/service/gateway/internal"
	"github.com/oldbai555/lbtool/log"
)

func main() {
	ctx := context.Background()
	err := internal.Server(ctx)
	if err != nil {
		log.Errorf("err is %v", err)
		return
	}
}
