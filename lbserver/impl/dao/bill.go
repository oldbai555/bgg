package dao

import (
	"context"
	"github.com/oldbai555/bgg/client/lb"
	"github.com/oldbai555/bgg/client/lbbill"
)

// IBillDao 目的是收回 SQL 的使用范围 避免 SQL 满天飞
type IBillDao interface {
	// 基础SQL - 自动生成

	Create(ctx context.Context, val *lbbill.ModelBill) error
	FirstOrCreate(ctx context.Context, val *lbbill.ModelBill, cand map[string]interface{}) error
	BatchCreate(ctx context.Context, valList []*lbbill.ModelBill) error
	UpdateOrCreate(ctx context.Context, candMap, attrMap map[string]interface{}, out *lbbill.ModelBill) error
	UpdateById(ctx context.Context, id uint64, updateMap map[string]interface{}) error
	UpdateByIdList(ctx context.Context, idList []uint64, updateMap map[string]interface{}) error
	DeleteById(ctx context.Context, id uint64) error
	DeleteByIdList(ctx context.Context, idList []uint64) error
	GetById(ctx context.Context, id uint64) (*lbbill.ModelBill, error)
	GetByIdList(ctx context.Context, idList []uint64) ([]*lbbill.ModelBill, error)
	FindPaginate(ctx context.Context, options *lb.Options) ([]*lbbill.ModelBill, *lb.Paginate, error)

	// 向下扩展业务SQL ......
}
