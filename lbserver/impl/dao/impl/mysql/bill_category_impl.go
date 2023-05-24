package mysql

import (
	"context"
	"github.com/mitchellh/mapstructure"
	"github.com/oldbai555/bgg/client/lb"
	"github.com/oldbai555/bgg/client/lbbill"
	"github.com/oldbai555/bgg/lbserver/impl/dao"
	"github.com/oldbai555/bgg/pkg/webtool"
	"github.com/oldbai555/gorm"
	"github.com/oldbai555/lbtool/log"
)

var _ dao.IBillCategoryDao = (*BillCategoryImpl)(nil)

type BillCategoryImpl struct {
	mysqlConn
}

func (a *BillCategoryImpl) mustGetConn(ctx context.Context) *gorm.DB {
	return a.mysqlConn.mustGetConn(ctx).Model(&lbbill.ModelBillCategory{})
}

func (a *BillCategoryImpl) FirstOrCreate(ctx context.Context, val *lbbill.ModelBillCategory, cand map[string]interface{}) error {
	return a.mysqlConn.mustGetConn(ctx).Model(&lbbill.ModelBillCategory{}).FirstOrCreate(val, cand).Error
}

func (a *BillCategoryImpl) UpdateByIdList(ctx context.Context, idList []uint64, updateMap map[string]interface{}) error {
	db := webtool.NewCondBuilder(a.mustGetConn(ctx))
	_, err := db.In(lbbill.FieldId_, idList).Update(ctx, updateMap)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

func (a *BillCategoryImpl) DeleteByIdList(ctx context.Context, idList []uint64) error {
	db := webtool.NewCondBuilder(a.mustGetConn(ctx))
	_, err := db.In(lbbill.FieldId_, idList).Delete(ctx, &lbbill.ModelBillCategory{})
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

func (a *BillCategoryImpl) GetByIdList(ctx context.Context, idList []uint64) ([]*lbbill.ModelBillCategory, error) {
	var valList []*lbbill.ModelBillCategory
	db := webtool.NewCondBuilder(a.mustGetConn(ctx))
	err := db.In(lbbill.FieldId_, idList).Find(ctx, &valList)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return valList, nil
}

func (a *BillCategoryImpl) UpdateById(ctx context.Context, id uint64, updateMap map[string]interface{}) error {
	db := webtool.NewCondBuilder(a.mustGetConn(ctx))
	_, err := db.Eq(lbbill.FieldId_, id).Update(ctx, updateMap)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

func (a *BillCategoryImpl) UpdateOrCreate(ctx context.Context, candMap, attrMap map[string]interface{}, out *lbbill.ModelBillCategory) error {
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

func (a *BillCategoryImpl) BatchCreate(ctx context.Context, valList []*lbbill.ModelBillCategory) error {
	res := a.mustGetConn(ctx).CreateInBatches(valList, len(valList))
	log.Infof("batch create rows_affected %d", res.RowsAffected)
	if res.Error != nil {
		log.Errorf("err:%v", res.Error)
		return res.Error
	}
	return nil
}

func (a *BillCategoryImpl) Create(ctx context.Context, val *lbbill.ModelBillCategory) error {
	db := webtool.NewCondBuilder(a.mustGetConn(ctx))
	_, err := db.Create(ctx, val)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

func (a *BillCategoryImpl) DeleteById(ctx context.Context, id uint64) error {
	db := webtool.NewCondBuilder(a.mustGetConn(ctx))
	_, err := db.Eq(lbbill.FieldId_, id).Delete(ctx, &lbbill.ModelBillCategory{})
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

func (a *BillCategoryImpl) GetById(ctx context.Context, id uint64) (*lbbill.ModelBillCategory, error) {
	var val lbbill.ModelBillCategory
	db := webtool.NewCondBuilder(a.mustGetConn(ctx))
	err := db.Eq(lbbill.FieldId_, id).First(ctx, &val)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return &val, nil
}

func (a *BillCategoryImpl) FindPaginate(ctx context.Context, options *lb.Options) ([]*lbbill.ModelBillCategory, *lb.Paginate, error) {
	var list []*lbbill.ModelBillCategory
	db := webtool.NewList(a.mustGetConn(ctx), options)
	err := webtool.ProcessDefaultOptions(options, db)
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, nil, err
	}
	err = lb.NewOptionsProcessor(options).
		AddString(lbbill.GetBillCategoryListReq_OptionLikeName, func(val string) error {
			db.Like(lbbill.FieldName_, val)
			return nil
		}).
		AddUint32(lbbill.GetBillCategoryListReq_OptionRootCategory, func(val uint32) error {
			db.Eq(lbbill.FieldRootCategory_, val)
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

func NewBillCategoryImpl(ctx context.Context, dsn string) (dao.IBillCategoryDao, error) {
	var d = &BillCategoryImpl{
		mysqlConn{
			dsn: dsn,
		},
	}

	err := d.mysqlConn.mustGetConn(ctx).Set(autoMigrateOptKey, autoMigrateOptValue).AutoMigrate(&lbbill.ModelBillCategory{})
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return d, nil
}
