package mysql

import (
	"context"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/oldbai555/bgg/client/lbblog"
	"github.com/oldbai555/bgg/client/lbconst"
	"github.com/oldbai555/bgg/lbserver/impl/dao"
	webtool "github.com/oldbai555/bgg/pkg/webtoolv2"
	"github.com/oldbai555/gorm"
	"github.com/oldbai555/lbtool/log"
)

var _ dao.ArticleDao = (*ArticleImpl)(nil)

type ArticleImpl struct {
	mysqlConn
}

func (a *ArticleImpl) UpdateByIdList(ctx context.Context, idList []uint64, updateMap map[string]interface{}) error {
	db := webtool.NewCondBuilder(a.mustGetConn(ctx))
	_, err := db.In(lbblog.FieldId_, idList).Update(ctx, updateMap)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

func (a *ArticleImpl) DeleteByIdList(ctx context.Context, idList []uint64) error {
	db := webtool.NewCondBuilder(a.mustGetConn(ctx))
	_, err := db.In(lbblog.FieldId_, idList).Delete(ctx, &lbblog.ModelArticle{})
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

func (a *ArticleImpl) GetByIdList(ctx context.Context, idList []uint64) ([]*lbblog.ModelArticle, error) {
	var valList []*lbblog.ModelArticle
	db := webtool.NewCondBuilder(a.mustGetConn(ctx))
	err := db.In(lbblog.FieldId_, idList).Find(ctx, &valList)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return valList, nil
}

func (a *ArticleImpl) UpdateById(ctx context.Context, id uint64, updateMap map[string]interface{}) error {
	db := webtool.NewCondBuilder(a.mustGetConn(ctx))
	_, err := db.Eq(lbblog.FieldId_, id).Update(ctx, updateMap)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

func (a *ArticleImpl) UpdateOrCreate(ctx context.Context, candMap, attrMap map[string]interface{}, out *lbblog.ModelArticle) error {
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

func (a *ArticleImpl) BatchCreate(ctx context.Context, valList []*lbblog.ModelArticle) error {
	res := a.mustGetConn(ctx).CreateInBatches(valList, len(valList))
	log.Infof("batch create rows_affected %d", res.RowsAffected)
	if res.Error != nil {
		log.Errorf("err:%v", res.Error)
		return res.Error
	}
	return nil
}

func (a *ArticleImpl) Create(ctx context.Context, val *lbblog.ModelArticle) error {
	db := webtool.NewCondBuilder(a.mustGetConn(ctx))
	_, err := db.Create(ctx, val)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

func (a *ArticleImpl) DeleteById(ctx context.Context, id uint64) error {
	db := webtool.NewCondBuilder(a.mustGetConn(ctx))
	_, err := db.Eq(lbblog.FieldId_, id).Delete(ctx, &lbblog.ModelArticle{})
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

func (a *ArticleImpl) GetById(ctx context.Context, id uint64) (*lbblog.ModelArticle, error) {
	var val lbblog.ModelArticle
	db := webtool.NewCondBuilder(a.mustGetConn(ctx))
	err := db.Eq(lbblog.FieldId_, id).First(ctx, &val)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return &val, nil
}

func (a *ArticleImpl) FindPage(ctx context.Context, listOption *lbconst.ListOption) ([]*lbblog.ModelArticle, *lbconst.Page, error) {
	var list []*lbblog.ModelArticle
	db := webtool.NewList(a.mustGetConn(ctx), listOption)
	err := lbconst.NewListOptionProcessor(listOption).
		AddUint64(lbblog.GetArticleListReq_ListOptionCategoryId, func(val uint64) error {
			db.Eq(lbblog.FieldCategoryId_, val)
			return nil
		}).
		AddString(lbblog.GetArticleListReq_ListOptionLikeTitle, func(val string) error {
			db.Like(lbblog.FieldTitle_, fmt.Sprintf("%%%s%%", val))
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

func NewArticleImpl(ctx context.Context, dsn string) (dao.ArticleDao, error) {
	var d = &ArticleImpl{
		mysqlConn{
			dsn: dsn,
		},
	}
	err := d.mustGetConn(ctx).Set(autoMigrateOptKey, autoMigrateOptValue).AutoMigrate(&lbblog.ModelArticle{})
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return d, nil
}
