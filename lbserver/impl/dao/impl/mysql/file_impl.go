package mysql

import (
	"context"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/oldbai555/bgg/client/lb"
	"github.com/oldbai555/bgg/client/lbstore"
	"github.com/oldbai555/bgg/lbserver/impl/dao"
	"github.com/oldbai555/bgg/pkg/webtool"
	"github.com/oldbai555/gorm"
	"github.com/oldbai555/lbtool/log"
)

var _ dao.IFileDao = (*FileImpl)(nil)

type FileImpl struct {
	mysqlConn
}

func (a *FileImpl) mustGetConn(ctx context.Context) *gorm.DB {
	return a.mysqlConn.mustGetConn(ctx).Model(&lbstore.ModelFile{})
}

func (a *FileImpl) FirstOrCreate(ctx context.Context, val *lbstore.ModelFile, cand map[string]interface{}) error {
	return a.mysqlConn.mustGetConn(ctx).Model(&lbstore.ModelFile{}).FirstOrCreate(val, cand).Error
}

func (a *FileImpl) UpdateByIdList(ctx context.Context, idList []uint64, updateMap map[string]interface{}) error {
	db := webtool.NewCondBuilder(a.mustGetConn(ctx))
	_, err := db.In(lbstore.FieldId_, idList).Update(ctx, updateMap)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

func (a *FileImpl) DeleteByIdList(ctx context.Context, idList []uint64) error {
	db := webtool.NewCondBuilder(a.mustGetConn(ctx))
	_, err := db.In(lbstore.FieldId_, idList).Delete(ctx, &lbstore.ModelFile{})
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

func (a *FileImpl) GetByIdList(ctx context.Context, idList []uint64) ([]*lbstore.ModelFile, error) {
	var valList []*lbstore.ModelFile
	db := webtool.NewCondBuilder(a.mustGetConn(ctx))
	err := db.In(lbstore.FieldId_, idList).Find(ctx, &valList)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return valList, nil
}

func (a *FileImpl) UpdateById(ctx context.Context, id uint64, updateMap map[string]interface{}) error {
	db := webtool.NewCondBuilder(a.mustGetConn(ctx))
	_, err := db.Eq(lbstore.FieldId_, id).Update(ctx, updateMap)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

func (a *FileImpl) UpdateOrCreate(ctx context.Context, candMap, attrMap map[string]interface{}, out *lbstore.ModelFile) error {
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

func (a *FileImpl) BatchCreate(ctx context.Context, valList []*lbstore.ModelFile) error {
	res := a.mustGetConn(ctx).CreateInBatches(valList, len(valList))
	log.Infof("batch create rows_affected %d", res.RowsAffected)
	if res.Error != nil {
		log.Errorf("err:%v", res.Error)
		return res.Error
	}
	return nil
}

func (a *FileImpl) Create(ctx context.Context, val *lbstore.ModelFile) error {
	db := webtool.NewCondBuilder(a.mustGetConn(ctx))
	_, err := db.Create(ctx, val)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

func (a *FileImpl) DeleteById(ctx context.Context, id uint64) error {
	db := webtool.NewCondBuilder(a.mustGetConn(ctx))
	_, err := db.Eq(lbstore.FieldId_, id).Delete(ctx, &lbstore.ModelFile{})
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

func (a *FileImpl) GetById(ctx context.Context, id uint64) (*lbstore.ModelFile, error) {
	var val lbstore.ModelFile
	db := webtool.NewCondBuilder(a.mustGetConn(ctx))
	err := db.Eq(lbstore.FieldId_, id).First(ctx, &val)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return &val, nil
}

func (a *FileImpl) FindPaginate(ctx context.Context, options *lb.Options) ([]*lbstore.ModelFile, *lb.Paginate, error) {
	var list []*lbstore.ModelFile
	db := webtool.NewList(a.mustGetConn(ctx), options)
	err := webtool.ProcessDefaultOptions(options, db)
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, nil, err
	}
	err = lb.NewOptionsProcessor(options).
		AddString(lbstore.GetFileListReq_OptionLikeFileName, func(val string) error {
			db.Like(lbstore.FieldFileName_, fmt.Sprintf("%%%s%%", val))
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

func NewFileImpl(ctx context.Context, dsn string) (dao.IFileDao, error) {
	var d = &FileImpl{
		mysqlConn{
			dsn: dsn,
		},
	}

	err := d.mysqlConn.mustGetConn(ctx).Set(autoMigrateOptKey, autoMigrateOptValue).AutoMigrate(&lbstore.ModelFile{})
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return d, nil
}
