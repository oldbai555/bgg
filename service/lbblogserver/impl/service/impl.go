package service

import (
	"context"
	"fmt"
	"github.com/oldbai555/bgg/service/lb"
	"github.com/oldbai555/bgg/service/lbblog"
	"github.com/oldbai555/bgg/service/lbblogserver/impl/mysql"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/utils"
)

var OnceSvrImpl = &LbblogServer{}

type LbblogServer struct {
	lbblog.UnimplementedLbblogServer
}

func (a *LbblogServer) AddArticle(ctx context.Context, req *lbblog.AddArticleReq) (*lbblog.AddArticleRsp, error) {
	var rsp lbblog.AddArticleRsp
	var err error

	_, err = mysql.Article.NewScope(ctx).Create(req.Data)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.Data = req.Data

	return &rsp, err
}
func (a *LbblogServer) DelArticleList(ctx context.Context, req *lbblog.DelArticleListReq) (*lbblog.DelArticleListRsp, error) {
	var rsp lbblog.DelArticleListRsp
	var err error

	listRsp, err := a.GetArticleList(ctx, &lbblog.GetArticleListReq{
		ListOption: req.ListOption.
			SetSkipTotal().
			AddOpt(lb.DefaultListOption_DefaultListOptionSelect, lbblog.FieldId_),
	})
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	if len(listRsp.List) == 0 {
		log.Infof("list is empty")
		return &rsp, nil
	}

	idList := utils.PluckUint64List(listRsp.List, lbblog.FieldId)
	_, err = mysql.Article.NewScope(ctx).In(lbblog.FieldId, idList).Delete(&lbblog.ModelArticle{})
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}
func (a *LbblogServer) UpdateArticle(ctx context.Context, req *lbblog.UpdateArticleReq) (*lbblog.UpdateArticleRsp, error) {
	var rsp lbblog.UpdateArticleRsp
	var err error

	var data lbblog.ModelArticle
	err = mysql.Article.NewScope(ctx).Where(lbblog.FieldId_, req.Data.Id).First(&data)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	_, err = mysql.Article.NewScope(ctx).Where(lbblog.FieldId_, data.Id).Update(utils.OrmStruct2Map(req.Data))
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}
func (a *LbblogServer) GetArticle(ctx context.Context, req *lbblog.GetArticleReq) (*lbblog.GetArticleRsp, error) {
	var rsp lbblog.GetArticleRsp
	var err error

	var data lbblog.ModelArticle
	err = mysql.Article.NewScope(ctx).Where(lbblog.FieldId_, req.Id).First(&data)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.Data = &data

	return &rsp, err
}
func (a *LbblogServer) GetArticleList(ctx context.Context, req *lbblog.GetArticleListReq) (*lbblog.GetArticleListRsp, error) {
	var rsp lbblog.GetArticleListRsp
	var err error

	db := mysql.Article.NewList(ctx, req.ListOption)
	err = lb.ProcessDefaultOptions(req.ListOption, db)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	err = lb.NewOptionsProcessor(req.ListOption).
		AddUint64(lbblog.GetArticleListReq_ListOptionCategoryId, func(val uint64) error {
			db.Eq(lbblog.FieldCategoryId_, val)
			return nil
		}).
		AddString(lbblog.GetArticleListReq_ListOptionLikeTitle, func(val string) error {
			db.Like(lbblog.FieldTitle_, fmt.Sprintf("%%%v%%", val))
			return nil
		}).
		Process()
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	rsp.Paginate, err = db.FindPaginate(&rsp.List)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}
func (a *LbblogServer) AddCategory(ctx context.Context, req *lbblog.AddCategoryReq) (*lbblog.AddCategoryRsp, error) {
	var rsp lbblog.AddCategoryRsp
	var err error

	_, err = mysql.Category.NewScope(ctx).Create(req.Data)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.Data = req.Data

	return &rsp, err
}
func (a *LbblogServer) DelCategoryList(ctx context.Context, req *lbblog.DelCategoryListReq) (*lbblog.DelCategoryListRsp, error) {
	var rsp lbblog.DelCategoryListRsp
	var err error

	listRsp, err := a.GetCategoryList(ctx, &lbblog.GetCategoryListReq{
		ListOption: req.ListOption.
			SetSkipTotal().
			AddOpt(lb.DefaultListOption_DefaultListOptionSelect, lbblog.FieldId_),
	})
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	if len(listRsp.List) == 0 {
		log.Infof("list is empty")
		return &rsp, nil
	}

	idList := utils.PluckUint64List(listRsp.List, lbblog.FieldId)
	_, err = mysql.Category.NewScope(ctx).In(lbblog.FieldId, idList).Delete(&lbblog.ModelCategory{})
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}
func (a *LbblogServer) UpdateCategory(ctx context.Context, req *lbblog.UpdateCategoryReq) (*lbblog.UpdateCategoryRsp, error) {
	var rsp lbblog.UpdateCategoryRsp
	var err error

	var data lbblog.ModelCategory
	err = mysql.Category.NewScope(ctx).Where(lbblog.FieldId_, req.Data.Id).First(&data)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	_, err = mysql.Category.NewScope(ctx).Where(lbblog.FieldId_, data.Id).Update(utils.OrmStruct2Map(req.Data))
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}
func (a *LbblogServer) GetCategory(ctx context.Context, req *lbblog.GetCategoryReq) (*lbblog.GetCategoryRsp, error) {
	var rsp lbblog.GetCategoryRsp
	var err error

	var data lbblog.ModelCategory
	err = mysql.Category.NewScope(ctx).Where(lbblog.FieldId_, req.Id).First(&data)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.Data = &data

	return &rsp, err
}
func (a *LbblogServer) GetCategoryList(ctx context.Context, req *lbblog.GetCategoryListReq) (*lbblog.GetCategoryListRsp, error) {
	var rsp lbblog.GetCategoryListRsp
	var err error

	db := mysql.Category.NewList(ctx, req.ListOption)
	err = lb.ProcessDefaultOptions(req.ListOption, db)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	err = lb.NewOptionsProcessor(req.ListOption).
		AddString(lbblog.GetCategoryListReq_ListOptionLikeName, func(val string) error {
			db.Like(lbblog.FieldName_, fmt.Sprintf("%%%v%%", val))
			return nil
		}).
		Process()

	rsp.Paginate, err = db.FindPaginate(&rsp.List)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}
func (a *LbblogServer) AddComment(ctx context.Context, req *lbblog.AddCommentReq) (*lbblog.AddCommentRsp, error) {
	var rsp lbblog.AddCommentRsp
	var err error

	_, err = mysql.Comment.NewScope(ctx).Create(req.Data)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.Data = req.Data

	return &rsp, err
}
func (a *LbblogServer) DelCommentList(ctx context.Context, req *lbblog.DelCommentListReq) (*lbblog.DelCommentListRsp, error) {
	var rsp lbblog.DelCommentListRsp
	var err error

	listRsp, err := a.GetCommentList(ctx, &lbblog.GetCommentListReq{
		ListOption: req.ListOption.
			SetSkipTotal().
			AddOpt(lb.DefaultListOption_DefaultListOptionSelect, lbblog.FieldId_),
	})
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	if len(listRsp.List) == 0 {
		log.Infof("list is empty")
		return &rsp, nil
	}

	idList := utils.PluckUint64List(listRsp.List, lbblog.FieldId)
	_, err = mysql.Comment.NewScope(ctx).In(lbblog.FieldId, idList).Delete(&lbblog.ModelComment{})
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}
func (a *LbblogServer) UpdateComment(ctx context.Context, req *lbblog.UpdateCommentReq) (*lbblog.UpdateCommentRsp, error) {
	var rsp lbblog.UpdateCommentRsp
	var err error

	var data lbblog.ModelComment
	err = mysql.Comment.NewScope(ctx).Where(lbblog.FieldId_, req.Data.Id).First(&data)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	_, err = mysql.Comment.NewScope(ctx).Where(lbblog.FieldId_, data.Id).Update(utils.OrmStruct2Map(req.Data))
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}
func (a *LbblogServer) GetComment(ctx context.Context, req *lbblog.GetCommentReq) (*lbblog.GetCommentRsp, error) {
	var rsp lbblog.GetCommentRsp
	var err error

	var data lbblog.ModelComment
	err = mysql.Comment.NewScope(ctx).Where(lbblog.FieldId_, req.Id).First(&data)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.Data = &data

	return &rsp, err
}
func (a *LbblogServer) GetCommentList(ctx context.Context, req *lbblog.GetCommentListReq) (*lbblog.GetCommentListRsp, error) {
	var rsp lbblog.GetCommentListRsp
	var err error

	db := mysql.Comment.NewList(ctx, req.ListOption)
	err = lb.ProcessDefaultOptions(req.ListOption, db)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	err = lb.NewOptionsProcessor(req.ListOption).
		AddString(lbblog.GetCommentListReq_ListOptionLikeContent, func(val string) error {
			db.Like(lbblog.FieldContent_, fmt.Sprintf("%%%v%%", val))
			return nil
		}).
		AddUint32(lbblog.GetCommentListReq_ListOptionStatus, func(val uint32) error {
			db.Eq(lbblog.FieldStatus_, val)
			return nil
		}).
		AddString(lbblog.GetCommentListReq_ListOptionLikeUserEmail, func(val string) error {
			db.Like(lbblog.FieldUserEmail_, fmt.Sprintf("%%%v%%", val))
			return nil
		}).
		AddUint64List(lbblog.GetCommentListReq_ListOptionArticleIdList, func(valList []uint64) error {
			if len(valList) == 1 {
				db.Eq(lbblog.FieldArticleId_, valList[0])
			} else {
				db.In(lbblog.FieldArticleId_, valList)
			}
			return nil
		}).
		Process()

	rsp.Paginate, err = db.FindPaginate(&rsp.List)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}
