package lbwxmpserver

import (
	"context"
	"fmt"
	"github.com/oldbai555/bgg/pkg/bctx"
	"github.com/oldbai555/bgg/pkg/syscfg"
	"github.com/oldbai555/bgg/service/lbbase"
	"github.com/oldbai555/bgg/service/lbsingleserver/cache"
	"github.com/oldbai555/bgg/service/lbwxmp"
	"github.com/oldbai555/bgg/service/lbwxmpserver/wxminiprogram"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/utils"
	"github.com/oldbai555/micro/core"
	"github.com/oldbai555/micro/gormx"
	"github.com/oldbai555/micro/uctx"
)

func (a *LbwxmpServer) WxMPAuthSendSmsLogin(ctx context.Context, req *lbwxmp.WxMPAuthSendSmsLoginReq) (*lbwxmp.WxMPAuthSendSmsLoginRsp, error) {
	var rsp lbwxmp.WxMPAuthSendSmsLoginRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	ctx2, ok := uCtx.(*bctx.Ctx)
	if !ok {
		return nil, lbwxmp.ErrMpCtxConvertFailed
	}

	code, err := cache.GetMpSmsCode(fmt.Sprintf("%s", req.Mobile))
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	if code != req.Code {
		return nil, lbwxmp.ErrMpSmsCodeNotEqual
	}
	cache.DelMpSmsCode(fmt.Sprintf("%s", req.Mobile))

	mpUser := &lbwxmp.ModelMpMemberUser{
		Nickname:    "微信用户_" + utils.GenUUID(),
		Mobile:      req.Mobile,
		RegisterIp:  ctx2.ClientIp,
		LastLoginAt: utils.TimeNow(),
		LastLoginIp: ctx2.ClientIp,
		LoginType:   req.From,
		MpOpenid:    req.Openid,
	}

	isEmpty, err := OrmMpMemberUser.NewBaseScope().Where(lbwxmp.FieldMobile_, mpUser.Mobile).IsEmpty(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	if isEmpty {
		err = OrmMpMemberUser.NewBaseScope().Create(uCtx, mpUser)
	}
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	rsp.UserInfo = mpUser
	sid := utils.GenUUID()
	err = cache.SetLoginInfo(sid, rsp.UserInfo.ToBaseUser())
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.AccessToken = sid

	return &rsp, err
}
func (a *LbwxmpServer) WxMPAuthSendSmsCode(ctx context.Context, req *lbwxmp.WxMPAuthSendSmsCodeReq) (*lbwxmp.WxMPAuthSendSmsCodeRsp, error) {
	var rsp lbwxmp.WxMPAuthSendSmsCodeRsp
	var err error

	code, err := cache.GetMpSmsCode(fmt.Sprintf("%s", req.Mobile))
	if err == nil && code != "" {
		return nil, lbwxmp.ErrMpSmsCodeNoExpired
	}

	rsp.Code = utils.GetRandomString(6, utils.RandomStringModNumber)
	err = cache.SetMpSmsCode(fmt.Sprintf("%s", req.Mobile), rsp.Code)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}
func (a *LbwxmpServer) WxMiniProgramAuthSession(ctx context.Context, req *lbwxmp.WxMiniProgramAuthSessionReq) (*lbwxmp.WxMiniProgramAuthSessionRsp, error) {
	var rsp lbwxmp.WxMiniProgramAuthSessionRsp
	var err error

	conf, err := syscfg.GetWxMiniProgramConf()
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	session, err := wxminiprogram.Code2Session(&lbbase.JsCodeToSessionReq{
		JsCode: req.Code,
		Appid:  conf.AppId,
		Secret: conf.Secret,
	})
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	rsp.Openid = session.Openid
	rsp.UserInfo, err = OrmMpMemberUser.NewBaseScope().Where(lbwxmp.FieldMpOpenid_, session.Openid).First(bctx.NewCtx(ctx))
	if err != nil {
		log.Errorf("err:%v", err)
		err = nil
	}
	if rsp.UserInfo != nil {
		sid := utils.GenUUID()
		err = cache.SetLoginInfo(sid, rsp.UserInfo.ToBaseUser())
		if err != nil {
			log.Errorf("err:%v", err)
			return nil, err
		}
		rsp.AccessToken = sid
	}

	return &rsp, err
}
func (a *LbwxmpServer) MPShopNearBy(ctx context.Context, req *lbwxmp.MPShopNearByReq) (*lbwxmp.MPShopNearByRsp, error) {
	var rsp lbwxmp.MPShopNearByRsp
	var err error
	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	db := OrmMpStoreShop.NewBaseScope()
	if req.ShopId != 0 {
		db.Where(lbwxmp.FieldId_, req.ShopId)
	}
	if req.Kw != "" {
		db.WhereLike(lbwxmp.FieldName_, fmt.Sprintf("%%%s%%", req.Kw))
	}
	rsp.Shop, err = db.Find(uCtx)
	return &rsp, err
}
func (a *LbwxmpServer) MPShopProduct(ctx context.Context, req *lbwxmp.MPShopProductReq) (*lbwxmp.MPShopProductRsp, error) {
	var rsp lbwxmp.MPShopProductRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	db := OrmMpStoreShop.NewBaseScope()
	if req.ShopId != 0 {
		db.Where(lbwxmp.FieldId_, req.ShopId)
	}
	storeShop, err := db.First(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.List, err = OrmMpProductCategory.NewBaseScope().Where(lbwxmp.FieldMpStoreShopId_, storeShop.Id).Find(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	if len(rsp.List) == 0 {
		return &rsp, nil
	}

	rsp.ProductMap = make(map[uint64]*lbwxmp.MPShopProductRsp_AppStoreProduct)
	for _, category := range rsp.List {
		var product lbwxmp.MPShopProductRsp_AppStoreProduct
		find, err := OrmMpStoreProduct.NewBaseScope().Where(lbwxmp.FieldMpStoreShopId_, storeShop.Id).Where(lbwxmp.FieldCateId_, category.Id).Find(uCtx)
		if err != nil {
			log.Errorf("err:%v", err)
			continue
		}
		product.List = find
		rsp.ProductMap[category.Id] = &product
	}

	return &rsp, err
}
func (a *LbwxmpServer) GetMpShopAdsListPublic(ctx context.Context, req *lbwxmp.GetMpShopAdsListPublicReq) (*lbwxmp.GetMpShopAdsListPublicRsp, error) {
	var rsp lbwxmp.GetMpShopAdsListPublicRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	db := OrmMpShopAds.NewList(req.ListOption)
	err = gormx.ProcessDefaultOptions[*lbwxmp.ModelMpShopAds](req.ListOption, db)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	err = core.NewOptionsProcessor(req.ListOption).
		Process()

	var list []*lbwxmp.ModelMpShopAds
	_, err = db.FindPaginate(uCtx, &list)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	if len(rsp.List) == 0 {
		rsp.List = append(rsp.List, "https://oldbai.top/oss/download/BUOZ74", "https://oldbai.top/oss/download/BUOZ74")
	}
	rsp.IsActive = true

	return &rsp, err
}
