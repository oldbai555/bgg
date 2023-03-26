package mysql

import (
	"context"
	"github.com/mitchellh/mapstructure"
	"github.com/oldbai555/bgg/client/lbconst"
	"github.com/oldbai555/bgg/client/lbim"
	"github.com/oldbai555/bgg/lbserver/impl/dao"
	webtool "github.com/oldbai555/bgg/pkg/webtoolv2"
	"github.com/oldbai555/gorm"
	"github.com/oldbai555/lbtool/log"
)

var _ dao.ChatDao = (*ChatImpl)(nil)

type ChatImpl struct {
	mysqlConn
}

func (a *ChatImpl) mustGetConn(ctx context.Context) *gorm.DB {
	return a.mysqlConn.mustGetConn(ctx).Model(&lbim.ModelChat{})
}

func (a *ChatImpl) FirstOrCreate(ctx context.Context, val *lbim.ModelChat, cand map[string]interface{}) error {
	return a.mysqlConn.mustGetConn(ctx).Model(&lbim.ModelChat{}).FirstOrCreate(val, cand).Error
}

func (a *ChatImpl) UpdateByIdList(ctx context.Context, idList []uint64, updateMap map[string]interface{}) error {
	db := webtool.NewCondBuilder(a.mustGetConn(ctx))
	_, err := db.In(lbim.FieldId_, idList).Update(ctx, updateMap)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

func (a *ChatImpl) DeleteByIdList(ctx context.Context, idList []uint64) error {
	db := webtool.NewCondBuilder(a.mustGetConn(ctx))
	_, err := db.In(lbim.FieldId_, idList).Delete(ctx, &lbim.ModelChat{})
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

func (a *ChatImpl) GetByIdList(ctx context.Context, idList []uint64) ([]*lbim.ModelChat, error) {
	var valList []*lbim.ModelChat
	db := webtool.NewCondBuilder(a.mustGetConn(ctx))
	err := db.In(lbim.FieldId_, idList).Find(ctx, &valList)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return valList, nil
}

func (a *ChatImpl) UpdateById(ctx context.Context, id uint64, updateMap map[string]interface{}) error {
	db := webtool.NewCondBuilder(a.mustGetConn(ctx))
	_, err := db.Eq(lbim.FieldId_, id).Update(ctx, updateMap)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

func (a *ChatImpl) UpdateOrCreate(ctx context.Context, candMap, attrMap map[string]interface{}, out *lbim.ModelChat) error {
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

func (a *ChatImpl) BatchCreate(ctx context.Context, valList []*lbim.ModelChat) error {
	res := a.mustGetConn(ctx).CreateInBatches(valList, len(valList))
	log.Infof("batch create rows_affected %d", res.RowsAffected)
	if res.Error != nil {
		log.Errorf("err:%v", res.Error)
		return res.Error
	}
	return nil
}

func (a *ChatImpl) Create(ctx context.Context, val *lbim.ModelChat) error {
	db := webtool.NewCondBuilder(a.mustGetConn(ctx))
	_, err := db.Create(ctx, val)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

func (a *ChatImpl) DeleteById(ctx context.Context, id uint64) error {
	db := webtool.NewCondBuilder(a.mustGetConn(ctx))
	_, err := db.Eq(lbim.FieldId_, id).Delete(ctx, &lbim.ModelChat{})
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

func (a *ChatImpl) GetById(ctx context.Context, id uint64) (*lbim.ModelChat, error) {
	var val lbim.ModelChat
	db := webtool.NewCondBuilder(a.mustGetConn(ctx))
	err := db.Eq(lbim.FieldId_, id).First(ctx, &val)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return &val, nil
}

func (a *ChatImpl) FindPage(ctx context.Context, listOption *lbconst.ListOption) ([]*lbim.ModelChat, *lbconst.Page, error) {
	var list []*lbim.ModelChat
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

func NewChatImpl(ctx context.Context, dsn string) (dao.ChatDao, error) {
	var d = &ChatImpl{
		mysqlConn{
			dsn: dsn,
		},
	}

	err := d.mysqlConn.mustGetConn(ctx).Set(autoMigrateOptKey, autoMigrateOptValue).AutoMigrate(&lbim.ModelChat{})
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return d, nil
}
