package server

import (
	"github.com/oldbai555/bgg/pkg/syscfg"
	"github.com/oldbai555/bgg/singlesrv/client"
	"github.com/oldbai555/bgg/singlesrv/client/lbsingledb"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/micro/gormx"
	"github.com/oldbai555/micro/gormx/egimpl"
	"github.com/oldbai555/micro/gormx/engine"
)

var (
	OrmFile *gormx.BaseModel[*client.ModelFile]
	OrmUser *gormx.BaseModel[*client.ModelUser]
)

func Init() (err error) {
	log.Infof("start init db orm......")
	mysqlConf := syscfg.NewGormMysqlConf("")
	gormEngine := egimpl.NewGormEngine(mysqlConf.Dsn())

	// 注入引擎
	engine.SetOrmEngine(gormEngine)

	// 注册表
	gormEngine.AutoMigrate([]interface{}{
		&client.ModelFile{},
		&client.ModelUser{},
	},
	)

	// 注入结构
	gormEngine.RegObjectType(
		lbsingledb.ModelFile,
		lbsingledb.ModelUser,
	)

	OrmFile = gormx.NewBaseModel[*client.ModelFile](gormx.ModelConfig{
		NotFoundErrCode: int32(client.ErrCode_ErrFileNotFound),
		Db:              "",
	})
	OrmUser = gormx.NewBaseModel[*client.ModelUser](gormx.ModelConfig{
		NotFoundErrCode: int32(client.ErrCode_ErrUserNotFound),
		Db:              "",
	})
	log.Infof("end init db orm......")
	return
}
