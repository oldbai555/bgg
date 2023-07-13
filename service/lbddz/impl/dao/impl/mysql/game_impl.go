package mysql

import (
	"context"
	"github.com/mitchellh/mapstructure"
	"github.com/oldbai555/bgg/client/lb"
	"github.com/oldbai555/bgg/client/lbddz"
	"github.com/oldbai555/bgg/pkg/webtool"
	"github.com/oldbai555/bgg/service/lbddz/impl/dao"
	"github.com/oldbai555/gorm"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/lberr"
)

var _ dao.IGameDao = (*GameImpl)(nil)

type GameImpl struct{}

func (a *GameImpl) GetOrmEngine(ctx context.Context) *gorm.DB {
	return MasterOrm.WithContext(ctx).Model(&lbddz.ModelGame{})
}

func (a *GameImpl) Create(ctx context.Context, val *lbddz.ModelGame) (*webtool.Result, error) {
	db := a.GetOrmCondBuilder(ctx)
	rows, err := db.Create(ctx, val)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return rows, nil
}

func (a *GameImpl) FirstOrCreate(ctx context.Context, cand map[string]interface{}, val *lbddz.ModelGame) (*webtool.Result, error) {
	result := a.GetOrmEngine(ctx).FirstOrCreate(val, cand)
	return webtool.NewResult(result.RowsAffected, result.RowsAffected > 0), result.Error
}

func (a *GameImpl) BatchCreate(ctx context.Context, valList []*lbddz.ModelGame) (*webtool.Result, error) {
	res := a.GetOrmEngine(ctx).CreateInBatches(valList, len(valList))
	log.Infof("batch create rows_affected %d", res.RowsAffected)
	if res.Error != nil {
		log.Errorf("err:%v", res.Error)
		return nil, res.Error
	}
	return webtool.NewResult(res.RowsAffected, true), nil
}

func (a *GameImpl) UpdateOrCreate(ctx context.Context, candMap, attrMap map[string]interface{}, out *lbddz.ModelGame) (*webtool.Result, error) {
	selectDb := a.GetOrmCondBuilder(ctx)
	err := selectDb.AndMap(candMap).First(ctx, out)
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Errorf("err:%v", err)
		return nil, err
	}

	optDb := a.GetOrmCondBuilder(ctx)
	if err == gorm.ErrRecordNotFound {
		for k, v := range candMap {
			attrMap[k] = v
		}
		err = mapstructure.Decode(attrMap, out)
		if err != nil {
			log.Errorf("err is : %v", err)
			return nil, err
		}
		rows, err := optDb.Create(ctx, out)
		if err != nil {
			log.Errorf("err is %v", err)
			return nil, err
		}
		return rows, nil
	}

	rows, err := optDb.AndMap(candMap).Update(ctx, attrMap)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return rows, nil
}

func (a *GameImpl) UpdateById(ctx context.Context, id uint64, updateMap map[string]interface{}) (*webtool.Result, error) {
	db := a.GetOrmCondBuilder(ctx)
	rows, err := db.Eq(lbddz.FieldId_, id).Update(ctx, updateMap)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return rows, nil
}

func (a *GameImpl) UpdateByIdList(ctx context.Context, idList []uint64, updateMap map[string]interface{}) (*webtool.Result, error) {
	db := a.GetOrmCondBuilder(ctx)
	rows, err := db.In(lbddz.FieldId_, idList).Update(ctx, updateMap)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return rows, nil
}

func (a *GameImpl) DeleteById(ctx context.Context, id uint64) (*webtool.Result, error) {
	db := a.GetOrmCondBuilder(ctx)
	rows, err := db.Eq(lbddz.FieldId_, id).Delete(ctx, &lbddz.ModelGame{})
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return rows, nil
}

func (a *GameImpl) DeleteByIdList(ctx context.Context, idList []uint64) (*webtool.Result, error) {
	db := a.GetOrmCondBuilder(ctx)
	rows, err := db.In(lbddz.FieldId_, idList).Delete(ctx, &lbddz.ModelGame{})
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return rows, nil
}

func (a *GameImpl) GetById(ctx context.Context, id uint64) (*lbddz.ModelGame, error) {
	var val lbddz.ModelGame
	db := a.GetOrmCondBuilder(ctx)
	err := db.Eq(lbddz.FieldId_, id).First(ctx, &val)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return &val, nil
}

func (a *GameImpl) GetByIdList(ctx context.Context, idList []uint64) ([]*lbddz.ModelGame, error) {
	var valList []*lbddz.ModelGame
	db := a.GetOrmCondBuilder(ctx)
	err := db.In(lbddz.FieldId_, idList).Find(ctx, &valList)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return valList, nil
}

func (a *GameImpl) FindPaginate(ctx context.Context, options *lb.Options) ([]*lbddz.ModelGame, *lb.Paginate, error) {
	var list []*lbddz.ModelGame
	db := webtool.NewList(a.GetOrmEngine(ctx), options)
	err := webtool.ProcessDefaultOptions(options, db)
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, nil, err
	}
	err = lb.NewOptionsProcessor(options).
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

func (a *GameImpl) Increment(ctx context.Context, field string, num uint32, candMap map[string]interface{}) error {
	if len(candMap) == 0 {
		return lberr.NewInvalidArg("candMap must be not nil")
	}
	db := a.GetOrmCondBuilder(ctx)
	_, err := db.AndMap(candMap).Update(ctx, map[string]interface{}{
		field: gorm.Expr("? + ?", field, num),
	})
	if err != nil {
		log.Errorf("err is %v", err)
		return err
	}
	return nil
}

func (a *GameImpl) Decrement(ctx context.Context, field string, num uint32, candMap map[string]interface{}) error {
	if len(candMap) == 0 {
		return lberr.NewInvalidArg("candMap must be not nil")
	}
	db := a.GetOrmCondBuilder(ctx)
	_, err := db.AndMap(candMap).Update(ctx, map[string]interface{}{
		field: gorm.Expr("? + ?", field, num),
	})
	if err != nil {
		log.Errorf("err is %v", err)
		return err
	}
	return nil
}

func (a *GameImpl) IsNotFoundErr(err error) bool {
	return err == gorm.ErrRecordNotFound
}

func (a *GameImpl) GetList(ctx context.Context, candMap map[string]interface{}, opts ...*webtool.Opt) ([]*lbddz.ModelGame, error) {
	var list []*lbddz.ModelGame
	db := a.GetOrmCondBuilder(ctx)
	webtool.ProcessOpts(db, opts...)
	err := db.Where(candMap).Find(ctx, &list)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return list, nil
}

func (a *GameImpl) GetOne(ctx context.Context, candMap map[string]interface{}, opts ...*webtool.Opt) (*lbddz.ModelGame, error) {
	var val lbddz.ModelGame
	db := a.GetOrmCondBuilder(ctx)
	webtool.ProcessOpts(db, opts...)
	err := db.Where(candMap).First(ctx, &val)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return &val, nil
}

func (a *GameImpl) GetOrmCondBuilder(ctx context.Context) *webtool.OrmCondBuilder {
	return webtool.NewCondBuilder(a.GetOrmEngine(ctx))
}

func NewGameImpl() (dao.IGameDao, error) {
	return &GameImpl{}, nil
}
