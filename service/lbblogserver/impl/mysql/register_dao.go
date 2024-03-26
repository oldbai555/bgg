package mysql

import (
	"github.com/oldbai555/bgg/pkg/syscfg"
	"github.com/oldbai555/bgg/service/lb"
	"github.com/oldbai555/bgg/service/lbblog"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/micro/bdb"
)

var (
	Article  *lb.Model
	Category *lb.Model
	Comment  *lb.Model
)

func RegisterOrm() (err error) {
	log.Infof("start init db orm......")
	mysqlConf := syscfg.NewGormMysqlConf("")

	// 注册表
	bdb.RegisterModel(
	// ...
	)

	err = bdb.InitMasterOrm(mysqlConf.Dsn())
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	bdb.AutoMigrate()
	Article = lb.NewModel(bdb.MasterOrm, &lbblog.ModelArticle{}, lbblog.ErrArticleNotFound)
	Category = lb.NewModel(bdb.MasterOrm, &lbblog.ModelCategory{}, lbblog.ErrCategoryNotFound)
	Comment = lb.NewModel(bdb.MasterOrm, &lbblog.ModelComment{}, lbblog.ErrCommentNotFound)

	log.Infof("end init db orm......")
	return
}
