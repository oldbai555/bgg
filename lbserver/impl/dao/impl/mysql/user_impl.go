package mysql

import (
	"context"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/oldbai555/bgg/client/lbconst"
	"github.com/oldbai555/bgg/client/lbuser"
	"github.com/oldbai555/bgg/lbserver/impl/dao"
	webtool "github.com/oldbai555/bgg/pkg/webtoolv2"
	"github.com/oldbai555/gorm"
	"github.com/oldbai555/lbtool/log"
)

var _ dao.UserDao = (*UserImpl)(nil)

type UserImpl struct {
	mysqlConn
}

func (a *UserImpl) mustGetConn(ctx context.Context) *gorm.DB {
	return a.mysqlConn.mustGetConn(ctx).Model(&lbuser.ModelUser{})
}

// id:需要排除的 ID
func (a *UserImpl) CheckUserNameExit(ctx context.Context, id uint64, username string) (bool, error) {
	db := webtool.NewCondBuilder(a.mustGetConn(ctx))

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
	db := webtool.NewCondBuilder(a.mustGetConn(ctx))
	err := db.Eq(lbuser.FieldRole_, uint32(lbuser.ModelUser_RoleAdmin)).First(ctx, &val)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return &val, nil
}

func (a *UserImpl) GetByUserName(ctx context.Context, username string) (*lbuser.ModelUser, error) {
	var val lbuser.ModelUser
	db := webtool.NewCondBuilder(a.mustGetConn(ctx))
	err := db.Eq(lbuser.FieldUsername_, username).First(ctx, &val)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return &val, nil
}

func (a *UserImpl) UpdateByIdList(ctx context.Context, idList []uint64, updateMap map[string]interface{}) error {
	db := webtool.NewCondBuilder(a.mustGetConn(ctx))
	_, err := db.In(lbuser.FieldId_, idList).Update(ctx, updateMap)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

func (a *UserImpl) DeleteByIdList(ctx context.Context, idList []uint64) error {
	db := webtool.NewCondBuilder(a.mustGetConn(ctx))
	_, err := db.In(lbuser.FieldId_, idList).Delete(ctx, &lbuser.ModelUser{})
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

func (a *UserImpl) GetByIdList(ctx context.Context, idList []uint64) ([]*lbuser.ModelUser, error) {
	var valList []*lbuser.ModelUser
	db := webtool.NewCondBuilder(a.mustGetConn(ctx))
	err := db.In(lbuser.FieldId_, idList).Find(ctx, &valList)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return valList, nil
}

func (a *UserImpl) UpdateById(ctx context.Context, id uint64, updateMap map[string]interface{}) error {
	db := webtool.NewCondBuilder(a.mustGetConn(ctx))
	_, err := db.Eq(lbuser.FieldId_, id).Update(ctx, updateMap)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

func (a *UserImpl) UpdateOrCreate(ctx context.Context, candMap, attrMap map[string]interface{}, out *lbuser.ModelUser) error {
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

func (a *UserImpl) BatchCreate(ctx context.Context, valList []*lbuser.ModelUser) error {
	res := a.mustGetConn(ctx).CreateInBatches(valList, len(valList))
	log.Infof("batch create rows_affected %d", res.RowsAffected)
	if res.Error != nil {
		log.Errorf("err:%v", res.Error)
		return res.Error
	}
	return nil
}

func (a *UserImpl) Create(ctx context.Context, val *lbuser.ModelUser) error {
	db := webtool.NewCondBuilder(a.mustGetConn(ctx))
	_, err := db.Create(ctx, val)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

func (a *UserImpl) DeleteById(ctx context.Context, id uint64) error {
	db := webtool.NewCondBuilder(a.mustGetConn(ctx))
	_, err := db.Eq(lbuser.FieldId_, id).Delete(ctx, &lbuser.ModelUser{})
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

func (a *UserImpl) GetById(ctx context.Context, id uint64) (*lbuser.ModelUser, error) {
	var val lbuser.ModelUser
	db := webtool.NewCondBuilder(a.mustGetConn(ctx))
	err := db.Eq(lbuser.FieldId_, id).First(ctx, &val)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return &val, nil
}

func (a *UserImpl) FindPage(ctx context.Context, listOption *lbconst.ListOption) ([]*lbuser.ModelUser, *lbconst.Page, error) {
	var list []*lbuser.ModelUser
	db := webtool.NewList(a.mustGetConn(ctx), listOption)
	err := lbconst.NewListOptionProcessor(listOption).
		AddString(lbuser.GetUserListReq_ListOptionLikeUsername, func(val string) error {
			db.Like(lbuser.FieldUsername_, fmt.Sprintf("%%%s%%", val))
			return nil
		}).
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

func NewUserImpl(ctx context.Context, dsn string) (dao.UserDao, error) {
	var d = &UserImpl{
		mysqlConn{
			dsn: dsn,
		},
	}
	err := d.mysqlConn.mustGetConn(ctx).Set(autoMigrateOptKey, autoMigrateOptValue).AutoMigrate(&lbuser.ModelUser{})
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return d, nil
}