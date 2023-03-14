package mysql

import (
	"context"
	"github.com/oldbai555/bgg/client/lbaccount"
	"github.com/oldbai555/bgg/client/lbconst"
	"github.com/oldbai555/bgg/lbserver/impl/dao"
	webtool "github.com/oldbai555/bgg/pkg/webtoolv2"
	"github.com/oldbai555/lbtool/log"
	"sync/atomic"
)

var _ dao.AccountDao = (*AccountImpl)(nil)

var migratedAccount atomic.Bool

type AccountImpl struct {
	mysqlConn
}

func (a *AccountImpl) Create(ctx context.Context, account *lbaccount.ModelAccount) error {
	db := webtool.NewCondBuilder(a.mustGetConn(ctx))
	_, err := db.Create(ctx, account)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

func (a *AccountImpl) DeleteById(ctx context.Context, id uint64) error {
	db := webtool.NewCondBuilder(a.mustGetConn(ctx))
	_, err := db.Eq(lbaccount.FieldId_, id).Delete(ctx, &lbaccount.ModelAccount{})
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

func (a *AccountImpl) GetById(ctx context.Context, id uint64) (*lbaccount.ModelAccount, error) {
	var account lbaccount.ModelAccount
	db := webtool.NewCondBuilder(a.mustGetConn(ctx))
	err := db.Eq(lbaccount.FieldId_, id).First(ctx, &lbaccount.ModelAccount{})
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return &account, nil
}

func (a *AccountImpl) FindPage(ctx context.Context, listOption *lbconst.ListOption) ([]*lbaccount.ModelAccount, *lbconst.Page, error) {
	var list []*lbaccount.ModelAccount
	db := webtool.NewList(a.mustGetConn(ctx), listOption)
	err := lbconst.NewListOptionProcessor(listOption).
		Process()
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, nil, err
	}

	page, err := db.FindPage(ctx, &list)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, nil, err
	}
	return list, page, nil
}

func NewAccountImpl(ctx context.Context, dsn string) (dao.AccountDao, error) {
	var d = &AccountImpl{
		mysqlConn{
			dsn: dsn,
		},
	}
	if !migratedAccount.Load() {
		err := d.mustGetConn(ctx).Set(autoMigrateOptKey, autoMigrateOptValue).AutoMigrate(&lbaccount.ModelAccount{})
		if err != nil {
			log.Errorf("err:%v", err)
			return nil, err
		}
	}
	return d, nil
}
