package impl

import (
	"context"
	"fmt"
	"github.com/oldbai555/bgg/lbblog"
	"github.com/oldbai555/bgg/lbconst"
	"github.com/oldbai555/bgg/lbuser"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/utils"
)

var lbblogServer LbblogServer

type LbblogServer struct {
	*lbblog.UnimplementedLbblogServer
}

func (a *LbblogServer) GetArticleList(ctx context.Context, req *lbblog.GetArticleListReq) (*lbblog.GetArticleListRsp, error) {
	var rsp lbblog.GetArticleListRsp

	db := ArticleOrm.NewList(req.ListOption)
	err := lbconst.NewListOptionProcessor(req.ListOption).
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
		log.Errorf("err is %v", err)
		return nil, err
	}

	rsp.Page, err = db.FindPage(ctx, &rsp.List)
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}

	if len(rsp.List) == 0 {
		return &rsp, nil
	}
	categoryIdList := utils.PluckUint64List(rsp.List, lbblog.FieldCategoryId)
	var categoryList []*lbblog.ModelCategory
	err = CategoryOrm.NewScope().In(lbblog.FieldId_, categoryIdList).Find(ctx, &categoryList)
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}
	rsp.CategoryMap = utils.Slice2MapKeyByStructField(categoryList, lbblog.FieldId).(map[uint64]*lbblog.ModelCategory)
	return &rsp, nil
}

func (a *LbblogServer) GetArticle(ctx context.Context, req *lbblog.GetArticleReq) (*lbblog.GetArticleRsp, error) {
	var rsp lbblog.GetArticleRsp

	err := ArticleOrm.NewScope().Eq(lbblog.FieldId_, req.Id).First(ctx, &rsp.Article)
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}

	return &rsp, nil
}

func (a *LbblogServer) UpdateArticle(ctx context.Context, req *lbblog.UpdateArticleReq) (*lbblog.UpdateArticleRsp, error) {
	var rsp lbblog.UpdateArticleRsp

	err := ArticleOrm.NewScope().Eq(lbblog.FieldId_, req.Article.Id).First(ctx, &lbblog.ModelArticle{})
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}

	err = ArticleOrm.NewScope().Eq(lbblog.FieldId_, req.Article.Id).Update(ctx, utils.OrmStruct2Map(req.Article))
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}

	return &rsp, nil
}

func (a *LbblogServer) DelArticle(ctx context.Context, req *lbblog.DelArticleReq) (*lbblog.DelArticleRsp, error) {
	var rsp lbblog.DelArticleRsp

	err := ArticleOrm.NewScope().Eq(lbblog.FieldId_, req.Id).Delete(ctx)
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}

	return &rsp, nil
}

func (a *LbblogServer) AddArticle(ctx context.Context, req *lbblog.AddArticleReq) (*lbblog.AddArticleRsp, error) {
	var rsp lbblog.AddArticleRsp

	err := ArticleOrm.NewScope().Create(ctx, req.Article)
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}

	return &rsp, nil
}

func (a *LbblogServer) GetCategoryList(ctx context.Context, req *lbblog.GetCategoryListReq) (*lbblog.GetCategoryListRsp, error) {
	var rsp lbblog.GetCategoryListRsp

	db := CategoryOrm.NewList(req.ListOption)
	err := lbconst.NewListOptionProcessor(req.ListOption).
		Process()
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}

	rsp.Page, err = db.FindPage(ctx, &rsp.List)
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}

	return &rsp, nil
}

func (a *LbblogServer) GetCategory(ctx context.Context, req *lbblog.GetCategoryReq) (*lbblog.GetCategoryRsp, error) {
	var rsp lbblog.GetCategoryRsp

	err := CategoryOrm.NewScope().Eq(lbblog.FieldId_, req.Id).First(ctx, &rsp.Category)
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}

	return &rsp, nil
}

func (a *LbblogServer) UpdateCategory(ctx context.Context, req *lbblog.UpdateCategoryReq) (*lbblog.UpdateCategoryRsp, error) {
	var rsp lbblog.UpdateCategoryRsp

	err := CategoryOrm.NewScope().Eq(lbblog.FieldId_, req.Category.Id).First(ctx, &lbblog.ModelCategory{})
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}

	err = CategoryOrm.NewScope().Eq(lbblog.FieldId_, req.Category.Id).Update(ctx, utils.OrmStruct2Map(req.Category))
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}

	return &rsp, nil
}

func (a *LbblogServer) DelCategory(ctx context.Context, req *lbblog.DelCategoryReq) (*lbblog.DelCategoryRsp, error) {
	var rsp lbblog.DelCategoryRsp

	err := CategoryOrm.NewScope().Eq(lbblog.FieldId_, req.Id).Delete(ctx)
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}

	return &rsp, nil
}

func (a *LbblogServer) AddCategory(ctx context.Context, req *lbblog.AddCategoryReq) (*lbblog.AddCategoryRsp, error) {
	var rsp lbblog.AddCategoryRsp

	err := CategoryOrm.NewScope().Create(ctx, req.Category)
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}

	return &rsp, nil
}

func (a *LbblogServer) GetCommentList(ctx context.Context, req *lbblog.GetCommentListReq) (*lbblog.GetCommentListRsp, error) {
	var rsp lbblog.GetCommentListRsp

	db := CommentOrm.NewList(req.ListOption)
	err := lbconst.NewListOptionProcessor(req.ListOption).
		AddUint64(lbblog.GetCommentListReq_ListOptionArticleId, func(val uint64) error {
			db.Eq(lbblog.FieldArticleId_, val)
			return nil
		}).
		Process()
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}

	rsp.Page, err = db.OrderByDesc(lbblog.FieldCreatedAt_).FindPage(ctx, &rsp.List)
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}

	articleIdList := utils.PluckUint64List(rsp.List, lbblog.FieldArticleId)
	var articleList []*lbblog.ModelArticle
	err = ArticleOrm.NewScope().In(lbblog.FieldId_, articleIdList).Find(ctx, &articleList)
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}
	rsp.ArticleMap = utils.Slice2MapKeyByStructField(articleList, lbblog.FieldId).(map[uint64]*lbblog.ModelArticle)

	userIdList := utils.PluckUint64List(rsp.List, lbblog.FieldUserId)
	var userList []*lbuser.ModelUser
	err = UserOrm.NewScope().In(lbuser.FieldId_, userIdList).Find(ctx, &userList)
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}
	rsp.UserMap = utils.Slice2MapKeyByStructField(userList, lbuser.FieldId).(map[uint64]*lbuser.ModelUser)

	return &rsp, nil
}

func (a *LbblogServer) GetComment(ctx context.Context, req *lbblog.GetCommentReq) (*lbblog.GetCommentRsp, error) {
	var rsp lbblog.GetCommentRsp

	err := CommentOrm.NewScope().Eq(lbblog.FieldId_, req.Id).First(ctx, &rsp.Comment)
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}

	return &rsp, nil
}

func (a *LbblogServer) UpdateComment(ctx context.Context, req *lbblog.UpdateCommentReq) (*lbblog.UpdateCommentRsp, error) {
	var rsp lbblog.UpdateCommentRsp

	err := CommentOrm.NewScope().Eq(lbblog.FieldId_, req.Comment.Id).First(ctx, &lbblog.ModelComment{})
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}

	err = CommentOrm.NewScope().Eq(lbblog.FieldId_, req.Comment.Id).Update(ctx, utils.OrmStruct2Map(req.Comment))
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}

	return &rsp, nil
}

func (a *LbblogServer) DelComment(ctx context.Context, req *lbblog.DelCommentReq) (*lbblog.DelCommentRsp, error) {
	var rsp lbblog.DelCommentRsp

	err := CommentOrm.NewScope().Eq(lbblog.FieldId_, req.Id).Delete(ctx)
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}

	return &rsp, nil
}

func (a *LbblogServer) AddComment(ctx context.Context, req *lbblog.AddCommentReq) (*lbblog.AddCommentRsp, error) {
	var rsp lbblog.AddCommentRsp

	err := CommentOrm.NewScope().Create(ctx, req.Comment)
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}

	return &rsp, nil
}
