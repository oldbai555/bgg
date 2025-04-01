package lbsingleserver

import (
	"context"
	"github.com/oldbai555/bgg/service/lbsingle"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/lberr"
	"github.com/oldbai555/lbtool/utils"
	"github.com/oldbai555/micro/core"
	"github.com/oldbai555/micro/gormx"
	"github.com/oldbai555/micro/uctx"
)

func (a *LbsingleServer) AddDictData(ctx context.Context, req *lbsingle.AddDictDataReq) (*lbsingle.AddDictDataRsp, error) {
	var rsp lbsingle.AddDictDataRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	err = OrmDictData.NewBaseScope().Create(uCtx, req.Data)
	if err != nil {
		return nil, lberr.Wrap(err)
	}
	rsp.Data = req.Data

	return &rsp, err
}

func (a *LbsingleServer) DelDictDataList(ctx context.Context, req *lbsingle.DelDictDataListReq) (*lbsingle.DelDictDataListRsp, error) {
	var rsp lbsingle.DelDictDataListRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	listRsp, err := a.GetDictDataList(ctx, &lbsingle.GetDictDataListReq{
		ListOption: req.ListOption.
			SetSkipTotal().
			AddOpt(core.DefaultListOption_DefaultListOptionSelect, lbsingle.FieldId_),
	})
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	if len(listRsp.List) == 0 {
		log.Infof("list is empty")
		return &rsp, nil
	}

	idList := utils.PluckUint64List(listRsp.List, lbsingle.FieldId)
	_, err = OrmDictData.NewBaseScope().WhereIn(lbsingle.FieldId_, idList).Delete(uCtx)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	return &rsp, err
}

func (a *LbsingleServer) UpdateDictData(ctx context.Context, req *lbsingle.UpdateDictDataReq) (*lbsingle.UpdateDictDataRsp, error) {
	var rsp lbsingle.UpdateDictDataRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	data, err := OrmDictData.NewBaseScope().Where(lbsingle.FieldId_, req.Data.Id).First(uCtx)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	_, err = OrmDictData.NewBaseScope().Where(lbsingle.FieldId_, data.Id).Update(uCtx, utils.OrmStruct2Map(req.Data))
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	return &rsp, err
}

func (a *LbsingleServer) GetDictData(ctx context.Context, req *lbsingle.GetDictDataReq) (*lbsingle.GetDictDataRsp, error) {
	var rsp lbsingle.GetDictDataRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	data, err := OrmDictData.NewBaseScope().Where(lbsingle.FieldId_, req.Id).First(uCtx)
	if err != nil {
		return nil, lberr.Wrap(err)
	}
	rsp.Data = data

	return &rsp, err
}

func (a *LbsingleServer) GetDictDataList(ctx context.Context, req *lbsingle.GetDictDataListReq) (*lbsingle.GetDictDataListRsp, error) {
	var rsp lbsingle.GetDictDataListRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	db := OrmDictData.NewList(req.ListOption)
	err = gormx.ProcessDefaultOptions[*lbsingle.ModelDictData](req.ListOption, db)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	err = core.NewOptionsProcessor(req.ListOption).
		AddUint64(lbsingle.GetDictDataListReq_ListOptionDictTypeId, func(val uint64) error {
			db.Where(lbsingle.FieldDictTypeId_, val)
			return nil
		}).
		Process()

	rsp.Paginate, err = db.FindPaginate(uCtx, &rsp.List)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	return &rsp, err
}
