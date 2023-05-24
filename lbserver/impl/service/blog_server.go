package service

import (
	"context"
	"github.com/oldbai555/bgg/client/lbblog"
	"github.com/oldbai555/bgg/client/lbuser"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/utils"
)

var BlogServer LbblogServer

type LbblogServer struct {
	*lbblog.UnimplementedLbblogServer
}

func (a *LbblogServer) GetArticleList(ctx context.Context, req *lbblog.GetArticleListReq) (*lbblog.GetArticleListRsp, error) {
	var rsp lbblog.GetArticleListRsp
	var err error

	rsp.List, rsp.Paginate, err = Article.FindPaginate(ctx, req.Options)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	if len(rsp.List) == 0 {
		return &rsp, nil
	}

	categoryIdList := utils.PluckUint64List(rsp.List, lbblog.FieldCategoryId)
	categoryList, err := Category.GetByIdList(ctx, categoryIdList)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.CategoryMap = utils.Slice2MapKeyByStructField(categoryList, lbblog.FieldId).(map[uint64]*lbblog.ModelCategory)
	return &rsp, nil
}

func (a *LbblogServer) GetArticle(ctx context.Context, req *lbblog.GetArticleReq) (*lbblog.GetArticleRsp, error) {
	var rsp lbblog.GetArticleRsp
	var err error

	rsp.Article, err = Article.GetById(ctx, req.Id)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, nil
}

func (a *LbblogServer) UpdateArticle(ctx context.Context, req *lbblog.UpdateArticleReq) (*lbblog.UpdateArticleRsp, error) {
	var rsp lbblog.UpdateArticleRsp

	err := Article.UpdateById(ctx, req.Article.Id, utils.OrmStruct2Map(req.Article))
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, nil
}

func (a *LbblogServer) DelArticle(ctx context.Context, req *lbblog.DelArticleReq) (*lbblog.DelArticleRsp, error) {
	var rsp lbblog.DelArticleRsp

	err := Article.DeleteById(ctx, req.Id)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, nil
}

func (a *LbblogServer) AddArticle(ctx context.Context, req *lbblog.AddArticleReq) (*lbblog.AddArticleRsp, error) {
	var rsp lbblog.AddArticleRsp

	err := Article.Create(ctx, req.Article)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, nil
}

func (a *LbblogServer) GetCategoryList(ctx context.Context, req *lbblog.GetCategoryListReq) (*lbblog.GetCategoryListRsp, error) {
	var rsp lbblog.GetCategoryListRsp
	var err error

	rsp.List, rsp.Paginate, err = Category.FindPaginate(ctx, req.Options)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, nil
}

func (a *LbblogServer) GetCategory(ctx context.Context, req *lbblog.GetCategoryReq) (*lbblog.GetCategoryRsp, error) {
	var rsp lbblog.GetCategoryRsp
	var err error

	rsp.Category, err = Category.GetById(ctx, req.Id)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, nil
}

func (a *LbblogServer) UpdateCategory(ctx context.Context, req *lbblog.UpdateCategoryReq) (*lbblog.UpdateCategoryRsp, error) {
	var rsp lbblog.UpdateCategoryRsp

	err := Category.UpdateById(ctx, req.Category.Id, utils.OrmStruct2Map(req.Category))
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}

	return &rsp, nil
}

func (a *LbblogServer) DelCategory(ctx context.Context, req *lbblog.DelCategoryReq) (*lbblog.DelCategoryRsp, error) {
	var rsp lbblog.DelCategoryRsp

	err := Category.DeleteById(ctx, req.Id)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, nil
}

func (a *LbblogServer) AddCategory(ctx context.Context, req *lbblog.AddCategoryReq) (*lbblog.AddCategoryRsp, error) {
	var rsp lbblog.AddCategoryRsp

	err := Category.Create(ctx, req.Category)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	rsp.Category = req.Category
	return &rsp, nil
}

func (a *LbblogServer) GetCommentList(ctx context.Context, req *lbblog.GetCommentListReq) (*lbblog.GetCommentListRsp, error) {
	var rsp lbblog.GetCommentListRsp
	var err error

	rsp.List, rsp.Paginate, err = Comment.FindPaginate(ctx, req.Options)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	articleIdList := utils.PluckUint64List(rsp.List, lbblog.FieldArticleId)
	articleList, err := Article.GetByIdList(ctx, articleIdList)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.ArticleMap = utils.Slice2MapKeyByStructField(articleList, lbblog.FieldId).(map[uint64]*lbblog.ModelArticle)

	userIdList := utils.PluckUint64List(rsp.List, lbblog.FieldUserId)
	userList, err := User.GetByIdList(ctx, userIdList)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.UserMap = utils.Slice2MapKeyByStructField(userList, lbuser.FieldId).(map[uint64]*lbuser.ModelUser)

	return &rsp, nil
}

func (a *LbblogServer) GetComment(ctx context.Context, req *lbblog.GetCommentReq) (*lbblog.GetCommentRsp, error) {
	var rsp lbblog.GetCommentRsp
	var err error

	rsp.Comment, err = Comment.GetById(ctx, req.Id)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, nil
}

func (a *LbblogServer) UpdateComment(ctx context.Context, req *lbblog.UpdateCommentReq) (*lbblog.UpdateCommentRsp, error) {
	var rsp lbblog.UpdateCommentRsp

	err := Comment.UpdateById(ctx, req.Comment.Id, utils.OrmStruct2Map(req.Comment))
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, nil
}

func (a *LbblogServer) DelComment(ctx context.Context, req *lbblog.DelCommentReq) (*lbblog.DelCommentRsp, error) {
	var rsp lbblog.DelCommentRsp

	err := Comment.DeleteById(ctx, req.Id)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, nil
}

func (a *LbblogServer) AddComment(ctx context.Context, req *lbblog.AddCommentReq) (*lbblog.AddCommentRsp, error) {
	var rsp lbblog.AddCommentRsp

	err := Comment.Create(ctx, req.Comment)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, nil
}
