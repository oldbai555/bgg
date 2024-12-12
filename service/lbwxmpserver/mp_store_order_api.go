package lbwxmpserver

import (
	"context"
	"github.com/oldbai555/bgg/service/lbbase"
	"github.com/oldbai555/bgg/service/lbsingle"
	"github.com/oldbai555/bgg/service/lbwxmp"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/utils"
	"github.com/oldbai555/micro/core"
	"github.com/oldbai555/micro/gormx"
	"github.com/oldbai555/micro/uctx"
	"time"
)

func (a *LbwxmpServer) AddMpStoreOrder(ctx context.Context, req *lbwxmp.AddMpStoreOrderReq) (*lbwxmp.AddMpStoreOrderRsp, error) {
	var rsp lbwxmp.AddMpStoreOrderRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	err = OrmMpStoreOrder.NewBaseScope().Create(uCtx, req.Data)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.Data = req.Data

	return &rsp, err
}
func (a *LbwxmpServer) DelMpStoreOrderList(ctx context.Context, req *lbwxmp.DelMpStoreOrderListReq) (*lbwxmp.DelMpStoreOrderListRsp, error) {
	var rsp lbwxmp.DelMpStoreOrderListRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	listRsp, err := a.GetMpStoreOrderList(ctx, &lbwxmp.GetMpStoreOrderListReq{
		ListOption: req.ListOption.
			SetSkipTotal().
			AddOpt(core.DefaultListOption_DefaultListOptionSelect, lbwxmp.FieldId_),
	})
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	if len(listRsp.List) == 0 {
		log.Infof("list is empty")
		return &rsp, nil
	}

	idList := utils.PluckUint64List(listRsp.List, lbwxmp.FieldId)
	_, err = OrmMpStoreOrder.NewBaseScope().WhereIn(lbwxmp.FieldId_, idList).Delete(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}
func (a *LbwxmpServer) UpdateMpStoreOrder(ctx context.Context, req *lbwxmp.UpdateMpStoreOrderReq) (*lbwxmp.UpdateMpStoreOrderRsp, error) {
	var rsp lbwxmp.UpdateMpStoreOrderRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	data, err := OrmMpStoreOrder.NewBaseScope().Where(lbwxmp.FieldId_, req.Data.Id).First(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	_, err = OrmMpStoreOrder.NewBaseScope().Where(lbwxmp.FieldId_, data.Id).Update(uCtx, utils.OrmStruct2Map(req.Data))
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}
func (a *LbwxmpServer) GetMpStoreOrder(ctx context.Context, req *lbwxmp.GetMpStoreOrderReq) (*lbwxmp.GetMpStoreOrderRsp, error) {
	var rsp lbwxmp.GetMpStoreOrderRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	data, err := OrmMpStoreOrder.NewBaseScope().Where(lbwxmp.FieldId_, req.Id).First(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.Data = data

	return &rsp, err
}
func (a *LbwxmpServer) GetMpStoreOrderList(ctx context.Context, req *lbwxmp.GetMpStoreOrderListReq) (*lbwxmp.GetMpStoreOrderListRsp, error) {
	var rsp lbwxmp.GetMpStoreOrderListRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	db := OrmMpStoreOrder.NewList(req.ListOption)
	err = gormx.ProcessDefaultOptions[*lbwxmp.ModelMpStoreOrder](req.ListOption, db)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	err = core.NewOptionsProcessor(req.ListOption).
		Process()

	rsp.Paginate, err = db.FindPaginate(uCtx, &rsp.List)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}
func (a *LbwxmpServer) MPShopOrderCreate(ctx context.Context, req *lbwxmp.MPShopOrderCreateReq) (*lbwxmp.MPShopOrderCreateRsp, error) {
	var rsp lbwxmp.MPShopOrderCreateRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	user, ok := uCtx.ExtInfo().(*lbbase.BaseUser)
	if !ok {
		return nil, lbsingle.ErrUserNotFound
	}

	storeShop, err := OrmMpStoreShop.NewBaseScope().Where(lbwxmp.FieldId_, req.ShopId).First(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	orderSn := GenerateOrderID()
	orderNumber := &lbwxmp.ModelMpOrderNumber{
		OrderSn: "orderSn",
	}
	err = OrmMpOrderNumber.Create(uCtx, orderNumber)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	// 生成一个订单号
	order := &lbwxmp.ModelMpStoreOrder{
		OrderSn:       orderSn,
		MpUid:         user.Id,
		NumberId:      int64(orderNumber.GetId()), // 取餐号
		MpStoreShopId: storeShop.Id,
		GetAt:         uint32(time.Now().Add(time.Duration(req.GetGetTime()) * time.Minute).Unix()),
		TotalNum:      1,
		PayType:       req.PayType,
		Mark:          req.Remark,
		ShippingType:  2, // 默认门店自提
		OrderType:     req.OrderType,
	}

	err = OrmMpStoreOrder.NewBaseScope().Create(uCtx, order)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	// 添加进入购物车
	for i, productId := range req.ProductId {
		product, err := OrmMpStoreProduct.NewBaseScope().Where(lbwxmp.FieldId_, productId).First(uCtx)
		if err != nil {
			log.Errorf("err:%v", err)
			return nil, err
		}
		var shopCar = &lbwxmp.ModelMpStoreOrderCartInfo{
			MpOrderId:    order.Id,
			OrderSn:      orderSn,
			ProductId:    productId,
			Unique:       utils.GenUUID(),
			IsAfterSales: 0,
			Title:        product.Name,
			Image:        product.Image,
			Number:       1,
			Spec:         req.Spec[i],
		}
		err = OrmMpStoreOrderCartInfo.NewBaseScope().Create(uCtx, shopCar)
		if err != nil {
			log.Errorf("err:%v", err)
			return nil, err
		}
	}

	// 订单状态
	var orderStatus = &lbwxmp.ModelMpStoreOrderStatus{
		Oid:           order.Id,
		ChangeType:    "yshop_create_order",
		ChangeMessage: "订单生成",
	}
	err = OrmMpStoreOrderStatus.NewBaseScope().Create(uCtx, orderStatus)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	rsp.OrderSn = orderSn

	return &rsp, err
}
func (a *LbwxmpServer) MPShopOrderList(ctx context.Context, req *lbwxmp.MPShopOrderListReq) (*lbwxmp.MPShopOrderListRsp, error) {
	var rsp lbwxmp.MPShopOrderListRsp
	var err error
	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	db := OrmMpStoreOrder.NewList(&core.ListOption{
		Offset:    uint32((req.Page - 1) * req.Limit),
		Limit:     uint32(req.Limit),
		SkipTotal: true,
	})
	if req.Type >= 0 {
		db.Where(lbwxmp.FieldStatus_, req.Type)
	}
	list, err := db.Find(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.List = list

	mpStoreShopIds := utils.PluckUint64List(rsp.List, lbwxmp.FieldMpStoreShopId)
	mpStoreShopList, err := OrmMpStoreShop.WhereIn(lbwxmp.FieldId_, mpStoreShopIds).Find(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.ShopMap = utils.Slice2MapKeyByStructField(mpStoreShopList, lbwxmp.FieldId).(map[uint64]*lbwxmp.ModelMpStoreShop)

	orderSnList := utils.PluckStringList(rsp.List, lbwxmp.FieldOrderSn)
	cartInfoList, err := OrmMpStoreOrderCartInfo.WhereIn(lbwxmp.FieldOrderSn_, orderSnList).Find(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	rsp.CartMap = make(map[string]*lbwxmp.MPShopOrderListRsp_CartInfo)
	for _, info := range cartInfoList {
		cartInfo, ok := rsp.CartMap[info.OrderSn]
		if !ok {
			rsp.CartMap[info.OrderSn] = &lbwxmp.MPShopOrderListRsp_CartInfo{}
			cartInfo = rsp.CartMap[info.OrderSn]
		}
		cartInfo.List = append(cartInfo.List, info)
	}
	return &rsp, err
}
func (a *LbwxmpServer) MPShopOrderDetail(ctx context.Context, req *lbwxmp.MPShopOrderDetailReq) (*lbwxmp.MPShopOrderDetailRsp, error) {
	var rsp lbwxmp.MPShopOrderDetailRsp
	var err error
	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	order, err := OrmMpStoreOrder.Where(lbwxmp.FieldOrderSn_, req.OrderSn).First(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	mpStoreShop, err := OrmMpStoreShop.Where(lbwxmp.FieldId_, order.MpStoreShopId).First(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	cartInfoList, err := OrmMpStoreOrderCartInfo.Where(lbwxmp.FieldOrderSn_, order.OrderSn).Find(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	rsp.Order = order
	rsp.ShopInfo = mpStoreShop
	rsp.CartList = cartInfoList
	return &rsp, err
}
func (a *LbwxmpServer) MPShopOrderTake(ctx context.Context, req *lbwxmp.MPShopOrderTakeReq) (*lbwxmp.MPShopOrderTakeRsp, error) {
	var rsp lbwxmp.MPShopOrderTakeRsp
	var err error

	return &rsp, err
}
func (a *LbwxmpServer) MPShopOrderRefund(ctx context.Context, req *lbwxmp.MPShopOrderRefundReq) (*lbwxmp.MPShopOrderRefundRsp, error) {
	var rsp lbwxmp.MPShopOrderRefundRsp
	var err error

	return &rsp, err
}
func (a *LbwxmpServer) MPShopOrderPay(ctx context.Context, req *lbwxmp.MPShopOrderPayReq) (*lbwxmp.MPShopOrderPayRsp, error) {
	var rsp lbwxmp.MPShopOrderPayRsp
	var err error
	rsp.Status = "ok"
	return &rsp, err
}
