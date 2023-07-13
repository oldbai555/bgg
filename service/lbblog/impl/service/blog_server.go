package service

import (
	"context"
	"github.com/oldbai555/bgg/client/lbblog"
	"github.com/oldbai555/lbtool/log"
)

var BlogServer LbblogServer

type LbblogServer struct {
	*lbblog.UnimplementedLbblogServer
}

func (a *LbblogServer) GetArticleList(ctx context.Context, req *lbblog.GetArticleListReq) (*lbblog.GetArticleListRsp, error) {
	var rsp lbblog.GetArticleListRsp
	var err error
	log.Infof("hello world")
	return &rsp, err
}
func (a *LbblogServer) GetArticle(ctx context.Context, req *lbblog.GetArticleReq) (*lbblog.GetArticleRsp, error) {
	var rsp lbblog.GetArticleRsp
	var err error

	return &rsp, err
}
func (a *LbblogServer) UpdateArticle(ctx context.Context, req *lbblog.UpdateArticleReq) (*lbblog.UpdateArticleRsp, error) {
	var rsp lbblog.UpdateArticleRsp
	var err error

	return &rsp, err
}
func (a *LbblogServer) DelArticle(ctx context.Context, req *lbblog.DelArticleReq) (*lbblog.DelArticleRsp, error) {
	var rsp lbblog.DelArticleRsp
	var err error

	return &rsp, err
}
func (a *LbblogServer) AddArticle(ctx context.Context, req *lbblog.AddArticleReq) (*lbblog.AddArticleRsp, error) {
	var rsp lbblog.AddArticleRsp
	var err error

	return &rsp, err
}
func (a *LbblogServer) GetCategoryList(ctx context.Context, req *lbblog.GetCategoryListReq) (*lbblog.GetCategoryListRsp, error) {
	var rsp lbblog.GetCategoryListRsp
	var err error

	return &rsp, err
}
func (a *LbblogServer) GetCategory(ctx context.Context, req *lbblog.GetCategoryReq) (*lbblog.GetCategoryRsp, error) {
	var rsp lbblog.GetCategoryRsp
	var err error

	return &rsp, err
}
func (a *LbblogServer) UpdateCategory(ctx context.Context, req *lbblog.UpdateCategoryReq) (*lbblog.UpdateCategoryRsp, error) {
	var rsp lbblog.UpdateCategoryRsp
	var err error

	return &rsp, err
}
func (a *LbblogServer) DelCategory(ctx context.Context, req *lbblog.DelCategoryReq) (*lbblog.DelCategoryRsp, error) {
	var rsp lbblog.DelCategoryRsp
	var err error

	return &rsp, err
}
func (a *LbblogServer) AddCategory(ctx context.Context, req *lbblog.AddCategoryReq) (*lbblog.AddCategoryRsp, error) {
	var rsp lbblog.AddCategoryRsp
	var err error

	return &rsp, err
}
func (a *LbblogServer) GetCommentList(ctx context.Context, req *lbblog.GetCommentListReq) (*lbblog.GetCommentListRsp, error) {
	var rsp lbblog.GetCommentListRsp
	var err error

	return &rsp, err
}
func (a *LbblogServer) GetComment(ctx context.Context, req *lbblog.GetCommentReq) (*lbblog.GetCommentRsp, error) {
	var rsp lbblog.GetCommentRsp
	var err error

	return &rsp, err
}
func (a *LbblogServer) UpdateComment(ctx context.Context, req *lbblog.UpdateCommentReq) (*lbblog.UpdateCommentRsp, error) {
	var rsp lbblog.UpdateCommentRsp
	var err error

	return &rsp, err
}
func (a *LbblogServer) DelComment(ctx context.Context, req *lbblog.DelCommentReq) (*lbblog.DelCommentRsp, error) {
	var rsp lbblog.DelCommentRsp
	var err error

	return &rsp, err
}
func (a *LbblogServer) AddComment(ctx context.Context, req *lbblog.AddCommentReq) (*lbblog.AddCommentRsp, error) {
	var rsp lbblog.AddCommentRsp
	var err error

	return &rsp, err
}
