package mysql

import (
	"context"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/oldbai555/bgg/client/lb"
	"github.com/oldbai555/bgg/client/lbbill"
	"github.com/oldbai555/bgg/lbserver/impl/dao"
	"github.com/oldbai555/bgg/pkg/webtool"
	"github.com/oldbai555/gorm"
	"github.com/oldbai555/lbtool/log"
)

var _ dao.IBillDao = (*BillImpl)(nil)

type BillImpl struct {
	mysqlConn
}

func (a *BillImpl) mustGetConn(ctx context.Context) *gorm.DB {
	return a.mysqlConn.mustGetConn(ctx).Model(&lbbill.ModelBill{})
}

func (a *BillImpl) FirstOrCreate(ctx context.Context, val *lbbill.ModelBill, cand map[string]interface{}) error {
	return a.mysqlConn.mustGetConn(ctx).Model(&lbbill.ModelBill{}).FirstOrCreate(val, cand).Error
}

func (a *BillImpl) UpdateByIdList(ctx context.Context, idList []uint64, updateMap map[string]interface{}) error {
	db := webtool.NewCondBuilder(a.mustGetConn(ctx))
	_, err := db.In(lbbill.FieldId_, idList).Update(ctx, updateMap)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

func (a *BillImpl) DeleteByIdList(ctx context.Context, idList []uint64) error {
	db := webtool.NewCondBuilder(a.mustGetConn(ctx))
	_, err := db.In(lbbill.FieldId_, idList).Delete(ctx, &lbbill.ModelBill{})
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

func (a *BillImpl) GetByIdList(ctx context.Context, idList []uint64) ([]*lbbill.ModelBill, error) {
	var valList []*lbbill.ModelBill
	db := webtool.NewCondBuilder(a.mustGetConn(ctx))
	err := db.In(lbbill.FieldId_, idList).Find(ctx, &valList)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return valList, nil
}

func (a *BillImpl) UpdateById(ctx context.Context, id uint64, updateMap map[string]interface{}) error {
	db := webtool.NewCondBuilder(a.mustGetConn(ctx))
	_, err := db.Eq(lbbill.FieldId_, id).Update(ctx, updateMap)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

func (a *BillImpl) UpdateOrCreate(ctx context.Context, candMap, attrMap map[string]interface{}, out *lbbill.ModelBill) error {
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

func (a *BillImpl) BatchCreate(ctx context.Context, valList []*lbbill.ModelBill) error {
	res := a.mustGetConn(ctx).CreateInBatches(valList, len(valList))
	log.Infof("batch create rows_affected %d", res.RowsAffected)
	if res.Error != nil {
		log.Errorf("err:%v", res.Error)
		return res.Error
	}
	return nil
}

func (a *BillImpl) Create(ctx context.Context, val *lbbill.ModelBill) error {
	db := webtool.NewCondBuilder(a.mustGetConn(ctx))
	_, err := db.Create(ctx, val)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

func (a *BillImpl) DeleteById(ctx context.Context, id uint64) error {
	db := webtool.NewCondBuilder(a.mustGetConn(ctx))
	_, err := db.Eq(lbbill.FieldId_, id).Delete(ctx, &lbbill.ModelBill{})
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

func (a *BillImpl) GetById(ctx context.Context, id uint64) (*lbbill.ModelBill, error) {
	var val lbbill.ModelBill
	db := webtool.NewCondBuilder(a.mustGetConn(ctx))
	err := db.Eq(lbbill.FieldId_, id).First(ctx, &val)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return &val, nil
}

func (a *BillImpl) FindPaginate(ctx context.Context, options *lb.Options) ([]*lbbill.ModelBill, *lb.Paginate, error) {
	var list []*lbbill.ModelBill
	db := webtool.NewList(a.mustGetConn(ctx), options)
	err := webtool.ProcessDefaultOptions(options, db)
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, nil, err
	}
	err = lb.NewOptionsProcessor(options).
		AddUint64(lbbill.GetBillListReq_OptionCategoryId, func(val uint64) error {
			db.Eq(lbbill.FieldCategoryId_, val)
			return nil
		}).
		AddUint32(lbbill.GetBillListReq_OptionRootCategory, func(val uint32) error {
			db.Eq(lbbill.FieldRootCategory_, val)
			return nil
		}).
		AddTimeStampRange(lbbill.GetBillListReq_OptionDateUnixRange, func(beginAt, endAt uint32) error {
			db.Between(lbbill.FieldDateUnix_, beginAt, endAt)
			return nil
		}).
		AddString(lbbill.GetBillListReq_OptionLikeRemark, func(val string) error {
			db.Like(lbbill.FieldRemark_, fmt.Sprintf("%%%s%%", val))
			return nil
		}).
		Process()
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, nil, err
	}

	paginate, err := db.FindPaginate(ctx, &list)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, nil, err
	}
	return list, paginate, nil
}

func NewBillImpl(ctx context.Context, dsn string) (dao.IBillDao, error) {
	var d = &BillImpl{
		mysqlConn{
			dsn: dsn,
		},
	}

	err := d.mysqlConn.mustGetConn(ctx).Set(autoMigrateOptKey, autoMigrateOptValue).AutoMigrate(&lbbill.ModelBill{})
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return d, nil
}
