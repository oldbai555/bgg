package dao

import (
	"context"
	"github.com/oldbai555/bgg/client/lbconst"
	"github.com/oldbai555/bgg/client/lbim"
)

// MessageDao 目的是收回 SQL 的使用范围 避免 SQL 满天飞
type MessageDao interface {
	// 基础SQL - 自动生成

	Create(ctx context.Context, val *lbim.ModelMessage) error
	BatchCreate(ctx context.Context, valList []*lbim.ModelMessage) error
	UpdateOrCreate(ctx context.Context, candMap, attrMap map[string]interface{}, out *lbim.ModelMessage) error
	UpdateById(ctx context.Context, id uint64, updateMap map[string]interface{}) error
	UpdateByIdList(ctx context.Context, idList []uint64, updateMap map[string]interface{}) error
	DeleteById(ctx context.Context, id uint64) error
	DeleteByIdList(ctx context.Context, idList []uint64) error
	GetById(ctx context.Context, id uint64) (*lbim.ModelMessage, error)
	GetByIdList(ctx context.Context, idList []uint64) ([]*lbim.ModelMessage, error)
	FindPage(ctx context.Context, listOption *lbconst.ListOption) ([]*lbim.ModelMessage, *lbconst.Page, error)

	// 向下扩展业务SQL ......
}
