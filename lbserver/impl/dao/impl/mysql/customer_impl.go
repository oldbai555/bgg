package mysql

import (
	"context"
	"github.com/mitchellh/mapstructure"
	"github.com/oldbai555/bgg/client/lbconst"
	"github.com/oldbai555/bgg/client/lbcustomer"
	"github.com/oldbai555/bgg/lbserver/impl/dao"
	webtool "github.com/oldbai555/bgg/pkg/webtoolv2"
	"github.com/oldbai555/gorm"
	"github.com/oldbai555/lbtool/log"
	"sync/atomic"
)

var _ dao.CustomerDao = (*CustomerImpl)(nil)

var migratedCustomer atomic.Bool

type CustomerImpl struct {
	mysqlConn
}

func (a *CustomerImpl) UpdateByIdList(ctx context.Context, idList []uint64, updateMap map[string]interface{}) error {
	db := webtool.NewCondBuilder(a.mustGetConn(ctx))
	_, err := db.In(lbcustomer.FieldId_, idList).Update(ctx, updateMap)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

func (a *CustomerImpl) DeleteByIdList(ctx context.Context, idList []uint64) error {
	db := webtool.NewCondBuilder(a.mustGetConn(ctx))
	_, err := db.In(lbcustomer.FieldId_, idList).Delete(ctx, &lbcustomer.ModelCustomer{})
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

func (a *CustomerImpl) GetByIdList(ctx context.Context, idList []uint64) ([]*lbcustomer.ModelCustomer, error) {
	var valList []*lbcustomer.ModelCustomer
	db := webtool.NewCondBuilder(a.mustGetConn(ctx))
	err := db.In(lbcustomer.FieldId_, idList).Find(ctx, &valList)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return valList, nil
}

func (a *CustomerImpl) UpdateById(ctx context.Context, id uint64, updateMap map[string]interface{}) error {
	db := webtool.NewCondBuilder(a.mustGetConn(ctx))
	_, err := db.Eq(lbcustomer.FieldId_, id).Update(ctx, updateMap)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

func (a *CustomerImpl) UpdateOrCreate(ctx context.Context, candMap, attrMap map[string]interface{}, out *lbcustomer.ModelCustomer) error {
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

func (a *CustomerImpl) BatchCreate(ctx context.Context, valList []*lbcustomer.ModelCustomer) error {
	res := a.mustGetConn(ctx).CreateInBatches(valList, len(valList))
	log.Infof("batch create rows_affected %d", res.RowsAffected)
	if res.Error != nil {
		log.Errorf("err:%v", res.Error)
		return res.Error
	}
	return nil
}

func (a *CustomerImpl) Create(ctx context.Context, val *lbcustomer.ModelCustomer) error {
	db := webtool.NewCondBuilder(a.mustGetConn(ctx))
	_, err := db.Create(ctx, val)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

func (a *CustomerImpl) DeleteById(ctx context.Context, id uint64) error {
	db := webtool.NewCondBuilder(a.mustGetConn(ctx))
	_, err := db.Eq(lbcustomer.FieldId_, id).Delete(ctx, &lbcustomer.ModelCustomer{})
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

func (a *CustomerImpl) GetById(ctx context.Context, id uint64) (*lbcustomer.ModelCustomer, error) {
	var val lbcustomer.ModelCustomer
	db := webtool.NewCondBuilder(a.mustGetConn(ctx))
	err := db.Eq(lbcustomer.FieldId_, id).First(ctx, &val)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return &val, nil
}

func (a *CustomerImpl) FindPage(ctx context.Context, listOption *lbconst.ListOption) ([]*lbcustomer.ModelCustomer, *lbconst.Page, error) {
	var list []*lbcustomer.ModelCustomer
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

func NewCustomerImpl(ctx context.Context, dsn string) (dao.CustomerDao, error) {
	var d = &CustomerImpl{
		mysqlConn{
			dsn: dsn,
		},
	}
	if !migratedCustomer.Load() {
		err := d.mustGetConn(ctx).Set(autoMigrateOptKey, autoMigrateOptValue).AutoMigrate(&lbcustomer.ModelCustomer{})
		if err != nil {
			log.Errorf("err:%v", err)
			return nil, err
		}
	}
	return d, nil
}
