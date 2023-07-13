package game

import (
	"context"
	"fmt"
	"github.com/name5566/leaf/gate"
	"github.com/oldbai555/bgg/client/lbddz"
	"github.com/oldbai555/bgg/service/lbddz/impl/dao/impl/mysql"
	mgr2 "github.com/oldbai555/bgg/service/lbddz/impl/service/mgr"
	"github.com/oldbai555/bgg/service/lbddz/impl/service/workers"
	"github.com/oldbai555/bgg/service/lbddz/impl/service/workers/webhook"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/lberr"
	"github.com/oldbai555/lbtool/utils"
	"time"
)

func RegisterHandler() {
	mgr2.GameHandlerMgr.Register(lbddz.ConsumeTypeMatch, Match)
	mgr2.GameHandlerMgr.Register(lbddz.ConsumeTypeStartGame, StartGame)
	mgr2.GameHandlerMgr.Register(lbddz.ConsumeTypeWantLandlord, WantLandlord)
	mgr2.GameHandlerMgr.Register(lbddz.ConsumeTypePlayCard, PlayCard)
}

// PlayCard
// params[0] gate.Agent
// params[1] gameId
// params[2] lbddz.PlayCardIn
func PlayCard(_ context.Context, params ...interface{}) error {
	a := params[0].(gate.Agent)
	gameId := params[1].(uint64)
	playCardIn := params[2].(lbddz.PlayCardIn)

	baseGame, gamePlayer, err := beforeGetBaseGameWithGamePlayer(a, gameId)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	game := baseGame.G
	curCard := playCardIn.GetCurCard()
	as := mgr2.AgentMgr.GetAgentList(baseGame.GetGamePlayerIds()...)

	if len(curCard.Cards) > 0 {
		// 判断是否符合出牌规则

		// 牌型检查
		cardType := pokerLogic.CalcPokerType(curCard.Cards)
		if uint32(cardType) != curCard.CardType { // 计算出的牌型与传过来的不匹配
			log.Warnf("计算出的类型为 %d,传过来的为 %d", cardType, curCard.CardType)
			// 告知玩家出牌失败
			webhook.PubByAck(uint32(lbddz.Webhook_TypeAckPlayCardFail), gamePlayer.Seq, a)
			return nil
		}

		// 头牌检查
		cardHeader := pokerLogic.CalcPokerHeader(curCard.Cards, lbddz.CardType(curCard.CardType))
		if uint32(cardHeader) != curCard.Header { // 计算出的头牌与传过来的不匹配
			log.Warnf("计算出的头牌为 %d,传过来的为 %d", cardHeader, curCard.Header)
			// 告知玩家出牌失败
			webhook.PubByAck(uint32(lbddz.Webhook_TypeAckPlayCardFail), gamePlayer.Seq, a)
			return nil
		}

		// 是否可出牌检查
		if !pokerLogic.CanOut(curCard, game.LastCards) {
			log.Warnf("当前牌型为 %d, 头牌为 %d", curCard.CardType, game.LastCards.Header)
			log.Warnf("新来的牌型为 %d, 头牌为 %d", curCard.CardType, curCard.Header)
			webhook.PubByAck(uint32(lbddz.Webhook_TypeAckPlayCardFail), gamePlayer.Seq, a)
			return nil
		}

		// 移除手中的牌
		canRemove, hasOut := gamePlayer.RemoveCards(curCard.Cards)
		if !canRemove {
			// 告知玩家出牌失败
			webhook.PubByAck(uint32(lbddz.Webhook_TypeAckPlayCardFail), gamePlayer.Seq, a)
			return nil
		}

		// 出完牌了
		if hasOut {
			game.State = uint32(lbddz.GameStateChange_GameStateChangeGameOver)
			// 通知一下游戏结束
			webhook.PubByStateChange(&lbddz.StateChange{
				StateChange: game.State,
			}, as...)
		}
	}

	// 发送ack
	webhook.PubByAck(uint32(lbddz.Webhook_TypeAckPlayCard), gamePlayer.Seq, a)

	if curCard.CardType == uint32(lbddz.CardType_CardTypePassCards) { // 如果是过牌
		game.PassNum++
		game.NowBiggerSeq = game.LastPlayerSeq
		game.AddCurIndex()
		webhook.SendPassMsg(game.State, game.CurPlayerSeq, as...)
		if game.PassNum == 2 {
			game.PassNum = 0
		}
	} else {
		// 处理新的最大牌
		game.LastCards = curCard
		game.PassNum = 0

		// 通知本轮出牌，以及下一个应出牌的玩家
		game.AddCurIndex()
		webhook.SendNextCardOut(game.State, game.CurPlayerSeq, curCard.CardType, curCard.Header, curCard.Cards, as...)
	}

	// 更新一下信息
	game.LastCards = curCard
	game.LastPlayerSeq = gamePlayer.Seq
	baseGame.ResetSave(game, baseGame.Gps)
	mgr2.GameMgr.Set(baseGame)
	workers.OrmWorker.Send(lbddz.OrmConsumeTypeSyncGameData, baseGame)
	return nil
}

// WantLandlord
// params[0] gate.Agent
// params[1] gameId
// params[2] score
func WantLandlord(_ context.Context, params ...interface{}) error {
	a := params[0].(gate.Agent)
	gameId := params[1].(uint64)
	score := params[2].(uint32)

	baseGame, gamePlayer, err := beforeGetBaseGameWithGamePlayer(a, gameId)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	game := baseGame.G
	gamePlayers := baseGame.Gps
	as := mgr2.AgentMgr.GetAgentList(baseGame.GetGamePlayerIds()...)

	// 叫地主次数加一
	game.WantDiZhuTimes++

	// 发送ack
	webhook.PubByAck(uint32(lbddz.Webhook_TypeAckWantLandlord), gamePlayer.Seq, a)

	if score == 3 { // 直接喊3分，成为地主
		game.CurLandlordScore = score
		game.LandlordSeq = gamePlayer.Seq

		gamePlayer.Cards = append(gamePlayer.Cards, game.LandlordCards...)
		gamePlayer.CurCards = append(gamePlayer.CurCards, game.LandlordCards...)

		// 告知地主结果
		webhook.PubByNotifyWantLandlordResult(&lbddz.WantLandlordResult{
			StateChange:   uint32(lbddz.GameStateChange_GameStateChangeWantLandlord),
			CurPlayerSeq:  game.CurPlayerSeq,
			NewScore:      game.CurLandlordScore,
			LandlordCards: game.LandlordCards,
		}, as...)

		// 状态变更
		game.State = uint32(lbddz.GameStateChange_GameStateChangeGaming)
		webhook.PubByStateChange(&lbddz.StateChange{
			StateChange: game.State,
		}, as...)

		baseGame.ResetSave(game, gamePlayers)
		mgr2.GameMgr.Set(baseGame)
		workers.OrmWorker.Send(lbddz.OrmConsumeTypeSyncGameData, baseGame)

		return nil
	} else if score > game.CurLandlordScore { // 叫了个更高分，更新
		game.CurLandlordScore = score
		game.LandlordSeq = gamePlayer.Seq
	}

	// 如果是第三次，表示每个人都表过态了
	if game.WantDiZhuTimes == 3 {
		game.LandlordSeq = gamePlayer.Seq

		gamePlayer.Cards = append(gamePlayer.Cards, game.LandlordCards...)
		gamePlayer.CurCards = append(gamePlayer.CurCards, game.LandlordCards...)

		// 告知地主结果
		webhook.PubByNotifyWantLandlordResult(&lbddz.WantLandlordResult{
			StateChange:   uint32(lbddz.GameStateChange_GameStateChangeWantLandlord),
			CurPlayerSeq:  game.CurPlayerSeq,
			NewScore:      game.CurLandlordScore,
			LandlordCards: game.LandlordCards,
		}, as...)

		// 状态变更
		game.State = uint32(lbddz.GameStateChange_GameStateChangeGaming)
		webhook.PubByStateChange(&lbddz.StateChange{
			StateChange: game.State,
		}, as...)
	} else {
		// 继续问下一个人
		game.AddCurIndex()

		webhook.SendToRoomPlayers(&lbddz.Webhook{
			Type: uint32(lbddz.Webhook_TypeWantLandlordOutput),
			WantLandlordOutput: &lbddz.WantLandlordOutput{
				StateChange:  uint32(lbddz.GameStateChange_GameStateChangeWantLandlord),
				CurPlayerSeq: game.CurPlayerSeq,
				NewScore:     game.CurLandlordScore,
			},
		}, as...)
	}

	baseGame.ResetSave(game, gamePlayers)
	mgr2.GameMgr.Set(baseGame)
	workers.OrmWorker.Send(lbddz.OrmConsumeTypeSyncGameData, baseGame)
	return nil
}

// Match
// params[0] playerId
// params[1] gate.Agent
func Match(ctx context.Context, params ...interface{}) error {
	pId := params[0].(uint64)
	a := params[1].(gate.Agent)
	item := NewMatchItem(pId, a)

	// 进行排队
	matchQueue.Push(item)

	// 每三个进行处理
	if matchQueue.queue.Len() < 3 {
		return nil
	}

	var players []*lbddz.ModelPlayer
	var pItems []*MatchItem

	// 拿出元素
	for i := 0; i < 3; i++ {

		p, b := matchQueue.Pop()
		if !b {
			log.Warnf("获取玩家失败")
			continue
		}

		player, ok := mgr2.PlayerMgr.Get(p.PId)
		if !ok {
			log.Warnf("玩家 %d 离线", p.PId)
			continue
		}

		players = append(players, player)
		pItems = append(pItems, p)
	}

	// 兜底校验
	if len(pItems) != 3 {
		// 重新放回队列
		matchQueue.Push(pItems...)
		return nil
	}

	// 开始游戏
	// 创建房间
	r := &lbddz.ModelRoom{
		CreatorId: players[0].Id,
		Name:      fmt.Sprintf("房间%d", time.Now().UnixMilli()),
		PlayerIds: utils.PluckUint64List(players, lbddz.FieldId),
	}
	_, err := mysql.RoomOrm.Create(ctx, r)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	mgr2.RoomMgr.Set(r)

	// 发送匹配成功消息
	webhook.PubByMatchResult(r, players)

	workers.GameWorker.Send(lbddz.ConsumeTypeStartGame, r)
	return nil
}

// StartGame
// params[0] is *lbddz.ModelRoom
func StartGame(ctx context.Context, params ...interface{}) error {
	if len(params) < 1 {
		return lberr.NewInvalidArg("params[0] must is *lbddz.ModelGame")
	}

	r := params[0].(*lbddz.ModelRoom)

	// 创建对局
	game := &lbddz.ModelGame{
		RoomId:        r.Id,
		PlayerIds:     r.PlayerIds,
		LandlordCards: make([]uint32, 3),
	}
	_, err := mysql.GameOrm.Create(ctx, game)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	var gps []*lbddz.ModelGamePlayer
	for i := range r.PlayerIds {
		pId := r.PlayerIds[i]
		var gp = &lbddz.ModelGamePlayer{
			RoomId:   r.Id,
			GameId:   game.Id,
			PlayerId: pId,
			Seq:      uint32(i) + 1,
		}
		_, err := mysql.GamePlayerOrm.Create(ctx, gp)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
		gps = append(gps, gp)
	}

	// 发牌
	cards54 := lbddz.GetNewCards54()
	for i := 0; i < 17; i++ {
		gps[0].CurCards = append(gps[0].CurCards, cards54[i])
		gps[0].Cards = append(gps[0].Cards, cards54[i])

		gps[1].CurCards = append(gps[1].CurCards, cards54[i+17])
		gps[1].Cards = append(gps[1].Cards, cards54[i+17])

		gps[2].CurCards = append(gps[2].CurCards, cards54[i+34])
		gps[2].Cards = append(gps[2].Cards, cards54[i+34])
	}

	// 地主牌
	for i := 0; i < 3; i++ {
		game.LandlordCards[i] = cards54[i+51]
	}

	game.State = uint32(lbddz.GameStateChange_GameStateChangeWantLandlord)
	var baseGame = &lbddz.BaseGame{
		G:   game,
		Gps: gps,
	}

	// 更新一下玩家信息和对局信息
	workers.OrmWorker.Send(lbddz.OrmConsumeTypeSyncGameData, baseGame)

	// 更新一下内存里的值
	mgr2.GameMgr.Set(baseGame)

	// 通知发牌结束 到抢地主阶段了
	var ns []gate.Agent
	for i := 0; i < len(r.PlayerIds); i++ {
		ns = append(ns, mgr2.AgentMgr.GetAgent(r.PlayerIds[i]))
	}

	// 通知发牌
	webhook.PubByGiveCard(&lbddz.GiveCard{
		BaseGame: baseGame,
	}, ns...)

	// 通知游戏状态转换为抢地主
	webhook.PubByStateChange(&lbddz.StateChange{
		StateChange: game.State,
	}, ns...)

	return nil
}
