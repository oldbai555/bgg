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
	OrmUser                *gormx.BaseModel[*lbsingle.ModelUser]
	OrmChat                *gormx.BaseModel[*lbsingle.ModelChat]
	OrmMessage             *gormx.BaseModel[*lbsingle.ModelMessage]
	OrmDailyShortSentences *gormx.BaseModel[*lbsingle.ModelDailyShortSentences]
	OrmOutsideWebSite      *gormx.BaseModel[*lbsingle.ModelOutsideWebSite]
	OrmArticle             *gormx.BaseModel[*lbsingle.ModelArticle]
	OrmArticleTag          *gormx.BaseModel[*lbsingle.ModelArticleTag]
	OrmArticleCategory     *gormx.BaseModel[*lbsingle.ModelArticleCategory]
	OrmDictType            *gormx.BaseModel[*lbsingle.ModelDictType]
	OrmDictData            *gormx.BaseModel[*lbsingle.ModelDictData]
)

func Init() (err error) {
	log.Infof("start init db orm......")
	mysqlConf := syscfg.NewGormMysqlConf("")
	gormEngine := egimpl.NewGormEngine(mysqlConf.Dsn())

	// 注入引擎
	engine.SetOrmEngine(gormEngine)

	// 注册表
	gormEngine.AutoMigrate([]interface{}{

		&lbsingle.ModelUser{},
		&lbsingle.ModelChat{},
		&lbsingle.ModelMessage{},
		&lbsingle.ModelDailyShortSentences{},
		&lbsingle.ModelOutsideWebSite{},
		&lbsingle.ModelArticle{},
		&lbsingle.ModelArticleTag{},
		&lbsingle.ModelArticleCategory{},
		&lbsingle.ModelDictType{},
		&lbsingle.ModelDictData{},
	},
	)
	// 注入结构
	gormEngine.RegObjectType(

		autodb.ModelUser,
		autodb.ModelChat,
		autodb.ModelMessage,
		autodb.ModelDailyShortSentences,
		autodb.ModelOutsideWebSite,
		autodb.ModelArticle,
		autodb.ModelArticleTag,
		autodb.ModelArticleCategory,
		autodb.ModelDictType,
		autodb.ModelDictData,
	)

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
	OrmDailyShortSentences = gormx.NewBaseModel[*lbsingle.ModelDailyShortSentences](gormx.ModelConfig{
		NotFoundErrCode: int32(lbsingle.ErrCode_ErrDailyShortSentencesNotFound),
		Db:              "biz",
	})
	OrmOutsideWebSite = gormx.NewBaseModel[*lbsingle.ModelOutsideWebSite](gormx.ModelConfig{
		NotFoundErrCode: int32(lbsingle.ErrCode_ErrOutsideWebSiteNotFound),
		Db:              "biz",
	})
	OrmArticle = gormx.NewBaseModel[*lbsingle.ModelArticle](gormx.ModelConfig{
		NotFoundErrCode: int32(lbsingle.ErrCode_ErrArticleNotFound),
		Db:              "biz",
	})
	OrmArticleTag = gormx.NewBaseModel[*lbsingle.ModelArticleTag](gormx.ModelConfig{
		NotFoundErrCode: int32(lbsingle.ErrCode_ErrArticleTagNotFound),
		Db:              "biz",
	})
	OrmArticleCategory = gormx.NewBaseModel[*lbsingle.ModelArticleCategory](gormx.ModelConfig{
		NotFoundErrCode: int32(lbsingle.ErrCode_ErrArticleCategoryNotFound),
		Db:              "biz",
	})
	OrmDictType = gormx.NewBaseModel[*lbsingle.ModelDictType](gormx.ModelConfig{
		NotFoundErrCode: int32(lbsingle.ErrCode_ErrDictTypeNotFound),
		Db:              "biz",
	})
	OrmDictData = gormx.NewBaseModel[*lbsingle.ModelDictData](gormx.ModelConfig{
		NotFoundErrCode: int32(lbsingle.ErrCode_ErrDictDataNotFound),
		Db:              "biz",
	})

	log.Infof("end init db orm......")
	return
}
