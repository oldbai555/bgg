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

func (a *LbwxmpServer) AddMpProductCategory(ctx context.Context, req *lbwxmp.AddMpProductCategoryReq) (*lbwxmp.AddMpProductCategoryRsp, error) {
	var rsp lbwxmp.AddMpProductCategoryRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	err = OrmMpProductCategory.NewBaseScope().Create(uCtx, req.Data)
	if err != nil {
		return nil, lberr.Wrap(err)
	}
	rsp.Data = req.Data

	return &rsp, err
}
func (a *LbwxmpServer) DelMpProductCategoryList(ctx context.Context, req *lbwxmp.DelMpProductCategoryListReq) (*lbwxmp.DelMpProductCategoryListRsp, error) {
	var rsp lbwxmp.DelMpProductCategoryListRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	listRsp, err := a.GetMpProductCategoryList(ctx, &lbwxmp.GetMpProductCategoryListReq{
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
	_, err = OrmMpProductCategory.NewBaseScope().WhereIn(lbwxmp.FieldId_, idList).Delete(uCtx)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	return &rsp, err
}
func (a *LbwxmpServer) UpdateMpProductCategory(ctx context.Context, req *lbwxmp.UpdateMpProductCategoryReq) (*lbwxmp.UpdateMpProductCategoryRsp, error) {
	var rsp lbwxmp.UpdateMpProductCategoryRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	data, err := OrmMpProductCategory.NewBaseScope().Where(lbwxmp.FieldId_, req.Data.Id).First(uCtx)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	_, err = OrmMpProductCategory.NewBaseScope().Where(lbwxmp.FieldId_, data.Id).Update(uCtx, utils.OrmStruct2Map(req.Data))
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	return &rsp, err
}
func (a *LbwxmpServer) GetMpProductCategory(ctx context.Context, req *lbwxmp.GetMpProductCategoryReq) (*lbwxmp.GetMpProductCategoryRsp, error) {
	var rsp lbwxmp.GetMpProductCategoryRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	data, err := OrmMpProductCategory.NewBaseScope().Where(lbwxmp.FieldId_, req.Id).First(uCtx)
	if err != nil {
		return nil, lberr.Wrap(err)
	}
	rsp.Data = data

	return &rsp, err
}
func (a *LbwxmpServer) GetMpProductCategoryList(ctx context.Context, req *lbwxmp.GetMpProductCategoryListReq) (*lbwxmp.GetMpProductCategoryListRsp, error) {
	var rsp lbwxmp.GetMpProductCategoryListRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	db := OrmMpProductCategory.NewList(req.ListOption)
	err = gormx.ProcessDefaultOptions[*lbwxmp.ModelMpProductCategory](req.ListOption, db)
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
