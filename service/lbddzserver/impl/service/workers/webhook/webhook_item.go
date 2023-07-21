package webhook

import (
	"github.com/name5566/leaf/gate"
	"github.com/oldbai555/bgg/service/lbddz"
	iface2 "github.com/oldbai555/bgg/service/lbddz/iface"
)

type Item struct {
	Notifies []gate.Agent   // 通知列表
	Wh       *lbddz.Webhook // 通知内容
}

func NewWebhookItem(wh *lbddz.Webhook, notifies ...gate.Agent) iface2.IItem {
	return iface2.NewBaseItem(lbddz.WebhookConsumeTypeWebHook, &Item{
		Notifies: notifies,
		Wh:       wh,
	})
}
