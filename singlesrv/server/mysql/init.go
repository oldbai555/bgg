package mysql

import (
	"github.com/oldbai555/bgg/pkg/syscfg"
	"github.com/oldbai555/bgg/service/lb"
	"github.com/oldbai555/bgg/singlesrv/client"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/micro/bdb"
)

var (
	File *lb.Model
	User *lb.Model
)

func Init() (err error) {
	log.Infof("start init db orm......")
	mysqlConf := syscfg.NewGormMysqlConf("")

	// 注册表
	bdb.RegisterModel(

		&client.ModelFile{},
		&client.ModelUser{},
	)

	err = bdb.InitMasterOrm(mysqlConf.Dsn())
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	bdb.AutoMigrate()

	File = lb.NewModel(bdb.MasterOrm, &client.ModelFile{}, client.ErrFileNotFound)
	User = lb.NewModel(bdb.MasterOrm, &client.ModelUser{}, client.ErrUserNotFound)
	log.Infof("end init db orm......")
	return
}
