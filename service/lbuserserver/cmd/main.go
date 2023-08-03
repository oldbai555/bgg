package main

import (
	"github.com/oldbai555/bgg/service/lbuserserver/impl"
	"github.com/oldbai555/lbtool/log"

	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "lbuser"
	app.Version = "v0.0.1"
	app.Description = "lbuserserver"
	app.Action = impl.Run
	err := app.Run(os.Args)
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
}