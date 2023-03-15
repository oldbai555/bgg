package dao

import (
	"context"
	"github.com/oldbai555/bgg/client/lbconst"
	"github.com/oldbai555/bgg/client/lbcustomer"
)

// CustomerDao 目的是收回 SQL 的使用范围 避免 SQL 满天飞
type CustomerDao interface {
	// 基础SQL - 自动生成

	Create(ctx context.Context, val *lbcustomer.ModelCustomer) error
	BatchCreate(ctx context.Context, valList []*lbcustomer.ModelCustomer) error
	UpdateOrCreate(ctx context.Context, candMap, attrMap map[string]interface{}, out *lbcustomer.ModelCustomer) error
	UpdateById(ctx context.Context, id uint64, updateMap map[string]interface{}) error
	UpdateByIdList(ctx context.Context, idList []uint64, updateMap map[string]interface{}) error
	DeleteById(ctx context.Context, id uint64) error
	DeleteByIdList(ctx context.Context, idList []uint64) error
	GetById(ctx context.Context, id uint64) (*lbcustomer.ModelCustomer, error)
	GetByIdList(ctx context.Context, idList []uint64) ([]*lbcustomer.ModelCustomer, error)
	FindPage(ctx context.Context, listOption *lbconst.ListOption) ([]*lbcustomer.ModelCustomer, *lbconst.Page, error)

	// 向下扩展业务SQL ......
}
