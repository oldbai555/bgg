package impl

import (
	"context"
	"fmt"
	"github.com/name5566/leaf"
	lconf "github.com/name5566/leaf/conf"
	"github.com/oldbai555/bgg/service/lbddz"
	"github.com/oldbai555/bgg/service/lbddzserver/impl/conf"
	"github.com/oldbai555/bgg/service/lbddzserver/impl/dao/impl/mysql"
	"github.com/oldbai555/bgg/service/lbddzserver/impl/moude"
	"github.com/oldbai555/bgg/service/lbddzserver/impl/service/workers"
	"github.com/oldbai555/bgg/service/lbddzserver/impl/service/workers/game"
	"github.com/oldbai555/bgg/service/lbddzserver/impl/service/workers/orm"
	"github.com/oldbai555/bgg/service/lbddzserver/impl/service/workers/webhook"
	"github.com/urfave/cli/v2"
)

func Run(ctx *cli.Context) error {
	conf.InitWebTool()

	lconf.LogLevel = conf.Server.LogLevel
	lconf.LogPath = conf.Server.LogPath
	lconf.LogFlag = conf.LogFlag
	lconf.ConsolePort = conf.Server.ConsolePort
	lconf.ProfilePath = conf.Server.ProfilePath

	RegisterWorker(ctx.Context)

	mysql.RegisterModel([]interface{}{
		&lbddz.ModelGame{},
		&lbddz.ModelGamePlayer{},
		&lbddz.ModelPlayer{},
		&lbddz.ModelRoom{},
	}...)
	if err := mysql.RegisterOrm(conf.Global.MysqlConf.Dsn()); err != nil {
		panic(fmt.Sprintf("err is %v", err))
	}

	leaf.Run(moude.Modules()...)
	return nil
}

func RegisterWorker(ctx context.Context) {
	game.RegisterHandler()
	webhook.RegisterHandler()
	orm.RegisterHandler()
	workers.RegisterWorker(ctx)
}
