package lbossserver

import (
	"github.com/oldbai555/bgg/pkg/syscfg"
	"github.com/oldbai555/bgg/service/lboss"
	"github.com/oldbai555/bgg/service/lbossserver/autodb"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/micro/gormx"
	"github.com/oldbai555/micro/gormx/egimpl"
	"github.com/oldbai555/micro/gormx/engine"
)

var (
	OrmFile *gormx.BaseModel[*lboss.ModelFile]
)

func Init() (err error) {
	log.Infof("start init db orm......")
	mysqlConf := syscfg.NewGormMysqlConf("")
	gormEngine := egimpl.NewGormEngine(mysqlConf.Dsn())

	// 注入引擎
	engine.SetOrmEngine(gormEngine)

	// 注册表
	gormEngine.AutoMigrate([]interface{}{

		&lboss.ModelFile{},
	},
	)
	// 注入结构
	gormEngine.RegObjectType(

		autodb.ModelFile,
	)

	OrmFile = gormx.NewBaseModel[*lboss.ModelFile](gormx.ModelConfig{
		NotFoundErrCode: int32(lboss.ErrCode_ErrFileNotFound),
		Db:              "biz",
	})

	log.Infof("end init db orm......")
	return
}
