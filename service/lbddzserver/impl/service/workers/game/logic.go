package game

import (
	"github.com/name5566/leaf/gate"
	"github.com/oldbai555/bgg/service/lbddz"
	"github.com/oldbai555/bgg/service/lbddzserver/impl/service/mgr"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/lberr"
)

func beforeGetBaseGameWithGamePlayer(a gate.Agent, gameId uint64) (*lbddz.BaseGame, *lbddz.ModelGamePlayer, error) {
	playerId := mgr.AgentMgr.GetPlayerId(a)

	baseGame, ok := mgr.GameMgr.Get(gameId)
	// 拿不到对局信息
	if !ok {
		return nil, nil, lberr.NewInvalidArg("not found game , gameId is %d , playerId is %d", gameId, playerId)
	}

	game := baseGame.G
	gamePlayer, err := baseGame.GetGamePlayer(playerId)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, nil, err
	}

	// 还没轮到它
	if !game.CheckCurPlayerSeq(gamePlayer.Seq) {
		return nil, nil, lberr.NewInvalidArg("not take turns game player , playerId is %d,current player seq is %d,player seq is %d", playerId, game.CurPlayerSeq, gamePlayer.Seq)
	}

	return baseGame, gamePlayer, nil
}
