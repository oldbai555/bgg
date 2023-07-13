package dao

import (
	"context"
	"github.com/oldbai555/bgg/client/lb"
	"github.com/oldbai555/bgg/client/lbbill"
	"github.com/oldbai555/bgg/pkg/webtool"
	"github.com/oldbai555/gorm"
)

// IBillDao 目的是收回 SQL 的使用范围 避免 SQL 满天飞
type IBillDao interface {
	// 基础SQL - 自动生成

	GetOrmEngine(ctx context.Context) *gorm.DB
	Create(ctx context.Context, val *lbbill.ModelBill) (*webtool.Result, error)
	FirstOrCreate(ctx context.Context, cand map[string]interface{}, val *lbbill.ModelBill) (*webtool.Result, error)
	BatchCreate(ctx context.Context, valList []*lbbill.ModelBill) (*webtool.Result, error)
	UpdateOrCreate(ctx context.Context, candMap, attrMap map[string]interface{}, out *lbbill.ModelBill) (*webtool.Result, error)
	UpdateById(ctx context.Context, id uint64, updateMap map[string]interface{}) (*webtool.Result, error)
	UpdateByIdList(ctx context.Context, idList []uint64, updateMap map[string]interface{}) (*webtool.Result, error)
	DeleteById(ctx context.Context, id uint64) (*webtool.Result, error)
	DeleteByIdList(ctx context.Context, idList []uint64) (*webtool.Result, error)
	GetById(ctx context.Context, id uint64) (*lbbill.ModelBill, error)
	GetByIdList(ctx context.Context, idList []uint64) ([]*lbbill.ModelBill, error)
	FindPaginate(ctx context.Context, options *lb.Options) ([]*lbbill.ModelBill, *lb.Paginate, error)
	Increment(ctx context.Context, field string, num uint32, candMap map[string]interface{}) error
	Decrement(ctx context.Context, field string, num uint32, candMap map[string]interface{}) error
	IsNotFoundErr(err error) bool
	GetList(ctx context.Context, candMap map[string]interface{}, opts ...*webtool.Opt) ([]*lbbill.ModelBill, error)
	GetOne(ctx context.Context, candMap map[string]interface{}, opts ...*webtool.Opt) (*lbbill.ModelBill, error)
	GetOrmCondBuilder(ctx context.Context) *webtool.OrmCondBuilder

	// 向下扩展业务SQL ......
}
