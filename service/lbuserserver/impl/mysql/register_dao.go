package mysql

import (
	"github.com/oldbai555/bgg/internal/bgorm"
	"github.com/oldbai555/bgg/pkg/syscfg"
	"github.com/oldbai555/bgg/service/lbuser"
	"github.com/oldbai555/lbtool/log"
)

var (
	User *bgorm.Model
)

func RegisterOrm() (err error) {
	log.Infof("start init db orm......")
	mysqlConf := syscfg.NewGormMysqlConf("")

	// 注册表
	bgorm.RegisterModel(
		&lbuser.ModelUser{},
	)

	err = bgorm.InitMasterOrm(mysqlConf.Dsn())
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	bgorm.AutoMigrate()

	User = bgorm.NewModel(bgorm.MasterOrm, &lbuser.ModelUser{}, lbuser.ErrUserNotFound)

	log.Infof("end init db orm......")
	return
}
