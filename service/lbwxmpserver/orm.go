package lbwxmpserver

import (
	"github.com/oldbai555/bgg/pkg/syscfg"
	"github.com/oldbai555/bgg/service/lbwxmp"
	"github.com/oldbai555/bgg/service/lbwxmpserver/autodb"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/micro/gormx"
	"github.com/oldbai555/micro/gormx/egimpl"
	"github.com/oldbai555/micro/gormx/engine"
)

var (
	OrmMpMemberUser         *gormx.BaseModel[*lbwxmp.ModelMpMemberUser]
	OrmMpStoreShop          *gormx.BaseModel[*lbwxmp.ModelMpStoreShop]
	OrmMpProductCategory    *gormx.BaseModel[*lbwxmp.ModelMpProductCategory]
	OrmMpStoreProduct       *gormx.BaseModel[*lbwxmp.ModelMpStoreProduct]
	OrmMpStoreOrderCartInfo *gormx.BaseModel[*lbwxmp.ModelMpStoreOrderCartInfo]
	OrmMpStoreOrder         *gormx.BaseModel[*lbwxmp.ModelMpStoreOrder]
	OrmMpService            *gormx.BaseModel[*lbwxmp.ModelMpService]
	OrmMpShopAds            *gormx.BaseModel[*lbwxmp.ModelMpShopAds]
	OrmMpStoreOrderStatus   *gormx.BaseModel[*lbwxmp.ModelMpStoreOrderStatus]
	OrmMpOrderNumber        *gormx.BaseModel[*lbwxmp.ModelMpOrderNumber]
)

func Init() (err error) {
	log.Infof("start init db orm......")
	mysqlConf := syscfg.NewGormMysqlConf("")
	gormEngine := egimpl.NewGormEngine(mysqlConf.Dsn())

	// 注入引擎
	engine.SetOrmEngine(gormEngine)

	// 注册表
	gormEngine.AutoMigrate([]interface{}{

		&lbwxmp.ModelMpMemberUser{},
		&lbwxmp.ModelMpStoreShop{},
		&lbwxmp.ModelMpProductCategory{},
		&lbwxmp.ModelMpStoreProduct{},
		&lbwxmp.ModelMpStoreOrderCartInfo{},
		&lbwxmp.ModelMpStoreOrder{},
		&lbwxmp.ModelMpService{},
		&lbwxmp.ModelMpShopAds{},
		&lbwxmp.ModelMpStoreOrderStatus{},
		&lbwxmp.ModelMpOrderNumber{},
	},
	)
	// 注入结构
	gormEngine.RegObjectType(

		autodb.ModelMpMemberUser,
		autodb.ModelMpStoreShop,
		autodb.ModelMpProductCategory,
		autodb.ModelMpStoreProduct,
		autodb.ModelMpStoreOrderCartInfo,
		autodb.ModelMpStoreOrder,
		autodb.ModelMpService,
		autodb.ModelMpShopAds,
		autodb.ModelMpStoreOrderStatus,
		autodb.ModelMpOrderNumber,
	)

	OrmMpMemberUser = gormx.NewBaseModel[*lbwxmp.ModelMpMemberUser](gormx.ModelConfig{
		NotFoundErrCode: int32(lbwxmp.ErrCode_ErrMpMemberUserNotFound),
		Db:              "biz",
	})
	OrmMpStoreShop = gormx.NewBaseModel[*lbwxmp.ModelMpStoreShop](gormx.ModelConfig{
		NotFoundErrCode: int32(lbwxmp.ErrCode_ErrMpStoreShopNotFound),
		Db:              "biz",
	})
	OrmMpProductCategory = gormx.NewBaseModel[*lbwxmp.ModelMpProductCategory](gormx.ModelConfig{
		NotFoundErrCode: int32(lbwxmp.ErrCode_ErrMpProductCategoryNotFound),
		Db:              "biz",
	})
	OrmMpStoreProduct = gormx.NewBaseModel[*lbwxmp.ModelMpStoreProduct](gormx.ModelConfig{
		NotFoundErrCode: int32(lbwxmp.ErrCode_ErrMpStoreProductNotFound),
		Db:              "biz",
	})
	OrmMpStoreOrderCartInfo = gormx.NewBaseModel[*lbwxmp.ModelMpStoreOrderCartInfo](gormx.ModelConfig{
		NotFoundErrCode: int32(lbwxmp.ErrCode_ErrMpStoreOrderCartInfoNotFound),
		Db:              "biz",
	})
	OrmMpStoreOrder = gormx.NewBaseModel[*lbwxmp.ModelMpStoreOrder](gormx.ModelConfig{
		NotFoundErrCode: int32(lbwxmp.ErrCode_ErrMpStoreOrderNotFound),
		Db:              "biz",
	})
	OrmMpService = gormx.NewBaseModel[*lbwxmp.ModelMpService](gormx.ModelConfig{
		NotFoundErrCode: int32(lbwxmp.ErrCode_ErrMpServiceNotFound),
		Db:              "biz",
	})
	OrmMpShopAds = gormx.NewBaseModel[*lbwxmp.ModelMpShopAds](gormx.ModelConfig{
		NotFoundErrCode: int32(lbwxmp.ErrCode_ErrMpShopAdsNotFound),
		Db:              "biz",
	})
	OrmMpStoreOrderStatus = gormx.NewBaseModel[*lbwxmp.ModelMpStoreOrderStatus](gormx.ModelConfig{
		NotFoundErrCode: int32(lbwxmp.ErrCode_ErrMpStoreOrderStatusNotFound),
		Db:              "biz",
	})
	OrmMpOrderNumber = gormx.NewBaseModel[*lbwxmp.ModelMpOrderNumber](gormx.ModelConfig{
		NotFoundErrCode: int32(lbwxmp.ErrCode_ErrMpOrderNumberNotFound),
		Db:              "biz",
	})

	log.Infof("end init db orm......")
	return
}
