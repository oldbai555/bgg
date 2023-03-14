package dao

import (
	"context"
	"github.com/oldbai555/bgg/client/lbaccount"
	"github.com/oldbai555/bgg/client/lbconst"
)

// AccountDao 目的是收回 SQL 的使用范围 避免 SQL 满天飞
type AccountDao interface {
	// 基础SQL - 自动生成

	Create(ctx context.Context, account *lbaccount.ModelAccount) error
	DeleteById(ctx context.Context, id uint64) error
	GetById(ctx context.Context, id uint64) (*lbaccount.ModelAccount, error)
	FindPage(ctx context.Context, listOption *lbconst.ListOption) ([]*lbaccount.ModelAccount, *lbconst.Page, error)

	// 向下扩展业务SQL ......
}
