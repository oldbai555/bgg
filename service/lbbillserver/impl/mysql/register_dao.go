package mysql

import (
	"github.com/oldbai555/bgg/pkg/syscfg"
	"github.com/oldbai555/bgg/service/lb"
	"github.com/oldbai555/bgg/service/lbbill"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/micro/bdb"
)

var (
	Bill         *lb.Model
	BillCategory *lb.Model
)

func RegisterOrm() (err error) {
	log.Infof("start init db orm......")
	mysqlConf := syscfg.NewGormMysqlConf("")

	// 注册表
	bdb.RegisterModel(
		&lbbill.ModelBill{},
		&lbbill.ModelBillCategory{},
	)

	err = bdb.InitMasterOrm(mysqlConf.Dsn())
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	bdb.AutoMigrate()

	Bill = lb.NewModel(bdb.MasterOrm, &lbbill.ModelBill{}, lbbill.ErrBillNotFound)
	BillCategory = lb.NewModel(bdb.MasterOrm, &lbbill.ModelBillCategory{}, lbbill.ErrCategoryNotFound)

	log.Infof("end init db orm......")
	return
}
