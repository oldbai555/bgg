package mysql

import (
	"github.com/oldbai555/bgg/internal/bgorm"
	"github.com/oldbai555/bgg/pkg/syscfg"
	"github.com/oldbai555/bgg/service/lbbill"
	"github.com/oldbai555/lbtool/log"
)

var (
	Bill         *bgorm.Model
	BillCategory *bgorm.Model
)

func RegisterOrm() (err error) {
	log.Infof("start init db orm......")
	mysqlConf := syscfg.NewGormMysqlConf("")

	// 注册表
	bgorm.RegisterModel(
		&lbbill.ModelBill{},
		&lbbill.ModelBillCategory{},
	)

	err = bgorm.InitMasterOrm(mysqlConf.Dsn())
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	bgorm.AutoMigrate()

	Bill = bgorm.NewModel(bgorm.MasterOrm, &lbbill.ModelBill{}, lbbill.ErrBillNotFound)
	BillCategory = bgorm.NewModel(bgorm.MasterOrm, &lbbill.ModelBillCategory{}, lbbill.ErrCategoryNotFound)

	log.Infof("end init db orm......")
	return
}
