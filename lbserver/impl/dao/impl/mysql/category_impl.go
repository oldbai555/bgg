package mysql

import (
	"context"
	"github.com/mitchellh/mapstructure"
	"github.com/oldbai555/bgg/client/lbblog"
	"github.com/oldbai555/bgg/client/lbconst"
	"github.com/oldbai555/bgg/lbserver/impl/dao"
	webtool "github.com/oldbai555/bgg/pkg/webtoolv2"
	"github.com/oldbai555/gorm"
	"github.com/oldbai555/lbtool/log"
)

var _ dao.CategoryDao = (*CategoryImpl)(nil)

type CategoryImpl struct {
	mysqlConn
}

func (a *CategoryImpl) mustGetConn(ctx context.Context) *gorm.DB {
	return a.mysqlConn.mustGetConn(ctx).Model(&lbblog.ModelCategory{})
}

func (a *CategoryImpl) UpdateByIdList(ctx context.Context, idList []uint64, updateMap map[string]interface{}) error {
	db := webtool.NewCondBuilder(a.mustGetConn(ctx))
	_, err := db.In(lbblog.FieldId_, idList).Update(ctx, updateMap)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

func (a *CategoryImpl) DeleteByIdList(ctx context.Context, idList []uint64) error {
	db := webtool.NewCondBuilder(a.mustGetConn(ctx))
	_, err := db.In(lbblog.FieldId_, idList).Delete(ctx, &lbblog.ModelCategory{})
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

func (a *CategoryImpl) GetByIdList(ctx context.Context, idList []uint64) ([]*lbblog.ModelCategory, error) {
	var valList []*lbblog.ModelCategory
	db := webtool.NewCondBuilder(a.mustGetConn(ctx))
	err := db.In(lbblog.FieldId_, idList).Find(ctx, &valList)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return valList, nil
}

func (a *CategoryImpl) UpdateById(ctx context.Context, id uint64, updateMap map[string]interface{}) error {
	db := webtool.NewCondBuilder(a.mustGetConn(ctx))
	_, err := db.Eq(lbblog.FieldId_, id).Update(ctx, updateMap)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

func (a *CategoryImpl) UpdateOrCreate(ctx context.Context, candMap, attrMap map[string]interface{}, out *lbblog.ModelCategory) error {
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

func (a *CategoryImpl) BatchCreate(ctx context.Context, valList []*lbblog.ModelCategory) error {
	res := a.mustGetConn(ctx).CreateInBatches(valList, len(valList))
	log.Infof("batch create rows_affected %d", res.RowsAffected)
	if res.Error != nil {
		log.Errorf("err:%v", res.Error)
		return res.Error
	}
	return nil
}

func (a *CategoryImpl) Create(ctx context.Context, val *lbblog.ModelCategory) error {
	db := webtool.NewCondBuilder(a.mustGetConn(ctx))
	_, err := db.Create(ctx, val)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

func (a *CategoryImpl) DeleteById(ctx context.Context, id uint64) error {
	db := webtool.NewCondBuilder(a.mustGetConn(ctx))
	_, err := db.Eq(lbblog.FieldId_, id).Delete(ctx, &lbblog.ModelCategory{})
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

func (a *CategoryImpl) GetById(ctx context.Context, id uint64) (*lbblog.ModelCategory, error) {
	var val lbblog.ModelCategory
	db := webtool.NewCondBuilder(a.mustGetConn(ctx))
	err := db.Eq(lbblog.FieldId_, id).First(ctx, &val)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return &val, nil
}

func (a *CategoryImpl) FindPage(ctx context.Context, listOption *lbconst.ListOption) ([]*lbblog.ModelCategory, *lbconst.Page, error) {
	var list []*lbblog.ModelCategory
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

func NewCategoryImpl(ctx context.Context, dsn string) (dao.CategoryDao, error) {
	var d = &CategoryImpl{
		mysqlConn{
			dsn: dsn,
		},
	}
	err := d.mysqlConn.mustGetConn(ctx).Set(autoMigrateOptKey, autoMigrateOptValue).AutoMigrate(&lbblog.ModelCategory{})
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return d, nil
}