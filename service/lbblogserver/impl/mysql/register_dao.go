package mysql

import (
	"github.com/oldbai555/bgg/internal/bgorm"
	"github.com/oldbai555/bgg/pkg/syscfg"
	"github.com/oldbai555/bgg/service/lbblog"
	"github.com/oldbai555/lbtool/log"
)

var (
	Article  *bgorm.Model
	Category *bgorm.Model
	Comment  *bgorm.Model
)

func RegisterOrm() (err error) {
	log.Infof("start init db orm......")
	mysqlConf := syscfg.NewGormMysqlConf("")

	// 注册表
	bgorm.RegisterModel()

	err = bgorm.InitMasterOrm(mysqlConf.Dsn())
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	bgorm.AutoMigrate()

	Article = bgorm.NewModel(bgorm.MasterOrm, &lbblog.ModelArticle{}, lbblog.ErrArticleNotFound)
	Category = bgorm.NewModel(bgorm.MasterOrm, &lbblog.ModelCategory{}, lbblog.ErrCategoryNotFound)
	Comment = bgorm.NewModel(bgorm.MasterOrm, &lbblog.ModelComment{}, lbblog.ErrCommentNotFound)

	log.Infof("end init db orm......")
	return
}
