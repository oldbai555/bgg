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
	OrmFile    *gormx.BaseModel[*client.ModelFile]
	OrmUser    *gormx.BaseModel[*client.ModelUser]
	OrmChat    *gormx.BaseModel[*client.ModelChat]
	OrmMessage *gormx.BaseModel[*client.ModelMessage]
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
		&client.ModelChat{},
		&client.ModelMessage{},
	},
	)
	// 注入结构
	gormEngine.RegObjectType(

		lbsingledb.ModelFile,
		lbsingledb.ModelUser,
		lbsingledb.ModelChat,
		lbsingledb.ModelMessage,
	)

	OrmFile = gormx.NewBaseModel[*client.ModelFile](gormx.ModelConfig{
		NotFoundErrCode: int32(client.ErrCode_ErrFileNotFound),
		Db:              "biz",
	})
	OrmUser = gormx.NewBaseModel[*client.ModelUser](gormx.ModelConfig{
		NotFoundErrCode: int32(client.ErrCode_ErrUserNotFound),
		Db:              "biz",
	})
	OrmChat = gormx.NewBaseModel[*client.ModelChat](gormx.ModelConfig{
		NotFoundErrCode: int32(client.ErrCode_ErrChatNotFound),
		Db:              "biz",
	})
	OrmMessage = gormx.NewBaseModel[*client.ModelMessage](gormx.ModelConfig{
		NotFoundErrCode: int32(client.ErrCode_ErrMessageNotFound),
		Db:              "biz",
	})

	log.Infof("end init db orm......")
	return
}
