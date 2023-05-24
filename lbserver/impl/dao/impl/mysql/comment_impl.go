package mysql

import (
	"context"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/oldbai555/bgg/client/lb"
	"github.com/oldbai555/bgg/client/lbblog"
	"github.com/oldbai555/bgg/lbserver/impl/dao"
	"github.com/oldbai555/bgg/pkg/webtool"
	"github.com/oldbai555/gorm"
	"github.com/oldbai555/lbtool/log"
)

var _ dao.ICommentDao = (*CommentImpl)(nil)

type CommentImpl struct {
	mysqlConn
}

func (a *CommentImpl) mustGetConn(ctx context.Context) *gorm.DB {
	return a.mysqlConn.mustGetConn(ctx).Model(&lbblog.ModelComment{})
}

func (a *CommentImpl) FirstOrCreate(ctx context.Context, val *lbblog.ModelComment, cand map[string]interface{}) error {
	return a.mysqlConn.mustGetConn(ctx).Model(&lbblog.ModelComment{}).FirstOrCreate(val, cand).Error
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

func (a *CommentImpl) FindPaginate(ctx context.Context, options *lb.Options) ([]*lbblog.ModelComment, *lb.Paginate, error) {
	var list []*lbblog.ModelComment
	db := webtool.NewList(a.mustGetConn(ctx), options)
	err := webtool.ProcessDefaultOptions(options, db)
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, nil, err
	}
	err = lb.NewOptionsProcessor(options).
		AddUint64(lbblog.GetCommentListReq_OptionArticleId, func(val uint64) error {
			db.Eq(lbblog.FieldArticleId_, val)
			return nil
		}).
		AddString(lbblog.GetCommentListReq_OptionLikeContent, func(val string) error {
			db.Like(lbblog.FieldContent_, fmt.Sprintf("%%%s%%", val))
			return nil
		}).
		AddUint32(lbblog.GetCommentListReq_OptionStatus, func(val uint32) error {
			db.Eq(lbblog.FieldStatus_, val)
			return nil
		}).
		AddString(lbblog.GetCommentListReq_OptionLikeUserEmail, func(val string) error {
			db.Like(lbblog.FieldUserEmail_, fmt.Sprintf("%%%s%%", val))
			return nil
		}).
		AddUint64List(lbblog.GetCommentListReq_OptionArticleIdList, func(valList []uint64) error {
			if len(valList) == 1 {
				db.Eq(lbblog.FieldArticleId_, valList[0])
			} else {
				db.In(lbblog.FieldArticleId_, valList)
			}
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

func NewCommentImpl(ctx context.Context, dsn string) (dao.ICommentDao, error) {
	var d = &CommentImpl{
		mysqlConn{
			dsn: dsn,
		},
	}

	err := d.mysqlConn.mustGetConn(ctx).Set(autoMigrateOptKey, autoMigrateOptValue).AutoMigrate(&lbblog.ModelComment{})
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return d, nil
}
