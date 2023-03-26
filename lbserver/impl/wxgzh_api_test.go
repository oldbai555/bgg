package impl

import (
	"context"
	"github.com/oldbai555/bgg/client/lbcustomer"
	"github.com/oldbai555/bgg/lbserver/impl/conf"
	"github.com/oldbai555/bgg/lbserver/impl/service"
	"github.com/oldbai555/lbtool/log"
	"testing"
)

func init() {
	initEnv()
}

func initEnv() {
	conf.InitWebTool()
	// 初始化数据库
	err := service.InitDao(context.Background(), conf.Global.MysqlConf.Dsn())
	if err != nil {
		log.Errorf("err:%v", err)
		panic(err)
	}
}

func TestWXMsgReceive(t *testing.T) {
	var customer lbcustomer.ModelCustomer
	err := service.Customer.FirstOrCreate(context.Background(), &customer, map[string]interface{}{
		lbcustomer.FieldSn_: "o_d7556oXQ9gnqV8X8CAByvi4ug8",
	})
	if err != nil {
		log.Errorf("err is %v", err)
		return
	}
}
