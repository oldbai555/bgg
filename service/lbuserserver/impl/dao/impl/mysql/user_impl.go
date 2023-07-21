package mysql

import (
	"context"
	"github.com/mitchellh/mapstructure"
	"github.com/oldbai555/bgg/pkg/webtool"
	"github.com/oldbai555/bgg/service/lb"
	"github.com/oldbai555/bgg/service/lbuser"
	"github.com/oldbai555/bgg/service/lbuserserver/impl/dao"
	"github.com/oldbai555/gorm"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/lberr"
)

var _ dao.IUserDao = (*UserImpl)(nil)

type UserImpl struct{}

// CheckUserNameExit id:需要排除的 ID
func (a *UserImpl) CheckUserNameExit(ctx context.Context, id uint64, username string) (bool, error) {
	db := a.GetOrmCondBuilder(ctx)

	if id > 0 {
		db.NotEq(lbuser.FieldId_, id)
	}

	err := db.Eq(lbuser.FieldUsername_, username).First(ctx, &lbuser.ModelUser{})
	if err != nil {
		log.Errorf("err:%v", err)
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (a *UserImpl) GetByAdmin(ctx context.Context) (*lbuser.ModelUser, error) {
	var val lbuser.ModelUser
	db := a.GetOrmCondBuilder(ctx)
	err := db.Eq(lbuser.FieldRole_, uint32(lbuser.ModelUser_RoleAdmin)).First(ctx, &val)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return &val, nil
}

func (a *UserImpl) GetByUserName(ctx context.Context, username string) (*lbuser.ModelUser, error) {
	var val lbuser.ModelUser
	db := a.GetOrmCondBuilder(ctx)
	err := db.Eq(lbuser.FieldUsername_, username).First(ctx, &val)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return &val, nil
}

func (a *UserImpl) GetOrmEngine(ctx context.Context) *gorm.DB {
	return MasterOrm.WithContext(ctx).Model(&lbuser.ModelUser{})
}

func (a *UserImpl) Create(ctx context.Context, val *lbuser.ModelUser) (*webtool.Result, error) {
	db := a.GetOrmCondBuilder(ctx)
	rows, err := db.Create(ctx, val)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return rows, nil
}

func (a *UserImpl) FirstOrCreate(ctx context.Context, cand map[string]interface{}, val *lbuser.ModelUser) (*webtool.Result, error) {
	result := a.GetOrmEngine(ctx).FirstOrCreate(val, cand)
	return webtool.NewResult(result.RowsAffected, result.RowsAffected > 0), result.Error
}

func (a *UserImpl) BatchCreate(ctx context.Context, valList []*lbuser.ModelUser) (*webtool.Result, error) {
	res := a.GetOrmEngine(ctx).CreateInBatches(valList, len(valList))
	log.Infof("batch create rows_affected %d", res.RowsAffected)
	if res.Error != nil {
		log.Errorf("err:%v", res.Error)
		return nil, res.Error
	}
	return webtool.NewResult(res.RowsAffected, true), nil
}

func (a *UserImpl) UpdateOrCreate(ctx context.Context, candMap, attrMap map[string]interface{}, out *lbuser.ModelUser) (*webtool.Result, error) {
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

func (a *UserImpl) UpdateById(ctx context.Context, id uint64, updateMap map[string]interface{}) (*webtool.Result, error) {
	db := a.GetOrmCondBuilder(ctx)
	rows, err := db.Eq(lbuser.FieldId_, id).Update(ctx, updateMap)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return rows, nil
}

func (a *UserImpl) UpdateByIdList(ctx context.Context, idList []uint64, updateMap map[string]interface{}) (*webtool.Result, error) {
	db := a.GetOrmCondBuilder(ctx)
	rows, err := db.In(lbuser.FieldId_, idList).Update(ctx, updateMap)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return rows, nil
}

func (a *UserImpl) DeleteById(ctx context.Context, id uint64) (*webtool.Result, error) {
	db := a.GetOrmCondBuilder(ctx)
	rows, err := db.Eq(lbuser.FieldId_, id).Delete(ctx, &lbuser.ModelUser{})
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return rows, nil
}

func (a *UserImpl) DeleteByIdList(ctx context.Context, idList []uint64) (*webtool.Result, error) {
	db := a.GetOrmCondBuilder(ctx)
	rows, err := db.In(lbuser.FieldId_, idList).Delete(ctx, &lbuser.ModelUser{})
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return rows, nil
}

func (a *UserImpl) GetById(ctx context.Context, id uint64) (*lbuser.ModelUser, error) {
	var val lbuser.ModelUser
	db := a.GetOrmCondBuilder(ctx)
	err := db.Eq(lbuser.FieldId_, id).First(ctx, &val)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return &val, nil
}

func (a *UserImpl) GetByIdList(ctx context.Context, idList []uint64) ([]*lbuser.ModelUser, error) {
	var valList []*lbuser.ModelUser
	db := a.GetOrmCondBuilder(ctx)
	err := db.In(lbuser.FieldId_, idList).Find(ctx, &valList)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return valList, nil
}

func (a *UserImpl) FindPaginate(ctx context.Context, options *lb.Options) ([]*lbuser.ModelUser, *lb.Paginate, error) {
	var list []*lbuser.ModelUser
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

func (a *UserImpl) Increment(ctx context.Context, field string, num uint32, candMap map[string]interface{}) error {
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

func (a *UserImpl) Decrement(ctx context.Context, field string, num uint32, candMap map[string]interface{}) error {
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

func (a *UserImpl) IsNotFoundErr(err error) bool {
	return err == gorm.ErrRecordNotFound
}

func (a *UserImpl) GetList(ctx context.Context, candMap map[string]interface{}, opts ...*webtool.Opt) ([]*lbuser.ModelUser, error) {
	var list []*lbuser.ModelUser
	db := a.GetOrmCondBuilder(ctx)
	webtool.ProcessOpts(db, opts...)
	err := db.Where(candMap).Find(ctx, &list)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return list, nil
}

func (a *UserImpl) GetOne(ctx context.Context, candMap map[string]interface{}, opts ...*webtool.Opt) (*lbuser.ModelUser, error) {
	var val lbuser.ModelUser
	db := a.GetOrmCondBuilder(ctx)
	webtool.ProcessOpts(db, opts...)
	err := db.Where(candMap).First(ctx, &val)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return &val, nil
}

func (a *UserImpl) GetOrmCondBuilder(ctx context.Context) *webtool.OrmCondBuilder {
	return webtool.NewCondBuilder(a.GetOrmEngine(ctx))
}

func NewUserImpl() (dao.IUserDao, error) {
	return &UserImpl{}, nil
}
