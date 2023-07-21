package workers

import (
	"context"
	"github.com/oldbai555/bgg/service/lbddz/iface"
	"github.com/oldbai555/bgg/service/lbddzserver/impl/service/mgr"
)

// 用于 worker间 通信

var (
	GameWorker    iface.IWorker
	WebhookWorker iface.IWorker
	OrmWorker     iface.IWorker
)

func RegisterWorker(ctx context.Context) {
	// 游戏逻辑工作者
	GameWorker = iface.NewBaseWorker(1024, "game")
	GameWorker.Start(ctx, mgr.GameHandlerMgr)
	// webhook 推送工作者
	WebhookWorker = iface.NewBaseWorker(1024, "webhook")
	WebhookWorker.Start(ctx, mgr.WebhookHandlerMgr)
	// 数据同步工作者
	OrmWorker = iface.NewBaseWorker(1024, "orm")
	OrmWorker.Start(ctx, mgr.OrmHandlerMgr)
}
