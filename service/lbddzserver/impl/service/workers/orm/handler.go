package orm

import (
	"context"
	"github.com/oldbai555/bgg/service/lbddz"
	"github.com/oldbai555/bgg/service/lbddzserver/impl/dao/impl/mysql"
	mgr2 "github.com/oldbai555/bgg/service/lbddzserver/impl/service/mgr"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/lberr"
	"github.com/oldbai555/lbtool/utils"
)

/*
在游戏服务器中，使用缓存和 MySQL 进行高性能存储的方法是将数据分为实时数据和非实时数据，并使用不同的存储方式进行处理。实时数据是指需要频繁读写的数据，如玩家位置、状态等，可以使用缓存来存储。而非实时数据是指不需要频繁读写的数据，如玩家装备、背包等，可以使用 MySQL 来存储。

具体实现步骤如下：

1. 将实时数据存储在 Redis 缓存中，通过异步方式将数据写入 MySQL 中。这样可以减少对 MySQL 的频繁读写，提高性能。同时，使用 Redis 可以提高读写速度和并发处理能力。

2. 将非实时数据存储在 MySQL 中，通过定时任务或者阈值触发的方式进行数据更新。定时任务可以定期更新数据，而阈值触发则是在数据量达到一定程度后进行更新。这样可以减少对 MySQL 的频繁读写，提高性能。

3. 使用事务处理和日志记录来确保数据的一致性和完整性。在更新数据时，使用事务处理可以确保数据的原子性和一致性。同时，使用日志记录可以记录每次数据更新的操作，以便出现问题时进行回溯和恢复。

通过以上步骤的实现，可以在不影响游戏性能的情况下，提高数据存储的性能和可靠性。
*/

func RegisterHandler() {
	mgr2.OrmHandlerMgr.Register(lbddz.OrmConsumeTypePlayerLogin, Login)
	mgr2.OrmHandlerMgr.Register(lbddz.OrmConsumeTypePlayerLogout, Logout)
	mgr2.OrmHandlerMgr.Register(lbddz.OrmConsumeTypeLoadPlayer, LoadPlayer)
	mgr2.OrmHandlerMgr.Register(lbddz.OrmConsumeTypeSyncGameData, SyncGameData)
}

// Login
// params[0] player.id
// params[1] cur_ip_addr
func Login(ctx context.Context, params ...interface{}) error {
	if len(params) < 2 {
		return lberr.NewInvalidArg("params[0] must player.id, params[1] must cur_ip_addr")
	}

	pId := params[0].(uint64)
	curIpAddr := params[1].(string)
	_, err := mysql.PlayerOrm.UpdateById(ctx, pId, map[string]interface{}{
		lbddz.DbCurIpAddr:   curIpAddr,
		lbddz.DbIsOnline:    true,
		lbddz.DbLastLoginAt: utils.TimeNow(),
	})
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

// Logout
// params[0] player.id
// params[1] last_ip_addr
func Logout(ctx context.Context, params ...interface{}) error {
	if len(params) < 2 {
		return lberr.NewInvalidArg("params[0] must player.id, params[1] must last_ip_addr")
	}

	pId := params[0].(uint64)
	lastIpAddr := params[1].(string)
	_, err := mysql.PlayerOrm.UpdateById(ctx, pId, map[string]interface{}{
		lbddz.DbLastIpAddr: lastIpAddr,
		lbddz.DbIsOnline:   false,
	})
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

// LoadPlayer
// params[0] player.id
func LoadPlayer(ctx context.Context, params ...interface{}) error {
	if len(params) < 1 {
		return lberr.NewInvalidArg("params[0] must player.id")
	}
	pId := params[0].(uint64)
	player, err := mysql.PlayerOrm.GetById(ctx, pId)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	mgr2.PlayerMgr.LoadPlayerInfo(player)
	return nil
}

// SyncGameData
// params[0] *lbddz.BaseGame
func SyncGameData(ctx context.Context, params ...interface{}) error {
	if len(params) < 1 {
		return lberr.NewInvalidArg("params[0] must *lbddzserver.BaseGame")
	}

	baseGame := params[0].(*lbddz.BaseGame)
	game := baseGame.G
	gamePlayers := baseGame.Gps
	_, err := mysql.GameOrm.UpdateById(ctx, game.Id, utils.OrmStruct2Map(game))
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	for i := range gamePlayers {
		gamePlayer := gamePlayers[i]
		_, err := mysql.GamePlayerOrm.UpdateById(ctx, gamePlayer.Id, utils.OrmStruct2Map(gamePlayer))
		if err != nil {
			log.Errorf("err:%v", err)
		}
	}

	// todo
	// mgr2.GameMgr.Set(baseGame)
	return nil
}
