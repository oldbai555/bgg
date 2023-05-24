package service

import (
	"context"
	"github.com/oldbai555/bgg/client/lbbill"
	"github.com/oldbai555/bgg/pkg/webtool"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/utils"
)

var BillServer LbbillServer

type LbbillServer struct {
	*lbbill.UnimplementedLbbillServer
}

func (a *LbbillServer) AddBill(ctx context.Context, req *lbbill.AddBillReq) (*lbbill.AddBillRsp, error) {
	var rsp lbbill.AddBillRsp
	var err error

	claims, err := webtool.GetClaimsWithCtx(ctx)
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}
	req.Bill.CreatorUid = claims.GetUserId()
	req.Bill.DateUnix = utils.TimeNow()

	err = Bill.Create(ctx, req.Bill)
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}
	rsp.Bill = req.Bill

	return &rsp, nil
}
func (a *LbbillServer) AddBillCategory(ctx context.Context, req *lbbill.AddBillCategoryReq) (*lbbill.AddBillCategoryRsp, error) {
	var rsp lbbill.AddBillCategoryRsp
	var err error

	err = BillCate.Create(ctx, req.Category)
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}

	rsp.Category = req.Category
	return &rsp, nil
}
func (a *LbbillServer) DelBillCategory(ctx context.Context, req *lbbill.DelBillCategoryReq) (*lbbill.DelBillCategoryRsp, error) {
	var rsp lbbill.DelBillCategoryRsp
	var err error

	err = BillCate.DeleteById(ctx, req.Id)
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}

	return &rsp, nil
}
func (a *LbbillServer) UpdateBillCategory(ctx context.Context, req *lbbill.UpdateBillCategoryReq) (*lbbill.UpdateBillCategoryRsp, error) {
	var rsp lbbill.UpdateBillCategoryRsp
	var err error

	err = BillCate.UpdateById(ctx, req.Category.Id, utils.OrmStruct2Map(req.Category))
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}

	return &rsp, nil
}
func (a *LbbillServer) GetBillCategoryList(ctx context.Context, req *lbbill.GetBillCategoryListReq) (*lbbill.GetBillCategoryListRsp, error) {
	var rsp lbbill.GetBillCategoryListRsp
	var err error

	rsp.List, rsp.Paginate, err = BillCate.FindPaginate(ctx, req.Options)
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}

	return &rsp, nil
}
func (a *LbbillServer) DelBill(ctx context.Context, req *lbbill.DelBillReq) (*lbbill.DelBillRsp, error) {
	var rsp lbbill.DelBillRsp
	var err error

	err = Bill.DeleteById(ctx, req.Id)
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}

	return &rsp, nil
}
func (a *LbbillServer) UpdateBill(ctx context.Context, req *lbbill.UpdateBillReq) (*lbbill.UpdateBillRsp, error) {
	var rsp lbbill.UpdateBillRsp
	var err error

	err = Bill.UpdateById(ctx, req.Bill.Id, utils.OrmStruct2Map(req.Bill))
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}

	return &rsp, nil
}
func (a *LbbillServer) GetBill(ctx context.Context, req *lbbill.GetBillReq) (*lbbill.GetBillRsp, error) {
	var rsp lbbill.GetBillRsp
	var err error

	rsp.Bill, err = Bill.GetById(ctx, req.Id)
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}

	return &rsp, nil
}
func (a *LbbillServer) GetBillList(ctx context.Context, req *lbbill.GetBillListReq) (*lbbill.GetBillListRsp, error) {
	var rsp lbbill.GetBillListRsp
	var err error

	rsp.List, rsp.Paginate, err = Bill.FindPaginate(ctx, req.Options)
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}

	categoryIdList := utils.PluckUint64List(rsp.List, lbbill.FieldCategoryId)
	billCategories, err := BillCate.GetByIdList(ctx, categoryIdList)
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}
	rsp.CategoryMap = utils.Slice2MapKeyByStructField(billCategories, lbbill.FieldId).(map[uint64]*lbbill.ModelBillCategory)

	return &rsp, nil
}
func (a *LbbillServer) GetBillCategory(ctx context.Context, req *lbbill.GetBillCategoryReq) (*lbbill.GetBillCategoryRsp, error) {
	var rsp lbbill.GetBillCategoryRsp
	var err error

	rsp.Category, err = BillCate.GetById(ctx, req.Id)
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}

	return &rsp, nil
}
