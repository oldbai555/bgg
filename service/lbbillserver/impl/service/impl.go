package service

import (
	"context"
	"github.com/oldbai555/bgg/service/lb"
	"github.com/oldbai555/bgg/service/lbbill"
	"github.com/oldbai555/bgg/service/lbbillserver/impl/mysql"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/utils"
)

var OnceSvrImpl = &LbbillServer{}

type LbbillServer struct {
	lbbill.UnimplementedLbbillServer
}

func (a *LbbillServer) AddBillSys(ctx context.Context, req *lbbill.AddBillSysReq) (*lbbill.AddBillSysRsp, error) {
	var rsp lbbill.AddBillSysRsp
	var err error

	_, err = mysql.Bill.NewScope(ctx).Create(req.Data)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.Data = req.Data

	return &rsp, err
}
func (a *LbbillServer) DelBillSysList(ctx context.Context, req *lbbill.DelBillSysListReq) (*lbbill.DelBillSysListRsp, error) {
	var rsp lbbill.DelBillSysListRsp
	var err error

	listRsp, err := a.GetBillSysList(ctx, &lbbill.GetBillSysListReq{
		ListOption: req.ListOption.
			SetSkipTotal().
			AddOpt(lb.DefaultListOption_DefaultListOptionSelect, lbbill.FieldId_),
	})
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	if len(listRsp.List) == 0 {
		log.Infof("list is empty")
		return &rsp, nil
	}

	idList := utils.PluckUint64List(listRsp.List, lbbill.FieldId)
	_, err = mysql.Bill.NewScope(ctx).In(lbbill.FieldId_, idList).Delete(&lbbill.ModelBill{})
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}
func (a *LbbillServer) UpdateBillSys(ctx context.Context, req *lbbill.UpdateBillSysReq) (*lbbill.UpdateBillSysRsp, error) {
	var rsp lbbill.UpdateBillSysRsp
	var err error

	var data lbbill.ModelBill
	err = mysql.Bill.NewScope(ctx).Where(lbbill.FieldId_, req.Data.Id).First(&data)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	_, err = mysql.Bill.NewScope(ctx).Where(lbbill.FieldId_, data.Id).Update(utils.OrmStruct2Map(req.Data))
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}
func (a *LbbillServer) GetBillSys(ctx context.Context, req *lbbill.GetBillSysReq) (*lbbill.GetBillSysRsp, error) {
	var rsp lbbill.GetBillSysRsp
	var err error

	var data lbbill.ModelBill
	err = mysql.Bill.NewScope(ctx).Where(lbbill.FieldId_, req.Id).First(&data)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.Data = &data

	return &rsp, err
}
func (a *LbbillServer) GetBillSysList(ctx context.Context, req *lbbill.GetBillSysListReq) (*lbbill.GetBillSysListRsp, error) {
	var rsp lbbill.GetBillSysListRsp
	var err error

	db := mysql.Bill.NewList(ctx, req.ListOption)
	err = lb.ProcessDefaultOptions(req.ListOption, db)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	err = lb.NewOptionsProcessor(req.ListOption).
		Process()

	rsp.Paginate, err = db.FindPaginate(&rsp.List)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}
func (a *LbbillServer) AddBillCategorySys(ctx context.Context, req *lbbill.AddBillCategorySysReq) (*lbbill.AddBillCategorySysRsp, error) {
	var rsp lbbill.AddBillCategorySysRsp
	var err error

	_, err = mysql.BillCategory.NewScope(ctx).Create(req.Data)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.Data = req.Data

	return &rsp, err
}
func (a *LbbillServer) DelBillCategorySysList(ctx context.Context, req *lbbill.DelBillCategorySysListReq) (*lbbill.DelBillCategorySysListRsp, error) {
	var rsp lbbill.DelBillCategorySysListRsp
	var err error

	listRsp, err := a.GetBillCategorySysList(ctx, &lbbill.GetBillCategorySysListReq{
		ListOption: req.ListOption.
			SetSkipTotal().
			AddOpt(lb.DefaultListOption_DefaultListOptionSelect, lbbill.FieldId_),
	})
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	if len(listRsp.List) == 0 {
		log.Infof("list is empty")
		return &rsp, nil
	}

	idList := utils.PluckUint64List(listRsp.List, lbbill.FieldId)
	_, err = mysql.BillCategory.NewScope(ctx).In(lbbill.FieldId_, idList).Delete(&lbbill.ModelBillCategory{})
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}
func (a *LbbillServer) UpdateBillCategorySys(ctx context.Context, req *lbbill.UpdateBillCategorySysReq) (*lbbill.UpdateBillCategorySysRsp, error) {
	var rsp lbbill.UpdateBillCategorySysRsp
	var err error

	var data lbbill.ModelBillCategory
	err = mysql.BillCategory.NewScope(ctx).Where(lbbill.FieldId_, req.Data.Id).First(&data)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	_, err = mysql.BillCategory.NewScope(ctx).Where(lbbill.FieldId_, data.Id).Update(utils.OrmStruct2Map(req.Data))
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}
func (a *LbbillServer) GetBillCategorySys(ctx context.Context, req *lbbill.GetBillCategorySysReq) (*lbbill.GetBillCategorySysRsp, error) {
	var rsp lbbill.GetBillCategorySysRsp
	var err error

	var data lbbill.ModelBillCategory
	err = mysql.BillCategory.NewScope(ctx).Where(lbbill.FieldId_, req.Id).First(&data)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.Data = &data

	return &rsp, err
}
func (a *LbbillServer) GetBillCategorySysList(ctx context.Context, req *lbbill.GetBillCategorySysListReq) (*lbbill.GetBillCategorySysListRsp, error) {
	var rsp lbbill.GetBillCategorySysListRsp
	var err error

	db := mysql.BillCategory.NewList(ctx, req.ListOption)
	err = lb.ProcessDefaultOptions(req.ListOption, db)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	err = lb.NewOptionsProcessor(req.ListOption).
		Process()

	rsp.Paginate, err = db.FindPaginate(&rsp.List)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}
