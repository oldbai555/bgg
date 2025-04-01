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

func (a *LbsingleServer) AddArticle(ctx context.Context, req *lbsingle.AddArticleReq) (*lbsingle.AddArticleRsp, error) {
	var rsp lbsingle.AddArticleRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	err = OrmArticle.NewBaseScope().Create(uCtx, req.Data)
	if err != nil {
		return nil, lberr.Wrap(err)
	}
	rsp.Data = req.Data

	return &rsp, err
}
func (a *LbsingleServer) DelArticleList(ctx context.Context, req *lbsingle.DelArticleListReq) (*lbsingle.DelArticleListRsp, error) {
	var rsp lbsingle.DelArticleListRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	listRsp, err := a.GetArticleList(ctx, &lbsingle.GetArticleListReq{
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
	_, err = OrmArticle.NewBaseScope().WhereIn(lbsingle.FieldId_, idList).Delete(uCtx)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	return &rsp, err
}
func (a *LbsingleServer) UpdateArticle(ctx context.Context, req *lbsingle.UpdateArticleReq) (*lbsingle.UpdateArticleRsp, error) {
	var rsp lbsingle.UpdateArticleRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	data, err := OrmArticle.NewBaseScope().Where(lbsingle.FieldId_, req.Data.Id).First(uCtx)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	_, err = OrmArticle.NewBaseScope().Where(lbsingle.FieldId_, data.Id).Update(uCtx, utils.OrmStruct2Map(req.Data))
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	return &rsp, err
}
func (a *LbsingleServer) GetArticle(ctx context.Context, req *lbsingle.GetArticleReq) (*lbsingle.GetArticleRsp, error) {
	var rsp lbsingle.GetArticleRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	data, err := OrmArticle.NewBaseScope().Where(lbsingle.FieldId_, req.Id).First(uCtx)
	if err != nil {
		return nil, lberr.Wrap(err)
	}
	rsp.Data = data

	return &rsp, err
}
func (a *LbsingleServer) GetArticleList(ctx context.Context, req *lbsingle.GetArticleListReq) (*lbsingle.GetArticleListRsp, error) {
	var rsp lbsingle.GetArticleListRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	db := OrmArticle.NewList(req.ListOption)
	err = gormx.ProcessDefaultOptions[*lbsingle.ModelArticle](req.ListOption, db)
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
