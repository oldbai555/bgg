package service

import (
	"context"
	"github.com/oldbai555/bgg/client/lbbill"
)

var BillServer LbbillServer

type LbbillServer struct {
	*lbbill.UnimplementedLbbillServer
}

func (a *LbbillServer) AddBill(ctx context.Context, req *lbbill.AddBillReq) (*lbbill.AddBillRsp, error) {
	var rsp lbbill.AddBillRsp
	var err error

	return &rsp, err
}
func (a *LbbillServer) DelBill(ctx context.Context, req *lbbill.DelBillReq) (*lbbill.DelBillRsp, error) {
	var rsp lbbill.DelBillRsp
	var err error

	return &rsp, err
}
func (a *LbbillServer) UpdateBill(ctx context.Context, req *lbbill.UpdateBillReq) (*lbbill.UpdateBillRsp, error) {
	var rsp lbbill.UpdateBillRsp
	var err error

	return &rsp, err
}
func (a *LbbillServer) GetBill(ctx context.Context, req *lbbill.GetBillReq) (*lbbill.GetBillRsp, error) {
	var rsp lbbill.GetBillRsp
	var err error

	return &rsp, err
}
func (a *LbbillServer) GetBillList(ctx context.Context, req *lbbill.GetBillListReq) (*lbbill.GetBillListRsp, error) {
	var rsp lbbill.GetBillListRsp
	var err error

	return &rsp, err
}
func (a *LbbillServer) AddBillCategory(ctx context.Context, req *lbbill.AddBillCategoryReq) (*lbbill.AddBillCategoryRsp, error) {
	var rsp lbbill.AddBillCategoryRsp
	var err error

	return &rsp, err
}
func (a *LbbillServer) DelBillCategory(ctx context.Context, req *lbbill.DelBillCategoryReq) (*lbbill.DelBillCategoryRsp, error) {
	var rsp lbbill.DelBillCategoryRsp
	var err error

	return &rsp, err
}
func (a *LbbillServer) UpdateBillCategory(ctx context.Context, req *lbbill.UpdateBillCategoryReq) (*lbbill.UpdateBillCategoryRsp, error) {
	var rsp lbbill.UpdateBillCategoryRsp
	var err error

	return &rsp, err
}
func (a *LbbillServer) GetBillCategory(ctx context.Context, req *lbbill.GetBillCategoryReq) (*lbbill.GetBillCategoryRsp, error) {
	var rsp lbbill.GetBillCategoryRsp
	var err error

	return &rsp, err
}
func (a *LbbillServer) GetBillCategoryList(ctx context.Context, req *lbbill.GetBillCategoryListReq) (*lbbill.GetBillCategoryListRsp, error) {
	var rsp lbbill.GetBillCategoryListRsp
	var err error

	return &rsp, err
}
