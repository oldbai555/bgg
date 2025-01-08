package lbwxmpserver

import (
	"context"
	"github.com/oldbai555/bgg/service/lbwxmp"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/lberr"
	"github.com/oldbai555/lbtool/utils"
	"github.com/oldbai555/micro/core"
	"github.com/oldbai555/micro/gormx"
	"github.com/oldbai555/micro/uctx"
)

func (a *LbwxmpServer) AddMpStoreShop(ctx context.Context, req *lbwxmp.AddMpStoreShopReq) (*lbwxmp.AddMpStoreShopRsp, error) {
	var rsp lbwxmp.AddMpStoreShopRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	err = OrmMpStoreShop.NewBaseScope().Create(uCtx, req.Data)
	if err != nil {
		return nil, lberr.Wrap(err)
	}
	rsp.Data = req.Data

	return &rsp, err
}
func (a *LbwxmpServer) DelMpStoreShopList(ctx context.Context, req *lbwxmp.DelMpStoreShopListReq) (*lbwxmp.DelMpStoreShopListRsp, error) {
	var rsp lbwxmp.DelMpStoreShopListRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	listRsp, err := a.GetMpStoreShopList(ctx, &lbwxmp.GetMpStoreShopListReq{
		ListOption: req.ListOption.
			SetSkipTotal().
			AddOpt(core.DefaultListOption_DefaultListOptionSelect, lbwxmp.FieldId_),
	})
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	if len(listRsp.List) == 0 {
		log.Infof("list is empty")
		return &rsp, nil
	}

	idList := utils.PluckUint64List(listRsp.List, lbwxmp.FieldId)
	_, err = OrmMpStoreShop.NewBaseScope().WhereIn(lbwxmp.FieldId_, idList).Delete(uCtx)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	return &rsp, err
}
func (a *LbwxmpServer) UpdateMpStoreShop(ctx context.Context, req *lbwxmp.UpdateMpStoreShopReq) (*lbwxmp.UpdateMpStoreShopRsp, error) {
	var rsp lbwxmp.UpdateMpStoreShopRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	data, err := OrmMpStoreShop.NewBaseScope().Where(lbwxmp.FieldId_, req.Data.Id).First(uCtx)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	_, err = OrmMpStoreShop.NewBaseScope().Where(lbwxmp.FieldId_, data.Id).Update(uCtx, utils.OrmStruct2Map(req.Data))
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	return &rsp, err
}
func (a *LbwxmpServer) GetMpStoreShop(ctx context.Context, req *lbwxmp.GetMpStoreShopReq) (*lbwxmp.GetMpStoreShopRsp, error) {
	var rsp lbwxmp.GetMpStoreShopRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	data, err := OrmMpStoreShop.NewBaseScope().Where(lbwxmp.FieldId_, req.Id).First(uCtx)
	if err != nil {
		return nil, lberr.Wrap(err)
	}
	rsp.Data = data

	return &rsp, err
}
func (a *LbwxmpServer) GetMpStoreShopList(ctx context.Context, req *lbwxmp.GetMpStoreShopListReq) (*lbwxmp.GetMpStoreShopListRsp, error) {
	var rsp lbwxmp.GetMpStoreShopListRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	db := OrmMpStoreShop.NewList(req.ListOption)
	err = gormx.ProcessDefaultOptions[*lbwxmp.ModelMpStoreShop](req.ListOption, db)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	err = core.NewOptionsProcessor(req.ListOption).
		Process()

	rsp.Paginate, err = db.FindPaginate(uCtx, &rsp.List)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	return &rsp, err
}
