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
	"sync/atomic"
)

var _ dao.MessageDao = (*MessageImpl)(nil)

var migratedMessage atomic.Bool

type MessageImpl struct {
	mysqlConn
}

func (a *MessageImpl) UpdateByIdList(ctx context.Context, idList []uint64, updateMap map[string]interface{}) error {
	db := webtool.NewCondBuilder(a.mustGetConn(ctx))
	_, err := db.In(lbim.FieldId_, idList).Update(ctx, updateMap)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

func (a *MessageImpl) DeleteByIdList(ctx context.Context, idList []uint64) error {
	db := webtool.NewCondBuilder(a.mustGetConn(ctx))
	_, err := db.In(lbim.FieldId_, idList).Delete(ctx, &lbim.ModelMessage{})
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

func (a *MessageImpl) GetByIdList(ctx context.Context, idList []uint64) ([]*lbim.ModelMessage, error) {
	var valList []*lbim.ModelMessage
	db := webtool.NewCondBuilder(a.mustGetConn(ctx))
	err := db.In(lbim.FieldId_, idList).Find(ctx, &valList)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return valList, nil
}

func (a *MessageImpl) UpdateById(ctx context.Context, id uint64, updateMap map[string]interface{}) error {
	db := webtool.NewCondBuilder(a.mustGetConn(ctx))
	_, err := db.Eq(lbim.FieldId_, id).Update(ctx, updateMap)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

func (a *MessageImpl) UpdateOrCreate(ctx context.Context, candMap, attrMap map[string]interface{}, out *lbim.ModelMessage) error {
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

func (a *MessageImpl) BatchCreate(ctx context.Context, valList []*lbim.ModelMessage) error {
	res := a.mustGetConn(ctx).CreateInBatches(valList, len(valList))
	log.Infof("batch create rows_affected %d", res.RowsAffected)
	if res.Error != nil {
		log.Errorf("err:%v", res.Error)
		return res.Error
	}
	return nil
}

func (a *MessageImpl) Create(ctx context.Context, val *lbim.ModelMessage) error {
	db := webtool.NewCondBuilder(a.mustGetConn(ctx))
	_, err := db.Create(ctx, val)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

func (a *MessageImpl) DeleteById(ctx context.Context, id uint64) error {
	db := webtool.NewCondBuilder(a.mustGetConn(ctx))
	_, err := db.Eq(lbim.FieldId_, id).Delete(ctx, &lbim.ModelMessage{})
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

func (a *MessageImpl) GetById(ctx context.Context, id uint64) (*lbim.ModelMessage, error) {
	var val lbim.ModelMessage
	db := webtool.NewCondBuilder(a.mustGetConn(ctx))
	err := db.Eq(lbim.FieldId_, id).First(ctx, &val)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return &val, nil
}

func (a *MessageImpl) FindPage(ctx context.Context, listOption *lbconst.ListOption) ([]*lbim.ModelMessage, *lbconst.Page, error) {
	var list []*lbim.ModelMessage
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

func NewMessageImpl(ctx context.Context, dsn string) (dao.MessageDao, error) {
	var d = &MessageImpl{
		mysqlConn{
			dsn: dsn,
		},
	}
	if !migratedMessage.Load() {
		err := d.mustGetConn(ctx).Set(autoMigrateOptKey, autoMigrateOptValue).AutoMigrate(&lbim.ModelMessage{})
		if err != nil {
			log.Errorf("err:%v", err)
			return nil, err
		}
	}
	return d, nil
}
