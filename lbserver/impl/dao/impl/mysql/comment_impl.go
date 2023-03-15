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

var _ dao.CommentDao = (*CommentImpl)(nil)

type CommentImpl struct {
	mysqlConn
}

func (a *CommentImpl) UpdateByIdList(ctx context.Context, idList []uint64, updateMap map[string]interface{}) error {
	db := webtool.NewCondBuilder(a.mustGetConn(ctx))
	_, err := db.In(lbblog.FieldId_, idList).Update(ctx, updateMap)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

func (a *CommentImpl) DeleteByIdList(ctx context.Context, idList []uint64) error {
	db := webtool.NewCondBuilder(a.mustGetConn(ctx))
	_, err := db.In(lbblog.FieldId_, idList).Delete(ctx, &lbblog.ModelComment{})
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

func (a *CommentImpl) GetByIdList(ctx context.Context, idList []uint64) ([]*lbblog.ModelComment, error) {
	var valList []*lbblog.ModelComment
	db := webtool.NewCondBuilder(a.mustGetConn(ctx))
	err := db.In(lbblog.FieldId_, idList).Find(ctx, &valList)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return valList, nil
}

func (a *CommentImpl) UpdateById(ctx context.Context, id uint64, updateMap map[string]interface{}) error {
	db := webtool.NewCondBuilder(a.mustGetConn(ctx))
	_, err := db.Eq(lbblog.FieldId_, id).Update(ctx, updateMap)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

func (a *CommentImpl) UpdateOrCreate(ctx context.Context, candMap, attrMap map[string]interface{}, out *lbblog.ModelComment) error {
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

func (a *CommentImpl) BatchCreate(ctx context.Context, valList []*lbblog.ModelComment) error {
	res := a.mustGetConn(ctx).CreateInBatches(valList, len(valList))
	log.Infof("batch create rows_affected %d", res.RowsAffected)
	if res.Error != nil {
		log.Errorf("err:%v", res.Error)
		return res.Error
	}
	return nil
}

func (a *CommentImpl) Create(ctx context.Context, val *lbblog.ModelComment) error {
	db := webtool.NewCondBuilder(a.mustGetConn(ctx))
	_, err := db.Create(ctx, val)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

func (a *CommentImpl) DeleteById(ctx context.Context, id uint64) error {
	db := webtool.NewCondBuilder(a.mustGetConn(ctx))
	_, err := db.Eq(lbblog.FieldId_, id).Delete(ctx, &lbblog.ModelComment{})
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

func (a *CommentImpl) GetById(ctx context.Context, id uint64) (*lbblog.ModelComment, error) {
	var val lbblog.ModelComment
	db := webtool.NewCondBuilder(a.mustGetConn(ctx))
	err := db.Eq(lbblog.FieldId_, id).First(ctx, &val)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return &val, nil
}

func (a *CommentImpl) FindPage(ctx context.Context, listOption *lbconst.ListOption) ([]*lbblog.ModelComment, *lbconst.Page, error) {
	var list []*lbblog.ModelComment
	db := webtool.NewList(a.mustGetConn(ctx), listOption)
	err := lbconst.NewListOptionProcessor(listOption).
		AddUint64(lbblog.GetCommentListReq_ListOptionArticleId, func(val uint64) error {
			db.Eq(lbblog.FieldArticleId_, val)
			return nil
		}).
		Process()
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, nil, err
	}

	page, err := db.OrderByDesc(lbblog.FieldCreatedAt_).FindPage(ctx, &list)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, nil, err
	}
	return list, page, nil
}

func NewCommentImpl(ctx context.Context, dsn string) (dao.CommentDao, error) {
	var d = &CommentImpl{
		mysqlConn{
			dsn: dsn,
		},
	}
	err := d.mustGetConn(ctx).Set(autoMigrateOptKey, autoMigrateOptValue).AutoMigrate(&lbblog.ModelComment{})
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return d, nil
}
