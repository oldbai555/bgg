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
	OrmFile                     *gormx.BaseModel[*client.ModelFile]
	OrmUser                     *gormx.BaseModel[*client.ModelUser]
	OrmChat                     *gormx.BaseModel[*client.ModelChat]
	OrmMessage                  *gormx.BaseModel[*client.ModelMessage]
	OrmMpMerchantDetails        *gormx.BaseModel[*client.ModelMpMerchantDetails]
	OrmMpMemberUser             *gormx.BaseModel[*client.ModelMpMemberUser]
	OrmMpUserAddress            *gormx.BaseModel[*client.ModelMpUserAddress]
	OrmMpUserBill               *gormx.BaseModel[*client.ModelMpUserBill]
	OrmMpProductCategory        *gormx.BaseModel[*client.ModelMpProductCategory]
	OrmMpStoreProduct           *gormx.BaseModel[*client.ModelMpStoreProduct]
	OrmMpStoreProductAttr       *gormx.BaseModel[*client.ModelMpStoreProductAttr]
	OrmMpStoreProductAttrResult *gormx.BaseModel[*client.ModelMpStoreProductAttrResult]
	OrmMpStoreProductAttrValue  *gormx.BaseModel[*client.ModelMpStoreProductAttrValue]
	OrmMpStoreProductReply      *gormx.BaseModel[*client.ModelMpStoreProductReply]
	OrmMpStoreProductRule       *gormx.BaseModel[*client.ModelMpStoreProductRule]
	OrmMpStoreShop              *gormx.BaseModel[*client.ModelMpStoreShop]
	OrmMpCoupon                 *gormx.BaseModel[*client.ModelMpCoupon]
	OrmMpCouponUser             *gormx.BaseModel[*client.ModelMpCouponUser]
	OrmMpOrderNumber            *gormx.BaseModel[*client.ModelMpOrderNumber]
	OrmMpStoreOrder             *gormx.BaseModel[*client.ModelMpStoreOrder]
	OrmMpStoreOrderCartInfo     *gormx.BaseModel[*client.ModelMpStoreOrderCartInfo]
	OrmMpStoreOrderStatus       *gormx.BaseModel[*client.ModelMpStoreOrderStatus]
	OrmMpMaterial               *gormx.BaseModel[*client.ModelMpMaterial]
	OrmMpMaterialGroup          *gormx.BaseModel[*client.ModelMpMaterialGroup]
	OrmMpService                *gormx.BaseModel[*client.ModelMpService]
	OrmMpShopAds                *gormx.BaseModel[*client.ModelMpShopAds]
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
		&client.ModelMpMerchantDetails{},
		&client.ModelMpMemberUser{},
		&client.ModelMpUserAddress{},
		&client.ModelMpUserBill{},
		&client.ModelMpProductCategory{},
		&client.ModelMpStoreProduct{},
		&client.ModelMpStoreProductAttr{},
		&client.ModelMpStoreProductAttrResult{},
		&client.ModelMpStoreProductAttrValue{},
		&client.ModelMpStoreProductReply{},
		&client.ModelMpStoreProductRule{},
		&client.ModelMpStoreShop{},
		&client.ModelMpCoupon{},
		&client.ModelMpCouponUser{},
		&client.ModelMpOrderNumber{},
		&client.ModelMpStoreOrder{},
		&client.ModelMpStoreOrderCartInfo{},
		&client.ModelMpStoreOrderStatus{},
		&client.ModelMpMaterial{},
		&client.ModelMpMaterialGroup{},
		&client.ModelMpService{},
		&client.ModelMpShopAds{},
	},
	)
	// 注入结构
	gormEngine.RegObjectType(

		lbsingledb.ModelFile,
		lbsingledb.ModelUser,
		lbsingledb.ModelChat,
		lbsingledb.ModelMessage,
		lbsingledb.ModelMpMerchantDetails,
		lbsingledb.ModelMpMemberUser,
		lbsingledb.ModelMpUserAddress,
		lbsingledb.ModelMpUserBill,
		lbsingledb.ModelMpProductCategory,
		lbsingledb.ModelMpStoreProduct,
		lbsingledb.ModelMpStoreProductAttr,
		lbsingledb.ModelMpStoreProductAttrResult,
		lbsingledb.ModelMpStoreProductAttrValue,
		lbsingledb.ModelMpStoreProductReply,
		lbsingledb.ModelMpStoreProductRule,
		lbsingledb.ModelMpStoreShop,
		lbsingledb.ModelMpCoupon,
		lbsingledb.ModelMpCouponUser,
		lbsingledb.ModelMpOrderNumber,
		lbsingledb.ModelMpStoreOrder,
		lbsingledb.ModelMpStoreOrderCartInfo,
		lbsingledb.ModelMpStoreOrderStatus,
		lbsingledb.ModelMpMaterial,
		lbsingledb.ModelMpMaterialGroup,
		lbsingledb.ModelMpService,
		lbsingledb.ModelMpShopAds,
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
	OrmMpMerchantDetails = gormx.NewBaseModel[*client.ModelMpMerchantDetails](gormx.ModelConfig{
		NotFoundErrCode: int32(client.ErrCode_ErrMpMerchantDetailsNotFound),
		Db:              "biz",
	})
	OrmMpMemberUser = gormx.NewBaseModel[*client.ModelMpMemberUser](gormx.ModelConfig{
		NotFoundErrCode: int32(client.ErrCode_ErrMpMemberUserNotFound),
		Db:              "biz",
	})
	OrmMpUserAddress = gormx.NewBaseModel[*client.ModelMpUserAddress](gormx.ModelConfig{
		NotFoundErrCode: int32(client.ErrCode_ErrMpUserAddressNotFound),
		Db:              "biz",
	})
	OrmMpUserBill = gormx.NewBaseModel[*client.ModelMpUserBill](gormx.ModelConfig{
		NotFoundErrCode: int32(client.ErrCode_ErrMpUserBillNotFound),
		Db:              "biz",
	})
	OrmMpProductCategory = gormx.NewBaseModel[*client.ModelMpProductCategory](gormx.ModelConfig{
		NotFoundErrCode: int32(client.ErrCode_ErrMpProductCategoryNotFound),
		Db:              "biz",
	})
	OrmMpStoreProduct = gormx.NewBaseModel[*client.ModelMpStoreProduct](gormx.ModelConfig{
		NotFoundErrCode: int32(client.ErrCode_ErrMpStoreProductNotFound),
		Db:              "biz",
	})
	OrmMpStoreProductAttr = gormx.NewBaseModel[*client.ModelMpStoreProductAttr](gormx.ModelConfig{
		NotFoundErrCode: int32(client.ErrCode_ErrMpStoreProductAttrNotFound),
		Db:              "biz",
	})
	OrmMpStoreProductAttrResult = gormx.NewBaseModel[*client.ModelMpStoreProductAttrResult](gormx.ModelConfig{
		NotFoundErrCode: int32(client.ErrCode_ErrMpStoreProductAttrResultNotFound),
		Db:              "biz",
	})
	OrmMpStoreProductAttrValue = gormx.NewBaseModel[*client.ModelMpStoreProductAttrValue](gormx.ModelConfig{
		NotFoundErrCode: int32(client.ErrCode_ErrMpStoreProductAttrValueNotFound),
		Db:              "biz",
	})
	OrmMpStoreProductReply = gormx.NewBaseModel[*client.ModelMpStoreProductReply](gormx.ModelConfig{
		NotFoundErrCode: int32(client.ErrCode_ErrMpStoreProductReplyNotFound),
		Db:              "biz",
	})
	OrmMpStoreProductRule = gormx.NewBaseModel[*client.ModelMpStoreProductRule](gormx.ModelConfig{
		NotFoundErrCode: int32(client.ErrCode_ErrMpStoreProductRuleNotFound),
		Db:              "biz",
	})
	OrmMpStoreShop = gormx.NewBaseModel[*client.ModelMpStoreShop](gormx.ModelConfig{
		NotFoundErrCode: int32(client.ErrCode_ErrMpStoreShopNotFound),
		Db:              "biz",
	})
	OrmMpCoupon = gormx.NewBaseModel[*client.ModelMpCoupon](gormx.ModelConfig{
		NotFoundErrCode: int32(client.ErrCode_ErrMpCouponNotFound),
		Db:              "biz",
	})
	OrmMpCouponUser = gormx.NewBaseModel[*client.ModelMpCouponUser](gormx.ModelConfig{
		NotFoundErrCode: int32(client.ErrCode_ErrMpCouponUserNotFound),
		Db:              "biz",
	})
	OrmMpOrderNumber = gormx.NewBaseModel[*client.ModelMpOrderNumber](gormx.ModelConfig{
		NotFoundErrCode: int32(client.ErrCode_ErrMpOrderNumberNotFound),
		Db:              "biz",
	})
	OrmMpStoreOrder = gormx.NewBaseModel[*client.ModelMpStoreOrder](gormx.ModelConfig{
		NotFoundErrCode: int32(client.ErrCode_ErrMpStoreOrderNotFound),
		Db:              "biz",
	})
	OrmMpStoreOrderCartInfo = gormx.NewBaseModel[*client.ModelMpStoreOrderCartInfo](gormx.ModelConfig{
		NotFoundErrCode: int32(client.ErrCode_ErrMpStoreOrderCartInfoNotFound),
		Db:              "biz",
	})
	OrmMpStoreOrderStatus = gormx.NewBaseModel[*client.ModelMpStoreOrderStatus](gormx.ModelConfig{
		NotFoundErrCode: int32(client.ErrCode_ErrMpStoreOrderStatusNotFound),
		Db:              "biz",
	})
	OrmMpMaterial = gormx.NewBaseModel[*client.ModelMpMaterial](gormx.ModelConfig{
		NotFoundErrCode: int32(client.ErrCode_ErrMpMaterialNotFound),
		Db:              "biz",
	})
	OrmMpMaterialGroup = gormx.NewBaseModel[*client.ModelMpMaterialGroup](gormx.ModelConfig{
		NotFoundErrCode: int32(client.ErrCode_ErrMpMaterialGroupNotFound),
		Db:              "biz",
	})
	OrmMpService = gormx.NewBaseModel[*client.ModelMpService](gormx.ModelConfig{
		NotFoundErrCode: int32(client.ErrCode_ErrMpServiceNotFound),
		Db:              "biz",
	})
	OrmMpShopAds = gormx.NewBaseModel[*client.ModelMpShopAds](gormx.ModelConfig{
		NotFoundErrCode: int32(client.ErrCode_ErrMpShopAdsNotFound),
		Db:              "biz",
	})

	log.Infof("end init db orm......")
	return
}
