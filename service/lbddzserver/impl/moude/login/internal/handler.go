package internal

import (
	"context"
	"fmt"
	"github.com/name5566/leaf/gate"
	"github.com/oldbai555/bgg/service/lbddz"
	"github.com/oldbai555/bgg/service/lbddzserver/impl/dao/impl/mysql"
	"github.com/oldbai555/bgg/service/lbddzserver/impl/moude/db"
	"github.com/oldbai555/bgg/service/lbddzserver/impl/service/mgr"
	"github.com/oldbai555/bgg/service/lbddzserver/impl/service/workers"
	"github.com/oldbai555/bgg/service/lbddzserver/impl/service/workers/webhook"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/utils"
	"reflect"
)

func handleMsg(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func init() {
	handleMsg(&lbddz.Register{}, handleRegister)
	handleMsg(&lbddz.Login{}, handleLogin)
}

func handleRegister(args []interface{}) {
	m := args[0].(*lbddz.Register)
	a := args[1].(gate.Agent)
	log.Infof("register msg is %v", m)

	ctx := context.Background()
	var player lbddz.ModelPlayer

	result, err := mysql.PlayerOrm.FirstOrCreate(ctx, map[string]interface{}{
		lbddz.DbUsername: m.Username,
	}, &player)
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}

	// 已经被注册过
	if !result.Created {
		webhook.PubByException(lbddz.ErrAlreadyRegister, a)
		return
	}
	err = db.ChanRPC.Call0(lbddz.OrmConsumeTypePlayerLogin)
	if err != nil {
		log.Errorf("err:%v", err)
	}
	// 更新密码和名称
	_, err = mysql.PlayerOrm.UpdateOrCreate(ctx, map[string]interface{}{
		lbddz.DbId: player.Id,
	}, map[string]interface{}{
		lbddz.DbPassword: utils.Md5(m.Password),
		lbddz.DbNickname: fmt.Sprintf("用户%d", utils.TimeNow()),
	}, &player)
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}

	webhook.PubByRegister(&player, a)
}

func handleLogin(args []interface{}) {
	m := args[0].(*lbddz.Login)
	a := args[1].(gate.Agent)
	log.Infof("login msg is %v", m)
	ctx := context.Background()

	player, err := mysql.PlayerOrm.GetOne(ctx, map[string]interface{}{
		lbddz.DbUsername: m.Username,
	})
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}

	if player.Password != utils.Md5(m.Password) {
		webhook.PubByException(lbddz.ErrPasswordMistake, a)
		return
	}
	player.CurIpAddr = a.RemoteAddr().String()

	// 记录管理
	mgr.PlayerMgr.Set(player)
	mgr.AgentMgr.Set(a, player.Id)

	// 告知 db 去更新数据
	workers.OrmWorker.Send(lbddz.OrmConsumeTypePlayerLogin, player.Id, a.RemoteAddr().String())

	// 拉取一下完整的信息
	workers.OrmWorker.Send(lbddz.OrmConsumeTypeLoadPlayer, player.Id)

	webhook.PubByLogin(player, a)
}
