package internal

import (
	"github.com/oldbai555/lbtool/pkg/lberr"
	"reflect"
)

func handleMsg(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func init() {
	//handleMsg(lbddz.OrmConsumeTypePlayerLogin, handleLogin)
	//handleMsg(lbddz.OrmConsumeTypePlayerLogout, handleLogout)
	//handleMsg(lbddz.OrmConsumeTypeLoadPlayer, handleLoadPlayer)
	//handleMsg(lbddz.OrmConsumeTypeSyncGameData, handleSyncGameData)
}

func handleLogin(args []interface{}) interface{} {
	return lberr.NewErr(1, "2")
}

func handleLogout(args []interface{}) {

}

func handleLoadPlayer(args []interface{}) {

}

func handleSyncGameData(args []interface{}) {

}
