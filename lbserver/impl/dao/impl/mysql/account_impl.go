package mysql

import (
	"context"
	"github.com/mitchellh/mapstructure"
	"github.com/oldbai555/bgg/client/lbaccount"
	"github.com/oldbai555/bgg/client/lbconst"
	"github.com/oldbai555/bgg/lbserver/impl/dao"
	webtool "github.com/oldbai555/bgg/pkg/webtoolv2"
	"github.com/oldbai555/gorm"
	"github.com/oldbai555/lbtool/log"
)

var _ dao.AccountDao = (*AccountImpl)(nil)

type AccountImpl struct {
	mysqlConn
}

func (a *AccountImpl) mustGetConn(ctx context.Context) *gorm.DB {
	return a.mysqlConn.mustGetConn(ctx).Model(&lbaccount.ModelAccount{})
}

func (a *AccountImpl) UpdateByIdList(ctx context.Context, idList []uint64, updateMap map[string]interface{}) error {
	db := webtool.NewCondBuilder(a.mustGetConn(ctx))
	_, err := db.In(lbaccount.FieldId_, idList).Update(ctx, updateMap)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

func (a *AccountImpl) DeleteByIdList(ctx context.Context, idList []uint64) error {
	db := webtool.NewCondBuilder(a.mustGetConn(ctx))
	_, err := db.In(lbaccount.FieldId_, idList).Delete(ctx, &lbaccount.ModelAccount{})
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

func (a *AccountImpl) GetByIdList(ctx context.Context, idList []uint64) ([]*lbaccount.ModelAccount, error) {
	var valList []*lbaccount.ModelAccount
	db := webtool.NewCondBuilder(a.mustGetConn(ctx))
	err := db.In(lbaccount.FieldId_, idList).Find(ctx, &valList)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return valList, nil
}

func (a *AccountImpl) UpdateById(ctx context.Context, id uint64, updateMap map[string]interface{}) error {
	db := webtool.NewCondBuilder(a.mustGetConn(ctx))
	_, err := db.Eq(lbaccount.FieldId_, id).Update(ctx, updateMap)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

func (a *AccountImpl) UpdateOrCreate(ctx context.Context, candMap, attrMap map[string]interface{}, out *lbaccount.ModelAccount) error {
	selectDb := webtool.NewCondBuilder(a.mustGetConn(ctx))
	err := selectDb.AndMap(candMap).First(ctx, out)
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Errorf("err:%v", err)
		return err
	}

	optDb := webtool.NewCondBuilder(a.mustGetConn(ctx))
	if err == nil {
		_, err := optDb.AndMap(candMap).Update(ctx, attrMap)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
	}

	if err == gorm.ErrRecordNotFound {
		for k, v := range candMap {
			attrMap[k] = v
		}
		err = mapstructure.Decode(attrMap, out)
		if err != nil {
			log.Errorf("err is : %v", err)
			return err
		}
		_, err = optDb.Create(ctx, out)
		if err != nil {
			log.Errorf("err is %v", err)
			return err
		}
		return nil
	}
	return nil
}

func (a *AccountImpl) BatchCreate(ctx context.Context, valList []*lbaccount.ModelAccount) error {
	res := a.mustGetConn(ctx).CreateInBatches(valList, len(valList))
	log.Infof("batch create rows_affected %d", res.RowsAffected)
	if res.Error != nil {
		log.Errorf("err:%v", res.Error)
		return res.Error
	}
	return nil
}

func (a *AccountImpl) Create(ctx context.Context, val *lbaccount.ModelAccount) error {
	db := webtool.NewCondBuilder(a.mustGetConn(ctx))
	_, err := db.Create(ctx, val)
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
	var val lbaccount.ModelAccount
	db := webtool.NewCondBuilder(a.mustGetConn(ctx))
	err := db.Eq(lbaccount.FieldId_, id).First(ctx, &val)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return &val, nil
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

	err := d.mysqlConn.mustGetConn(ctx).Set(autoMigrateOptKey, autoMigrateOptValue).AutoMigrate(&lbaccount.ModelAccount{})
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return d, nil
}
