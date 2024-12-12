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

func (a *LbwxmpServer) AddMpService(ctx context.Context, req *lbwxmp.AddMpServiceReq) (*lbwxmp.AddMpServiceRsp, error) {
	var rsp lbwxmp.AddMpServiceRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	err = OrmMpService.NewBaseScope().Create(uCtx, req.Data)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.Data = req.Data

	return &rsp, err
}
func (a *LbwxmpServer) DelMpServiceList(ctx context.Context, req *lbwxmp.DelMpServiceListReq) (*lbwxmp.DelMpServiceListRsp, error) {
	var rsp lbwxmp.DelMpServiceListRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	listRsp, err := a.GetMpServiceList(ctx, &lbwxmp.GetMpServiceListReq{
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
	_, err = OrmMpService.NewBaseScope().WhereIn(lbwxmp.FieldId_, idList).Delete(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}
func (a *LbwxmpServer) UpdateMpService(ctx context.Context, req *lbwxmp.UpdateMpServiceReq) (*lbwxmp.UpdateMpServiceRsp, error) {
	var rsp lbwxmp.UpdateMpServiceRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	data, err := OrmMpService.NewBaseScope().Where(lbwxmp.FieldId_, req.Data.Id).First(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	_, err = OrmMpService.NewBaseScope().Where(lbwxmp.FieldId_, data.Id).Update(uCtx, utils.OrmStruct2Map(req.Data))
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}
func (a *LbwxmpServer) GetMpService(ctx context.Context, req *lbwxmp.GetMpServiceReq) (*lbwxmp.GetMpServiceRsp, error) {
	var rsp lbwxmp.GetMpServiceRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	data, err := OrmMpService.NewBaseScope().Where(lbwxmp.FieldId_, req.Id).First(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.Data = data

	return &rsp, err
}
func (a *LbwxmpServer) GetMpServiceList(ctx context.Context, req *lbwxmp.GetMpServiceListReq) (*lbwxmp.GetMpServiceListRsp, error) {
	var rsp lbwxmp.GetMpServiceListRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	db := OrmMpService.NewList(req.ListOption)
	err = gormx.ProcessDefaultOptions[*lbwxmp.ModelMpService](req.ListOption, db)
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
