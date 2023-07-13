package internal

import (
	"context"
	"github.com/name5566/leaf/gate"
	"github.com/oldbai555/bgg/client/lbddz"
	"github.com/oldbai555/bgg/service/lbddz/impl/service/workers"
	"github.com/oldbai555/lbtool/log"
)

type handleFn func(ctx context.Context, a gate.Agent, e *lbddz.Event)

var fnMap = map[uint32]handleFn{
	uint32(lbddz.Event_TypeMatchPlayer):  matchPlayerFn,
	uint32(lbddz.Event_TypeWantLandlord): wantLandlordFn,
	uint32(lbddz.Event_TypePlayCardIn):   playCardIn,
}

func playCardIn(_ context.Context, a gate.Agent, e *lbddz.Event) {
	return
}

func matchPlayerFn(_ context.Context, a gate.Agent, e *lbddz.Event) {
	matchPlayer := e.GetMatchPlayer()
	log.Infof("matchPlayer is %v", matchPlayer)

	// 加入匹配队列
	workers.GameWorker.Send(lbddz.ConsumeTypeMatch, matchPlayer.PlayerId, a)
}

func wantLandlordFn(_ context.Context, a gate.Agent, e *lbddz.Event) {
	wantLandlord := e.GetWantLandlord()
	log.Infof("wantLandlord is %v", wantLandlord)

	workers.GameWorker.Send(lbddz.ConsumeTypeWantLandlord, a, wantLandlord.GameId, wantLandlord.Score)
}
