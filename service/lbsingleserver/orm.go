package lbsingleserver

import (
	"github.com/oldbai555/bgg/pkg/syscfg"
	"github.com/oldbai555/bgg/service/lbsingle"
	"github.com/oldbai555/bgg/service/lbsingleserver/autodb"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/micro/gormx"
	"github.com/oldbai555/micro/gormx/egimpl"
	"github.com/oldbai555/micro/gormx/engine"
)

var (
	OrmFile    *gormx.BaseModel[*lbsingle.ModelFile]
	OrmUser    *gormx.BaseModel[*lbsingle.ModelUser]
	OrmChat    *gormx.BaseModel[*lbsingle.ModelChat]
	OrmMessage *gormx.BaseModel[*lbsingle.ModelMessage]
)

func Init() (err error) {
	log.Infof("start init db orm......")
	mysqlConf := syscfg.NewGormMysqlConf("")
	gormEngine := egimpl.NewGormEngine(mysqlConf.Dsn())

	// 注入引擎
	engine.SetOrmEngine(gormEngine)

	// 注册表
	gormEngine.AutoMigrate([]interface{}{

		&lbsingle.ModelFile{},
		&lbsingle.ModelUser{},
		&lbsingle.ModelChat{},
		&lbsingle.ModelMessage{},
	},
	)
	// 注入结构
	gormEngine.RegObjectType(

		autodb.ModelFile,
		autodb.ModelUser,
		autodb.ModelChat,
		autodb.ModelMessage,
	)

	OrmFile = gormx.NewBaseModel[*lbsingle.ModelFile](gormx.ModelConfig{
		NotFoundErrCode: int32(lbsingle.ErrCode_ErrFileNotFound),
		Db:              "biz",
	})
	OrmUser = gormx.NewBaseModel[*lbsingle.ModelUser](gormx.ModelConfig{
		NotFoundErrCode: int32(lbsingle.ErrCode_ErrUserNotFound),
		Db:              "biz",
	})
	OrmChat = gormx.NewBaseModel[*lbsingle.ModelChat](gormx.ModelConfig{
		NotFoundErrCode: int32(lbsingle.ErrCode_ErrChatNotFound),
		Db:              "biz",
	})
	OrmMessage = gormx.NewBaseModel[*lbsingle.ModelMessage](gormx.ModelConfig{
		NotFoundErrCode: int32(lbsingle.ErrCode_ErrMessageNotFound),
		Db:              "biz",
	})

	log.Infof("end init db orm......")
	return
}
