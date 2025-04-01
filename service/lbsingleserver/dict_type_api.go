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

func (a *LbsingleServer) AddDictType(ctx context.Context, req *lbsingle.AddDictTypeReq) (*lbsingle.AddDictTypeRsp, error) {
	var rsp lbsingle.AddDictTypeRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	err = OrmDictType.NewBaseScope().Create(uCtx, req.Data)
	if err != nil {
		return nil, lberr.Wrap(err)
	}
	rsp.Data = req.Data

	return &rsp, err
}
func (a *LbsingleServer) DelDictTypeList(ctx context.Context, req *lbsingle.DelDictTypeListReq) (*lbsingle.DelDictTypeListRsp, error) {
	var rsp lbsingle.DelDictTypeListRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	listRsp, err := a.GetDictTypeList(ctx, &lbsingle.GetDictTypeListReq{
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
	_, err = OrmDictType.NewBaseScope().WhereIn(lbsingle.FieldId_, idList).Delete(uCtx)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	return &rsp, err
}
func (a *LbsingleServer) UpdateDictType(ctx context.Context, req *lbsingle.UpdateDictTypeReq) (*lbsingle.UpdateDictTypeRsp, error) {
	var rsp lbsingle.UpdateDictTypeRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	data, err := OrmDictType.NewBaseScope().Where(lbsingle.FieldId_, req.Data.Id).First(uCtx)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	_, err = OrmDictType.NewBaseScope().Where(lbsingle.FieldId_, data.Id).Update(uCtx, utils.OrmStruct2Map(req.Data))
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	return &rsp, err
}
func (a *LbsingleServer) GetDictType(ctx context.Context, req *lbsingle.GetDictTypeReq) (*lbsingle.GetDictTypeRsp, error) {
	var rsp lbsingle.GetDictTypeRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	data, err := OrmDictType.NewBaseScope().Where(lbsingle.FieldId_, req.Id).First(uCtx)
	if err != nil {
		return nil, lberr.Wrap(err)
	}
	rsp.Data = data

	return &rsp, err
}
func (a *LbsingleServer) GetDictTypeList(ctx context.Context, req *lbsingle.GetDictTypeListReq) (*lbsingle.GetDictTypeListRsp, error) {
	var rsp lbsingle.GetDictTypeListRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	db := OrmDictType.NewList(req.ListOption)
	err = gormx.ProcessDefaultOptions[*lbsingle.ModelDictType](req.ListOption, db)
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
