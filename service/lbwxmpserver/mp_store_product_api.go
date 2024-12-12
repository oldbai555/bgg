package lbwxmpserver

import (
	"context"
	"github.com/oldbai555/bgg/service/lbwxmp"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/utils"
	"github.com/oldbai555/micro/core"
	"github.com/oldbai555/micro/gormx"
	"github.com/oldbai555/micro/uctx"
)

func (a *LbwxmpServer) AddMpStoreProduct(ctx context.Context, req *lbwxmp.AddMpStoreProductReq) (*lbwxmp.AddMpStoreProductRsp, error) {
	var rsp lbwxmp.AddMpStoreProductRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	err = OrmMpStoreProduct.NewBaseScope().Create(uCtx, req.Data)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.Data = req.Data

	return &rsp, err
}
func (a *LbwxmpServer) DelMpStoreProductList(ctx context.Context, req *lbwxmp.DelMpStoreProductListReq) (*lbwxmp.DelMpStoreProductListRsp, error) {
	var rsp lbwxmp.DelMpStoreProductListRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	listRsp, err := a.GetMpStoreProductList(ctx, &lbwxmp.GetMpStoreProductListReq{
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
	_, err = OrmMpStoreProduct.NewBaseScope().WhereIn(lbwxmp.FieldId_, idList).Delete(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}
func (a *LbwxmpServer) UpdateMpStoreProduct(ctx context.Context, req *lbwxmp.UpdateMpStoreProductReq) (*lbwxmp.UpdateMpStoreProductRsp, error) {
	var rsp lbwxmp.UpdateMpStoreProductRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	data, err := OrmMpStoreProduct.NewBaseScope().Where(lbwxmp.FieldId_, req.Data.Id).First(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	_, err = OrmMpStoreProduct.NewBaseScope().Where(lbwxmp.FieldId_, data.Id).Update(uCtx, utils.OrmStruct2Map(req.Data))
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}
func (a *LbwxmpServer) GetMpStoreProduct(ctx context.Context, req *lbwxmp.GetMpStoreProductReq) (*lbwxmp.GetMpStoreProductRsp, error) {
	var rsp lbwxmp.GetMpStoreProductRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	data, err := OrmMpStoreProduct.NewBaseScope().Where(lbwxmp.FieldId_, req.Id).First(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.Data = data

	return &rsp, err
}
func (a *LbwxmpServer) GetMpStoreProductList(ctx context.Context, req *lbwxmp.GetMpStoreProductListReq) (*lbwxmp.GetMpStoreProductListRsp, error) {
	var rsp lbwxmp.GetMpStoreProductListRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	db := OrmMpStoreProduct.NewList(req.ListOption)
	err = gormx.ProcessDefaultOptions[*lbwxmp.ModelMpStoreProduct](req.ListOption, db)
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
