package service

import (
	"context"
	"github.com/oldbai555/bgg/service/lbbill"
	"github.com/oldbai555/bgg/service/lbbillserver/impl/dao/impl/mysql"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/utils"
)

var ServerImpl LbbillServer

type LbbillServer struct {
	*lbbill.UnimplementedLbbillServer
}

func (a *LbbillServer) AddBill(ctx context.Context, req *lbbill.AddBillReq) (*lbbill.AddBillRsp, error) {
	var rsp lbbill.AddBillRsp
	var err error

	req.Bill.DateUnix = utils.TimeNow()
	_, err = mysql.BillOrm.Create(ctx, req.Bill)
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}
	rsp.Bill = req.Bill

	return &rsp, err
}
func (a *LbbillServer) DelBill(ctx context.Context, req *lbbill.DelBillReq) (*lbbill.DelBillRsp, error) {
	var rsp lbbill.DelBillRsp
	var err error

	_, err = mysql.BillOrm.DeleteById(ctx, req.Id)
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}

	return &rsp, err
}
func (a *LbbillServer) UpdateBill(ctx context.Context, req *lbbill.UpdateBillReq) (*lbbill.UpdateBillRsp, error) {
	var rsp lbbill.UpdateBillRsp
	var err error

	_, err = mysql.BillOrm.UpdateById(ctx, req.Bill.Id, utils.OrmStruct2Map(req.Bill))
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}

	return &rsp, err
}
func (a *LbbillServer) GetBill(ctx context.Context, req *lbbill.GetBillReq) (*lbbill.GetBillRsp, error) {
	var rsp lbbill.GetBillRsp
	var err error

	rsp.Bill, err = mysql.BillOrm.GetById(ctx, req.Id)
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}

	return &rsp, err
}
func (a *LbbillServer) GetBillList(ctx context.Context, req *lbbill.GetBillListReq) (*lbbill.GetBillListRsp, error) {
	var rsp lbbill.GetBillListRsp
	var err error

	rsp.List, rsp.Paginate, err = mysql.BillOrm.FindPaginate(ctx, req.Options)
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}

	categoryIdList := utils.PluckUint64List(rsp.List, lbbill.FieldCategoryId)
	billCategories, err := mysql.BillCategoryOrm.GetByIdList(ctx, categoryIdList)
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}
	rsp.CategoryMap = utils.Slice2MapKeyByStructField(billCategories, lbbill.FieldId).(map[uint64]*lbbill.ModelBillCategory)

	return &rsp, err
}
func (a *LbbillServer) AddBillCategory(ctx context.Context, req *lbbill.AddBillCategoryReq) (*lbbill.AddBillCategoryRsp, error) {
	var rsp lbbill.AddBillCategoryRsp
	var err error

	_, err = mysql.BillCategoryOrm.Create(ctx, req.Category)
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}

	rsp.Category = req.Category

	return &rsp, err
}
func (a *LbbillServer) DelBillCategory(ctx context.Context, req *lbbill.DelBillCategoryReq) (*lbbill.DelBillCategoryRsp, error) {
	var rsp lbbill.DelBillCategoryRsp
	var err error

	_, err = mysql.BillCategoryOrm.DeleteById(ctx, req.Id)
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}

	return &rsp, err
}
func (a *LbbillServer) UpdateBillCategory(ctx context.Context, req *lbbill.UpdateBillCategoryReq) (*lbbill.UpdateBillCategoryRsp, error) {
	var rsp lbbill.UpdateBillCategoryRsp
	var err error

	_, err = mysql.BillCategoryOrm.UpdateById(ctx, req.Category.Id, utils.OrmStruct2Map(req.Category))
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}

	return &rsp, err
}
func (a *LbbillServer) GetBillCategory(ctx context.Context, req *lbbill.GetBillCategoryReq) (*lbbill.GetBillCategoryRsp, error) {
	var rsp lbbill.GetBillCategoryRsp
	var err error

	rsp.Category, err = mysql.BillCategoryOrm.GetById(ctx, req.Id)
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}

	return &rsp, err
}
func (a *LbbillServer) GetBillCategoryList(ctx context.Context, req *lbbill.GetBillCategoryListReq) (*lbbill.GetBillCategoryListRsp, error) {
	var rsp lbbill.GetBillCategoryListRsp
	var err error

	rsp.List, rsp.Paginate, err = mysql.BillCategoryOrm.FindPaginate(ctx, req.Options)
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}

	return &rsp, err
}
