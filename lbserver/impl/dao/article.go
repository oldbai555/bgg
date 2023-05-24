package dao

import (
	"context"
	"github.com/oldbai555/bgg/client/lb"
	"github.com/oldbai555/bgg/client/lbblog"
)

// IArticleDao 目的是收回 SQL 的使用范围 避免 SQL 满天飞
type IArticleDao interface {
	// 基础SQL - 自动生成

	Create(ctx context.Context, val *lbblog.ModelArticle) error
	FirstOrCreate(ctx context.Context, val *lbblog.ModelArticle, cand map[string]interface{}) error
	BatchCreate(ctx context.Context, valList []*lbblog.ModelArticle) error
	UpdateOrCreate(ctx context.Context, candMap, attrMap map[string]interface{}, out *lbblog.ModelArticle) error
	UpdateById(ctx context.Context, id uint64, updateMap map[string]interface{}) error
	UpdateByIdList(ctx context.Context, idList []uint64, updateMap map[string]interface{}) error
	DeleteById(ctx context.Context, id uint64) error
	DeleteByIdList(ctx context.Context, idList []uint64) error
	GetById(ctx context.Context, id uint64) (*lbblog.ModelArticle, error)
	GetByIdList(ctx context.Context, idList []uint64) ([]*lbblog.ModelArticle, error)
	FindPaginate(ctx context.Context, options *lb.Options) ([]*lbblog.ModelArticle, *lb.Paginate, error)

	// 向下扩展业务SQL ......
}
