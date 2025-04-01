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

func (a *LbsingleServer) AddArticleCategory(ctx context.Context, req *lbsingle.AddArticleCategoryReq) (*lbsingle.AddArticleCategoryRsp, error) {
	var rsp lbsingle.AddArticleCategoryRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	err = OrmArticleCategory.NewBaseScope().Create(uCtx, req.Data)
	if err != nil {
		return nil, lberr.Wrap(err)
	}
	rsp.Data = req.Data

	return &rsp, err
}
func (a *LbsingleServer) DelArticleCategoryList(ctx context.Context, req *lbsingle.DelArticleCategoryListReq) (*lbsingle.DelArticleCategoryListRsp, error) {
	var rsp lbsingle.DelArticleCategoryListRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	listRsp, err := a.GetArticleCategoryList(ctx, &lbsingle.GetArticleCategoryListReq{
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
	_, err = OrmArticleCategory.NewBaseScope().WhereIn(lbsingle.FieldId_, idList).Delete(uCtx)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	return &rsp, err
}
func (a *LbsingleServer) UpdateArticleCategory(ctx context.Context, req *lbsingle.UpdateArticleCategoryReq) (*lbsingle.UpdateArticleCategoryRsp, error) {
	var rsp lbsingle.UpdateArticleCategoryRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	data, err := OrmArticleCategory.NewBaseScope().Where(lbsingle.FieldId_, req.Data.Id).First(uCtx)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	_, err = OrmArticleCategory.NewBaseScope().Where(lbsingle.FieldId_, data.Id).Update(uCtx, utils.OrmStruct2Map(req.Data))
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	return &rsp, err
}
func (a *LbsingleServer) GetArticleCategory(ctx context.Context, req *lbsingle.GetArticleCategoryReq) (*lbsingle.GetArticleCategoryRsp, error) {
	var rsp lbsingle.GetArticleCategoryRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	data, err := OrmArticleCategory.NewBaseScope().Where(lbsingle.FieldId_, req.Id).First(uCtx)
	if err != nil {
		return nil, lberr.Wrap(err)
	}
	rsp.Data = data

	return &rsp, err
}
func (a *LbsingleServer) GetArticleCategoryList(ctx context.Context, req *lbsingle.GetArticleCategoryListReq) (*lbsingle.GetArticleCategoryListRsp, error) {
	var rsp lbsingle.GetArticleCategoryListRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	db := OrmArticleCategory.NewList(req.ListOption)
	err = gormx.ProcessDefaultOptions[*lbsingle.ModelArticleCategory](req.ListOption, db)
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
