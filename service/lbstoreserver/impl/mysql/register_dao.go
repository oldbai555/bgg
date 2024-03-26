package mysql

import (
	"github.com/oldbai555/bgg/pkg/syscfg"
	"github.com/oldbai555/bgg/service/lb"
	"github.com/oldbai555/bgg/service/lbstore"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/micro/bdb"
)

var File *lb.Model

func RegisterOrm() (err error) {
	log.Infof("start init db orm......")
	mysqlConf := syscfg.NewGormMysqlConf("")

	// 注册表
	bdb.RegisterModel(
		// ...
		&lbstore.ModelFile{},
	)

	err = bdb.InitMasterOrm(mysqlConf.Dsn())
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	bdb.AutoMigrate()
	File = lb.NewModel(bdb.MasterOrm, &lbstore.ModelFile{}, lbstore.ErrFileNotFound)

	log.Infof("end init db orm......")
	return
}
