package service

import (
	"context"
	"github.com/oldbai555/bgg/service/lb"
	"github.com/oldbai555/bgg/service/lbblog"
	"github.com/oldbai555/bgg/service/lbblogserver/impl/dao/impl/mysql"
	"github.com/oldbai555/bgg/service/lbuser"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/utils"
)

var ServerImpl LbblogServer

type LbblogServer struct {
	*lbblog.UnimplementedLbblogServer
}

func (a *LbblogServer) GetArticleList(ctx context.Context, req *lbblog.GetArticleListReq) (*lbblog.GetArticleListRsp, error) {
	var rsp lbblog.GetArticleListRsp
	var err error

	rsp.List, rsp.Paginate, err = mysql.ArticleOrm.FindPaginate(ctx, req.Options)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	if len(rsp.List) == 0 {
		return &rsp, nil
	}

	categoryIdList := utils.PluckUint64List(rsp.List, lbblog.FieldCategoryId)
	categoryList, err := mysql.CategoryOrm.GetByIdList(ctx, categoryIdList)
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

	rsp.Article, err = mysql.ArticleOrm.GetById(ctx, req.Id)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, nil
}

func (a *LbblogServer) UpdateArticle(ctx context.Context, req *lbblog.UpdateArticleReq) (*lbblog.UpdateArticleRsp, error) {
	var rsp lbblog.UpdateArticleRsp

	_, err := mysql.ArticleOrm.UpdateById(ctx, req.Article.Id, utils.OrmStruct2Map(req.Article))
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, nil
}

func (a *LbblogServer) DelArticle(ctx context.Context, req *lbblog.DelArticleReq) (*lbblog.DelArticleRsp, error) {
	var rsp lbblog.DelArticleRsp

	_, err := mysql.ArticleOrm.DeleteById(ctx, req.Id)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, nil
}

func (a *LbblogServer) AddArticle(ctx context.Context, req *lbblog.AddArticleReq) (*lbblog.AddArticleRsp, error) {
	var rsp lbblog.AddArticleRsp

	_, err := mysql.ArticleOrm.Create(ctx, req.Article)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, nil
}

func (a *LbblogServer) GetCategoryList(ctx context.Context, req *lbblog.GetCategoryListReq) (*lbblog.GetCategoryListRsp, error) {
	var rsp lbblog.GetCategoryListRsp
	var err error

	rsp.List, rsp.Paginate, err = mysql.CategoryOrm.FindPaginate(ctx, req.Options)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, nil
}

func (a *LbblogServer) GetCategory(ctx context.Context, req *lbblog.GetCategoryReq) (*lbblog.GetCategoryRsp, error) {
	var rsp lbblog.GetCategoryRsp
	var err error

	rsp.Category, err = mysql.CategoryOrm.GetById(ctx, req.Id)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, nil
}

func (a *LbblogServer) UpdateCategory(ctx context.Context, req *lbblog.UpdateCategoryReq) (*lbblog.UpdateCategoryRsp, error) {
	var rsp lbblog.UpdateCategoryRsp

	_, err := mysql.CategoryOrm.UpdateById(ctx, req.Category.Id, utils.OrmStruct2Map(req.Category))
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}

	return &rsp, nil
}

func (a *LbblogServer) DelCategory(ctx context.Context, req *lbblog.DelCategoryReq) (*lbblog.DelCategoryRsp, error) {
	var rsp lbblog.DelCategoryRsp

	_, err := mysql.CategoryOrm.DeleteById(ctx, req.Id)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, nil
}

func (a *LbblogServer) AddCategory(ctx context.Context, req *lbblog.AddCategoryReq) (*lbblog.AddCategoryRsp, error) {
	var rsp lbblog.AddCategoryRsp

	_, err := mysql.CategoryOrm.Create(ctx, req.Category)
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

	rsp.List, rsp.Paginate, err = mysql.CommentOrm.FindPaginate(ctx, req.Options)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	articleIdList := utils.PluckUint64List(rsp.List, lbblog.FieldArticleId)
	articleList, err := mysql.ArticleOrm.GetByIdList(ctx, articleIdList)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.ArticleMap = utils.Slice2MapKeyByStructField(articleList, lbblog.FieldId).(map[uint64]*lbblog.ModelArticle)

	userIdList := utils.PluckUint64List(rsp.List, lbblog.FieldUserId)
	listRsp, err := lbuser.GetUserList(ctx, &lbuser.GetUserListReq{
		Options: lb.NewOptions().AddOpt(lb.DefaultListOption_DefaultListOptionIdList, userIdList).SetSkipTotal(),
	})
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	rsp.UserMap = utils.Slice2MapKeyByStructField(listRsp.List, lbuser.FieldId).(map[uint64]*lbuser.ModelUser)

	return &rsp, nil
}

func (a *LbblogServer) GetComment(ctx context.Context, req *lbblog.GetCommentReq) (*lbblog.GetCommentRsp, error) {
	var rsp lbblog.GetCommentRsp
	var err error

	rsp.Comment, err = mysql.CommentOrm.GetById(ctx, req.Id)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, nil
}

func (a *LbblogServer) UpdateComment(ctx context.Context, req *lbblog.UpdateCommentReq) (*lbblog.UpdateCommentRsp, error) {
	var rsp lbblog.UpdateCommentRsp

	_, err := mysql.CommentOrm.UpdateById(ctx, req.Comment.Id, utils.OrmStruct2Map(req.Comment))
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, nil
}

func (a *LbblogServer) DelComment(ctx context.Context, req *lbblog.DelCommentReq) (*lbblog.DelCommentRsp, error) {
	var rsp lbblog.DelCommentRsp

	_, err := mysql.CommentOrm.DeleteById(ctx, req.Id)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, nil
}

func (a *LbblogServer) AddComment(ctx context.Context, req *lbblog.AddCommentReq) (*lbblog.AddCommentRsp, error) {
	var rsp lbblog.AddCommentRsp

	_, err := mysql.CommentOrm.Create(ctx, req.Comment)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, nil
}
