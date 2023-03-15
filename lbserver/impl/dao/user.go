package dao

import (
	"context"
	"github.com/oldbai555/bgg/client/lbconst"
	"github.com/oldbai555/bgg/client/lbuser"
)

// UserDao 目的是收回 SQL 的使用范围 避免 SQL 满天飞
type UserDao interface {
	// 基础SQL - 自动生成

	Create(ctx context.Context, val *lbuser.ModelUser) error
	BatchCreate(ctx context.Context, valList []*lbuser.ModelUser) error
	UpdateOrCreate(ctx context.Context, candMap, attrMap map[string]interface{}, out *lbuser.ModelUser) error
	UpdateById(ctx context.Context, id uint64, updateMap map[string]interface{}) error
	UpdateByIdList(ctx context.Context, idList []uint64, updateMap map[string]interface{}) error
	DeleteById(ctx context.Context, id uint64) error
	DeleteByIdList(ctx context.Context, idList []uint64) error
	GetById(ctx context.Context, id uint64) (*lbuser.ModelUser, error)
	GetByIdList(ctx context.Context, idList []uint64) ([]*lbuser.ModelUser, error)
	FindPage(ctx context.Context, listOption *lbconst.ListOption) ([]*lbuser.ModelUser, *lbconst.Page, error)

	// 向下扩展业务SQL ......
	GetByUserName(ctx context.Context, username string) (*lbuser.ModelUser, error)
	GetByAdmin(ctx context.Context) (*lbuser.ModelUser, error)
	CheckUserNameExit(ctx context.Context, id uint64, username string) (bool, error)
}
