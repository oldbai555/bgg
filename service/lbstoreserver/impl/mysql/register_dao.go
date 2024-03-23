package mysql

import (
	"github.com/oldbai555/bgg/internal/bgorm"
	"github.com/oldbai555/bgg/pkg/syscfg"
	"github.com/oldbai555/bgg/service/lbstore"
	"github.com/oldbai555/lbtool/log"
)

var File *bgorm.Model

func RegisterOrm() (err error) {
	log.Infof("start init db orm......")
	mysqlConf := syscfg.NewGormMysqlConf("")

	// 注册表
	bgorm.RegisterModel(
		lbstore.ModelFile{},
	)

	err = bgorm.InitMasterOrm(mysqlConf.Dsn())
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	bgorm.AutoMigrate()

	File = bgorm.NewModel(bgorm.MasterOrm, &lbstore.ModelFile{}, lbstore.ErrFileNotFound)

	log.Infof("end init db orm......")
	return
}
