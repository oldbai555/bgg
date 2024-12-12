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

func (a *LbwxmpServer) AddMpStoreOrderCartInfo(ctx context.Context, req *lbwxmp.AddMpStoreOrderCartInfoReq) (*lbwxmp.AddMpStoreOrderCartInfoRsp, error) {
	var rsp lbwxmp.AddMpStoreOrderCartInfoRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	err = OrmMpStoreOrderCartInfo.NewBaseScope().Create(uCtx, req.Data)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.Data = req.Data

	return &rsp, err
}
func (a *LbwxmpServer) DelMpStoreOrderCartInfoList(ctx context.Context, req *lbwxmp.DelMpStoreOrderCartInfoListReq) (*lbwxmp.DelMpStoreOrderCartInfoListRsp, error) {
	var rsp lbwxmp.DelMpStoreOrderCartInfoListRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	listRsp, err := a.GetMpStoreOrderCartInfoList(ctx, &lbwxmp.GetMpStoreOrderCartInfoListReq{
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
	_, err = OrmMpStoreOrderCartInfo.NewBaseScope().WhereIn(lbwxmp.FieldId_, idList).Delete(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}
func (a *LbwxmpServer) UpdateMpStoreOrderCartInfo(ctx context.Context, req *lbwxmp.UpdateMpStoreOrderCartInfoReq) (*lbwxmp.UpdateMpStoreOrderCartInfoRsp, error) {
	var rsp lbwxmp.UpdateMpStoreOrderCartInfoRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	data, err := OrmMpStoreOrderCartInfo.NewBaseScope().Where(lbwxmp.FieldId_, req.Data.Id).First(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	_, err = OrmMpStoreOrderCartInfo.NewBaseScope().Where(lbwxmp.FieldId_, data.Id).Update(uCtx, utils.OrmStruct2Map(req.Data))
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}
func (a *LbwxmpServer) GetMpStoreOrderCartInfo(ctx context.Context, req *lbwxmp.GetMpStoreOrderCartInfoReq) (*lbwxmp.GetMpStoreOrderCartInfoRsp, error) {
	var rsp lbwxmp.GetMpStoreOrderCartInfoRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	data, err := OrmMpStoreOrderCartInfo.NewBaseScope().Where(lbwxmp.FieldId_, req.Id).First(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.Data = data

	return &rsp, err
}
func (a *LbwxmpServer) GetMpStoreOrderCartInfoList(ctx context.Context, req *lbwxmp.GetMpStoreOrderCartInfoListReq) (*lbwxmp.GetMpStoreOrderCartInfoListRsp, error) {
	var rsp lbwxmp.GetMpStoreOrderCartInfoListRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	db := OrmMpStoreOrderCartInfo.NewList(req.ListOption)
	err = gormx.ProcessDefaultOptions[*lbwxmp.ModelMpStoreOrderCartInfo](req.ListOption, db)
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
