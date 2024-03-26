package mysql

import (
	"github.com/oldbai555/bgg/pkg/syscfg"
	"github.com/oldbai555/bgg/service/lb"
	"github.com/oldbai555/bgg/service/lbuser"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/micro/bdb"
)

var (
	User *lb.Model
)

func RegisterOrm() (err error) {
	log.Infof("start init db orm......")
	mysqlConf := syscfg.NewGormMysqlConf("")

	// 注册表
	bdb.RegisterModel(
		// ...
		&lbuser.ModelUser{},
	)

	err = bdb.InitMasterOrm(mysqlConf.Dsn())
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	bdb.AutoMigrate()
	User = lb.NewModel(bdb.MasterOrm, &lbuser.ModelUser{}, lbuser.ErrUserNotFound)

	log.Infof("end init db orm......")
	return
}
