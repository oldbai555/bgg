package dao

import (
	"context"
	"github.com/oldbai555/bgg/client/lb"
	"github.com/oldbai555/bgg/client/lbblog"
)

// ICommentDao 目的是收回 SQL 的使用范围 避免 SQL 满天飞
type ICommentDao interface {
	// 基础SQL - 自动生成

	Create(ctx context.Context, val *lbblog.ModelComment) error
	FirstOrCreate(ctx context.Context, val *lbblog.ModelComment, cand map[string]interface{}) error
	BatchCreate(ctx context.Context, valList []*lbblog.ModelComment) error
	UpdateOrCreate(ctx context.Context, candMap, attrMap map[string]interface{}, out *lbblog.ModelComment) error
	UpdateById(ctx context.Context, id uint64, updateMap map[string]interface{}) error
	UpdateByIdList(ctx context.Context, idList []uint64, updateMap map[string]interface{}) error
	DeleteById(ctx context.Context, id uint64) error
	DeleteByIdList(ctx context.Context, idList []uint64) error
	GetById(ctx context.Context, id uint64) (*lbblog.ModelComment, error)
	GetByIdList(ctx context.Context, idList []uint64) ([]*lbblog.ModelComment, error)
	FindPaginate(ctx context.Context, options *lb.Options) ([]*lbblog.ModelComment, *lb.Paginate, error)

	// 向下扩展业务SQL ......
}
