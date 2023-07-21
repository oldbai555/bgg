package webhook

import (
	"context"
	"github.com/name5566/leaf/gate"
	"github.com/oldbai555/bgg/service/lbddz"
	"github.com/oldbai555/bgg/service/lbddzserver/impl/service/mgr"
	"github.com/oldbai555/bgg/service/lbddzserver/impl/service/workers"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/lberr"
	"github.com/pkg/errors"
)

func RegisterHandler() {
	mgr.WebhookHandlerMgr.Register(lbddz.WebhookConsumeTypeWebHook, SendWebhook)
}

func SendWebhook(ctx context.Context, params ...interface{}) error {
	item := params[0].(*Item)
	for _, agent := range item.Notifies {
		log.Infof("send webhook , agent local addr %s, remote addr %s ,wh is %+v", agent.LocalAddr(), agent.RemoteAddr(), item.Wh)
		agent.WriteMsg(item.Wh)
	}
	return nil
}

// 针对对应的逻辑封装的推送

func PubByException(err error, notifies ...gate.Agent) {
	// 获取根错误
	rootErr := errors.Cause(err)

	var code = -1
	var msg = err.Error()
	if e, ok := rootErr.(*lberr.LbErr); ok {
		code = int(e.Code())
		msg = e.Message()
	}

	workers.WebhookWorker.SendItem(NewWebhookItem(&lbddz.Webhook{
		Type: uint32(lbddz.Webhook_TypeException),
		Exception: &lbddz.Exception{
			Code:    int32(code),
			Message: msg,
		},
	}, notifies...))
}

func PubByRegister(player *lbddz.ModelPlayer, notifies ...gate.Agent) {
	workers.WebhookWorker.SendItem(NewWebhookItem(&lbddz.Webhook{
		Type: uint32(lbddz.Webhook_TypeRegisterResult),
		Register: &lbddz.RegisterResult{
			Player: player,
		},
	}, notifies...))
}

func PubByLogin(player *lbddz.ModelPlayer, notifies ...gate.Agent) {
	workers.WebhookWorker.SendItem(NewWebhookItem(&lbddz.Webhook{
		Type: uint32(lbddz.Webhook_TypeLoginResult),
		Login: &lbddz.LoginResult{
			Player: player,
		},
	}, notifies...))
}

func PubByMatchResult(room *lbddz.ModelRoom, players []*lbddz.ModelPlayer, notifies ...gate.Agent) {
	workers.WebhookWorker.SendItem(NewWebhookItem(&lbddz.Webhook{
		Type: uint32(lbddz.Webhook_TypeMatchResult),
		Match: &lbddz.MatchResult{
			Players: players,
			Room:    room,
		},
	}, notifies...))
}

func PubByGiveCard(giveCard *lbddz.GiveCard, notifies ...gate.Agent) {
	workers.WebhookWorker.SendItem(NewWebhookItem(&lbddz.Webhook{
		Type:     uint32(lbddz.Webhook_TypeGiveCard),
		GiveCard: giveCard,
	}, notifies...))
}

func PubByStateChange(stateChange *lbddz.StateChange, notifies ...gate.Agent) {
	workers.WebhookWorker.SendItem(NewWebhookItem(&lbddz.Webhook{
		Type:        uint32(lbddz.Webhook_TypeStateChange),
		StateChange: stateChange,
	}, notifies...))
}

func PubByAck(ackType uint32, seq uint32, notifies ...gate.Agent) {
	workers.WebhookWorker.SendItem(NewWebhookItem(&lbddz.Webhook{
		Type: ackType,
		AckWantLandlord: &lbddz.Ack{
			Seq: seq,
		},
	}, notifies...))
}

func PubByNotifyWantLandlordResult(res *lbddz.WantLandlordResult, notifies ...gate.Agent) {
	workers.WebhookWorker.SendItem(NewWebhookItem(&lbddz.Webhook{
		Type:               uint32(lbddz.Webhook_TypeWantLandlordResult),
		WantLandlordResult: res,
	}, notifies...))
}

func SendToRoomPlayers(webhook *lbddz.Webhook, notifies ...gate.Agent) {
	workers.WebhookWorker.SendItem(NewWebhookItem(webhook, notifies...))
}

func SendPassMsg(gameState, nextSeq uint32, notifies ...gate.Agent) {
	workers.WebhookWorker.SendItem(NewWebhookItem(&lbddz.Webhook{
		Type: uint32(lbddz.Webhook_TypeAckPlayCardOut),
		PlayCardOut: &lbddz.PlayCardOut{
			StateChange: gameState,
			NextSeq:     nextSeq,
			CurCard: &lbddz.CardSet{
				CardType: uint32(lbddz.CardType_CardTypePassCards),
				Cards:    []uint32{},
				Header:   0,
			},
		},
	}, notifies...))
}

func SendNextCardOut(gameState, nextSeq uint32, cardType, header uint32, cards []uint32, notifies ...gate.Agent) {
	workers.WebhookWorker.SendItem(NewWebhookItem(&lbddz.Webhook{
		Type: uint32(lbddz.Webhook_TypeAckPlayCardOut),
		PlayCardOut: &lbddz.PlayCardOut{
			StateChange: gameState,
			NextSeq:     nextSeq,
			CurCard: &lbddz.CardSet{
				CardType: cardType,
				Cards:    cards,
				Header:   header,
			},
		},
	}, notifies...))
}
